package main

import (
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdV1Alpha1 "k8s.tars.io/api/crd/v1alpha1"
)

const (
	ResourceOutControlReason = "ResourceOutControlReason"

	ResourceDeleteReason = "ResourceDeleteReason"

	ResourceGetReason = "ResourceGetReason"
)

const (
	// ResourceOutControlError = "kind namespace/name already exists but not managed by namespace/name"
	ResourceOutControlError = "%s %s/%s already exists but not managed by %s/%s"

	//ResourceDeleteError = "delete kind namespace/name err: errMsg"
	ResourceDeleteError = "delete %s %s/%s err: %s"

	//ResourceGetError = "get kind namespace/name err: errMsg"
	ResourceGetError = "get %s %s:%s error: %s"

	//ResourceCreateError = "create kind namespace/name err: errMsg"
	ResourceCreateError = "create %s %s/%s error: %s"

	//ResourceUpdateError = "update kind namespace/name err: errMsg"
	ResourceUpdateError = "patch %s %s/%s error: %s"

	//ResourcePatchError = "update kind namespace/name err: errMsg"
	ResourcePatchError = "patch %s %s/%s error: %s"

	//ResourceSelectorError = "selector namespace/kind err: errMsg"
	ResourceSelectorError = "selector %s/%s error: %s"
)

const ServiceImagePlaceholder = " "

const TarsNodeLabelPrefix = "tars.io/node." // 此标签表示 该节点可以被 tars 使用

const TarsAbilityNodeLabelPrefix = "tars.io/ability." // 此标签表示 该节点可以被 tars 当做 App节点池使用
const TarsPublicNodeLabel = "tars.io/public"          // 此标签表示 该节点可以被 tars 当做 公用节点池使用

const TemplateLabel = "tars.io/Template"
const TSubTypeLabel = "tars.io/SubType"

const TServerAppLabel = "tars.io/ServerApp"
const TServerNameLabel = "tars.io/ServerName"
const TServerTagLabel = "tars.io/ServerTag"
const TConfigNameLabel = "tars.io/ConfigName"

const NodeServantName = "nodeobj"
const NodeServantPort = 19385

const TServerAPIVersion = "k8s.tars.io/v1alpha1"
const TServerKind = "TServer"

const TPodReadinessGate = "tars.io/active"

const K8SHostNameLabel = "kubernetes.io/hostname"

const ReleaseSourceLabel = "tars.io/ReleaseSource"
const ReleaseTagLabel = "tars.io/ReleaseTag"

const WebhookCertFile = "/etc/tarscontrol-cert/cert.pem"
const WebhookCertKey = "/etc/tarscontrol-cert/cert.key"

const TarsControlServiceAccount = "tars-tarscontrol"

const TarsTreeResourceName = "tars-tree"

func isOwnerByTServer(server *crdV1Alpha1.TServer, object k8sMetaV1.Object) bool {
	if ownerRef := object.GetOwnerReferences(); ownerRef != nil {
		for i := range ownerRef {
			if ownerRef[i].UID == server.UID {
				return true
			}
		}
	}
	return false
}
