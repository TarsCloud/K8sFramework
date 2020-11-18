package main

import (
	"fmt"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"net/http"
	"time"
)

type Webhook struct {
	validating Validating
	mutating   Mutating
	k8sOption  *K8SOption
	watcher    *Watcher
	srv        *http.Server
}

func NewWebhook(k8sOption *K8SOption, watcher *Watcher) *Webhook {
	webhook := &Webhook{
		k8sOption: k8sOption,
		watcher:   watcher,
		validating: Validating{
			k8sOption:               k8sOption,
			watcher:                 watcher,
			controlAccount:          fmt.Sprintf("system:serviceaccount:%s:%s", k8sOption.namespace, TafControlServiceAccount),
			garbageCollectorAccount: "system:serviceaccount:kube-system:generic-garbage-collector",
		},
	}
	return webhook
}

func (h *Webhook) Start(stopCh chan struct{}) {
	go wait.Until(func() {
		validatingFunc := func(writer http.ResponseWriter, r *http.Request) {
			writer.Header().Add("Content-Type", "application/json")
			writer.Header().Add("Connection", "keep-alive")
			h.validating.handle(r, writer)
		}

		mutatingFunc := func(writer http.ResponseWriter, r *http.Request) {
			writer.Header().Add("Content-Type", "application/json")
			writer.Header().Add("Connection", "keep-alive")
			h.mutating.handle(r, writer)
		}

		mux := http.NewServeMux()

		mux.HandleFunc("/validating", validatingFunc)
		mux.HandleFunc("/mutating", mutatingFunc)

		srv := &http.Server{
			Addr:              ":443",
			Handler:           mux,
			ReadTimeout:       2 * time.Second,
			ReadHeaderTimeout: 1 * time.Second,
			WriteTimeout:      2 * time.Second,
			IdleTimeout:       120 * time.Second,
		}
		// ListenAndServe always returns a non-nil error. After Shutdown or Close,
		// the returned error is ErrServerClosed.
		err := srv.ListenAndServeTLS(WebhookCertFile, WebhookCertKey)
		if err != nil {
			utilRuntime.HandleError(fmt.Errorf("will exist because : %s \n", err.Error()))
			notifyStop()
		}
	}, time.Second, stopCh)
}
