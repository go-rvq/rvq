package activity

type Messages struct {
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	ActivityLogs      string `i18n:"hint='Name of the activity log model in the plural (menu, listing title).'"`
	Activities        string `i18n:"hint='Title of a record\\'s activity tab.'"`
	ActionAll         string `i18n:"label='Action: all', hint='Filter option for every kind of action.'"`
	ActionView        string `i18n:"label='Action: view', hint='Kind of action: a record was seen.'"`
	ActionEdit        string `i18n:"label='Action: edit', hint='Kind of action: a record was changed.'"`
	ActionCreate      string `i18n:"label='Action: create', hint='Kind of action: a record was created.'"`
	ActionDelete      string `i18n:"label='Action: delete', hint='Kind of action: a record was deleted.'"`

	ModelUserID    string `i18n:"label='Creator ID', hint='Label of the id of who did an action.'"`
	ModelCreatedAt string `i18n:"label='Date time', hint='Label of when an action was done.'"`
	ModelAction    string `i18n:"label='Action', hint='Label of the kind of an action.'"`
	ModelCreator   string `i18n:"label='Creator', hint='Label of who did an action.'"`
	ModelKeys      string `i18n:"label='Keys', hint='Label of the keys of the record an action was done on.'"`
	ModelName      string `i18n:"label='Table name', hint='Label of the table of the record an action was done on.'"`
	ModelLabel     string `i18n:"label='Menu name', hint='Label of the model of the record an action was done on, as the menu names it.'"`
	ModelLink      string `i18n:"label='Link', hint='Label of the link to the record an action was done on.'"`
	ModelDiffs     string `i18n:"label='Diffs', hint='Label of what an action changed.'"`
	ModelIP        string `i18n:"label='IP address', hint='Label of the address an action came from.'"`
	ModelUserAgent string `i18n:"label='Browser', hint='Label of the browser an action came from.'"`

	LogAction string `i18n:"label='Activity log', hint='Action that shows the activity log of a record.'"`
	LogEmpty  string `i18n:"hint='Shown when no activity was recorded for a record.'"`

	FilterAction    string `i18n:"label='Filter: action', hint='Filter of the activity log by kind of action.'"`
	FilterCreatedAt string `i18n:"label='Filter: date', hint='Filter of the activity log by date.'"`
	FilterCreator   string `i18n:"label='Filter: creator', hint='Filter of the activity log by who did the action.'"`
	FilterModel     string `i18n:"label='Filter: model', hint='Filter of the activity log by model.'"`

	DiffDetail  string `i18n:"label='Diff: detail', hint='Title of the detail of a change.'"`
	DiffNew     string `i18n:"label='Diff: new', hint='Title of the values of a created record.'"`
	DiffDelete  string `i18n:"label='Diff: delete', hint='Title of the values of a deleted record.'"`
	DiffChanges string `i18n:"label='Diff: changes', hint='Title of the values a change changed.'"`
	DiffField   string `i18n:"label='Diff: field', hint='Column with the field that changed.'"`
	DiffOld     string `i18n:"label='Diff: old', hint='Column with the value before the change.'"`
	DiffNow     string `i18n:"label='Diff: now', hint='Column with the value after the change.'"`
	DiffValue   string `i18n:"label='Diff: value', hint='Column with a value of a created or deleted record.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription: "The activity log: who did what, and when.",
	ActivityLogs:      "Activity Logs",
	Activities:        "Activities",
	ActionAll:         "All",
	ActionView:        "View",
	ActionEdit:        "Edit",
	ActionCreate:      "Create",
	ActionDelete:      "Delete",

	ModelUserID:    "Creator ID",
	ModelCreatedAt: "Date Time",
	ModelAction:    "Action",
	ModelCreator:   "Creator",
	ModelKeys:      "Keys",
	ModelName:      "Table Name",
	ModelLabel:     "Menu Name",
	ModelLink:      "Link",
	ModelDiffs:     "Diffs",
	ModelIP:        "IP Address",
	ModelUserAgent: "Browser",

	LogAction: "Activity Log",
	LogEmpty:  "No activity recorded for this record.",

	FilterAction:    "Action",
	FilterCreatedAt: "Create Time",
	FilterCreator:   "Creator",
	FilterModel:     "Model Name",

	DiffDetail:  "Detail",
	DiffNew:     "New",
	DiffDelete:  "Delete",
	DiffChanges: "Changes",
	DiffField:   "Filed",
	DiffOld:     "Old",
	DiffNow:     "Now",
	DiffValue:   "Value",
}

var Messages_zh_CN = &Messages{
	ModuleDescription: "活动日志：谁在何时做了什么。",
	Activities:        "活动",
	ActionAll:         "全部",
	ActionView:        "查看",
	ActionEdit:        "编辑",
	ActionCreate:      "创建",
	ActionDelete:      "删除",

	ModelUserID:    "操作者ID",
	ModelCreatedAt: "日期时间",
	ModelAction:    "操作",
	ModelCreator:   "操作者",
	ModelKeys:      "表的主键值",
	ModelName:      "表名",
	ModelLabel:     "菜单名",
	ModelLink:      "链接",
	ModelDiffs:     "差异",
	ModelIP:        "IP 地址",
	ModelUserAgent: "浏览器",

	LogAction: "活动日志",
	LogEmpty:  "该记录没有活动记录。",

	FilterAction:    "操作类型",
	FilterCreatedAt: "操作时间",
	FilterCreator:   "操作人",
	FilterModel:     "操作对象",
	DiffDetail:      "详情",
	DiffNew:         "新加",
	DiffDelete:      "删除",
	DiffChanges:     "修改",
	DiffField:       "字段",
	DiffOld:         "之前的值",
	DiffNow:         "当前的值",
	DiffValue:       "值",
}

var Messages_pt_BR = &Messages{
	ModuleDescription: "O registro de atividades: quem fez o quê, e quando.",
	ActivityLogs:      "Registros de atividade",
	Activities:        "Atividades",
	ActionAll:         "Todas",
	ActionView:        "Consulta",
	ActionEdit:        "Alteração",
	ActionCreate:      "Criação",
	ActionDelete:      "Exclusão",

	ModelUserID:    "ID do autor",
	ModelCreatedAt: "Data e hora",
	ModelAction:    "Ação",
	ModelCreator:   "Autor",
	ModelKeys:      "Chaves",
	ModelName:      "Tabela",
	ModelLabel:     "Menu",
	ModelLink:      "Link",
	ModelDiffs:     "Diferenças",
	ModelIP:        "Endereço IP",
	ModelUserAgent: "Navegador",

	LogAction: "Registro de atividade",
	LogEmpty:  "Nenhuma atividade registrada para este registro.",

	FilterAction:    "Ação",
	FilterCreatedAt: "Data",
	FilterCreator:   "Autor",
	FilterModel:     "Modelo",

	DiffDetail:  "Detalhe",
	DiffNew:     "Novo",
	DiffDelete:  "Excluído",
	DiffChanges: "Alterações",
	DiffField:   "Campo",
	DiffOld:     "Antes",
	DiffNow:     "Agora",
	DiffValue:   "Valor",
}
