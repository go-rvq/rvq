package microsite

type Messages struct {
	CurrentPackage string `i18n:"hint='Label of the package a microsite currently serves.'"`
}

var Messages_en_US = &Messages{
	CurrentPackage: "Current Package",
}

var Messages_zh_CN = &Messages{
	CurrentPackage: "当前压缩包",
}

var Messages_pt_BR = &Messages{
	CurrentPackage: "Pacote atual",
}
