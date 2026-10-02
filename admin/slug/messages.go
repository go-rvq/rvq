package slug

type Messages struct {
	Sync string `i18n:"label='Auto sync', hint='Label of the switch that keeps a slug in step with another field.', fields=(;'%s'='that field\\'s label, in lower case')"`
}

var Messages_en_US = &Messages{
	Sync: "Auto sync from %s",
}

var Messages_zh_CN = &Messages{
	Sync: "从%s自动同步",
}

var Messages_pt_BR = &Messages{
	Sync: "Sincronizar automaticamente a partir de %s",
}
