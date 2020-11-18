package k8s

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"io/ioutil"
	"net/http"
	"tafadmin/openapi/models"
	"tafadmin/openapi/restapi/operations/agent"
)

type SelectAvailHostPortHandler struct {}

func (s *SelectAvailHostPortHandler) Handle(params agent.SelectAvailHostPortParams) middleware.Responder {
	pod, ok := getDaemonPodByName(*params.NodeName)
	if !ok {
		return agent.NewSelectAvailHostPortInternalServerError().WithPayload(&models.Error{Code: -1, Message: fmt.Sprintf("Can Not Find DaemonSet In %s.", *params.NodeName)})
	}
	containers := pod.Spec.Containers
	if len(containers) <= 0 {
		return agent.NewSelectAvailHostPortInternalServerError().WithPayload(&models.Error{Code: -1, Message: fmt.Sprintf("%s.%s Has No Container.", *params.NodeName, TafAgentDaemonSetName)})
	}
	ports := containers[0].Ports
	if len(ports) <= 0 {
		return agent.NewSelectAvailHostPortInternalServerError().WithPayload(&models.Error{Code: -1, Message: fmt.Sprintf("%s.%s Has No Host Port.", *params.NodeName, TafAgentDaemonSetName)})
	}

	hostIp := pod.Status.HostIP
	hostPort := containers[0].Ports[0].HostPort

	// proxy forward
	rsp, err := http.Get(fmt.Sprintf("http://%s:%d/port?host=%s&port=%d", hostIp, hostPort, hostIp, *params.Port))
	if err != nil {
		return agent.NewSelectAvailHostPortInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return agent.NewSelectAvailHostPortInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	if rsp.StatusCode == 500 {
		return agent.NewSelectAvailHostPortInternalServerError().WithPayload(&models.Error{Code: -1, Message: string(body)})
	} else {
		return agent.NewSelectAvailHostPortOK().WithPayload(&agent.SelectAvailHostPortOKBody{Result: string(body)})
	}
}
