package slug

type Messages struct {
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Sync              string `i18n:"label='Auto sync', hint='Label of the switch that keeps a slug in step with another field.', fields=(;'%s'='that field\\'s label, in lower case')"`
}

var Messages_en_US = &Messages{
	ModuleDescription: "The slugs: the part of an address made from a field.",
	Sync:              "Auto sync from %s",
}

var Messages_zh_CN = &Messages{
	ModuleDescription: "Slug：由字段生成的地址部分。",
	Sync:              "从%s自动同步",
}

var Messages_pt_BR = &Messages{
	ModuleDescription: "Os slugs: a parte de um endereço gerada a partir de um campo.",
	Sync:              "Sincronizar automaticamente a partir de %s",
}
