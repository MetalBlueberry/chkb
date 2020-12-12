module github.com/MetalBlueberry/chkb

go 1.14

require (
	github.com/benbjohnson/clock v1.1.0
	github.com/bendahl/uinput v1.4.0
	github.com/gvalkov/golang-evdev v0.0.0-20191114124502-287e62b94bcb
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/afero v1.4.1
	github.com/stretchr/testify v1.6.1
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

// Patch until merge is accepted, fix dead lock
replace github.com/benbjohnson/clock => github.com/MetalBlueberry/clock v1.1.1-0.20201212194419-31ee72d00441
