package publish

import (
	"context"
	"strings"

	"github.com/go-rvq/rvq/x/i18n"
)

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, I18nPublishKey, DefaultMessages).(*Messages)
}

type Messages struct {
	ModuleDescription                       string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	StatusDraft                             string `i18n:"label='Status: draft', hint='Status of a record not published yet.'"`
	StatusOnline                            string `i18n:"label='Status: online', hint='Status of a published record.'"`
	StatusOffline                           string `i18n:"label='Status: offline', hint='Status of a record taken off the site.'"`
	Publication                             string `i18n:"hint='Title of the publication section of a record.'"`
	Publish                                 string `i18n:"hint='Action that publishes a record.'"`
	PublishHelp                             string `i18n:"label='Publish: hint', hint='When Publish applies: to a draft or offline record.'"`
	PublishOrRepublish                      string `i18n:"label='Publish or republish', hint='Action that publishes a draft or offline record, or republishes an online one.'"`
	PublishOrRepublishHelp                  string `i18n:"label='Publish or republish: hint', hint='What Publish or republish does depending on the status.'"`
	Unpublish                               string `i18n:"hint='Action that takes a record off the site.'"`
	UnpublishHelp                           string `i18n:"label='Unpublish: hint', hint='When Unpublish applies: to an online record.'"`
	Republish                               string `i18n:"hint='Action that publishes again an online record.'"`
	RepublishHelp                           string `i18n:"label='Republish: hint', hint='When Republish applies: to an online record.'"`
	Areyousure                              string `i18n:"label='Are you sure', hint='Asks to confirm an action.'"`
	ScheduledStartAt                        string `i18n:"label='Scheduled start', hint='Label of when a scheduled publication starts.'"`
	ScheduledEndAt                          string `i18n:"label='Scheduled end', hint='Label of when a scheduled publication ends.'"`
	ScheduledStartAtShouldLaterThanNow      string `i18n:"label='Start in the past', hint='Error when the scheduled start is not in the future.'"`
	ScheduledEndAtShouldLaterThanNowOrEmpty string `i18n:"label='End in the past', hint='Error when the scheduled end is set and not in the future.'"`
	ScheduledEndAtShouldLaterThanStartAt    string `i18n:"label='End before start', hint='Error when the scheduled end is not after the start.'"`
	ScheduledStartAtShouldNotEmpty          string `i18n:"label='Start missing', hint='Error when the scheduled start is empty.'"`
	PublishedAt                             string `i18n:"label='Published at', hint='Label of when a record was published.'"`
	UnPublishedAt                           string `i18n:"label='Unpublished at', hint='Label of when a record was taken off the site.'"`
	ActualPublishTime                       string `i18n:"hint='Title of when a record actually was published.'"`
	SchedulePublishTime                     string `i18n:"hint='Title of when a record is scheduled to be published.'"`
	NotSet                                  string `i18n:"hint='Shown for a time not set.'"`
	WhenDoYouWantToPublish                  string `i18n:"label='When to publish', hint='Title of the dialog that schedules a publication.'"`
	PublishScheduleTip                      string `i18n:"label='Schedule: hint', hint='How scheduling works.', fields=(;'{SchedulePublishTime}'='the title of the scheduled time')"`
	DateTimePickerClearText                 string `i18n:"label='Date picker: clear', hint='Button of the date and time picker that clears the value.'"`
	DateTimePickerOkText                    string `i18n:"label='Date picker: OK', hint='Button of the date and time picker that confirms the value.'"`
	SaveAsNewVersion                        string `i18n:"hint='Action that saves the changes as a new version.'"`
	SwitchedToNewVersion                    string `i18n:"hint='Shown once the record switched to the new version.'"`
	SuccessfullyCreated                     string `i18n:"hint='Shown once a version was created.'"`
	SuccessfullyRename                      string `i18n:"label='Successfully renamed', hint='Shown once a version was renamed.'"`
	SuccessfullyPublished                   string `i18n:"hint='Shown once a record was published.'"`
	SuccessfullyPublishedOrRepublished      string `i18n:"hint='Shown once records were published or republished.'"`
	SuccessfullyUnpublished                 string `i18n:"hint='Shown once a record was taken off the site.'"`
	SuccessfullyRepublished                 string `i18n:"hint='Shown once a record was republished.'"`
	OnlineVersion                           string `i18n:"hint='Mark of the version that is published.'"`
	VersionsList                            string `i18n:"hint='Title of the list of a record\\'s versions.'"`
	AllVersions                             string `i18n:"hint='Filter of the versions list: every version.'"`
	NamedVersions                           string `i18n:"hint='Filter of the versions list: the versions given a name.'"`
	RenameVersion                           string `i18n:"hint='Action that renames a version.'"`
	DeleteVersionConfirmationTextTemplate   string `i18n:"label='Delete version: confirmation', hint='Asks to confirm the deletion of a version.', fields=(;'{VersionName}'='its name')"`
	BulkActionConfirmationTextTemplate      string `i18n:"type=html, label='Bulk action: confirmation', hint='Asks to confirm an action on the records listed (HTML).', fields=(;'{Action}'='the action')"`
	BulkActionNoRecordsTextTemplate         string `i18n:"type=html, label='Bulk action: no records', hint='Shown when there is no record to act on (HTML).', fields=(;'{Action}'='the action')"`

	FilterTabAllVersions   string `i18n:"label='Tab: all versions', hint='Tab that shows every version.'"`
	FilterTabOnlineVersion string `i18n:"label='Tab: online versions', hint='Tab that shows the published versions.'"`
	FilterTabNamedVersions string `i18n:"label='Tab: named versions', hint='Tab that shows the versions given a name.'"`
	Rename                 string `i18n:"hint='Action that renames.'"`
	PageOverView           string `i18n:"label='Page overview', hint='Title of the summary of a page.'"`
	Duplicate              string `i18n:"hint='Action that copies a record into a new one.'"`
}

