module tafadmin

require (
	base v0.0.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	k8s.io/api v0.16.9
	k8s.io/apimachinery v0.16.9
	k8s.io/client-go v0.16.9
	k8s/client v0.0.0
	k8s/watch v0.0.0
	rpc v0.0.0
)

replace (
	base v0.0.0 => ./base
	k8s/client v0.0.0 => ./k8s/client/
	k8s/watch v0.0.0 => ./k8s/watch/
	rpc v0.0.0 => ./rpc/
)

go 1.14
