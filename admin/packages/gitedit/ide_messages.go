package gitedit

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// IdeMessagesKey is the module of the words of the gad IDE's Changes and Git
// panels (and of the diff browser in them), as the IDE is given them
// (GadIde's messages: IdeMessages of @gad-lang/ide-vuetify — the same keys).
const IdeMessagesKey i18n.ModuleKey = "rvq-admin/gitedit-ide"

// GetIdeMessages are the words of the IDE's panels in ctx's language.
func GetIdeMessages(ctx context.Context) *IdeMessages {
	return i18n.MustGetModuleMessages(ctx, IdeMessagesKey, IdeMessages_en_US).(*IdeMessages)
}

// ConfigureIdeMessages registers the words of the IDE's panels.
func ConfigureIdeMessages(b *i18n.Builder) {
	b.RegisterForModule(language.English, IdeMessagesKey, IdeMessages_en_US).
		RegisterForModule(language.BrazilianPortuguese, IdeMessagesKey, IdeMessages_pt_BR)
}

// IdeMessages are the texts of the IDE — its own (gadIde) and its Changes and
// Git panels' (gadGit) —: their JSON goes to the admin's vue-i18n (gadide
// merges it under both namespaces, in the page's locale).
type IdeMessages struct {
	ModuleDescription   string `json:"-" i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Changes             string `json:"changes" i18n:"hint='Title of the panel of the IDE with the files changed and not committed (its tab and its button in the header).'"`
	Git                 string `json:"git" i18n:"hint='Title of the panel of the IDE with the branches and their commits.'"`
	Preview             string `json:"preview" i18n:"hint='Button of the header of the IDE that opens the panel of the preview of the file.'"`
	Settings            string `json:"settings" i18n:"hint='Button of the header of the IDE that opens its settings.'"`
	Refresh             string `json:"refresh" i18n:"hint='Title of the button that reads the changes, or the branches, again.'"`
	Files               string `json:"files" i18n:"hint='Title of the tree of the files of a diff.'"`
	Branches            string `json:"branches" i18n:"hint='Title of the list of the local branches.'"`
	RemoteBranches      string `json:"remoteBranches" i18n:"hint='Title of the list of the branches of the remotes.'"`
	More                string `json:"more" i18n:"hint='Button at the end of the log that loads more commits.'"`
	ChooseCommit        string `json:"chooseCommit" i18n:"hint='Shown in the Git panel with no commit chosen.'"`
	Patch               string `json:"patch" i18n:"hint='Button that downloads the patch of the commit chosen.'"`
	DownloadCommitPatch string `json:"downloadCommitPatch" i18n:"hint='Title of the button that downloads the patch of the commit chosen.'"`
	DownloadFile        string `json:"downloadFile" i18n:"hint='Title of the button of the diff of a file that downloads it as the commit has it.'"`
	DownloadPatch       string `json:"downloadPatch" i18n:"hint='Title of the button of the diff of a file that downloads its patch from the commit before.'"`
	NoParent            string `json:"noParent" i18n:"hint='Title of the old side of a file of the first commit (none before it).'"`
	Expand              string `json:"expand" i18n:"hint='Title of the button that expands the panel over the others.'"`
	Restore             string `json:"restore" i18n:"hint='Title of the button that brings the panel expanded back.'"`
	Old                 string `json:"old" i18n:"hint='Title of the side of a file as it was.'"`
	Current             string `json:"current" i18n:"hint='Title of the side of a file as it is.'"`
	Binary              string `json:"binary" i18n:"hint='Shown for a binary file changed: it is not compared.'"`
	NoChanges           string `json:"noChanges" i18n:"hint='Shown when no file was changed.'"`
	RenamedTo           string `json:"renamedTo" i18n:"hint='Precedes the new path of a file renamed or moved.'"`
	Unchanged           string `json:"unchanged" i18n:"hint='Said of a file renamed or moved whose content did not change.'"`
	Loading             string `json:"loading" i18n:"hint='Shown in the tab of a file while its diff is loaded.'"`
	LoadError           string `json:"loadError" i18n:"hint='Shown in the tab of a file whose diff could not be loaded.'"`
	Undo                string `json:"undo" i18n:"hint='Button of the diff of a file: undoes the last edit of its current side.'"`
	Redo                string `json:"redo" i18n:"hint='Button of the diff of a file: redoes the edit undone.'"`
	Save                string `json:"save" i18n:"hint='Button of the diff of a file: saves its current side as edited.'"`
	Saving              string `json:"saving" i18n:"hint='Said while the file edited in the diff is saved.'"`
	Saved               string `json:"saved" i18n:"hint='Said once the file edited in the diff was saved.'"`
	Unsaved             string `json:"unsaved" i18n:"hint='Said of a file edited in the diff and not saved.'"`
	Revert              string `json:"revert" i18n:"hint='Title of the button of a change: takes the old part into the current.'"`
	Prev                string `json:"prev" i18n:"hint='Title of the button that goes to the previous change of the diff (Shift+F7).'"`
	Next                string `json:"next" i18n:"hint='Title of the button that goes to the next change of the diff (F7).'"`
	Log                 string `json:"log" i18n:"hint='The mode of the Git panel with the log of the branch.'"`
	FileHistory         string `json:"fileHistory" i18n:"hint='The mode of the Git panel with the history of a file.'"`
	Compare             string `json:"compare" i18n:"hint='The mode of the Git panel that compares two refs.'"`
	Branch              string `json:"branch" i18n:"hint='Label of the choice of the branch navigated.'"`
	DefaultBranch       string `json:"defaultBranch" i18n:"hint='Marks the default branch of the repository.'"`
	File                string `json:"file" i18n:"hint='Label of the file whose history is shown.'"`
	UseOpenFile         string `json:"useOpenFile" i18n:"hint='Title of the button that takes the file open in the editor.'"`
	ChooseFile          string `json:"chooseFile" i18n:"hint='Shown with no file chosen for the history.'"`
	CompareFrom         string `json:"compareFrom" i18n:"hint='Label of the side compared from.'"`
	CompareTo           string `json:"compareTo" i18n:"hint='Label of the side compared to.'"`
	Swap                string `json:"swap" i18n:"hint='Title of the button that swaps the sides of the compare.'"`
	WorkingFiles        string `json:"workingFiles" i18n:"hint='The working files, as a side of a compare.'"`
	CompareWith         string `json:"compareWith" i18n:"hint='Label of the choice of what the commit of a history is compared with.'"`
	CompareHint         string `json:"compareHint" i18n:"hint='Says how two commits are compared.'"`
	ChooseRefs          string `json:"chooseRefs" i18n:"hint='Shown in the compare with nothing chosen.'"`
	Filters             string `json:"filters" i18n:"hint='Title of the button and dialog of the filters of the log.'"`
	Search              string `json:"search" i18n:"hint='Label of the search of the log filters.'"`
	SearchHint          string `json:"searchHint" i18n:"hint='Placeholder of the search of the log filters.'"`
	Authors             string `json:"authors" i18n:"hint='Label of the choice of authors of the log filters.'"`
	FileName            string `json:"fileName" i18n:"hint='Label of the filter of the name of a file changed.'"`
	FileNameHint        string `json:"fileNameHint" i18n:"hint='Placeholder of the filter of the name of a file.'"`
	Content             string `json:"content" i18n:"hint='Label of the filter of the content changed.'"`
	ContentHint         string `json:"contentHint" i18n:"hint='Placeholder of the filter of the content changed.'"`
	CommitMessage       string `json:"commitMessage" i18n:"hint='Label of the filter of the message of the commit.'"`
	Period              string `json:"period" i18n:"hint='Label of the filter of the period of the commits.'"`
	Since               string `json:"since" i18n:"hint='Label of the start of the period.'"`
	Until               string `json:"until" i18n:"hint='Label of the end of the period.'"`
	Last7               string `json:"last7" i18n:"hint='Button: the period of the last 7 days.'"`
	Last30              string `json:"last30" i18n:"hint='Button: the period of the last 30 days.'"`
	LastYear            string `json:"lastYear" i18n:"hint='Button: the period of the last year.'"`
	Apply               string `json:"apply" i18n:"hint='Button that applies the filters.'"`
	Clear               string `json:"clear" i18n:"hint='Button that clears the filters.'"`
	NoCommits           string `json:"noCommits" i18n:"hint='Shown when no commit matches the filters.'"`
	Blame               string `json:"blame" i18n:"hint='Title of the button of a file that shows who changed each line last.'"`
	NotCommitted        string `json:"notCommitted" i18n:"hint='Said of a line of the blame not committed yet.'"`
	HelpFilters         string `json:"helpFilters" i18n:"hint='Help of the filters of the log: a paragraph a line, lines starting with - a list.'"`
	Explorer            string `json:"explorer" i18n:"hint='Title of the panel of the files.'"`
	Editor              string `json:"editor" i18n:"hint='Title of the panel of the editor.'"`
	Structure           string `json:"structure" i18n:"hint='Title of the panel of the declarations of the code.'"`
	Output              string `json:"output" i18n:"hint='Title of the panel of the output of a run.'"`
	CallStack           string `json:"callStack" i18n:"hint='Title of the panel of the call stack of the debugger.'"`
	Locals              string `json:"locals" i18n:"hint='Title of the panel of the local variables of the debugger.'"`
	Breakpoints         string `json:"breakpoints" i18n:"hint='Title of the panel of the breakpoints.'"`
	Close               string `json:"close" i18n:"hint='Button or title that closes.'"`
	Cancel              string `json:"cancel" i18n:"hint='Button that cancels a dialog.'"`
	Ok                  string `json:"ok" i18n:"hint='Button that confirms a dialog.'"`
	Format              string `json:"format" i18n:"hint='Button of the editor that formats the code.'"`
	ReloadFromDisk      string `json:"reloadFromDisk" i18n:"hint='Button of the editor that reads the file again.'"`
	Run                 string `json:"run" i18n:"hint='Button of the editor that runs the file.'"`
	Debug               string `json:"debug" i18n:"hint='Button of the editor that debugs the file.'"`
	Restart             string `json:"restart" i18n:"hint='Button of the editor that restarts the debugger.'"`
	CurrentFile         string `json:"currentFile" i18n:"hint='The file open, as what a run runs.'"`
	Profiles            string `json:"profiles" i18n:"hint='Title of the list of the run profiles.'"`
	DeleteProfile       string `json:"deleteProfile" i18n:"hint='Title of the button that deletes a run profile.'"`
	NewProfile          string `json:"newProfile" i18n:"hint='Item that creates a run profile.'"`
	ContinueRun         string `json:"continueRun" i18n:"hint='Button of the debugger: continues.'"`
	StepOver            string `json:"stepOver" i18n:"hint='Button of the debugger: the next line.'"`
	StepIn              string `json:"stepIn" i18n:"hint='Button of the debugger: into the call.'"`
	StepOut             string `json:"stepOut" i18n:"hint='Button of the debugger: out of the function.'"`
	NoImagePreview      string `json:"noImagePreview" i18n:"hint='Shown for an image that cannot be shown.'"`
	SelectOrCreate      string `json:"selectOrCreate" i18n:"hint='Shown in the editor with no file open.'"`
	NoFile              string `json:"noFile" i18n:"hint='The status bar with no file open.'"`
	DecreaseFont        string `json:"decreaseFont" i18n:"hint='Title of the button that makes the font smaller.'"`
	IncreaseFont        string `json:"increaseFont" i18n:"hint='Title of the button that makes the font larger.'"`
	NewFile             string `json:"newFile" i18n:"hint='Button and dialog that create a file.'"`
	NewFolder           string `json:"newFolder" i18n:"hint='Button and dialog that create a folder.'"`
	UploadFiles         string `json:"uploadFiles" i18n:"hint='Button that uploads files.'"`
	ImportFromUrl       string `json:"importFromUrl" i18n:"hint='Button and dialog that download a file of a URL.'"`
	RenameOpenFile      string `json:"renameOpenFile" i18n:"hint='Title of the button that renames the file open.'"`
	MoveOpenFile        string `json:"moveOpenFile" i18n:"hint='Title of the button that moves the file open.'"`
	DeleteOpenFile      string `json:"deleteOpenFile" i18n:"hint='Title of the button that deletes the file open.'"`
	HideHidden          string `json:"hideHidden" i18n:"hint='Title of the button that hides the hidden files.'"`
	ShowHidden          string `json:"showHidden" i18n:"hint='Title of the button that shows the hidden files.'"`
	ResetChanges        string `json:"resetChanges" i18n:"hint='Title of the button that resets the workspace.'"`
	NoStructure         string `json:"noStructure" i18n:"hint='Shown when the file has no structure (not code).'"`
	NoDeclarations      string `json:"noDeclarations" i18n:"hint='Shown when the code declares nothing.'"`
	NoDocumentation     string `json:"noDocumentation" i18n:"hint='Shown when the file has no documentation.'"`
	GoToLine            string `json:"goToLine" i18n:"hint='Title of a breakpoint: goes to its line.'"`
	EditCondition       string `json:"editCondition" i18n:"hint='Title of the button that edits the condition of a breakpoint.'"`
	Remove              string `json:"remove" i18n:"hint='Title of the button that removes a breakpoint.'"`
	All                 string `json:"all" i18n:"hint='Tab of all the breakpoints.'"`
	NoFileOpen          string `json:"noFileOpen" i18n:"hint='Shown with no file open.'"`
	NoBreakpoints       string `json:"noBreakpoints" i18n:"hint='Shown with no breakpoints.'"`
	Evaluate            string `json:"evaluate" i18n:"hint='Label of the field that evaluates an expression in the debugger.'"`
	InspectValue        string `json:"inspectValue" i18n:"hint='Title of the button that inspects a value.'"`
	RunToSeeOutput      string `json:"runToSeeOutput" i18n:"hint='Shown in the panel of the output before a run.'"`
	Panels              string `json:"panels" i18n:"hint='Tab of the settings: the panels shown.'"`
	Formatter           string `json:"formatter" i18n:"hint='Tab of the settings: the formatter.'"`
	Transpile           string `json:"transpile" i18n:"hint='Tab of the settings: the transpiler.'"`
	Template            string `json:"template" i18n:"hint='Tab of the settings: the templates.'"`
	ArrayItemNewLine    string `json:"arrayItemNewLine" i18n:"hint='Option of the formatter.'"`
	DictItemNewLine     string `json:"dictItemNewLine" i18n:"hint='Option of the formatter.'"`
	CallParamsNewLine   string `json:"callParamsNewLine" i18n:"hint='Option of the formatter.'"`
	ForceMultiline      string `json:"forceMultiline" i18n:"hint='Option of the formatter.'"`
	MaxColumns          string `json:"maxColumns" i18n:"hint='Option of the formatter.'"`
	KeepBackup          string `json:"keepBackup" i18n:"hint='Option of the formatter.'"`
	Move                string `json:"move" i18n:"hint='Button of the dialog that moves a file.'"`
	UploadFile          string `json:"uploadFile" i18n:"hint='Title of the dialog that uploads one file.'"`
	UploadNFiles        string `json:"uploadNFiles" i18n:"hint='Title of the dialog that uploads several files (%d: how many).'"`
	DestinationName     string `json:"destinationName" i18n:"hint='Label of the name of an uploaded file.'"`
	ExtractArchiveOf    string `json:"extractArchiveOf" i18n:"hint='Option of an upload of an archive (%s: its kind).'"`
	TargetFolder        string `json:"targetFolder" i18n:"hint='Label of the folder of an upload or import.'"`
	ChooseFolder        string `json:"chooseFolder" i18n:"hint='Title of the button that chooses the folder.'"`
	ReplaceExisting     string `json:"replaceExisting" i18n:"hint='Option of an upload over a file that exists.'"`
	Upload              string `json:"upload" i18n:"hint='Button of the dialog that uploads.'"`
	Url                 string `json:"url" i18n:"hint='Label of the URL of an import.'"`
	ExtractArchive      string `json:"extractArchive" i18n:"hint='Option of an import of an archive.'"`
	ImportAction        string `json:"importAction" i18n:"hint='Button of the dialog that imports.'"`
	NewRunProfile       string `json:"newRunProfile" i18n:"hint='Title of the dialog that creates a run profile.'"`
	Name                string `json:"name" i18n:"hint='Label of a name.'"`
	Arguments           string `json:"arguments" i18n:"hint='Label of the arguments of a run profile.'"`
	ArgumentsHint       string `json:"argumentsHint" i18n:"hint='Hint of the arguments of a run profile.'"`
	GadxTagOutput       string `json:"gadxTagOutput" i18n:"hint='Label of the output of the tags of a gadx run.'"`
	RenderHtml          string `json:"renderHtml" i18n:"hint='An output of the tags of a gadx run.'"`
	Create              string `json:"create" i18n:"hint='Button of a dialog that creates.'"`
	Condition           string `json:"condition" i18n:"hint='Label of the condition of a breakpoint.'"`
	ConditionHint       string `json:"conditionHint" i18n:"hint='Hint of the condition of a breakpoint.'"`
	Disabled            string `json:"disabled" i18n:"hint='Option that disables a breakpoint.'"`
	FileNamePrompt      string `json:"fileNamePrompt" i18n:"hint='Label of the name of a new file.'"`
	FolderName          string `json:"folderName" i18n:"hint='Label of the name of a new folder.'"`
	DeleteFile          string `json:"deleteFile" i18n:"hint='Title of the dialog that deletes a file.'"`
	DeleteConfirm       string `json:"deleteConfirm" i18n:"hint='Question of the dialog that deletes a file (%s: its path).'"`
	Rename              string `json:"rename" i18n:"hint='Title of the dialog that renames a file.'"`
	NewName             string `json:"newName" i18n:"hint='Label of the new name of a file.'"`
	ResetWorkspace      string `json:"resetWorkspace" i18n:"hint='Title of the dialog that resets the workspace.'"`
	ResetConfirm        string `json:"resetConfirm" i18n:"hint='Question of the dialog that resets the workspace.'"`
	ReplaceOne          string `json:"replaceOne" i18n:"hint='Option of an upload over a file that exists (%s: its name).'"`
	ReplaceMany         string `json:"replaceMany" i18n:"hint='Option of an upload over files that exist (%d: how many).'"`
	Help                string `json:"help" i18n:"hint='Title of the help of a part of a panel.'"`
	HelpChanges         string `json:"helpChanges" i18n:"hint='Help of the Changes panel: a paragraph a line, lines starting with - a list.'"`
	HelpBranches        string `json:"helpBranches" i18n:"hint='Help of the list of branches.'"`
	HelpLog             string `json:"helpLog" i18n:"hint='Help of the log.'"`
	HelpFileHistory     string `json:"helpFileHistory" i18n:"hint='Help of the history of a file.'"`
	HelpCommit          string `json:"helpCommit" i18n:"hint='Help of the detail of what was chosen in the Git panel.'"`
	HelpCompare         string `json:"helpCompare" i18n:"hint='Help of the compare of two refs.'"`
}

var IdeMessages_en_US = &IdeMessages{
	ModuleDescription:   "The Changes and Git panels of the editor of the files (the IDE).",
	Changes:             "Changes",
	Git:                 "Git",
	Preview:             "Preview",
	Settings:            "Settings",
	Refresh:             "Refresh",
	Files:               "Files",
	Branches:            "Branches",
	RemoteBranches:      "Remote",
	More:                "More",
	ChooseCommit:        "Choose a commit.",
	Patch:               "Patch",
	DownloadCommitPatch: "Download the patch of the commit",
	DownloadFile:        "Download the file (this commit's)",
	DownloadPatch:       "Download the patch (from the commit before)",
	NoParent:            "(none)",
	Expand:              "Expand",
	Restore:             "Restore",
	Old:                 "Old",
	Current:             "Current",
	Binary:              "A binary file: it is not compared.",
	NoChanges:           "No changes.",
	RenamedTo:           "Renamed to",
	Unchanged:           "Its content did not change.",
	Loading:             "Loading…",
	LoadError:           "The diff could not be loaded.",
	Undo:                "Undo",
	Redo:                "Redo",
	Save:                "Save",
	Saving:              "Saving…",
	Saved:               "Saved",
	Unsaved:             "Not saved",
	Revert:              "Revert: the old part into the current",
	Prev:                "Previous change (Shift+F7)",
	Next:                "Next change (F7)",
	Log:                 "Log",
	FileHistory:         "File history",
	Compare:             "Compare",
	Branch:              "Branch",
	DefaultBranch:       "default",
	File:                "File",
	UseOpenFile:         "Use the file open in the editor",
	ChooseFile:          "Type or choose a file.",
	CompareFrom:         "From",
	CompareTo:           "To",
	Swap:                "Swap",
	WorkingFiles:        "Working files",
	CompareWith:         "Compare with",
	CompareHint:         "Ctrl+click another commit to compare the two.",
	ChooseRefs:          "Choose what to compare: branches, commits or the working files.",
	Filters:             "Filters",
	Search:              "Search",
	SearchHint:          "In the message, or a hash",
	Authors:             "Authors",
	FileName:            "File name",
	FileNameHint:        "Part of the name, or a pattern (*.gadx)",
	Content:             "Content",
	ContentHint:         "Text of a line added or taken out",
	CommitMessage:       "Commit message",
	Period:              "Period",
	Since:               "From",
	Until:               "To",
	Last7:               "7 days",
	Last30:              "30 days",
	LastYear:            "1 year",
	Apply:               "Apply",
	Clear:               "Clear",
	NoCommits:           "No commit matches.",
	Blame:               "Blame: who changed each line last",
	NotCommitted:        "Not committed yet",
	HelpFilters:         "Chooses the commits shown. The search applies as it is typed; the other filters, by Apply. Case aside.\n- Search: text in the message, or the start of a commit's hash.\n- Authors: one or more (any of them).\n- File name (log only): commits that changed a file whose name has the text, or matches a pattern (*.gadx).\n- Content: commits whose changes added or took out a line with the text.\n- Commit message: text in the whole message.\n- Period: from and to (dates).\nThe filters applied show as chips above the commits — the x of one takes it out —; filtered, the log has no graph.",
	Explorer:            "Explorer",
	Editor:              "Editor",
	Structure:           "Structure",
	Output:              "Output",
	CallStack:           "Call Stack",
	Locals:              "Locals",
	Breakpoints:         "Breakpoints",
	Close:               "Close",
	Cancel:              "Cancel",
	Ok:                  "OK",
	Format:              "Format",
	ReloadFromDisk:      "Reload from disk",
	Run:                 "Run",
	Debug:               "Debug",
	Restart:             "Restart",
	CurrentFile:         "Current file",
	Profiles:            "Profiles",
	DeleteProfile:       "Delete profile",
	NewProfile:          "New profile…",
	ContinueRun:         "Continue",
	StepOver:            "Step Over",
	StepIn:              "Step In",
	StepOut:             "Step Out",
	NoImagePreview:      "No preview of images here.",
	SelectOrCreate:      "Select or create a file to begin.",
	NoFile:              "(no file)",
	DecreaseFont:        "Decrease font size",
	IncreaseFont:        "Increase font size",
	NewFile:             "New file",
	NewFolder:           "New folder",
	UploadFiles:         "Upload files",
	ImportFromUrl:       "Import from URL",
	RenameOpenFile:      "Rename open file",
	MoveOpenFile:        "Move open file",
	DeleteOpenFile:      "Delete open file",
	HideHidden:          "Hide hidden files",
	ShowHidden:          "Show hidden files",
	ResetChanges:        "Reset changes",
	NoStructure:         "No structure here.",
	NoDeclarations:      "No declarations.",
	NoDocumentation:     "No documentation.",
	GoToLine:            "Go to line",
	EditCondition:       "Edit condition",
	Remove:              "Remove",
	All:                 "All",
	NoFileOpen:          "No file open.",
	NoBreakpoints:       "No breakpoints.",
	Evaluate:            "Evaluate",
	InspectValue:        "Inspect value",
	RunToSeeOutput:      "Run or debug to see output.",
	Panels:              "Panels",
	Formatter:           "Formatter",
	Transpile:           "Transpile",
	Template:            "Template",
	ArrayItemNewLine:    "Each array item on its own line",
	DictItemNewLine:     "Each dict item on its own line",
	CallParamsNewLine:   "Each call argument on its own line",
	ForceMultiline:      "Force the full multi-line layout",
	MaxColumns:          "Max columns (blank uses the default)",
	KeepBackup:          "Keep .backup on format",
	Move:                "Move",
	UploadFile:          "Upload file",
	UploadNFiles:        "Upload %d files",
	DestinationName:     "Destination name",
	ExtractArchiveOf:    "Extract archive (%s)",
	TargetFolder:        "Target folder",
	ChooseFolder:        "Choose folder",
	ReplaceExisting:     "Replace existing",
	Upload:              "Upload",
	Url:                 "URL",
	ExtractArchive:      "Extract archive (ZIP / TAR / TAR.GZ)",
	ImportAction:        "Import",
	NewRunProfile:       "New run profile",
	Name:                "Name",
	Arguments:           "Arguments",
	ArgumentsHint:       "Press Enter after each argument (e.g. --count=5, a value, --flag)",
	GadxTagOutput:       "GADX tag output",
	RenderHtml:          "Render (HTML)",
	Create:              "Create",
	Condition:           "Condition (Gad expression)",
	ConditionHint:       "Pause only when this is truthy, e.g. i > 10",
	Disabled:            "Disabled",
	FileNamePrompt:      "File name (e.g. hello.gad)",
	FolderName:          "Folder name",
	DeleteFile:          "Delete file",
	DeleteConfirm:       "Delete %s?",
	Rename:              "Rename",
	NewName:             "New name",
	ResetWorkspace:      "Reset workspace",
	ResetConfirm:        "Discard all your changes and restore the samples?",
	ReplaceOne:          "\"%s\" already exists — replace it",
	ReplaceMany:         "%d files already exist — replace them",
	Help:                "Help",
	HelpChanges:         "The files changed and not committed, each compared with its last commit.\n- Click a file in the tree to see its diff.\n- Edit the current side (right) and save it; revert a change with the arrow between the sides.\n- F7 / Shift+F7 go to the next / previous change.",
	HelpBranches:        "The branches of the repository: local ones, then the remotes'. The one marked \"default\" is the repository's default branch, navigated at first.\n- Click a branch to navigate it: the log, or the file's history, shows its commits.\n- The branch can be chosen in the bar too.",
	HelpLog:             "The commits of the branch navigated, newest first, as a graph of their parents.\n- Click a commit to see its message and the files it changed.\n- Ctrl+click (⌘+click) another commit to compare the two: every file changed between them.\n- \"More\" loads older commits.\n- The filter icon chooses the commits shown (search, authors, file name, content, message, period).",
	HelpFileHistory:     "The commits of the branch that changed the file — the one open in the editor, or typed in the bar —, renames followed (R: the commit renamed it).\n- Click a commit to see the file's change in it.\n- Ctrl+click (⌘+click) another commit to compare the file in the two.\n- \"Compare with\" compares the file, as the commit has it, with another branch's or the working one.\n- The filter icon chooses the commits shown.",
	HelpCommit:          "What was chosen: a commit's message and files, or a compare. Each file in the tree opens its diff — old at the left, new at the right —, scrolled together.\n- F7 / Shift+F7 go to the next / previous change.\n- A commit's file can be downloaded, as the commit has it, and its patch.\n- The blame button of a file shows who changed each line last (as IntelliJ's annotate): a click on a line's commit shows it.\n- The edge above the files resizes the message.",
	HelpCompare:         "What changed from a ref to another, file by file: \"From\" and \"To\" take a branch, a commit's hash or the working files (\"To\" only).\n- Type or choose the two refs; the swap button exchanges them.\n- Click a file in the tree to see its diff.",
}

var IdeMessages_pt_BR = &IdeMessages{
	ModuleDescription:   "Os painéis Alterações e Git do editor dos arquivos (a IDE).",
	Changes:             "Alterações",
	Git:                 "Git",
	Preview:             "Visualizar",
	Settings:            "Configurações",
	Refresh:             "Atualizar",
	Files:               "Arquivos",
	Branches:            "Branches",
	RemoteBranches:      "Remotas",
	More:                "Mais",
	ChooseCommit:        "Escolha um commit.",
	Patch:               "Patch",
	DownloadCommitPatch: "Baixar o patch do commit",
	DownloadFile:        "Baixar o arquivo (deste commit)",
	DownloadPatch:       "Baixar o patch (desde o commit anterior)",
	NoParent:            "(nenhum)",
	Expand:              "Expandir",
	Restore:             "Restaurar",
	Old:                 "Antigo",
	Current:             "Atual",
	Binary:              "Um arquivo binário: não é comparado.",
	NoChanges:           "Nenhuma alteração.",
	RenamedTo:           "Renomeado para",
	Unchanged:           "O conteúdo não mudou.",
	Loading:             "Carregando…",
	LoadError:           "Não foi possível carregar o diff.",
	Undo:                "Desfazer",
	Redo:                "Refazer",
	Save:                "Salvar",
	Saving:              "Salvando…",
	Saved:               "Salvo",
	Unsaved:             "Não salvo",
	Revert:              "Desfazer: a parte antiga no atual",
	Prev:                "Alteração anterior (Shift+F7)",
	Next:                "Próxima alteração (F7)",
	Log:                 "Log",
	FileHistory:         "Histórico do arquivo",
	Compare:             "Comparar",
	Branch:              "Branch",
	DefaultBranch:       "padrão",
	File:                "Arquivo",
	UseOpenFile:         "Usar o arquivo aberto no editor",
	ChooseFile:          "Digite ou escolha um arquivo.",
	CompareFrom:         "De",
	CompareTo:           "Para",
	Swap:                "Inverter",
	WorkingFiles:        "Arquivos de trabalho",
	CompareWith:         "Comparar com",
	CompareHint:         "Ctrl+clique em outro commit para comparar os dois.",
	ChooseRefs:          "Escolha o que comparar: branches, commits ou os arquivos de trabalho.",
	Filters:             "Filtros",
	Search:              "Pesquisar",
	SearchHint:          "Na mensagem, ou um hash",
	Authors:             "Autores",
	FileName:            "Nome do arquivo",
	FileNameHint:        "Parte do nome, ou um padrão (*.gadx)",
	Content:             "Conteúdo",
	ContentHint:         "Texto de uma linha adicionada ou removida",
	CommitMessage:       "Mensagem do commit",
	Period:              "Período",
	Since:               "De",
	Until:               "Até",
	Last7:               "7 dias",
	Last30:              "30 dias",
	LastYear:            "1 ano",
	Apply:               "Aplicar",
	Clear:               "Limpar",
	NoCommits:           "Nenhum commit encontrado.",
	Blame:               "Blame: quem alterou cada linha por último",
	NotCommitted:        "Ainda não commitado",
	HelpFilters:         "Escolhe os commits mostrados. A pesquisa é aplicada ao digitar; os demais filtros, por Aplicar. Sem diferenciar maiúsculas.\n- Pesquisar: texto na mensagem, ou o início do hash de um commit.\n- Autores: um ou mais (qualquer deles).\n- Nome do arquivo (só no log): commits que alteraram um arquivo cujo nome tem o texto, ou casa com um padrão (*.gadx).\n- Conteúdo: commits cujas alterações adicionaram ou removeram uma linha com o texto.\n- Mensagem do commit: texto na mensagem inteira.\n- Período: de e até (datas).\nOs filtros aplicados aparecem como chips acima dos commits — o x de um o remove —; filtrado, o log fica sem o grafo.",
	Explorer:            "Explorador",
	Editor:              "Editor",
	Structure:           "Estrutura",
	Output:              "Saída",
	CallStack:           "Pilha de chamadas",
	Locals:              "Variáveis locais",
	Breakpoints:         "Pontos de parada",
	Close:               "Fechar",
	Cancel:              "Cancelar",
	Ok:                  "OK",
	Format:              "Formatar",
	ReloadFromDisk:      "Recarregar do disco",
	Run:                 "Executar",
	Debug:               "Depurar",
	Restart:             "Reiniciar",
	CurrentFile:         "Arquivo atual",
	Profiles:            "Perfis",
	DeleteProfile:       "Excluir perfil",
	NewProfile:          "Novo perfil…",
	ContinueRun:         "Continuar",
	StepOver:            "Passar por cima",
	StepIn:              "Entrar",
	StepOut:             "Sair",
	NoImagePreview:      "Sem visualização de imagens aqui.",
	SelectOrCreate:      "Selecione ou crie um arquivo para começar.",
	NoFile:              "(nenhum arquivo)",
	DecreaseFont:        "Diminuir a fonte",
	IncreaseFont:        "Aumentar a fonte",
	NewFile:             "Novo arquivo",
	NewFolder:           "Nova pasta",
	UploadFiles:         "Enviar arquivos",
	ImportFromUrl:       "Importar de URL",
	RenameOpenFile:      "Renomear o arquivo aberto",
	MoveOpenFile:        "Mover o arquivo aberto",
	DeleteOpenFile:      "Excluir o arquivo aberto",
	HideHidden:          "Ocultar arquivos ocultos",
	ShowHidden:          "Mostrar arquivos ocultos",
	ResetChanges:        "Desfazer as alterações",
	NoStructure:         "Nenhuma estrutura aqui.",
	NoDeclarations:      "Nenhuma declaração.",
	NoDocumentation:     "Nenhuma documentação.",
	GoToLine:            "Ir para a linha",
	EditCondition:       "Editar condição",
	Remove:              "Remover",
	All:                 "Todos",
	NoFileOpen:          "Nenhum arquivo aberto.",
	NoBreakpoints:       "Nenhum ponto de parada.",
	Evaluate:            "Avaliar",
	InspectValue:        "Inspecionar o valor",
	RunToSeeOutput:      "Execute ou depure para ver a saída.",
	Panels:              "Painéis",
	Formatter:           "Formatador",
	Transpile:           "Transpilar",
	Template:            "Template",
	ArrayItemNewLine:    "Cada item de array numa linha",
	DictItemNewLine:     "Cada item de dict numa linha",
	CallParamsNewLine:   "Cada argumento de chamada numa linha",
	ForceMultiline:      "Forçar o layout de várias linhas",
	MaxColumns:          "Máximo de colunas (vazio: o padrão)",
	KeepBackup:          "Manter .backup ao formatar",
	Move:                "Mover",
	UploadFile:          "Enviar arquivo",
	UploadNFiles:        "Enviar %d arquivos",
	DestinationName:     "Nome de destino",
	ExtractArchiveOf:    "Extrair o arquivo compactado (%s)",
	TargetFolder:        "Pasta de destino",
	ChooseFolder:        "Escolher a pasta",
	ReplaceExisting:     "Substituir o existente",
	Upload:              "Enviar",
	Url:                 "URL",
	ExtractArchive:      "Extrair o arquivo compactado (ZIP / TAR / TAR.GZ)",
	ImportAction:        "Importar",
	NewRunProfile:       "Novo perfil de execução",
	Name:                "Nome",
	Arguments:           "Argumentos",
	ArgumentsHint:       "Tecle Enter após cada argumento (ex.: --count=5, um valor, --flag)",
	GadxTagOutput:       "Saída das tags GADX",
	RenderHtml:          "Renderizar (HTML)",
	Create:              "Criar",
	Condition:           "Condição (expressão Gad)",
	ConditionHint:       "Pausa só quando isto for verdadeiro, ex.: i > 10",
	Disabled:            "Desativado",
	FileNamePrompt:      "Nome do arquivo (ex.: hello.gad)",
	FolderName:          "Nome da pasta",
	DeleteFile:          "Excluir arquivo",
	DeleteConfirm:       "Excluir %s?",
	Rename:              "Renomear",
	NewName:             "Novo nome",
	ResetWorkspace:      "Restaurar o espaço de trabalho",
	ResetConfirm:        "Descartar todas as suas alterações e restaurar os exemplos?",
	ReplaceOne:          "\"%s\" já existe — substituí-lo",
	ReplaceMany:         "%d arquivos já existem — substituí-los",
	Help:                "Ajuda",
	HelpChanges:         "Os arquivos alterados e ainda não commitados, cada um comparado com o seu último commit.\n- Clique num arquivo da árvore para ver o diff.\n- Edite o lado atual (à direita) e salve; reverta uma alteração com a seta entre os lados.\n- F7 / Shift+F7 vão para a próxima / anterior alteração.",
	HelpBranches:        "As branches do repositório: as locais, depois as dos remotos. A marcada \"padrão\" é a branch padrão do repositório, navegada ao abrir.\n- Clique numa branch para navegá-la: o log, ou o histórico do arquivo, mostra os commits dela.\n- A branch também pode ser escolhida na barra.",
	HelpLog:             "Os commits da branch navegada, do mais novo ao mais antigo, como um grafo dos seus pais.\n- Clique num commit para ver a mensagem e os arquivos que ele alterou.\n- Ctrl+clique (⌘+clique) em outro commit para comparar os dois: todos os arquivos alterados entre eles.\n- \"Mais\" carrega commits mais antigos.\n- O ícone de filtro escolhe os commits mostrados (pesquisa, autores, nome do arquivo, conteúdo, mensagem, período).",
	HelpFileHistory:     "Os commits da branch que alteraram o arquivo — o aberto no editor, ou digitado na barra —, seguindo renomeações (R: o commit o renomeou).\n- Clique num commit para ver a alteração do arquivo nele.\n- Ctrl+clique (⌘+clique) em outro commit para comparar o arquivo nos dois.\n- \"Comparar com\" compara o arquivo, como o commit o tem, com o de outra branch ou com o de trabalho.\n- O ícone de filtro escolhe os commits mostrados.",
	HelpCommit:          "O que foi escolhido: a mensagem e os arquivos de um commit, ou uma comparação. Cada arquivo da árvore abre o seu diff — o antigo à esquerda, o novo à direita —, com rolagem sincronizada.\n- F7 / Shift+F7 vão para a próxima / anterior alteração.\n- O arquivo de um commit pode ser baixado, como o commit o tem, e também o seu patch.\n- O botão de blame de um arquivo mostra quem alterou cada linha por último (como o annotate do IntelliJ): um clique no commit de uma linha o mostra.\n- A borda acima dos arquivos redimensiona a mensagem.",
	HelpCompare:         "O que mudou de uma ref para outra, arquivo por arquivo: \"De\" e \"Para\" aceitam uma branch, o hash de um commit ou os arquivos de trabalho (só \"Para\").\n- Digite ou escolha as duas refs; o botão de inverter troca os lados.\n- Clique num arquivo da árvore para ver o diff.",
}
