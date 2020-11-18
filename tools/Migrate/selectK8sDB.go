package main

import (
	"database/sql"
	"encoding/json"
	"github.com/elgris/sqrl"
	"strconv"
)

type SelectResult struct {
	Data  []map[string]interface{} `json:"data"`
	Count map[string]int64         `json:"count,omitempty"`
}


type Error struct {
	ErrorMessage string `json:"message,omitempty"`
	ErrorCode    int    `json:"code,omitempty"`
}

func (e Error) Error() string {
	return e.ErrorMessage
}

func (e Error) Code() int {
	return e.ErrorCode
}

type SqlColumn struct {
	ColumnName string
	ColumnType string
}

type RequestColumnSqlColumnMap map[string]SqlColumn

func replaceRequestColumn(requestColumnSqlColumnMap RequestColumnSqlColumnMap) ([]string, []string) {
	resultColumns := make([]string, 0, len(requestColumnSqlColumnMap))
	selectColumns := make([]string, 0, len(requestColumnSqlColumnMap))

	for k, v := range requestColumnSqlColumnMap {
		resultColumns = append(resultColumns, k)
		selectColumns = append(selectColumns, v.ColumnName)
	}
	return resultColumns, selectColumns
}

func execSelectSql(db sqrl.BaseRunner, from string, requestColumnSqlColumnMap RequestColumnSqlColumnMap) ([]map[string]interface{}, map[string]int64, error) {
	var err error
	var rows *sql.Rows
	var data []map[string]interface{}
	var count map[string]int64

	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()

	for {
		resultColumns, selectColumns := replaceRequestColumn(requestColumnSqlColumnMap)
		if len(resultColumns) == 0 {
			data = make([]map[string]interface{}, 0)
			break
		}

		selectBuilder := sqrl.Select(selectColumns...).From(from)
		if rows, err = selectBuilder.RunWith(db).Query(); err != nil {
			break
		}

		columns := make([]interface{}, len(selectColumns))
		columnPointers := make([]interface{}, len(selectColumns))

		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		data = make([]map[string]interface{}, 0, 30)
		for rows.Next() {
			m := make(map[string]interface{}, len(selectColumns))
			if err := rows.Scan(columnPointers...); err != nil {
				break
			}
			for i, columnName := range resultColumns {
				value := columns[i]
				dataType := requestColumnSqlColumnMap[columnName].ColumnType
				selectColumnName := requestColumnSqlColumnMap[columnName].ColumnName
				switch dataType {
				case "json":
					if value == nil {
						m[selectColumnName] = nil
					} else {
						m[selectColumnName] = json.RawMessage(value.([]byte))
					}
				case "string":
					if value == nil {
						m[selectColumnName] = ""
					} else {
						m[selectColumnName] = string(value.([]byte))
					}
				case "bool":
					if value == nil {
						m[selectColumnName] = false
					} else {
						switch value.(type) {
						case []uint8:
							v := value.([]byte)[0]
							m[selectColumnName] = v == '1'
						case int64:
							v := value.(int64)
							m[selectColumnName] = v != 0
						}
					}
				case "int":
					if value == nil {
						m[selectColumnName] = nil
					} else {
						switch value.(type) {
						case int, int64, int32:
							m[selectColumnName] = value
						case []uint8:
							vStr := string(value.([]byte))
							vInt, _ := strconv.Atoi(vStr)
							m[selectColumnName] = vInt
						}
					}

				default:
					m[selectColumnName] = string(value.([]byte))
				}
			}
			data = append(data, m)
		}

		break
	}

	if err != nil {
		return nil, nil, err
	}

	return data, count, nil
}
func selectServer(tafDb *sql.DB) (SelectResult, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"ServerId": SqlColumn{
			ColumnName: "f_server_id",
			ColumnType: "int",
		},
		"ServerApp": SqlColumn{
			ColumnName: "f_server_app",
			ColumnType: "string",
		},
		"ServerName": SqlColumn{
			ColumnName: "f_server_name",
			ColumnType: "string",
		},
		"ServerMark": SqlColumn{
			ColumnName: "f_server_mark",
			ColumnType: "string",
		},
		"ServerType": SqlColumn{
			ColumnName: "f_server_type",
			ColumnType: "string",
		},
		"DeployPerson": SqlColumn{
			ColumnName: "f_deploy_person",
			ColumnType: "string",
		},
		"DeployTime": SqlColumn{
			ColumnName: "f_deploy_time",
			ColumnType: "string",
		},
	}
	const from = "t_server"

	var err error
	selectResult := SelectResult{
		Data:  nil,
		Count: nil,
	}

	if selectResult.Data, selectResult.Count, err = execSelectSql(tafDb, from, requestColumnsSqlColumnsMap); err != nil {
		return SelectResult{}, Error{err.Error(), -1}
	}

	return selectResult, Error{"", 0}
}

