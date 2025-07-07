package worker

import (
	"context"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "presets/admin/i18n"

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_pt_BR).
		RegisterForModules(language.SimplifiedChinese, MessagesKey, Messages_zh_CN).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

type Messages struct {
	StatusNew                string
	StatusScheduled          string
	StatusRunning            string
	StatusCancelled          string
	StatusDone               string
	StatusException          string
	StatusKilled             string
	FilterTabAll             string
	FilterTabRunning         string
	FilterTabScheduled       string
	FilterTabDone            string
	FilterTabErrors          string
	ActionCancelJob          string
	ActionAbortJob           string
	ActionUpdateJob          string
	ActionRerunJob           string
	DetailTitleStatus        string
	DetailTitleLog           string
	NoticeJobCannotBeAborted string
	NoticeJobWontBeExecuted  string
	ScheduleTime             string
	DateTimePickerClearText  string
	DateTimePickerOkText     string
	PleaseSelectJob          string
	QorJobs                  string
	QorJob                   string
	WorkersJob               string
	ErrJobRunsOnce           i18n.ErrorString
}

var Messages_en_US = &Messages{
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
	QorJobs:                  "Jobs",
	QorJob:                   "Job",
	WorkersJob:               "Job",
	ErrJobRunsOnce:           "This job runs once",
}

var Messages_zh_CN = &Messages{
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
	QorJobs:                  "Processos de Sistema",
	QorJob:                   "Processo de Sistema",
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
