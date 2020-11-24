package main

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"net/http"
	"strings"

	k8sAdmissionV1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilRuntime "k8s.io/apimachinery/pkg/util/runtime"
	crdV1Alpha1 "k8s.tars.io/api/crd/v1alpha1"
)

func validTDeploy(newTdeploy *crdV1Alpha1.TDeploy, oldTDeploy *crdV1Alpha1.TDeploy, option *K8SOption, watcher *Watcher) error {

	if oldTDeploy != nil {
		if oldTDeploy.Approve != nil {
			if !equality.Semantic.DeepEqual(newTdeploy.Apply, oldTDeploy.Apply) {
				return fmt.Errorf("the value of /apply cannot be changed")
			}

			if !equality.Semantic.DeepEqual(newTdeploy.Approve, oldTDeploy.Approve) {
				return fmt.Errorf("the value of /approve cannot be changed")
			}
		}

		if oldTDeploy.Deployed != nil {
			if !equality.Semantic.DeepEqual(newTdeploy.Deployed, oldTDeploy.Deployed) {
				return fmt.Errorf("the value of /deployed cannot be changed")
			}
		}
	}

	targetTServerName := fmt.Sprintf("%s-%s", strings.ToLower(newTdeploy.Apply.App), strings.ToLower(newTdeploy.Apply.Server))
	_, err := watcher.tServerLister.TServers(option.namespace).Get(targetTServerName)

	if err == nil {
		return fmt.Errorf("tserver/%s areadly exist", targetTServerName)
	}

	if !errors.IsNotFound(err) {
		return fmt.Errorf("get tserver/%s error : %s, try again latter ", targetTServerName, err.Error())
	}

	fakeTServer := &crdV1Alpha1.TServer{
		ObjectMeta: k8sMetaV1.ObjectMeta{
			Name:      targetTServerName,
			Namespace: option.namespace,
		},
		Spec: newTdeploy.Apply,
	}
	return validTServer(fakeTServer, nil, option, watcher)
}

func (v *Validating) validCreateTDeploy(requestAdmissionView *k8sAdmissionV1.AdmissionReview) error {
	tdeploy := &crdV1Alpha1.TDeploy{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, tdeploy)
	return validTDeploy(tdeploy, nil, v.k8sOption, v.watcher)
}

func (v *Validating) validUpdateTDeploy(requestAdmissionView *k8sAdmissionV1.AdmissionReview) error {
	newTDeploy := &crdV1Alpha1.TDeploy{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, newTDeploy)

	oldTDeploy := &crdV1Alpha1.TDeploy{}
	_ = json.Unmarshal(requestAdmissionView.Request.OldObject.Raw, oldTDeploy)

	if oldTDeploy.Approve == nil && newTDeploy.Approve != nil {
		//todo checkout account
		return nil
	}

	if oldTDeploy.Deployed == nil && newTDeploy.Deployed != nil {
		if requestAdmissionView.Request.UserInfo.Username != v.controlAccount {
			return fmt.Errorf("only use authorizedAccount can update /deployed")
		}
	}
	return validTDeploy(newTDeploy, oldTDeploy, v.k8sOption, v.watcher)
}

func (v *Validating) validDeleteTDeploy(view *k8sAdmissionV1.AdmissionReview) error {
	return nil
}

