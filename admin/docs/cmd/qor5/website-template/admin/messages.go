package admin

import (
	"github.com/go-rvq/rvq/x/i18n"
)

const I18nExampleKey i18n.ModuleKey = "I18nExampleKey"

type Messages struct {
	FilterTabsAll            string
	FilterTabsHasUnreadNotes string
	FilterTabsActive         string
}

var Messages_en_US = &Messages{
	FilterTabsAll:            "All",
	FilterTabsHasUnreadNotes: "Has Unread Notes",
	FilterTabsActive:         "Active",
}

var Messages_zh_CN = &Messages{
	FilterTabsAll:            "全部",
	FilterTabsHasUnreadNotes: "未读备注",
	FilterTabsActive:         "有效",
}

type Messages_ModelsI18nModuleKey struct {
	QOR5Example string
	Roles       string
	Users       string

	Posts          string
	PostsID        string
	PostsTitle     string
	PostsHeroImage string
	PostsBody      string
	Example        string
	Settings       string
	Post           string
	PostsBodyImage string

	SeoPost             string
	SeoVariableTitle    string
	SeoVariableSiteName string

	PageBuilder              string
	Pages                    string
	SharedContainers         string
	DemoContainers           string
	Templates                string
	Categories               string
	ECManagement             string
	ECDashboard              string
	Orders                   string
	Products                 string
	SiteManagement           string
	SEO                      string
	UserManagement           string
	Profile                  string
	FeaturedModelsManagement string
	InputHarnesses           string
	ListEditorExample        string
	Customers                string
	ListModels               string
	MicrositeModels          string
	Workers                  string
	ActivityLogs             string
	MediaLibrary             string

	PagesID         string
	PagesTitle      string
	PagesSlug       string
	PagesLocale     string
	PagesNotes      string
	PagesDraftCount string
	PagesOnline     string

	Page                   string
	PagesStatus            string
	PagesSchedule          string
	PagesCategoryID        string
	PagesTemplateSelection string
	PagesEditContainer     string

	WebHeader       string
	WebHeadersColor string
	Header          string

	WebFooter             string
	WebFootersEnglishUrl  string
	WebFootersJapaneseUrl string
	Footer                string

	VideoBanner                       string
	VideoBannersAddTopSpace           string
	VideoBannersAddBottomSpace        string
	VideoBannersAnchorID              string
	VideoBannersVideo                 string
	VideoBannersBackgroundVideo       string
	VideoBannersMobileBackgroundVideo string
	VideoBannersVideoCover            string
	VideoBannersMobileVideoCover      string
	VideoBannersHeading               string
	VideoBannersPopupText             string
	VideoBannersText                  string
	VideoBannersLinkText              string
	VideoBannersLink                  string

	Heading                   string
	HeadingsAddTopSpace       string
	HeadingsAddBottomSpace    string
	HeadingsAnchorID          string
	HeadingsHeading           string
	HeadingsFontColor         string
	HeadingsBackgroundColor   string
	HeadingsLink              string
	HeadingsLinkText          string
	HeadingsLinkDisplayOption string
	HeadingsText              string

	BrandGrid                string
	BrandGridsAddTopSpace    string
	BrandGridsAddBottomSpace string
	BrandGridsAnchorID       string
	BrandGridsBrands         string

	ListContent                   string
	ListContentsAddTopSpace       string
	ListContentsAddBottomSpace    string
	ListContentsAnchorID          string
	ListContentsBackgroundColor   string
	ListContentsItems             string
	ListContentsLink              string
	ListContentsLinkText          string
	ListContentsLinkDisplayOption string

	ImageContainer                           string
	ImageContainersAddTopSpace               string
	ImageContainersAddBottomSpace            string
	ImageContainersAnchorID                  string
	ImageContainersBackgroundColor           string
	ImageContainersTransitionBackgroundColor string
	ImageContainersImage                     string
	Image                                    string

	InNumber                string
	InNumbersAddTopSpace    string
	InNumbersAddBottomSpace string
	InNumbersAnchorID       string
	InNumbersHeading        string
	InNumbersItems          string
	InNumbers               string

	ContactForm                    string
	ContactFormsAddTopSpace        string
	ContactFormsAddBottomSpace     string
	ContactFormsAnchorID           string
	ContactFormsHeading            string
	ContactFormsText               string
	ContactFormsSendButtonText     string
	ContactFormsFormButtonText     string
	ContactFormsMessagePlaceholder string
	ContactFormsNamePlaceholder    string
	ContactFormsEmailPlaceholder   string
	ContactFormsThankyouMessage    string
	ContactFormsActionUrl          string
	ContactFormsPrivacyPolicy      string
}

