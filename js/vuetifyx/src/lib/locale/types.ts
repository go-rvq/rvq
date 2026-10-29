// VxMessages are the words of the vuetifyx components: the namespace `vx` of
// the page's vue-i18n (the admin merges them, in the user's language, once a
// page: rvq x/ui/vuetifyx Messages), over the defaults (en.ts, pt.ts).
export interface VxMessages {
  /** Button that goes back. */
  back: string
  /** Button that saves. */
  save: string
  /** Button that applies. */
  apply: string
  /** Button that cancels. */
  cancel: string
  /** Button that confirms. */
  ok: string
  /** Button that closes. */
  close: string
  /** Button that clears a value or a filter. */
  clear: string
  /** Button that undoes the last edit. */
  undo: string
  /** Button that redoes the edit undone. */
  redo: string
  /** Button that moves. */
  move: string
  /** Button that changes the order of items. */
  changeOrder: string
  /** Button that selects. */
  select: string
  /** The first page. */
  home: string
  /** Button that previews. */
  preview: string
  /** Button and title of a help. */
  help: string
  /** Button that expands. */
  expand: string
  /** Button that collapses. */
  collapse: string
  /** Title of the settings of a table. */
  settings: string
  /** Title of the columns of a table. */
  columns: string
  /** Title of the button that turns the theme light. */
  themeLight: string
  /** Title of the button that turns the theme dark. */
  themeDark: string
  /** The theme that follows the system's (light or dark). */
  themeAuto: string
  /** Takes a tag off (the trash of a chip). */
  tagRemove: string
  /** No tag to choose from. */
  tagsNone: string
  /** Label of a year (a month picker). */
  year: string
  /** Label of a month (a month picker). */
  month: string
  /** The names of the months, January first (the admin's i18n default
   * module's; out of it, Intl's in the vue-i18n's locale) */
  monthNames: string[]
  /** their short names */
  abbrMonthNames: string[]
  /** Button that removes every filter. */
  filterClearAll: string
  /** Button that adds a filter. */
  filterAdd: string
  /** Title of a filter ({filter}: its name). */
  filterBy: string
  /** Between the two values of a range of a filter. */
  filterTo: string
  /** Modifier of a number filter. */
  filterNumberEquals: string
  /** Modifier of a number filter. */
  filterNumberBetween: string
  /** Modifier of a number filter. */
  filterNumberGreaterThan: string
  /** Modifier of a number filter. */
  filterNumberLessThan: string
  /** Between the two numbers of a number filter. */
  filterNumberAnd: string
  /** Modifier of a text filter. */
  filterStringEquals: string
  /** Modifier of a text filter. */
  filterStringContains: string
  /** Modifier of a choice filter: one of them. */
  filterIn: string
  /** Modifier of a choice filter: none of them. */
  filterNotIn: string
  /** Button that copies. */
  copy: string
  /** Said once copied. */
  copied: string
  /** Said when the copy failed. */
  copyError: string
  /** Button that pastes. */
  paste: string
  /** Menu of the copy and paste of a content. */
  copyOrPasteContent: string
  /** Item that copies the content. */
  copyContent: string
  /** Item that pastes the content. */
  pasteContent: string
  /** Said once copied. */
  copiedToClipboard: string
  /** Said once pasted. */
  pastedFromClipboard: string
  /** Said when what was pasted is not valid. */
  invalidContent: string
  /** Said when the clipboard is empty. */
  pastedIsEmpty: string
  /** Menu of tools. */
  tools: string
  /** Asks for a code. */
  pleaseEnterTheCode: string
  /** A script. */
  script: string
  /** The result of a script. */
  scriptResult: string
  /** Button that shows the detail. */
  showDetail: string
  /** Button of the editor. */
  insertTextBlockBefore: string
  /** Button of the editor. */
  insertTextBlockAfter: string
  /** A pre formatted text. */
  preFormattedText: string
  /** Title of the dialog of the HTML of a content. */
  sourceCodeHtml: string
  /** Label of the choice of an image. */
  selectImage: string
  /** Label of the alternative text of an image. */
  imageAlt: string
  /** Option of an image. */
  lockAspectRatio: string
  /** Placeholder of a header of the block editor. */
  headerPlaceholder: string
  /** Said when an image could not be loaded. */
  imageLoadError: string
  /** Said while a file is dragged over. */
  imageReleaseHere: string
  /** Said of a file not supported. */
  imageUnsupported: string
  /** Button that adds an image of a link. */
  imageAddFromLink: string
  /** Button that takes a picture. */
  imageTakePicture: string
  /** Button that edits the picture. */
  imageEdit: string
  /** Button that adds an image of a file. */
  imageAddFromFile: string
  /** Button that shows the image whole. */
  imageFullView: string
  /** Button of the camera. */
  takePicture: string
  /** Button of the camera. */
  flipCamera: string
  /** Button of the camera. */
  loadCameras: string
  /** Said when the browser denied the camera. */
  cameraDenied: string
  /** Said when there is no camera. */
  cameraNotFound: string
  /** Said when the camera failed ({error}: why). */
  cameraError: string
  /** Said when the map could not be loaded. */
  mapLoadError: string
  /** Said when Google Maps could not be loaded. */
  mapsLoadError: string
  /** Said when there is no key of Google Maps. */
  mapsNoKey: string
  /** Column of the name of a SEO variable. */
  seoName: string
  /** Column of the value of a SEO variable. */
  seoValue: string
  /** Button that adds a SEO variable. */
  seoAdd: string
  /** A SEO variable of text. */
  seoText: string
  /** Title of the button that makes a value one line. */
  seoOneLine: string
  /** Title of the button that makes a value several lines. */
  seoMultiline: string
  /** Button that opens the map. */
  seoMap: string
  /** Placeholder of the places of a ZIP variable. */
  seoPickPlaces: string
  /** Title of the map of the ZIP codes. */
  seoZipcodesTitle: string
  /** Button that tests the key of Google Maps. */
  mapsTest: string
  /** Title of the test of the key. */
  mapsTestTitle: string
  /** Said with no key. */
  mapsTestNoKey: string
  /** An API of Google. */
  mapsTestMaps: string
  /** An API of Google. */
  mapsTestPlaces: string
  /** An API of Google. */
  mapsTestGeocoding: string
  /** Said while testing. */
  mapsTestWaiting: string
  /** Said of an API that works. */
  mapsTestOk: string
  /** Said of an API that failed. */
  mapsTestFailed: string
  /** Hint of the test. */
  mapsTestTryIt: string
  /** Said when another key is loaded. */
  mapsTestReload: string
  /** Said when Google refused the key ({code}: its error). */
  mapsTestRejected: string
  /** An error of Google Maps. */
  mapsErrRefererNotAllowed: string
  /** An error of Google Maps. */
  mapsErrApiNotActivated: string
  /** An error of Google Maps. */
  mapsErrInvalidKey: string
  /** An error of Google Maps. */
  mapsErrBillingNotEnabled: string
  /** An error of Google Maps. */
  mapsErrExpiredKey: string
  /** An error of Google Maps. */
  mapsErrZeroResults: string
  /** An error of Google Maps. */
  mapsErrRequestDenied: string
  /** An error of Google Maps. */
  mapsErrOverQueryLimit: string
  /** Button that formats the code (gad fmt). */
  format: string
  /** Title of the panel of the declarations of the code. */
  structure: string
  /** Title of the panel of the documentation of the code. */
  doc: string
  /** Button that reads again. */
  refresh: string
  /** Said when the code declares nothing. */
  noStructure: string
  /** Said when the code has no documentation. */
  noDoc: string
  /** Said when the code does not parse. */
  notFormatted: string
  /** Button that compares the code with its initial value. */
  compare: string
  /** Button that leaves the compare. */
  edit: string
  /** Button of the compare. */
  prevChange: string
  /** Button of the compare. */
  nextChange: string
  /** Button of a change of the compare: takes the old part back. */
  revert: string
  /** Title of the tree of the files of a diff. */
  diffFiles: string
  /** Title of the old side of a diff. */
  diffOld: string
  /** Title of the current side of a diff. */
  diffCurrent: string
  /** Said of a binary file. */
  diffBinary: string
  /** Said with no changes. */
  diffNoChanges: string
  /** Before the new path of a renamed file. */
  diffRenamedTo: string
  /** Said of a renamed file of the same content. */
  diffUnchanged: string
  /** Said while loading. */
  loading: string
  /** Said when a diff could not be loaded. */
  diffLoadError: string
  /** Said while saving. */
  saving: string
  /** Said once saved. */
  saved: string
  /** Said of a change not saved. */
  unsaved: string
  /** Title of the button of a change of a diff. */
  diffRevert: string
  /** Title of the button of a diff. */
  diffPrev: string
  /** Title of the button of a diff. */
  diffNext: string
  /** Editors of a text (mdkit), rich: bold. */
  mdkitRichBold: string
  /** Editors of a text (mdkit), rich: italic. */
  mdkitRichItalic: string
  /** Editors of a text (mdkit), rich: strike. */
  mdkitRichStrike: string
  /** Editors of a text (mdkit), rich: heading. */
  mdkitRichHeading: string
  /** Editors of a text (mdkit), rich: subheading. */
  mdkitRichSubheading: string
  /** Editors of a text (mdkit), rich: bullets. */
  mdkitRichBullets: string
  /** Editors of a text (mdkit), rich: numbers. */
  mdkitRichNumbers: string
  /** Editors of a text (mdkit), rich: quote. */
  mdkitRichQuote: string
  /** Editors of a text (mdkit), rich: code. */
  mdkitRichCode: string
  /** Editors of a text (mdkit), rich: link. */
  mdkitRichLink: string
  /** Editors of a text (mdkit), rich: linkPrompt. */
  mdkitRichLinkPrompt: string
  /** Editors of a text (mdkit), rich: linkSet. */
  mdkitRichLinkSet: string
  /** Editors of a text (mdkit), rich: cancel. */
  mdkitRichCancel: string
  /** Editors of a text (mdkit), rich: image. */
  mdkitRichImage: string
  /** Editors of a text (mdkit), rich: undo. */
  mdkitRichUndo: string
  /** Editors of a text (mdkit), rich: redo. */
  mdkitRichRedo: string
  /** Editors of a text (mdkit), rich: notImage. */
  mdkitRichNotImage: string
  /** Editors of a text (mdkit), rich: uploadError. */
  mdkitRichUploadError: string
  /** Editors of a text (mdkit), rich: placeholder. */
  mdkitRichPlaceholder: string
  /** Editors of a text (mdkit), rich: placeholderText. */
  mdkitRichPlaceholderText: string
  /** Editors of a text (mdkit), rich: footnote. */
  mdkitRichFootnote: string
  /** Editors of a text (mdkit), rich: size_small. */
  mdkitRichSizeSmall: string
  /** Editors of a text (mdkit), rich: size_medium. */
  mdkitRichSizeMedium: string
  /** Editors of a text (mdkit), rich: size_large. */
  mdkitRichSizeLarge: string
  /** Editors of a text (mdkit), rich: size_original. */
  mdkitRichSizeOriginal: string
  /** Editors of a text (mdkit), table: title. */
  mdkitTableTitle: string
  /** Editors of a text (mdkit), table: rowAbove. */
  mdkitTableRowAbove: string
  /** Editors of a text (mdkit), table: rowBelow. */
  mdkitTableRowBelow: string
  /** Editors of a text (mdkit), table: rowRemove. */
  mdkitTableRowRemove: string
  /** Editors of a text (mdkit), table: colLeft. */
  mdkitTableColLeft: string
  /** Editors of a text (mdkit), table: colRight. */
  mdkitTableColRight: string
  /** Editors of a text (mdkit), table: colRemove. */
  mdkitTableColRemove: string
  /** Editors of a text (mdkit), table: merge. */
  mdkitTableMerge: string
  /** Editors of a text (mdkit), table: unmerge. */
  mdkitTableUnmerge: string
  /** Editors of a text (mdkit), table: alignLeft. */
  mdkitTableAlignLeft: string
  /** Editors of a text (mdkit), table: alignCenter. */
  mdkitTableAlignCenter: string
  /** Editors of a text (mdkit), table: alignRight. */
  mdkitTableAlignRight: string
  /** Editors of a text (mdkit), table: header. */
  mdkitTableHeader: string
  /** Editors of a text (mdkit), table: footnote. */
  mdkitTableFootnote: string
  /** Editors of a text (mdkit), table: preview. */
  mdkitTablePreview: string
  /** Editors of a text (mdkit), table: resize. */
  mdkitTableResize: string
  /** Editors of a text (mdkit), table: widthFixed. */
  mdkitTableWidthFixed: string
  /** Editors of a text (mdkit), table: widthPercent. */
  mdkitTableWidthPercent: string
  /** Editors of a text (mdkit), table: widthAuto. */
  mdkitTableWidthAuto: string
  /** Editors of a text (mdkit), table: fluid. */
  mdkitTableFluid: string
  /** Editors of a text (mdkit), table: positionLeft. */
  mdkitTablePositionLeft: string
  /** Editors of a text (mdkit), table: positionCenter. */
  mdkitTablePositionCenter: string
  /** Editors of a text (mdkit), table: positionRight. */
  mdkitTablePositionRight: string
  /** Editors of a text (mdkit), table: valignTop. */
  mdkitTableValignTop: string
  /** Editors of a text (mdkit), table: valignMiddle. */
  mdkitTableValignMiddle: string
  /** Editors of a text (mdkit), table: valignBottom. */
  mdkitTableValignBottom: string
  /** Editors of a text (mdkit), table: cellPlace. */
  mdkitTableCellPlace: string
  /** Editors of a text (mdkit), table: expand. */
  mdkitTableExpand: string
  /** Editors of a text (mdkit), table: shrink. */
  mdkitTableShrink: string
  /** Editors of a text (mdkit), table: paste. */
  mdkitTablePaste: string
  /** Editors of a text (mdkit), table: pasteTitle. */
  mdkitTablePasteTitle: string
  /** Editors of a text (mdkit), table: pasteHint. */
  mdkitTablePasteHint: string
  /** Editors of a text (mdkit), table: pastePlaceholder. */
  mdkitTablePastePlaceholder: string
  /** Editors of a text (mdkit), table: pasteFile. */
  mdkitTablePasteFile: string
  /** Editors of a text (mdkit), table: separator. */
  mdkitTableSeparator: string
  /** Editors of a text (mdkit), table: sepAuto. */
  mdkitTableSepAuto: string
  /** Editors of a text (mdkit), table: sepTab. */
  mdkitTableSepTab: string
  /** Editors of a text (mdkit), table: sepComma. */
  mdkitTableSepComma: string
  /** Editors of a text (mdkit), table: sepSemicolon. */
  mdkitTableSepSemicolon: string
  /** Editors of a text (mdkit), table: sepPipe. */
  mdkitTableSepPipe: string
  /** Editors of a text (mdkit), table: sepOther. */
  mdkitTableSepOther: string
  /** Editors of a text (mdkit), table: sepOtherChar. */
  mdkitTableSepOtherChar: string
  /** Editors of a text (mdkit), table: pasteHTML. */
  mdkitTablePasteHTML: string
  /** Editors of a text (mdkit), table: pasteReplace. */
  mdkitTablePasteReplace: string
  /** Editors of a text (mdkit), table: pasteAt. */
  mdkitTablePasteAt: string
  /** Editors of a text (mdkit), table: pasteFound. */
  mdkitTablePasteFound: string
  /** Editors of a text (mdkit), table: pasteDo. */
  mdkitTablePasteDo: string
  /** Editors of a text (mdkit), table: modeSheet. */
  mdkitTableModeSheet: string
  /** Editors of a text (mdkit), table: modeSplit. */
  mdkitTableModeSplit: string
  /** Editors of a text (mdkit), table: modePreview. */
  mdkitTableModePreview: string
  /** Editors of a text (mdkit), table: htmlNote. */
  mdkitTableHtmlNote: string
  /** Editors of a text (mdkit), table: insert. */
  mdkitTableInsert: string
  /** Editors of a text (mdkit), table: update. */
  mdkitTableUpdate: string
  /** Editors of a text (mdkit), table: cancel. */
  mdkitTableCancel: string
  /** Editors of a text (mdkit), footnote: title. */
  mdkitFootnoteTitle: string
  /** Editors of a text (mdkit), footnote: add. */
  mdkitFootnoteAdd: string
  /** Editors of a text (mdkit), footnote: empty. */
  mdkitFootnoteEmpty: string
  /** Editors of a text (mdkit), footnote: noNote. */
  mdkitFootnoteNoNote: string
  /** Editors of a text (mdkit), footnote: noRef. */
  mdkitFootnoteNoRef: string
  /** Editors of a text (mdkit), footnote: refs. */
  mdkitFootnoteRefs: string
  /** Editors of a text (mdkit), footnote: label. */
  mdkitFootnoteLabel: string
  /** Editors of a text (mdkit), footnote: text. */
  mdkitFootnoteText: string
  /** Editors of a text (mdkit), footnote: rename. */
  mdkitFootnoteRename: string
  /** Editors of a text (mdkit), footnote: remove. */
  mdkitFootnoteRemove: string
  /** Editors of a text (mdkit), footnote: removeAsk. */
  mdkitFootnoteRemoveAsk: string
  /** Editors of a text (mdkit), footnote: removeYes. */
  mdkitFootnoteRemoveYes: string
  /** Editors of a text (mdkit), footnote: insert. */
  mdkitFootnoteInsert: string
  /** Editors of a text (mdkit), footnote: useHere. */
  mdkitFootnoteUseHere: string
  /** Editors of a text (mdkit), footnote: here. */
  mdkitFootnoteHere: string
  /** Editors of a text (mdkit), footnote: save. */
  mdkitFootnoteSave: string
  /** Editors of a text (mdkit), footnote: choose. */
  mdkitFootnoteChoose: string
  /** Editors of a text (mdkit), footnote: taken. */
  mdkitFootnoteTaken: string
  /** Editors of a text (mdkit), footnote: modeWrite. */
  mdkitFootnoteModeWrite: string
  /** Editors of a text (mdkit), footnote: modeSplit. */
  mdkitFootnoteModeSplit: string
  /** Editors of a text (mdkit), footnote: modePreview. */
  mdkitFootnoteModePreview: string
  /** Editors of a text (mdkit), footnote: done. */
  mdkitFootnoteDone: string
  /** Editors of a text (mdkit), footnote: cancel. */
  mdkitFootnoteCancel: string
  /** Editors of a text (mdkit), markdown: write. */
  mdkitMarkdownWrite: string
  /** Editors of a text (mdkit), markdown: preview. */
  mdkitMarkdownPreview: string
  /** Editors of a text (mdkit), markdown: hint. */
  mdkitMarkdownHint: string
  /** Editors of a text (mdkit), markdown: heading. */
  mdkitMarkdownHeading: string
  /** Editors of a text (mdkit), markdown: headingLevel. */
  mdkitMarkdownHeadingLevel: string
  /** Editors of a text (mdkit), markdown: paragraph. */
  mdkitMarkdownParagraph: string
  /** Editors of a text (mdkit), markdown: help. */
  mdkitMarkdownHelp: string
  /** Editors of a text (mdkit), markdown: footnote. */
  mdkitMarkdownFootnote: string
  /** Editors of a text (mdkit), markdown: include. */
  mdkitMarkdownInclude: string
  /** Editors of a text (mdkit), markdown: shortcuts. */
  mdkitMarkdownShortcuts: string
  /** Editors of a text (mdkit), markdown: guide. */
  mdkitMarkdownGuide: string
  /** Editors of a text (mdkit), markdown: close. */
  mdkitMarkdownClose: string
  /** Editors of a text (mdkit), markdown: headingSmaller. */
  mdkitMarkdownHeadingSmaller: string
  /** Editors of a text (mdkit), markdown: headingBigger. */
  mdkitMarkdownHeadingBigger: string
  /** Editors of a text (mdkit), markdown: headingLevels. */
  mdkitMarkdownHeadingLevels: string
  /** Editors of a text (mdkit), markdown: cleanBlock. */
  mdkitMarkdownCleanBlock: string
  /** Editors of a text (mdkit), markdown: taskList. */
  mdkitMarkdownTaskList: string
  /** Editors of a text (mdkit), markdown: exitFullscreen. */
  mdkitMarkdownExitFullscreen: string
  /** Editors of a text (mdkit), markdown: tipPictures. */
  mdkitMarkdownTipPictures: string
  /** Editors of a text (mdkit), markdown: tipFootnote. */
  mdkitMarkdownTipFootnote: string
  /** Editors of a text (mdkit), markdown: tipLists. */
  mdkitMarkdownTipLists: string
  /** Editors of a text (mdkit), markdown: tipTables. */
  mdkitMarkdownTipTables: string
  /** Editors of a text (mdkit), markdown: sideBySide. */
  mdkitMarkdownSideBySide: string
  /** Editors of a text (mdkit), markdown: previewOnly. */
  mdkitMarkdownPreviewOnly: string
  /** Editors of a text (mdkit), markdown: table. */
  mdkitMarkdownTable: string
  /** Editors of a text (mdkit), markdown: rule. */
  mdkitMarkdownRule: string
  /** Editors of a text (mdkit), markdown: fullscreen. */
  mdkitMarkdownFullscreen: string
  /** Editors of a text (mdkit), markdown: helpUrl. */
  mdkitMarkdownHelpUrl: string
  /** Editors of a text (mdkit), card: button. */
  mdkitCardButton: string
  /** Editors of a text (mdkit), card: title. */
  mdkitCardTitle: string
  /** Editors of a text (mdkit), card: search. */
  mdkitCardSearch: string
  /** Editors of a text (mdkit), card: useAddress. */
  mdkitCardUseAddress: string
  /** Editors of a text (mdkit), card: none. */
  mdkitCardNone: string
  /** Editors of a text (mdkit), card: hint. */
  mdkitCardHint: string
  /** Editors of a text (mdkit), card: cancel. */
  mdkitCardCancel: string
  /** Editors of a text (mdkit), card: remove. */
  mdkitCardRemove: string
  /** Editors of a text (mdkit), card: pages. */
  mdkitCardPages: string
  /** Editors of a text (mdkit), card: posts. */
  mdkitCardPosts: string
}
