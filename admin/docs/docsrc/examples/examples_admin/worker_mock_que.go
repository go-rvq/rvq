package examples_admin

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/admin/worker"
	"github.com/go-rvq/rvq/admin/worker/mock"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func WorkerExample(b *presets.Builder, db *gorm.DB) http.Handler {
	b.DataOperator(gorm2op.DataOperator(db))

	wb := worker.New(i18n.New(), db, worker.WithQueue(Que))
	b.Use(wb)
	addJobs(wb)
	wb.Listen()
	return b
}

var Que = &mock.QueueMock{
	AddFunc: func(ctx context.Context, job worker.QueJobInterface) error {
		jobInfo, err := job.GetJobInfo()
		if err != nil {
			return err
		}
		if scheduler, ok := jobInfo.Argument.(worker.Scheduler); ok && scheduler.GetScheduleTime() != nil {
			job.SetStatus(worker.JobStatusScheduled)
			go func() {
				time.Sleep(time.Until(*scheduler.GetScheduleTime()))
				ConsumeQueItem(job)
			}()
		} else {
			go func() {
				ConsumeQueItem(job)
			}()
		}
		return nil
	},
	KillFunc: func(ctx context.Context, job worker.QueJobInterface) error {
		return job.SetStatus(worker.JobStatusKilled)
	},
	ListenFunc: func(jobDefs []*worker.JobDefinition, getJob func(jobID uuid.UUID) (worker.QueJobInterface, error)) error {
		return nil
	},
	RemoveFunc: func(ctx context.Context, job worker.QueJobInterface) error {
		return job.SetStatus(worker.JobStatusCancelled)
	},
}

func ConsumeQueItem(job worker.QueJobInterface) (err error) {
	defer func() {
		if r := recover(); r != nil {
			job.AddLog(string(debug.Stack()))
			job.SetProgressText(fmt.Sprint(r))
			job.SetStatus(worker.JobStatusException)
			job.StopRefresh()
		}
	}()

	if job.GetStatus() == worker.JobStatusCancelled {
		return
	}
	if job.GetStatus() != worker.JobStatusNew && job.GetStatus() != worker.JobStatusScheduled {
		job.SetStatus(worker.JobStatusKilled)
		return errors.New("invalid job status, current status: " + job.GetStatus())
	}

	err = job.SetStatus(worker.JobStatusRunning)
	if err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond)

	hctx, cf := context.WithCancel(context.Background())
	hDoneC := make(chan struct{})
	isAborted := false
	go func() {
		timer := time.NewTicker(time.Second)
		for {
			select {
			case <-hDoneC:
				return
			case <-timer.C:
				status, _ := job.FetchAndSetStatus()
				if status == worker.JobStatusKilled {
					isAborted = true
					cf()
					return
				}
			}
		}
	}()
	job.StartRefresh()
	err = job.GetHandler()(hctx, job)
	job.StopRefresh()
	if !isAborted {
		hDoneC <- struct{}{}
	}
	if err != nil {
		job.SetProgressText(err.Error())
		job.SetStatus(worker.JobStatusException)
		return err
	}
	if isAborted {
		return
	}

	err = job.SetStatus(worker.JobStatusDone)
	if err != nil {
		return err
	}
	return
}
