module github.com/msmsny/gopt

go 1.16

require (
	github.com/spf13/cobra v1.1.3
	github.com/stretchr/testify v1.3.0
)

// [CVE-2020-15114] In etcd before versions 3.3.23 and 3.4.10, the etcd gateway is a simple TCP prox...
replace github.com/coreos/etcd => github.com/coreos/etcd v3.3.24+incompatible
