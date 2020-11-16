package rpc

import (
	"base"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
)

type UpdateRequestTarget map[string]*json.RawMessage

type SelectRequestColumns []string

type SelectRequestCount []string

type SelectRequestOrderElem struct {
	Column string `json:"column"`
	Order  string `json:"order"`
}

type SelectRequestOrder []SelectRequestOrderElem

type SelectRequestFilter struct {
	EQ   map[string]interface{} `json:"eq"`
	NE   map[string]interface{} `json:"ne"`
	GT   map[string]interface{} `json:"gt"`
	GE   map[string]interface{} `json:"ge"`
	LT   map[string]interface{} `json:"lt"`
	LE   map[string]interface{} `json:"le"`
	IN   map[string]interface{} `json:"in"`
	LIKE map[string]string      `json:"like"`
}

type DoRequestAction struct {
	Key   string           `json:"key"`
	Value *json.RawMessage `json:"value"`
}

type SelectRequestLimiter struct {
	Offset uint64 `json:"offset"`
	Rows   uint64 `json:"rows"`
}

type RequestParams struct {
	Kind         string                `json:"kind"`
	Token        string                `json:"token"`
	Metadata     *json.RawMessage      `json:"metadata"`
	Columns      *SelectRequestColumns `json:"columns"`
	Count        *SelectRequestCount   `json:"count"`
	Filter       *SelectRequestFilter  `json:"filter"`
	Order        *SelectRequestOrder   `json:"order"`
	Target       *UpdateRequestTarget  `json:"target"`
	Limiter      *SelectRequestLimiter `json:"limiter"`
	Action       *DoRequestAction      `json:"action"`
	Confirmation bool                  `json:"confirmation"`
}

type Request struct {
	Version        string         `json:"json"`
	ID             string         `json:"id"`
	Method         string         `json:"method"`
	Params         *RequestParams `json:"params"`
	RequestAccount *base.RequestAccount
}

type SelectResult struct {
	Data  []map[string]interface{} `json:"data"`
	Count map[string]int64         `json:"count,omitempty"`
}

type Result interface{}