var Messages_zh_CN_ModelsI18nModuleKey = &Messages_ModelsI18nModuleKey{
	Posts:          "帖子 示例",
	PostsID:        "ID",
	PostsTitle:     "标题",
	PostsHeroImage: "主图",
	PostsBody:      "内容",
	Example:        "QOR5演示",
	Settings:       "SEO 设置",
	Post:           "帖子",
	PostsBodyImage: "内容图片",

	SeoPost:             "帖子",
	SeoVariableTitle:    "标题",
	SeoVariableSiteName: "站点名称",

	QOR5Example: "QOR5 示例",
	Roles:       "权限管理",
	Users:       "用户管理",

	PageBuilder:              "页面管理菜单",
	Pages:                    "页面管理",
	SharedContainers:         "公用组件",
	DemoContainers:           "示例组件",
	Templates:                "模板页面",
	Categories:               "目录管理",
	ECManagement:             "电子商务管理",
	ECDashboard:              "电子商务仪表盘",
	Orders:                   "订单管理",
	Products:                 "产品管理",
	SiteManagement:           "站点管理菜单",
	SEO:                      "SEO 管理",
	UserManagement:           "用户管理菜单",
	Profile:                  "个人页面",
	FeaturedModelsManagement: "特色模块管理菜单",
	InputHarnesses:           "Input 示例",
	ListEditorExample:        "ListEditor 示例",
	Customers:                "Customers 示例",
	ListModels:               "发布带排序及分页模块 示例",
	MicrositeModels:          "Microsite 示例",
	Workers:                  "后台工作进程管理",
	ActivityLogs:             "操作日志",
	MediaLibrary:             "媒体库",

	PagesID:         "ID",
	PagesTitle:      "标题",
	PagesSlug:       "Slug",
	PagesLocale:     "地区",
	PagesNotes:      "备注",
	PagesDraftCount: "草稿数",
	PagesOnline:     "在线",

	Page:                   "Page",
	PagesStatus:            "PagesStatus",
	PagesSchedule:          "PagesSchedule",
	PagesCategoryID:        "PagesCategoryID",
	PagesTemplateSelection: "PagesTemplateSelection",
	PagesEditContainer:     "PagesEditContainer",

	WebHeader:       "WebHeader",
	WebHeadersColor: "WebHeadersColor",
	Header:          "Header",

	WebFooter:             "WebFooter",
	WebFootersEnglishUrl:  "WebFootersEnglishUrl",
	WebFootersJapaneseUrl: "WebFootersJapaneseUrl",
	Footer:                "Footer",

	VideoBanner:                       "VideoBanner",
	VideoBannersAddTopSpace:           "VideoBannersAddTopSpace",
	VideoBannersAddBottomSpace:        "VideoBannersAddBottomSpace",
	VideoBannersAnchorID:              "VideoBannersAnchorID",
	VideoBannersVideo:                 "VideoBannersVideo",
	VideoBannersBackgroundVideo:       "VideoBannersBackgroundVideo",
	VideoBannersMobileBackgroundVideo: "VideoBannersMobileBackgroundVideo",
	VideoBannersVideoCover:            "VideoBannersVideoCover",
	VideoBannersMobileVideoCover:      "VideoBannersMobileVideoCover",
	VideoBannersHeading:               "VideoBannersHeading",
	VideoBannersPopupText:             "VideoBannersPopupText",
	VideoBannersText:                  "VideoBannersText",
	VideoBannersLinkText:              "VideoBannersLinkText",
	VideoBannersLink:                  "VideoBannersLink",

	Heading:                   "Heading",
	HeadingsAddTopSpace:       "HeadingsAddTopSpace",
	HeadingsAddBottomSpace:    "HeadingsAddBottomSpace",
	HeadingsAnchorID:          "HeadingsAnchorID",
	HeadingsHeading:           "HeadingsHeading",
	HeadingsFontColor:         "HeadingsFontColor",
	HeadingsBackgroundColor:   "HeadingsBackgroundColor",
	HeadingsLink:              "HeadingsLink",
	HeadingsLinkText:          "HeadingsLinkText",
	HeadingsLinkDisplayOption: "HeadingsLinkDisplayOption",
	HeadingsText:              "HeadingsText",

	BrandGrid:                "BrandGrid",
	BrandGridsAddTopSpace:    "BrandGridsAddTopSpace",
	BrandGridsAddBottomSpace: "BrandGridsAddBottomSpace",
	BrandGridsAnchorID:       "BrandGridsAnchorID",
	BrandGridsBrands:         "BrandGridsBrands",

	ListContent:                   "ListContent",
	ListContentsAddTopSpace:       "ListContentsAddTopSpace",
	ListContentsAddBottomSpace:    "ListContentsAddBottomSpace",
	ListContentsAnchorID:          "ListContentsAnchorID",
	ListContentsBackgroundColor:   "ListContentsBackgroundColor",
	ListContentsItems:             "ListContentsItems",
	ListContentsLink:              "ListContentsLink",
	ListContentsLinkText:          "ListContentsLinkText",
	ListContentsLinkDisplayOption: "ListContentsLinkDisplayOption",

	ImageContainer:                           "ImageContainer",
	ImageContainersAddTopSpace:               "ImageContainersAddTopSpace",
	ImageContainersAddBottomSpace:            "ImageContainersAddBottomSpace",
	ImageContainersAnchorID:                  "ImageContainersAnchorID",
	ImageContainersBackgroundColor:           "ImageContainersBackgroundColor",
	ImageContainersTransitionBackgroundColor: "ImageContainersTransitionBackgroundColor",
	ImageContainersImage:                     "ImageContainersImage",
	Image:                                    "Image",

	InNumber:                "InNumber",
	InNumbersAddTopSpace:    "InNumbersAddTopSpace",
	InNumbersAddBottomSpace: "InNumbersAddBottomSpace",
	InNumbersAnchorID:       "InNumbersAnchorID",
	InNumbersHeading:        "InNumbersHeading",
	InNumbersItems:          "InNumbersItems",
	InNumbers:               "InNumbers",

	ContactForm:                    "ContactForm",
	ContactFormsAddTopSpace:        "ContactFormsAddTopSpace",
	ContactFormsAddBottomSpace:     "ContactFormsAddBottomSpace",
	ContactFormsAnchorID:           "ContactFormsAnchorID",
	ContactFormsHeading:            "ContactFormsHeading",
	ContactFormsText:               "ContactFormsText",
	ContactFormsSendButtonText:     "ContactFormsSendButtonText",
	ContactFormsFormButtonText:     "ContactFormsFormButtonText",
	ContactFormsMessagePlaceholder: "ContactFormsMessagePlaceholder",
	ContactFormsNamePlaceholder:    "ContactFormsNamePlaceholder",
	ContactFormsEmailPlaceholder:   "ContactFormsEmailPlaceholder",
	ContactFormsThankyouMessage:    "ContactFormsThankyouMessage",
	ContactFormsActionUrl:          "ContactFormsActionUrl",
	ContactFormsPrivacyPolicy:      "ContactFormsPrivacyPolicy",
}
