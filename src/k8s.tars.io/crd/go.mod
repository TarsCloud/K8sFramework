module k8s.tars.io/crd

go 1.14

require (
	golang.org/x/oauth2 v0.0.0-20201109201403-9fd604954f58 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	k8s.io/api v0.18.8
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/code-generator v0.18.8
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
)

replace (
	k8s.io/client-go => k8s.io/client-go v0.18.2
)