func (msgr *Messages) Status(status string) string {
	switch status {
	case "draft":
		return msgr.StatusDraft
	case "online":
		return msgr.StatusOnline
	case "offline":
		return msgr.StatusOffline
	default:
		return status
	}
}

func (msgr *Messages) ActivityTitle(activity string) string {
	switch activity {
	case ActivityPublish:
		return msgr.Publish
	case ActivityPublishOrRepublish:
		return msgr.PublishOrRepublish
	case ActivityRepublish:
		return msgr.Republish
	case ActivityUnPublish:
		return msgr.Unpublish
	default:
		return ""
	}
}

func (msgr *Messages) ActivityHelp(activity string) string {
	switch activity {
	case ActivityPublish:
		return msgr.PublishHelp
	case ActivityPublishOrRepublish:
		return msgr.PublishOrRepublishHelp
	case ActivityRepublish:
		return msgr.RepublishHelp
	case ActivityUnPublish:
		return msgr.UnpublishHelp
	default:
		return ""
	}
}

func (msgr *Messages) DeleteVersionConfirmationText(versionName string) string {
	return strings.NewReplacer("{VersionName}", versionName).
		Replace(msgr.DeleteVersionConfirmationTextTemplate)
}

func (msgr *Messages) BulkActionConfirmationText(action string) string {
	return strings.NewReplacer("{Action}", action).
		Replace(msgr.BulkActionConfirmationTextTemplate)
}

func (msgr *Messages) BulkActionNoRecordsText(action string) string {
	return strings.NewReplacer("{Action}", action).
		Replace(msgr.BulkActionNoRecordsTextTemplate)
}

var Messages_en_US = &Messages{
	ModuleDescription:                       "Publishing: versions, schedules and putting records online.",
	StatusDraft:                             "Draft",
	StatusOnline:                            "Online",
	StatusOffline:                           "Offline",
	Publication:                             "Publication",
	Publish:                                 "Publish",
	PublishHelp:                             "Publish only if status is 'draft' or 'offline'",
	PublishOrRepublish:                      "Publish/Republish",
	PublishOrRepublishHelp:                  "If status is 'draft' or 'offline', publish then, other else 'republish'",
	Unpublish:                               "UnPublish",
	UnpublishHelp:                           "Unpublish if status is 'online'",
	Republish:                               "Republish",
	RepublishHelp:                           "Republish if status is 'online'",
	Areyousure:                              "Are you sure?",
	ScheduledStartAt:                        "Start at",
	ScheduledEndAt:                          "End at",
	ScheduledStartAtShouldLaterThanNow:      "Start at should be later than now",
	ScheduledEndAtShouldLaterThanNowOrEmpty: "End at should be later than now or empty",
	ScheduledEndAtShouldLaterThanStartAt:    "End at should be later than start at",
	ScheduledStartAtShouldNotEmpty:          "Start at should not be empty",
	PublishedAt:                             "Start at",
	UnPublishedAt:                           "End at",
	ActualPublishTime:                       "Actual Publish Time",
	SchedulePublishTime:                     "Schedule Publish Time",
	NotSet:                                  "Not set",
	WhenDoYouWantToPublish:                  "When do you want to publish?",
	PublishScheduleTip:                      "After you set the {SchedulePublishTime}, the system will automatically publish/unpublish it.",
	DateTimePickerClearText:                 "Clear",
	DateTimePickerOkText:                    "OK",
	SaveAsNewVersion:                        "Save As New Version",
	SwitchedToNewVersion:                    "Switched To New Version",
	SuccessfullyCreated:                     "Successfully Created",
	SuccessfullyRename:                      "Successfully Rename",
	SuccessfullyPublished:                   "Successfully published",
	SuccessfullyPublishedOrRepublished:      "Successfully published or republished",
	SuccessfullyUnpublished:                 "Successfully unpublished",
	SuccessfullyRepublished:                 "Successfully republished",
	OnlineVersion:                           "Online Version",
	VersionsList:                            "Versions List",
	AllVersions:                             "All versions",
	NamedVersions:                           "Named versions",
	RenameVersion:                           "Rename Version",
	DeleteVersionConfirmationTextTemplate:   "Are you sure you want to delete version {VersionName} ?",
	BulkActionConfirmationTextTemplate:      "Are you sure you want to <b>{Action}</b> below records?",
	BulkActionNoRecordsTextTemplate:         "No records to <b>{Action}</b>.",
	FilterTabAllVersions:                    "All Versions",
	FilterTabOnlineVersion:                  "Online Versions",
	FilterTabNamedVersions:                  "Named Versions",
	Rename:                                  "Rename",
	PageOverView:                            "Page Overview",
	Duplicate:                               "Duplicate",
}