func validTServer(newTServer *crdV1Alpha1.TServer, oldTServer *crdV1Alpha1.TServer, option *K8SOption, watcher *Watcher) error {

	if oldTServer != nil {
		if newTServer.Spec.App != oldTServer.Spec.App {
			return fmt.Errorf("the value of /spec/app cannot be changed")
		}

		if newTServer.Spec.Server != oldTServer.Spec.Server {
			return fmt.Errorf("the value of /spec/server cannot be changed")
		}

		if newTServer.Spec.SubType != oldTServer.Spec.SubType {
			return fmt.Errorf("the value of /spec/subType cannot be changed")
		}

		if oldTServer.Spec.Tars != nil {
			if newTServer.Spec.Tars == nil {
				return fmt.Errorf("the value of /spec/tars cannot be changed")
			}
		}

		if oldTServer.Spec.Normal != nil {
			if newTServer.Spec.Normal == nil {
				return fmt.Errorf("the value of /spec/normal cannot be changed")
			}
		}
	}

	if newTServer.Name != strings.ToLower(newTServer.Spec.App)+"-"+strings.ToLower(newTServer.Spec.Server) {
		return fmt.Errorf("unexpected resources name")
	}

	var portNames map[string]interface{}
	var portValues map[int32]interface{}

	if newTServer.Spec.Tars != nil {
		servants := newTServer.Spec.Tars.Servants

		if servants == nil || len(servants) < 1 {
			return fmt.Errorf("servants should not empty")
		}

		portNames = make(map[string]interface{}, len(servants)+1)
		portValues = make(map[int32]interface{}, len(servants)+1)

		for i := range servants {

			servantName := strings.ToLower(servants[i].Name)

			if servantName == NodeServantName {
				return fmt.Errorf("servants name value should not equal %s", NodeServantName)
			}

			if _, ok := portNames[servantName]; ok {
				return fmt.Errorf("duplicate servants name value")
			}

			if servants[i].Port == NodeServantPort {
				return fmt.Errorf("servants port value should not equal %d ", NodeServantPort)
			}

			if _, ok := portValues[servants[i].Port]; ok {
				return fmt.Errorf("duplicate servants port value")
			}

			portNames[servantName] = nil
			portValues[servants[i].Port] = nil
		}

		templateName := newTServer.Spec.Tars.Template
		_, err := watcher.tTemplateLister.TTemplates(option.namespace).Get(templateName)
		if err != nil {
			if !errors.IsNotFound(err) {
				return fmt.Errorf("get ttemplate/%s error: %s, try again latter ", templateName, err.Error())
			}
			return fmt.Errorf("ttemplate/%s not exist", templateName)
		}
	} else if newTServer.Spec.Normal != nil {

		ports := newTServer.Spec.Normal.Ports

		if ports == nil || len(ports) == 0 {
			return fmt.Errorf("ports should not empty")
		}

		portNames = make(map[string]interface{}, len(ports)+1)
		portValues = make(map[int32]interface{}, len(ports)+1)

		for i := range ports {
			name := strings.ToLower(ports[i].Name)

			if _, ok := portNames[name]; ok {
				return fmt.Errorf("duplicate ports name value")
			}

			if _, ok := portValues[ports[i].Port]; ok {
				return fmt.Errorf("duplicate ports port value")
			}

			portNames[name] = nil
			portValues[ports[i].Port] = nil
		}
	}

	if newTServer.Spec.K8S.HostPorts != nil {
		hostPorts := newTServer.Spec.K8S.HostPorts
		hostPortPorts := make(map[int32]interface{}, len(hostPorts))

		for i := range hostPorts {
			nameRef := strings.ToLower(hostPorts[i].NameRef)
			if _, ok := portNames[nameRef]; !ok {
				return fmt.Errorf("k8s.hostPort[%d].objRef value should in servants", i)
			}

			if _, ok := hostPortPorts[hostPorts[i].Port]; ok {
				return fmt.Errorf("duplicate hostPort.port value %d", hostPorts[i].Port)
			}
		}
	}

	if newTServer.Spec.Release != nil {

		release, err := watcher.tReleaseLister.TReleases(option.namespace).Get(newTServer.Spec.Release.Source)

		if err != nil {
			if !errors.IsNotFound(err) {
				return fmt.Errorf("get trelease/%s error: %s, try again latter ", newTServer.Name, err.Error())
			}
			return fmt.Errorf("trelease/%s not exist", newTServer.Spec.Release.Source)
		}

		hadCheckTag := false

		for _, version := range release.Spec.List {
			if !hadCheckTag && newTServer.Spec.Release.Tag == version.Tag {
				if newTServer.Spec.Release.Image != version.Image {
					return fmt.Errorf("")
				}
				if newTServer.Spec.Release.ImagePullSecret != version.ImagePullSecret {
				}
				hadCheckTag = true
			}
		}

		if !hadCheckTag {
			return fmt.Errorf("trelease/%s[%s] not exist", newTServer.Spec.Release.Source, newTServer.Spec.Release.Tag)
		}
	}

	return nil
}

func (v *Validating) validCreateTServer(requestAdmissionView *k8sAdmissionV1.AdmissionReview) error {
	newTServer := &crdV1Alpha1.TServer{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, newTServer)
	return validTServer(newTServer, nil, v.k8sOption, v.watcher)
}

