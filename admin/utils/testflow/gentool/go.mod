module github.com/go-rvq/rvq/admin/utils/testflow/gentool

go 1.24

toolchain go1.24.3

require (
	github.com/go-rvq/rvq v0.0.0-20250707114503-7cc586165fd6
	github.com/gobuffalo/flect v1.0.2
	github.com/pkg/errors v0.9.1
	github.com/sergi/go-diff v1.3.1
	mvdan.cc/gofumpt v0.6.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/mod v0.18.0 // indirect
	golang.org/x/tools v0.22.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/go-rvq/rvq => ../../../..
