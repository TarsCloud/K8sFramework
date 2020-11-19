package main

import (
	"fmt"
	k8sAppsV1 "k8s.io/api/apps/v1"
	k8sCoreV1 "k8s.io/api/core/v1"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	crdV1Alpha1 "k8s.tars.io/crd/v1alpha1"
	"strings"
)

func buildPodAffinity(server *crdV1Alpha1.TServer) *k8sCoreV1.Affinity {
	affinity := &k8sCoreV1.Affinity{}
	if server.Spec.K8S.NodeSelector.NodeBind != nil {
		affinity.NodeAffinity = &k8sCoreV1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &k8sCoreV1.NodeSelector{
				NodeSelectorTerms: []k8sCoreV1.NodeSelectorTerm{
					{
						MatchExpressions: []k8sCoreV1.NodeSelectorRequirement{
							{
								Key:      TarsNodeLabelPrefix + server.Namespace,
								Operator: k8sCoreV1.NodeSelectorOpExists,
							},
							{
								Key:      TarsAbilityNodeLabelPrefix + server.Spec.App,
								Operator: k8sCoreV1.NodeSelectorOpExists,
							},
						},
					},
				},
			},
		}
	} else if server.Spec.K8S.NodeSelector.AbilityPool != nil {
		affinity.NodeAffinity = &k8sCoreV1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &k8sCoreV1.NodeSelector{
				NodeSelectorTerms: []k8sCoreV1.NodeSelectorTerm{
					{
						MatchExpressions: []k8sCoreV1.NodeSelectorRequirement{
							{
								Key:      TarsNodeLabelPrefix + server.Namespace,
								Operator: k8sCoreV1.NodeSelectorOpExists,
							},
							{
								Key:      TarsAbilityNodeLabelPrefix + server.Spec.App,
								Operator: k8sCoreV1.NodeSelectorOpExists,
							},
						},
					},
				},
			},
		}
	} else if server.Spec.K8S.NodeSelector.PublicPool != nil {
		affinity.NodeAffinity = &k8sCoreV1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &k8sCoreV1.NodeSelector{
				NodeSelectorTerms: []k8sCoreV1.NodeSelectorTerm{
					{
						MatchExpressions: []k8sCoreV1.NodeSelectorRequirement{
							{
								Key:      TarsNodeLabelPrefix + server.Namespace,
								Operator: k8sCoreV1.NodeSelectorOpExists,
							},
							{
								Key:      TarsPublicNodeLabel,
								Operator: k8sCoreV1.NodeSelectorOpExists,
							},
						},
					},
				},
			},
		}
	} else if server.Spec.K8S.NodeSelector.DaemonSet != nil {
		affinity.NodeAffinity = &k8sCoreV1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &k8sCoreV1.NodeSelector{
				NodeSelectorTerms: []k8sCoreV1.NodeSelectorTerm{
					{
						MatchExpressions: []k8sCoreV1.NodeSelectorRequirement{
							{
								Key:      TarsNodeLabelPrefix + server.Namespace,
								Operator: k8sCoreV1.NodeSelectorOpExists,
							},
						},
					},
				},
			},
		}
	}

	if server.Spec.K8S.NotStacked != nil && *server.Spec.K8S.NotStacked {
		affinity.PodAntiAffinity = &k8sCoreV1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []k8sCoreV1.PodAffinityTerm{
				{
					LabelSelector: &k8sMetaV1.LabelSelector{
						MatchLabels: map[string]string{
							TServerAppLabel:  server.Spec.App,
							TServerNameLabel: server.Spec.Server,
						},
					},
					Namespaces:  []string{server.Namespace},
					TopologyKey: "kubernetes.io/hostname",
				},
			},
		}
	}

	return affinity
}

func buildImagePullSecrets(server *crdV1Alpha1.TServer) []k8sCoreV1.LocalObjectReference {

	if server.Spec.Release == nil || server.Spec.Release.ImagePullSecret == "" {
		return nil
	}

	return []k8sCoreV1.LocalObjectReference{
		{
			Name: server.Spec.Release.ImagePullSecret,
		},
	}
}