func (v *Validating) validUpdateTServer(requestAdmissionView *k8sAdmissionV1.AdmissionReview) error {
	newTServer := &crdV1Alpha1.TServer{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, newTServer)

	oldTServer := &crdV1Alpha1.TServer{}
	_ = json.Unmarshal(requestAdmissionView.Request.OldObject.Raw, oldTServer)

	return validTServer(newTServer, oldTServer, v.k8sOption, v.watcher)
}

func (v *Validating) validDeleteTServer(view *k8sAdmissionV1.AdmissionReview) error {
	return nil
}

func (v *Validating) validCreateTEndpoint(requestAdmissionView *k8sAdmissionV1.AdmissionReview) error {
	if requestAdmissionView.Request.UserInfo.Username == v.controlAccount {
		return nil
	}
	return fmt.Errorf("only use authorizedAccount can create tendpoints")
}

func (v *Validating) validUpdateTEndpoint(requestAdmissionView *k8sAdmissionV1.AdmissionReview) error {
	if requestAdmissionView.Request.UserInfo.Username == v.controlAccount {
		return nil
	}
	return fmt.Errorf("only use authorizedAccount can update tendpoints")
}

func (v *Validating) validDeleteTEndpoint(requestAdmissionView *k8sAdmissionV1.AdmissionReview) error {
	if requestAdmissionView.Request.UserInfo.Username == v.controlAccount || requestAdmissionView.Request.UserInfo.Username != v.garbageCollectorAccount {
		return nil
	}
	return fmt.Errorf("only use authorizedAccount can delate tendpoints")
}

func validTConfig(newTConfig *crdV1Alpha1.TConfig, oldTConfig *crdV1Alpha1.TConfig, option *K8SOption) error {

	if newTConfig.AppConfig != nil {
		if newTConfig.Name != fmt.Sprintf("%s-%d", strings.ToLower(newTConfig.AppConfig.App), crc32.ChecksumIEEE([]byte(newTConfig.AppConfig.ConfigName))) {
			return fmt.Errorf("unexpected resources name")
		}
		return nil
	}

	if newTConfig.ServerConfig != nil {
		if newTConfig.ServerConfig.PodSeq == nil {
			if newTConfig.Name != fmt.Sprintf(
				"%s-%s-%d",
				strings.ToLower(newTConfig.ServerConfig.App),
				strings.ToLower(newTConfig.ServerConfig.Server),
				crc32.ChecksumIEEE([]byte(newTConfig.ServerConfig.ConfigName)),
			) {
				return fmt.Errorf("unexpected resources name")
			}
			return nil
		}

		if newTConfig.Name != fmt.Sprintf(
			"%s-%s-%d-%s",
			strings.ToLower(newTConfig.ServerConfig.App),
			strings.ToLower(newTConfig.ServerConfig.Server),
			crc32.ChecksumIEEE([]byte(newTConfig.ServerConfig.ConfigName)),
			*newTConfig.ServerConfig.PodSeq,
		) {
			return fmt.Errorf("unexpected resources name")
		}

		masterConfigName := fmt.Sprintf(
			"%s-%s-%d",
			strings.ToLower(newTConfig.ServerConfig.App),
			strings.ToLower(newTConfig.ServerConfig.Server),
			crc32.ChecksumIEEE([]byte(newTConfig.ServerConfig.ConfigName)),
		)
		_, err := option.crdClientSet.CrdV1alpha1().TConfigs(option.namespace).Get(context.TODO(), masterConfigName, k8sMetaV1.GetOptions{})
		if err != nil {
			return fmt.Errorf(ResourceGetError, "tconfig", option.namespace, masterConfigName, err)
		}
		return nil
	}

	return fmt.Errorf("should not reach place")
}

func (v *Validating) validCreateTConfig(requestAdmissionView *k8sAdmissionV1.AdmissionReview) error {
	newTConfig := &crdV1Alpha1.TConfig{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, newTConfig)
	return validTConfig(newTConfig, nil, v.k8sOption)
}

func (v *Validating) validUpdateTConfig(requestAdmissionView *k8sAdmissionV1.AdmissionReview) error {
	newTConfig := &crdV1Alpha1.TConfig{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, newTConfig)

	oldTConfig := &crdV1Alpha1.TConfig{}
	_ = json.Unmarshal(requestAdmissionView.Request.OldObject.Raw, oldTConfig)

	return validTConfig(newTConfig, oldTConfig, v.k8sOption)
}