type NormalResponse struct {
	Version string      `json:"json"`
	ID      string      `json:"id"`
	Result  interface{} `json:"result"`
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

type ErrorResponse struct {
	Version string `json:"json"`
	ID      string `json:"id"`
	Error   Error  `json:"error"`
}

const RequestMethodCreate string = "create"
const RequestMethodSelect string = "select"
const RequestMethodUpdate string = "update"
const RequestMethodDelete string = "delete"
const RequestMethodDo string = "do"

func verifyRequest(request *Request) (bool, string) {

	if request.Params == nil {
		return false, "Bad Params : Params Should Not Null"
	}

	if request.Params.Kind == "" {
		return false, "Bad Params : Params.Kind Should Not Empty"
	}

	switch request.Method {
	case RequestMethodCreate:
		if request.Params.Metadata == nil {
			return false, "Bad Schema : Params.Metadata Should Not Null"
		}
	case RequestMethodDelete:
		if request.Params.Metadata == nil {
			return false, "Bad Schema : Params.Metadata Should Not Null"
		}
	case RequestMethodSelect:

		if request.Params.Columns == nil && request.Params.Count == nil {
			return false, "Bad Schema : Params.Columns Or Params.Count Should Not Null"
		}

		if request.Params.Columns != nil && len(*request.Params.Columns) == 0 {
			return false, "Bad Schema : Params.Columns Should Not Empty"
		}

		if request.Params.Count != nil && len(*request.Params.Count) == 0 {
			return false, "Bad Schema : Params.Count Should Not Empty"
		}
	case RequestMethodUpdate:
		if request.Params.Metadata == nil {
			return false, "Bad Schema : Params.Metadata Should Not Null"
		}
		if request.Params.Target == nil {
			return false, "Bad Schema : Params.Target Should Not Null"
		}
	case RequestMethodDo:
		if request.Params.Action == nil {
			return false, "Bad Schema : Params.Action Should Not Null"
		}
		if request.Params.Action.Key == "" {
			return false, "Bad Schema : Params.Action.Key Should Not Null"
		}
		if request.Params.Action.Value == nil {
			return false, "Bad Schema : Params.Action.Value Should Not Null"
		}
	default:
		return false, "Bad Schema : Unsupported Method"
	}

	return true, ""
}

func GenerateRequest(in io.ReadCloser) (*Request, error) {
	request := &Request{
		Version:        "",
		ID:             "-1",
		Method:         "",
		Params:         nil,
		RequestAccount: nil,
	}

	if err := json.NewDecoder(in).Decode(request); err != nil {
		return request, errors.New("Bad Schema : " + err.Error())
	}

	if ok, message := verifyRequest(request); ok != true {
		return request, errors.New(message)
	}

	return request, nil
}

func HandleRequest(request *Request) (Result, Error) {

	var errMessage string
	var errCode int

	//todo Token 管理
	//var account *RequestAccount

	//if request.Params.Token == "" {
	//	if request.Method != RequestMethodCreate && request.Params.Kind != "Token" {
	//		return nil, "Bad Schema : Params.Token Should Not Empty", -1
	//	}
	//	return tokenHandler.CreateToken(request)
	//}
	//
	//if account, errMessage, errCode = tokenHandler.LoadAccountFormToken(request.Params.Token); errorCode != Success {
	//	return nil, errMessage, errCode
	//}
	//
	//request.RequestAccount = account

	request.RequestAccount = &base.RequestAccount{
		Name: "admin",
	}

	var fun func(*Request) (Result, Error)
	var ok bool
	if request.Method != RequestMethodDo {
		fun, ok = handlerFunctionTable[request.Method+request.Params.Kind]
		errMessage = "Bad Schema : Unsupported Params.Kind+Method"
		errCode = UnSupported
	} else {
		fun, ok = actionFunctionTable[request.Params.Kind+request.Params.Action.Key]
		errMessage = "Bad Schema : Unsupported Params.Kind+Params.Action.Key"
		errCode = UnSupported
	}
	if !ok {
		return nil, Error{errMessage, errCode}
	}

	return fun(request)
}

func GenerateENormalResponse(request *Request, result interface{}) []byte {
	response := NormalResponse{
		Version: request.Version,
		ID:      request.ID,
		Result:  result,
	}
	bytes, _ := json.Marshal(response)
	return bytes
}

func GenerateErrorResponse(request *Request, error Error) []byte {
	response := ErrorResponse{
		Version: request.Version,
		ID:      request.ID,
		Error: Error{
			ErrorCode:    error.Code(),
			ErrorMessage: error.Error(),
		},
	}
	bytes, _ := json.Marshal(response)
	return bytes
}

type handlerFunction func(*Request) (Result, Error)

var handlerFunctionTable map[string]handlerFunction
var actionFunctionTable map[string]handlerFunction

func registryHandle(methodWithKind string, function handlerFunction) {
	if handlerFunctionTable == nil {
		handlerFunctionTable = make(map[string]handlerFunction, 50)
	}
	if methodWithKind == "" {
		panic("Bad MethodWithKind Value :" + methodWithKind)
	}
	handlerFunctionTable[methodWithKind] = function
}
func registryAction(kindWithAction string, function handlerFunction) {
	if actionFunctionTable == nil {
		actionFunctionTable = make(map[string]handlerFunction, 30)
	}
	if kindWithAction == "" {
		panic("Bad KindWithAction Value :" + kindWithAction)
	}
	actionFunctionTable[kindWithAction] = function
}

type Imp struct {
}

func (rpc Imp) SetK8SWatchImp(watchImp base.K8SWatchInterface) {
	k8sWatchImp = watchImp
}

func (rpc Imp) SetTafDb(db *sql.DB) {
	tafDb = db
}

func (rpc Imp) SetK8SClientImp(imp base.K8SClientInterface) {
	k8sClientImp = imp
}

func (rpc Imp) Handler(io io.ReadCloser) []byte {
	rpcRequest, err := GenerateRequest(io)
	if err != nil {
		return GenerateErrorResponse(rpcRequest, Error{err.Error(), BadRequestSchema})
	}

	result, err_ := HandleRequest(rpcRequest)
	if err_.Code() != 0 {
		return GenerateErrorResponse(rpcRequest, err_)
	}
	return GenerateENormalResponse(rpcRequest, result)
}
