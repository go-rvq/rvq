package media

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
)

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, I18nMediaLibraryKey, Messages_en_US).(*Messages)
}

type Messages struct {
	Crop                        string `i18n:"hint='Button that crops an image.'"`
	CropImage                   string `i18n:"hint='Title of the dialog that crops an image.'"`
	ChooseFile                  string `i18n:"hint='Button that opens the choice of a file of the media library.'"`
	Delete                      string `i18n:"hint='Button that removes a file from a field.'"`
	CopyLink                    string `i18n:"hint='Button that copies the address of a file.'"`
	LinkCopied                  string `i18n:"hint='Shown once the address of a file was copied.'"`
	Download                    string `i18n:"hint='Button that downloads a file.'"`
	ChooseAFile                 string `i18n:"label='Choose a file', hint='Title of the dialog that chooses a file of the media library.'"`
	ChosenFileIsGone            string `i18n:"hint='Shown when the chosen file was removed from the media library meanwhile.'"`
	RecommendedOriginalSize     string `i18n:"hint='The best size for the original of an image. %d and %d are its width and height in pixels, %s its proportion (16:9).'"`
	Search                      string `i18n:"hint='Placeholder of the search of the media library.'"`
	UploadFiles                 string `i18n:"hint='Button that uploads files to the media library.'"`
	Cropping                    string `i18n:"hint='Shown while an image is being cropped.'"`
	DescriptionUpdated          string `i18n:"hint='Shown once the description of a file was saved.'"`
	DescriptionForAccessibility string `i18n:"hint='Placeholder of a file\\'s description, read by screen readers.'"`
	OrderBy                     string `i18n:"hint='Label of the order of the media library.'"`
	UploadedAt                  string `i18n:"label='Date uploaded', hint='Order of the media library by upload date, oldest first.'"`
	UploadedAtDESC              string `i18n:"label='Date uploaded (newest first)', hint='Order of the media library by upload date, newest first.'"`
	All                         string `i18n:"hint='Filter of the media library that shows every file.'"`
	Images                      string `i18n:"hint='Filter of the media library that shows the images.'"`
	Videos                      string `i18n:"hint='Filter of the media library that shows the videos.'"`
	Files                       string `i18n:"hint='Filter of the media library that shows the files that are not images or videos.'"`
	MediaLibrary                string `i18n:"hint='Name of the media library model in the singular.'"`
	MediaLibraries              string `i18n:"hint='Name of the media library model in the plural (menu, listing title).'"`
	ShowHidden                  string `i18n:"hint='Filter of the hidden files of the media library.'"`
	OnlyHidden                  string `i18n:"hint='Option of the hidden filter: only the hidden files.'"`
	IncludeHidden               string `i18n:"hint='Option of the hidden filter: the hidden files too.'"`
	NotHidden                   string `i18n:"hint='Option of the hidden filter: only the files that are not hidden.'"`
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
