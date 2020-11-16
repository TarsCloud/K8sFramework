package client

import (
	"base"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sClient "k8s.io/client-go/kubernetes"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

var k8sNameSpace string
var k8sClientSet *k8sClient.Clientset
var k8sWatchImp base.K8SWatchInterface

const ServiceImagePlaceholder = " "

func buildAffinity(serverApp, serverName string, serverK8S *base.ServerK8S) *apiv1.Affinity {
	var affinity = &apiv1.Affinity{}

	switch serverK8S.NodeSelector.Kind {
	case base.NodeBind:
		{
			affinity.NodeAffinity = &apiv1.NodeAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: &apiv1.NodeSelector{
					NodeSelectorTerms: []apiv1.NodeSelectorTerm{
						{
							MatchExpressions: []apiv1.NodeSelectorRequirement{
								{
									Key:      base.TafNodeLabel,
									Operator: apiv1.NodeSelectorOpExists,
								},
								{
									Key:      "kubernetes.io/hostname",
									Operator: apiv1.NodeSelectorOpIn,
									Values:   serverK8S.NodeSelector.Value,
								},
							},
						},
					},
				},
			}
		}

	case base.AbilityPool:
		{
			affinity.NodeAffinity = &apiv1.NodeAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: &apiv1.NodeSelector{
					NodeSelectorTerms: []apiv1.NodeSelectorTerm{
						{
							MatchExpressions: []apiv1.NodeSelectorRequirement{
								{
									Key:      base.TafNodeLabel,
									Operator: apiv1.NodeSelectorOpExists,
								},
								{
									Key:      base.TafAbilityNodeLabelPrefix + serverApp,
									Operator: apiv1.NodeSelectorOpExists,
								},
							},
						},
					},
				},
			}
		}
	case base.PublicPool:
		{
			affinity.NodeAffinity = &apiv1.NodeAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: &apiv1.NodeSelector{
					NodeSelectorTerms: []apiv1.NodeSelectorTerm{
						{
							MatchExpressions: []apiv1.NodeSelectorRequirement{
								{
									Key:      base.TafNodeLabel,
									Operator: apiv1.NodeSelectorOpExists,
								},
								{
									Key:      base.TafPublicNodeLabel,
									Operator: apiv1.NodeSelectorOpExists,
								},
							},
						},
					},
				},
			}
		}
	}

	if serverK8S.NotStacked {
		affinity.PodAntiAffinity = &apiv1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []apiv1.PodAffinityTerm{
				{
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							base.TafServerAppLabel:  serverApp,
							base.TafServerNameLabel: serverName,
						},
					},
					Namespaces:  []string{k8sNameSpace},
					TopologyKey: "kubernetes.io/hostname",
				},
			},
		}
	}

	return affinity
}

