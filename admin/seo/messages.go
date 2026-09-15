package seo

type Messages struct {
	Variable             string
	VariableDescription  string
	Basic                string
	Title                string
	Description          string
	Keywords             string
	OpenGraphInformation string
	OpenGraphTitle       string
	OpenGraphDescription string
	OpenGraphURL         string
	OpenGraphType        string
	OpenGraphImageURL    string
	OpenGraphImage       string
	OpenGraphMetadata    string
	Seo                  string
	Customize            string

	// Help overlay (see help.go). Body fields are HTML.
	HelpTooltip          string
	HelpTitle            string
	HelpIntroTitle       string
	HelpIntro            string
	HelpVariablesTitle   string
	HelpVariables        string
	HelpInheritanceTitle string
	HelpInheritance      string

	// Setting-variable labels: the label of a setting variable is looked up as
	// strcase.ToCamel("SettingVar " + varName), so the built-in SiteName variable
	// maps to SettingVarSiteName. An app that adds variables adds fields here (or
	// registers them in this module).
	SettingVarSiteName string
}

var Messages_en_US = &Messages{
	Variable:             "Variables Setting",
	Basic:                "Basic",
	Title:                "Title",
	Description:          "Description",
	Keywords:             "Keywords",
	OpenGraphInformation: "Open Graph Information",
	OpenGraphTitle:       "Open Graph Title",
	OpenGraphDescription: "Open Graph Description",
	OpenGraphURL:         "Open Graph URL",
	OpenGraphType:        "Open Graph Type",
	OpenGraphImageURL:    "Open Graph Image URL",
	OpenGraphImage:       "Open Graph Image",
	OpenGraphMetadata:    "Open Graph Metadata",
	Seo:                  "SEO",
	Customize:            "Customize",

	HelpTooltip:          "How SEO works",
	HelpTitle:            "How SEO works",
	HelpIntroTitle:       "The SEO fields",
	HelpIntro:            "These fields become the page's metadata — the <b>Title</b> shown in the browser tab and search results, the <b>Description</b> and <b>Keywords</b>, and the <b>Open Graph</b> fields that social networks read when the page is shared. Turn on <b>Customize</b> to give this record its own SEO; leave a field empty to inherit it (see below).",
	HelpVariablesTitle:   "Variables",
	HelpVariables:        "A field value may contain variables written as <code>{{Name}}</code> — for example <code>{{SiteName}}</code> — replaced when the page is rendered. Click a variable chip above a field to insert it.",
	HelpInheritanceTitle: "Inheritance &amp; fallback",
	HelpInheritance:      "SEO is layered: an empty field takes its value from the more general SEO in the chain, ending at the Global SEO. So at each level you set only what should differ from the level above.",

	SettingVarSiteName: "Site Name",
}

var Messages_zh_CN = &Messages{
	Variable:             "变量设置",
	Basic:                "基本信息",
	Title:                "标题",
	Description:          "描述",
	Keywords:             "关键词",
	OpenGraphInformation: "OG 信息",
	OpenGraphTitle:       "OG 标题",
	OpenGraphDescription: "OG 描述",
	OpenGraphURL:         "OG 链接",
	OpenGraphType:        "OG 类型",
	OpenGraphImageURL:    "OG 图片链接",
	OpenGraphImage:       "OG 图片",
	OpenGraphMetadata:    "OG 元数据",
	Seo:                  "搜索引擎优化",
	Customize:            "自定义",

	HelpTooltip:          "SEO 工作原理",
	HelpTitle:            "SEO 工作原理",
	HelpIntroTitle:       "SEO 字段",
	HelpIntro:            "这些字段构成页面的元数据 —— 浏览器标签页和搜索结果中显示的<b>标题</b>、<b>描述</b>和<b>关键词</b>，以及分享页面时社交网络读取的 <b>Open Graph</b> 字段。打开<b>自定义</b>为此记录设置独立的 SEO；留空的字段将被继承（见下文）。",
	HelpVariablesTitle:   "变量",
	HelpVariables:        "字段值可以包含以 <code>{{Name}}</code> 形式书写的变量 —— 例如 <code>{{SiteName}}</code> —— 在页面渲染时被替换。点击字段上方的变量标签即可插入。",
	HelpInheritanceTitle: "继承与回退",
	HelpInheritance:      "SEO 是分层的：留空的字段会从链中更通用的 SEO 取值，最终回退到全局 SEO。因此每一层只需设置与上一层不同的部分。",

	SettingVarSiteName: "站点名称",
}

var Messages_pt_BR = &Messages{
	Variable:             "Configuração das variáveis",
	Basic:                "Básico",
	Title:                "Título",
	Description:          "Descrição",
	Keywords:             "Palavras-chave",
	OpenGraphInformation: "Informações do Open Graph",
	OpenGraphTitle:       "Título do Open Graph",
	OpenGraphDescription: "Descrição do Open Graph",
	OpenGraphURL:         "URL do Open Graph",
	OpenGraphType:        "Tipo do Open Graph",
	OpenGraphImageURL:    "URL da imagem do Open Graph",
	OpenGraphImage:       "Imagem do Open Graph",
	OpenGraphMetadata:    "Metadados do Open Graph",
	Seo:                  "SEO",
	Customize:            "Personalizar",

	HelpTooltip:          "Como o SEO funciona",
	HelpTitle:            "Como o SEO funciona",
	HelpIntroTitle:       "Os campos de SEO",
	HelpIntro:            "Estes campos viram os metadados da página — o <b>Título</b> mostrado na aba do navegador e nos resultados de busca, a <b>Descrição</b> e as <b>Palavras-chave</b>, e os campos de <b>Open Graph</b> que as redes sociais leem quando a página é compartilhada. Ative <b>Personalizar</b> para dar a este registro o SEO próprio; deixe um campo vazio para herdá-lo (veja abaixo).",
	HelpVariablesTitle:   "Variáveis",
	HelpVariables:        "O valor de um campo pode conter variáveis escritas como <code>{{Nome}}</code> — por exemplo <code>{{SiteName}}</code> — substituídas quando a página é renderizada. Clique num chip de variável acima do campo para inseri-la.",
	HelpInheritanceTitle: "Herança e fallback",
	HelpInheritance:      "O SEO é em camadas: um campo vazio assume o valor do SEO mais geral na cadeia, terminando no SEO Global. Assim, em cada nível você define apenas o que deve diferir do nível acima.",

	SettingVarSiteName: "Nome do Site",
}
