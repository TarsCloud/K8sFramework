module tarsadmin

require (
	github.com/go-sql-driver/mysql v1.5.0
	k8s.io/client-go v11.0.0+incompatible
	tarsadmin/handler v0.0.0
)

replace (
	k8s.io/client-go => k8s.io/client-go v0.18.2
	k8s.tars.io v0.0.1 => ../k8s.tars.io
	tarsadmin/handler v0.0.0 => ./handler
	tarsadmin/openapi v0.0.0 => ./openapi
)

go 1.14
