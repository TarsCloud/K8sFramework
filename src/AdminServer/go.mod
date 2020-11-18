module tafadmin

require (
	github.com/go-sql-driver/mysql v1.5.0
	k8s.io/client-go v0.18.6
	tafadmin/handler v0.0.0
)

replace (
	k8s.taf.io/crd v0.0.1 => ./k8s.taf.io/crd/
	tafadmin/handler v0.0.0 => ./handler
	tafadmin/openapi v0.0.0 => ./openapi
)

go 1.14
