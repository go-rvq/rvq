package seo

type Messages struct {
	ModuleDescription    string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Variable             string `i18n:"label='Variables setting', hint='Title of the SEO variables settings.'"`
	VariableDescription  string `i18n:"hint='Description of the SEO variables, under their title.'"`
	Basic                string `i18n:"hint='Title of the basic SEO fields (title, description, keywords).'"`
	Title                string `i18n:"hint='Label of the page title for search engines and the browser tab.'"`
	Description          string `i18n:"hint='Label of the page description for search results.'"`
	Keywords             string `i18n:"hint='Label of the page keywords.'"`
	OpenGraphInformation string `i18n:"hint='Title of the Open Graph fields (how a page shows when shared).'"`
	OpenGraphTitle       string `i18n:"hint='Label of the title shown when the page is shared.'"`
	OpenGraphDescription string `i18n:"hint='Label of the description shown when the page is shared.'"`
	OpenGraphURL         string `i18n:"label='Open Graph URL', hint='Label of the address shown when the page is shared.'"`
	OpenGraphType        string `i18n:"hint='Label of the Open Graph type of the page (website, article…).'"`
	OpenGraphImageURL    string `i18n:"label='Open Graph image URL', hint='Label of the address of the image shown when the page is shared.'"`
	OpenGraphImage       string `i18n:"hint='Label of the image shown when the page is shared.'"`
	OpenGraphMetadata    string `i18n:"hint='Label of extra Open Graph tags.'"`
	Seo                  string `i18n:"label='SEO', hint='Name of the SEO section and tab.'"`
	Customize            string `i18n:"hint='Switch that customizes the SEO of a record instead of inheriting it.'"`
	// AddVariable labels the "+ Variable" menu button in the SEO editor.
	AddVariable string `i18n:"label='Add variable', hint='Button that inserts a variable into an SEO field.'"`

	// Help overlay (see help.go). Body fields are HTML.
	HelpTooltip          string `i18n:"label='Help: tooltip', hint='Tooltip of the button that opens the SEO help.'"`
	HelpTitle            string `i18n:"label='Help: title', hint='Title of the SEO help.'"`
	HelpIntroTitle       string `i18n:"label='Help: intro title', hint='Title of the first section of the SEO help.'"`
	HelpIntro            string `i18n:"type=html, label='Help: intro', hint='First section of the SEO help: what the fields become (HTML).'"`
	HelpVariablesTitle   string `i18n:"label='Help: variables title', hint='Title of the section of the SEO help about variables.'"`
	HelpVariables        string `i18n:"type=html, label='Help: variables', hint='Section of the SEO help about the variables, {{Name}} (HTML).'"`
	HelpInheritanceTitle string `i18n:"type=html, label='Help: inheritance title', hint='Title of the section of the SEO help about inheritance (HTML).'"`
	HelpInheritance      string `i18n:"type=html, label='Help: inheritance', hint='Section of the SEO help about how an empty field inherits (HTML).'"`

	// Setting-variable labels: the label of a setting variable is looked up as
	// strcase.ToCamel("SettingVar " + varName), so the built-in SiteName variable
	// maps to SettingVarSiteName. An app that adds variables adds fields here (or
	// registers them in this module).
	SettingVarSiteName string `i18n:"label='Variable: site name', hint='Name of the SiteName variable.'"`

	// Custom variables editor.
	CustomVars      string `i18n:"hint='Title of the variables the user defines.'"`                                                                     // section heading
	CustomVarsGroup string `i18n:"label='Custom variables group', hint='Title of the group of the variables the user defines, in the variables menu.'"` // "+ Variable" menu group label

	// SEO configuration: the Google Maps API key field and its "?" help. The Help
	// body is HTML.
	SEOConfig          string `i18n:"label='SEO settings', hint='Name of the SEO settings model.'"`
	MapsKey            string `i18n:"label='Google Maps API key', hint='Label of the Google Maps key used by the ZIP-codes variable.'"`
	MapsKeyHint        string `i18n:"label='Google Maps API key: hint', hint='Hint of the Google Maps key: which APIs it needs.'"`
	MapsKeyHelpTooltip string `i18n:"label='Google Maps API key: help tooltip', hint='Tooltip of the button that opens the help on the Google Maps key.'"`
	MapsKeyHelpTitle   string `i18n:"label='Google Maps API key: help title', hint='Title of the help on how to get a Google Maps key.'"`
	MapsKeyHelp        string `i18n:"type=html, label='Google Maps API key: help', hint='Steps to get a Google Maps key (HTML).'"`
}