func buildContainerPort(server *crdV1Alpha1.TServer) []k8sCoreV1.ContainerPort {

	if server.Spec.K8S.HostPorts == nil {
		return nil
	}

	var containerPorts []k8sCoreV1.ContainerPort

	hostPorts := server.Spec.K8S.HostPorts
	if server.Spec.Tars != nil {
		servants := server.Spec.Tars.Servants
		containerPorts = make([]k8sCoreV1.ContainerPort, 0, len(servants))
		for i := range hostPorts {
			for j := range servants {
				if servants[j].Name == hostPorts[i].NameRef {
					containersPort := k8sCoreV1.ContainerPort{
						Name:          fmt.Sprintf("p%d-%d", hostPorts[i].Port, servants[i].Port),
						HostPort:      hostPorts[i].Port,
						ContainerPort: servants[i].Port,
						HostIP:        "",
					}
					if servants[j].IsTcp {
						containersPort.Protocol = k8sCoreV1.ProtocolTCP
					} else {
						containersPort.Protocol = k8sCoreV1.ProtocolUDP
					}
					containerPorts = append(containerPorts, containersPort)
				}
			}
		}
	} else if server.Spec.Normal.Ports != nil {
		ports := server.Spec.Normal.Ports
		for i := range hostPorts {
			for j := range ports {
				if ports[j].Name == hostPorts[i].NameRef {
					containersPort := k8sCoreV1.ContainerPort{
						Name:          fmt.Sprintf("%d-%d", hostPorts[i].Port, ports[i].Port),
						HostPort:      hostPorts[i].Port,
						ContainerPort: ports[i].Port,
						HostIP:        "",
					}
					if ports[j].IsTcp {
						containersPort.Protocol = k8sCoreV1.ProtocolTCP
					} else {
						containersPort.Protocol = k8sCoreV1.ProtocolUDP
					}
					containerPorts = append(containerPorts, containersPort)
				}
			}
		}
	}
	return containerPorts
}

func buildReadinessGates(server *crdV1Alpha1.TServer) []k8sCoreV1.PodReadinessGate {
	if server.Spec.K8S.ReadinessGate == nil || *server.Spec.K8S.ReadinessGate == "" {
		return nil
	}
	return []k8sCoreV1.PodReadinessGate{
		{
			ConditionType: TPodReadinessGate,
		},
	}
}

func buildMounts(Server *crdV1Alpha1.TServer) ([]k8sCoreV1.Volume, []k8sCoreV1.VolumeMount) {
	if Server.Spec.K8S.Mounts == nil || len(Server.Spec.K8S.Mounts) == 0 {
		return nil, nil
	}
	mounts := Server.Spec.K8S.Mounts
	volumes := make([]k8sCoreV1.Volume, 0, len(mounts))
	volumeMounts := make([]k8sCoreV1.VolumeMount, 0, len(mounts))
	for i := range mounts {
		volume := k8sCoreV1.Volume{
			Name:         mounts[i].Name,
			VolumeSource: mounts[i].Source,
		}
		volumes = append(volumes, volume)
		volumeMount := k8sCoreV1.VolumeMount{
			Name:             mounts[i].Name,
			ReadOnly:         mounts[i].ReadOnly,
			MountPath:        mounts[i].MountPath,
			SubPath:          mounts[i].SubPath,
			MountPropagation: (*k8sCoreV1.MountPropagationMode)(mounts[i].MountPropagation),
			SubPathExpr:      mounts[i].SubPathExpr,
		}
		volumeMounts = append(volumeMounts, volumeMount)
	}
	return volumes, volumeMounts
}