func buildPodTemplate(serverApp string, serverName string, serverK8S *base.ServerK8S, serverServant base.ServerServant) apiv1.PodTemplateSpec {

	var k8sResourceName = strings.ToLower(serverApp + "-" + serverName)

	var hostPathCreateType = apiv1.HostPathDirectoryOrCreate
	var enableServiceLinks = false
	var FixedDNSConfigNDOTS = "2"

	var dnsPolicy apiv1.DNSPolicy

	if serverK8S.HostNetwork {
		dnsPolicy = apiv1.DNSClusterFirstWithHostNet
	} else {
		dnsPolicy = apiv1.DNSClusterFirst
	}

	if serverK8S.Image == "" {
		serverK8S.Image = ServiceImagePlaceholder
	}

	spec := apiv1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name: k8sResourceName,
			Labels: map[string]string{
				base.TafServerAppLabel:     serverApp,
				base.TafServerNameLabel:    serverName,
				base.TafServerVersionLabel: serverK8S.Version,
			},
		},
		Spec: apiv1.PodSpec{
			Volumes: []apiv1.Volume{
				{
					Name: "host-log-path",
					VolumeSource: apiv1.VolumeSource{
						HostPath: &apiv1.HostPathVolumeSource{
							Path: "/usr/local/app/taf/app_log",
							Type: &hostPathCreateType,
						},
					},
				},
				{
					Name: "host-data-path",
					VolumeSource: apiv1.VolumeSource{
						HostPath: &apiv1.HostPathVolumeSource{
							Path: "/usr/local/app/taf/data",
							Type: &hostPathCreateType,
						},
					},
				},
			},
			Containers: []apiv1.Container{
				{
					Name:  k8sResourceName,
					Image: serverK8S.Image,
					Env: []apiv1.EnvVar{
						{
							Name: "ServerApp",
							ValueFrom: &apiv1.EnvVarSource{
								FieldRef: &apiv1.ObjectFieldSelector{
									FieldPath: fmt.Sprintf("metadata.labels['%s']", base.TafServerAppLabel),
								},
							},
						},
						{
							Name: "PodName",
							ValueFrom: &apiv1.EnvVarSource{
								FieldRef: &apiv1.ObjectFieldSelector{
									FieldPath: "metadata.name",
								},
							},
						},
					},
					Resources: apiv1.ResourceRequirements{},
					VolumeMounts: []apiv1.VolumeMount{
						{
							Name:      "host-log-path",
							MountPath: "/host-log-path",
						},
						{
							Name:      "host-data-path",
							MountPath: "/host-data-path",
						},
					},
					ImagePullPolicy: apiv1.PullAlways,
				},
			},
			DNSPolicy:   dnsPolicy,
			HostNetwork: serverK8S.HostNetwork,
			HostIPC:     serverK8S.HostIpc,
			ImagePullSecrets: []apiv1.LocalObjectReference{
				{
					Name: "taf-image-secret",
				},
			},
			Affinity: buildAffinity(serverApp, serverName, serverK8S),
			DNSConfig: &apiv1.PodDNSConfig{
				Options: []apiv1.PodDNSConfigOption{
					{
						Name:  "ndots",
						Value: &FixedDNSConfigNDOTS,
					},
				},
			},
			ReadinessGates: []apiv1.PodReadinessGate{
				{
					ConditionType: "taf.io/service",
				},
			},
			EnableServiceLinks: &enableServiceLinks,
		},
	}

	if !serverK8S.HostNetwork && serverK8S.HostPort != nil {
		for hostPortName, hostPortValue := range serverK8S.HostPort {

			if hostPortValue == 0 {
				delete(serverK8S.HostPort, hostPortName)
				continue
			}

			lowerCaseHostPortName := strings.ToLower(hostPortName)
			serverServantElem, ok := serverServant[lowerCaseHostPortName]

			if !ok {
				delete(serverK8S.HostPort, hostPortName)
				continue
			}
			containersPort := apiv1.ContainerPort{
				Name:          lowerCaseHostPortName,
				HostPort:      hostPortValue,
				ContainerPort: int32(serverServantElem.Port),
			}
			if serverServantElem.IsTcp {
				containersPort.Protocol = apiv1.ProtocolTCP
			} else {
				containersPort.Protocol = apiv1.ProtocolUDP
			}
			if spec.Spec.Containers[0].Ports == nil {
				spec.Spec.Containers[0].Ports = make([]apiv1.ContainerPort, 0, 1)
			}
			spec.Spec.Containers[0].Ports = append(spec.Spec.Containers[0].Ports, containersPort)
		}
	}
	return spec
}

func buildStatefulSet(serverApp string, serverName string, serverServant base.ServerServant, serverK8S *base.ServerK8S) *appsv1.StatefulSet {
	var k8sResourceName = strings.ToLower(serverApp + "-" + serverName)
	nodeSelectorBytes, _ := json.Marshal(serverK8S.NodeSelector)

	var statefulSet = &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      k8sResourceName,
			Namespace: k8sNameSpace,
			Annotations: map[string]string{
				base.TafNodeSelectorLabel: string(nodeSelectorBytes),
				base.TafNotStackedLabel:   strconv.FormatBool(serverK8S.NotStacked),
			},
			Labels: map[string]string{
				base.TafServerAppLabel:  serverApp,
				base.TafServerNameLabel: serverName,
			},
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName:         k8sResourceName,
			PodManagementPolicy: appsv1.ParallelPodManagement,
			Replicas:            &serverK8S.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					base.TafServerAppLabel:  serverApp,
					base.TafServerNameLabel: serverName,
				},
			},
			Template: buildPodTemplate(serverApp, serverName, serverK8S, serverServant),
		},
	}
	return statefulSet
}

