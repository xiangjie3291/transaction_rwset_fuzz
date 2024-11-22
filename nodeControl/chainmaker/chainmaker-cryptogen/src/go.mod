module chainmaker.org/chainmaker-cryptogen

go 1.15

require (
	chainmaker.org/chainmaker/common/v2 v2.3.0
	github.com/mr-tron/base58 v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.9.0
)

replace (
	github.com/spf13/afero => github.com/spf13/afero v1.5.1 //for go1.15 build
	github.com/spf13/viper => github.com/spf13/viper v1.7.1 //for go1.15 build
)
