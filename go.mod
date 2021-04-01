module github.com/msmsny/gopt

go 1.16

require (
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.3.0
)

replace (
	// [CVE-2020-15114] In etcd before versions 3.3.23 and 3.4.10, the etcd gateway is a simple TCP prox...
	github.com/coreos/etcd => github.com/coreos/etcd v3.3.24+incompatible
	// [CVE-2021-3121] An issue was discovered in GoGo Protobuf before 1.3.2. plugin/unmarshal/unmarsha...
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
)
