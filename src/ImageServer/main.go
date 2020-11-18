package main

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var gBuildTaskManger BuildTaskManager

func main() {

	gBuildTaskManger.init()

	// ConfigMap挂载目录
	_, err := os.Stat(DefaultConfigMapPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/upload", func(writer http.ResponseWriter, r *http.Request) {
		type UploadResponse struct {
			UploadStatus  string `json:"UploadStatus"`
			UploadMessage string `json:"UploadMessage"`
			FileName      string `json:"FileName"`
		}

		var uploadResponse UploadResponse
		var err error

		for true {
			if err = r.ParseMultipartForm(32 << 20); err != nil {
				uploadResponse.UploadStatus = "error"
				uploadResponse.UploadMessage = "错误的请求格式"
				break
			}

			file, handler, err := r.FormFile("uploadfile")
			if err != nil {
				uploadResponse.UploadStatus = "error"
				uploadResponse.UploadMessage = "错误的请求格式"
				break
			}
			fileExt := filepath.Ext(handler.Filename)
			fileName := strconv.FormatInt(time.Now().UnixNano(), 10) + fileExt
			absoluteServerFile := DefaultUploadDir + fileName

			f, _ := os.OpenFile(absoluteServerFile, os.O_CREATE|os.O_WRONLY, 0666)

			if _, err = io.Copy(f, file); err != nil {
				uploadResponse.UploadStatus = "error"
				uploadResponse.UploadMessage = "内部错误"
				break
			}

			go func() {
				time.Sleep(time.Minute * 45)
				_ = os.Remove(absoluteServerFile)
			}()

			uploadResponse.UploadStatus = "done"
			uploadResponse.FileName = fileName
			break
		}

		byteData, _ := json.Marshal(uploadResponse)
		writer.Header().Add("Content-Type", "application/json")
		_, _ = fmt.Fprint(writer, string(byteData))
	})

	http.HandleFunc("/build", func(writer http.ResponseWriter, r *http.Request) {
		var detail BuildRequest
		var buildStatus *BuildStatus
		var err error

		for true {
			if err = json.NewDecoder(r.Body).Decode(&detail); err != nil {
				buildStatus = &BuildStatus{
					BuildId:      "",
					BuildStatus:  "error",
					BuildMessage: "任务启动失败: " + err.Error(),
				}
				break
			}

			if _, err = govalidator.ValidateStruct(detail); err != nil {
				buildStatus = &BuildStatus{
					BuildId:      "",
					BuildStatus:  "error",
					BuildMessage: "任务启动失败: " + err.Error(),
				}
				break
			}

			buildStatus = gBuildTaskManger.CreateBuildTask(detail)
			break
		}

		byteData, _ := json.Marshal(buildStatus)
		writer.Header().Add("Content-Type", "application/json")
		_, _ = fmt.Fprint(writer, string(byteData))
	})

	http.HandleFunc("/buildStatus", func(writer http.ResponseWriter, r *http.Request) {

		type BuildStatusRequest struct {
			BuildId string `json:"BuildId"`
		}

		var request BuildStatusRequest
		var response *BuildStatus
		var err error

		for true {
			if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
				response = &BuildStatus{
					BuildId:      "",
					BuildStatus:  "error",
					BuildMessage: "错误的请求参数",
				}
				break
			}

			if _, err = govalidator.ValidateStruct(request); err != nil {
				response = &BuildStatus{
					BuildId:      "",
					BuildStatus:  "error",
					BuildMessage: "错误的请求参数",
				}
				break
			}

			response = gBuildTaskManger.GetStatus(request.BuildId)

			break
		}
		byteData, _ := json.Marshal(response)
		writer.Header().Add("Content-Type", "application/json")
		fmt.Print("will response : ", string(byteData))
		_, _ = fmt.Fprint(writer, string(byteData))
	})

	err = http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println(err)
	}
}
