package rpc

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/elgris/sqrl"
	"runtime"
)

const _RequestKindDefaultValue = "DefaultValue"

func selectDefaultValue(rpcRequest *Request) (Result, Error) {

	selectBuilder := sqrl.Select("f_label", "f_value").From("t_default_value")
	selectBuilder.Where(sqrl.Eq{"f_label": *rpcRequest.Params.Columns})

	var err error
	var rows *sql.Rows
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()

	if rows, err = selectBuilder.RunWith(tafDb).Query(); err != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

		return nil, Error{"内部错误", -1}
	}

	selectResult := make(map[string]json.RawMessage, len(*rpcRequest.Params.Columns))
	for rows.Next() {
		var label string
		var value []byte
		if err = rows.Scan(&label, &value); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))

			return nil, Error{"内部错误", -1}
		}
		selectResult[label] = value
	}
	return selectResult, Error{"", Success}
}

func init() {
	registryHandle(RequestMethodSelect+_RequestKindDefaultValue, selectDefaultValue)
}