func (v *Validating) validDeleteTConfig(requestAdmissionView *k8sAdmissionV1.AdmissionReview) error {
	tconfig := &crdV1Alpha1.TConfig{}
	_ = json.Unmarshal(requestAdmissionView.Request.OldObject.Raw, tconfig)

	if tconfig.AppConfig != nil {
		return nil
	}

	if tconfig.ServerConfig != nil && tconfig.ServerConfig.PodSeq != nil {
		return nil
	}

	appRequirement, _ := labels.NewRequirement(TServerAppLabel, "==", []string{tconfig.ServerConfig.App})
	serverRequirement, _ := labels.NewRequirement(TServerNameLabel, "==", []string{tconfig.ServerConfig.Server})
	configNameRequirement, _ := labels.NewRequirement(TConfigNameLabel, "==", []string{tconfig.ServerConfig.ConfigName})

	labelsSelector := labels.NewSelector().Add(*appRequirement).Add(*serverRequirement).Add(*configNameRequirement).String()
	tconfigs, err := v.k8sOption.crdClientSet.CrdV1alpha1().TConfigs(tconfig.Namespace).List(context.TODO(), k8sMetaV1.ListOptions{
		LabelSelector: labelsSelector,
	})

	if err != nil {
		return fmt.Errorf("selector tconfig error: %s, try again later", err.Error())
	}

	for i := range tconfigs.Items {
		if tconfigs.Items[i].ServerConfig == nil {
			return fmt.Errorf("should not reach place")
		}
		if tconfigs.Items[i].ServerConfig.PodSeq != nil {
			return fmt.Errorf("cannot delete tconfig/%s because it is referance by another tconfig", requestAdmissionView.Request.Name)
		}
	}

	return nil
}

func validTTemplate(newTTemplate *crdV1Alpha1.TTemplate, oldTTemplate *crdV1Alpha1.TTemplate, option *K8SOption, watcher *Watcher) error {

	parentName := newTTemplate.Spec.Parent
	if parentName == "" {
		return fmt.Errorf("the valule /spec/parent should not empty ")
	}

	if newTTemplate.Name == newTTemplate.Spec.Parent {
		return nil
	}

	if _, err := watcher.tTemplateLister.TTemplates(option.namespace).Get(parentName); err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("ttemplate/%s not exist", parentName)
		}
		return fmt.Errorf("get ttemplate/%s error, try again latter ", parentName)
	}
	return nil
}

func (v *Validating) validCreateTTemplate(requestAdmissionView *k8sAdmissionV1.AdmissionReview) error {
	newTTemplate := &crdV1Alpha1.TTemplate{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, newTTemplate)
	return validTTemplate(newTTemplate, nil, v.k8sOption, v.watcher)
}

func (v *Validating) validUpdateTTemplate(requestAdmissionView *k8sAdmissionV1.AdmissionReview) error {
	newTTemplate := &crdV1Alpha1.TTemplate{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, newTTemplate)

	oldTTemplate := &crdV1Alpha1.TTemplate{}
	_ = json.Unmarshal(requestAdmissionView.Request.OldObject.Raw, oldTTemplate)

	return validTTemplate(newTTemplate, oldTTemplate, v.k8sOption, v.watcher)
}

func (v *Validating) validDeleteTTemplate(requestAdmissionView *k8sAdmissionV1.AdmissionReview) error {
	ttemplate := &crdV1Alpha1.TTemplate{}
	_ = json.Unmarshal(requestAdmissionView.Request.OldObject.Raw, ttemplate)
	requirement, _ := labels.NewRequirement(TemplateLabel, "==", []string{ttemplate.Name})
	tServers, err := v.watcher.tServerLister.List(labels.NewSelector().Add(*requirement))
	if err != nil {
		utilRuntime.HandleError(err)
		return err
	}
	if tServers != nil && len(tServers) != 0 {
		return fmt.Errorf("cannot delete ttemplate/%s because it is reference by some tserver", requestAdmissionView.Request.Name)
	}
	return nil
}