func buildPodTemplate(server *crdV1Alpha1.TServer) k8sCoreV1.PodTemplateSpec {

	var enableServiceLinks = false
	var FixedDNSConfigNDOTS = "2"

	var dnsPolicy = k8sCoreV1.DNSClusterFirst

	serverImage := ServiceImagePlaceholder
	serverTag := ""

	if server.Spec.Release != nil {
		serverImage = server.Spec.Release.Image
		serverTag = server.Spec.Release.Tag
	}

	volumes, volumesMounts := buildMounts(server)
	spec := k8sCoreV1.PodTemplateSpec{
		ObjectMeta: k8sMetaV1.ObjectMeta{
			Name: server.Name,
			Labels: map[string]string{
				TServerAppLabel:  server.Spec.App,
				TServerNameLabel: server.Spec.Server,
				TServerTagLabel:  serverTag,
			},
		},
		Spec: k8sCoreV1.PodSpec{
			Volumes:        volumes,
			InitContainers: nil,
			Containers: []k8sCoreV1.Container{
				{
					Image:           serverImage,
					Name:            server.Name,
					EnvFrom:         server.Spec.K8S.EnvFrom,
					Env:             server.Spec.K8S.Env,
					Resources:       k8sCoreV1.ResourceRequirements{},
					VolumeMounts:    volumesMounts,
					ImagePullPolicy: k8sCoreV1.PullAlways,
					Ports:           buildContainerPort(server),
				},
			},
			DNSPolicy:          dnsPolicy,
			ServiceAccountName: server.Spec.K8S.ServiceAccount,
			HostIPC:            server.Spec.K8S.HostIPC,
			HostNetwork:        server.Spec.K8S.HostNetwork,
			ImagePullSecrets:   buildImagePullSecrets(server),
			Affinity:           buildPodAffinity(server),
			DNSConfig: &k8sCoreV1.PodDNSConfig{
				Options: []k8sCoreV1.PodDNSConfigOption{
					{
						Name:  "ndots",
						Value: &FixedDNSConfigNDOTS,
					},
				},
			},
			ReadinessGates:     buildReadinessGates(server),
			EnableServiceLinks: &enableServiceLinks,
		},
	}
	return spec
}

func buildStatefulSet(server *crdV1Alpha1.TServer) *k8sAppsV1.StatefulSet {
	var statefulSet = &k8sAppsV1.StatefulSet{
		ObjectMeta: k8sMetaV1.ObjectMeta{
			Name:      server.Name,
			Namespace: server.Namespace,
			Labels: map[string]string{
				TServerAppLabel:  server.Spec.App,
				TServerNameLabel: server.Spec.Server,
			},
			OwnerReferences: []k8sMetaV1.OwnerReference{
				{
					APIVersion: TServerAPIVersion,
					Kind:       TServerKind,
					Name:       server.Name,
					UID:        server.UID,
				},
			},
		},
		Spec: k8sAppsV1.StatefulSetSpec{
			ServiceName: server.Name,
			Replicas:    &server.Spec.K8S.Replicas,
			Selector: &k8sMetaV1.LabelSelector{
				MatchLabels: map[string]string{
					TServerAppLabel:  server.Spec.App,
					TServerNameLabel: server.Spec.Server,
				},
			},
			Template: buildPodTemplate(server),
		},
	}
	statefulSet.Spec.PodManagementPolicy = k8sAppsV1.PodManagementPolicyType(server.Spec.K8S.PodManagementPolicy)
	return statefulSet
}

func buildServicePortsFromTServant(server *crdV1Alpha1.TServer) []k8sCoreV1.ServicePort {

	serverServant := server.Spec.Tars.Servants
	ports := make([]k8sCoreV1.ServicePort, 0, len(serverServant)+1)
	for _, v := range serverServant {
		var port k8sCoreV1.ServicePort
		port.Name = strings.ToLower(v.Name)
		port.Port = v.Port
		port.TargetPort = intstr.FromInt(int(v.Port))
		if v.IsTcp {
			port.Protocol = k8sCoreV1.ProtocolTCP
		} else {
			port.Protocol = k8sCoreV1.ProtocolUDP
		}
		ports = append(ports, port)
	}

	if !server.Spec.Tars.Foreground {
		ports = append(ports, k8sCoreV1.ServicePort{
			Protocol:   k8sCoreV1.ProtocolTCP,
			Port:       NodeServantPort,
			Name:       NodeServantName,
			TargetPort: intstr.FromInt(NodeServantPort),
		})
	}
	return ports
}

