/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Tag 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	k8sCoreV1 "k8s.io/api/core/v1"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TServant struct {
	Name       string `json:"name"`
	Port       int32  `json:"port"`
	Thread     int32    `json:"thread"`
	Connection int32    `json:"connection"`
	Capacity   int32    `json:"capacity"`
	Timeout    int32    `json:"timeout"`
	IsTaf      bool   `json:"isTaf"`
	IsTcp      bool   `json:"isTcp"`
}

type TK8SMount struct {
	Name string `json:"name"`

	Source k8sCoreV1.VolumeSource `json:"source"`

	// Mounted read-only if true, read-write otherwise (false or unspecified).
	// Defaults to false.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty"`
	// Path within the container at which the volume should be mounted.  Must
	// not contain ':'.
	MountPath string `json:"mountPath"`
	// Path within the volume from which the container's volume should be mounted.
	// Defaults to "" (volume's root).
	// +optional
	SubPath string `json:"subPath,omitempty"`
	//// mountPropagation determines how mounts are propagated from the host
	//// to container and the other way around.
	//// When not set, MountPropagationNone is used.
	//// This field is beta in 1.10.
	//// +optional
	MountPropagation *string `json:"mountPropagation,omitempty"`
	// Expanded path within the volume from which the container's volume should be mounted.
	// Behaves similarly to SubPath but environment variable references $(VAR_NAME) are expanded using the container's environment.
	// Defaults to "" (volume's root).
	// SubPathExpr and SubPath are mutually exclusive.
	// +optional
	SubPathExpr string `json:"subPathExpr,omitempty"`
}

type TServerRelease struct {
	ServerType      string         `json:"serverType"`
	Source          string         `json:"source"`
	Tag             string         `json:"tag"`
	Image           string         `json:"image"`
	ImagePullSecret string         `json:"imagePullSecret,omitempty"`
	ActiveTime      k8sMetaV1.Time `json:"activeTime,omitempty"`
	ActiveReason    string         `json:"activeReason,omitempty"`
	ActivePerson    string         `json:"activePerson,omitempty"`
}

type TServerK8S struct {
	ServiceAccount string `json:"serviceAccount,omitempty"`

	Env []k8sCoreV1.EnvVar `json:"env,omitempty"`

	EnvFrom []k8sCoreV1.EnvFromSource `json:"envFrom,omitempty"`

	HostIPC bool `json:"hostIPC,omitempty"`

	HostNetwork bool `json:"hostNetwork,omitempty"`

	HostPorts []TK8SHostPort `json:"hostPorts,omitempty"`

	Mounts []TK8SMount `json:"mounts,omitempty"`

	NodeSelector TK8SNodeSelector `json:"nodeSelector"`

	NotStacked *bool `json:"notStacked,omitempty"`

	PodManagementPolicy string `json:"podManagementPolicy,omitempty"`

	Replicas int32 `json:"replicas"`

	ReadinessGate *string `json:"readinessGate,omitempty"`
}

type TK8SNodeSelectorKind struct {
	Values []string `json:"values"`
}

type TK8SNodeSelector struct {
	NodeBind    *TK8SNodeSelectorKind `json:"nodeBind,omitempty"`
	PublicPool  *TK8SNodeSelectorKind `json:"publicPool,omitempty"`
	AbilityPool *TK8SNodeSelectorKind `json:"abilityPool,omitempty"`
	DaemonSet   *TK8SNodeSelectorKind `json:"daemonSet,,omitempty"`
}

type TK8SHostPort struct {
	NameRef string `json:"nameRef"`
	Port    int32  `json:"port"`
}

type TServerExternalAddress struct {
	IP   string `json:"ip"`
	Port int32  `json:"port"`
}

type TServerExternalUPStream struct {
	Name      string                   `json:"name"`
	IsTcp     bool                     `json:"isTcp"`
	Addresses []TServerExternalAddress `json:"addresses"`
}

type TServerExternal struct {
	Upstreams []TServerExternalUPStream `json:"upstreams"`
}

