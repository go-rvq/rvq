package media

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
)

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, I18nMediaLibraryKey, Messages_en_US).(*Messages)
}

type Messages struct {
	ModuleDescription           string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Crop                        string `i18n:"hint='Button that crops an image.'"`
	CropImage                   string `i18n:"hint='Title of the dialog that crops an image.'"`
	ChooseFile                  string `i18n:"hint='Button that opens the choice of a file of the media library.'"`
	Delete                      string `i18n:"hint='Button that removes a file from a field.'"`
	CopyLink                    string `i18n:"hint='Button that copies the address of a file.'"`
	LinkCopied                  string `i18n:"hint='Shown once the address of a file was copied.'"`
	Download                    string `i18n:"hint='Button that downloads a file.'"`
	ChooseAFile                 string `i18n:"label='Choose a file', hint='Title of the dialog that chooses a file of the media library.'"`
	ChosenFileIsGone            string `i18n:"hint='Shown when the chosen file was removed from the media library meanwhile.'"`
	FileInUse                   string `i18n:"hint='Shown when a file asked to be deleted is in use: where it is used follows.', fields=(;'%d'='how many uses')"`
	MediaBoxUsage               string `i18n:"hint='Name, in the singular, of the sizes of the images of a field.'"`
	MediaBoxUsages              string `i18n:"hint='Name, in the plural, of the sizes of the images of the fields (menu, listing title).'"`
	MediaBoxUsageModel          string `i18n:"hint='Label of the model of a field of images.'"`
	MediaBoxUsageLabel          string `i18n:"hint='Label of the name of a field of images.'"`
	MediaBoxUsageSizes          string `i18n:"hint='Label of the sizes of the crops of a field of images.'"`
	MediaBoxUsageRecropState    string `i18n:"hint='Label of how the crops of a field are being made again.'"`
	MediaBoxUsageRecropError    string `i18n:"hint='Label of what went wrong making the crops of a field again.'"`
	SizesOfTheAdmin             string `i18n:"hint='Said of sizes set in the admin, not the code\\'s.'"`
	Recrop                      string `i18n:"hint='Action that makes the crops of the images of a field again.'"`
	RecropStarted               string `i18n:"hint='Shown once the crops of a field started being made again.'"`
	RecropBusy                  string `i18n:"hint='Shown when the crops of a field are already being made.'"`
	RecropRunning               string `i18n:"hint='The crops of a field being made.', fields=(;'%d'='done, of the total that follows')"`
	RecropMade                  string `i18n:"hint='The crops of a field made.', fields=(;'%d'='how many files')"`
	RecropFailedText            string `i18n:"hint='The crops of a field could not all be made.'"`
	SizeNameRequired            string `i18n:"hint='Error: a size of a field with no name.'"`
	SizeNameTwice               string `i18n:"hint='Error: two sizes of a field of the same name.', fields=(;'%s'='the name')"`
	SizeNotPositive             string `i18n:"hint='Error: a size with no width or height.', fields=(;'%s'='its name')"`
	SizesConflict               string `i18n:"hint='Error: a size has the name of a size of another field in other measures.', fields=(;'%s'='the other field and its size')"`
	RecommendedOriginalSize     string `i18n:"hint='The best size for the original of an image.', fields=(;'%[1]d'='its width, in pixels', '%[2]d'='its height, in pixels', '%[3]s'='its proportion (16:9)')"`
	Search                      string `i18n:"hint='Placeholder of the search of the media library.'"`
	UploadFiles                 string `i18n:"hint='Button that uploads files to the media library.'"`
	TakePhoto                   string `i18n:"hint='Button that takes a photo with the camera into the media library.'"`
	Cropping                    string `i18n:"hint='Shown while an image is being cropped.'"`
	DescriptionUpdated          string `i18n:"hint='Shown once the description of a file was saved.'"`
	DescriptionForAccessibility string `i18n:"hint='Placeholder of a file\\'s description, read by screen readers.'"`
	OrderBy                     string `i18n:"hint='Label of the order of the media library.'"`
	FileType                    string `i18n:"hint='Label of the filter of the type of the files: images, videos, files.'"`
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
	MediaLibraryDeletedByID     string `i18n:"hint='Label of the deleted by field of a media file.'"`
	MediaLibraryDeletedOrigin   string `i18n:"hint='Label of the deleted from field of a media file.'"`
	MediaLibraryHidden          string `i18n:"hint='Label of the hidden field of a media file.'"`
	MediaLibrarySelectedType    string `i18n:"hint='Label of the selected type field of a media file.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription:           "The media library: uploading, choosing and cropping files and images.",
	Crop:                        "Crop",
	CropImage:                   "Crop Image",
	ChooseFile:                  "Choose File",
	Delete:                      "Delete",
	ChooseAFile:                 "Choose a File",
	ChosenFileIsGone:            "That file is no longer in the media library. Close this dialog and open it again.",
	FileInUse:                   "This file is in use (%d) and was not deleted. Take it out of:",
	MediaBoxUsage:               "Image sizes of a field",
	MediaBoxUsages:              "Image sizes",
	MediaBoxUsageModel:          "Model",
	MediaBoxUsageLabel:          "Field",
	MediaBoxUsageSizes:          "Sizes",
	MediaBoxUsageRecropState:    "Crops",
	MediaBoxUsageRecropError:    "What went wrong",
	SizesOfTheAdmin:             "the admin's",
	Recrop:                      "Make the crops again",
	RecropStarted:               "The crops are being made again.",
	RecropBusy:                  "The crops of this field are being made already.",
	RecropRunning:               "being made: %d of %d",
	RecropMade:                  "made (%d files)",
	RecropFailedText:            "some could not be made",
	SizeNameRequired:            "A size has no name.",
	SizeNameTwice:               "Two sizes are named %s.",
	SizeNotPositive:             "The size %s has no width or no height.",
	SizesConflict:               "Another field has a size of this name in other measures (%s): a file used in both would be cropped in one, then the other.",
	RecommendedOriginalSize:     "Recommended original: at least %d×%d px, %s.",
	CopyLink:                    "Copy Link",
	LinkCopied:                  "Link Copied!",
	Download:                    "Download",
	Search:                      "Search",
	UploadFiles:                 "Upload files",
	TakePhoto:                   "Take a photo",
	Cropping:                    "Cropping",
	DescriptionUpdated:          "Description Updated",
	DescriptionForAccessibility: "description for accessibility",
	OrderBy:                     "Order By",
	FileType:                    "Type",
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
	MediaLibraryDeletedByID:     "Deleted by",
	MediaLibraryDeletedOrigin:   "Deleted from",
	MediaLibraryHidden:          "Hidden",
	MediaLibrarySelectedType:    "Selected type",
}

var Messages_zh_CN = &Messages{
	ModuleDescription:           "媒体库：上传、选择和裁剪文件与图片。",
	Crop:                        "剪裁",
	CropImage:                   "剪裁图片",
	ChooseFile:                  "选择文件",
	Delete:                      "删除",
	ChooseAFile:                 "选择一个文件",
	ChosenFileIsGone:            "该文件已不在媒体库中。请关闭此对话框后重新打开。",
	FileInUse:                   "该文件正在使用中（%d 处），未删除。请先从以下位置移除：",
	MediaBoxUsage:               "字段的图片尺寸",
	MediaBoxUsages:              "图片尺寸",
	MediaBoxUsageModel:          "模型",
	MediaBoxUsageLabel:          "字段",
	MediaBoxUsageSizes:          "尺寸",
	MediaBoxUsageRecropState:    "裁剪",
	MediaBoxUsageRecropError:    "出错原因",
	SizesOfTheAdmin:             "管理员设置",
	Recrop:                      "重新裁剪",
	RecropStarted:               "正在重新裁剪。",
	RecropBusy:                  "该字段正在裁剪中。",
	RecropRunning:               "处理中：%d / %d",
	RecropMade:                  "已完成（%d 个文件）",
	RecropFailedText:            "部分未能完成",
	SizeNameRequired:            "有尺寸没有名称。",
	SizeNameTwice:               "有两个尺寸都叫 %s。",
	SizeNotPositive:             "尺寸 %s 没有宽度或高度。",
	SizesConflict:               "另一个字段有同名但不同尺寸的设置（%s）：两处都使用的文件会被来回裁剪。",
	RecommendedOriginalSize:     "建议原图：至少 %d×%d 像素，比例 %s。",
	Search:                      "搜索",
	UploadFiles:                 "上传多个文件",
	TakePhoto:                   "拍照",
	Cropping:                    "正在剪裁...",
	DescriptionUpdated:          "描述更新成功",
	DescriptionForAccessibility: "图片描述",
	OrderBy:                     "排序",
	FileType:                    "类型",
	UploadedAt:                  "上传时间",
	UploadedAtDESC:              "上传时间 (降序)",
	All:                         "全部",
	Images:                      "图片",
	Videos:                      "视频",
	Files:                       "文件",
}

var Messages_ja_JP = &Messages{
	ModuleDescription:           "メディアライブラリ：ファイルと画像のアップロード、選択、トリミング。",
	Crop:                        "トリミング",
	CropImage:                   "画像をトリミング",
	ChooseFile:                  "ファイルを選択",
	Delete:                      "削除",
	ChooseAFile:                 "ファイルを選択",
	ChosenFileIsGone:            "そのファイルはメディアライブラリにありません。ダイアログを閉じて開き直してください。",
	FileInUse:                   "このファイルは使用中（%d 件）のため削除されませんでした。次の場所から外してください：",
	MediaBoxUsage:               "フィールドの画像サイズ",
	MediaBoxUsages:              "画像サイズ",
	MediaBoxUsageModel:          "モデル",
	MediaBoxUsageLabel:          "フィールド",
	MediaBoxUsageSizes:          "サイズ",
	MediaBoxUsageRecropState:    "トリミング",
	MediaBoxUsageRecropError:    "エラーの内容",
	SizesOfTheAdmin:             "管理画面の設定",
	Recrop:                      "トリミングをやり直す",
	RecropStarted:               "トリミングをやり直しています。",
	RecropBusy:                  "このフィールドはすでにトリミング中です。",
	RecropRunning:               "処理中: %d / %d",
	RecropMade:                  "完了（%d ファイル）",
	RecropFailedText:            "一部を作成できませんでした",
	SizeNameRequired:            "名前のないサイズがあります。",
	SizeNameTwice:               "%s という名前のサイズが 2 つあります。",
	SizeNotPositive:             "サイズ %s に幅または高さがありません。",
	SizesConflict:               "別のフィールドに同じ名前で寸法の違うサイズがあります（%s）: 両方で使うファイルが交互にトリミングされます。",
	RecommendedOriginalSize:     "推奨する元画像: %d×%d px 以上、比率 %s。",
	Search:                      "検索",
	UploadFiles:                 "ファイルをアップロード",
	TakePhoto:                   "写真を撮る",
	Cropping:                    "トリミング中",
	DescriptionUpdated:          "説明を更新しました",
	DescriptionForAccessibility: "画像の説明",
	OrderBy:                     "並び替え",
	FileType:                    "種類",
	UploadedAt:                  "アップロード日時",
	UploadedAtDESC:              "アップロード日時 (降順)",
	All:                         "すべて",
	Images:                      "画像",
	Videos:                      "動画",
	Files:                       "ファイル",
}

var Messages_pt_BR = &Messages{
	ModuleDescription:           "A biblioteca de mídia: envio, escolha e recorte de arquivos e imagens.",
	Crop:                        "Recortar",
	CropImage:                   "Recortar imagem",
	ChooseFile:                  "Escolher arquivo",
	Delete:                      "Excluir",
	ChooseAFile:                 "Escolha um arquivo",
	ChosenFileIsGone:            "Este arquivo não está mais na biblioteca de mídia. Feche esta janela e abra de novo.",
	FileInUse:                   "Este arquivo está em uso (%d) e não foi excluído. Tire-o de:",
	MediaBoxUsage:               "Tamanhos das imagens de um campo",
	MediaBoxUsages:              "Tamanhos das imagens",
	MediaBoxUsageModel:          "Modelo",
	MediaBoxUsageLabel:          "Campo",
	MediaBoxUsageSizes:          "Tamanhos",
	MediaBoxUsageRecropState:    "Recortes",
	MediaBoxUsageRecropError:    "O que deu errado",
	SizesOfTheAdmin:             "do admin",
	Recrop:                      "Refazer os recortes",
	RecropStarted:               "Os recortes estão sendo refeitos.",
	RecropBusy:                  "Os recortes deste campo já estão sendo refeitos.",
	RecropRunning:               "em andamento: %d de %d",
	RecropMade:                  "prontos (%d arquivos)",
	RecropFailedText:            "alguns não puderam ser feitos",
	SizeNameRequired:            "Um tamanho está sem nome.",
	SizeNameTwice:               "Dois tamanhos se chamam %s.",
	SizeNotPositive:             "O tamanho %s está sem largura ou altura.",
	SizesConflict:               "Outro campo tem um tamanho com este nome em outras medidas (%s): um arquivo usado nos dois seria recortado em um, depois no outro.",
	RecommendedOriginalSize:     "Original recomendada: pelo menos %d×%d px, proporção %s.",
	CopyLink:                    "Copiar link",
	LinkCopied:                  "Link copiado!",
	Download:                    "Baixar",
	Search:                      "Pesquisar",
	UploadFiles:                 "Enviar arquivos",
	TakePhoto:                   "Tirar foto",
	Cropping:                    "Recortando",
	DescriptionUpdated:          "Descrição atualizada",
	DescriptionForAccessibility: "descrição para acessibilidade",
	OrderBy:                     "Ordenar por",
	FileType:                    "Tipo",
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
	MediaLibraryDeletedByID:     "Excluído por",
	MediaLibraryDeletedOrigin:   "Excluído de",
	MediaLibraryHidden:          "Oculto",
	MediaLibrarySelectedType:    "Tipo selecionado",
}
