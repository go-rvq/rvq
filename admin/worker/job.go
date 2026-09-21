package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	. "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"gorm.io/gorm"
)

//go:generate moq -pkg mock -out mock/job.go . JobInterface

type JobCronConfig struct {
	Arg  interface{} `json:"arg"`
	Spec string      `json:"spec"`
	Once bool        `json:"once"`
}

func (c *JobCronConfig) Valid() bool {
	return len(c.Spec) > 0
}

type JobBuilder struct {
	b              *Builder
	name           string
	r              interface{}
	rmb            *presets.ModelBuilder
	h              JobHandler
	contextHandler func(*web.EventContext) map[string]interface{} // optional
	global         bool
	system         bool
	title          func(ctx *web.EventContext) string
	cronConfig     JobCronConfig
}

func newJob(b *Builder, name string) *JobBuilder {
	if b == nil {
		panic("builder is nil")
	}
	if strings.TrimSpace(name) == "" {
		panic("name is empty")
	}

	return &JobBuilder{
		b:      b,
		name:   name,
		global: true,
	}
}

func (jb *JobBuilder) Title(f func(ctx *web.EventContext) string) *JobBuilder {
	jb.title = f
	return jb
}

func (jb *JobBuilder) GetTitle(ctx *web.EventContext) (t string) {
	if jb.title != nil {
		t = jb.title(ctx)
	} else {
		t = getTJob(ctx.Context(), jb.name)
	}
	return
}

func (jb *JobBuilder) System(v bool) *JobBuilder {
	jb.system = v
	return jb
}

func (jb *JobBuilder) IsSystem() bool {
	return jb.system
}

func (jb *JobBuilder) CronConfig(config JobCronConfig) *JobBuilder {
	jb.cronConfig = config
	return jb
}

func (jb *JobBuilder) GetCronConfig() JobCronConfig {
	return jb.cronConfig
}

type JobHandler func(context.Context, JobInterface) error

// r should be ptr to struct
func (jb *JobBuilder) Resource(r interface{}, do ...func(mb *presets.ModelBuilder)) *JobBuilder {
	{
		v := reflect.TypeOf(r)
		if v.Kind() != reflect.Ptr {
			panic("resource is not ptr to struct")
		}
		if v.Elem().Kind() != reflect.Struct {
			panic("resource is not ptr to struct")
		}
	}

	jb.r = r
	// jpb only renders the job form; its models are not menu entries, and two
	// job resources may well share a type name.
	jb.rmb = jb.b.jpb.Model(r, presets.ModelNotInMenu())

	if _, ok := r.(Scheduler); ok {
		jb.rmb.Editing().Field("ScheduleTime").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) HTMLComponent {
			msgr := i18n.MustGetModuleMessages(ctx.Context(), MessagesKey, Messages_en_US).(*Messages)
			t := field.Obj.(Scheduler).GetScheduleTime()
			var v string
			if t != nil {
				v = t.Local().Format("2006-01-02 15:04")
			}
			return vx.VXDateTimePicker().Attr(web.VField(field.FormKey, v)...).Label(msgr.ScheduleTime).
				TimePickerProps(vx.TimePickerProps{
					Format:     "24hr",
					Scrollable: true,
				}).
				ClearText(msgr.DateTimePickerClearText).OkText(msgr.DateTimePickerOkText)
		}).SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
			v := ctx.R.Form.Get(field.Name)
			if v == "" {
				return nil
			}
			t, err := time.ParseInLocation("2006-01-02 15:04", v, time.Local)
			if err != nil {
				return err
			}
			obj.(Scheduler).SetScheduleTime(&t)
			return nil
		})
	}

	for _, f := range do {
		f(jb.rmb)
	}
	return jb
}

func (jb *JobBuilder) GetResourceBuilder() *presets.ModelBuilder {
	return jb.rmb
}

func (jb *JobBuilder) Handler(h JobHandler) *JobBuilder {
	jb.h = h
	return jb
}

func (jb *JobBuilder) ContextHandler(handler func(*web.EventContext) map[string]interface{}) *JobBuilder {
	jb.contextHandler = handler
	return jb
}

func (jb *JobBuilder) newResourceObject() interface{} {
	if jb.r == nil {
		return nil
	}
	return reflect.New(reflect.TypeOf(jb.r).Elem()).Interface()
}

func (jb *JobBuilder) unmarshalForm(ctx *web.EventContext) (args interface{}, vErr web.ValidationErrors) {
	args = jb.newResourceObject()
	if args != nil {
		vErr = jb.rmb.Editing().RunSetterFunc(nil, ctx, false, args)
	}

	return args, vErr
}

func (jb *JobBuilder) parseArgs(in string) (args interface{}, err error) {
	if jb.r == nil {
		return nil, nil
	}
	args = jb.newResourceObject()
	err = json.Unmarshal([]byte(in), args)
	if err != nil {
		return nil, err
	}

	return args, nil
}

func getModelJobInstance(db *gorm.DB, jobID uint) (*JobInstance, error) {
	var insts []*JobInstance
	err := db.Where("job_id = ?", jobID).
		Order("created_at desc").
		Limit(1).
		Find(&insts).
		Error
	if err != nil {
		return nil, err
	}
	if len(insts) == 0 {
		return nil, errors.New("no qor job instance")
	}

	return insts[0], nil
}

func (jb *JobBuilder) getJobInstance(jobID uint) (*JobInstance, error) {
	inst, err := getModelJobInstance(jb.b.db, jobID)
	if err != nil {
		return nil, err
	}

	inst.jb = jb

	return inst, nil
}