var DefaultMessages = Messages_en_US

var Messages_zh_CN = &Messages{
	ModuleDescription:                       "发布：版本、计划和上线记录。",
	StatusDraft:                             "草稿",
	StatusOnline:                            "在线",
	StatusOffline:                           "离线",
	Publish:                                 "发布",
	Unpublish:                               "取消发布",
	Republish:                               "重新发布",
	Areyousure:                              "你确定吗?",
	ScheduledStartAt:                        "发布时间",
	ScheduledEndAt:                          "下线时间",
	ScheduledStartAtShouldLaterThanNow:      "计划发布时间应当晚于现在时间",
	ScheduledEndAtShouldLaterThanNowOrEmpty: "计划下线时间应当晚于现在时间",
	ScheduledEndAtShouldLaterThanStartAt:    "计划下线时间应当晚于计划发布时间",
	ScheduledStartAtShouldNotEmpty:          "计划发布时间不得为空",
	PublishedAt:                             "发布时间",
	UnPublishedAt:                           "下线时间",
	ActualPublishTime:                       "实际发布时间",
	SchedulePublishTime:                     "计划发布时间",
	NotSet:                                  "未设定",
	WhenDoYouWantToPublish:                  "你希望什么时候发布？",
	PublishScheduleTip:                      "设定好 {SchedulePublishTime} 之后, 系统会按照时间自动将它发布/下线。",
	DateTimePickerClearText:                 "清空",
	DateTimePickerOkText:                    "确定",
	SaveAsNewVersion:                        "保存为一个新版本",
	SwitchedToNewVersion:                    "切换到新版本",
	SuccessfullyCreated:                     "成功创建",
	SuccessfullyRename:                      "成功命名",
	OnlineVersion:                           "在线版本",
	VersionsList:                            "版本列表",
	AllVersions:                             "所有版本",
	NamedVersions:                           "已命名版本",
	RenameVersion:                           "命名版本",
	DeleteVersionConfirmationTextTemplate:   "你确定你要删除此版本 {VersionName} 吗？",

	FilterTabAllVersions:   "所有版本",
	FilterTabOnlineVersion: "在线版本",
	FilterTabNamedVersions: "已命名版本",
	Rename:                 "重命名",
	PageOverView:           "页面概览",
	Duplicate:              "复制",
}