func buildServicePortsFromNPorts(server *crdV1Alpha1.TServer) []k8sCoreV1.ServicePort {
	nPorts := server.Spec.Normal.Ports
	ports := make([]k8sCoreV1.ServicePort, 0, len(nPorts))
	for _, v := range nPorts {
		var port k8sCoreV1.ServicePort
		port.Name = strings.ToLower(v.Name)
		port.Port = v.Port
		port.TargetPort = intstr.FromInt(int(v.Port))
		if v.IsTcp {
			port.Protocol = k8sCoreV1.ProtocolTCP
		} else {
			port.Protocol = k8sCoreV1.ProtocolUDP
		}
		ports = append(ports, port)
	}
	return ports
}

func buildService(server *crdV1Alpha1.TServer) *k8sCoreV1.Service {
	service := &k8sCoreV1.Service{
		ObjectMeta: k8sMetaV1.ObjectMeta{
			Name:      server.Name,
			Namespace: server.Namespace,
			Labels: map[string]string{
				TServerAppLabel:  server.Spec.App,
				TServerNameLabel: server.Spec.Server,
			},
			OwnerReferences: []k8sMetaV1.OwnerReference{
				{
					APIVersion: TServerAPIVersion,
					Kind:       TServerKind,
					Name:       server.Name,
					UID:        server.UID,
				},
			},
		},
		Spec: k8sCoreV1.ServiceSpec{
			Selector: map[string]string{
				TServerAppLabel:  server.Spec.App,
				TServerNameLabel: server.Spec.Server,
			},
			ClusterIP: k8sCoreV1.ClusterIPNone,
			Type:      k8sCoreV1.ServiceTypeClusterIP,
		},
	}
	if server.Spec.Tars != nil {
		service.Spec.Ports = buildServicePortsFromTServant(server)
	} else if server.Spec.Normal != nil {
		service.Spec.Ports = buildServicePortsFromNPorts(server)
	}
	return service
}

func buildDaemonSet(server *crdV1Alpha1.TServer) *k8sAppsV1.DaemonSet {
	daemonSet := &k8sAppsV1.DaemonSet{
		TypeMeta: k8sMetaV1.TypeMeta{},
		ObjectMeta: k8sMetaV1.ObjectMeta{
			Name:      server.Name,
			Namespace: server.Namespace,
			OwnerReferences: []k8sMetaV1.OwnerReference{
				{
					APIVersion: TServerAPIVersion,
					Kind:       TServerKind,
					Name:       server.Name,
					UID:        server.UID,
				},
			},
			Labels: map[string]string{
				TServerAppLabel:  server.Spec.App,
				TServerNameLabel: server.Spec.Server,
			},
		},
		Spec: k8sAppsV1.DaemonSetSpec{
			Selector: &k8sMetaV1.LabelSelector{
				MatchLabels: map[string]string{
					TServerAppLabel:  server.Spec.App,
					TServerNameLabel: server.Spec.Server,
				},
			},
			Template: buildPodTemplate(server),
		},
	}
	return daemonSet
}

func buildTEndpoint(server *crdV1Alpha1.TServer) *crdV1Alpha1.TEndpoint {
	endpoint := &crdV1Alpha1.TEndpoint{
		ObjectMeta: k8sMetaV1.ObjectMeta{
			Name:      server.Name,
			Namespace: server.Namespace,
			OwnerReferences: []k8sMetaV1.OwnerReference{
				{
					APIVersion: TServerAPIVersion,
					Kind:       TServerKind,
					Name:       server.Name,
					UID:        server.UID,
				},
			},
			Labels: map[string]string{
				TServerAppLabel:  server.Spec.App,
				TServerNameLabel: server.Spec.Server,
			},
		},
		Spec: crdV1Alpha1.TEndpointSpec{
			App:       server.Spec.App,
			Server:    server.Spec.Server,
			SubType:   server.Spec.SubType,
			Important: server.Spec.Important,
			Tars:       server.Spec.Tars,
			Normal:    server.Spec.Normal,
			HostPorts: server.Spec.K8S.HostPorts,
		},
	}
	return endpoint
}

