package base

type DaemonPodK8S struct {
	PodName			string	`json:"pon_name"`
	ContainerName 	string 	`json:"container_name"`
	NodeName		string	`json:"node_name"`
	HostIP 			string	`json:"ip"`
	HostPort 		int32	`json:"port"`
	Ready			bool	`json:"ready"`
}

func NewDaemonPodK8S() *DaemonPodK8S {
	return &DaemonPodK8S{Ready: false}
}
