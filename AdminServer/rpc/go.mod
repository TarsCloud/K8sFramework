module rpc

require (
	github.com/TarsCloud/TarsGo v1.1.3
	github.com/asaskevich/govalidator v0.0.0-20200428143746-21a406dcc535
	github.com/elgris/sqrl v0.0.0-20190909141434-5a439265eeec
	github.com/go-sql-driver/mysql v1.5.0
	base v0.0.0
)

replace  base v0.0.0 => ./../base/
go 1.14