var Messages_en_US = &Messages{
	ModuleDescription:    "The SEO settings: titles, descriptions and the Open Graph of the pages.",
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
	AddVariable:          "Variable",

	HelpTooltip:          "How SEO works",
	HelpTitle:            "How SEO works",
	HelpIntroTitle:       "The SEO fields",
	HelpIntro:            "These fields become the page's metadata — the <b>Title</b> shown in the browser tab and search results, the <b>Description</b> and <b>Keywords</b>, and the <b>Open Graph</b> fields that social networks read when the page is shared. Turn on <b>Customize</b> to give this record its own SEO; leave a field empty to inherit it (see below).",
	HelpVariablesTitle:   "Variables",
	HelpVariables:        "A field value may contain variables written as <code>{{Name}}</code> — for example <code>{{SiteName}}</code> — replaced when the page is rendered. Click a variable chip above a field to insert it.",
	HelpInheritanceTitle: "Inheritance &amp; fallback",
	HelpInheritance:      "SEO is layered: an empty field takes its value from the more general SEO in the chain, ending at the Global SEO. So at each level you set only what should differ from the level above.",

	SettingVarSiteName: "Site Name",

	CustomVars:      "Custom variables",
	CustomVarsGroup: "Variables",

	SEOConfig:          "Settings",
	MapsKey:            "Google Maps API Key",
	MapsKeyHint:        "A browser API key (Maps JavaScript + Places) used by the ZIP-codes variable.",
	MapsKeyHelpTooltip: "How to get the API key",
	MapsKeyHelpTitle:   "How to generate the Google Maps API Key",
	MapsKeyHelp: "<p>The ZIP-codes variable opens Google Maps to pick places. For the map to work, provide a Google API key (Maps JavaScript API + Places API):</p>" +
		"<ol>" +
		"<li>Go to console.cloud.google.com and sign in.</li>" +
		"<li>Create (or select) a project.</li>" +
		"<li>Under \"APIs &amp; Services\" → \"Library\", enable \"Maps JavaScript API\" and \"Places API\".</li>" +
		"<li>Under \"APIs &amp; Services\" → \"Credentials\", click \"Create credentials\" → \"API key\".</li>" +
		"<li>Restrict the key: \"Application restrictions\" → \"HTTP referrers\" with your site domain; \"API restrictions\" limited to the two APIs above.</li>" +
		"<li>Enable billing on the project (Google requires it; there is a free monthly quota).</li>" +
		"<li>Copy the key and paste it in the \"Google Maps API Key\" field.</li>" +
		"</ol>" +
		"<p>Without the key, the ZIP-codes variable still works as plain text.</p>",
}

var Messages_zh_CN = &Messages{
	ModuleDescription:    "SEO 设置：页面的标题、描述和 Open Graph。",
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
	AddVariable:          "变量",

	HelpTooltip:          "SEO 工作原理",
	HelpTitle:            "SEO 工作原理",
	HelpIntroTitle:       "SEO 字段",
	HelpIntro:            "这些字段构成页面的元数据 —— 浏览器标签页和搜索结果中显示的<b>标题</b>、<b>描述</b>和<b>关键词</b>，以及分享页面时社交网络读取的 <b>Open Graph</b> 字段。打开<b>自定义</b>为此记录设置独立的 SEO；留空的字段将被继承（见下文）。",
	HelpVariablesTitle:   "变量",
	HelpVariables:        "字段值可以包含以 <code>{{Name}}</code> 形式书写的变量 —— 例如 <code>{{SiteName}}</code> —— 在页面渲染时被替换。点击字段上方的变量标签即可插入。",
	HelpInheritanceTitle: "继承与回退",
	HelpInheritance:      "SEO 是分层的：留空的字段会从链中更通用的 SEO 取值，最终回退到全局 SEO。因此每一层只需设置与上一层不同的部分。",

	SettingVarSiteName: "站点名称",

	CustomVars:      "自定义变量",
	CustomVarsGroup: "变量",

	SEOConfig:          "SEO 配置",
	MapsKey:            "Google Maps API 密钥",
	MapsKeyHint:        "浏览器 API 密钥（Maps JavaScript + Places），供邮编变量使用。",
	MapsKeyHelpTooltip: "如何获取 API 密钥",
	MapsKeyHelpTitle:   "如何生成 Google Maps API 密钥",
	MapsKeyHelp: "<p>邮编变量会打开 Google 地图来选择地点。要使地图工作，请提供 Google API 密钥（Maps JavaScript API + Places API）：</p>" +
		"<ol>" +
		"<li>访问 console.cloud.google.com 并登录。</li>" +
		"<li>创建（或选择）一个项目。</li>" +
		"<li>在“API 和服务”→“库”中，启用“Maps JavaScript API”和“Places API”。</li>" +
		"<li>在“API 和服务”→“凭据”中，点击“创建凭据”→“API 密钥”。</li>" +
		"<li>限制密钥：“应用限制”→“HTTP 引荐来源网址”填写站点域名；“API 限制”仅限上述两个 API。</li>" +
		"<li>为项目启用结算（Google 要求，但有每月免费额度）。</li>" +
		"<li>复制密钥并粘贴到“Google Maps API 密钥”字段。</li>" +
		"</ol>" +
		"<p>没有密钥时，邮编变量仍可作为纯文本使用。</p>",
}

