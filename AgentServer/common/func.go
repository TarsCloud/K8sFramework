package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorInfo struct {
	ErrCode int 	`json:"err_code"`
	ErrMsg	string 	`json:"err_msg"`
}

func (e *ErrorInfo) Error() string {
	return e.ErrMsg
}

func WriteJsonRsp(writer http.ResponseWriter, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		errInfo := ErrorInfo{ErrCode: -1, ErrMsg: err.Error()}
		_ = json.NewEncoder(writer).Encode(errInfo)
	}
}
//对产生的任何error进行处理
func JSONAppErrorReporter() gin.HandlerFunc {
	return jsonAppErrorReporterT(gin.ErrorTypeAny)
}

func jsonAppErrorReporterT(errType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		detectedErrors := c.Errors.ByType(errType)
		if len(detectedErrors) > 0 {
			err := detectedErrors[0].Err
			var parsedError *ErrorInfo
			switch err.(type) {
			//如果产生的error是自定义的结构体,转换error,返回自定义的code和msg
			case *ErrorInfo:
				parsedError = err.(*ErrorInfo)
			default:
				parsedError = &ErrorInfo{
					ErrCode:    http.StatusInternalServerError,
					ErrMsg: err.Error(),
				}
			}
			c.IndentedJSON(parsedError.ErrCode, parsedError)
			return
		}

	}
}

//设置所有跨域请求
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