func updateStatefulSet(statefulSet *appsv1.StatefulSet, serverApp, serverName string, serverServant base.ServerServant, serverK8S *base.ServerK8S) {
	nodeSelectorBytes, _ := json.Marshal(serverK8S.NodeSelector)
	statefulSet.Annotations[base.TafNodeSelectorLabel] = string(nodeSelectorBytes)
	statefulSet.Annotations[base.TafNotStackedLabel] = strconv.FormatBool(serverK8S.NotStacked)
	statefulSet.Spec.Replicas = &serverK8S.Replicas
	statefulSet.Spec.Template = buildPodTemplate(serverApp, serverName, serverK8S, serverServant)
}

func buildService(serverApp string, serverName string, serverServant base.ServerServant) *apiv1.Service {

	servantByte, _ := json.Marshal(serverServant)
	k8sService := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.ToLower(serverApp + "-" + serverName),
			Namespace: k8sNameSpace,
			Labels: map[string]string{
				base.TafServerAppLabel:  serverApp,
				base.TafServerNameLabel: serverName,
			},
			Annotations: map[string]string{
				base.TafServantLabel: string(servantByte),
			},
		},
	}
	k8sService.Spec.Selector = map[string]string{
		base.TafServerAppLabel:  serverApp,
		base.TafServerNameLabel: serverName,
	}
	for _, v := range serverServant {
		var port apiv1.ServicePort
		port.Name = strings.ToLower(v.Name)
		port.Port = int32(v.Port)
		if v.IsTcp {
			port.Protocol = apiv1.ProtocolTCP
		} else {
			port.Protocol = apiv1.ProtocolUDP
		}
		k8sService.Spec.Ports = append(k8sService.Spec.Ports, port)
	}

	// tafnode 的对外管理端口 ,固定值
	k8sService.Spec.Ports = append(k8sService.Spec.Ports, apiv1.ServicePort{
		Protocol: apiv1.ProtocolTCP,
		Port:     19385,
		Name:     "nodeobj",
	})

	k8sService.Spec.Type = apiv1.ServiceTypeClusterIP
	k8sService.Spec.ClusterIP = apiv1.ClusterIPNone

	return k8sService
}

func updateService(service *apiv1.Service, serverServant base.ServerServant) {
	servantByte, _ := json.Marshal(serverServant)
	service.Spec.Ports = make([]apiv1.ServicePort, 0, len(serverServant)+1)
	service.Annotations[base.TafServantLabel] = string(servantByte)
	for _, v := range serverServant {
		var port apiv1.ServicePort
		port.Name = strings.ToLower(v.Name)
		port.Port = int32(v.Port)
		if v.IsTcp {
			port.Protocol = apiv1.ProtocolTCP
		} else {
			port.Protocol = apiv1.ProtocolUDP
		}
		service.Spec.Ports = append(service.Spec.Ports, port)
	}
	service.Spec.Ports = append(service.Spec.Ports, apiv1.ServicePort{
		Protocol: apiv1.ProtocolTCP,
		Port:     19385,
		Name:     "nodeobj",
	})
}

func createServer(serverApp string, serverName string, serverServant base.ServerServant, serverK8S *base.ServerK8S) error {

	allOk := false

	k8sService := buildService(serverApp, serverName, serverServant)
	serviceInterface := k8sClientSet.CoreV1().Services(k8sNameSpace)
	if _, err := serviceInterface.Create(k8sService); err != nil {
		return err
	}
	defer func() {
		if !allOk {
			_ = serviceInterface.Delete(k8sService.Name, nil)
		}
	}()

	statefulSet := buildStatefulSet(serverApp, serverName, serverServant, serverK8S)
	statefulSetInterface := k8sClientSet.AppsV1().StatefulSets(k8sNameSpace)
	if _, err := statefulSetInterface.Create(statefulSet); err != nil {
		return err
	}
	defer func() {
		if !allOk {
			_ = statefulSetInterface.Delete(statefulSet.Name, nil)
		}
	}()

	allOk = true
	return nil
}

