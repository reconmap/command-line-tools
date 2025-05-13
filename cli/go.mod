module github.com/reconmap/cli

go 1.24.1

require (
	github.com/fatih/color v1.18.0
	github.com/rodaine/table v1.3.0
	github.com/urfave/cli/v2 v2.27.6
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
)

require (
	github.com/coreos/go-oidc v2.3.0+incompatible
	github.com/reconmap/shared-lib v0.0.0-20220910165932-7d018d9111fc
	golang.org/x/oauth2 v0.30.0
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pquerna/cachecontrol v0.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	gopkg.in/go-jose/go-jose.v2 v2.6.3 // indirect
)

replace github.com/reconmap/shared-lib => ../shared-lib
