package vuetifyx

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// MessagesKey is the module of the words of the vuetifyx components — their
// buttons, titles, the filter, the pickers, the editors —: the namespace
// MessagesNamespace of the page's vue-i18n, merged once a page in the user's
// language (presets' layout), over the components' defaults.
const MessagesKey i18n.ModuleKey = "rvq-admin/vuetifyx"

// MessagesNamespace is their namespace in the vue-i18n: `{ vx: { … } }`.
const MessagesNamespace = "vx"

// GetMessages are the words in ctx's language.
func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

// ConfigureMessages registers the words.
func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModule(language.English, MessagesKey, Messages_en_US).
		RegisterForModule(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}

// Messages are the words of the vuetifyx components; their JSON is the
// namespace MessagesNamespace of the vue-i18n (the same keys as the
// components' VxMessages).
type Messages struct {
	ModuleDescription        string `json:"-" i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Back                     string `json:"back" i18n:"hint='Button that goes back.'"`
	Save                     string `json:"save" i18n:"hint='Button that saves.'"`
	Apply                    string `json:"apply" i18n:"hint='Button that applies.'"`
	Cancel                   string `json:"cancel" i18n:"hint='Button that cancels.'"`
	Ok                       string `json:"ok" i18n:"hint='Button that confirms.'"`
	Close                    string `json:"close" i18n:"hint='Button that closes.'"`
	Clear                    string `json:"clear" i18n:"hint='Button that clears a value or a filter.'"`
	Undo                     string `json:"undo" i18n:"hint='Button that undoes the last edit.'"`
	Redo                     string `json:"redo" i18n:"hint='Button that redoes the edit undone.'"`
	Move                     string `json:"move" i18n:"hint='Button that moves.'"`
	ChangeOrder              string `json:"changeOrder" i18n:"hint='Button that changes the order of items.'"`
	Select                   string `json:"select" i18n:"hint='Button that selects.'"`
	Home                     string `json:"home" i18n:"hint='The first page.'"`
	Preview                  string `json:"preview" i18n:"hint='Button that previews.'"`
	Help                     string `json:"help" i18n:"hint='Button and title of a help.'"`
	Expand                   string `json:"expand" i18n:"hint='Button that expands.'"`
	Collapse                 string `json:"collapse" i18n:"hint='Button that collapses.'"`
	Settings                 string `json:"settings" i18n:"hint='Title of the settings of a table.'"`
	Columns                  string `json:"columns" i18n:"hint='Title of the columns of a table.'"`
	ThemeLight               string `json:"themeLight" i18n:"hint='Title of the button that turns the theme light.'"`
	ThemeDark                string `json:"themeDark" i18n:"hint='Title of the button that turns the theme dark.'"`
	Year                     string `json:"year" i18n:"hint='Label of a year (a month picker).'"`
	Month                    string `json:"month" i18n:"hint='Label of a month (a month picker).'"`
	FilterClearAll           string `json:"filterClearAll" i18n:"hint='Button that removes every filter.'"`
	FilterAdd                string `json:"filterAdd" i18n:"hint='Button that adds a filter.'"`
	FilterBy                 string `json:"filterBy" i18n:"hint='Title of a filter ({filter}: its name).'"`
	FilterTo                 string `json:"filterTo" i18n:"hint='Between the two values of a range of a filter.'"`
	FilterNumberEquals       string `json:"filterNumberEquals" i18n:"hint='Modifier of a number filter.'"`
	FilterNumberBetween      string `json:"filterNumberBetween" i18n:"hint='Modifier of a number filter.'"`
	FilterNumberGreaterThan  string `json:"filterNumberGreaterThan" i18n:"hint='Modifier of a number filter.'"`
	FilterNumberLessThan     string `json:"filterNumberLessThan" i18n:"hint='Modifier of a number filter.'"`
	FilterNumberAnd          string `json:"filterNumberAnd" i18n:"hint='Between the two numbers of a number filter.'"`
	FilterStringEquals       string `json:"filterStringEquals" i18n:"hint='Modifier of a text filter.'"`
	FilterStringContains     string `json:"filterStringContains" i18n:"hint='Modifier of a text filter.'"`
	FilterIn                 string `json:"filterIn" i18n:"hint='Modifier of a choice filter: one of them.'"`
	FilterNotIn              string `json:"filterNotIn" i18n:"hint='Modifier of a choice filter: none of them.'"`
	Copy                     string `json:"copy" i18n:"hint='Button that copies.'"`
	Copied                   string `json:"copied" i18n:"hint='Said once copied.'"`
	CopyError                string `json:"copyError" i18n:"hint='Said when the copy failed.'"`
	Paste                    string `json:"paste" i18n:"hint='Button that pastes.'"`
	CopyOrPasteContent       string `json:"copyOrPasteContent" i18n:"hint='Menu of the copy and paste of a content.'"`
	CopyContent              string `json:"copyContent" i18n:"hint='Item that copies the content.'"`
	PasteContent             string `json:"pasteContent" i18n:"hint='Item that pastes the content.'"`
	CopiedToClipboard        string `json:"copiedToClipboard" i18n:"hint='Said once copied.'"`
	PastedFromClipboard      string `json:"pastedFromClipboard" i18n:"hint='Said once pasted.'"`
	InvalidContent           string `json:"invalidContent" i18n:"hint='Said when what was pasted is not valid.'"`
	PastedIsEmpty            string `json:"pastedIsEmpty" i18n:"hint='Said when the clipboard is empty.'"`
	Tools                    string `json:"tools" i18n:"hint='Menu of tools.'"`
	PleaseEnterTheCode       string `json:"pleaseEnterTheCode" i18n:"hint='Asks for a code.'"`
	Script                   string `json:"script" i18n:"hint='A script.'"`
	ScriptResult             string `json:"scriptResult" i18n:"hint='The result of a script.'"`
	ShowDetail               string `json:"showDetail" i18n:"hint='Button that shows the detail.'"`
	InsertTextBlockBefore    string `json:"insertTextBlockBefore" i18n:"hint='Button of the editor.'"`
	InsertTextBlockAfter     string `json:"insertTextBlockAfter" i18n:"hint='Button of the editor.'"`
	PreFormattedText         string `json:"preFormattedText" i18n:"hint='A pre formatted text.'"`
	SourceCodeHtml           string `json:"sourceCodeHtml" i18n:"hint='Title of the dialog of the HTML of a content.'"`
	SelectImage              string `json:"selectImage" i18n:"hint='Label of the choice of an image.'"`
	ImageAlt                 string `json:"imageAlt" i18n:"hint='Label of the alternative text of an image.'"`
	LockAspectRatio          string `json:"lockAspectRatio" i18n:"hint='Option of an image.'"`
	HeaderPlaceholder        string `json:"headerPlaceholder" i18n:"hint='Placeholder of a header of the block editor.'"`
	ImageLoadError           string `json:"imageLoadError" i18n:"hint='Said when an image could not be loaded.'"`
	ImageReleaseHere         string `json:"imageReleaseHere" i18n:"hint='Said while a file is dragged over.'"`
	ImageUnsupported         string `json:"imageUnsupported" i18n:"hint='Said of a file not supported.'"`
	ImageAddFromLink         string `json:"imageAddFromLink" i18n:"hint='Button that adds an image of a link.'"`
	ImageTakePicture         string `json:"imageTakePicture" i18n:"hint='Button that takes a picture.'"`
	ImageEdit                string `json:"imageEdit" i18n:"hint='Button that edits the picture.'"`
	ImageAddFromFile         string `json:"imageAddFromFile" i18n:"hint='Button that adds an image of a file.'"`
	ImageFullView            string `json:"imageFullView" i18n:"hint='Button that shows the image whole.'"`
	TakePicture              string `json:"takePicture" i18n:"hint='Button of the camera.'"`
	FlipCamera               string `json:"flipCamera" i18n:"hint='Button of the camera.'"`
	LoadCameras              string `json:"loadCameras" i18n:"hint='Button of the camera.'"`
	CameraDenied             string `json:"cameraDenied" i18n:"hint='Said when the browser denied the camera.'"`
	CameraNotFound           string `json:"cameraNotFound" i18n:"hint='Said when there is no camera.'"`
	CameraError              string `json:"cameraError" i18n:"hint='Said when the camera failed ({error}: why).'"`
	MapLoadError             string `json:"mapLoadError" i18n:"hint='Said when the map could not be loaded.'"`
	MapsLoadError            string `json:"mapsLoadError" i18n:"hint='Said when Google Maps could not be loaded.'"`
	MapsNoKey                string `json:"mapsNoKey" i18n:"hint='Said when there is no key of Google Maps.'"`
	SeoName                  string `json:"seoName" i18n:"hint='Column of the name of a SEO variable.'"`
	SeoValue                 string `json:"seoValue" i18n:"hint='Column of the value of a SEO variable.'"`
	SeoAdd                   string `json:"seoAdd" i18n:"hint='Button that adds a SEO variable.'"`
	SeoText                  string `json:"seoText" i18n:"hint='A SEO variable of text.'"`
	SeoOneLine               string `json:"seoOneLine" i18n:"hint='Title of the button that makes a value one line.'"`
	SeoMultiline             string `json:"seoMultiline" i18n:"hint='Title of the button that makes a value several lines.'"`
	SeoMap                   string `json:"seoMap" i18n:"hint='Button that opens the map.'"`
	SeoPickPlaces            string `json:"seoPickPlaces" i18n:"hint='Placeholder of the places of a ZIP variable.'"`
	SeoZipcodesTitle         string `json:"seoZipcodesTitle" i18n:"hint='Title of the map of the ZIP codes.'"`
	MapsTest                 string `json:"mapsTest" i18n:"hint='Button that tests the key of Google Maps.'"`
	MapsTestTitle            string `json:"mapsTestTitle" i18n:"hint='Title of the test of the key.'"`
	MapsTestNoKey            string `json:"mapsTestNoKey" i18n:"hint='Said with no key.'"`
	MapsTestMaps             string `json:"mapsTestMaps" i18n:"hint='An API of Google.'"`
	MapsTestPlaces           string `json:"mapsTestPlaces" i18n:"hint='An API of Google.'"`
	MapsTestGeocoding        string `json:"mapsTestGeocoding" i18n:"hint='An API of Google.'"`
	MapsTestWaiting          string `json:"mapsTestWaiting" i18n:"hint='Said while testing.'"`
	MapsTestOk               string `json:"mapsTestOk" i18n:"hint='Said of an API that works.'"`
	MapsTestFailed           string `json:"mapsTestFailed" i18n:"hint='Said of an API that failed.'"`
	MapsTestTryIt            string `json:"mapsTestTryIt" i18n:"hint='Hint of the test.'"`
	MapsTestReload           string `json:"mapsTestReload" i18n:"hint='Said when another key is loaded.'"`
	MapsTestRejected         string `json:"mapsTestRejected" i18n:"hint='Said when Google refused the key ({code}: its error).'"`
	MapsErrRefererNotAllowed string `json:"mapsErrRefererNotAllowed" i18n:"hint='An error of Google Maps.'"`
	MapsErrApiNotActivated   string `json:"mapsErrApiNotActivated" i18n:"hint='An error of Google Maps.'"`
	MapsErrInvalidKey        string `json:"mapsErrInvalidKey" i18n:"hint='An error of Google Maps.'"`
	MapsErrBillingNotEnabled string `json:"mapsErrBillingNotEnabled" i18n:"hint='An error of Google Maps.'"`
	MapsErrExpiredKey        string `json:"mapsErrExpiredKey" i18n:"hint='An error of Google Maps.'"`
	MapsErrZeroResults       string `json:"mapsErrZeroResults" i18n:"hint='An error of Google Maps.'"`
	MapsErrRequestDenied     string `json:"mapsErrRequestDenied" i18n:"hint='An error of Google Maps.'"`
	MapsErrOverQueryLimit    string `json:"mapsErrOverQueryLimit" i18n:"hint='An error of Google Maps.'"`
	Format                   string `json:"format" i18n:"hint='Button that formats the code (gad fmt).'"`
	Structure                string `json:"structure" i18n:"hint='Title of the panel of the declarations of the code.'"`
	Doc                      string `json:"doc" i18n:"hint='Title of the panel of the documentation of the code.'"`
	Refresh                  string `json:"refresh" i18n:"hint='Button that reads again.'"`
	NoStructure              string `json:"noStructure" i18n:"hint='Said when the code declares nothing.'"`
	NoDoc                    string `json:"noDoc" i18n:"hint='Said when the code has no documentation.'"`
	NotFormatted             string `json:"notFormatted" i18n:"hint='Said when the code does not parse.'"`
	Compare                  string `json:"compare" i18n:"hint='Button that compares the code with its initial value.'"`
	Edit                     string `json:"edit" i18n:"hint='Button that leaves the compare.'"`
	PrevChange               string `json:"prevChange" i18n:"hint='Button of the compare.'"`
	NextChange               string `json:"nextChange" i18n:"hint='Button of the compare.'"`
	Revert                   string `json:"revert" i18n:"hint='Button of a change of the compare: takes the old part back.'"`
	DiffFiles                string `json:"diffFiles" i18n:"hint='Title of the tree of the files of a diff.'"`
	DiffOld                  string `json:"diffOld" i18n:"hint='Title of the old side of a diff.'"`
	DiffCurrent              string `json:"diffCurrent" i18n:"hint='Title of the current side of a diff.'"`
	DiffBinary               string `json:"diffBinary" i18n:"hint='Said of a binary file.'"`
	DiffNoChanges            string `json:"diffNoChanges" i18n:"hint='Said with no changes.'"`
	DiffRenamedTo            string `json:"diffRenamedTo" i18n:"hint='Before the new path of a renamed file.'"`
	DiffUnchanged            string `json:"diffUnchanged" i18n:"hint='Said of a renamed file of the same content.'"`
	Loading                  string `json:"loading" i18n:"hint='Said while loading.'"`
	DiffLoadError            string `json:"diffLoadError" i18n:"hint='Said when a diff could not be loaded.'"`
	Saving                   string `json:"saving" i18n:"hint='Said while saving.'"`
	Saved                    string `json:"saved" i18n:"hint='Said once saved.'"`
	Unsaved                  string `json:"unsaved" i18n:"hint='Said of a change not saved.'"`
	DiffRevert               string `json:"diffRevert" i18n:"hint='Title of the button of a change of a diff.'"`
	DiffPrev                 string `json:"diffPrev" i18n:"hint='Title of the button of a diff.'"`
	DiffNext                 string `json:"diffNext" i18n:"hint='Title of the button of a diff.'"`

	// the words of the editors of a text (mdkit: js/mdkit), by editor
	MdkitRichBold               string `json:"mdkitRichBold" i18n:"hint='Editors of a text (mdkit), rich: bold.'"`
	MdkitRichItalic             string `json:"mdkitRichItalic" i18n:"hint='Editors of a text (mdkit), rich: italic.'"`
	MdkitRichStrike             string `json:"mdkitRichStrike" i18n:"hint='Editors of a text (mdkit), rich: strike.'"`
	MdkitRichHeading            string `json:"mdkitRichHeading" i18n:"hint='Editors of a text (mdkit), rich: heading.'"`
	MdkitRichSubheading         string `json:"mdkitRichSubheading" i18n:"hint='Editors of a text (mdkit), rich: subheading.'"`
	MdkitRichBullets            string `json:"mdkitRichBullets" i18n:"hint='Editors of a text (mdkit), rich: bullets.'"`
	MdkitRichNumbers            string `json:"mdkitRichNumbers" i18n:"hint='Editors of a text (mdkit), rich: numbers.'"`
	MdkitRichQuote              string `json:"mdkitRichQuote" i18n:"hint='Editors of a text (mdkit), rich: quote.'"`
	MdkitRichCode               string `json:"mdkitRichCode" i18n:"hint='Editors of a text (mdkit), rich: code.'"`
	MdkitRichLink               string `json:"mdkitRichLink" i18n:"hint='Editors of a text (mdkit), rich: link.'"`
	MdkitRichLinkPrompt         string `json:"mdkitRichLinkPrompt" i18n:"hint='Editors of a text (mdkit), rich: linkPrompt.'"`
	MdkitRichLinkSet            string `json:"mdkitRichLinkSet" i18n:"hint='Editors of a text (mdkit), rich: linkSet.'"`
	MdkitRichCancel             string `json:"mdkitRichCancel" i18n:"hint='Editors of a text (mdkit), rich: cancel.'"`
	MdkitRichImage              string `json:"mdkitRichImage" i18n:"hint='Editors of a text (mdkit), rich: image.'"`
	MdkitRichUndo               string `json:"mdkitRichUndo" i18n:"hint='Editors of a text (mdkit), rich: undo.'"`
	MdkitRichRedo               string `json:"mdkitRichRedo" i18n:"hint='Editors of a text (mdkit), rich: redo.'"`
	MdkitRichNotImage           string `json:"mdkitRichNotImage" i18n:"hint='Editors of a text (mdkit), rich: notImage.'"`
	MdkitRichUploadError        string `json:"mdkitRichUploadError" i18n:"hint='Editors of a text (mdkit), rich: uploadError.'"`
	MdkitRichPlaceholder        string `json:"mdkitRichPlaceholder" i18n:"hint='Editors of a text (mdkit), rich: placeholder.'"`
	MdkitRichPlaceholderText    string `json:"mdkitRichPlaceholderText" i18n:"hint='Editors of a text (mdkit), rich: placeholderText.'"`
	MdkitRichFootnote           string `json:"mdkitRichFootnote" i18n:"hint='Editors of a text (mdkit), rich: footnote.'"`
	MdkitRichSizeSmall          string `json:"mdkitRichSizeSmall" i18n:"hint='Editors of a text (mdkit), rich: size_small.'"`
	MdkitRichSizeMedium         string `json:"mdkitRichSizeMedium" i18n:"hint='Editors of a text (mdkit), rich: size_medium.'"`
	MdkitRichSizeLarge          string `json:"mdkitRichSizeLarge" i18n:"hint='Editors of a text (mdkit), rich: size_large.'"`
	MdkitRichSizeOriginal       string `json:"mdkitRichSizeOriginal" i18n:"hint='Editors of a text (mdkit), rich: size_original.'"`
	MdkitTableTitle             string `json:"mdkitTableTitle" i18n:"hint='Editors of a text (mdkit), table: title.'"`
	MdkitTableRowAbove          string `json:"mdkitTableRowAbove" i18n:"hint='Editors of a text (mdkit), table: rowAbove.'"`
	MdkitTableRowBelow          string `json:"mdkitTableRowBelow" i18n:"hint='Editors of a text (mdkit), table: rowBelow.'"`
	MdkitTableRowRemove         string `json:"mdkitTableRowRemove" i18n:"hint='Editors of a text (mdkit), table: rowRemove.'"`
	MdkitTableColLeft           string `json:"mdkitTableColLeft" i18n:"hint='Editors of a text (mdkit), table: colLeft.'"`
	MdkitTableColRight          string `json:"mdkitTableColRight" i18n:"hint='Editors of a text (mdkit), table: colRight.'"`
	MdkitTableColRemove         string `json:"mdkitTableColRemove" i18n:"hint='Editors of a text (mdkit), table: colRemove.'"`
	MdkitTableMerge             string `json:"mdkitTableMerge" i18n:"hint='Editors of a text (mdkit), table: merge.'"`
	MdkitTableUnmerge           string `json:"mdkitTableUnmerge" i18n:"hint='Editors of a text (mdkit), table: unmerge.'"`
	MdkitTableAlignLeft         string `json:"mdkitTableAlignLeft" i18n:"hint='Editors of a text (mdkit), table: alignLeft.'"`
	MdkitTableAlignCenter       string `json:"mdkitTableAlignCenter" i18n:"hint='Editors of a text (mdkit), table: alignCenter.'"`
	MdkitTableAlignRight        string `json:"mdkitTableAlignRight" i18n:"hint='Editors of a text (mdkit), table: alignRight.'"`
	MdkitTableHeader            string `json:"mdkitTableHeader" i18n:"hint='Editors of a text (mdkit), table: header.'"`
	MdkitTableFootnote          string `json:"mdkitTableFootnote" i18n:"hint='Editors of a text (mdkit), table: footnote.'"`
	MdkitTablePreview           string `json:"mdkitTablePreview" i18n:"hint='Editors of a text (mdkit), table: preview.'"`
	MdkitTableResize            string `json:"mdkitTableResize" i18n:"hint='Editors of a text (mdkit), table: resize.'"`
	MdkitTableWidthFixed        string `json:"mdkitTableWidthFixed" i18n:"hint='Editors of a text (mdkit), table: widthFixed.'"`
	MdkitTableWidthPercent      string `json:"mdkitTableWidthPercent" i18n:"hint='Editors of a text (mdkit), table: widthPercent.'"`
	MdkitTableWidthAuto         string `json:"mdkitTableWidthAuto" i18n:"hint='Editors of a text (mdkit), table: widthAuto.'"`
	MdkitTableFluid             string `json:"mdkitTableFluid" i18n:"hint='Editors of a text (mdkit), table: fluid.'"`
	MdkitTablePositionLeft      string `json:"mdkitTablePositionLeft" i18n:"hint='Editors of a text (mdkit), table: positionLeft.'"`
	MdkitTablePositionCenter    string `json:"mdkitTablePositionCenter" i18n:"hint='Editors of a text (mdkit), table: positionCenter.'"`
	MdkitTablePositionRight     string `json:"mdkitTablePositionRight" i18n:"hint='Editors of a text (mdkit), table: positionRight.'"`
	MdkitTableValignTop         string `json:"mdkitTableValignTop" i18n:"hint='Editors of a text (mdkit), table: valignTop.'"`
	MdkitTableValignMiddle      string `json:"mdkitTableValignMiddle" i18n:"hint='Editors of a text (mdkit), table: valignMiddle.'"`
	MdkitTableValignBottom      string `json:"mdkitTableValignBottom" i18n:"hint='Editors of a text (mdkit), table: valignBottom.'"`
	MdkitTableCellPlace         string `json:"mdkitTableCellPlace" i18n:"hint='Editors of a text (mdkit), table: cellPlace.'"`
	MdkitTableExpand            string `json:"mdkitTableExpand" i18n:"hint='Editors of a text (mdkit), table: expand.'"`
	MdkitTableShrink            string `json:"mdkitTableShrink" i18n:"hint='Editors of a text (mdkit), table: shrink.'"`
	MdkitTablePaste             string `json:"mdkitTablePaste" i18n:"hint='Editors of a text (mdkit), table: paste.'"`
	MdkitTablePasteTitle        string `json:"mdkitTablePasteTitle" i18n:"hint='Editors of a text (mdkit), table: pasteTitle.'"`
	MdkitTablePasteHint         string `json:"mdkitTablePasteHint" i18n:"hint='Editors of a text (mdkit), table: pasteHint.'"`
	MdkitTablePastePlaceholder  string `json:"mdkitTablePastePlaceholder" i18n:"hint='Editors of a text (mdkit), table: pastePlaceholder.'"`
	MdkitTablePasteFile         string `json:"mdkitTablePasteFile" i18n:"hint='Editors of a text (mdkit), table: pasteFile.'"`
	MdkitTableSeparator         string `json:"mdkitTableSeparator" i18n:"hint='Editors of a text (mdkit), table: separator.'"`
	MdkitTableSepAuto           string `json:"mdkitTableSepAuto" i18n:"hint='Editors of a text (mdkit), table: sepAuto.'"`
	MdkitTableSepTab            string `json:"mdkitTableSepTab" i18n:"hint='Editors of a text (mdkit), table: sepTab.'"`
	MdkitTableSepComma          string `json:"mdkitTableSepComma" i18n:"hint='Editors of a text (mdkit), table: sepComma.'"`
	MdkitTableSepSemicolon      string `json:"mdkitTableSepSemicolon" i18n:"hint='Editors of a text (mdkit), table: sepSemicolon.'"`
	MdkitTableSepPipe           string `json:"mdkitTableSepPipe" i18n:"hint='Editors of a text (mdkit), table: sepPipe.'"`
	MdkitTableSepOther          string `json:"mdkitTableSepOther" i18n:"hint='Editors of a text (mdkit), table: sepOther.'"`
	MdkitTableSepOtherChar      string `json:"mdkitTableSepOtherChar" i18n:"hint='Editors of a text (mdkit), table: sepOtherChar.'"`
	MdkitTablePasteHTML         string `json:"mdkitTablePasteHTML" i18n:"hint='Editors of a text (mdkit), table: pasteHTML.'"`
	MdkitTablePasteReplace      string `json:"mdkitTablePasteReplace" i18n:"hint='Editors of a text (mdkit), table: pasteReplace.'"`
	MdkitTablePasteAt           string `json:"mdkitTablePasteAt" i18n:"hint='Editors of a text (mdkit), table: pasteAt.'"`
	MdkitTablePasteFound        string `json:"mdkitTablePasteFound" i18n:"hint='Editors of a text (mdkit), table: pasteFound.'"`
	MdkitTablePasteDo           string `json:"mdkitTablePasteDo" i18n:"hint='Editors of a text (mdkit), table: pasteDo.'"`
	MdkitTableModeSheet         string `json:"mdkitTableModeSheet" i18n:"hint='Editors of a text (mdkit), table: modeSheet.'"`
	MdkitTableModeSplit         string `json:"mdkitTableModeSplit" i18n:"hint='Editors of a text (mdkit), table: modeSplit.'"`
	MdkitTableModePreview       string `json:"mdkitTableModePreview" i18n:"hint='Editors of a text (mdkit), table: modePreview.'"`
	MdkitTableHtmlNote          string `json:"mdkitTableHtmlNote" i18n:"hint='Editors of a text (mdkit), table: htmlNote.'"`
	MdkitTableInsert            string `json:"mdkitTableInsert" i18n:"hint='Editors of a text (mdkit), table: insert.'"`
	MdkitTableUpdate            string `json:"mdkitTableUpdate" i18n:"hint='Editors of a text (mdkit), table: update.'"`
	MdkitTableCancel            string `json:"mdkitTableCancel" i18n:"hint='Editors of a text (mdkit), table: cancel.'"`
	MdkitFootnoteTitle          string `json:"mdkitFootnoteTitle" i18n:"hint='Editors of a text (mdkit), footnote: title.'"`
	MdkitFootnoteAdd            string `json:"mdkitFootnoteAdd" i18n:"hint='Editors of a text (mdkit), footnote: add.'"`
	MdkitFootnoteEmpty          string `json:"mdkitFootnoteEmpty" i18n:"hint='Editors of a text (mdkit), footnote: empty.'"`
	MdkitFootnoteNoNote         string `json:"mdkitFootnoteNoNote" i18n:"hint='Editors of a text (mdkit), footnote: noNote.'"`
	MdkitFootnoteNoRef          string `json:"mdkitFootnoteNoRef" i18n:"hint='Editors of a text (mdkit), footnote: noRef.'"`
	MdkitFootnoteRefs           string `json:"mdkitFootnoteRefs" i18n:"hint='Editors of a text (mdkit), footnote: refs.'"`
	MdkitFootnoteLabel          string `json:"mdkitFootnoteLabel" i18n:"hint='Editors of a text (mdkit), footnote: label.'"`
	MdkitFootnoteText           string `json:"mdkitFootnoteText" i18n:"hint='Editors of a text (mdkit), footnote: text.'"`
	MdkitFootnoteRename         string `json:"mdkitFootnoteRename" i18n:"hint='Editors of a text (mdkit), footnote: rename.'"`
	MdkitFootnoteRemove         string `json:"mdkitFootnoteRemove" i18n:"hint='Editors of a text (mdkit), footnote: remove.'"`
	MdkitFootnoteRemoveAsk      string `json:"mdkitFootnoteRemoveAsk" i18n:"hint='Editors of a text (mdkit), footnote: removeAsk.'"`
	MdkitFootnoteRemoveYes      string `json:"mdkitFootnoteRemoveYes" i18n:"hint='Editors of a text (mdkit), footnote: removeYes.'"`
	MdkitFootnoteInsert         string `json:"mdkitFootnoteInsert" i18n:"hint='Editors of a text (mdkit), footnote: insert.'"`
	MdkitFootnoteUseHere        string `json:"mdkitFootnoteUseHere" i18n:"hint='Editors of a text (mdkit), footnote: useHere.'"`
	MdkitFootnoteHere           string `json:"mdkitFootnoteHere" i18n:"hint='Editors of a text (mdkit), footnote: here.'"`
	MdkitFootnoteSave           string `json:"mdkitFootnoteSave" i18n:"hint='Editors of a text (mdkit), footnote: save.'"`
	MdkitFootnoteChoose         string `json:"mdkitFootnoteChoose" i18n:"hint='Editors of a text (mdkit), footnote: choose.'"`
	MdkitFootnoteTaken          string `json:"mdkitFootnoteTaken" i18n:"hint='Editors of a text (mdkit), footnote: taken.'"`
	MdkitFootnoteModeWrite      string `json:"mdkitFootnoteModeWrite" i18n:"hint='Editors of a text (mdkit), footnote: modeWrite.'"`
	MdkitFootnoteModeSplit      string `json:"mdkitFootnoteModeSplit" i18n:"hint='Editors of a text (mdkit), footnote: modeSplit.'"`
	MdkitFootnoteModePreview    string `json:"mdkitFootnoteModePreview" i18n:"hint='Editors of a text (mdkit), footnote: modePreview.'"`
	MdkitFootnoteDone           string `json:"mdkitFootnoteDone" i18n:"hint='Editors of a text (mdkit), footnote: done.'"`
	MdkitFootnoteCancel         string `json:"mdkitFootnoteCancel" i18n:"hint='Editors of a text (mdkit), footnote: cancel.'"`
	MdkitMarkdownWrite          string `json:"mdkitMarkdownWrite" i18n:"hint='Editors of a text (mdkit), markdown: write.'"`
	MdkitMarkdownPreview        string `json:"mdkitMarkdownPreview" i18n:"hint='Editors of a text (mdkit), markdown: preview.'"`
	MdkitMarkdownHint           string `json:"mdkitMarkdownHint" i18n:"hint='Editors of a text (mdkit), markdown: hint.'"`
	MdkitMarkdownHeading        string `json:"mdkitMarkdownHeading" i18n:"hint='Editors of a text (mdkit), markdown: heading.'"`
	MdkitMarkdownHeadingLevel   string `json:"mdkitMarkdownHeadingLevel" i18n:"hint='Editors of a text (mdkit), markdown: headingLevel.'"`
	MdkitMarkdownParagraph      string `json:"mdkitMarkdownParagraph" i18n:"hint='Editors of a text (mdkit), markdown: paragraph.'"`
	MdkitMarkdownHelp           string `json:"mdkitMarkdownHelp" i18n:"hint='Editors of a text (mdkit), markdown: help.'"`
	MdkitMarkdownFootnote       string `json:"mdkitMarkdownFootnote" i18n:"hint='Editors of a text (mdkit), markdown: footnote.'"`
	MdkitMarkdownInclude        string `json:"mdkitMarkdownInclude" i18n:"hint='Editors of a text (mdkit), markdown: include.'"`
	MdkitMarkdownShortcuts      string `json:"mdkitMarkdownShortcuts" i18n:"hint='Editors of a text (mdkit), markdown: shortcuts.'"`
	MdkitMarkdownGuide          string `json:"mdkitMarkdownGuide" i18n:"hint='Editors of a text (mdkit), markdown: guide.'"`
	MdkitMarkdownClose          string `json:"mdkitMarkdownClose" i18n:"hint='Editors of a text (mdkit), markdown: close.'"`
	MdkitMarkdownHeadingSmaller string `json:"mdkitMarkdownHeadingSmaller" i18n:"hint='Editors of a text (mdkit), markdown: headingSmaller.'"`
	MdkitMarkdownHeadingBigger  string `json:"mdkitMarkdownHeadingBigger" i18n:"hint='Editors of a text (mdkit), markdown: headingBigger.'"`
	MdkitMarkdownHeadingLevels  string `json:"mdkitMarkdownHeadingLevels" i18n:"hint='Editors of a text (mdkit), markdown: headingLevels.'"`
	MdkitMarkdownCleanBlock     string `json:"mdkitMarkdownCleanBlock" i18n:"hint='Editors of a text (mdkit), markdown: cleanBlock.'"`
	MdkitMarkdownTaskList       string `json:"mdkitMarkdownTaskList" i18n:"hint='Editors of a text (mdkit), markdown: taskList.'"`
	MdkitMarkdownExitFullscreen string `json:"mdkitMarkdownExitFullscreen" i18n:"hint='Editors of a text (mdkit), markdown: exitFullscreen.'"`
	MdkitMarkdownTipPictures    string `json:"mdkitMarkdownTipPictures" i18n:"hint='Editors of a text (mdkit), markdown: tipPictures.'"`
	MdkitMarkdownTipFootnote    string `json:"mdkitMarkdownTipFootnote" i18n:"hint='Editors of a text (mdkit), markdown: tipFootnote.'"`
	MdkitMarkdownTipLists       string `json:"mdkitMarkdownTipLists" i18n:"hint='Editors of a text (mdkit), markdown: tipLists.'"`
	MdkitMarkdownTipTables      string `json:"mdkitMarkdownTipTables" i18n:"hint='Editors of a text (mdkit), markdown: tipTables.'"`
	MdkitMarkdownSideBySide     string `json:"mdkitMarkdownSideBySide" i18n:"hint='Editors of a text (mdkit), markdown: sideBySide.'"`
	MdkitMarkdownPreviewOnly    string `json:"mdkitMarkdownPreviewOnly" i18n:"hint='Editors of a text (mdkit), markdown: previewOnly.'"`
	MdkitMarkdownTable          string `json:"mdkitMarkdownTable" i18n:"hint='Editors of a text (mdkit), markdown: table.'"`
	MdkitMarkdownRule           string `json:"mdkitMarkdownRule" i18n:"hint='Editors of a text (mdkit), markdown: rule.'"`
	MdkitMarkdownFullscreen     string `json:"mdkitMarkdownFullscreen" i18n:"hint='Editors of a text (mdkit), markdown: fullscreen.'"`
	MdkitMarkdownHelpUrl        string `json:"mdkitMarkdownHelpUrl" i18n:"hint='Editors of a text (mdkit), markdown: helpUrl.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription:           "The words of the components of the admin (vuetifyx): buttons, the filter of the lists, the pickers, the code editors, the diff.",
	Back:                        "Back",
	Save:                        "Save",
	Apply:                       "Apply",
	Cancel:                      "Cancel",
	Ok:                          "OK",
	Close:                       "Close",
	Clear:                       "Clear",
	Undo:                        "Undo",
	Redo:                        "Redo",
	Move:                        "Move",
	ChangeOrder:                 "Change order",
	Select:                      "Select",
	Home:                        "Home",
	Preview:                     "Preview",
	Help:                        "Help",
	Expand:                      "Expand",
	Collapse:                    "Collapse",
	Settings:                    "Settings",
	Columns:                     "Columns",
	ThemeLight:                  "Light",
	ThemeDark:                   "Dark",
	Year:                        "Year",
	Month:                       "Month",
	FilterClearAll:              "Clear Filters",
	FilterAdd:                   "Add Filters",
	FilterBy:                    "Filter by {filter}",
	FilterTo:                    "to",
	FilterNumberEquals:          "is equal to",
	FilterNumberBetween:         "between",
	FilterNumberGreaterThan:     "is greater than",
	FilterNumberLessThan:        "is less than",
	FilterNumberAnd:             "and",
	FilterStringEquals:          "is equal to",
	FilterStringContains:        "contains",
	FilterIn:                    "in",
	FilterNotIn:                 "not in",
	Copy:                        "Copy",
	Copied:                      "Copied!",
	CopyError:                   "Press Ctrl+C to copy",
	Paste:                       "Paste",
	CopyOrPasteContent:          "Copy or paste content",
	CopyContent:                 "Copy the content to the clipboard",
	PasteContent:                "Paste the content of the clipboard",
	CopiedToClipboard:           "Copied to the clipboard",
	PastedFromClipboard:         "Pasted from the clipboard",
	InvalidContent:              "This is not a valid content",
	PastedIsEmpty:               "Nothing to paste",
	Tools:                       "Tools",
	PleaseEnterTheCode:          "Please enter the code",
	Script:                      "Script",
	ScriptResult:                "Script Result",
	ShowDetail:                  "Show detail",
	InsertTextBlockBefore:       "Insert text block before",
	InsertTextBlockAfter:        "Insert text block after",
	PreFormattedText:            "Pre formatted text",
	SourceCodeHtml:              "Code (HTML)",
	SelectImage:                 "Select Image",
	ImageAlt:                    "Alternative text",
	LockAspectRatio:             "Lock original aspect ratio",
	HeaderPlaceholder:           "Enter a header",
	ImageLoadError:              "An error occurred on load image.",
	ImageReleaseHere:            "Release here...",
	ImageUnsupported:            "Unsupported file type.",
	ImageAddFromLink:            "Add from link",
	ImageTakePicture:            "Take picture from camera",
	ImageEdit:                   "Edit Picture",
	ImageAddFromFile:            "Add from local file",
	ImageFullView:               "Full view",
	TakePicture:                 "Take picture",
	FlipCamera:                  "Flip camera",
	LoadCameras:                 "Load cameras",
	CameraDenied:                "Permission to use the camera was denied: allow it in the browser.",
	CameraNotFound:              "No camera found.",
	CameraError:                 "The camera could not be opened: {error}",
	MapLoadError:                "The map could not be loaded.",
	MapsLoadError:               "Could not load Google Maps.",
	MapsNoKey:                   "Set the Google Maps API Key in SEOConfig.",
	SeoName:                     "Name",
	SeoValue:                    "Value",
	SeoAdd:                      "Add",
	SeoText:                     "Text",
	SeoOneLine:                  "One line",
	SeoMultiline:                "Multiple lines",
	SeoMap:                      "Map",
	SeoPickPlaces:               "Select places on the map",
	SeoZipcodesTitle:            "Select ZIP codes on the map",
	MapsTest:                    "Test",
	MapsTestTitle:               "Test of the Google Maps key",
	MapsTestNoKey:               "Type the key first.",
	MapsTestMaps:                "Maps JavaScript API",
	MapsTestPlaces:              "Places API (New) (search of places)",
	MapsTestGeocoding:           "Geocoding API (addresses)",
	MapsTestWaiting:             "Testing…",
	MapsTestOk:                  "Working",
	MapsTestFailed:              "Failed",
	MapsTestTryIt:               "Search an address: it shows on the map.",
	MapsTestReload:              "The page already loaded Google Maps with another key: save the settings and reload the page to test this one.",
	MapsTestRejected:            "Google refused the key: {code}.",
	MapsErrRefererNotAllowed:    "The key does not allow this site: add its address to the website restrictions of the key.",
	MapsErrApiNotActivated:      "The Maps JavaScript API is not enabled in the project of the key.",
	MapsErrInvalidKey:           "The key is not valid.",
	MapsErrBillingNotEnabled:    "The project of the key has no billing account: enable billing in the Google Cloud console.",
	MapsErrExpiredKey:           "The key expired.",
	MapsErrZeroResults:          "Google found nothing for the search.",
	MapsErrRequestDenied:        "The API is not enabled in the project of the key, or the key does not allow it.",
	MapsErrOverQueryLimit:       "The quota of the key is over.",
	Format:                      "Format",
	Structure:                   "Structure",
	Doc:                         "Documentation",
	Refresh:                     "Refresh",
	NoStructure:                 "No declarations.",
	NoDoc:                       "No documentation.",
	NotFormatted:                "The code could not be formatted: fix what is wrong in it.",
	Compare:                     "Compare with the initial value",
	Edit:                        "Back to the editor",
	PrevChange:                  "Previous change",
	NextChange:                  "Next change",
	Revert:                      "Revert",
	DiffFiles:                   "Files",
	DiffOld:                     "Old",
	DiffCurrent:                 "Current",
	DiffBinary:                  "A binary file: it is not compared.",
	DiffNoChanges:               "No changes.",
	DiffRenamedTo:               "Renamed to",
	DiffUnchanged:               "Its content did not change.",
	Loading:                     "Loading…",
	DiffLoadError:               "The diff could not be loaded.",
	Saving:                      "Saving…",
	Saved:                       "Saved",
	Unsaved:                     "Not saved",
	DiffRevert:                  "Revert: the old part into the current",
	DiffPrev:                    "Previous change (Shift+F7)",
	DiffNext:                    "Next change (F7)",
	MdkitRichBold:               "Bold",
	MdkitRichItalic:             "Italic",
	MdkitRichStrike:             "Strikethrough",
	MdkitRichHeading:            "Heading",
	MdkitRichSubheading:         "Subheading",
	MdkitRichBullets:            "List",
	MdkitRichNumbers:            "Numbered list",
	MdkitRichQuote:              "Quote",
	MdkitRichCode:               "Code",
	MdkitRichLink:               "Link",
	MdkitRichLinkPrompt:         "Address of the link (empty: none)",
	MdkitRichLinkSet:            "OK",
	MdkitRichCancel:             "Cancel",
	MdkitRichImage:              "Picture",
	MdkitRichUndo:               "Undo",
	MdkitRichRedo:               "Redo",
	MdkitRichNotImage:           "Only pictures (PNG, JPEG, GIF, WebP).",
	MdkitRichUploadError:        "The picture could not be sent.",
	MdkitRichPlaceholder:        "Describe it. Paste or drop pictures here.",
	MdkitRichPlaceholderText:    "Describe it.",
	MdkitRichFootnote:           "Footnote",
	MdkitRichSizeSmall:          "Small",
	MdkitRichSizeMedium:         "Medium",
	MdkitRichSizeLarge:          "Large",
	MdkitRichSizeOriginal:       "Original",
	MdkitTableTitle:             "Table",
	MdkitTableRowAbove:          "Row above",
	MdkitTableRowBelow:          "Row below",
	MdkitTableRowRemove:         "Remove rows",
	MdkitTableColLeft:           "Column to the left",
	MdkitTableColRight:          "Column to the right",
	MdkitTableColRemove:         "Remove columns",
	MdkitTableMerge:             "Merge cells",
	MdkitTableUnmerge:           "Unmerge cells",
	MdkitTableAlignLeft:         "Align left",
	MdkitTableAlignCenter:       "Center",
	MdkitTableAlignRight:        "Align right",
	MdkitTableHeader:            "Header row",
	MdkitTableFootnote:          "Footnote",
	MdkitTablePreview:           "Preview",
	MdkitTableResize:            "Drag to size the column; a double click: automatic",
	MdkitTableWidthFixed:        "Widths in pixels",
	MdkitTableWidthPercent:      "Widths in parts of the table",
	MdkitTableWidthAuto:         "Automatic widths",
	MdkitTableFluid:             "Fluid table (as wide as the text)",
	MdkitTablePositionLeft:      "Table on the left",
	MdkitTablePositionCenter:    "Table in the center",
	MdkitTablePositionRight:     "Table on the right",
	MdkitTableValignTop:         "Content at the top",
	MdkitTableValignMiddle:      "Content in the middle",
	MdkitTableValignBottom:      "Content at the bottom",
	MdkitTableCellPlace:         "Place the content of the cells",
	MdkitTableExpand:            "Expand to the whole screen",
	MdkitTableShrink:            "Back to the dialog",
	MdkitTablePaste:             "Paste CSV, Excel or HTML",
	MdkitTablePasteTitle:        "Paste a table",
	MdkitTablePasteHint:         "Paste here (Ctrl+V) the cells copied from a spreadsheet (Excel, LibreOffice) or a page, or a CSV text; or open a CSV file.",
	MdkitTablePastePlaceholder:  "Ctrl+V",
	MdkitTablePasteFile:         "Open a CSV file",
	MdkitTableSeparator:         "Column separator",
	MdkitTableSepAuto:           "Automatic",
	MdkitTableSepTab:            "Tab",
	MdkitTableSepComma:          "Comma ,",
	MdkitTableSepSemicolon:      "Semicolon ;",
	MdkitTableSepPipe:           "Vertical bar |",
	MdkitTableSepOther:          "Other",
	MdkitTableSepOtherChar:      "Separator",
	MdkitTablePasteHTML:         "A table of HTML: its cells as they are (no separator).",
	MdkitTablePasteReplace:      "Replace the table",
	MdkitTablePasteAt:           "From the chosen cell",
	MdkitTablePasteFound:        "{rows} × {cols} cells",
	MdkitTablePasteDo:           "Paste",
	MdkitTableModeSheet:         "Sheet",
	MdkitTableModeSplit:         "Side by side",
	MdkitTableModePreview:       "Preview only",
	MdkitTableHtmlNote:          "With merged cells, or no header row, the table is written as HTML.",
	MdkitTableInsert:            "Insert",
	MdkitTableUpdate:            "Update",
	MdkitTableCancel:            "Cancel",
	MdkitFootnoteTitle:          "Footnotes",
	MdkitFootnoteAdd:            "New note",
	MdkitFootnoteEmpty:          "No footnotes yet.",
	MdkitFootnoteNoNote:         "No note",
	MdkitFootnoteNoRef:          "No reference",
	MdkitFootnoteRefs:           "References",
	MdkitFootnoteLabel:          "Label",
	MdkitFootnoteText:           "The note",
	MdkitFootnoteRename:         "Rename",
	MdkitFootnoteRemove:         "Delete",
	MdkitFootnoteRemoveAsk:      "Delete this note and its references?",
	MdkitFootnoteRemoveYes:      "Delete",
	MdkitFootnoteInsert:         "Insert the reference",
	MdkitFootnoteUseHere:        "Use it in this reference",
	MdkitFootnoteHere:           "This reference",
	MdkitFootnoteSave:           "Keep the note",
	MdkitFootnoteChoose:         "Choose a note, or make a new one.",
	MdkitFootnoteTaken:          "This label is taken.",
	MdkitFootnoteModeWrite:      "Write",
	MdkitFootnoteModeSplit:      "Side by side",
	MdkitFootnoteModePreview:    "Preview only",
	MdkitFootnoteDone:           "Done",
	MdkitFootnoteCancel:         "Cancel",
	MdkitMarkdownWrite:          "Write",
	MdkitMarkdownPreview:        "Preview",
	MdkitMarkdownHint:           "Markdown: **bold**, _italic_, lists, links, `code`.",
	MdkitMarkdownHeading:        "Heading",
	MdkitMarkdownHeadingLevel:   "Heading {n}",
	MdkitMarkdownParagraph:      "Normal text",
	MdkitMarkdownHelp:           "Markdown help",
	MdkitMarkdownFootnote:       "Footnote",
	MdkitMarkdownInclude:        "Include Markdown files",
	MdkitMarkdownShortcuts:      "Keyboard shortcuts",
	MdkitMarkdownGuide:          "Markdown guide",
	MdkitMarkdownClose:          "Close",
	MdkitMarkdownHeadingSmaller: "Smaller heading",
	MdkitMarkdownHeadingBigger:  "Bigger heading",
	MdkitMarkdownHeadingLevels:  "Heading 1 … 6",
	MdkitMarkdownCleanBlock:     "Clear the block",
	MdkitMarkdownTaskList:       "Task list",
	MdkitMarkdownExitFullscreen: "Leave full screen",
	MdkitMarkdownTipPictures:    "Paste or drop pictures to put them in the text.",
	MdkitMarkdownTipFootnote:    "Ctrl+click (⌘+click on a Mac) on a reference [^…] opens its note.",
	MdkitMarkdownTipLists:       "In a list, Enter goes on with it; Tab and Shift+Tab indent it.",
	MdkitMarkdownTipTables:      "A table under the cursor (or selected) opens in the table builder.",
	MdkitMarkdownSideBySide:     "Side by side",
	MdkitMarkdownPreviewOnly:    "Preview only",
	MdkitMarkdownTable:          "Table",
	MdkitMarkdownRule:           "Horizontal line",
	MdkitMarkdownFullscreen:     "Full screen",
	MdkitMarkdownHelpUrl:        "https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax",
}