func buildTExitedRecord(server *crdV1Alpha1.TServer) *crdV1Alpha1.TExitedRecord {
	tExitedPod := &crdV1Alpha1.TExitedRecord{
		TypeMeta: k8sMetaV1.TypeMeta{},
		ObjectMeta: k8sMetaV1.ObjectMeta{
			Name:      server.Name,
			Namespace: server.Namespace,
			OwnerReferences: []k8sMetaV1.OwnerReference{
				{
					APIVersion: TServerAPIVersion,
					Kind:       TServerKind,
					Name:       server.Name,
					UID:        server.UID,
				},
			},
			Labels: map[string]string{
				TServerAppLabel:  server.Spec.App,
				TServerNameLabel: server.Spec.Server,
			},
		},
		App:    server.Spec.App,
		Server: server.Spec.Server,
		Pods:   []crdV1Alpha1.TExitedPod{},
	}
	return tExitedPod
}

func syncTEndpoint(server *crdV1Alpha1.TServer, endpoint *crdV1Alpha1.TEndpoint) {
	endpoint.Labels = map[string]string{
		TServerAppLabel:  server.Spec.App,
		TServerNameLabel: server.Spec.Server,
	}
	endpoint.OwnerReferences = []k8sMetaV1.OwnerReference{{
		APIVersion: TServerAPIVersion,
		Kind:       TServerKind,
		Name:       server.Name,
		UID:        server.UID,
	},
	}
	endpoint.Spec.App = server.Spec.App
	endpoint.Spec.Server = server.Spec.Server
	endpoint.Spec.SubType = server.Spec.SubType
	endpoint.Spec.Important = server.Spec.Important
	endpoint.Spec.Tars = server.Spec.Tars
	endpoint.Spec.Normal = server.Spec.Normal
	endpoint.Spec.HostPorts = server.Spec.K8S.HostPorts
}

func syncService(server *crdV1Alpha1.TServer, service *k8sCoreV1.Service) {
	if server.Spec.Tars != nil {
		service.Spec.Ports = buildServicePortsFromTServant(server)
	} else if server.Spec.Normal != nil {
		service.Spec.Ports = buildServicePortsFromNPorts(server)
	}
}

func syncStatefulSet(server *crdV1Alpha1.TServer, statefulSet *k8sAppsV1.StatefulSet) {
	statefulSet.Labels = map[string]string{
		TServerAppLabel:  server.Spec.App,
		TServerNameLabel: server.Spec.Server,
	}
	statefulSet.OwnerReferences = []k8sMetaV1.OwnerReference{{
		APIVersion: TServerAPIVersion,
		Kind:       TServerKind,
		Name:       server.Name,
		UID:        server.UID,
	},
	}
	statefulSet.Spec.Replicas = &server.Spec.K8S.Replicas
	statefulSet.Spec.Selector = &k8sMetaV1.LabelSelector{
		MatchLabels: map[string]string{
			TServerAppLabel:  server.Spec.App,
			TServerNameLabel: server.Spec.Server,
		},
	}

	statefulSet.Spec.Template = buildPodTemplate(server)
}

func syncDaemonSet(server *crdV1Alpha1.TServer, daemonSet *k8sAppsV1.DaemonSet) {
	daemonSet.Labels = map[string]string{
		TServerAppLabel:  server.Spec.App,
		TServerNameLabel: server.Spec.Server,
	}
	daemonSet.OwnerReferences = []k8sMetaV1.OwnerReference{{
		APIVersion: TServerAPIVersion,
		Kind:       TServerKind,
		Name:       server.Name,
		UID:        server.UID,
	},
	}
	daemonSet.Spec.Selector = &k8sMetaV1.LabelSelector{
		MatchLabels: map[string]string{
			TServerAppLabel:  server.Spec.App,
			TServerNameLabel: server.Spec.Server,
		},
	}
	daemonSet.Spec.Template = buildPodTemplate(server)
}
