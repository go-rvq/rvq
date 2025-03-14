package presets

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/qor5/x/v3/i18n"
)

func MustGetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, CoreI18nModuleKey, DefaultMessages).(*Messages)
}

type TimeFormatMessages struct {
	Date     string
	Time     string
	DateTime string
}

type StrMap map[string]string

func (m StrMap) Get(key string) string {
	return m[key]
}

func (m StrMap) Set(pair ...string) {
	if len(pair)%2 != 0 {
		panic("pairs must have pairs")
	}
	for i := 0; i < len(pair); i += 2 {
		m[pair[i]] = pair[i+1]
	}
}

type Messages struct {
	SuccessfullyUpdated string
	SuccessfullyCreated string
	SuccessfullyDeleted string
	Search              string
	TheFemaleTitle      string
	TheMaleTitle        string

	YouAreHere                                 string
	New                                        string
	Update                                     string
	Execute                                    string
	Delete                                     string
	Edit                                       string
	FormTitle                                  string
	OK                                         string
	Cancel                                     string
	Clear                                      string
	Create                                     string
	DeleteConfirmationTextTemplate             string
	CreatingFemaleObjectTitleTemplate          string
	CreatingObjectTitleTemplate                string
	EditingObjectTitleTemplate                 string
	ListingObjectTitleTemplate                 string
	DetailingObjectTitleTemplate               string
	FiltersClear                               string
	FiltersAdd                                 string
	FilterApply                                string
	FilterByTemplate                           string
	FiltersDateInTheLast                       string
	FiltersDateEquals                          string
	FiltersDateBetween                         string
	FiltersDateIsAfter                         string
	FiltersDateIsAfterOrOn                     string
	FiltersDateIsBefore                        string
	FiltersDateIsBeforeOrOn                    string
	FiltersDateDays                            string
	FiltersDateMonths                          string
	FiltersDateAnd                             string
	FiltersTo                                  string
	FiltersNumberEquals                        string
	FiltersNumberBetween                       string
	FiltersNumberGreaterThan                   string
	FiltersNumberLessThan                      string
	FiltersNumberAnd                           string
	FiltersStringEquals                        string
	FiltersStringContains                      string
	FiltersMultipleSelectIn                    string
	FiltersMultipleSelectNotIn                 string
	Month                                      string
	MonthNames                                 [time.December + 1]string
	Year                                       string
	PaginationRowsPerPage                      string
	PaginationPageInfo                         string
	PaginationPage                             string
	PaginationOfPage                           string
	ListingNoRecordToShow                      string
	ListingSelectedCountNotice                 string
	ListingClearSelection                      string
	BulkActionNoAvailableRecords               string
	BulkActionSelectedIdsProcessNoticeTemplate string
	ConfirmDialogPromptText                    string
	Language                                   string
	Colon                                      string
	NotFoundPageNotice                         string
	AddRow                                     string

	BulkActionConfirmationTextTemplate string

	TimeFormats TimeFormatMessages

	Common StrMap

	Error           string
	ErrEmptyParamID error

	CopiedToClipboard string
}

func (msgr *Messages) TheTitle(female bool, title string, args ...string) string {
	if female {
		return fmt.Sprintf(msgr.TheFemaleTitle, title)
	}
	return fmt.Sprintf(msgr.TheMaleTitle, title)
}

func (msgr *Messages) DeleteConfirmationText(model, theModelTitle, title string) string {
	return strings.NewReplacer("{model}", model, "{the_model}", theModelTitle, "{title}", title).
		Replace(msgr.DeleteConfirmationTextTemplate)
}

func (msgr *Messages) CreatingObjectTitle(modelName string, female bool) string {
	tmpl := msgr.CreatingObjectTitleTemplate
	if female && msgr.CreatingFemaleObjectTitleTemplate != "" {
		tmpl = msgr.CreatingFemaleObjectTitleTemplate
	}
	return strings.NewReplacer("{modelName}", modelName).
		Replace(tmpl)
}

func (msgr *Messages) EditingObjectTitle(label string, name string) string {
	return strings.NewReplacer("{id}", name, "{modelName}", label).
		Replace(msgr.EditingObjectTitleTemplate)
}

func (msgr *Messages) ListingObjectTitle(label string) string {
	return strings.NewReplacer("{modelName}", label).
		Replace(msgr.ListingObjectTitleTemplate)
}

func (msgr *Messages) DetailingObjectTitle(label string, name string) string {
	return strings.NewReplacer("{id}", name, "{modelName}", label).
		Replace(msgr.DetailingObjectTitleTemplate)
}

