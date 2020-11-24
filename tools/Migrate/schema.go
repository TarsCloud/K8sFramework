package main

type Servant struct {
	Name  string `yaml:"name"`
	Port  int    `yaml:"port"`
	IsTars bool   `yaml:"isTars"`
}

type Hostport struct {
	NameRef string `yaml:"nameRef"`
	Port    int    `yaml:"port"`
}

type NodeSelectorValues struct {
	Values []string `yaml:"values"`
}

type NodeSelector struct {
	AbilityPool NodeSelectorValues `yaml:"abilityPool"` // 目前ability没有填值，但一定必须有oneOf
	NodeBind NodeSelectorValues `yaml:"nodeBind"`
}

type TServer struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		App     string `yaml:"app"`
		Server  string `yaml:"server"`
		SubType string `yaml:"subType"`
		Tars     struct {
			Template string `yaml:"template"`
			Profile  string `yaml:"profile"`
			Servants []Servant `yaml:"servants"`
		} `yaml:"tars"`
		K8S struct {
			Replicas     int `yaml:"replicas"`
			//NodeSelector nodeSelector `yaml:"nodeSelector"`
			NodeSelector map[string]interface{} `yaml:"nodeSelector"`
			HostPorts []Hostport `yaml:"hostPorts,omitempty"`
			Env []struct {
				Name      string `yaml:"name"`
				ValueFrom struct {
					FieldRef struct {
						FieldPath string `yaml:"fieldPath"`
					} `yaml:"fieldRef"`
				} `yaml:"valueFrom"`
			} `yaml:"env"`
			Mounts []struct {
				Name   string `yaml:"name"`
				Source struct {
					HostPath struct {
						Path string `yaml:"path"`
						Type string `yaml:"type"`
					} `yaml:"hostPath"`
				} `yaml:"source"`
				MountPath   string `yaml:"mountPath"`
				SubPathExpr string `yaml:"subPathExpr"`
			} `yaml:"mounts"`
		} `yaml:"k8s"`
		Release struct {
			Source          string `yaml:"source"`
			Tag             string `yaml:"tag"`
			Image           string `yaml:"image"`
			ImagePullSecret string `yaml:"imagePullSecret"`
			ServerType		string `yaml:"serverType"`
		} `yaml:"release"`
	} `yaml:"spec"`
}

type ReleaseImageItem struct {
	Image           string `yaml:"image"`
	ImagePullSecret string `yaml:"imagePullSecret"`
	Tag             string `yaml:"tag"`
	ServerType 		string `yaml:"serverType"`
}

type TRelease struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		List []ReleaseImageItem `yaml:"list"`
	} `yaml:"spec"`
}

type TTemplate struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		Content string `yaml:"content"`
		Parent  string `yaml:"parent"`
	} `yaml:"spec"`
}

type TConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	ServerConfig struct {
		App     string `yaml:"app"`
		Server  string `yaml:"server"`
		ConfigName     string `yaml:"configName"`
		ConfigContent  string `yaml:"configContent"`
		UpdatePerson	string `yaml:"updatePerson"`
		UpdateTime		string `yaml:"updateTime"`
		UpdateReason	string `yaml:"updateReason"`
	} `yaml:"serverConfig"`
}


type HelmChart struct {
	APIVersion   string   `yaml:"apiVersion"`
	Name         string   `yaml:"name"`
	Version      string   `yaml:"version"`
	KubeVersion  string   `yaml:"kubeVersion"`
	Type         string   `yaml:"type"`
	Keywords     []string `yaml:"keywords"`
	Dependencies []struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"dependencies"`
	AppVersion string `yaml:"appVersion"`
}