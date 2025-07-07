package pagebuilder

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
)

const I18nPageBuilderKey i18n.ModuleKey = "I18nPageBuilderKey"

func MustGetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, I18nPageBuilderKey, Messages_en_US).(*Messages)
}

type Messages struct {
	Category                       string
	Preview                        string
	Containers                     string
	AddContainers                  string
	New                            string
	Shared                         string
	Select                         string
	SelectedTemplateLabel          string
	CreateFromTemplate             string
	ChangeTemplate                 string
	RelatedOnlinePages             string
	RepublishAllRelatedOnlinePages string
	Unnamed                        string
	NotDescribed                   string
	Blank                          string
	NewPage                        string
	FilterTabAllVersions           string
	FilterTabOnlineVersion         string
	FilterTabNamedVersions         string
	Rename                         string
	PageOverView                   string
	ErrPermissionDenied            i18n.ErrorString
}

var Messages_en_US = &Messages{
	Category:                       "Category",
	Preview:                        "Preview",
	Containers:                     "Containers",
	AddContainers:                  "Add Containers",
	New:                            "New",
	Shared:                         "Shared",
	Select:                         "Select",
	SelectedTemplateLabel:          "Template",
	CreateFromTemplate:             "Create From Template",
	ChangeTemplate:                 "Change Template",
	RelatedOnlinePages:             "Related Online Pages",
	RepublishAllRelatedOnlinePages: "Republish All",
	Unnamed:                        "Unnamed",
	NotDescribed:                   "Not Described",
	Blank:                          "Blank",
	NewPage:                        "New Page",
	FilterTabAllVersions:           "All Versions",
	FilterTabOnlineVersion:         "Online Versions",
	FilterTabNamedVersions:         "Named Versions",
	Rename:                         "Rename",
	PageOverView:                   "Page Overview",
	ErrPermissionDenied:            "Permission Denied",
}

var Messages_zh_CN = &Messages{
	Category:                       "目录",
	Preview:                        "预览",
	Containers:                     "组件",
	AddContainers:                  "增加组件",
	New:                            "新增",
	Shared:                         "公用的",
	Select:                         "选择",
	SelectedTemplateLabel:          "模板",
	CreateFromTemplate:             "从模板中创建",
	ChangeTemplate:                 "更改模版",
	RelatedOnlinePages:             "相关在线页面",
	RepublishAllRelatedOnlinePages: "重新发布所有页面",
	Unnamed:                        "未命名",
	NotDescribed:                   "未描述",
	Blank:                          "空白",
	NewPage:                        "新页面",
	FilterTabAllVersions:           "所有版本",
	FilterTabOnlineVersion:         "在线版本",
	FilterTabNamedVersions:         "已命名版本",
	Rename:                         "重命名",
	PageOverView:                   "页面概览",
	ErrPermissionDenied:            "沒有權限",
}

var Messages_ja_JP = &Messages{
	Category:                       "カテゴリー",
	Preview:                        "プレビュー",
	Containers:                     "コンテナ",
	AddContainers:                  "コンテナを追加する",
	New:                            "新規",
	Shared:                         "共有",
	Select:                         "選択する",
	SelectedTemplateLabel:          "テンプレート",
	CreateFromTemplate:             "テンプレートから新規作成する",
	ChangeTemplate:                 "テンプレートを変更する",
	RelatedOnlinePages:             "関連オンラインページ",
	RepublishAllRelatedOnlinePages: "すべて再公開",
	Unnamed:                        "名前なし",
	NotDescribed:                   "記述されていません",
	Blank:                          "空白",
	NewPage:                        "新しいページ",
	FilterTabAllVersions:           "全てのバージョン",
	FilterTabOnlineVersion:         "オンラインバージョン",
	FilterTabNamedVersions:         "名付け済みバージョン",
	Rename:                         "名前の変更",
	PageOverView:                   "ページ概要",
	ErrPermissionDenied:            "許可が拒否されました",
}
