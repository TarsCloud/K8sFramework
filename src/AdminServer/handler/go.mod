module tarsadmin/handler

go 1.14

require (
	github.com/Azure/go-autorest/autorest v0.9.0 // indirect
	github.com/TarsCloud/TarsGo v1.1.5
	github.com/elgris/sqrl v0.0.0-20190909141434-5a439265eeec
	github.com/evanphx/json-patch v4.2.0+incompatible // indirect
	github.com/go-openapi/errors v0.19.7
	github.com/go-openapi/loads v0.19.5
	github.com/go-openapi/runtime v0.19.22
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gophercloud/gophercloud v0.1.0 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/gregjones/httpcache v0.0.0-20180305231024-9cad4c3443a7 // indirect
	github.com/imdario/mergo v0.3.5 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	golang.org/x/net v0.0.0-20201006153459-a7d1128ccaa0
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0
	k8s.tars.io/crd v0.0.1
	tarsadmin/openapi v0.0.0
)

replace (
	k8s.tars.io/crd v0.0.1 => ../../k8s.tars.io/crd
	tarsadmin/openapi v0.0.0 => ../openapi
)