func updateServerK8S(serverApp string, serverName string, params map[base.UpdateK8SKey]interface{}) error {

	currentK8S := k8sWatchImp.GetServerK8S(serverApp, serverName)
	if currentK8S == nil {
		return errors.New("内部错误")
	}

	anyParamChange := false
	for targetKey, targetInterface := range params {
		switch targetKey {
		case base.Replicas, base.Image:

			targetReplicasInterface, replicasOk := params[base.Replicas]

			targetImageInterface, imageOk := params[base.Image]

			if replicasOk {
				targetReplicas := targetReplicasInterface.(int32)
				if targetReplicas == currentK8S.Replicas {
					delete(params, base.Replicas)
					break
				}

				if !imageOk {
					if currentK8S.Image == ServiceImagePlaceholder {
						return errors.New("未发布版本时不允许设置 Replicas 参数")
					}
					currentK8S.Replicas = targetReplicas
					anyParamChange = true
					break
				}

				targetImageValue := targetImageInterface.(string)
				if targetImageValue == "" && targetReplicas != 0 {
					return errors.New(" 不允许设置Image参数为\"\"")
				}

				currentK8S.Replicas = targetReplicas
				currentK8S.Image = targetImageValue
				anyParamChange = true
				delete(params, base.Replicas)
				delete(params, base.Image)
				break
			}

			targetImageValue := targetImageInterface.(string)
			if targetImageValue == "" && currentK8S.Replicas != 0 {
				return errors.New(" 不允许设置Image参数为\"\"")
			}
			currentK8S.Image = targetImageValue
			anyParamChange = true

		case base.NotStacked:
			targetValue := targetInterface.(bool)
			if targetValue != currentK8S.NotStacked {
				currentK8S.NotStacked = targetValue
				anyParamChange = true
			}
		case base.HostIpc:
			targetValue := targetInterface.(bool)
			if targetValue != currentK8S.HostIpc {
				currentK8S.HostIpc = targetValue
				anyParamChange = true
			}
		case base.HostNetwork:
			targetValue := targetInterface.(bool)
			if targetValue != currentK8S.HostNetwork {
				currentK8S.HostNetwork = targetValue
				anyParamChange = true
			}
		case base.Version:
			targetValue := targetInterface.(string)
			if targetValue != currentK8S.Version {
				currentK8S.Version = targetValue
				anyParamChange = true
			}
		case base.NodeSelect:
			targetValue := targetInterface.(base.NodeSelector)
			if !reflect.DeepEqual(targetValue, currentK8S.NodeSelector) {
				currentK8S.NodeSelector = targetValue
				anyParamChange = true
			}
		case base.HostPort:
			targetValue := targetInterface.(map[string]int32)
			if !reflect.DeepEqual(targetValue, currentK8S.HostPort) {
				currentK8S.HostPort = targetValue
				anyParamChange = true
			}
		}
	}

	if !anyParamChange {
		return nil
	}

	serverK8SCheckFun, _ := govalidator.CustomTypeTagMap.Get("matches-ServerK8S")
	if !serverK8SCheckFun(*currentK8S, nil) {
		return errors.New("内部错误")
	}

	currentServant := k8sWatchImp.GetServerServant(serverApp, serverName)
	if currentServant == nil {
		return errors.New("内部错误")
	}

	k8sResourceName := strings.ToLower(serverApp + "-" + serverName)

	statefulSetInterface := k8sClientSet.AppsV1().StatefulSets(k8sNameSpace)
	currentStatefulSet, _ := statefulSetInterface.Get(k8sResourceName, metav1.GetOptions{})
	updateStatefulSet(currentStatefulSet, serverApp, serverName, currentServant, currentK8S)
	if _, err := statefulSetInterface.Update(currentStatefulSet); err != nil {
		return err
	}

	return nil
}

func deleteServer(serverApp string, serverName string) {
	k8sResourceName := strings.ToLower(serverApp + "-" + serverName)
	_ = k8sClientSet.CoreV1().Services(k8sNameSpace).Delete(k8sResourceName, nil)
	_ = k8sClientSet.AppsV1().StatefulSets(k8sNameSpace).Delete(k8sResourceName, nil)
}

