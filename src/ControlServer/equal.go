package main

import (
	k8sAppsV1 "k8s.io/api/apps/v1"
	k8sCoreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	crdV1Alpha1 "k8s.tars.io/api/crd/v1alpha1"
)

func equalServicePort(l, r []k8sCoreV1.ServicePort) bool {

	if len(l) != len(r) {
		return false
	}

	for i := range l {
		if l[i].Name != r[i].Name {
			return false
		}
		if l[i].Port != r[i].Port {
			return false
		}
		if l[i].Protocol != r[i].Protocol {
			return false
		}
		if l[i].NodePort != r[i].NodePort {
			return false
		}
		if l[i].TargetPort.Type != r[i].TargetPort.Type {
			return false
		}
		if l[i].TargetPort.IntVal != r[i].TargetPort.IntVal {
			return false
		}
		if l[i].TargetPort.StrVal != r[i].TargetPort.StrVal {
			return false
		}
	}
	return true
}

func equalLabelSelector(l, r map[string]string) bool {
	if len(l) != len(r) {
		return false
	}
	for lk, lv := range l {
		if rv, ok := r[lk]; !ok || rv != lv {
			return false
		}
	}
	return true
}

func equalEnvFieldRef(l, r *k8sCoreV1.ObjectFieldSelector) bool {
	if l == nil {
		if r == nil {
			return true
		}
		return false
	}
	if r == nil {
		return false
	}

	if l.FieldPath != r.FieldPath {
		return false
	}
	return true
}

func equalEnvConfigMap(l, r *k8sCoreV1.ConfigMapKeySelector) bool {
	if l == nil {
		if r == nil {
			return true
		}
		return false
	}
	if r == nil {
		return false
	}
	if l.Key != r.Key {
		return false
	}
	if l.Name != r.Name {
		return false
	}
	return true
}

func equalEnvResourceFieldRef(l, r *k8sCoreV1.ResourceFieldSelector) bool {
	if l == nil {
		if r == nil {
			return true
		}
		return false
	}
	if r == nil {
		return false
	}
	if l.ContainerName != r.ContainerName {
		return false
	}
	if l.Resource != r.Resource {
		return false
	}
	if l.ContainerName != r.ContainerName {
		return false
	}
	return true
}

func equalEnvSecretKeyRef(l, r *k8sCoreV1.SecretKeySelector) bool {
	if l == nil {
		if r == nil {
			return true
		}
		return false
	}
	if r == nil {
		return false
	}
	if l.Key != r.Key {
		return false
	}
	if l.Name != r.Name {
		return false
	}
	return true
}

func equalEnv(l, r []k8sCoreV1.EnvVar) bool {
	if len(l) != len(r) {
		return false
	}
	for i := range l {
		if l[i].Name != r[i].Name {
			return false
		}
		if l[i].Value != r[i].Value {
			return false
		}

		if l[i].ValueFrom == nil {
			if r[i].ValueFrom == nil {
				continue
			}
			return false
		}

		if r[i].ValueFrom == nil {
			return false
		}

		if !equalEnvConfigMap(l[i].ValueFrom.ConfigMapKeyRef, r[i].ValueFrom.ConfigMapKeyRef) {
			return false
		}
		if !equalEnvFieldRef(l[i].ValueFrom.FieldRef, r[i].ValueFrom.FieldRef) {
			return false
		}
		if !equalEnvResourceFieldRef(l[i].ValueFrom.ResourceFieldRef, r[i].ValueFrom.ResourceFieldRef) {
			return false
		}
		if !equalEnvSecretKeyRef(l[i].ValueFrom.SecretKeyRef, r[i].ValueFrom.SecretKeyRef) {
			return false
		}
	}
	return true
}

func equalEnvFrom(l, r []k8sCoreV1.EnvFromSource) bool {
	if len(l) != len(r) {
		return false
	}
	for i := range l {
		if l[i].Prefix != r[i].Prefix {
			return false
		}
		if !equality.Semantic.DeepEqual(l[i].ConfigMapRef, r[i].ConfigMapRef) {
			return false
		}
		if !equality.Semantic.DeepEqual(l[i].SecretRef, r[i].SecretRef) {
			return false
		}
	}
	return true
}

func equalVolumesDownwardAPI(l, r *k8sCoreV1.DownwardAPIVolumeSource) bool {
	if l == nil {
		if r == nil {
			return true
		}
		return false
	}

	if r == nil {
		return false
	}

	if len(l.Items) != len(r.Items) {
		return false
	}
	for i := range l.Items {
		if l.Items[i].Path != r.Items[i].Path {
			return false
		}
		if !equality.Semantic.DeepEqual(l.Items[i].Mode, r.Items[i].Mode) {
			return false
		}
	}
	return true
}