var Messages_pt_BR = &Messages{
	ModuleDescription:           "As palavras dos componentes do admin (vuetifyx): botões, o filtro das listas, os seletores, os editores de código, o diff.",
	Back:                        "Voltar",
	Save:                        "Salvar",
	Apply:                       "Aplicar",
	Cancel:                      "Cancelar",
	Ok:                          "OK",
	Close:                       "Fechar",
	Clear:                       "Limpar",
	Undo:                        "Desfazer",
	Redo:                        "Refazer",
	Move:                        "Mover",
	ChangeOrder:                 "Alterar ordem",
	Select:                      "Selecionar",
	Home:                        "Início",
	Preview:                     "Pré-visualizar",
	Help:                        "Ajuda",
	Expand:                      "Expandir",
	Collapse:                    "Recolher",
	Settings:                    "Configurações",
	Columns:                     "Colunas",
	ThemeLight:                  "Claro",
	ThemeDark:                   "Escuro",
	Year:                        "Ano",
	Month:                       "Mês",
	FilterClearAll:              "Limpar Filtros",
	FilterAdd:                   "Adicionar Filtro",
	FilterBy:                    "Filtrar por {filter}",
	FilterTo:                    "até",
	FilterNumberEquals:          "é igual a",
	FilterNumberBetween:         "entre",
	FilterNumberGreaterThan:     "é maior que",
	FilterNumberLessThan:        "é menor que",
	FilterNumberAnd:             "e",
	FilterStringEquals:          "é igual a",
	FilterStringContains:        "contém",
	FilterIn:                    "em",
	FilterNotIn:                 "fora de",
	Copy:                        "Copiar",
	Copied:                      "Copiado!",
	CopyError:                   "Tecle Ctrl+C para copiar",
	Paste:                       "Colar",
	CopyOrPasteContent:          "Copiar ou colar conteúdo",
	CopyContent:                 "Copia o conteúdo para a Área de Transferência",
	PasteContent:                "Cola o conteúdo da Área de Transferência",
	CopiedToClipboard:           "Copiado para a Área de Transferência",
	PastedFromClipboard:         "Colado da Área de Transferência",
	InvalidContent:              "Este não é um conteúdo válido",
	PastedIsEmpty:               "Nada a ser colado",
	Tools:                       "Ferramentas",
	PleaseEnterTheCode:          "Por favor insira o código",
	Script:                      "Script",
	ScriptResult:                "Resultado de Script",
	ShowDetail:                  "Ver detalhes",
	InsertTextBlockBefore:       "Inserir bloco de texto antes",
	InsertTextBlockAfter:        "Inserir bloco de texto depois",
	PreFormattedText:            "Texto pré-formatado",
	SourceCodeHtml:              "Código (HTML)",
	SelectImage:                 "Selecionar imagem",
	ImageAlt:                    "Texto alternativo",
	LockAspectRatio:             "Manter a proporção original",
	HeaderPlaceholder:           "Digite um título",
	ImageLoadError:              "Ocorreu um erro ao carregar a imagem.",
	ImageReleaseHere:            "Solte aqui...",
	ImageUnsupported:            "Tipo de arquivo não suportado.",
	ImageAddFromLink:            "Adicionar de um link",
	ImageTakePicture:            "Tirar foto pela câmera",
	ImageEdit:                   "Editar imagem",
	ImageAddFromFile:            "Adicionar de um arquivo local",
	ImageFullView:               "Ver inteira",
	TakePicture:                 "Tirar foto",
	FlipCamera:                  "Virar a câmera",
	LoadCameras:                 "Carregar câmeras",
	CameraDenied:                "A permissão para usar a câmera foi negada: libere-a no navegador.",
	CameraNotFound:              "Nenhuma câmera encontrada.",
	CameraError:                 "A câmera não pôde ser aberta: {error}",
	MapLoadError:                "O mapa não pôde ser carregado.",
	MapsLoadError:               "Não foi possível carregar o Google Maps.",
	MapsNoKey:                   "Configure a Google Maps API Key em SEOConfig.",
	SeoName:                     "Nome",
	SeoValue:                    "Valor",
	SeoAdd:                      "Adicionar",
	SeoText:                     "Texto",
	SeoOneLine:                  "Uma linha",
	SeoMultiline:                "Múltiplas linhas",
	SeoMap:                      "Mapa",
	SeoPickPlaces:               "Selecione lugares no mapa",
	SeoZipcodesTitle:            "Selecionar CEPs no mapa",
	MapsTest:                    "Testar",
	MapsTestTitle:               "Teste da chave do Google Maps",
	MapsTestNoKey:               "Digite a chave primeiro.",
	MapsTestMaps:                "Maps JavaScript API",
	MapsTestPlaces:              "Places API (New) (busca de lugares)",
	MapsTestGeocoding:           "Geocoding API (endereços)",
	MapsTestWaiting:             "Testando…",
	MapsTestOk:                  "Funcionando",
	MapsTestFailed:              "Falhou",
	MapsTestTryIt:               "Busque um endereço: ele aparece no mapa.",
	MapsTestReload:              "A página já carregou o Google Maps com outra chave: salve as configurações e recarregue a página para testar esta.",
	MapsTestRejected:            "O Google recusou a chave: {code}.",
	MapsErrRefererNotAllowed:    "A chave não permite este site: inclua o endereço dele nas restrições de sites da chave.",
	MapsErrApiNotActivated:      "A Maps JavaScript API não está ativada no projeto da chave.",
	MapsErrInvalidKey:           "A chave não é válida.",
	MapsErrBillingNotEnabled:    "O projeto da chave não tem conta de faturamento: ative o faturamento no console do Google Cloud.",
	MapsErrExpiredKey:           "A chave expirou.",
	MapsErrZeroResults:          "O Google não encontrou nada para a busca.",
	MapsErrRequestDenied:        "A API não está ativada no projeto da chave, ou a chave não a permite.",
	MapsErrOverQueryLimit:       "A cota da chave acabou.",
	Format:                      "Formatar",
	Structure:                   "Estrutura",
	Doc:                         "Documentação",
	Refresh:                     "Atualizar",
	NoStructure:                 "Nenhuma declaração.",
	NoDoc:                       "Nenhuma documentação.",
	NotFormatted:                "O código não pôde ser formatado: corrija o que está errado nele.",
	Compare:                     "Comparar com o valor inicial",
	Edit:                        "Voltar ao editor",
	PrevChange:                  "Mudança anterior",
	NextChange:                  "Próxima mudança",
	Revert:                      "Reverter",
	DiffFiles:                   "Arquivos",
	DiffOld:                     "Antigo",
	DiffCurrent:                 "Atual",
	DiffBinary:                  "Um arquivo binário: não é comparado.",
	DiffNoChanges:               "Nenhuma alteração.",
	DiffRenamedTo:               "Renomeado para",
	DiffUnchanged:               "O conteúdo não mudou.",
	Loading:                     "Carregando…",
	DiffLoadError:               "O diff não pôde ser carregado.",
	Saving:                      "Salvando…",
	Saved:                       "Salvo",
	Unsaved:                     "Não salvo",
	DiffRevert:                  "Reverter: a parte antiga na atual",
	DiffPrev:                    "Alteração anterior (Shift+F7)",
	DiffNext:                    "Próxima alteração (F7)",
	MdkitRichBold:               "Negrito",
	MdkitRichItalic:             "Itálico",
	MdkitRichStrike:             "Tachado",
	MdkitRichHeading:            "Título",
	MdkitRichSubheading:         "Subtítulo",
	MdkitRichBullets:            "Lista",
	MdkitRichNumbers:            "Lista numerada",
	MdkitRichQuote:              "Citação",
	MdkitRichCode:               "Código",
	MdkitRichLink:               "Link",
	MdkitRichLinkPrompt:         "Endereço do link (vazio: nenhum)",
	MdkitRichLinkSet:            "OK",
	MdkitRichCancel:             "Cancelar",
	MdkitRichImage:              "Imagem",
	MdkitRichUndo:               "Desfazer",
	MdkitRichRedo:               "Refazer",
	MdkitRichNotImage:           "Só imagens (PNG, JPEG, GIF, WebP).",
	MdkitRichUploadError:        "A imagem não pôde ser enviada.",
	MdkitRichPlaceholder:        "Descreva. Cole ou solte imagens aqui.",
	MdkitRichPlaceholderText:    "Descreva.",
	MdkitRichFootnote:           "Nota de rodapé",
	MdkitRichSizeSmall:          "Pequena",
	MdkitRichSizeMedium:         "Média",
	MdkitRichSizeLarge:          "Grande",
	MdkitRichSizeOriginal:       "Original",
	MdkitTableTitle:             "Tabela",
	MdkitTableRowAbove:          "Linha acima",
	MdkitTableRowBelow:          "Linha abaixo",
	MdkitTableRowRemove:         "Remover linhas",
	MdkitTableColLeft:           "Coluna à esquerda",
	MdkitTableColRight:          "Coluna à direita",
	MdkitTableColRemove:         "Remover colunas",
	MdkitTableMerge:             "Mesclar células",
	MdkitTableUnmerge:           "Desfazer a mesclagem",
	MdkitTableAlignLeft:         "Alinhar à esquerda",
	MdkitTableAlignCenter:       "Centralizar",
	MdkitTableAlignRight:        "Alinhar à direita",
	MdkitTableHeader:            "Linha de cabeçalho",
	MdkitTableFootnote:          "Nota de rodapé",
	MdkitTablePreview:           "Prévia",
	MdkitTableResize:            "Arraste para dimensionar a coluna; duplo clique: automática",
	MdkitTableWidthFixed:        "Larguras em pixels",
	MdkitTableWidthPercent:      "Larguras em partes da tabela",
	MdkitTableWidthAuto:         "Larguras automáticas",
	MdkitTableFluid:             "Tabela fluida (da largura do texto)",
	MdkitTablePositionLeft:      "Tabela à esquerda",
	MdkitTablePositionCenter:    "Tabela ao centro",
	MdkitTablePositionRight:     "Tabela à direita",
	MdkitTableValignTop:         "Conteúdo no alto",
	MdkitTableValignMiddle:      "Conteúdo no meio",
	MdkitTableValignBottom:      "Conteúdo embaixo",
	MdkitTableCellPlace:         "Posição do conteúdo das células",
	MdkitTableExpand:            "Expandir para a tela inteira",
	MdkitTableShrink:            "Voltar ao diálogo",
	MdkitTablePaste:             "Colar CSV, Excel ou HTML",
	MdkitTablePasteTitle:        "Colar uma tabela",
	MdkitTablePasteHint:         "Cole aqui (Ctrl+V) as células copiadas de uma planilha (Excel, LibreOffice) ou de uma página, ou um texto CSV; ou abra um arquivo CSV.",
	MdkitTablePastePlaceholder:  "Ctrl+V",
	MdkitTablePasteFile:         "Abrir um arquivo CSV",
	MdkitTableSeparator:         "Separador de colunas",
	MdkitTableSepAuto:           "Automático",
	MdkitTableSepTab:            "Tabulação",
	MdkitTableSepComma:          "Vírgula ,",
	MdkitTableSepSemicolon:      "Ponto e vírgula ;",
	MdkitTableSepPipe:           "Barra vertical |",
	MdkitTableSepOther:          "Outro",
	MdkitTableSepOtherChar:      "Separador",
	MdkitTablePasteHTML:         "Uma tabela HTML: as células como estão (sem separador).",
	MdkitTablePasteReplace:      "Substituir a tabela",
	MdkitTablePasteAt:           "A partir da célula escolhida",
	MdkitTablePasteFound:        "{rows} × {cols} células",
	MdkitTablePasteDo:           "Colar",
	MdkitTableModeSheet:         "Planilha",
	MdkitTableModeSplit:         "Lado a lado",
	MdkitTableModePreview:       "Só a prévia",
	MdkitTableHtmlNote:          "Com células mescladas, larguras ou sem linha de cabeçalho, a tabela é escrita em HTML.",
	MdkitTableInsert:            "Inserir",
	MdkitTableUpdate:            "Atualizar",
	MdkitTableCancel:            "Cancelar",
	MdkitFootnoteTitle:          "Notas de rodapé",
	MdkitFootnoteAdd:            "Nova nota",
	MdkitFootnoteEmpty:          "Nenhuma nota ainda.",
	MdkitFootnoteNoNote:         "Sem nota",
	MdkitFootnoteNoRef:          "Sem referência",
	MdkitFootnoteRefs:           "Referências",
	MdkitFootnoteLabel:          "Rótulo",
	MdkitFootnoteText:           "A nota",
	MdkitFootnoteRename:         "Renomear",
	MdkitFootnoteRemove:         "Apagar",
	MdkitFootnoteRemoveAsk:      "Apagar esta nota e suas referências?",
	MdkitFootnoteRemoveYes:      "Apagar",
	MdkitFootnoteInsert:         "Inserir a referência",
	MdkitFootnoteUseHere:        "Usar nesta referência",
	MdkitFootnoteHere:           "Esta referência",
	MdkitFootnoteSave:           "Guardar a nota",
	MdkitFootnoteChoose:         "Escolha uma nota, ou crie uma.",
	MdkitFootnoteTaken:          "Este rótulo já é usado.",
	MdkitFootnoteModeWrite:      "Escrever",
	MdkitFootnoteModeSplit:      "Lado a lado",
	MdkitFootnoteModePreview:    "Só a prévia",
	MdkitFootnoteDone:           "Concluir",
	MdkitFootnoteCancel:         "Cancelar",
	MdkitMarkdownWrite:          "Escrever",
	MdkitMarkdownPreview:        "Prévia",
	MdkitMarkdownHint:           "Markdown: **negrito**, _itálico_, listas, links, `código`.",
	MdkitMarkdownHeading:        "Título",
	MdkitMarkdownHeadingLevel:   "Título {n}",
	MdkitMarkdownParagraph:      "Texto normal",
	MdkitMarkdownHelp:           "Ajuda do Markdown",
	MdkitMarkdownFootnote:       "Nota de rodapé",
	MdkitMarkdownInclude:        "Incluir arquivos Markdown",
	MdkitMarkdownShortcuts:      "Atalhos de teclado",
	MdkitMarkdownGuide:          "Guia do Markdown",
	MdkitMarkdownClose:          "Fechar",
	MdkitMarkdownHeadingSmaller: "Título menor",
	MdkitMarkdownHeadingBigger:  "Título maior",
	MdkitMarkdownHeadingLevels:  "Título 1 … 6",
	MdkitMarkdownCleanBlock:     "Limpar o bloco",
	MdkitMarkdownTaskList:       "Lista de tarefas",
	MdkitMarkdownExitFullscreen: "Sair da tela cheia",
	MdkitMarkdownTipPictures:    "Cole ou solte imagens para pô-las no texto.",
	MdkitMarkdownTipFootnote:    "Ctrl+clique (⌘+clique no Mac) numa referência [^…] abre a sua nota.",
	MdkitMarkdownTipLists:       "Numa lista, Enter continua; Tab e Shift+Tab a recuam.",
	MdkitMarkdownTipTables:      "Uma tabela sob o cursor (ou selecionada) abre no construtor de tabelas.",
	MdkitMarkdownSideBySide:     "Lado a lado",
	MdkitMarkdownPreviewOnly:    "Só a prévia",
	MdkitMarkdownTable:          "Tabela",
	MdkitMarkdownRule:           "Linha horizontal",
	MdkitMarkdownFullscreen:     "Tela cheia",
	MdkitMarkdownHelpUrl:        "https://docs.github.com/pt/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax",
}

// VueMessages are the words of the vue-i18n's namespace vx in ctx's
// language: the module's, and the names of the months — the i18n's default
// module's (i18n:default: one place for them, translatable) —, January first.
func VueMessages(ctx context.Context) any {
	d := i18n.GetMessages(ctx)
	var months, abbr [12]string
	copy(months[:], d.MonthNames[1:])
	copy(abbr[:], d.AbbrMonthNames[1:])
	return &struct {
		*Messages
		MonthNames     [12]string `json:"monthNames"`
		AbbrMonthNames [12]string `json:"abbrMonthNames"`
	}{GetMessages(ctx), months, abbr}
}
