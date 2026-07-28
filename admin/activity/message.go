package activity

type Messages struct {
	ActivityLogs string
	Activities   string
	ActionAll    string
	ActionView   string
	ActionEdit   string
	ActionCreate string
	ActionDelete string

	ModelUserID    string
	ModelCreatedAt string
	ModelAction    string
	ModelCreator   string
	ModelKeys      string
	ModelName      string
	ModelLabel     string
	ModelLink      string
	ModelDiffs     string
	ModelIP        string
	ModelUserAgent string

	LogAction string
	LogEmpty  string

	FilterAction    string
	FilterCreatedAt string
	FilterCreator   string
	FilterModel     string

	DiffDetail  string
	DiffNew     string
	DiffDelete  string
	DiffChanges string
	DiffField   string
	DiffOld     string
	DiffNow     string
	DiffValue   string
}

var Messages_en_US = &Messages{
	ActivityLogs: "Activity Logs",
	Activities:   "Activities",
	ActionAll:    "All",
	ActionView:   "View",
	ActionEdit:   "Edit",
	ActionCreate: "Create",
	ActionDelete: "Delete",

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
	Activities:   "活动",
	ActionAll:    "全部",
	ActionView:   "查看",
	ActionEdit:   "编辑",
	ActionCreate: "创建",
	ActionDelete: "删除",

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
	ActivityLogs: "Registros de atividade",
	Activities:   "Atividades",
	ActionAll:    "Todas",
	ActionView:   "Consulta",
	ActionEdit:   "Alteração",
	ActionCreate: "Criação",
	ActionDelete: "Exclusão",

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
