package rpc

import "encoding/json"

type extractFunction func(key string, value *json.RawMessage) (columnName string, columnValue interface{}, err Error)

var extractFunctionTable map[string]extractFunction

func registryExtract(k string, function extractFunction) {
	if extractFunctionTable == nil {
		extractFunctionTable = make(map[string]extractFunction, 100)
	}
	extractFunctionTable[k] = function
}