func equalVolumesPVC(l, r *k8sCoreV1.PersistentVolumeClaimVolumeSource) bool {
	if l == nil {
		if r == nil {
			return true
		}
		return false
	}
	if r == nil {
		return false
	}

	if l.ClaimName != r.ClaimName {
		return false
	}

	if l.ReadOnly == r.ReadOnly {
		return false
	}

	return true
}

func equalVolumesHostPath(l, r *k8sCoreV1.HostPathVolumeSource) bool {
	if l == nil {
		if r == nil {
			return true
		}
		return false
	}
	if r == nil {
		return false
	}
	if l.Path != r.Path {
		return false
	}
	if !equality.Semantic.DeepEqual(l.Type, r.Type) {
		return false
	}
	return true
}

func equalVolumesEmptyDir(l, r *k8sCoreV1.EmptyDirVolumeSource) bool {
	if l == nil {
		if r == nil {
			return true
		}
		return false
	}
	if r == nil {
		return false
	}

	if l.Medium != r.Medium {
		return false
	}

	if !equality.Semantic.DeepEqual(l.SizeLimit, r.SizeLimit) {
		return false
	}

	return true
}

func equalVolumesSecret(l, r *k8sCoreV1.SecretVolumeSource) bool {
	if l == nil {
		if r == nil {
			return true
		}
		return false
	}
	if r == nil {
		return false
	}
	if len(l.Items) != len(r.Items) {
		return false
	}
	for i := range l.Items {
		if l.Items[i].Key != r.Items[i].Key {
			return false
		}
		if l.Items[i].Path != r.Items[i].Path {
			return false
		}
		if equality.Semantic.DeepEqual(l.Items[i].Mode, r.Items[i].Mode) {
			return false
		}
	}
	return true
}

func equalVolumesConfigMap(l, r *k8sCoreV1.ConfigMapVolumeSource) bool {
	if l == nil {
		if r == nil {
			return true
		}
		return false
	}
	if r == nil {
		return false
	}
	for i := range l.Items {
		if l.Items[i].Key != r.Items[i].Key {
			return false
		}
		if l.Items[i].Path != r.Items[i].Path {
			return false
		}
		if equality.Semantic.DeepEqual(l.Items[i].Mode, r.Items[i].Mode) {
			return false
		}
	}
	return true
}

func equalVolumes(l, r []k8sCoreV1.Volume) bool {
	if len(l) != len(r) {
		return false
	}

	for i := range l {
		if l[i].Name != r[i].Name {
			return false
		}
		if !equalVolumesConfigMap(l[i].ConfigMap, r[i].ConfigMap) {
			return false
		}
		if !equalVolumesSecret(l[i].Secret, r[i].Secret) {
			return false
		}
		if !equalVolumesEmptyDir(l[i].EmptyDir, r[i].EmptyDir) {
			return false
		}
		if !equalVolumesHostPath(l[i].HostPath, r[i].HostPath) {
			return false
		}
		if !equalVolumesPVC(l[i].PersistentVolumeClaim, r[i].PersistentVolumeClaim) {
			return false
		}
		if !equalVolumesDownwardAPI(l[i].DownwardAPI, r[i].DownwardAPI) {
			return false
		}
	}
	return true
}

func equalVolumesMounts(l, r []k8sCoreV1.VolumeMount) bool {
	if len(l) != len(r) {
		return false
	}
	for i := range l {
		if l[i].Name != r[i].Name {
			return false
		}
		if l[i].MountPath != r[i].MountPath {
			return false
		}
		if l[i].SubPath != r[i].SubPath {
			return false
		}
		if l[i].SubPathExpr != r[i].SubPathExpr {
			return false
		}
		if l[i].ReadOnly != r[i].ReadOnly {
			return false
		}
		if !equality.Semantic.DeepEqual(l[i].MountPropagation, r[i].MountPropagation) {
			return false
		}
	}
	return true
}

func equalContainerPorts(l, r []k8sCoreV1.ContainerPort) bool {
	if len(l) != len(r) {
		return false
	}
	for i := range l {
		if l[i].Name != r[i].Name {
			return false
		}
		if l[i].Protocol != r[i].Protocol {
			return false
		}
		if l[i].HostPort != r[i].HostPort {
			return false
		}
		if l[i].HostIP != r[i].HostIP {
			return false
		}
		if l[i].ContainerPort != r[i].ContainerPort {
			return false
		}
	}
	return true
}