func deleteNodeAbility(node string, apps ...string) error {

	if !k8sWatchImp.IsClusterHadNode(node) {
		return errors.New(fmt.Sprintf("%s不属于 Taf 管理的节点", node))
	}

	nodeInterface := k8sClientSet.CoreV1().Nodes()

	var err error
	var k8sNode *apiv1.Node

	if k8sNode, err = nodeInterface.Get(node, metav1.GetOptions{}); err != nil {
		return errors.New("内部错误")
	}

	deletedAnyLabel := false

	for _, v := range apps {
		abilityLabel := base.TafAbilityNodeLabelPrefix + v
		if _, ok := k8sNode.Labels[abilityLabel]; ok {
			deletedAnyLabel = true
			delete(k8sNode.Labels, abilityLabel)
		}
	}

	if deletedAnyLabel {
		_, updateErr := nodeInterface.Update(k8sNode)
		return updateErr
	}

	return nil
}

func addNodeAbility(node string, apps ...string) error {

	if !k8sWatchImp.IsClusterHadNode(node) {
		return errors.New(fmt.Sprintf("%s不属于 Taf 管理的节点", node))
	}

	nodeInterface := k8sClientSet.CoreV1().Nodes()

	var err error
	var k8sNode *apiv1.Node

	if k8sNode, err = nodeInterface.Get(node, metav1.GetOptions{}); err != nil {
		return errors.New("内部错误")
	}

	addAnyLabel := false

	for _, v := range apps {
		abilityLabel := base.TafAbilityNodeLabelPrefix + v
		if _, ok := k8sNode.Labels[abilityLabel]; !ok {
			addAnyLabel = true
			k8sNode.Labels[abilityLabel] = ""
		}
	}

	if addAnyLabel {
		_, updateErr := nodeInterface.Update(k8sNode)
		return updateErr
	}

	return nil
}

func setNodePublic(nodes ...string) error {
	nodeInterface := k8sClientSet.CoreV1().Nodes()

	var err error
	var k8sNode *apiv1.Node

	for _, node := range nodes {
		if !k8sWatchImp.IsClusterHadNode(node) {
			continue
		}

		if k8sNode, err = nodeInterface.Get(node, metav1.GetOptions{}); err != nil {
			return errors.New("内部错误")
		}

		if _, ok := k8sNode.Labels[base.TafPublicNodeLabel]; !ok {
			k8sNode.Labels[base.TafPublicNodeLabel] = ""
			_, updateErr := nodeInterface.Update(k8sNode)
			return updateErr
		}
	}
	return nil
}

func deleteNodePublic(nodes ...string) error {
	nodeInterface := k8sClientSet.CoreV1().Nodes()

	var err error
	var k8sNode *apiv1.Node

	for _, node := range nodes {
		if !k8sWatchImp.IsClusterHadNode(node) {
			continue
		}

		if k8sNode, err = nodeInterface.Get(node, metav1.GetOptions{}); err != nil {
			return errors.New("内部错误")
		}

		if _, ok := k8sNode.Labels[base.TafPublicNodeLabel]; ok {
			delete(k8sNode.Labels, base.TafPublicNodeLabel)
			_, updateErr := nodeInterface.Update(k8sNode)
			return updateErr
		}
	}
	return nil
}

func appendServant(serverApp string, serverName string, serverServant base.ServerServant) error {
	if serverServant == nil || len(serverServant) == 0 {
		return nil
	}

	currentServant := k8sWatchImp.GetServerServant(serverApp, serverName)
	if currentServant == nil {
		return errors.New("内部错误")
	}

	for k, v := range serverServant {
		if _, ok := currentServant[k]; ok {
			return errors.New(fmt.Sprintf("已存在同名servant:%s", k))
		}
		currentServant[k] = v
	}

	checkServantFun, _ := govalidator.CustomTypeTagMap.Get("matches-ServerServant")
	if !checkServantFun(currentServant, nil) {
		return errors.New("错误参数")
	}

	k8sResourceName := strings.ToLower(serverApp + "-" + serverName)
	serviceInterface := k8sClientSet.CoreV1().Services(k8sNameSpace)
	currentService, _ := serviceInterface.Get(k8sResourceName, metav1.GetOptions{})
	updateService(currentService, currentServant)

	if _, err := serviceInterface.Update(currentService); err != nil {
		return err
	}

	return nil
}

