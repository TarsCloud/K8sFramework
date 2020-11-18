module tafcontrol

go 1.14

require (
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	k8s.io/klog v1.0.0
	k8s.taf.io/crd v0.0.1
)

replace k8s.taf.io/crd v0.0.1 => ./k8s.taf.io/crd/
