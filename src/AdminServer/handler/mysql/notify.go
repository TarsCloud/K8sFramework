package mysql

import (
	"github.com/go-openapi/runtime/middleware"
	"tafadmin/openapi/models"
	"tafadmin/openapi/restapi/operations/notify"
)

var notifyColumnsSqlColumnsMap = RequestColumnSqlColumnMap{
	"NotifyId": SqlColumn{
		ColumnName: "f_notify_id",
		ColumnType: "int",
	},
	"AppServer": SqlColumn{
		ColumnName: "f_app_server",
		ColumnType: "string",
	},
	"PodName": SqlColumn{
		ColumnName: "f_pod_name",
		ColumnType: "string",
	},
	"NotifyLevel": SqlColumn{
		ColumnName: "f_notify_level",
		ColumnType: "string",
	},
	"NotifyMessage": SqlColumn{
		ColumnName: "f_notify_message",
		ColumnType: "string",
	},
	"NotifyTime": SqlColumn{
		ColumnName: "f_notify_time",
		ColumnType: "string",
	},
	"NotifyThread": SqlColumn{
		ColumnName: "f_notify_thread",
		ColumnType: "string",
	},
	"NotifySource": SqlColumn{
		ColumnName: "f_notify_source",
		ColumnType: "string",
	},
}

type SelectNotifyHandler struct {}

func (s *SelectNotifyHandler) Handle(params notify.SelectNotifyParams) middleware.Responder {
	const from = "t_server_notify"

	result, err := SelectQueryResult(from, params.Filter, params.Limiter, params.Order, notifyColumnsSqlColumnsMap)
	if err != nil {
		return notify.NewSelectNotifyInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return notify.NewSelectNotifyOK().WithPayload(result)
}