func equalTarsServants(l, r []crdV1Alpha1.TServant) bool {
	if len(l) != len(r) {
		return false
	}

	for i := range l {
		if l[i].Name != r[i].Name {
			return false
		}
		if l[i].Timeout != r[i].Timeout {
			return false
		}
		if l[i].Connection != r[i].Connection {
			return false
		}
		if l[i].Thread != r[i].Thread {
			return false
		}
		if l[i].Port != r[i].Port {
			return false
		}
		if l[i].IsTars != r[i].IsTars {
			return false
		}
		if l[i].IsTcp != r[i].IsTcp {
			return false
		}
	}

	return true
}

func equalTars(l, r *crdV1Alpha1.TServerTars) bool {

	if l == nil {
		if r == nil {
			return true
		}
		return false
	}

	if r == nil {
		return false
	}

	if l.AsyncThread != r.AsyncThread {
		return false
	}

	if l.Foreground != r.Foreground {
		return false
	}
	if l.Profile != r.Profile {
		return false
	}

	if l.Template != r.Template {
		return false
	}

	if !equalTarsServants(l.Servants, r.Servants) {
		return false
	}

	return true
}

func equalTServerPorts(l, r []crdV1Alpha1.TServerNormalPort) bool {

	if len(l) != len(r) {
		return false
	}

	for i := range l {
		if l[i].Name != r[i].Name {
			return false
		}
		if l[i].IsTcp != r[i].IsTcp {
			return false
		}
	}
	return true
}

func equalNormal(l, r *crdV1Alpha1.TServerNormal) bool {
	if l == nil {
		if r == nil {
			return true
		}
		return false
	}

	if r == nil {
		return false
	}

	if !equalTServerPorts(l.Ports, r.Ports) {
		return false
	}

	return true
}

func equalK8SHostPorts(l, r []crdV1Alpha1.TK8SHostPort) bool {
	if len(l) != len(r) {
		return false
	}
	for i := range l {
		if l[i].Port != r[i].Port {
			return false
		}
		if l[i].NameRef != r[i].NameRef {
			return false
		}
	}
	return true
}

func equalTServerAndTEndpoint(server *crdV1Alpha1.TServer, endpoint *crdV1Alpha1.TEndpoint) bool {

	if server.Spec.App != endpoint.Spec.App {
		return false
	}

	if server.Spec.Server != endpoint.Spec.Server {
		return false
	}

	if server.Spec.Important != endpoint.Spec.Important {
		return false
	}

	if !equalK8SHostPorts(server.Spec.K8S.HostPorts, endpoint.Spec.HostPorts) {
		return false
	}

	switch server.Spec.SubType {
	case crdV1Alpha1.TARS:
		return equalTars(server.Spec.Tars, endpoint.Spec.Tars)
	case crdV1Alpha1.Normal:
		return equalNormal(server.Spec.Normal, endpoint.Spec.Normal)
	}

	//should not reach here
	return false
}

func equalTServerAndService(server *crdV1Alpha1.TServer, service *k8sCoreV1.Service) bool {
	serviceSpec := &service.Spec
	tarServerSpec := &server.Spec

	if serviceSpec.Type != k8sCoreV1.ServiceTypeClusterIP || serviceSpec.ClusterIP != k8sCoreV1.ClusterIPNone {
		return false
	}

	targetSelector := map[string]string{
		TServerAppLabel:  tarServerSpec.App,
		TServerNameLabel: tarServerSpec.Server,
	}

	if !equalLabelSelector(targetSelector, serviceSpec.Selector) {
		return false
	}

	if server.Spec.Tars != nil {
		targetPorts := buildServicePortsFromTServant(server)
		if !equalServicePort(targetPorts, serviceSpec.Ports) {
			return false
		}

	} else if server.Spec.Normal != nil {
		targetPorts := buildServicePortsFromNPorts(server)
		if !equalServicePort(targetPorts, serviceSpec.Ports) {
			return false
		}
	}
	return true
}

