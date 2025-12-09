package locale

type Message struct {
	Badge       string `json:"badge"`
	Open        string `json:"open"`
	Close       string `json:"close"`
	Dismiss     string `json:"dismiss"`
	ConfirmEdit struct {
		Ok     string `json:"ok"`
		Cancel string `json:"cancel"`
	} `json:"confirmEdit"`
	DataIterator struct {
		NoResultsText string `json:"noResultsText"`
		LoadingText   string `json:"loadingText"`
	} `json:"dataIterator"`
	DataTable struct {
		ItemsPerPageText string `json:"itemsPerPageText"`
		AriaLabel        struct {
			SortDescending     string `json:"sortDescending"`
			SortAscending      string `json:"sortAscending"`
			SortNone           string `json:"sortNone"`
			ActivateNone       string `json:"activateNone"`
			ActivateDescending string `json:"activateDescending"`
			ActivateAscending  string `json:"activateAscending"`
		} `json:"ariaLabel"`
		SortBy string `json:"sortBy"`
	} `json:"dataTable"`
	DataFooter struct {
		ItemsPerPageText string `json:"itemsPerPageText"`
		ItemsPerPageAll  string `json:"itemsPerPageAll"`
		NextPage         string `json:"nextPage"`
		PrevPage         string `json:"prevPage"`
		FirstPage        string `json:"firstPage"`
		LastPage         string `json:"lastPage"`
		PageText         string `json:"pageText"`
	} `json:"dataFooter"`
	DateRangeInput struct {
		Divider string `json:"divider"`
	} `json:"dateRangeInput"`
	DatePicker struct {
		ItemsSelected string `json:"itemsSelected"`
		Range         struct {
			Title  string `json:"title"`
			Header string `json:"header"`
		} `json:"range"`
		Title  string `json:"title"`
		Header string `json:"header"`
		Input  struct {
			Placeholder string `json:"placeholder"`
		} `json:"input"`
	} `json:"datePicker"`
	NoDataText string `json:"noDataText"`
	Carousel   struct {
		Prev      string `json:"prev"`
		Next      string `json:"next"`
		AriaLabel struct {
			Delimiter string `json:"delimiter"`
		} `json:"ariaLabel"`
	} `json:"carousel"`
	Calendar struct {
		MoreEvents string `json:"moreEvents"`
		Today      string `json:"today"`
	} `json:"calendar"`
	Input struct {
		Clear         string `json:"clear"`
		PrependAction string `json:"prependAction"`
		AppendAction  string `json:"appendAction"`
		Otp           string `json:"otp"`
	} `json:"input"`
	FileInput struct {
		Counter     string `json:"counter"`
		CounterSize string `json:"counterSize"`
	} `json:"fileInput"`
	FileUpload struct {
		Title   string `json:"title"`
		Divider string `json:"divider"`
		Browse  string `json:"browse"`
	} `json:"fileUpload"`
	TimePicker struct {
		Am    string `json:"am"`
		Pm    string `json:"pm"`
		Title string `json:"title"`
	} `json:"timePicker"`
	Pagination struct {
		AriaLabel struct {
			Root        string `json:"root"`
			Next        string `json:"next"`
			Previous    string `json:"previous"`
			Page        string `json:"page"`
			CurrentPage string `json:"currentPage"`
			First       string `json:"first"`
			Last        string `json:"last"`
		} `json:"ariaLabel"`
	} `json:"pagination"`
	Stepper struct {
		Next string `json:"next"`
		Prev string `json:"prev"`
	} `json:"stepper"`
	Rating struct {
		AriaLabel struct {
			Item string `json:"item"`
		} `json:"ariaLabel"`
	} `json:"rating"`
	Loading        string `json:"loading"`
	InfiniteScroll struct {
		LoadMore string `json:"loadMore"`
		Empty    string `json:"empty"`
	} `json:"infiniteScroll"`
	Rules struct {
		Required     string `json:"required"`
		Email        string `json:"email"`
		Number       string `json:"number"`
		Integer      string `json:"integer"`
		Capital      string `json:"capital"`
		MaxLength    string `json:"maxLength"`
		MinLength    string `json:"minLength"`
		StrictLength string `json:"strictLength"`
		Exclude      string `json:"exclude"`
		NotEmpty     string `json:"notEmpty"`
		Pattern      string `json:"pattern"`
	} `json:"rules"`
}
