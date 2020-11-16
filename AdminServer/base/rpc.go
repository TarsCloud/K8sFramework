package base

import (
	"database/sql"
	"io"
)

type RPCInterface interface {
	Handler(in io.ReadCloser) []byte
	SetTafDb(db *sql.DB)
	SetK8SClientImp(k8sClientImp K8SClientInterface)
	SetK8SWatchImp(k8sWatchImp K8SWatchInterface)
}
