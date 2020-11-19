package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func StartHTTPServer()  {
	CPUFunc := func(writer http.ResponseWriter, r *http.Request) {
		data, err := GetCpuInfo()
		writeJsonRsp(writer, data, err)
	}

	MemoryFunc := func(writer http.ResponseWriter, r *http.Request) {
		data, err := GetMemInfo()
		writeJsonRsp(writer, data, err)
	}

	DiskFunc := func(writer http.ResponseWriter, r *http.Request) {
		data, err := GetDiskInfo()
		writeJsonRsp(writer, data, err)
	}

	HostFunc := func(writer http.ResponseWriter, r *http.Request) {
		data, err := GetHostInfo()
		writeJsonRsp(writer, data, err)
	}

	NetFunc := func(writer http.ResponseWriter, r *http.Request) {
		data, err := GetNetInfo()
		writeJsonRsp(writer, data, err)
	}

	PortFunc := func(writer http.ResponseWriter, r *http.Request) {
		writer.Header().Add("Content-Type", "application/json")

		host := r.URL.Query().Get("host")
		if host == "" {
			host = "127.0.0.1"
		}
		port, _ := strconv.Atoi(r.URL.Query().Get("port"))

		data, err := GetPortInfo(host, port)
		writeJsonRsp(writer, data, err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/cpu",  	CPUFunc)
	mux.HandleFunc("/memory", MemoryFunc)
	mux.HandleFunc("/disk", 	DiskFunc)
	mux.HandleFunc("/host", 	HostFunc)
	mux.HandleFunc("/net", 	NetFunc)
	mux.HandleFunc("/port", 	PortFunc)

	port := os.Getenv("TAF_AGENT_PORT")
	if port == "" {
		port = "8000"
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", "0.0.0.0", port),
		Handler:           mux,
		ReadTimeout:       2 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      2 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(fmt.Sprintf("tars-agent exited with %s.", err))
	}
}

func writeJsonRsp(writer http.ResponseWriter, data interface{}, err error) {
	writer.Header().Add("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(writer).Encode(err)
	} else {
		writer.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(writer).Encode(data)
	}
}