func eraseServant(serverApp string, serverName string, servantElemName string) error {

	currentServant := k8sWatchImp.GetServerServant(serverApp, serverName)
	if currentServant == nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, fmt.Sprintf("GetServerServant(%s,%s return nil)", serverApp, serverName)))
		return errors.New("内部错误")
	}

	lowerServantElemName := strings.ToLower(servantElemName)

	_, ok := currentServant[lowerServantElemName]
	if !ok {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, "servantElemName does't exist"))
		return nil
	}
	if len(currentServant) == 1 {
		return errors.New("不能删除唯一的 servant ")
	}

	delete(currentServant, lowerServantElemName)

	k8sResourceName := strings.ToLower(serverApp + "-" + serverName)
	serviceInterface := k8sClientSet.CoreV1().Services(k8sNameSpace)
	currentService, _ := serviceInterface.Get(k8sResourceName, metav1.GetOptions{})
	updateService(currentService, currentServant)
	if _, err := serviceInterface.Update(currentService); err != nil {
		return err
	}

	shouldUpdateStatefulSet := false

	serverK8S := k8sWatchImp.GetServerK8S(serverApp, serverName)

	if serverK8S == nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, fmt.Sprintf("GetServerK8S(%s,%s return nil)", serverApp, serverName)))
		return errors.New("内部错误")
	}

	if !serverK8S.HostNetwork {
		if hostPortValue, ok := serverK8S.HostPort[lowerServantElemName]; ok && hostPortValue != 0 {
			shouldUpdateStatefulSet = true
			delete(serverK8S.HostPort, lowerServantElemName)
		}
	}

	if shouldUpdateStatefulSet {
		statefulSetInterface := k8sClientSet.AppsV1().StatefulSets(k8sNameSpace)
		currentStatefulSet, _ := statefulSetInterface.Get(k8sResourceName, metav1.GetOptions{})
		updateStatefulSet(currentStatefulSet, serverApp, serverName, currentServant, serverK8S)
		if _, err := statefulSetInterface.Update(currentStatefulSet); err != nil {
			return err
		}
	}

	return nil
}

func updateServant(serverApp string, serverName string, servantElemName string, params map[base.UpdateServantKey]interface{}) error {

	lowerServantElemName := strings.ToLower(servantElemName)

	currentServant := k8sWatchImp.GetServerServant(serverApp, serverName)
	if currentServant == nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, fmt.Sprintf("GetServerServant(%s,%s return nil)", serverApp, serverName)))
		return errors.New("内部错误")
	}

	servantElem, ok := currentServant[lowerServantElemName]
	if !ok {
		_, file, line, _ := runtime.Caller(0)
		fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, "servantElemName does't exist"))
		return nil
	}

	anyParamsChanged := false
	shouldCheckStatefulSet := false

	for k, v := range params {
		switch k {
		case base.ServantName:
			vv := v.(string)
			if strings.ToLower(vv) != strings.ToLower(servantElem.Name) {
				servantElem.Name = vv
				delete(currentServant, lowerServantElemName)
				currentServant[strings.ToLower(vv)] = servantElem
				shouldCheckStatefulSet = true
				anyParamsChanged = true
			}
		case base.ServantPort:
			vv := v.(int)
			if vv != servantElem.Port {
				servantElem.Port = vv
				shouldCheckStatefulSet = true
				anyParamsChanged = true
			}
		case base.ServantIsTcp:
			vv := v.(bool)
			if vv != servantElem.IsTcp {
				servantElem.IsTcp = vv
				shouldCheckStatefulSet = true
				anyParamsChanged = true
			}
		case base.ServantCapacity:
			vv := v.(int)
			if vv != servantElem.Capacity {
				servantElem.Capacity = vv
				anyParamsChanged = true
			}
		case base.ServantConnections:
			vv := v.(int)
			if vv != servantElem.Connections {
				servantElem.Connections = vv
				anyParamsChanged = true
			}
		case base.ServantIsTaf:
			vv := v.(bool)
			if vv != servantElem.IsTaf {
				servantElem.IsTaf = vv
				anyParamsChanged = true
			}
		case base.ServantThreads:
			vv := v.(int)
			if vv != servantElem.Threads {
				servantElem.Threads = vv
				anyParamsChanged = true
			}
		case base.ServantTimeout:
			vv := v.(int)
			if vv != servantElem.Timeout {
				servantElem.Timeout = vv
				anyParamsChanged = true
			}
		}
	}

	if !anyParamsChanged {
		return nil
	}

	checkServantFun, _ := govalidator.CustomTypeTagMap.Get("matches-ServerServant")
	if !checkServantFun(currentServant, nil) {
		return errors.New("错误参数")
	}

	k8sResourceName := strings.ToLower(serverApp + "-" + serverName)

	serviceInterface := k8sClientSet.CoreV1().Services(k8sNameSpace)
	currentService, _ := serviceInterface.Get(k8sResourceName, metav1.GetOptions{})
	updateService(currentService, currentServant)
	if _, err := serviceInterface.Update(currentService); err != nil {
		return err
	}

	if shouldCheckStatefulSet {
		for {
			currentK8S := k8sWatchImp.GetServerK8S(serverApp, serverName)
			if currentK8S == nil {
				_, file, line, _ := runtime.Caller(0)
				fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, fmt.Sprintf("GetServerK8S(%s,%s return nil)", serverApp, serverName)))
				return errors.New("内部错误")
			}

			if len(currentK8S.HostPort) == 0 {
				break
			}

			currentHostPort, ok := currentK8S.HostPort[lowerServantElemName]
			if !ok {
				break
			}

			//防止出现 update servantElemName 的情况
			delete(currentK8S.HostPort, lowerServantElemName)
			currentK8S.HostPort[strings.ToLower(servantElem.Name)] = currentHostPort

			statefulSetInterface := k8sClientSet.AppsV1().StatefulSets(k8sNameSpace)
			currentStatefulSet, _ := statefulSetInterface.Get(k8sResourceName, metav1.GetOptions{})

			updateStatefulSet(currentStatefulSet, serverApp, serverName, currentServant, currentK8S)
			if _, err := statefulSetInterface.Update(currentStatefulSet); err != nil {
				return err
			}

			break
		}
	}
	return nil
}