func (msgr *Messages) BulkActionSelectedIdsProcessNotice(ids string) string {
	return strings.NewReplacer("{ids}", ids).
		Replace(msgr.BulkActionSelectedIdsProcessNoticeTemplate)
}

func (msgr *Messages) FilterBy(filter string) string {
	return strings.NewReplacer("{filter}", filter).
		Replace(msgr.FilterByTemplate)
}

func (msgr *Messages) BulkActionConfirmationText(action string) string {
	return strings.NewReplacer("{Action}", action).
		Replace(msgr.BulkActionConfirmationTextTemplate)
}

var Messages_en_US = &Messages{
	YouAreHere:                        "You Are Here",
	SuccessfullyUpdated:               "Successfully Updated",
	SuccessfullyCreated:               "Successfully Created",
	SuccessfullyDeleted:               "Successfully Deleted",
	Search:                            "Search",
	New:                               "New",
	Update:                            "Update",
	Execute:                           "Execute",
	Delete:                            "Delete",
	Edit:                              "Edit",
	FormTitle:                         "Form",
	OK:                                "OK",
	Cancel:                            "Cancel",
	Clear:                             "Clear",
	Create:                            "Create",
	DeleteConfirmationTextTemplate:    "Are you sure you want to delete {the_model}: {title}?",
	CreatingObjectTitleTemplate:       "New {modelName}",
	CreatingFemaleObjectTitleTemplate: "New {modelName}",
	EditingObjectTitleTemplate:        "Editing {modelName} {id}",
	ListingObjectTitleTemplate:        "Listing {modelName}",
	DetailingObjectTitleTemplate:      "{modelName} {id}",
	FiltersClear:                      "Clear Filters",
	FiltersAdd:                        "Add Filters",
	FilterApply:                       "Apply",
	FilterByTemplate:                  "Filter by {filter}",
	FiltersDateInTheLast:              "is in the last",
	FiltersDateEquals:                 "is equal to",
	FiltersDateBetween:                "is between",
	FiltersDateIsAfter:                "is after",
	FiltersDateIsAfterOrOn:            "is on or after",
	FiltersDateIsBefore:               "is before",
	FiltersDateIsBeforeOrOn:           "is before or on",
	FiltersDateDays:                   "days",
	FiltersDateMonths:                 "months",
	FiltersDateAnd:                    "and",
	FiltersTo:                         "to",
	FiltersNumberEquals:               "is equal to",
	FiltersNumberBetween:              "between",
	FiltersNumberGreaterThan:          "is greater than",
	FiltersNumberLessThan:             "is less than",
	FiltersNumberAnd:                  "and",
	FiltersStringEquals:               "is equal to",
	FiltersStringContains:             "contains",
	FiltersMultipleSelectIn:           "in",
	FiltersMultipleSelectNotIn:        "not in",
	Month:                             "Month",
	MonthNames: [time.December + 1]string{
		"", "January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	},
	Year:                                       "Year",
	PaginationRowsPerPage:                      "Rows per page: ",
	ListingNoRecordToShow:                      "No records to show",
	ListingSelectedCountNotice:                 "{count} records are selected. ",
	ListingClearSelection:                      "clear selection",
	BulkActionNoAvailableRecords:               "None of the selected records can be executed with this action.",
	BulkActionSelectedIdsProcessNoticeTemplate: "Partially selected records cannot be executed with this action: {ids}.",
	ConfirmDialogPromptText:                    "Are you sure?",
	Language:                                   "Language",
	Colon:                                      ":",
	NotFoundPageNotice:                         "Sorry, the requested page cannot be found. Please check the URL.",
	AddRow:                                     "Add Row",

	BulkActionConfirmationTextTemplate: "Are you sure you want to <b>{Action}</b> below records?",

	Error:             "ERROR",
	ErrEmptyParamID:   errors.New("Empty param ID"),
	CopiedToClipboard: "Copied to clipboard",
}

var DefaultMessages = Messages_en_US

var Messages_zh_CN = &Messages{
	SuccessfullyUpdated:            "成功更新了",
	Search:                         "搜索",
	New:                            "新建",
	Update:                         "更新",
	Delete:                         "删除",
	Edit:                           "编辑",
	FormTitle:                      "表单",
	OK:                             "确定",
	Cancel:                         "取消",
	Clear:                          "清空",
	Create:                         "创建",
	DeleteConfirmationTextTemplate: "你确定你要删除这个对象吗，对象: {title}?",
	CreatingObjectTitleTemplate:    "新建{modelName}",
	EditingObjectTitleTemplate:     "编辑{modelName} {id}",
	ListingObjectTitleTemplate:     "{modelName}列表",
	DetailingObjectTitleTemplate:   "{modelName} {id}",
	FiltersClear:                   "清空筛选器",
	FiltersAdd:                     "添加筛选器",
	FilterApply:                    "应用",
	FilterByTemplate:               "按{filter}筛选",
	FiltersDateInTheLast:           "过去",
	FiltersDateEquals:              "等于",
	FiltersDateBetween:             "之间",
	FiltersDateIsAfter:             "之后",
	FiltersDateIsAfterOrOn:         "当天或之后",
	FiltersDateIsBefore:            "之前",
	FiltersDateIsBeforeOrOn:        "当天或之前",
	FiltersDateDays:                "天",
	FiltersDateMonths:              "月",
	FiltersDateAnd:                 "和",
	FiltersTo:                      "至",
	FiltersNumberEquals:            "等于",
	FiltersNumberBetween:           "之间",
	FiltersNumberGreaterThan:       "大于",
	FiltersNumberLessThan:          "小于",
	FiltersNumberAnd:               "和",
	FiltersStringEquals:            "等于",
	FiltersStringContains:          "包含",
	FiltersMultipleSelectIn:        "包含",
	FiltersMultipleSelectNotIn:     "不包含",
	PaginationRowsPerPage:          "每页: ",
	ListingNoRecordToShow:          "没有可显示的记录",
	ListingSelectedCountNotice:     "{count}条记录被选中。",
	ListingClearSelection:          "清除选择",
	BulkActionNoAvailableRecords:   "所有选中的记录均无法执行这个操作。",
	BulkActionSelectedIdsProcessNoticeTemplate: "部分选中的记录无法被执行这个操作: {ids}。",
	ConfirmDialogPromptText:                    "你确定吗?",
	Language:                                   "语言",
	Colon:                                      "：",
	NotFoundPageNotice:                         "很抱歉，所请求的页面不存在，请检查URL。",
}

var Messages_ja_JP = &Messages{
	SuccessfullyUpdated:            "更新に成功しました",
	Search:                         "検索する",
	New:                            "新規",
	Update:                         "更新する",
	Delete:                         "削除する",
	Edit:                           "編集する",
	FormTitle:                      "フォーム",
	OK:                             "OK",
	Cancel:                         "キャンセル",
	Create:                         "新規作成",
	DeleteConfirmationTextTemplate: ": {title}を削除して本当によろしいですか？",
	CreatingObjectTitleTemplate:    "{modelName} を作る",
	EditingObjectTitleTemplate:     "{modelName} {id} を編集する",
	ListingObjectTitleTemplate:     "リスティング {modelName} ",
	DetailingObjectTitleTemplate:   "{modelName} {id} ",
	FiltersClear:                   "フィルターをクリアする",
	FiltersAdd:                     "フィルターを追加する",
	FilterApply:                    "適用する",
	FilterByTemplate:               "{filter} でフィルターする",
	FiltersDateInTheLast:           "がサイト",
	FiltersDateEquals:              "と同等",
	FiltersDateBetween:             "の間",
	FiltersDateIsAfter:             "が後",
	FiltersDateIsAfterOrOn:         "が同時または後",
	FiltersDateIsBefore:            "が前",
	FiltersDateIsBeforeOrOn:        "が前または同時",
	FiltersDateDays:                "日間",
	FiltersDateMonths:              "月数",
	FiltersDateAnd:                 "＆",
	FiltersTo:                      "から",
	FiltersNumberEquals:            "と同等",
	FiltersNumberBetween:           "間",
	FiltersNumberGreaterThan:       "より大きい",
	FiltersNumberLessThan:          "より少ない",
	FiltersNumberAnd:               "＆",
	FiltersStringEquals:            "と同等",
	FiltersStringContains:          "を含む",
	FiltersMultipleSelectIn:        "中",
	FiltersMultipleSelectNotIn:     "以外",
	PaginationRowsPerPage:          "行 / ページ",
	ListingNoRecordToShow:          "表示できるデータはありません",
	ListingSelectedCountNotice:     "{count} レコードが選択されています",
	ListingClearSelection:          "選択をクリア",
	BulkActionNoAvailableRecords:   "この機能はご利用いただけません",
	BulkActionSelectedIdsProcessNoticeTemplate: "この一部の機能はご利用いただけません: {ids}",
	ConfirmDialogPromptText:                    "よろしいですか？",
	Language:                                   "言語",
	Colon:                                      ":",
	NotFoundPageNotice:                         "申し訳ありませんが、リクエストされたページは見つかりませんでした。URLを確認してください。",
}