func equalTServerAndDaemonSet(server *crdV1Alpha1.TServer, daemonset *k8sAppsV1.DaemonSet) bool {
	daemonSetSpec := &daemonset.Spec
	tServerSpec := &server.Spec

	targetSelector := map[string]string{
		TServerAppLabel:  server.Spec.App,
		TServerNameLabel: server.Spec.Server,
	}

	if !equalLabelSelector(targetSelector, daemonSetSpec.Selector.MatchLabels) {
		return false
	}

	daemonSetSpecTemplateSpec := &daemonSetSpec.Template.Spec

	if tServerSpec.K8S.HostIPC != daemonSetSpecTemplateSpec.HostIPC {
		return false
	}

	if tServerSpec.K8S.HostNetwork != daemonSetSpecTemplateSpec.HostNetwork {
		return false
	}

	if tServerSpec.K8S.ServiceAccount != daemonSetSpecTemplateSpec.ServiceAccountName {
		return false
	}

	targetImagePullSecrets := buildImagePullSecrets(server)
	if !equality.Semantic.DeepEqual(targetImagePullSecrets, daemonSetSpecTemplateSpec.ImagePullSecrets) {
		return false
	}

	targetAffinity := buildPodAffinity(server)
	if !equality.Semantic.DeepEqual(targetAffinity, daemonSetSpecTemplateSpec.Affinity) {
		return false
	}

	targetVolumes, targetVolumeMounts := buildMounts(server)
	if !equalVolumes(targetVolumes, daemonSetSpecTemplateSpec.Volumes) {
		return false
	}

	if len(daemonSetSpecTemplateSpec.Containers) == 0 {
		return false
	}

	if !equalVolumesMounts(targetVolumeMounts, daemonSetSpecTemplateSpec.Containers[0].VolumeMounts) {
		return false
	}

	targetReadinessGates := buildReadinessGates(server)
	if !equality.Semantic.DeepEqual(targetReadinessGates, daemonSetSpecTemplateSpec.ReadinessGates) {
		return false
	}

	serverImage := ServiceImagePlaceholder
	if server.Spec.Release != nil {
		serverImage = server.Spec.Release.Image
	}

	if serverImage != daemonSetSpecTemplateSpec.Containers[0].Image {
		return false
	}

	if !equalEnv(tServerSpec.K8S.Env, daemonSetSpecTemplateSpec.Containers[0].Env) {
		return false
	}

	if !equalEnvFrom(tServerSpec.K8S.EnvFrom, daemonSetSpecTemplateSpec.Containers[0].EnvFrom) {
		return false
	}

	targetContainerPorts := buildContainerPort(server)
	if !equalContainerPorts(targetContainerPorts, daemonSetSpecTemplateSpec.Containers[0].Ports) {
		return false
	}
	return true
}

func equalTServerAndStatefulSet(server *crdV1Alpha1.TServer, statefulSet *k8sAppsV1.StatefulSet) bool {
	statefulSetSpec := &statefulSet.Spec
	tServerSpec := &server.Spec

	if tServerSpec.K8S.Replicas != *statefulSetSpec.Replicas {
		return false
	}

	targetSelector := map[string]string{
		TServerAppLabel:  server.Spec.App,
		TServerNameLabel: server.Spec.Server,
	}

	if !equalLabelSelector(targetSelector, statefulSetSpec.Selector.MatchLabels) {
		return false
	}

	statefulSetSpecTemplateSpec := &statefulSetSpec.Template.Spec

	if tServerSpec.K8S.HostIPC != statefulSetSpecTemplateSpec.HostIPC {
		return false
	}

	if tServerSpec.K8S.HostNetwork != statefulSetSpecTemplateSpec.HostNetwork {
		return false
	}

	if tServerSpec.K8S.ServiceAccount != statefulSetSpecTemplateSpec.ServiceAccountName {
		return false
	}

	targetImagePullSecrets := buildImagePullSecrets(server)
	if !equality.Semantic.DeepEqual(targetImagePullSecrets, statefulSet.Spec.Template.Spec.ImagePullSecrets) {
		return false
	}

	targetAffinity := buildPodAffinity(server)
	if !equality.Semantic.DeepEqual(targetAffinity, statefulSet.Spec.Template.Spec.Affinity) {
		return false
	}

	targetVolumes, targetVolumeMounts := buildMounts(server)
	if !equalVolumes(targetVolumes, statefulSetSpecTemplateSpec.Volumes) {
		return false
	}

	if len(statefulSetSpec.Template.Spec.Containers) == 0 {
		return false
	}

	if !equalVolumesMounts(targetVolumeMounts, statefulSetSpecTemplateSpec.Containers[0].VolumeMounts) {
		return false
	}

	tarReadinessGates := buildReadinessGates(server)
	if !equality.Semantic.DeepEqual(tarReadinessGates, statefulSetSpecTemplateSpec.ReadinessGates) {
		return false
	}

	serverImage := ServiceImagePlaceholder
	if server.Spec.Release != nil {
		serverImage = server.Spec.Release.Image
	}

	if serverImage != statefulSetSpecTemplateSpec.Containers[0].Image {
		return false
	}

	if !equalEnv(tServerSpec.K8S.Env, statefulSetSpecTemplateSpec.Containers[0].Env) {
		return false
	}

	if !equalEnvFrom(tServerSpec.K8S.EnvFrom, statefulSetSpecTemplateSpec.Containers[0].EnvFrom) {
		return false
	}

	targetContainerPorts := buildContainerPort(server)
	if !equalContainerPorts(targetContainerPorts, statefulSetSpecTemplateSpec.Containers[0].Ports) {
		return false
	}
	return true
}