func selectServerAdapter(tafDb *sql.DB) (SelectResult, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"AdapterId": SqlColumn{
			ColumnName: "f_adapter_id",
			ColumnType: "int",
		},
		"ServerId": SqlColumn{
			ColumnName: "f_server_id",
			ColumnType: "int",
		},
		"Name": SqlColumn{
			ColumnName: "f_name",
			ColumnType: "string",
		},
		"Threads": SqlColumn{
			ColumnName: "f_threads",
			ColumnType: "int",
		},
		"Connections": SqlColumn{
			ColumnName: "f_connections",
			ColumnType: "int",
		},
		"Port": SqlColumn{
			ColumnName: "f_port",
			ColumnType: "int",
		},
		"Capacity": SqlColumn{
			ColumnName: "f_capacity",
			ColumnType: "int",
		},

		"Timeout": SqlColumn{
			ColumnName: "f_timeout",
			ColumnType: "int",
		},

		"IsTaf": SqlColumn{
			ColumnName: "f_is_taf",
			ColumnType: "bool",
		},
		"IsTcp": SqlColumn{
			ColumnName: "f_is_tcp",
			ColumnType: "bool",
		},
	}

	const from = "t_server_adapter"
	var err error
	selectResult := SelectResult{
		Data:  nil,
		Count: nil,
	}

	if selectResult.Data, selectResult.Count, err = execSelectSql(tafDb, from, requestColumnsSqlColumnsMap); err != nil {
		return SelectResult{}, Error{err.Error(), -1}
	}

	return selectResult, Error{"", 0}
}

func selectServicePool(tafDb *sql.DB) (SelectResult, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"ServiceId": SqlColumn{
			ColumnName: "f_service_id",
			ColumnType: "int",
		},
		"ServiceVersion": SqlColumn{
			ColumnName: "f_service_version",
			ColumnType: "int",
		},
		"ServiceMark": SqlColumn{
			ColumnName: "f_service_mark",
			ColumnType: "string",
		},
		"ServiceImage": SqlColumn{
			ColumnName: "f_service_image",
			ColumnType: "string",
		},
		"ImageDetail": SqlColumn{
			ColumnName: "f_image_detail",
			ColumnType: "json",
		},
		"ServerId": SqlColumn{
			ColumnName: "f_server_id",
			ColumnType: "int",
		},
		"ServerApp": SqlColumn{
			ColumnName: "f_server_app",
			ColumnType: "string",
		},
		"ServerName": SqlColumn{
			ColumnName: "f_server_name",
			ColumnType: "string",
		},
		"CreateTime": SqlColumn{
			ColumnName: "f_create_time",
			ColumnType: "string",
		},
		"CreatePerson": SqlColumn{
			ColumnName: "f_create_person",
			ColumnType: "string",
		},
	}

	const from = "t_service_pool"

	var err error
	selectResult := SelectResult{
		Data:  nil,
		Count: nil,
	}
	if selectResult.Data, selectResult.Count, err = execSelectSql(tafDb, from, requestColumnsSqlColumnsMap); err != nil {
		return SelectResult{}, Error{err.Error(), -1}
	}

	return selectResult, Error{"", 0}
}

func selectServerOption(tafDb *sql.DB) (SelectResult, Error) {

	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"ServerId": SqlColumn{
			ColumnName: "f_server_id",
			ColumnType: "int",
		},
		"ServerTemplate": SqlColumn{
			ColumnName: "f_server_template",
			ColumnType: "string",
		},
		"ServerProfile": SqlColumn{
			ColumnName: "f_server_profile",
			ColumnType: "string",
		},
		"StartScript": SqlColumn{
			ColumnName: "f_start_script_path",
			ColumnType: "string",
		},
		"StopScript": SqlColumn{
			ColumnName: "f_stop_script_path",
			ColumnType: "string",
		},
		"MonitorScript": SqlColumn{
			ColumnName: "f_monitor_script_path",
			ColumnType: "string",
		},
		"AsyncThread": SqlColumn{
			ColumnName: "f_async_thread",
			ColumnType: "int",
		},
		"ServerImportant": SqlColumn{
			ColumnName: "f_important_type",
			ColumnType: "int",
		},
		"RemoteLogType": SqlColumn{
			ColumnName: "f_remote_log_type",
			ColumnType: "int",
		},
		"RemoteLogReserve": SqlColumn{
			ColumnName: "f_remote_log_reserve_time",
			ColumnType: "int",
		},
		"RemoteLogCompress": SqlColumn{
			ColumnName: "f_remote_log_compress_time",
			ColumnType: "int",
		},
	}
	const from = "t_server_option"
	var err error
	selectResult := SelectResult{
		Data:  nil,
		Count: nil,
	}

	if selectResult.Data, selectResult.Count, err = execSelectSql(tafDb, from, requestColumnsSqlColumnsMap); err != nil {
		return SelectResult{}, Error{err.Error(), -1}
	}

	return selectResult, Error{"", 0}
}

