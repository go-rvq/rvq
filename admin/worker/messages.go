package worker

import (
	"context"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "presets/admin/i18n"

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.SimplifiedChinese, MessagesKey, Messages_zh_CN).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

type Messages struct {
	ModuleDescription        string           `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	StatusNew                string           `i18n:"hint='Status of a job created and not yet scheduled or started.'"`
	StatusScheduled          string           `i18n:"hint='Status of a job waiting for its scheduled time.'"`
	StatusRunning            string           `i18n:"hint='Status of a job being executed.'"`
	StatusCancelled          string           `i18n:"hint='Status of a job cancelled before it ran.'"`
	StatusDone               string           `i18n:"hint='Status of a job that finished successfully.'"`
	StatusException          string           `i18n:"hint='Status of a job that stopped with an error.'"`
	StatusKilled             string           `i18n:"hint='Status of a job aborted while running.'"`
	FilterTabAll             string           `i18n:"label='Tab: all jobs', hint='Tab of the jobs listing that shows every job.'"`
	FilterTabRunning         string           `i18n:"label='Tab: running', hint='Tab of the jobs listing that shows the running jobs.'"`
	FilterTabScheduled       string           `i18n:"label='Tab: scheduled', hint='Tab of the jobs listing that shows the scheduled jobs.'"`
	FilterTabDone            string           `i18n:"label='Tab: done', hint='Tab of the jobs listing that shows the finished jobs.'"`
	FilterTabErrors          string           `i18n:"label='Tab: errors', hint='Tab of the jobs listing that shows the jobs that failed.'"`
	ActionCancelJob          string           `i18n:"hint='Button that cancels a job not started yet.'"`
	ActionAbortJob           string           `i18n:"hint='Button that aborts a running job.'"`
	ActionUpdateJob          string           `i18n:"hint='Button that saves the changes to a scheduled job.'"`
	ActionRerunJob           string           `i18n:"hint='Button that runs a finished job again.'"`
	DetailTitleStatus        string           `i18n:"label='Detail: status title', hint='Title of the status section of a job\\'s detail.'"`
	DetailTitleLog           string           `i18n:"label='Detail: log title', hint='Title of the log section of a job\\'s detail.'"`
	NoticeJobCannotBeAborted string           `i18n:"hint='Shown when a job changed status and can no longer be aborted, cancelled or updated.'"`
	NoticeJobWontBeExecuted  string           `i18n:"hint='Shown when the code of a job was removed or changed, so the job will not run.'"`
	ScheduleTime             string           `i18n:"hint='Label of the field with the date and time a job is scheduled to run.'"`
	DateTimePickerClearText  string           `i18n:"label='Date picker: clear', hint='Button of the date and time picker that clears the value.'"`
	DateTimePickerOkText     string           `i18n:"label='Date picker: OK', hint='Button of the date and time picker that confirms the value.'"`
	PleaseSelectJob          string           `i18n:"hint='Shown when an action needs a job and none was selected.'"`
	Jobs                     string           `i18n:"hint='Name of the jobs model in the plural (menu, listing title).'"`
	Job                      string           `i18n:"hint='Name of the jobs model in the singular (detail and form titles).'"`
	WorkersJob               string           `i18n:"label='Job (field)', hint='Label of the field that says which job a worker record runs.'"`
	ErrJobRunsOnce           i18n.ErrorString `i18n:"hint='Error shown when a job that runs only once is started again.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription:        "The jobs run in the background, and their progress.",
	StatusNew:                "New",
	StatusScheduled:          "Scheduled",
	StatusRunning:            "Running",
	StatusCancelled:          "Cancelled",
	StatusDone:               "Done",
	StatusException:          "Exception",
	StatusKilled:             "Killed",
	FilterTabAll:             "All Jobs",
	FilterTabRunning:         "Running",
	FilterTabScheduled:       "Scheduled",
	FilterTabDone:            "Done",
	FilterTabErrors:          "Errors",
	ActionCancelJob:          "Cancel Job",
	ActionAbortJob:           "Abort Job",
	ActionUpdateJob:          "Update Job",
	ActionRerunJob:           "Rerun Job",
	DetailTitleStatus:        "Status",
	DetailTitleLog:           "Log",
	NoticeJobCannotBeAborted: "This job cannot be aborted/canceled/updated due to its status change",
	NoticeJobWontBeExecuted:  "This job won't be executed due to code being deleted/modified",
	ScheduleTime:             "Schedule Time",
	DateTimePickerClearText:  "Clear",
	DateTimePickerOkText:     "OK",
	PleaseSelectJob:          "Please select job",
	Jobs:                     "Jobs",
	Job:                      "Job",
	WorkersJob:               "Job",
	ErrJobRunsOnce:           "This job runs once",
}

