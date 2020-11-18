module tafadmin/handler

go 1.14

require (
	github.com/TarsCloud/TarsGo v1.1.5
	github.com/elgris/sqrl v0.0.0-20190909141434-5a439265eeec
	github.com/go-openapi/errors v0.19.7
	github.com/go-openapi/loads v0.19.5
	github.com/go-openapi/runtime v0.19.22
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gorilla/websocket v1.4.2
	golang.org/x/net v0.0.0-20201006153459-a7d1128ccaa0
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	k8s.io/klog v1.0.0
	k8s.taf.io/crd v0.0.1
	tafadmin/openapi v0.0.0
)

replace (
	k8s.taf.io/crd v0.0.1 => ../k8s.taf.io/crd/
	tafadmin/openapi v0.0.0 => ../openapi
)
