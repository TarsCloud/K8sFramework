package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	k8sAdmissionV1 "k8s.io/api/admission/v1"
	crdV1Alpha1 "k8s.taf.io/crd/v1alpha1"
	"net/http"
)

type Mutating struct {
	k8sOption *K8SOption
	watcher   *Watcher
}

func mutatingCreateTDeploy(requestAdmissionView *k8sAdmissionV1.AdmissionReview) []byte {
	tdeploy := &crdV1Alpha1.TDeploy{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, tdeploy)

	var patchContents = make([][]byte, 0, 10)
	patchContents = append(patchContents, []byte{'['})

	patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/taf.io~1Approve\",\"value\":\"Pending\"}")))
	patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"remove\",\"path\":\"/approve\"}")))
	patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"remove\",\"path\":\"/deployed\"}")))

	totalPatchContent := bytes.Join(patchContents, []byte{','})
	totalPatchContent[1] = ' '
	totalPatchContent = append(totalPatchContent, ']')
	return totalPatchContent
}

func mutatingUpdateTDeploy(requestAdmissionView *k8sAdmissionV1.AdmissionReview) []byte {
	tdeploy := &crdV1Alpha1.TDeploy{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, tdeploy)

	var patchContents = make([][]byte, 0, 10)

	patchContents = append(patchContents, []byte{'{'})

	if tdeploy.Approve != nil {
		if !tdeploy.Approve.Result {
			patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/taf.io~1Approve\",\"value\":\"Reject\"}")))
		} else {
			patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/taf.io~1Approve\",\"value\":\"Approved\"}")))
		}
	}

	if len(patchContents) == 1 {
		return nil
	}

	totalPatchContent := bytes.Join(patchContents, []byte{','})
	totalPatchContent[1] = ' '
	totalPatchContent = append(totalPatchContent, ']')
	return totalPatchContent
}

func mutatingCreateTServer(requestAdmissionView *k8sAdmissionV1.AdmissionReview) []byte {
	server := &crdV1Alpha1.TServer{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, server)

	var patchContents = make([][]byte, 0, 10)
	patchContents = append(patchContents, []byte{'['})

	patchContents = append(patchContents, []byte(fmt.Sprintf(",{\"op\": \"add\", \"path\": \"/metadata/labels/taf.io~1ServerApp\", \"value\": \"%s\"}", server.Spec.App)))

	patchContents = append(patchContents, []byte(fmt.Sprintf(",{\"op\": \"add\", \"path\": \"/metadata/labels/taf.io~1ServerName\", \"value\": \"%s\"}", server.Spec.Server)))

	patchContents = append(patchContents, []byte(fmt.Sprintf(",{\"op\": \"add\", \"path\": \"/metadata/labels/taf.io~1SubType\", \"value\": \"%s\"}", server.Spec.SubType)))

	if server.Spec.Taf != nil {
		patchContents = append(patchContents, []byte(fmt.Sprintf(",{\"op\": \"add\", \"path\": \"/metadata/labels/taf.io~1Template\", \"value\": \"%s\"}", server.Spec.Taf.Template)))
		if !server.Spec.Taf.Foreground && server.Spec.K8S.ReadinessGate == nil {
			patchContents = append(patchContents, []byte(fmt.Sprintf(",{\"op\": \"add\", \"path\":\"/spec/k8s/readinessGate\",\"value\":\"%s\"}", TPodReadinessGate)))
		}
	}

	if server.Spec.Release == nil {
		patchContents = append(patchContents, []byte(fmt.Sprintf(",{\"op\": \"add\", \"path\":\"/spec/k8s/replicas\",\"value\": %d}", 0)))
	} else {
		patchContents = append(patchContents, []byte(fmt.Sprintf(",{\"op\": \"add\", \"path\": \"/metadata/labels/taf.io~1ReleaseSource\", \"value\": \"%s\"}", server.Spec.Release.Source)))
		patchContents = append(patchContents, []byte(fmt.Sprintf(",{\"op\": \"add\", \"path\": \"/metadata/labels/taf.io~1ReleaseTag\", \"value\": \"%s\"}", server.Spec.Release.Tag)))
	}

	if len(server.Spec.K8S.HostPorts) > 0 || server.Spec.K8S.HostIPC {
		patchContents = append(patchContents, []byte(fmt.Sprintf(",{\"op\": \"add\", \"path\":\"/spec/k8s/notStacked\",\"value\":%t}", true)))
	}

	totalPatchContent := bytes.Join(patchContents, []byte{','})
	totalPatchContent[1] = ' '
	totalPatchContent = append(totalPatchContent, ']')
	return totalPatchContent
}