func (jb *JobBuilder) newJobInstance(
	r *http.Request,
	jobID uint,
	jobName string,
	once bool,
	args interface{},
	context interface{},
) (*JobInstance, error) {
	var mArgs string
	if v, ok := args.(string); ok {
		mArgs = v
	} else {
		bArgs, err := json.Marshal(args)
		if err != nil {
			return nil, err
		}
		mArgs = string(bArgs)
	}

	var ctx string
	if v, ok := context.(string); ok {
		ctx = v
	} else {
		bArgs, err := json.Marshal(context)
		if err != nil {
			return nil, err
		}
		ctx = string(bArgs)
	}

	inst := JobInstance{
		JobID:   jobID,
		Args:    mArgs,
		Context: ctx,
		Job:     jobName,
		Status:  JobStatusNew,
		Once:    once,
	}
	if jb.b.getCurrentUserIDFunc != nil {
		inst.Operator = jb.b.getCurrentUserIDFunc(r)
	}
	err := jb.b.db.Create(&inst).Error
	if err != nil {
		return nil, err
	}

	return jb.getJobInstance(jobID)
}

type QueJobInterface interface {
	JobInterface

	GetStatus() string
	FetchAndSetStatus() (string, error)
	SetStatus(string) error

	StartRefresh()
	StopRefresh()

	GetHandler() JobHandler
}

type JobInfo struct {
	JobID    string
	JobName  string
	Operator string
	Argument interface{}
	Context  map[string]interface{}
}

// for job handler
type JobInterface interface {
	GetJobInfo() (*JobInfo, error)
	SetProgress(uint) error
	SetProgressText(string) error
	AddLog(string) error
	AddLogf(format string, a ...interface{}) error
}

var _ QueJobInterface = (*JobInstance)(nil)

func (job *JobInstance) GetJobInfo() (ji *JobInfo, err error) {
	arg, err := job.getArgument()
	if err != nil {
		return
	}

	context, err := job.getContext()
	if err != nil {
		return
	}

	return &JobInfo{
		JobID:    fmt.Sprint(job.JobID),
		JobName:  job.Job,
		Operator: job.Operator,
		Argument: arg,
		Context:  context,
	}, nil
}

func (job *JobInstance) GetStatus() string {
	return job.Status
}

func (job *JobInstance) FetchAndSetStatus() (string, error) {
	var status string
	{
		db, err := job.jb.b.db.DB()
		if err != nil {
			return job.Status, err
		}

		err = db.QueryRow("select status from job_instances where id = $1", job.ID).Scan(&status)
		if err != nil {
			return job.Status, err
		}
		if status == "" {
			return job.Status, errors.New("failed to fetch job_instance status")
		}
	}

	if job.Status != status {
		err := job.SetStatus(status)
		if err != nil {
			return job.Status, err
		}
	}

	return job.Status, nil
}

func (job *JobInstance) SetStatus(status string) error {
	job.mutex.Lock()
	defer job.mutex.Unlock()

	job.Status = status
	if status == JobStatusDone {
		job.Progress = 100
	}

	if job.shouldCallSave() {
		return job.callSave()
	}

	return nil
}

func (job *JobInstance) SetProgress(progress uint) error {
	job.mutex.Lock()
	defer job.mutex.Unlock()

	if progress > 100 {
		progress = 100
	}
	job.Progress = progress

	if job.shouldCallSave() {
		return job.callSave()
	}

	return nil
}

func (job *JobInstance) SetProgressText(s string) error {
	job.mutex.Lock()
	defer job.mutex.Unlock()

	job.ProgressText = s
	if job.shouldCallSave() {
		return job.callSave()
	}

	return nil
}

func (job *JobInstance) AddLog(log string) error {
	if err := job.jb.b.db.Create(&JobLog{
		JobInstanceID: job.ID,
		Log:           log,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (job *JobInstance) AddLogf(format string, a ...interface{}) error {
	return job.AddLog(fmt.Sprintf(format, a...))
}

func (job *JobInstance) StartRefresh() {
	job.mutex.Lock()
	defer job.mutex.Unlock()
	if !job.inRefresh {
		job.inRefresh = true
		job.stopRefresh = false

		go func() {
			job.refresh()
		}()
	}
}

func (job *JobInstance) StopRefresh() {
	job.mutex.Lock()
	defer job.mutex.Unlock()

	err := job.callSave()
	if err != nil {
		log.Println(err)
	}

	job.stopRefresh = true
}

func (job *JobInstance) GetHandler() JobHandler {
	return job.jb.h
}

func (job *JobInstance) getArgument() (interface{}, error) {
	return job.jb.parseArgs(job.Args)
}

func (job *JobInstance) getContext() (map[string]interface{}, error) {
	context := make(map[string]interface{})
	err := json.Unmarshal([]byte(job.Context), &context)
	return context, err
}

func (job *JobInstance) shouldCallSave() bool {
	return !job.inRefresh || job.stopRefresh
}

func (job *JobInstance) callSave() error {
	err := job.jb.b.setStatus(job.JobID, job.Status)
	if err != nil {
		return err
	}
	return job.jb.b.db.Save(job).Error
}

func (job *JobInstance) refresh() {
	job.mutex.Lock()
	defer job.mutex.Unlock()

	err := job.callSave()
	if err != nil {
		log.Println(err)
	}

	if job.stopRefresh {
		job.inRefresh = false
		job.stopRefresh = false
	} else {
		time.AfterFunc(5*time.Second, job.refresh)
	}
}