type TServerTaf struct {
	Template    string     `json:"template"`
	Profile     string     `json:"profile"`
	Foreground  bool       `json:"foreground"`
	AsyncThread int32        `json:"asyncThread"`
	Servants    []TServant `json:"servants"`
}

type TServerNormalPort struct {
	Name  string `json:"name"`
	Port  int32  `json:"port"`
	IsTcp bool   `json:"isTcp"`
}

type TServerNormal struct {
	Ports []TServerNormalPort `json:"ports"`
}

type TServerSubType string

const (
	TAF    TServerSubType = "taf"
	Normal TServerSubType = "normal"
)

type TServerSpec struct {
	App       string          `json:"app"`
	Server    string          `json:"server"`
	SubType   TServerSubType  `json:"subType"`
	Important int32             `json:"important"`
	Taf       *TServerTaf     `json:"taf,omitempty"`
	Normal    *TServerNormal  `json:"normal,omitempty"`
	K8S       TServerK8S      `json:"k8s"`
	Release   *TServerRelease `json:"release,omitempty"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TServer struct {
	k8sMetaV1.TypeMeta   `json:",inline"`
	k8sMetaV1.ObjectMeta `json:"metadata,omitempty"`
	Spec                 TServerSpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TServerList struct {
	k8sMetaV1.TypeMeta `json:",inline"`
	k8sMetaV1.ListMeta `json:"metadata"`
	Items              []TServer `json:"items"`
}

type TEndpointSpec struct {
	App       string         `json:"app"`
	Server    string         `json:"server"`
	SubType   TServerSubType `json:"subType"`
	Important int32            `json:"important"`
	Taf       *TServerTaf    `json:"taf,omitempty"`
	Normal    *TServerNormal `json:"normal,omitempty"`
	HostPorts []TK8SHostPort `json:"hostPorts,omitempty"`
}

type TEndpointPodStatus struct {
	UID               string                      `json:"uid"`
	Name              string                      `json:"name"`
	PodIP             string                      `json:"podIP"`
	HostIP            string                      `json:"hostIP"`
	StartTime         k8sMetaV1.Time              `json:"startTime,omitempty"`
	ContainerStatuses []k8sCoreV1.ContainerStatus `json:"containerStatuses,omitempty"`
	SettingState      string                      `json:"settingState"`
	PresentState      string                      `json:"presentState"`
	PresentMessage    string                      `json:"presentMessage"`
	Tag               string                      `json:"tag"`
}

type TEndpointStatus struct {
	PodStatus []*TEndpointPodStatus `json:"pods"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TEndpoint struct {
	k8sMetaV1.TypeMeta   `json:",inline"`
	k8sMetaV1.ObjectMeta `json:"metadata,omitempty"`
	Spec                 TEndpointSpec   `json:"spec"`
	Status               TEndpointStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TEndpointList struct {
	k8sMetaV1.TypeMeta `json:",inline"`
	k8sMetaV1.ListMeta `json:"metadata"`
	Items              []TEndpoint `json:"items"`
}

type TTemplateSpec struct {
	Content string `json:"content"`
	Parent  string `json:"parent"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TTemplate struct {
	k8sMetaV1.TypeMeta   `json:",inline"`
	k8sMetaV1.ObjectMeta `json:"metadata,omitempty"`
	Spec                 TTemplateSpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TTemplateList struct {
	k8sMetaV1.TypeMeta `json:",inline"`
	k8sMetaV1.ListMeta `json:"metadata"`

	Items []TTemplate `json:"items"`
}

type TReleaseVersion struct {
	ServerType      string         `json:"serverType"`
	Image           string         `json:"image"`
	Tag             string         `json:"tag"`
	ImagePullSecret string         `json:"imagePullSecret"`
	CreatePerson    string         `json:"createPerson"`
	CreateTime      k8sMetaV1.Time `json:"createTime"`
}

type TReleaseSpec struct {
	List []*TReleaseVersion `json:"list"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TRelease struct {
	k8sMetaV1.TypeMeta   `json:",inline"`
	k8sMetaV1.ObjectMeta `json:"metadata,omitempty"`
	Spec                 TReleaseSpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TReleaseList struct {
	k8sMetaV1.TypeMeta `json:",inline"`
	k8sMetaV1.ListMeta `json:"metadata"`

	Items []TRelease `json:"items"`
}

type TTreeBusiness struct {
	Name         string         `json:"name"`
	Show         string         `json:"show"`
	Weight       int32            `json:"weight"`
	Mark         string         `json:"mark"`
	CreatePerson string         `json:"createPerson"`
	CreateTime   k8sMetaV1.Time `json:"createTime"`
}

type TTreeApps struct {
	Name         string         `json:"name"`
	BusinessRef  string         `json:"businessRef"`
	CreatePerson string         `json:"createPerson"`
	CreateTime   k8sMetaV1.Time `json:"createTime"`
	Mark         string         `json:"mark"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TTree struct {
	k8sMetaV1.TypeMeta   `json:",inline"`
	k8sMetaV1.ObjectMeta `json:"metadata,omitempty"`
	Businesses           []TTreeBusiness `json:"businesses"`
	Apps                 []TTreeApps     `json:"apps"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TTreeList struct {
	k8sMetaV1.TypeMeta `json:",inline"`
	k8sMetaV1.ListMeta `json:"metadata"`

	Items []TTree `json:"items"`
}

type TExitedPod struct {
	UID        string         `json:"uid"`
	Name       string         `json:"name"`
	Tag        string         `json:"tag"`
	NodeIP     string         `json:"nodeIP"`
	PodIP      string         `json:"podIP"`
	CreateTime k8sMetaV1.Time `json:"createTime"`
	DeleteTime k8sMetaV1.Time `json:"deleteTime"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TExitedRecord struct {
	k8sMetaV1.TypeMeta   `json:",inline"`
	k8sMetaV1.ObjectMeta `json:"metadata,omitempty"`
	App                  string       `json:"app"`
	Server               string       `json:"server"`
	Pods                 []TExitedPod `json:"pods"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TExitedRecordList struct {
	k8sMetaV1.TypeMeta `json:",inline"`
	k8sMetaV1.ListMeta `json:"metadata"`

	Items []TExitedRecord `json:"items"`
}

type TConfigApp struct {
	App           string         `json:"app"`
	ConfigName    string         `json:"configName"`
	ConfigContent string         `json:"configContent"`
	UpdateTime    k8sMetaV1.Time `json:"updateTime"`
	UpdatePerson  string         `json:"updatePerson"`
	UpdateReason  string         `json:"updateReason"`
}

type TConfigServer struct {
	App           string         `json:"app"`
	Server        string         `json:"server"`
	ConfigName    string         `json:"configName"`
	ConfigContent string         `json:"configContent"`
	PodSeq        *string        `json:"podSeq,omitempty"`
	UpdateTime    k8sMetaV1.Time `json:"updateTime"`
	UpdatePerson  string         `json:"updatePerson"`
	UpdateReason  string         `json:"updateReason"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TConfig struct {
	k8sMetaV1.TypeMeta   `json:",inline"`
	k8sMetaV1.ObjectMeta `json:"metadata,omitempty"`
	AppConfig            *TConfigApp    `json:"appConfig,omitempty"`
	ServerConfig         *TConfigServer `json:"serverConfig,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TConfigList struct {
	k8sMetaV1.TypeMeta `json:",inline"`
	k8sMetaV1.ListMeta `json:"metadata"`
	Items              []TConfig `json:"items"`
}

type TDeployApprove struct {
	Person string         `json:"person"`
	Time   k8sMetaV1.Time `json:"time"`
	Reason string         `json:"reason"`
	Result bool           `json:"result"`
}

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TDeploy struct {
	k8sMetaV1.TypeMeta   `json:",inline"`
	k8sMetaV1.ObjectMeta `json:"metadata,omitempty"`
	Apply                TServerSpec     `json:"apply"`
	Approve              *TDeployApprove `json:"approve,omitempty"`
	Deployed             *bool           `json:"deployed,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TDeployList struct {
	k8sMetaV1.TypeMeta `json:",inline"`
	k8sMetaV1.ListMeta `json:"metadata"`
	Items              []TDeploy `json:"items"`
}
