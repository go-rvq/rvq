package web

// StdContentClass default class for standard content like Post Body, Page Content etc.
const StdContentClass = "std-content"

func StdContentStyles() ComponentsPack {
	return ComponentsPack(`
.std-content {
	p {
		margin: 8px 0;
	}
	li {
		margin-left: 10px;
	}
}
`)
}