var Messages_ja_JP = &Messages{
	ModuleDescription:                       "公開：バージョン、スケジュール、レコードの公開。",
	StatusDraft:                             "下書き",
	StatusOnline:                            "公開中",
	StatusOffline:                           "非公開中",
	Publish:                                 "公開する",
	Unpublish:                               "非公開",
	Republish:                               "再公開",
	Areyousure:                              "よろしいですか？",
	ScheduledStartAt:                        "公開開始日時",
	ScheduledEndAt:                          "公開終了日時",
	ScheduledStartAtShouldLaterThanNow:      "予定開始時間は現在時刻よりも後でなければなりません",
	ScheduledEndAtShouldLaterThanNowOrEmpty: "予定終了時間は現在時刻よりも後、または空でなければなりません",
	ScheduledEndAtShouldLaterThanStartAt:    "予定終了時間は予定開始時間よりも後でなければなりません",
	ScheduledStartAtShouldNotEmpty:          "予定開始時間は空ではない必要があります",
	PublishedAt:                             "開始日時",
	UnPublishedAt:                           "公開終了日時",
	ActualPublishTime:                       "投稿日時",
	SchedulePublishTime:                     "公開日時を設定する",
	NotSet:                                  "未セット",
	WhenDoYouWantToPublish:                  "公開日時を設定してください",
	PublishScheduleTip:                      "{SchedulePublishTime} 設定後、システムが自動で当該記事の公開・非公開を行います。",
	DateTimePickerClearText:                 "クリア",
	DateTimePickerOkText:                    "OK",
	SaveAsNewVersion:                        "新規バージョンとして保存する",
	SwitchedToNewVersion:                    "新規バージョンに変更する",
	SuccessfullyCreated:                     "作成に成功しました",
	SuccessfullyRename:                      "名付けに成功しました",
	OnlineVersion:                           "オンラインバージョン",
	VersionsList:                            "バージョンリスト",
	AllVersions:                             "全てのバージョン",
	NamedVersions:                           "名付け済みバージョン",
	RenameVersion:                           "バージョンの名前を変更する",
	DeleteVersionConfirmationTextTemplate:   "このバージョン {VersionName} を削除してもよろしいですか？",

	FilterTabAllVersions:   "全てのバージョン",
	FilterTabOnlineVersion: "オンラインバージョン",
	FilterTabNamedVersions: "名付け済みバージョン",
	Rename:                 "名前の変更",
	PageOverView:           "ページ概要",
}

var Messages_pt_BR = &Messages{
	ModuleDescription:                       "A publicação: versões, agendamentos e colocar registros no ar.",
	StatusDraft:                             "Rascunho",
	StatusOnline:                            "No ar",
	StatusOffline:                           "Fora do ar",
	Publication:                             "Publicação",
	Publish:                                 "Publicar",
	PublishHelp:                             "Publica somente se a situação for 'rascunho' ou 'fora do ar'",
	PublishOrRepublish:                      "Publicar/Republicar",
	PublishOrRepublishHelp:                  "Se a situação for 'rascunho' ou 'fora do ar', publica; caso contrário, republica",
	Unpublish:                               "Despublicar",
	UnpublishHelp:                           "Despublica se a situação for 'no ar'",
	Republish:                               "Republicar",
	RepublishHelp:                           "Republica se a situação for 'no ar'",
	Areyousure:                              "Tem certeza?",
	ScheduledStartAt:                        "Início em",
	ScheduledEndAt:                          "Fim em",
	ScheduledStartAtShouldLaterThanNow:      "O início precisa ser depois de agora",
	ScheduledEndAtShouldLaterThanNowOrEmpty: "O fim precisa ser depois de agora, ou ficar em branco",
	ScheduledEndAtShouldLaterThanStartAt:    "O fim precisa ser depois do início",
	ScheduledStartAtShouldNotEmpty:          "O início não pode ficar em branco",
	PublishedAt:                             "Início em",
	UnPublishedAt:                           "Fim em",
	ActualPublishTime:                       "Publicação efetiva",
	SchedulePublishTime:                     "Publicação agendada",
	NotSet:                                  "Não definido",
	WhenDoYouWantToPublish:                  "Quando você quer publicar?",
	PublishScheduleTip:                      "Depois de definir a {SchedulePublishTime}, o sistema publica/despublica sozinho.",
	DateTimePickerClearText:                 "Limpar",
	DateTimePickerOkText:                    "OK",
	SaveAsNewVersion:                        "Salvar como nova versão",
	SwitchedToNewVersion:                    "Trocado para a nova versão",
	SuccessfullyCreated:                     "Criado com sucesso",
	SuccessfullyRename:                      "Renomeado com sucesso",
	SuccessfullyPublished:                   "Publicado com sucesso",
	SuccessfullyPublishedOrRepublished:      "Publicado/republicado com sucesso",
	SuccessfullyUnpublished:                 "Despublicado com sucesso",
	SuccessfullyRepublished:                 "Republicado com sucesso",
	OnlineVersion:                           "Versão no ar",
	VersionsList:                            "Lista de versões",
	AllVersions:                             "Todas as versões",
	NamedVersions:                           "Versões nomeadas",
	RenameVersion:                           "Renomear versão",
	DeleteVersionConfirmationTextTemplate:   "Tem certeza de que quer excluir a versão {VersionName}?",
	BulkActionConfirmationTextTemplate:      "Tem certeza de que quer <b>{Action}</b> os registros abaixo?",
	BulkActionNoRecordsTextTemplate:         "Nenhum registro para <b>{Action}</b>.",
	FilterTabAllVersions:                    "Todas as versões",
	FilterTabOnlineVersion:                  "Versões no ar",
	FilterTabNamedVersions:                  "Versões nomeadas",
	Rename:                                  "Renomear",
	PageOverView:                            "Visão geral da página",
	Duplicate:                               "Duplicar",
}