func validTRelease(newRelease *crdV1Alpha1.TRelease, oldRelease *crdV1Alpha1.TRelease, option *K8SOption, watcher *Watcher) error {
	newTReleaseVersionMap := make(map[string]*crdV1Alpha1.TReleaseVersion, len(newRelease.Spec.List))
	for _, pos := range newRelease.Spec.List {
		if _, ok := newTReleaseVersionMap[pos.Tag]; ok {
			return fmt.Errorf("duplicate tag value : %s", pos.Tag)
		}
		newTReleaseVersionMap[pos.Tag] = pos
	}

	if oldRelease == nil {
		return nil
	}

	for _, pos := range oldRelease.Spec.List {
		versionInNewTRelease, ok := newTReleaseVersionMap[pos.Tag]
		if ok {
			if pos.Image != versionInNewTRelease.Image {
				return fmt.Errorf("the value of /spec/list/tag/image cannot be changed")
			}
			if pos.ImagePullSecret != versionInNewTRelease.ImagePullSecret {
				return fmt.Errorf("the value of /spec/list/tag/imagePullSecret cannot be changed")
			}
		} else {
			releaseSourceMatch, _ := labels.NewRequirement(ReleaseSourceLabel, "==", []string{newRelease.Name})
			releaseTagMatch, _ := labels.NewRequirement(ReleaseTagLabel, "==", []string{pos.Tag})

			tservers, err := watcher.tServerLister.List(labels.NewSelector().Add(*releaseSourceMatch).Add(*releaseTagMatch))

			if err != nil {
				return fmt.Errorf("selector tserver error: %s, try again later", err.Error())
			}

			if tservers != nil && len(tservers) > 0 {
				return fmt.Errorf("cannot delete trelease/%s/spec/list/tag/%s ,because it is reference by some tserver", newRelease.Name, pos.Tag)
			}
		}
	}
	return nil
}

func (v *Validating) validCreateTRelease(requestAdmissionView *k8sAdmissionV1.AdmissionReview) error {
	newTRelease := &crdV1Alpha1.TRelease{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, newTRelease)
	return validTRelease(newTRelease, nil, v.k8sOption, v.watcher)
}

func (v *Validating) validUpdateTRelease(requestAdmissionView *k8sAdmissionV1.AdmissionReview) error {
	newTRelease := &crdV1Alpha1.TRelease{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, newTRelease)

	oldTRelease := &crdV1Alpha1.TRelease{}
	_ = json.Unmarshal(requestAdmissionView.Request.OldObject.Raw, oldTRelease)

	return validTRelease(newTRelease, oldTRelease, v.k8sOption, v.watcher)
}

func (v *Validating) validDeleteTRelease(requestAdmissionView *k8sAdmissionV1.AdmissionReview) error {
	oldTrelease := &crdV1Alpha1.TRelease{}
	_ = json.Unmarshal(requestAdmissionView.Request.OldObject.Raw, oldTrelease)

	fakeNewTRelease := &crdV1Alpha1.TRelease{
		Spec: crdV1Alpha1.TReleaseSpec{
			List: []*crdV1Alpha1.TReleaseVersion{},
		},
	}
	return validTRelease(fakeNewTRelease, oldTrelease, v.k8sOption, v.watcher)
}

func validTTree(newTTree *crdV1Alpha1.TTree, oldTTree *crdV1Alpha1.TTree, option *K8SOption, watcher *Watcher) error {
	businessMap := make(map[string]interface{}, len(newTTree.Businesses))
	for _, business := range newTTree.Businesses {
		if _, ok := businessMap[business.Name]; ok {
			return fmt.Errorf("duplicate business name : %s", business.Name)
		}
		businessMap[business.Name] = nil
	}

	appMap := make(map[string]interface{}, len(newTTree.Apps))
	for _, app := range newTTree.Apps {
		if _, ok := appMap[app.Name]; ok {
			return fmt.Errorf("duplicate app name : %s", app.Name)
		}
		if app.BusinessRef != "" {
			if _, ok := businessMap[app.BusinessRef]; !ok {
				return fmt.Errorf("business/%s is not exist", app.BusinessRef)
			}
		}
		appMap[app.Name] = nil
	}

	if oldTTree == nil {
		return nil
	}

	for i := range oldTTree.Apps {
		appName := oldTTree.Apps[i].Name
		if _, ok := appMap[appName]; !ok {
			requirement, _ := labels.NewRequirement(TServerAppLabel, "==", []string{appName})
			tservers, err := watcher.tServerLister.List(labels.NewSelector().Add(*requirement))
			if err != nil {
				utilRuntime.HandleError(err)
				return err
			}
			if tservers != nil && len(tservers) != 0 {
				return fmt.Errorf("cannot delete ttree/apps[%s] because it is reference by some tserver", appName)
			}
		}
	}
	return nil
}