func selectK8S(tafDb *sql.DB) (SelectResult, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"ServerId": SqlColumn{
			ColumnName: "f_server_id",
			ColumnType: "int",
		},
		"ServerApp": SqlColumn{
			ColumnName: "f_server_app",
			ColumnType: "string",
		},
		"ServerName": SqlColumn{
			ColumnName: "f_server_name",
			ColumnType: "string",
		},
		"HostPort": SqlColumn{
			ColumnName: "f_host_port",
			ColumnType: "json",
		},
		"HostIpc": SqlColumn{
			ColumnName: "f_host_ipc",
			ColumnType: "bool",
		},
		"HostNetwork": SqlColumn{
			ColumnName: "f_host_network",
			ColumnType: "bool",
		},
		"Replicas": SqlColumn{
			ColumnName: "f_replicas",
			ColumnType: "int",
		},
		"NodeSelector": SqlColumn{
			ColumnName: "f_node_selector",
			ColumnType: "json",
		},
	}

	const from string = "t_server_k8s"

	var err error
	selectResult := SelectResult{
		Data:  nil,
		Count: nil,
	}
	if selectResult.Data, selectResult.Count, err = execSelectSql(tafDb, from, requestColumnsSqlColumnsMap); err != nil {
		return SelectResult{}, Error{err.Error(), -1}
	}
	return selectResult, Error{"", 0}
}

func selectTemplate(tafDb *sql.DB) (SelectResult, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"TemplateId": SqlColumn{
			ColumnName: "f_template_id",
			ColumnType: "int",
		},
		"TemplateName": SqlColumn{
			ColumnName: "f_template_name",
			ColumnType: "string",
		},
		"TemplateParent": SqlColumn{
			ColumnName: "f_template_parent",
			ColumnType: "string",
		},

		"TemplateContent": SqlColumn{
			ColumnName: "f_template_content",
			ColumnType: "string",
		},
		"CreatePerson": SqlColumn{
			ColumnName: "f_create_person",
			ColumnType: "string",
		},
		"CreateTime": SqlColumn{
			ColumnName: "f_create_time",
			ColumnType: "string",
		},

		"CreateMark": SqlColumn{
			ColumnName: "f_create_mark",
			ColumnType: "string",
		},

		"UpdatePerson": SqlColumn{
			ColumnName: "f_update_person",
			ColumnType: "string",
		},

		"UpdateMark": SqlColumn{
			ColumnName: "f_update_mark",
			ColumnType: "string",
		},

		"UpdateTime": SqlColumn{
			ColumnName: "f_update_time",
			ColumnType: "string",
		},
	}
	const from = "t_template"
	var err error
	selectResult := SelectResult{
		Data:  nil,
		Count: nil,
	}
	if selectResult.Data, selectResult.Count, err = execSelectSql(tafDb, from, requestColumnsSqlColumnsMap); err != nil {
		return SelectResult{}, Error{err.Error(), -1}
	}
	return selectResult, Error{"", 0}
}

func selectConfig(tafDb *sql.DB) (SelectResult, Error) {
	requestColumnsSqlColumnsMap := RequestColumnSqlColumnMap{
		"ConfigId": SqlColumn{
			ColumnName: "f_config_id",
			ColumnType: "int",
		},
		"ConfigName": SqlColumn{
			ColumnName: "f_config_name",
			ColumnType: "string",
		},
		"ConfigVersion": SqlColumn{
			ColumnName: "f_config_version",
			ColumnType: "int",
		},
		"ConfigContent": SqlColumn{
			ColumnName: "f_config_content",
			ColumnType: "string",
		},
		"CreateTime": SqlColumn{
			ColumnName: "f_create_time",
			ColumnType: "string",
		},
		"CreatePerson": SqlColumn{
			ColumnName: "f_create_person",
			ColumnType: "string",
		},
		"ConfigMark": SqlColumn{
			ColumnName: "f_config_mark",
			ColumnType: "string",
		},
		"AppServer": SqlColumn{
			ColumnName: "f_app_server",
			ColumnType: "string",
		},
		"PodSeq": SqlColumn{
			ColumnName: "f_pod_seq",
			ColumnType: "int",
		},
	}

	const from = "t_config"

	var err error

	selectResult := SelectResult{
		Data:  nil,
		Count: nil,
	}

	if selectResult.Data, selectResult.Count, err = execSelectSql(tafDb, from, requestColumnsSqlColumnsMap); err != nil {
		return SelectResult{}, Error{err.Error(), -1}
	}

	return selectResult, Error{"", 0}
}