var Messages_zh_CN = &Messages{
	ModuleDescription:        "后台运行的任务及其进度。",
	StatusNew:                "新建",
	StatusScheduled:          "计划",
	StatusRunning:            "运行中",
	StatusCancelled:          "取消",
	StatusDone:               "完成",
	StatusException:          "错误",
	StatusKilled:             "中止",
	FilterTabAll:             "全部",
	FilterTabRunning:         "运行中",
	FilterTabScheduled:       "计划",
	FilterTabDone:            "完成",
	FilterTabErrors:          "错误",
	ActionCancelJob:          "取消Job",
	ActionAbortJob:           "中止Job",
	ActionUpdateJob:          "更新Job",
	ActionRerunJob:           "重跑Job",
	DetailTitleStatus:        "状态",
	DetailTitleLog:           "日志",
	NoticeJobCannotBeAborted: "Job状态已经改变，不能被中止/取消/更新",
	NoticeJobWontBeExecuted:  "Job代码被删除/修改, 这个Job不会被执行",
	ScheduleTime:             "执行时间",
	DateTimePickerClearText:  "清空",
	DateTimePickerOkText:     "确定",
	PleaseSelectJob:          "请选择Job",
	ErrJobRunsOnce:           "这个工作一次",
}

var Messages_pt_BR = &Messages{
	ModuleDescription:        "Os trabalhos executados em segundo plano, e seu andamento.",
	StatusNew:                "Novas",
	StatusScheduled:          "Agendada",
	StatusRunning:            "Executando",
	StatusCancelled:          "Cancelada",
	StatusDone:               "Concluída",
	StatusException:          "Exception",
	StatusKilled:             "Morta",
	FilterTabAll:             "Todas",
	FilterTabRunning:         "Executando",
	FilterTabScheduled:       "Agendadas",
	FilterTabDone:            "Encerradas",
	FilterTabErrors:          "Com Erros",
	ActionCancelJob:          "Cancelar Tarefa",
	ActionAbortJob:           "Abortar Tarefa",
	ActionUpdateJob:          "Atualizar Tarefa",
	ActionRerunJob:           "Executar Novamente",
	DetailTitleStatus:        "Situação",
	DetailTitleLog:           "Registro",
	NoticeJobCannotBeAborted: "Esta tarefa não pode ser abortada/cancelada/atualizada devido à mudança de status",
	NoticeJobWontBeExecuted:  "Esta tarefa não será executada porque o código foi excluído/modificado",
	ScheduleTime:             "Horário de Agendamento",
	DateTimePickerClearText:  "Limpar",
	DateTimePickerOkText:     "OK",
	PleaseSelectJob:          "Por favor selecione uma tarefa",
	Jobs:                     "Processos de Sistema",
	Job:                      "Processo de Sistema",
	WorkersJob:               "Tarefa",
	ErrJobRunsOnce:           "Esta tarefa só pode ser executada uma única vez",
}

func (m *Messages) GetStatus(status string) string {
	switch status {
	case JobStatusNew:
		return m.StatusNew
	case JobStatusScheduled:
		return m.StatusScheduled
	case JobStatusRunning:
		return m.StatusRunning
	case JobStatusCancelled:
		return m.StatusCancelled
	case JobStatusDone:
		return m.StatusDone
	case JobStatusException:
		return m.StatusException
	case JobStatusKilled:
		return m.StatusKilled
	}
	return status
}

func getTJob(ctx context.Context, v string) string {
	return i18n.PT(ctx, presets.ModelsI18nModuleKey, "WorkerJob", v)
}
