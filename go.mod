module github.com/hashicorp/terraform-ls

go 1.13

require (
	github.com/apparentlymart/go-textseg v1.0.0
	github.com/creachadair/jrpc2 v0.10.1
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gammazero/workerpool v1.0.0
	github.com/google/go-cmp v0.5.1
	github.com/hashicorp/go-multierror v1.1.0
	github.com/hashicorp/go-version v1.2.1
	github.com/hashicorp/hcl-lang v0.0.0-20201027200521-9c11e0151346
	github.com/hashicorp/hcl/v2 v2.6.0
	github.com/hashicorp/terraform-exec v0.11.1-0.20201007122305-ea2094d52cb5
	github.com/hashicorp/terraform-json v0.6.0
	github.com/hashicorp/terraform-schema v0.0.0-20201027194524-5093f7354c6b
	github.com/mh-cbon/go-fmt-fail v0.0.0-20160815164508-67765b3fbcb5
	github.com/mitchellh/cli v1.1.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.3.2
	github.com/pmezard/go-difflib v1.0.0
	github.com/sourcegraph/go-lsp v0.0.0-20200117082640-b19bb38222e2
	github.com/spf13/afero v1.3.2
	github.com/stretchr/testify v1.4.0
	github.com/vektra/mockery/v2 v2.3.0
	github.com/zclconf/go-cty v1.6.1
)

replace github.com/sourcegraph/go-lsp => github.com/radeksimko/go-lsp v0.1.0