func (v *Validating) validCreateTTree(view *k8sAdmissionV1.AdmissionReview) error {
	return fmt.Errorf("create ttree operation is defined")
}

func (v *Validating) validUpdateTTree(requestAdmissionView *k8sAdmissionV1.AdmissionReview) error {
	newTTree := &crdV1Alpha1.TTree{}
	_ = json.Unmarshal(requestAdmissionView.Request.Object.Raw, newTTree)

	oldTTree := &crdV1Alpha1.TTree{}
	_ = json.Unmarshal(requestAdmissionView.Request.OldObject.Raw, oldTTree)

	return validTTree(newTTree, oldTTree, v.k8sOption, v.watcher)
}

func (v *Validating) validDeleteTTree(view *k8sAdmissionV1.AdmissionReview) error {
	return fmt.Errorf("delete ttree operation is defined")
}

type Validating struct {
	k8sOption               *K8SOption
	watcher                 *Watcher
	controlAccount          string
	garbageCollectorAccount string
}

func (v Validating) handle(r *http.Request, w http.ResponseWriter) {
	requestAdmissionView := &k8sAdmissionV1.AdmissionReview{}

	err := json.NewDecoder(r.Body).Decode(requestAdmissionView)
	if err != nil {
		return
	}

	key := fmt.Sprintf("%s/%s", string(requestAdmissionView.Request.Operation), requestAdmissionView.Request.Kind.Kind)

	switch key {
	case "CREATE/TDeploy":
		err = v.validCreateTDeploy(requestAdmissionView)
	case "Update/TDeploy":
		err = v.validUpdateTDeploy(requestAdmissionView)
	case "Delete/TDeploy":
		err = v.validDeleteTDeploy(requestAdmissionView)

	case "CREATE/TServer":
		err = v.validCreateTServer(requestAdmissionView)
	case "UPDATE/TServer":
		err = v.validUpdateTServer(requestAdmissionView)
	case "DELETE/TServer":
		err = v.validDeleteTServer(requestAdmissionView)

	case "CREATE/TEndpoint":
		err = v.validCreateTEndpoint(requestAdmissionView)
	case "UPDATE/TEndpoint":
		err = v.validUpdateTEndpoint(requestAdmissionView)
	case "DELETE/TEndpoint":
		err = v.validDeleteTEndpoint(requestAdmissionView)

	case "CREATE/TConfig":
		err = v.validCreateTConfig(requestAdmissionView)
	case "UPDATE/TConfig":
		err = v.validUpdateTConfig(requestAdmissionView)
	case "DELETE/TConfig":
		err = v.validDeleteTConfig(requestAdmissionView)

	case "CREATE/TTemplate":
		err = v.validCreateTTemplate(requestAdmissionView)
	case "UPDATE/TTemplate":
		err = v.validUpdateTTemplate(requestAdmissionView)
	case "DELETE/TTemplate":
		err = v.validDeleteTTemplate(requestAdmissionView)

	case "CREATE/TRelease":
		err = v.validCreateTRelease(requestAdmissionView)
	case "UPDATE/TRelease":
		err = v.validUpdateTRelease(requestAdmissionView)
	case "DELETE/TRelease":
		err = v.validDeleteTRelease(requestAdmissionView)

	case "CREATE/TTree":
		err = v.validCreateTTree(requestAdmissionView)
	case "UPDATE/TTree":
		err = v.validUpdateTTree(requestAdmissionView)
	case "DELETE/TTree":
		err = v.validDeleteTTree(requestAdmissionView)
	}

	var responseAdmissionView *k8sAdmissionV1.AdmissionReview
	if err == nil {
		responseAdmissionView = &k8sAdmissionV1.AdmissionReview{
			Response: &k8sAdmissionV1.AdmissionResponse{
				UID:     requestAdmissionView.Request.UID,
				Allowed: true,
			},
		}
	} else {
		responseAdmissionView = &k8sAdmissionV1.AdmissionReview{
			Response: &k8sAdmissionV1.AdmissionResponse{
				UID:     requestAdmissionView.Request.UID,
				Allowed: false,
				Result: &k8sMetaV1.Status{
					Status:  "Failure",
					Message: err.Error(),
				},
			},
		}
	}
	responseBytes, _ := json.Marshal(responseAdmissionView)
	_, _ = w.Write(responseBytes)
}
