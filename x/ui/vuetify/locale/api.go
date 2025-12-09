package locale

type Messages struct {
	Badge string `json:"badge"`
	Open string `json:"open"`
	Close string `json:"close"`
	Dismiss string `json:"dismiss"`
	ConfirmEdit ConfirmEdit `json:"confirmEdit"`
	DataIterator DataIterator `json:"dataIterator"`
	DataTable DataTable `json:"dataTable"`
	DataFooter DataFooter `json:"dataFooter"`
	DateRangeInput DateRangeInput `json:"dateRangeInput"`
	DatePicker DatePicker `json:"datePicker"`
	NoDataText string `json:"noDataText"`
	Carousel Carousel `json:"carousel"`
	Calendar Calendar `json:"calendar"`
	Input Input `json:"input"`
	FileInput FileInput `json:"fileInput"`
	FileUpload FileUpload `json:"fileUpload"`
	TimePicker TimePicker `json:"timePicker"`
	Pagination Pagination `json:"pagination"`
	Stepper Stepper `json:"stepper"`
	Rating Rating `json:"rating"`
	Loading string `json:"loading"`
	InfiniteScroll InfiniteScroll `json:"infiniteScroll"`
	Rules Rules `json:"rules"`
}
type ConfirmEdit struct {
	Ok string `json:"ok"`
	Cancel string `json:"cancel"`
}
type DataIterator struct {
	NoResultsText string `json:"noResultsText"`
	LoadingText string `json:"loadingText"`
}
type AriaLabel struct {
	SortDescending string `json:"sortDescending"`
	SortAscending string `json:"sortAscending"`
	SortNone string `json:"sortNone"`
	ActivateNone string `json:"activateNone"`
	ActivateDescending string `json:"activateDescending"`
	ActivateAscending string `json:"activateAscending"`
}
type DataTable struct {
	ItemsPerPageText string `json:"itemsPerPageText"`
	AriaLabel AriaLabel `json:"ariaLabel"`
	SortBy string `json:"sortBy"`
}
type DataFooter struct {
	ItemsPerPageText string `json:"itemsPerPageText"`
	ItemsPerPageAll string `json:"itemsPerPageAll"`
	NextPage string `json:"nextPage"`
	PrevPage string `json:"prevPage"`
	FirstPage string `json:"firstPage"`
	LastPage string `json:"lastPage"`
	PageText string `json:"pageText"`
}
type DateRangeInput struct {
	Divider string `json:"divider"`
}
type Range struct {
	Title string `json:"title"`
	Header string `json:"header"`
}
type Input struct {
	Placeholder string `json:"placeholder"`
}
type DatePicker struct {
	ItemsSelected string `json:"itemsSelected"`
	Range Range `json:"range"`
	Title string `json:"title"`
	Header string `json:"header"`
	Input Input `json:"input"`
}
type AriaLabel struct {
	Delimiter string `json:"delimiter"`
}
type Carousel struct {
	Prev string `json:"prev"`
	Next string `json:"next"`
	AriaLabel AriaLabel `json:"ariaLabel"`
}
type Calendar struct {
	MoreEvents string `json:"moreEvents"`
	Today string `json:"today"`
}
type Input struct {
	Clear string `json:"clear"`
	PrependAction string `json:"prependAction"`
	AppendAction string `json:"appendAction"`
	Otp string `json:"otp"`
}
type FileInput struct {
	Counter string `json:"counter"`
	CounterSize string `json:"counterSize"`
}
type FileUpload struct {
	Title string `json:"title"`
	Divider string `json:"divider"`
	Browse string `json:"browse"`
}
type TimePicker struct {
	Am string `json:"am"`
	Pm string `json:"pm"`
	Title string `json:"title"`
}
type AriaLabel struct {
	Root string `json:"root"`
	Next string `json:"next"`
	Previous string `json:"previous"`
	Page string `json:"page"`
	CurrentPage string `json:"currentPage"`
	First string `json:"first"`
	Last string `json:"last"`
}
type Pagination struct {
	AriaLabel AriaLabel `json:"ariaLabel"`
}
type Stepper struct {
	Next string `json:"next"`
	Prev string `json:"prev"`
}
type AriaLabel struct {
	Item string `json:"item"`
}
type Rating struct {
	AriaLabel AriaLabel `json:"ariaLabel"`
}
type InfiniteScroll struct {
	LoadMore string `json:"loadMore"`
	Empty string `json:"empty"`
}
type Rules struct {
	Required string `json:"required"`
	Email string `json:"email"`
	Number string `json:"number"`
	Integer string `json:"integer"`
	Capital string `json:"capital"`
	MaxLength string `json:"maxLength"`
	MinLength string `json:"minLength"`
	StrictLength string `json:"strictLength"`
	Exclude string `json:"exclude"`
	NotEmpty string `json:"notEmpty"`
	Pattern string `json:"pattern"`
}
