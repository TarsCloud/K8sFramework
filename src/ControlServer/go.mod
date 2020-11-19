module tarscontrol

go 1.14

require (
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0
	k8s.tars.io/crd v0.0.1
)

replace (
	k8s.io/client-go => k8s.io/client-go v0.18.6
	k8s.tars.io/crd v0.0.1 => ../k8s.tars.io/crd/
)