func mutatingUpdateTServer(requestAdmissionView *k8sAdmissionV1.AdmissionReview) []byte {
	return mutatingCreateTServer(requestAdmissionView)
}

func mutatingCreateTConfig(requestAdmissionView *k8sAdmissionV1.AdmissionReview) []byte {
	tconfig := &crdV1Alpha1.TConfig{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, tconfig)

	var patchContents = make([][]byte, 0, 4)
	patchContents = append(patchContents, []byte{'['})

	if tconfig.ServerConfig != nil {
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/taf.io~1ServerApp\",\"value\":%s}", tconfig.ServerConfig.App)))
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/taf.io~1ServerName\",\"value\":%s}", tconfig.ServerConfig.Server)))
	}

	if tconfig.AppConfig != nil {
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/taf.io~1ServerApp\",\"value\":%s}", tconfig.ServerConfig.App)))
	}

	if len(patchContents) == 1 {
		return nil
	}

	return bytes.Join(patchContents, []byte{','})
}

func mutatingUpdateTConfig(requestAdmissionView *k8sAdmissionV1.AdmissionReview) []byte {
	return mutatingCreateTConfig(requestAdmissionView)
}

func mutatingCreateTTree(requestAdmissionView *k8sAdmissionV1.AdmissionReview) []byte {
	newTTree := &crdV1Alpha1.TTree{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, newTTree)

	businessMap := make(map[string]interface{}, len(newTTree.Businesses))
	for _, business := range newTTree.Businesses {
		businessMap[business.Name] = nil
	}

	var patchContents = make([][]byte, 0, 5)
	patchContents = append(patchContents, []byte{'['})

	for i, app := range newTTree.Apps {
		if app.BusinessRef != "" {
			if _, ok := businessMap[app.BusinessRef]; !ok {
				newTTreeApps := &crdV1Alpha1.TTreeApp{
					Name:         app.Name,
					BusinessRef:  "",
					CreatePerson: app.CreatePerson,
					CreateTime:   app.CreateTime,
					Mark:         app.Mark,
				}
				bs, _ := json.Marshal(newTTreeApps)
				patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"replace\",\"path\":\"/apps/%d\",\"value\":%s}", i, bs)))
			}
		}
	}

	totalPatchContent := bytes.Join(patchContents, []byte{','})
	totalPatchContent[1] = ' '
	totalPatchContent = append(totalPatchContent, ']')
	return totalPatchContent
}

func mutatingUpdateTTree(requestAdmissionView *k8sAdmissionV1.AdmissionReview) []byte {
	return mutatingCreateTTree(requestAdmissionView)
}

var mutatingFunctions map[string]func(*k8sAdmissionV1.AdmissionReview) []byte

func (v Mutating) handle(r *http.Request, w http.ResponseWriter) {
	requestAdmissionView := &k8sAdmissionV1.AdmissionReview{}
	err := json.NewDecoder(r.Body).Decode(&requestAdmissionView)
	if err != nil {
		return
	}
	responseAdmissionView := k8sAdmissionV1.AdmissionReview{
		Response: &k8sAdmissionV1.AdmissionResponse{
			UID: requestAdmissionView.Request.UID,
		},
	}
	var patchContent []byte

	key := fmt.Sprintf("%s/%s", string(requestAdmissionView.Request.Operation), requestAdmissionView.Request.Kind.Kind)
	if fun, ok := mutatingFunctions[key]; ok {
		patchContent = fun(requestAdmissionView)
	}

	if patchContent != nil {
		responseAdmissionView.Response.Patch = patchContent
		patchType := k8sAdmissionV1.PatchTypeJSONPatch
		responseAdmissionView.Response.PatchType = &patchType
	}

	responseAdmissionView.Response.Allowed = true
	responseBytes, _ := json.Marshal(responseAdmissionView)
	_, _ = w.Write(responseBytes)
}

func init() {
	mutatingFunctions = map[string]func(*k8sAdmissionV1.AdmissionReview) []byte{
		"CREATE/TDeploy": mutatingCreateTDeploy,
		"UPDATE/TDeploy": mutatingUpdateTDeploy,

		"CREATE/TServer": mutatingCreateTServer,
		"UPDATE/TServer": mutatingUpdateTServer,

		"CREATE/TConfig": mutatingCreateTConfig,
		"UPDATE/TConfig": mutatingUpdateTConfig,

		"CREATE/TTree": mutatingCreateTTree,
		"UPDATE/TTree": mutatingUpdateTTree,
	}
}
