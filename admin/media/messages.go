package media

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
)

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, I18nMediaLibraryKey, Messages_en_US).(*Messages)
}

type Messages struct {
	ModuleDescription             string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Crop                          string `i18n:"hint='Button that crops an image.'"`
	CropImage                     string `i18n:"hint='Title of the dialog that crops an image.'"`
	EditImage                     string `i18n:"hint='Title of the editor of an image (crop, colors, annotations), and of the click on a thumbnail that opens it.'"`
	EditBeforeUpload              string `i18n:"hint='Switch of the upload: the images chosen open in the editor before they are sent.'"`
	DropImages                    string `i18n:"hint='Placeholder of the place to drop the images of a list on.'"`
	DropFiles                     string `i18n:"hint='Placeholder of the place to drop the files of a list on.'"`
	AttachmentImages              string `i18n:"hint='Label of the images attached to a record.'"`
	AttachmentFiles               string `i18n:"hint='Label of the files attached to a record.'"`
	AttachmentImage               string `i18n:"hint='Label of an image attached to a record.'"`
	AttachmentName                string `i18n:"hint='Column of the name of a file attached to a record.'"`
	AttachmentDescription         string `i18n:"hint='Column of the description of a file attached to a record.'"`
	AttachmentSize                string `i18n:"hint='Column of the size of a file attached to a record.'"`
	DropFile                      string `i18n:"hint='Placeholder of an empty field of a file: drop it here, or click to choose it.'"`
	DropImage                     string `i18n:"hint='Placeholder of an empty field of an image: drop it here, or click to choose it.'"`
	ImageDescription              string `i18n:"hint='Placeholder of the description of an image of a list.'"`
	ChooseFile                    string `i18n:"hint='Button that opens the choice of a file of the media library.'"`
	Delete                        string `i18n:"hint='Button that removes a file from a field.'"`
	CopyLink                      string `i18n:"hint='Button that copies the address of a file.'"`
	LinkCopied                    string `i18n:"hint='Shown once the address of a file was copied.'"`
	Download                      string `i18n:"hint='Button that downloads a file.'"`
	ChooseAFile                   string `i18n:"label='Choose a file', hint='Title of the dialog that chooses a file of the media library.'"`
	ChosenFileIsGone              string `i18n:"hint='Shown when the chosen file was removed from the media library meanwhile.'"`
	FileInUse                     string `i18n:"hint='Shown when a file asked to be deleted is in use: where it is used follows.', fields=(;'%d'='how many uses')"`
	MediaBoxUsage                 string `i18n:"hint='Name, in the singular, of the sizes of the images of a field.'"`
	MediaBoxUsages                string `i18n:"hint='Name, in the plural, of the sizes of the images of the fields (menu, listing title).'"`
	MediaBoxUsageModel            string `i18n:"hint='Label of the model of a field of images.'"`
	MediaBoxUsageLabel            string `i18n:"hint='Label of the name of a field of images.'"`
	MediaBoxUsageSizes            string `i18n:"hint='Label of the sizes of the crops of a field of images.'"`
	Format                        string `i18n:"hint='Label of the select of the format an image of a field is cropped in.'"`
	FormatDefault                 string `i18n:"hint='The format of a field of images that is its own sizes.'"`
	MediaBoxUsageFieldFormats     string `i18n:"hint='Label of the formats a field of images offers: its sizes in other proportions.'"`
	MediaFormat                   string `i18n:"hint='Name, in the singular, of a format of the images: a proportion with a name.'"`
	MediaFormats                  string `i18n:"hint='Name, in the plural, of the formats of the images (menu, listing title).'"`
	MediaFormatName               string `i18n:"hint='Label of the name of a format: the suffix of its crops.'"`
	MediaFormatName_Hint          string `i18n:"hint='Hint of the name of a format.'"`
	MediaFormatLabel              string `i18n:"hint='Label of the label of a format: what a person reads of it.'"`
	MediaFormatWidth              string `i18n:"hint='Label of the width of the proportion of a format.'"`
	MediaFormatHeight             string `i18n:"hint='Label of the height of the proportion of a format.'"`
	MediaFormatRatio              string `i18n:"hint='Label of the proportion of a format (width:height).'"`
	MediaFormatRatio_Hint         string `i18n:"hint='Hint of the proportion of a format.'"`
	MediaFormatCode               string `i18n:"hint='Label of whether a format is declared by the code.'"`
	MediaFormatUsages             string `i18n:"hint='Label of the fields of images that offer a format.'"`
	MediaBoxUsageFormat           string `i18n:"hint='Name, in the singular, of a format a field of images offers.'"`
	MediaBoxUsageFormats          string `i18n:"hint='Name, in the plural, of the formats the fields of images offer (menu, listing title).'"`
	MediaBoxUsageFormatUsage      string `i18n:"hint='Label of the field of images of a format it offers.'"`
	MediaBoxUsageFormatFormat     string `i18n:"hint='Label of the format a field of images offers.'"`
	MediaBoxUsageFormatPosition   string `i18n:"hint='Label of the place of a format among the ones of its field.'"`
	MediaBoxUsageFormatWidth      string `i18n:"hint='Label of the width of the largest size of a format in its field.'"`
	MediaBoxUsageFormatWidth_Hint string `i18n:"hint='Hint of the width of the largest size of a format in its field.'"`
	MediaBoxUsageFormatSizes      string `i18n:"hint='Label of the sizes of a format in its field.'"`
	FormatNameInvalid             string `i18n:"hint='Error: a format with a name not of lower-case letters, digits and -.'"`
	FormatNameTaken               string `i18n:"hint='Error: a format with the name of another.'"`
	FormatNotPositive             string `i18n:"hint='Error: a format with no width or no height of its proportion.'"`
	FormatInUse                   string `i18n:"hint='Error: a format deleted while fields offer it (the fields).'"`
	FormatOfCode                  string `i18n:"hint='Error: a format of the code deleted.'"`
	FormatUsageRequired           string `i18n:"hint='Error: a format of a field with no field chosen.'"`
	FormatRequired                string `i18n:"hint='Error: a format of a field with no format chosen.'"`
	FormatTwice                   string `i18n:"hint='Error: a field offering a format twice.'"`
	FormatWidthNegative           string `i18n:"hint='Error: a negative width of the largest size of a format in a field.'"`
	MediaBoxUsageRecropState      string `i18n:"hint='Label of how the crops of a field are being made again.'"`
	MediaBoxUsageRecropError      string `i18n:"hint='Label of what went wrong making the crops of a field again.'"`
	MediaBoxUsageOriginal         string `i18n:"hint='Label of the largest an original image of a field is kept in.'"`
	MediaBoxUsageOriginalWidth    string `i18n:"hint='Label of the largest width of an original image of a field.'"`
	MediaBoxUsageOriginalHeight   string `i18n:"hint='Label of the largest height of an original image of a field.'"`
	OriginalOfTheCode             string `i18n:"hint='Hint of the limit of the original: empty, the code\\'s.', fields=(;'%[1]d'='its width', '%[2]d'='its height')"`
	OriginalOfTheCodeNone         string `i18n:"hint='Hint of the limit of the original: empty, the media library\\'s (the code has none).'"`
	OriginalBoth                  string `i18n:"hint='Error: the limit of the original needs both its width and its height, or neither.'"`
	MediaBoxUsageThumb            string `i18n:"hint='Label of the largest the thumbnail of an image of a field is made in.'"`
	MediaBoxUsageThumbWidth       string `i18n:"hint='Label of the largest width of the thumbnail of an image of a field.'"`
	MediaBoxUsageThumbHeight      string `i18n:"hint='Label of the largest height of the thumbnail of an image of a field.'"`
	ThumbOfTheCode                string `i18n:"hint='Hint of the limit of the thumbnail: empty, the code\\'s.', fields=(;'%[1]d'='its width', '%[2]d'='its height')"`
	ThumbOfTheCodeNone            string `i18n:"hint='Hint of the limit of the thumbnail: empty, none (the code has none).'"`
	SizesOfTheAdmin               string `i18n:"hint='Said of sizes set in the admin, not the code\\'s.'"`
	Recrop                        string `i18n:"hint='Action that makes the crops of the images of a field again.'"`
	RecropStarted                 string `i18n:"hint='Shown once the crops of a field started being made again.'"`
	RecropBusy                    string `i18n:"hint='Shown when the crops of a field are already being made.'"`
	RecropRunning                 string `i18n:"hint='The crops of a field being made.', fields=(;'%d'='done, of the total that follows')"`
	RecropMade                    string `i18n:"hint='The crops of a field made.', fields=(;'%d'='how many files')"`
	RecropFailedText              string `i18n:"hint='The crops of a field could not all be made.'"`
	SizeNameRequired              string `i18n:"hint='Error: a size of a field with no name.'"`
	SizeNameTwice                 string `i18n:"hint='Error: two sizes of a field of the same name.', fields=(;'%s'='the name')"`
	SizeNotPositive               string `i18n:"hint='Error: a size with no width or height.', fields=(;'%s'='its name')"`
	SizesConflict                 string `i18n:"hint='Error: a size has the name of a size of another field in other measures.', fields=(;'%s'='the other field and its size')"`
	RecommendedOriginalSize       string `i18n:"hint='The best size for the original of an image.', fields=(;'%[1]d'='its width, in pixels', '%[2]d'='its height, in pixels', '%[3]s'='its proportion (16:9)')"`
	Search                        string `i18n:"hint='Placeholder of the search of the media library.'"`
	UploadFiles                   string `i18n:"hint='Button that uploads files to the media library.'"`
	TakePhoto                     string `i18n:"hint='Button that takes a photo with the camera into the media library.'"`
	Cropping                      string `i18n:"hint='Shown while an image is being cropped.'"`
	DescriptionUpdated            string `i18n:"hint='Shown once the description of a file was saved.'"`
	DescriptionForAccessibility   string `i18n:"hint='Placeholder of a file\\'s description, read by screen readers.'"`
	OrderBy                       string `i18n:"hint='Label of the order of the media library.'"`
	FileType                      string `i18n:"hint='Label of the filter of the type of the files: images, videos, files.'"`
	UploadedAt                    string `i18n:"label='Date uploaded', hint='Order of the media library by upload date, oldest first.'"`
	UploadedAtDESC                string `i18n:"label='Date uploaded (newest first)', hint='Order of the media library by upload date, newest first.'"`
	All                           string `i18n:"hint='Filter of the media library that shows every file.'"`
	Images                        string `i18n:"hint='Filter of the media library that shows the images.'"`
	Videos                        string `i18n:"hint='Filter of the media library that shows the videos.'"`
	Files                         string `i18n:"hint='Filter of the media library that shows the files that are not images or videos.'"`
	MediaLibrary                  string `i18n:"hint='Name of the media library model in the singular.'"`
	MediaLibraries                string `i18n:"hint='Name of the media library model in the plural (menu, listing title).'"`
	ShowHidden                    string `i18n:"hint='Filter of the hidden files of the media library.'"`
	OnlyHidden                    string `i18n:"hint='Option of the hidden filter: only the hidden files.'"`
	IncludeHidden                 string `i18n:"hint='Option of the hidden filter: the hidden files too.'"`
	NotHidden                     string `i18n:"hint='Option of the hidden filter: only the files that are not hidden.'"`
	MediaLibraryDeletedByID       string `i18n:"hint='Label of the deleted by field of a media file.'"`
	MediaLibraryDeletedOrigin     string `i18n:"hint='Label of the deleted from field of a media file.'"`
	MediaLibraryHidden            string `i18n:"hint='Label of the hidden field of a media file.'"`
	MediaLibrarySelectedType      string `i18n:"hint='Label of the selected type field of a media file.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription:             "The media library: uploading, choosing and cropping files and images.",
	Crop:                          "Crop",
	CropImage:                     "Crop Image",
	EditImage:                     "Edit image",
	EditBeforeUpload:              "Edit before uploading",
	DropImages:                    "Drop the images here or click to choose",
	DropFiles:                     "Drop the files here or click to choose",
	AttachmentImages:              "Images",
	AttachmentFiles:               "Files",
	AttachmentImage:               "Image",
	AttachmentName:                "Name",
	AttachmentDescription:         "Description",
	AttachmentSize:                "Size",
	DropFile:                      "Drop the file here or click to choose",
	DropImage:                     "Drop the image here or click to choose",
	ImageDescription:              "Description",
	ChooseFile:                    "Choose File",
	Delete:                        "Delete",
	ChooseAFile:                   "Choose a File",
	ChosenFileIsGone:              "That file is no longer in the media library. Close this dialog and open it again.",
	FileInUse:                     "This file is in use (%d) and was not deleted. Take it out of:",
	MediaBoxUsage:                 "Image sizes of a field",
	MediaBoxUsages:                "Image sizes",
	MediaBoxUsageModel:            "Model",
	MediaBoxUsageLabel:            "Field",
	MediaBoxUsageSizes:            "Sizes",
	Format:                        "Format",
	FormatDefault:                 "Default",
	MediaBoxUsageFieldFormats:     "Formats",
	MediaFormat:                   "Image format",
	MediaFormats:                  "Image formats",
	MediaFormatName:               "Name",
	MediaFormatName_Hint:          "Lower-case letters, digits and -: the suffix of its crops (site.full_portrait). Not changed once made.",
	MediaFormatLabel:              "Label",
	MediaFormatWidth:              "Width",
	MediaFormatHeight:             "Height",
	MediaFormatRatio:              "Proportion",
	MediaFormatRatio_Hint:         "Width:height — 3:4 portrait, 1:1 square. Each size of a field in it keeps its shorter side.",
	MediaFormatCode:               "Of the code",
	MediaFormatUsages:             "Fields",
	MediaBoxUsageFormat:           "Format of a field",
	MediaBoxUsageFormats:          "Formats of the fields",
	MediaBoxUsageFormatUsage:      "Field",
	MediaBoxUsageFormatFormat:     "Format",
	MediaBoxUsageFormatPosition:   "Position",
	MediaBoxUsageFormatWidth:      "Width of the largest",
	MediaBoxUsageFormatWidth_Hint: "Empty: each size keeps its shorter side. A width: the largest size that wide, the others in proportion.",
	MediaBoxUsageFormatSizes:      "Sizes in the field",
	FormatNameInvalid:             "The name of a format is lower-case letters, digits and - (%s).",
	FormatNameTaken:               "Another format is named %s.",
	FormatNotPositive:             "The proportion needs a width and a height.",
	FormatInUse:                   "Fields offer this format: %s. Take it off them first.",
	FormatOfCode:                  "A format of the code is not deleted.",
	FormatUsageRequired:           "Choose the field.",
	FormatRequired:                "Choose the format.",
	FormatTwice:                   "The field already offers the format %s.",
	FormatWidthNegative:           "The width is 0 or more.",
	MediaBoxUsageRecropState:      "Crops",
	MediaBoxUsageRecropError:      "What went wrong",
	MediaBoxUsageOriginal:         "Largest original",
	MediaBoxUsageOriginalWidth:    "Largest width of the original",
	MediaBoxUsageOriginalHeight:   "Largest height of the original",
	OriginalOfTheCode:             "A larger image is made smaller to fit. Empty: the code's, %[1]d×%[2]d.",
	OriginalOfTheCodeNone:         "A larger image is made smaller to fit. Empty: the media library's.",
	OriginalBoth:                  "Give both the width and the height, or neither.",
	MediaBoxUsageThumb:            "Largest thumbnail",
	MediaBoxUsageThumbWidth:       "Largest width of the thumbnail",
	MediaBoxUsageThumbHeight:      "Largest height of the thumbnail",
	ThumbOfTheCode:                "The thumbnail of an image is made to fit it. Empty: the code's, %[1]d×%[2]d.",
	ThumbOfTheCodeNone:            "The thumbnail of an image is made to fit it. Empty: no thumbnail.",
	SizesOfTheAdmin:               "the admin's",
	Recrop:                        "Make the crops again",
	RecropStarted:                 "The crops are being made again.",
	RecropBusy:                    "The crops of this field are being made already.",
	RecropRunning:                 "being made: %d of %d",
	RecropMade:                    "made (%d files)",
	RecropFailedText:              "some could not be made",
	SizeNameRequired:              "A size has no name.",
	SizeNameTwice:                 "Two sizes are named %s.",
	SizeNotPositive:               "The size %s has no width or no height.",
	SizesConflict:                 "Another field has a size of this name in other measures (%s): a file used in both would be cropped in one, then the other.",
	RecommendedOriginalSize:       "Recommended original: at least %d×%d px, %s.",
	CopyLink:                      "Copy Link",
	LinkCopied:                    "Link Copied!",
	Download:                      "Download",
	Search:                        "Search",
	UploadFiles:                   "Upload files",
	TakePhoto:                     "Take a photo",
	Cropping:                      "Cropping",
	DescriptionUpdated:            "Description Updated",
	DescriptionForAccessibility:   "description for accessibility",
	OrderBy:                       "Order By",
	FileType:                      "Type",
	UploadedAt:                    "Date Uploaded",
	UploadedAtDESC:                "Date Uploaded (DESC)",
	All:                           "All",
	Images:                        "Images",
	Videos:                        "Videos",
	Files:                         "Files",
	ShowHidden:                    "Show hidden",
	OnlyHidden:                    "Only hidden",
	IncludeHidden:                 "Include hidden",
	NotHidden:                     "Not hidden",
	MediaLibraryDeletedByID:       "Deleted by",
	MediaLibraryDeletedOrigin:     "Deleted from",
	MediaLibraryHidden:            "Hidden",
	MediaLibrarySelectedType:      "Selected type",
}

var Messages_zh_CN = &Messages{
	ModuleDescription:             "媒体库：上传、选择和裁剪文件与图片。",
	Crop:                          "剪裁",
	CropImage:                     "剪裁图片",
	EditImage:                     "编辑图片",
	EditBeforeUpload:              "上传前编辑",
	DropImages:                    "将图片拖放到此处或点击选择",
	DropFiles:                     "将文件拖放到此处或点击选择",
	AttachmentImages:              "图片",
	AttachmentFiles:               "文件",
	AttachmentImage:               "图片",
	AttachmentName:                "名称",
	AttachmentDescription:         "描述",
	AttachmentSize:                "大小",
	DropFile:                      "将文件拖放到此处或点击选择",
	DropImage:                     "将图片拖放到此处或点击选择",
	ImageDescription:              "描述",
	ChooseFile:                    "选择文件",
	Delete:                        "删除",
	ChooseAFile:                   "选择一个文件",
	ChosenFileIsGone:              "该文件已不在媒体库中。请关闭此对话框后重新打开。",
	FileInUse:                     "该文件正在使用中（%d 处），未删除。请先从以下位置移除：",
	MediaBoxUsage:                 "字段的图片尺寸",
	MediaBoxUsages:                "图片尺寸",
	MediaBoxUsageModel:            "模型",
	MediaBoxUsageLabel:            "字段",
	MediaBoxUsageSizes:            "尺寸",
	Format:                        "格式",
	FormatDefault:                 "默认",
	MediaBoxUsageFieldFormats:     "格式",
	MediaFormat:                   "图片格式",
	MediaFormats:                  "图片格式",
	MediaFormatName:               "名称",
	MediaFormatName_Hint:          "小写字母、数字和 -：裁剪的后缀。创建后不可更改。",
	MediaFormatLabel:              "标签",
	MediaFormatWidth:              "宽度",
	MediaFormatHeight:             "高度",
	MediaFormatRatio:              "比例",
	MediaFormatRatio_Hint:         "宽:高 — 3:4 纵向，1:1 方形。",
	MediaFormatCode:               "代码定义",
	MediaFormatUsages:             "字段",
	MediaBoxUsageFormat:           "字段的格式",
	MediaBoxUsageFormats:          "字段的格式",
	MediaBoxUsageFormatUsage:      "字段",
	MediaBoxUsageFormatFormat:     "格式",
	MediaBoxUsageFormatPosition:   "位置",
	MediaBoxUsageFormatWidth:      "最大宽度",
	MediaBoxUsageFormatWidth_Hint: "留空：每个尺寸保持其短边。",
	MediaBoxUsageFormatSizes:      "字段中的尺寸",
	FormatNameInvalid:             "格式名称只能包含小写字母、数字和 -（%s）。",
	FormatNameTaken:               "已有名为 %s 的格式。",
	FormatNotPositive:             "比例需要宽度和高度。",
	FormatInUse:                   "以下字段使用此格式：%s。请先移除。",
	FormatOfCode:                  "代码定义的格式不能删除。",
	FormatUsageRequired:           "请选择字段。",
	FormatRequired:                "请选择格式。",
	FormatTwice:                   "该字段已提供格式 %s。",
	FormatWidthNegative:           "宽度须为 0 或更大。",
	MediaBoxUsageRecropState:      "裁剪",
	MediaBoxUsageRecropError:      "出错原因",
	MediaBoxUsageOriginal:         "原图最大尺寸",
	MediaBoxUsageOriginalWidth:    "原图最大宽度",
	MediaBoxUsageOriginalHeight:   "原图最大高度",
	OriginalOfTheCode:             "更大的图片会按比例缩小。留空：使用代码设置 %[1]d×%[2]d。",
	OriginalOfTheCodeNone:         "更大的图片会按比例缩小。留空：使用媒体库设置。",
	OriginalBoth:                  "宽度和高度需同时填写或都不填。",
	MediaBoxUsageThumb:            "缩略图最大尺寸",
	MediaBoxUsageThumbWidth:       "缩略图最大宽度",
	MediaBoxUsageThumbHeight:      "缩略图最大高度",
	ThumbOfTheCode:                "图片的缩略图按此尺寸生成。留空：使用代码设置 %[1]d×%[2]d。",
	ThumbOfTheCodeNone:            "图片的缩略图按此尺寸生成。留空：不生成缩略图。",
	SizesOfTheAdmin:               "管理员设置",
	Recrop:                        "重新裁剪",
	RecropStarted:                 "正在重新裁剪。",
	RecropBusy:                    "该字段正在裁剪中。",
	RecropRunning:                 "处理中：%d / %d",
	RecropMade:                    "已完成（%d 个文件）",
	RecropFailedText:              "部分未能完成",
	SizeNameRequired:              "有尺寸没有名称。",
	SizeNameTwice:                 "有两个尺寸都叫 %s。",
	SizeNotPositive:               "尺寸 %s 没有宽度或高度。",
	SizesConflict:                 "另一个字段有同名但不同尺寸的设置（%s）：两处都使用的文件会被来回裁剪。",
	RecommendedOriginalSize:       "建议原图：至少 %d×%d 像素，比例 %s。",
	Search:                        "搜索",
	UploadFiles:                   "上传多个文件",
	TakePhoto:                     "拍照",
	Cropping:                      "正在剪裁...",
	DescriptionUpdated:            "描述更新成功",
	DescriptionForAccessibility:   "图片描述",
	OrderBy:                       "排序",
	FileType:                      "类型",
	UploadedAt:                    "上传时间",
	UploadedAtDESC:                "上传时间 (降序)",
	All:                           "全部",
	Images:                        "图片",
	Videos:                        "视频",
	Files:                         "文件",
}

var Messages_ja_JP = &Messages{
	ModuleDescription:             "メディアライブラリ：ファイルと画像のアップロード、選択、トリミング。",
	Crop:                          "トリミング",
	CropImage:                     "画像をトリミング",
	EditImage:                     "画像を編集",
	EditBeforeUpload:              "アップロード前に編集",
	DropImages:                    "ここに画像をドロップするか、クリックして選択",
	DropFiles:                     "ここにファイルをドロップするか、クリックして選択",
	AttachmentImages:              "画像",
	AttachmentFiles:               "ファイル",
	AttachmentImage:               "画像",
	AttachmentName:                "名前",
	AttachmentDescription:         "説明",
	AttachmentSize:                "サイズ",
	DropFile:                      "ここにファイルをドロップするか、クリックして選択",
	DropImage:                     "ここに画像をドロップするか、クリックして選択",
	ImageDescription:              "説明",
	ChooseFile:                    "ファイルを選択",
	Delete:                        "削除",
	ChooseAFile:                   "ファイルを選択",
	ChosenFileIsGone:              "そのファイルはメディアライブラリにありません。ダイアログを閉じて開き直してください。",
	FileInUse:                     "このファイルは使用中（%d 件）のため削除されませんでした。次の場所から外してください：",
	MediaBoxUsage:                 "フィールドの画像サイズ",
	MediaBoxUsages:                "画像サイズ",
	MediaBoxUsageModel:            "モデル",
	MediaBoxUsageLabel:            "フィールド",
	MediaBoxUsageSizes:            "サイズ",
	Format:                        "フォーマット",
	FormatDefault:                 "デフォルト",
	MediaBoxUsageFieldFormats:     "フォーマット",
	MediaFormat:                   "画像フォーマット",
	MediaFormats:                  "画像フォーマット",
	MediaFormatName:               "名前",
	MediaFormatName_Hint:          "小文字・数字・-：切り抜きの接尾辞。作成後は変更不可。",
	MediaFormatLabel:              "ラベル",
	MediaFormatWidth:              "幅",
	MediaFormatHeight:             "高さ",
	MediaFormatRatio:              "比率",
	MediaFormatRatio_Hint:         "幅:高さ — 3:4 縦長、1:1 正方形。",
	MediaFormatCode:               "コード定義",
	MediaFormatUsages:             "フィールド",
	MediaBoxUsageFormat:           "フィールドのフォーマット",
	MediaBoxUsageFormats:          "フィールドのフォーマット",
	MediaBoxUsageFormatUsage:      "フィールド",
	MediaBoxUsageFormatFormat:     "フォーマット",
	MediaBoxUsageFormatPosition:   "位置",
	MediaBoxUsageFormatWidth:      "最大の幅",
	MediaBoxUsageFormatWidth_Hint: "空：各サイズは短辺を保ちます。",
	MediaBoxUsageFormatSizes:      "フィールドでのサイズ",
	FormatNameInvalid:             "フォーマット名は小文字・数字・- のみです（%s）。",
	FormatNameTaken:               "%s という名前のフォーマットがあります。",
	FormatNotPositive:             "比率には幅と高さが必要です。",
	FormatInUse:                   "このフォーマットを使うフィールド：%s。先に外してください。",
	FormatOfCode:                  "コード定義のフォーマットは削除できません。",
	FormatUsageRequired:           "フィールドを選択してください。",
	FormatRequired:                "フォーマットを選択してください。",
	FormatTwice:                   "このフィールドは既にフォーマット %s を提供しています。",
	FormatWidthNegative:           "幅は 0 以上です。",
	MediaBoxUsageRecropState:      "トリミング",
	MediaBoxUsageRecropError:      "エラーの内容",
	MediaBoxUsageOriginal:         "元画像の最大サイズ",
	MediaBoxUsageOriginalWidth:    "元画像の最大幅",
	MediaBoxUsageOriginalHeight:   "元画像の最大高さ",
	OriginalOfTheCode:             "これより大きい画像は縮小されます。空欄：コードの値 %[1]d×%[2]d。",
	OriginalOfTheCodeNone:         "これより大きい画像は縮小されます。空欄：メディアライブラリの値。",
	OriginalBoth:                  "幅と高さは両方入力するか、両方空欄にしてください。",
	MediaBoxUsageThumb:            "サムネイルの最大サイズ",
	MediaBoxUsageThumbWidth:       "サムネイルの最大幅",
	MediaBoxUsageThumbHeight:      "サムネイルの最大高さ",
	ThumbOfTheCode:                "画像のサムネイルはこのサイズに収まるように作られます。空欄：コードの値 %[1]d×%[2]d。",
	ThumbOfTheCodeNone:            "画像のサムネイルはこのサイズに収まるように作られます。空欄：サムネイルなし。",
	SizesOfTheAdmin:               "管理画面の設定",
	Recrop:                        "トリミングをやり直す",
	RecropStarted:                 "トリミングをやり直しています。",
	RecropBusy:                    "このフィールドはすでにトリミング中です。",
	RecropRunning:                 "処理中: %d / %d",
	RecropMade:                    "完了（%d ファイル）",
	RecropFailedText:              "一部を作成できませんでした",
	SizeNameRequired:              "名前のないサイズがあります。",
	SizeNameTwice:                 "%s という名前のサイズが 2 つあります。",
	SizeNotPositive:               "サイズ %s に幅または高さがありません。",
	SizesConflict:                 "別のフィールドに同じ名前で寸法の違うサイズがあります（%s）: 両方で使うファイルが交互にトリミングされます。",
	RecommendedOriginalSize:       "推奨する元画像: %d×%d px 以上、比率 %s。",
	Search:                        "検索",
	UploadFiles:                   "ファイルをアップロード",
	TakePhoto:                     "写真を撮る",
	Cropping:                      "トリミング中",
	DescriptionUpdated:            "説明を更新しました",
	DescriptionForAccessibility:   "画像の説明",
	OrderBy:                       "並び替え",
	FileType:                      "種類",
	UploadedAt:                    "アップロード日時",
	UploadedAtDESC:                "アップロード日時 (降順)",
	All:                           "すべて",
	Images:                        "画像",
	Videos:                        "動画",
	Files:                         "ファイル",
}

var Messages_pt_BR = &Messages{
	ModuleDescription:             "A biblioteca de mídia: envio, escolha e recorte de arquivos e imagens.",
	Crop:                          "Recortar",
	CropImage:                     "Recortar imagem",
	EditImage:                     "Editar imagem",
	EditBeforeUpload:              "Editar antes de enviar",
	DropImages:                    "Solte as imagens aqui ou clique para escolher",
	DropFiles:                     "Solte os arquivos aqui ou clique para escolher",
	AttachmentImages:              "Imagens",
	AttachmentFiles:               "Arquivos",
	AttachmentImage:               "Imagem",
	AttachmentName:                "Nome",
	AttachmentDescription:         "Descrição",
	AttachmentSize:                "Tamanho",
	DropFile:                      "Solte o arquivo aqui ou clique para escolher",
	DropImage:                     "Solte a imagem aqui ou clique para escolher",
	ImageDescription:              "Descrição",
	ChooseFile:                    "Escolher arquivo",
	Delete:                        "Excluir",
	ChooseAFile:                   "Escolha um arquivo",
	ChosenFileIsGone:              "Este arquivo não está mais na biblioteca de mídia. Feche esta janela e abra de novo.",
	FileInUse:                     "Este arquivo está em uso (%d) e não foi excluído. Tire-o de:",
	MediaBoxUsage:                 "Tamanhos das imagens de um campo",
	MediaBoxUsages:                "Tamanhos das imagens",
	MediaBoxUsageModel:            "Modelo",
	MediaBoxUsageLabel:            "Campo",
	MediaBoxUsageSizes:            "Tamanhos",
	Format:                        "Formato",
	FormatDefault:                 "Padrão",
	MediaBoxUsageFieldFormats:     "Formatos",
	MediaFormat:                   "Formato de imagem",
	MediaFormats:                  "Formatos de imagem",
	MediaFormatName:               "Nome",
	MediaFormatName_Hint:          "Letras minúsculas, dígitos e -: o sufixo dos recortes (site.full_retrato). Não muda depois de criado.",
	MediaFormatLabel:              "Rótulo",
	MediaFormatWidth:              "Largura",
	MediaFormatHeight:             "Altura",
	MediaFormatRatio:              "Proporção",
	MediaFormatRatio_Hint:         "Largura:altura — 3:4 retrato, 1:1 quadrado. Cada tamanho de um campo mantém nele o seu lado menor.",
	MediaFormatCode:               "Do código",
	MediaFormatUsages:             "Campos",
	MediaBoxUsageFormat:           "Formato de um campo",
	MediaBoxUsageFormats:          "Formatos dos campos",
	MediaBoxUsageFormatUsage:      "Campo",
	MediaBoxUsageFormatFormat:     "Formato",
	MediaBoxUsageFormatPosition:   "Posição",
	MediaBoxUsageFormatWidth:      "Largura do maior",
	MediaBoxUsageFormatWidth_Hint: "Vazio: cada tamanho mantém o seu lado menor. Uma largura: o maior tamanho com ela, os outros na proporção.",
	MediaBoxUsageFormatSizes:      "Tamanhos no campo",
	FormatNameInvalid:             "O nome de um formato é de letras minúsculas, dígitos e - (%s).",
	FormatNameTaken:               "Outro formato tem o nome %s.",
	FormatNotPositive:             "A proporção precisa de largura e altura.",
	FormatInUse:                   "Campos oferecem este formato: %s. Tire-o deles antes.",
	FormatOfCode:                  "Um formato do código não é excluído.",
	FormatUsageRequired:           "Escolha o campo.",
	FormatRequired:                "Escolha o formato.",
	FormatTwice:                   "O campo já oferece o formato %s.",
	FormatWidthNegative:           "A largura é 0 ou mais.",
	MediaBoxUsageRecropState:      "Recortes",
	MediaBoxUsageRecropError:      "O que deu errado",
	MediaBoxUsageOriginal:         "Original máximo",
	MediaBoxUsageOriginalWidth:    "Largura máxima do original",
	MediaBoxUsageOriginalHeight:   "Altura máxima do original",
	OriginalOfTheCode:             "Uma imagem maior é reduzida para caber, proporcionalmente. Vazio: o do código, %[1]d×%[2]d.",
	OriginalOfTheCodeNone:         "Uma imagem maior é reduzida para caber, proporcionalmente. Vazio: o da biblioteca de mídia.",
	OriginalBoth:                  "Informe a largura e a altura, ou nenhuma das duas.",
	MediaBoxUsageThumb:            "Miniatura máxima",
	MediaBoxUsageThumbWidth:       "Largura máxima da miniatura",
	MediaBoxUsageThumbHeight:      "Altura máxima da miniatura",
	ThumbOfTheCode:                "A miniatura de uma imagem é feita para caber nela, proporcionalmente. Vazio: a do código, %[1]d×%[2]d.",
	ThumbOfTheCodeNone:            "A miniatura de uma imagem é feita para caber nela, proporcionalmente. Vazio: sem miniatura.",
	SizesOfTheAdmin:               "do admin",
	Recrop:                        "Refazer os recortes",
	RecropStarted:                 "Os recortes estão sendo refeitos.",
	RecropBusy:                    "Os recortes deste campo já estão sendo refeitos.",
	RecropRunning:                 "em andamento: %d de %d",
	RecropMade:                    "prontos (%d arquivos)",
	RecropFailedText:              "alguns não puderam ser feitos",
	SizeNameRequired:              "Um tamanho está sem nome.",
	SizeNameTwice:                 "Dois tamanhos se chamam %s.",
	SizeNotPositive:               "O tamanho %s está sem largura ou altura.",
	SizesConflict:                 "Outro campo tem um tamanho com este nome em outras medidas (%s): um arquivo usado nos dois seria recortado em um, depois no outro.",
	RecommendedOriginalSize:       "Original recomendada: pelo menos %d×%d px, proporção %s.",
	CopyLink:                      "Copiar link",
	LinkCopied:                    "Link copiado!",
	Download:                      "Baixar",
	Search:                        "Pesquisar",
	UploadFiles:                   "Enviar arquivos",
	TakePhoto:                     "Tirar foto",
	Cropping:                      "Recortando",
	DescriptionUpdated:            "Descrição atualizada",
	DescriptionForAccessibility:   "descrição para acessibilidade",
	OrderBy:                       "Ordenar por",
	FileType:                      "Tipo",
	UploadedAt:                    "Data de envio",
	UploadedAtDESC:                "Data de envio (decrescente)",
	All:                           "Todos",
	Images:                        "Imagens",
	Videos:                        "Vídeos",
	Files:                         "Arquivos",
	MediaLibrary:                  "Biblioteca de mídia",
	MediaLibraries:                "Bibliotecas de mídia",
	ShowHidden:                    "Mostrar ocultos",
	OnlyHidden:                    "Somente ocultos",
	IncludeHidden:                 "Incluir ocultos",
	NotHidden:                     "Não ocultos",
	MediaLibraryDeletedByID:       "Excluído por",
	MediaLibraryDeletedOrigin:     "Excluído de",
	MediaLibraryHidden:            "Oculto",
	MediaLibrarySelectedType:      "Tipo selecionado",
}
