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
	Category                       string           `i18n:"hint='Label of a page\\'s category.'"`
	Preview                        string           `i18n:"hint='Button that shows the page as it will look.'"`
	Containers                     string           `i18n:"hint='Title of the blocks (containers) a page is built of.'"`
	AddContainers                  string           `i18n:"hint='Button that adds blocks to a page.'"`
	New                            string           `i18n:"hint='Button that creates a new item of the page builder.'"`
	Shared                         string           `i18n:"hint='Mark of a block shared between pages.'"`
	Select                         string           `i18n:"hint='Button that chooses a block or a template.'"`
	SelectedTemplateLabel          string           `i18n:"label='Selected template', hint='Label of the template a page was made from.'"`
	CreateFromTemplate             string           `i18n:"hint='Action that creates a page from a template.'"`
	ChangeTemplate                 string           `i18n:"hint='Action that changes the template of a page.'"`
	RelatedOnlinePages             string           `i18n:"hint='Title of the published pages that use a shared block.'"`
	RepublishAllRelatedOnlinePages string           `i18n:"label='Republish all related pages', hint='Button that publishes again every page that uses a shared block.'"`
	Unnamed                        string           `i18n:"hint='Shown for a page or block with no name.'"`
	NotDescribed                   string           `i18n:"hint='Shown for a page or block with no description.'"`
	Blank                          string           `i18n:"hint='The template of an empty page.'"`
	NewPage                        string           `i18n:"hint='Title of the creation of a page.'"`
	FilterTabAllVersions           string           `i18n:"label='Tab: all versions', hint='Tab that shows every version of a page.'"`
	FilterTabOnlineVersion         string           `i18n:"label='Tab: online versions', hint='Tab that shows the published versions of a page.'"`
	FilterTabNamedVersions         string           `i18n:"label='Tab: named versions', hint='Tab that shows the versions given a name.'"`
	Rename                         string           `i18n:"hint='Action that renames a version of a page.'"`
	PageOverView                   string           `i18n:"label='Page overview', hint='Title of the summary of a page.'"`
	ErrPermissionDenied            i18n.ErrorString `i18n:"hint='Error when the user may not do what was asked.'"`
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

var Messages_pt_BR = &Messages{
	Category:                       "Categoria",
	Preview:                        "Visualizar",
	Containers:                     "Blocos",
	AddContainers:                  "Adicionar blocos",
	New:                            "Novo",
	Shared:                         "Compartilhado",
	Select:                         "Selecionar",
	SelectedTemplateLabel:          "Modelo",
	CreateFromTemplate:             "Criar a partir de um modelo",
	ChangeTemplate:                 "Trocar o modelo",
	RelatedOnlinePages:             "Páginas no ar relacionadas",
	RepublishAllRelatedOnlinePages: "Republicar todas",
	Unnamed:                        "Sem nome",
	NotDescribed:                   "Sem descrição",
	Blank:                          "Em branco",
	NewPage:                        "Nova página",
	FilterTabAllVersions:           "Todas as versões",
	FilterTabOnlineVersion:         "Versões no ar",
	FilterTabNamedVersions:         "Versões nomeadas",
	Rename:                         "Renomear",
	PageOverView:                   "Visão geral da página",
	ErrPermissionDenied:            "Permissão negada",
}
