package k8s

import (
	"math/rand"
	"regexp"
	"strings"
	"tarsadmin/openapi/models"
	"time"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	crdv1alpha1 "k8s.tars.io/api/crd/v1alpha1"
)

var K8sOption *K8SOption
var K8sWatcher *Watcher

var KeyLabel = map[string]string{
	"ServerApp":  TServerAppLabel,
	"ServerName": TServerNameLabel,
}

func LikeMatch(key string, filter models.MapString, field string) (bool, error) {
	like, ok := filter[key]
	if !ok {
		return true, nil
	}
	//return regexp.MatchString(like2Regex(like), field)
	return regexp.MatchString(like, field)
}

func PageList(itemsLen int, limiter *models.SelectRequestLimiter) (int32, int32) {
	start := *limiter.Offset
	if start < 0 {
		start = 0
	} else if start >= int32(itemsLen) {
		start = int32(itemsLen)
	}

	stop := *limiter.Offset + limiter.Rows - 1
	if stop < 0 || stop >= int32(itemsLen) {
		stop = int32(itemsLen)
	}

	return start, stop
}

func BuildDoubleEqualSelector(filter *models.SelectRequestFilter, keyLabel map[string]string) []labels.Requirement {
	requirements := make([]labels.Requirement, 0, len(keyLabel))
	for key, label := range keyLabel {
		pattern, ok := filter.Eq[key]
		if ok && pattern != "" {
			requirement, _ := labels.NewRequirement(label, selection.DoubleEquals, []string{pattern.(string)})
			requirements = append(requirements, *requirement)
		}
	}
	return requirements
}

func BuildSubTypeTafSelector() []labels.Requirement {
	requirement, _ := labels.NewRequirement(TSubTypeLabel, selection.DoubleEquals, []string{string(crdv1alpha1.TARS)})
	return []labels.Requirement{*requirement}
}

func BuildTafAppSelector(appName string) []labels.Requirement {
	requirements := make([]labels.Requirement, 0, 2)

	r1, _ := labels.NewRequirement(TSubTypeLabel, selection.DoubleEquals, []string{string(crdv1alpha1.TARS)})
	requirements = append(requirements, *r1)

	r2, _ := labels.NewRequirement(TServerAppLabel, selection.DoubleEquals, []string{appName})
	requirements = append(requirements, *r2)

	return requirements
}

func like2Regex(like string) string {
	r1 := []byte("^")
	r2 := []byte("$")
	r3 := []byte(".*")

	lc := []byte(like)

	if !strings.Contains(like, "%") {
		lc = append(r1, lc...)
		lc = append(lc, r2...)
	} else {
		if lc[0] == '%' {
			lc = append(r3, lc[1:]...)
		} else {
			lc = append(r1, lc...)
		}

		if lc[len(lc)-1] == '%' {
			lc = append(lc[0:len(lc)-1], r3...)
		} else {
			lc = append(lc, r2...)
		}
	}

	return strings.ReplaceAll(string(lc), "%", ".*")
}

// 冗余代码：否则排序时都要遍历一轮将：[]Type -> []interface{}
type TEndpointWrapper struct {
	Endpoint []*crdv1alpha1.TEndpoint
	By       func(e1, e2 *crdv1alpha1.TEndpoint) bool
}

func (e TEndpointWrapper) Len() int {
	return len(e.Endpoint)
}

func (e TEndpointWrapper) Swap(i, j int) {
	e.Endpoint[i], e.Endpoint[j] = e.Endpoint[j], e.Endpoint[i]
}

func (e TEndpointWrapper) Less(i, j int) bool {
	return e.By(e.Endpoint[i], e.Endpoint[j])
}

type AppWrapper struct {
	App []*models.AppElem
	By  func(e1, e2 *models.AppElem) bool
}

func (e AppWrapper) Len() int {
	return len(e.App)
}

func (e AppWrapper) Swap(i, j int) {
	e.App[i], e.App[j] = e.App[j], e.App[i]
}

func (e AppWrapper) Less(i, j int) bool {
	return e.By(e.App[i], e.App[j])
}

type ServerElemWrapper struct {
	Server []*models.ServerElem
	By     func(e1, e2 *models.ServerElem) bool
}

func (e ServerElemWrapper) Len() int {
	return len(e.Server)
}

func (e ServerElemWrapper) Swap(i, j int) {
	e.Server[i], e.Server[j] = e.Server[j], e.Server[i]
}

func (e ServerElemWrapper) Less(i, j int) bool {
	return e.By(e.Server[i], e.Server[j])
}

type ExitedPodElemEx struct {
	AppName    string
	ServerName string
	crdv1alpha1.TExitedPod
}

type TExitedPodWrapper struct {
	ExitedPod []ExitedPodElemEx
	By        func(e1, e2 *ExitedPodElemEx) bool
}

func (e TExitedPodWrapper) Len() int {
	return len(e.ExitedPod)
}

func (e TExitedPodWrapper) Swap(i, j int) {
	e.ExitedPod[i], e.ExitedPod[j] = e.ExitedPod[j], e.ExitedPod[i]
}

func (e TExitedPodWrapper) Less(i, j int) bool {
	return e.By(&e.ExitedPod[i], &e.ExitedPod[j])
}

type TDeployWrapper struct {
	Deploy []*crdv1alpha1.TDeploy
	By     func(e1, e2 *crdv1alpha1.TDeploy) bool
}

func (e TDeployWrapper) Len() int {
	return len(e.Deploy)
}

func (e TDeployWrapper) Swap(i, j int) {
	e.Deploy[i], e.Deploy[j] = e.Deploy[j], e.Deploy[i]
}

func (e TDeployWrapper) Less(i, j int) bool {
	return e.By(e.Deploy[i], e.Deploy[j])
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandStringRunes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
