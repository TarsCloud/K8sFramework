package rpc

import (
	"base"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var tafDb *sql.DB
var k8sClientImp base.K8SClientInterface
var k8sWatchImp base.K8SWatchInterface
