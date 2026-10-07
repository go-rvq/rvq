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
	SEOConfig                 string `i18n:"label='SEO settings', hint='Name of the SEO settings model.'"`
	MapsKey                   string `i18n:"label='Google Maps API key', hint='Label of the Google Maps key used by the ZIP-codes variable.'"`
	MapsKeyHint               string `i18n:"label='Google Maps API key: hint', hint='Hint of the Google Maps key: which APIs it needs.'"`
	MapsKeyHelpTooltip        string `i18n:"label='Google Maps API key: help tooltip', hint='Tooltip of the button that opens the help on the Google Maps key.'"`
	MapsKeyHelpTitle          string `i18n:"label='Google Maps API key: help title', hint='Title of the help on how to get a Google Maps key.'"`
	MapsKeyHelp               string `i18n:"type=html, label='Google Maps API key: help', hint='Steps to get a Google Maps key (HTML).'"`
	MapsKeyTest               string `i18n:"label='Google Maps key: test', hint='Button that tests the Google Maps key: the map, the search of places, of addresses.'"`
	MapsKeyTestTitle          string `i18n:"label='Google Maps key: test title', hint='Title of the dialog of the test of the Google Maps key.'"`
	MapsKeyTestNoKey          string `i18n:"label='Google Maps key: test with no key', hint='Said when the key is tested with no key typed.'"`
	MapsKeyTestTry            string `i18n:"label='Google Maps key: test by hand', hint='Above the address field of the test: search an address to see it on the map.'"`
	MapsKeyTestReload         string `i18n:"label='Google Maps key: test of another key', hint='Said when the page loaded Google Maps with another key: save and reload the page.'"`
	MapsKeyTestRejected       string `i18n:"label='Google Maps key: refused', hint='Said when Google refuses the key.', fields=(;'{code}'='the code of Google\\'s error')"`
	MapsKeyTestOK             string `i18n:"label='Google Maps key: test passed', hint='Said of a part of the test that works.'"`
	MapsKeyTestFailed         string `i18n:"label='Google Maps key: test failed', hint='Said of a part of the test that fails.'"`
	MapsKeyTestBilling        string `i18n:"label='Google Maps key: no billing', hint='Said when the project of the key has no billing account (BillingNotEnabledMapError).'"`
	MapsKeyTestNotEnabled     string `i18n:"label='Google Maps key: API not enabled', hint='Said when an API is not enabled for the key (ApiNotActivatedMapError, REQUEST_DENIED).'"`
	MapsKeyTestReferer        string `i18n:"label='Google Maps key: site not allowed', hint='Said when the key does not allow the site (RefererNotAllowedMapError).'"`
	MapsKeyTestInvalid        string `i18n:"label='Google Maps key: invalid', hint='Said when the key is not valid (InvalidKeyMapError).'"`
	SEOConfigGoogleMapsAPIKey string `i18n:"hint='Label of the Google Maps API key field of the SEO settings.'"`
	SEOConfigYAML             string `i18n:"hint='Label of the YAML field of the SEO settings.'"`
	SEOSetting                string `i18n:"hint='Name of the SEO setting of a record.'"`
	SEOVariables              string `i18n:"hint='Label of the variables field of the SEO settings.'"`
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

	SEOConfig:             "Settings",
	MapsKey:               "Google Maps API Key",
	MapsKeyHint:           "A browser API key (Maps JavaScript, Places (New), Geocoding): the address fields and the ZIP-codes variable.",
	MapsKeyHelpTooltip:    "How to get the API key",
	MapsKeyTest:           "Test",
	MapsKeyTestTitle:      "Test of the Google Maps key",
	MapsKeyTestNoKey:      "Type the key first.",
	MapsKeyTestTry:        "Search an address: it shows on the map.",
	MapsKeyTestReload:     "This page already loaded Google Maps with another key: save the settings and reload the page to test this one.",
	MapsKeyTestRejected:   "Google refused the key: {code}.",
	MapsKeyTestOK:         "Working",
	MapsKeyTestFailed:     "Failed",
	MapsKeyTestBilling:    "The key's project has no billing account: link one in the Google Cloud console (Billing).",
	MapsKeyTestNotEnabled: "The API is not enabled in the key's project (APIs & Services → Library), or the key's API restrictions leave it out.",
	MapsKeyTestReferer:    "The key does not allow this site: add its address to the key's website restrictions.",
	MapsKeyTestInvalid:    "The key is not valid: copy it again from the console (Credentials).",
	MapsKeyHelpTitle:      "How to generate the Google Maps API Key",
	MapsKeyHelp: "<p>The key is used by the address fields of the forms (the suggestions of Google and the map below them) and by the ZIP-codes variable. It is a <b>browser key</b> of a project of Google Cloud:</p>" +
		"<ol>" +
		"<li>Go to <a href=\"https://console.cloud.google.com\" target=\"_blank\" rel=\"noopener\">console.cloud.google.com</a> and sign in; create (or select) a project.</li>" +
		"<li><b>Billing</b>: link the project to a billing account (\"Billing\" → \"Link a billing account\"). Google requires it even within the free monthly quota — without it, it refuses the searches (<code>BillingNotEnabledMapError</code>).</li>" +
		"<li><b>APIs</b>: under \"APIs &amp; Services\" → \"Library\", enable <b>Maps JavaScript API</b>, <b>Places API (New)</b> and <b>Geocoding API</b>. The legacy \"Places API\" is not enough: since March 2025 Google does not give its autocomplete to new projects.</li>" +
		"<li><b>The key</b>: under \"APIs &amp; Services\" → \"Credentials\", \"Create credentials\" → \"API key\".</li>" +
		"<li><b>Its restrictions</b> (\"Edit API key\"): \"Application restrictions\" → \"Websites\", with the addresses of the site and of the admin (<code>https://your-site.com/*</code>, <code>https://www.your-site.com/*</code>; to try it on a computer, <code>http://localhost:*/*</code> and <code>http://127.0.0.1:*/*</code>); \"API restrictions\" → \"Restrict key\", the three APIs above.</li>" +
		"<li>Copy the key, paste it in the \"Google Maps API Key\" field and <b>Test</b> it: the map, the search of places and of addresses must say \"Working\". A change in the console takes a few minutes to hold.</li>" +
		"</ol>" +
		"<p>Without the key, an address is only typed — no suggestions, no map — and the ZIP-codes variable is plain text.</p>",
	SEOConfigGoogleMapsAPIKey: "Google Maps API key",
	SEOConfigYAML:             "YAML",
	SEOSetting:                "Setting",
	SEOVariables:              "Variables",
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

	SEOConfig:             "SEO 配置",
	MapsKey:               "Google Maps API 密钥",
	MapsKeyHint:           "浏览器 API 密钥（Maps JavaScript、Places (New)、Geocoding）：地址字段和邮编变量使用。",
	MapsKeyHelpTooltip:    "如何获取 API 密钥",
	MapsKeyTest:           "测试",
	MapsKeyTestTitle:      "测试 Google Maps 密钥",
	MapsKeyTestNoKey:      "请先输入密钥。",
	MapsKeyTestTry:        "搜索一个地址：它会显示在地图上。",
	MapsKeyTestReload:     "此页面已使用另一个密钥加载 Google Maps：请保存设置并重新加载页面以测试此密钥。",
	MapsKeyTestRejected:   "Google 拒绝了该密钥：{code}。",
	MapsKeyTestOK:         "正常",
	MapsKeyTestFailed:     "失败",
	MapsKeyTestBilling:    "该密钥的项目没有结算账号：请在 Google Cloud 控制台中关联一个（结算）。",
	MapsKeyTestNotEnabled: "该 API 未在密钥的项目中启用（API 和服务 → 库），或者密钥的 API 限制不包括它。",
	MapsKeyTestReferer:    "该密钥不允许此网站：请将其地址添加到密钥的网站限制中。",
	MapsKeyTestInvalid:    "该密钥无效：请从控制台（凭据）重新复制。",
	MapsKeyHelpTitle:      "如何生成 Google Maps API 密钥",
	MapsKeyHelp: "<p>该密钥用于表单的地址字段（Google 的建议及其下方的地图）和邮编变量。它是 Google Cloud 项目的<b>浏览器密钥</b>：</p>" +
		"<ol>" +
		"<li>访问 <a href=\"https://console.cloud.google.com\" target=\"_blank\" rel=\"noopener\">console.cloud.google.com</a> 并登录；创建（或选择）一个项目。</li>" +
		"<li><b>结算</b>：将项目关联到结算账号。即使在每月免费额度内，Google 也要求这样做——否则会拒绝搜索（<code>BillingNotEnabledMapError</code>）。</li>" +
		"<li><b>API</b>：在“API 和服务”→“库”中启用 <b>Maps JavaScript API</b>、<b>Places API (New)</b> 和 <b>Geocoding API</b>。旧版“Places API”不够：自 2025 年 3 月起，Google 不再向新项目提供其自动补全。</li>" +
		"<li><b>密钥</b>：在“API 和服务”→“凭据”中，“创建凭据”→“API 密钥”。</li>" +
		"<li><b>限制</b>（“修改 API 密钥”）：“应用限制”→“网站”，填写网站和后台的地址（<code>https://your-site.com/*</code>；在电脑上测试时加 <code>http://localhost:*/*</code> 和 <code>http://127.0.0.1:*/*</code>）；“API 限制”→“限制密钥”，选择上述三个 API。</li>" +
		"<li>复制密钥，粘贴到“Google Maps API Key”字段并<b>测试</b>：地图、地点搜索和地址搜索都应显示“正常”。控制台中的更改需要几分钟才能生效。</li>" +
		"</ol>" +
		"<p>没有密钥时，地址只能手动输入——没有建议，也没有地图——邮编变量为纯文本。</p>",
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

	SEOConfig:             "Configurações",
	MapsKey:               "Google Maps API Key",
	MapsKeyHint:           "Chave de navegador (Maps JavaScript, Places (New), Geocoding): os campos de endereço e a variável de CEPs.",
	MapsKeyHelpTooltip:    "Como gerar a API Key",
	MapsKeyTest:           "Testar",
	MapsKeyTestTitle:      "Teste da chave do Google Maps",
	MapsKeyTestNoKey:      "Digite a chave primeiro.",
	MapsKeyTestTry:        "Busque um endereço: ele aparece no mapa.",
	MapsKeyTestReload:     "Esta página já carregou o Google Maps com outra chave: salve as configurações e recarregue a página para testar esta.",
	MapsKeyTestRejected:   "O Google recusou a chave: {code}.",
	MapsKeyTestOK:         "Funcionando",
	MapsKeyTestFailed:     "Falhou",
	MapsKeyTestBilling:    "O projeto da chave não tem conta de faturamento: vincule uma no console do Google Cloud (Faturamento).",
	MapsKeyTestNotEnabled: "A API não está ativada no projeto da chave (APIs e serviços → Biblioteca), ou as restrições de API da chave a deixam de fora.",
	MapsKeyTestReferer:    "A chave não permite este site: adicione o endereço dele às restrições de sites da chave.",
	MapsKeyTestInvalid:    "A chave não é válida: copie-a de novo do console (Credenciais).",
	MapsKeyHelpTitle:      "Como gerar a Google Maps API Key",
	MapsKeyHelp: "<p>A chave é usada pelos campos de endereço dos formulários (as sugestões do Google e o mapa abaixo deles) e pela variável de CEPs. É uma <b>chave de navegador</b> de um projeto do Google Cloud:</p>" +
		"<ol>" +
		"<li>Acesse <a href=\"https://console.cloud.google.com\" target=\"_blank\" rel=\"noopener\">console.cloud.google.com</a> e faça login; crie (ou selecione) um projeto.</li>" +
		"<li><b>Faturamento</b>: vincule o projeto a uma conta de faturamento (\"Faturamento\" → \"Vincular uma conta de faturamento\"). O Google exige, mesmo dentro da cota gratuita mensal — sem ele, recusa as buscas (<code>BillingNotEnabledMapError</code>).</li>" +
		"<li><b>APIs</b>: em \"APIs e serviços\" → \"Biblioteca\", ative a <b>Maps JavaScript API</b>, a <b>Places API (New)</b> e a <b>Geocoding API</b>. A \"Places API\" antiga não basta: desde março de 2025 o Google não dá o seu autocomplete a projetos novos.</li>" +
		"<li><b>A chave</b>: em \"APIs e serviços\" → \"Credenciais\", \"Criar credenciais\" → \"Chave de API\".</li>" +
		"<li><b>As restrições</b> (\"Editar chave de API\"): \"Restrições de aplicativo\" → \"Sites\", com os endereços do site e do admin (<code>https://seu-site.com/*</code>, <code>https://www.seu-site.com/*</code>; para testar num computador, <code>http://localhost:*/*</code> e <code>http://127.0.0.1:*/*</code>); \"Restrições de API\" → \"Restringir chave\", as três APIs acima.</li>" +
		"<li>Copie a chave, cole no campo \"Google Maps API Key\" e <b>Teste</b>: o mapa, a busca de lugares e a de endereços devem dizer \"Funcionando\". Uma mudança no console leva alguns minutos para valer.</li>" +
		"</ol>" +
		"<p>Sem a chave, um endereço é só digitado — sem sugestões, sem mapa — e a variável de CEPs é texto simples.</p>",
	SEOConfigGoogleMapsAPIKey: "Chave da API do Google Maps",
	SEOConfigYAML:             "YAML",
	SEOSetting:                "Configuração",
	SEOVariables:              "Variáveis",
}
