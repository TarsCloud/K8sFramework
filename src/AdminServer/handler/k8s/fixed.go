package k8s

const (
	ResourceOutControlReason = "ResourceOutControlReason"
	ResourceDeleteReason = "ResourceDeleteReason"
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

const TafNodeLabelPrefix = "taf.io/node." // 此标签表示 该节点可以被 taf 使用

const TafAbilityNodeLabelPrefix = "taf.io/ability." // 此标签表示 该节点可以被 taf 当做 App节点池使用
const TafPublicNodeLabel = "taf.io/public"          // 此标签表示 该节点可以被 taf 当做 公用节点池使用

const TemplateLabel = "taf.io/Template"
const TSubTypeLabel = "taf.io/SubType"

const TServerAppLabel = "taf.io/ServerApp"
const TServerNameLabel = "taf.io/ServerName"
const TServerTagLabel = "taf.io/ServerTag"
const TConfigNameLabel = "taf.io/ConfigName"
const TDeployApproveLabel = "taf.io/Approve"

const NodeServantName = "nodeobj"
const NodeServantPort = 19385

const TServerAPIVersion = "k8s.taf.io/v1alpha1"
const TServerKind = "TServer"

const TPodReadinessGate = "taf.io/active"

const K8SHostNameLabel = "kubernetes.io/hostname"

const ReleaseSourceLabel = "taf.io/ReleaseSource"
const ReleaseTagLabel = "taf.io/ReleaseTag"

//const WebhookCertFile = "/etc/tafoperator-cert/cert.pem"
//const WebhookCertKey = "/etc/tafoperator-cert/cert.key"

const WebhookCertFile = "/home/adugeek/Cert/cert.pem"
const WebhookCertKey = "/home/adugeek/Cert/cert.key"

const TafControlServiceAccount = "taf-tafcontrol"

// -----------------------------------------------------

const TafAgentDaemonSetName = "taf-tafagent"

const TafTreeName = "taf-tree"

