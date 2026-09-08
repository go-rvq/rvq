package media

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
)

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, I18nMediaLibraryKey, Messages_en_US).(*Messages)
}

type Messages struct {
	Crop                        string
	CropImage                   string
	ChooseFile                  string
	Delete                      string
	CopyLink                    string
	LinkCopied                  string
	Download                    string
	ChooseAFile                 string
	ChosenFileIsGone            string
	RecommendedOriginalSize     string
	Search                      string
	UploadFiles                 string
	Cropping                    string
	DescriptionUpdated          string
	DescriptionForAccessibility string
	OrderBy                     string
	UploadedAt                  string
	UploadedAtDESC              string
	All                         string
	Images                      string
	Videos                      string
	Files                       string
	MediaLibrary                string
	MediaLibraries              string
	ShowHidden                  string
	OnlyHidden                  string
	IncludeHidden               string
	NotHidden                   string
}

var Messages_en_US = &Messages{
	Crop:                        "Crop",
	CropImage:                   "Crop Image",
	ChooseFile:                  "Choose File",
	Delete:                      "Delete",
	ChooseAFile:                 "Choose a File",
	ChosenFileIsGone:            "That file is no longer in the media library. Close this dialog and open it again.",
	RecommendedOriginalSize:     "Recommended original: at least %d×%d px, %s.",
	CopyLink:                    "Copy Link",
	LinkCopied:                  "Link Copied!",
	Download:                    "Download",
	Search:                      "Search",
	UploadFiles:                 "Upload files",
	Cropping:                    "Cropping",
	DescriptionUpdated:          "Description Updated",
	DescriptionForAccessibility: "description for accessibility",
	OrderBy:                     "Order By",
	UploadedAt:                  "Date Uploaded",
	UploadedAtDESC:              "Date Uploaded (DESC)",
	All:                         "All",
	Images:                      "Images",
	Videos:                      "Videos",
	Files:                       "Files",
	ShowHidden:                  "Show hidden",
	OnlyHidden:                  "Only hidden",
	IncludeHidden:               "Include hidden",
	NotHidden:                   "Not hidden",
}

var Messages_zh_CN = &Messages{
	Crop:                        "剪裁",
	CropImage:                   "剪裁图片",
	ChooseFile:                  "选择文件",
	Delete:                      "删除",
	ChooseAFile:                 "选择一个文件",
	ChosenFileIsGone:            "该文件已不在媒体库中。请关闭此对话框后重新打开。",
	RecommendedOriginalSize:     "建议原图：至少 %d×%d 像素，比例 %s。",
	Search:                      "搜索",
	UploadFiles:                 "上传多个文件",
	Cropping:                    "正在剪裁...",
	DescriptionUpdated:          "描述更新成功",
	DescriptionForAccessibility: "图片描述",
	OrderBy:                     "排序",
	UploadedAt:                  "上传时间",
	UploadedAtDESC:              "上传时间 (降序)",
	All:                         "全部",
	Images:                      "图片",
	Videos:                      "视频",
	Files:                       "文件",
}

var Messages_ja_JP = &Messages{
	Crop:                        "トリミング",
	CropImage:                   "画像をトリミング",
	ChooseFile:                  "ファイルを選択",
	Delete:                      "削除",
	ChooseAFile:                 "ファイルを選択",
	ChosenFileIsGone:            "そのファイルはメディアライブラリにありません。ダイアログを閉じて開き直してください。",
	RecommendedOriginalSize:     "推奨する元画像: %d×%d px 以上、比率 %s。",
	Search:                      "検索",
	UploadFiles:                 "ファイルをアップロード",
	Cropping:                    "トリミング中",
	DescriptionUpdated:          "説明を更新しました",
	DescriptionForAccessibility: "画像の説明",
	OrderBy:                     "並び替え",
	UploadedAt:                  "アップロード日時",
	UploadedAtDESC:              "アップロード日時 (降順)",
	All:                         "すべて",
	Images:                      "画像",
	Videos:                      "動画",
	Files:                       "ファイル",
}

var Messages_pt_BR = &Messages{
	Crop:                        "Recortar",
	CropImage:                   "Recortar imagem",
	ChooseFile:                  "Escolher arquivo",
	Delete:                      "Excluir",
	ChooseAFile:                 "Escolha um arquivo",
	ChosenFileIsGone:            "Este arquivo não está mais na biblioteca de mídia. Feche esta janela e abra de novo.",
	RecommendedOriginalSize:     "Original recomendada: pelo menos %d×%d px, proporção %s.",
	CopyLink:                    "Copiar link",
	LinkCopied:                  "Link copiado!",
	Download:                    "Baixar",
	Search:                      "Pesquisar",
	UploadFiles:                 "Enviar arquivos",
	Cropping:                    "Recortando",
	DescriptionUpdated:          "Descrição atualizada",
	DescriptionForAccessibility: "descrição para acessibilidade",
	OrderBy:                     "Ordenar por",
	UploadedAt:                  "Data de envio",
	UploadedAtDESC:              "Data de envio (decrescente)",
	All:                         "Todos",
	Images:                      "Imagens",
	Videos:                      "Vídeos",
	Files:                       "Arquivos",
	MediaLibrary:                "Biblioteca de mídia",
	MediaLibraries:              "Bibliotecas de mídia",
	ShowHidden:                  "Mostrar ocultos",
	OnlyHidden:                  "Somente ocultos",
	IncludeHidden:               "Incluir ocultos",
	NotHidden:                   "Não ocultos",
}