var Messages_pt_BR = &Messages{
	ModuleDescription:    "As configurações de SEO: títulos, descrições e o Open Graph das páginas.",
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
	AddVariable:          "Variável",

	HelpTooltip:          "Como o SEO funciona",
	HelpTitle:            "Como o SEO funciona",
	HelpIntroTitle:       "Os campos de SEO",
	HelpIntro:            "Estes campos viram os metadados da página — o <b>Título</b> mostrado na aba do navegador e nos resultados de busca, a <b>Descrição</b> e as <b>Palavras-chave</b>, e os campos de <b>Open Graph</b> que as redes sociais leem quando a página é compartilhada. Ative <b>Personalizar</b> para dar a este registro o SEO próprio; deixe um campo vazio para herdá-lo (veja abaixo).",
	HelpVariablesTitle:   "Variáveis",
	HelpVariables:        "O valor de um campo pode conter variáveis escritas como <code>{{Nome}}</code> — por exemplo <code>{{SiteName}}</code> — substituídas quando a página é renderizada. Clique num chip de variável acima do campo para inseri-la.",
	HelpInheritanceTitle: "Herança e fallback",
	HelpInheritance:      "O SEO é em camadas: um campo vazio assume o valor do SEO mais geral na cadeia, terminando no SEO Global. Assim, em cada nível você define apenas o que deve diferir do nível acima.",

	SettingVarSiteName: "Nome do Site",

	CustomVars:      "Variáveis personalizadas",
	CustomVarsGroup: "Variáveis",

	SEOConfig:          "Configurações",
	MapsKey:            "Google Maps API Key",
	MapsKeyHint:        "Chave de navegador (Maps JavaScript + Places), usada pela variável de CEPs.",
	MapsKeyHelpTooltip: "Como gerar a API Key",
	MapsKeyHelpTitle:   "Como gerar a Google Maps API Key",
	MapsKeyHelp: "<p>A variável de CEPs abre o Google Maps para selecionar lugares. Para o mapa funcionar, informe uma chave de API do Google (Maps JavaScript API + Places API):</p>" +
		"<ol>" +
		"<li>Acesse console.cloud.google.com e faça login.</li>" +
		"<li>Crie (ou selecione) um projeto.</li>" +
		"<li>Em \"APIs e serviços\" → \"Biblioteca\", ative \"Maps JavaScript API\" e \"Places API\".</li>" +
		"<li>Em \"APIs e serviços\" → \"Credenciais\", clique \"Criar credenciais\" → \"Chave de API\".</li>" +
		"<li>Restrinja a chave: \"Restrições de aplicativo\" → \"Referenciadores HTTP\" com o domínio do site; \"Restrições de API\" limitada às duas APIs acima.</li>" +
		"<li>Ative o faturamento no projeto (o Google exige, mas há cota gratuita mensal).</li>" +
		"<li>Copie a chave e cole no campo \"Google Maps API Key\".</li>" +
		"</ol>" +
		"<p>Sem a chave, a variável de CEPs continua funcionando como texto simples.</p>",
}
