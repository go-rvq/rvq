package utils

type Messages struct {
	OK     string `i18n:"hint='Button that confirms a confirmation dialog.'"`
	Cancel string `i18n:"hint='Button that closes a confirmation dialog without doing anything.'"`
}

var Messages_en_US = &Messages{
	OK:     "OK",
	Cancel: "Cancel",
}

var Messages_zh_CN = &Messages{
	OK:     "确定",
	Cancel: "取消",
}

var Messages_ja_JP = &Messages{
	OK:     "OK",
	Cancel: "キャンセル",
}

var Messages_pt_BR = &Messages{
	OK:     "OK",
	Cancel: "Cancelar",
}
