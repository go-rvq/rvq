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
}

var Messages_en_US = &Messages{
	Crop:                        "Crop",
	CropImage:                   "Crop Image",
	ChooseFile:                  "Choose File",
	Delete:                      "Delete",
	ChooseAFile:                 "Choose a File",
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
}

var Messages_zh_CN = &Messages{
	Crop:                        "剪裁",
	CropImage:                   "剪裁图片",
	ChooseFile:                  "选择文件",
	Delete:                      "删除",
	ChooseAFile:                 "选择一个文件",
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
