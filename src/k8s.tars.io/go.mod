module k8s.tars.io

go 1.14

require (
	github.com/google/go-cmp v0.5.1 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/tools v0.0.0-20200825202427-b303f430e36d // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/code-generator v0.18.8
)

replace k8s.io/client-go => k8s.io/client-go v0.18.2