type K8SClientImp struct {
}

func (k K8SClientImp) AppendServant(serverApp string, serverName string, serverServant base.ServerServant) error {
	return appendServant(serverApp, serverName, serverServant)
}

func (k K8SClientImp) EraseServant(serverApp string, serverName string, adapterName string) error {
	return eraseServant(serverApp, serverName, adapterName)
}

func (k K8SClientImp) UpdateServant(serverApp string, serverName string, adapterName string, params map[base.UpdateServantKey]interface{}) error {
	return updateServant(serverApp, serverName, adapterName, params)
}

func (k K8SClientImp) SetK8SClient(clientSet *k8sClient.Clientset) {
	k8sClientSet = clientSet
}

func (k K8SClientImp) SetWorkNamespace(namespace string) {
	k8sNameSpace = namespace
}

func (k K8SClientImp) SetK8SWatchImp(watchImp base.K8SWatchInterface) {
	k8sWatchImp = watchImp
}

func (k K8SClientImp) DeleteNodeAbility(node string, apps ...string) error {
	return deleteNodeAbility(node, apps...)
}

func (k K8SClientImp) AddNodeAbility(node string, apps ...string) error {
	return addNodeAbility(node, apps...)
}

func (k K8SClientImp) SetPublicNode(nodes ...string) error {
	return setNodePublic(nodes...)
}

func (k K8SClientImp) DeletePublicNode(nodes ...string) error {
	return deleteNodePublic(nodes...)
}

func (k K8SClientImp) CreateServer(serverApp string, serverName string, serverServant base.ServerServant, serverK8S *base.ServerK8S) error {
	// 新建 Server时,将部分参数重设为默认值
	serverK8S.Replicas = 0
	serverK8S.Image = ""
	serverK8S.Version = ""
	if err := createServer(serverApp, serverName, serverServant, serverK8S); err != nil {
		return err
	}
	return nil
}

func (k K8SClientImp) UpdateServerK8S(serverApp string, serverName string, params map[base.UpdateK8SKey]interface{}) error {
	return updateServerK8S(serverApp, serverName, params)
}

func (k K8SClientImp) DeleteServer(serverApp string, serverName string) {
	deleteServer(serverApp, serverName)
}
