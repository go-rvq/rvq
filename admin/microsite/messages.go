package microsite

type Messages struct {
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	CurrentPackage    string `i18n:"hint='Label of the package a microsite currently serves.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription: "The microsites: sites uploaded as a package of files.",
	CurrentPackage:    "Current Package",
}

var Messages_zh_CN = &Messages{
	ModuleDescription: "微型网站：以文件包上传的网站。",
	CurrentPackage:    "当前压缩包",
}

var Messages_pt_BR = &Messages{
	ModuleDescription: "Os microsites: sites enviados como um pacote de arquivos.",
	CurrentPackage:    "Pacote atual",
}
