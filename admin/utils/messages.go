package utils

type Messages struct {
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	OK                string `i18n:"hint='Button that confirms a confirmation dialog.'"`
	Cancel            string `i18n:"hint='Button that closes a confirmation dialog without doing anything.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription: "Words of the confirmation dialogs.",
	OK:                "OK",
	Cancel:            "Cancel",
}

var Messages_zh_CN = &Messages{
	ModuleDescription: "确认对话框的文字。",
	OK:                "确定",
	Cancel:            "取消",
}

var Messages_ja_JP = &Messages{
	ModuleDescription: "確認ダイアログの文言。",
	OK:                "OK",
	Cancel:            "キャンセル",
}

var Messages_pt_BR = &Messages{
	ModuleDescription: "Palavras dos diálogos de confirmação.",
	OK:                "OK",
	Cancel:            "Cancelar",
}
