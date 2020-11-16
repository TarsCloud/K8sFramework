module k8s/watch

go 1.14

require (
	k8s.io/api v0.16.9
	k8s.io/client-go v0.16.9
	base v0.0.0
)
replace (
		base v0.0.0 => ./../../base/
)