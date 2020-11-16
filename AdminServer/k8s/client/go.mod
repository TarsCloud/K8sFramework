module k8s/client

require (
	k8s.io/api v0.16.9
	k8s.io/apimachinery v0.16.9
	k8s.io/client-go v0.16.9
    base v0.0.0
)

replace (
  base v0.0.0 => ./../../base/
)
go 1.14
