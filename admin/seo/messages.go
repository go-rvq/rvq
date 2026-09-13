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
	SettingVarSiteName:   "Site Name",
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
	SettingVarSiteName:   "站点名称",
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
	SettingVarSiteName:   "Nome do Site",
}
