module github.com/talos-systems/sidero

go 1.16

replace (
	github.com/pensando/goipmi v0.0.0-20200303170213-e858ec1cf0b5 => github.com/talos-systems/goipmi v0.0.0-20210504182258-b54796c8d678
	sigs.k8s.io/cluster-api v0.3.12 => sigs.k8s.io/cluster-api v0.3.9
)

require (
	github.com/evanphx/json-patch v5.6.0+incompatible
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v1.2.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	github.com/pensando/goipmi v0.0.0-20200303170213-e858ec1cf0b5
	github.com/pin/tftp v2.1.1-0.20200117065540-2f79be2dba4e+incompatible
	github.com/pkg/errors v0.9.1
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.7.1
	github.com/talos-systems/cluster-api-bootstrap-provider-talos v0.5.4
	github.com/talos-systems/cluster-api-control-plane-provider-talos v0.4.7
	github.com/talos-systems/go-blockdevice v0.3.3
	github.com/talos-systems/go-debug v0.2.1
	github.com/talos-systems/go-kmsg v0.1.1
	github.com/talos-systems/go-procfs v0.0.0-20210108152626-8cbc42d3dc24
	github.com/talos-systems/go-retry v0.3.1
	github.com/talos-systems/go-smbios v0.0.0-20210422124317-d3a32bea731a
	github.com/talos-systems/net v0.3.2
	github.com/talos-systems/talos/pkg/machinery v1.1.1
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd
	golang.org/x/sync v0.0.0-20220513210516-0976fa681c29
	golang.org/x/sys v0.0.0-20220408201424-a24fb2fb8a0f
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.28.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	k8s.io/api v0.23.5
	k8s.io/apiextensions-apiserver v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v0.23.5
	k8s.io/utils v0.0.0-20211116205334-6203023598ed
	sigs.k8s.io/cluster-api v1.1.4
	sigs.k8s.io/controller-runtime v0.11.2
)
