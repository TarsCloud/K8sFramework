package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	k8sAdmissionV1 "k8s.io/api/admission/v1"
	crdV1Alpha1 "k8s.tars.io/api/crd/v1alpha1"
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

	if tdeploy.Labels == nil {
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels\",\"value\":{}}")))
	}
	patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/tars.io~1Approve\",\"value\":\"Pending\"}")))

	if tdeploy.Approve != nil {
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"remove\",\"path\":\"/approve\"}")))
	}

	if tdeploy.Deployed != nil {
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"remove\",\"path\":\"/deployed\"}")))
	}

	if len(patchContents) == 1 {
		return nil
	}

	totalPatchContent := bytes.Join(patchContents, []byte{','})
	totalPatchContent[1] = ' '
	totalPatchContent = append(totalPatchContent, ']')
	return totalPatchContent
}

func mutatingUpdateTDeploy(requestAdmissionView *k8sAdmissionV1.AdmissionReview) []byte {
	tdeploy := &crdV1Alpha1.TDeploy{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, tdeploy)

	var patchContents = make([][]byte, 0, 10)

	patchContents = append(patchContents, []byte{'['})

	if tdeploy.Labels == nil {
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels\",\"value\":{}}")))
	}

	if tdeploy.Approve == nil {
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/tars.io~1Approve\",\"value\":\"Pending\"}")))
	} else if tdeploy.Approve.Result {
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/tars.io~1Approve\",\"value\":\"Approved\"}")))
	} else {
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/tars.io~1Approve\",\"value\":\"Reject\"}")))
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
	tserver := &crdV1Alpha1.TServer{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, tserver)

	var patchContents = make([][]byte, 0, 10)
	patchContents = append(patchContents, []byte{'['})

	if tserver.Labels == nil {
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels\",\"value\":{}}")))
	}

	patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\": \"add\", \"path\": \"/metadata/labels/tars.io~1ServerApp\", \"value\": \"%s\"}", tserver.Spec.App)))

	patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\": \"add\", \"path\": \"/metadata/labels/tars.io~1ServerName\", \"value\": \"%s\"}", tserver.Spec.Server)))

	patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\": \"add\", \"path\": \"/metadata/labels/tars.io~1SubType\", \"value\": \"%s\"}", tserver.Spec.SubType)))

	if tserver.Spec.Tars != nil {
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\": \"add\", \"path\": \"/metadata/labels/tars.io~1Template\", \"value\": \"%s\"}", tserver.Spec.Tars.Template)))
		if !tserver.Spec.Tars.Foreground && tserver.Spec.K8S.ReadinessGate == nil {
			patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\": \"add\", \"path\":\"/spec/k8s/readinessGate\",\"value\":\"%s\"}", TPodReadinessGate)))
		}
	}

	if tserver.Spec.Release == nil {
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\": \"add\", \"path\":\"/spec/k8s/replicas\",\"value\": %d}", 0)))
	} else {
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\": \"add\", \"path\": \"/metadata/labels/tars.io~1ReleaseSource\", \"value\": \"%s\"}", tserver.Spec.Release.Source)))
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\": \"add\", \"path\": \"/metadata/labels/tars.io~1ReleaseTag\", \"value\": \"%s\"}", tserver.Spec.Release.Tag)))
	}

	if len(tserver.Spec.K8S.HostPorts) > 0 || tserver.Spec.K8S.HostIPC {
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\": \"add\", \"path\":\"/spec/k8s/notStacked\",\"value\":%t}", true)))
	}

	if len(patchContents) == 1 {
		return nil
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

	var patchContents = make([][]byte, 0, 5)
	patchContents = append(patchContents, []byte{'['})

	if tconfig.Labels == nil {
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels\",\"value\":{}}")))
	}

	if tconfig.ServerConfig != nil {
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/tars.io~1ServerApp\",\"value\":\"%s\"}", tconfig.ServerConfig.App)))
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/tars.io~1ServerName\",\"value\":\"%s\"}", tconfig.ServerConfig.Server)))
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/tars.io~1ConfigName\",\"value\":\"%s\"}", tconfig.ServerConfig.ConfigName)))
		if tconfig.ServerConfig.PodSeq == nil {
			patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/tars.io~1PodSeq\",\"value\":\"%s\"}", "m")))
		} else {
			patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/tars.io~1PodSeq\",\"value\":\"%s\"}", *tconfig.ServerConfig.PodSeq)))
		}
	}

	if tconfig.AppConfig != nil {
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/tars.io~1ServerApp\",\"value\":\"%s\"}", tconfig.AppConfig.App)))
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/tars.io~1ServerName\",\"value\":\"%s\"}", "")))
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/tars.io~1ConfigName\",\"value\":\"%s\"}", tconfig.AppConfig.ConfigName)))
		patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"add\",\"path\":\"/metadata/labels/tars.io~1PodSeq\",\"value\":\"%s\"}", "m")))
	}

	if len(patchContents) == 1 {
		return nil
	}

	totalPatchContent := bytes.Join(patchContents, []byte{','})
	totalPatchContent[1] = ' '
	totalPatchContent = append(totalPatchContent, ']')
	return totalPatchContent
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
				patchContents = append(patchContents, []byte(fmt.Sprintf("{\"op\":\"replace\",\"path\":\"/apps/%d\",\"value\":\"%s\"}", i, bs)))
			}
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
