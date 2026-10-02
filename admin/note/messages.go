package note

type Messages struct {
	SuccessfullyCreated string `i18n:"hint='Shown once a note was added to a record.'"`
	Item                string `i18n:"hint='Label of the note\\'s text field.'"`
	Notes               string `i18n:"hint='Title of a record\\'s notes tab and column.'"`
	NewNote             string `i18n:"hint='Button that adds a note to a record.'"`
}

var Messages_en_US = &Messages{
	SuccessfullyCreated: "Successfully Created",
	Item:                "Item",
	Notes:               "Notes",
	NewNote:             "New Note",
}

var Messages_zh_CN = &Messages{
	SuccessfullyCreated: "成功创建",
	Item:                "记录",
	Notes:               "备注",
	NewNote:             "新建备注",
}

var Messages_ja_JP = &Messages{
	SuccessfullyCreated: "作成に成功しました",
	Item:                "アイテム",
	Notes:               "ノート",
	NewNote:             "新規ノート",
}

var Messages_pt_BR = &Messages{
	SuccessfullyCreated: "Criado com sucesso",
	Item:                "Item",
	Notes:               "Anotações",
	NewNote:             "Nova anotação",
}
