package worker

import (
	"sync"
	"time"

	"github.com/go-rvq/rvq/admin/utils/uuidkey"
	"github.com/go-rvq/rvq/web"
	"github.com/google/uuid"
)

type Job struct {
	uuidkey.Model

	Job    string
	Status string      `sql:"default:'new'"`
	Once   bool        `sql:"not null;default:false"`
	Args   interface{} `sql:"-" gorm:"-"`
}

type JobInstance struct {
	uuidkey.Model

	JobID uuid.UUID `gorm:"type:uuid;index"`

	Operator string

	Job     string
	Status  string `sql:"default:'new'"`
	Args    string
	Context string

	Progress     uint
	ProgressText string

	jb          *JobBuilder `sql:"-"`
	mutex       sync.Mutex  `sql:"-"`
	stopRefresh bool        `sql:"-"`
	inRefresh   bool        `sql:"-"`

	Once bool `sql:"not null;default:false"`
}

type JobLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `gorm:"index"`

	JobInstanceID uuid.UUID `gorm:"type:uuid;index"`
	Log           string
}

type Scheduler interface {
	GetScheduleTime() *time.Time
	SetScheduleTime(t *time.Time)
}

// Schedule could be embedded as job argument, then the job will get run as scheduled feature
type Schedule struct {
	ScheduleTime *time.Time
}

// GetScheduleTime get scheduled time
func (schedule *Schedule) GetScheduleTime() *time.Time {
	if scheduleTime := schedule.ScheduleTime; scheduleTime != nil {
		if scheduleTime.After(time.Now().Add(time.Minute)) {
			return scheduleTime
		}
	}
	return nil
}

func (schedule *Schedule) SetScheduleTime(t *time.Time) {
	schedule.ScheduleTime = t
}

type GoQueError struct {
	uuidkey.Model
	Error string
}

type CronScheduler interface {
	GetCronRule() string
	SetCronRule(v string)
}

type CronSchedule struct {
	CronRule string
}

func (s *CronSchedule) SetCronRule(v string) {
	s.CronRule = v
}

func (s *CronSchedule) GetCronRule() string {
	return s.CronRule
}

// paramJobID is the job id a request names ("jobID"); uuid.Nil when it names
// none, or not a UUID.
func paramJobID(ctx *web.EventContext) uuid.UUID {
	id, _ := uuid.Parse(ctx.Param("jobID"))
	return id
}
