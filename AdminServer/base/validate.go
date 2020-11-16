package base

import (
	"github.com/asaskevich/govalidator"
	"math"
	"regexp"
	"strings"
	"unicode/utf8"
)

func init() {
	govalidator.TagMap["matches-ServiceId"] = func(str string) bool {
		value64, _ := govalidator.ToInt(str)
		return value64 > 0 && value64 <= math.MaxInt32
	}

	govalidator.TagMap["matches-DeployId"] = func(str string) bool {
		value64, _ := govalidator.ToInt(str)
		return value64 > 0 && value64 <= math.MaxInt32
	}

	govalidator.TagMap["matches-ServerId"] = func(str string) bool {
		value64, _ := govalidator.ToInt(str)
		return value64 > 0 && value64 <= math.MaxInt32
	}

	govalidator.CustomTypeTagMap.Set("each-matches-ServerId", func(i interface{}, o interface{}) bool {
		switch i.(type) {
		case []int:
			var array = i.([]int)
			var matchesFun = govalidator.TagMap["matches-ServerId"]
			for _, v := range array {
				if !matchesFun(govalidator.ToString(v)) {
					return false
				}
			}
			return true
		default:
			return false
		}
	})

	govalidator.TagMap["matches-ServerApp"] = func(str string) bool {
		match, _ := regexp.MatchString(`^[a-zA-Z0-9]{1,24}\z`, str)
		return match
	}

	govalidator.CustomTypeTagMap.Set("each-matches-ServerApp", func(i interface{}, o interface{}) bool {
		switch i.(type) {
		case []string:
			var array = i.([]string)
			var matchesFun = govalidator.TagMap["matches-ServerApp"]
			for _, v := range array {
				if !matchesFun(v) {
					return false
				}
			}
			return true
		default:
			return false
		}
	})

	govalidator.TagMap["matches-ServerName"] = func(str string) bool {
		match, _ := regexp.MatchString(`^[a-zA-Z0-9]{1,64}\z`, str)
		return match
	}

	govalidator.TagMap["matches-ServerType"] = func(str string) bool {
		if len(str) >= 16 {
			return false
		}
		return str == "taf_cpp" || str == "taf_java_jar" || str == "taf_java_war" || str == "taf_node" || str == "taf_node8" || str == "taf_node10" || str == "taf_node_pkg"
	}

	govalidator.TagMap["matches-ServerApps"] = func(str string) bool {
		if len(str) >= 57 {
			return false
		}
		v := strings.Split(str, `.`)
		if len(v) == 1 {
			return govalidator.TagMap["matches-ServerApp"](str)
		}

		if len(v) == 2 {
			return govalidator.TagMap["matches-ServerApp"](v[0]) && govalidator.TagMap["matches-ServerName"](v[1])
		}
		return false
	}

	govalidator.CustomTypeTagMap.Set("each-matches-ServerApps", func(i interface{}, o interface{}) bool {
		switch i.(type) {
		case []string:
			var array = i.([]string)
			var matchesFun = govalidator.TagMap["matches-ServerApps"]
			for _, v := range array {
				if !matchesFun(v) {
					return false
				}
			}
			return true
		default:
			return false
		}
	})

	govalidator.TagMap["matches-BusinessName"] = func(str string) bool {
		if str == "" {
			return true
		}
		match, _ := regexp.MatchString(`^[a-zA-Z0-9.:_-]{1,128}\z`, str)
		return match
	}

	govalidator.TagMap["matches-ServiceImage"] = func(str string) bool {
		if str == "" || str == " " {
			return true
		}
		match, _ := regexp.MatchString(`^[a-zA-Z0-9.:/-]{1,256}\z`, str)
		return match
	}

	govalidator.TagMap["matches-ServiceVersion"] = func(str string) bool {
		if str == "" || str == " " {
			return true
		}
		match, _ := regexp.MatchString(`^1[0-9]{4}\z`, str)
		return match
	}

	govalidator.CustomTypeTagMap.Set("each-matches-BusinessName", func(i interface{}, o interface{}) bool {
		switch i.(type) {
		case []string:
			var array = i.([]string)
			var matchesFun = govalidator.TagMap["matches-BusinessName"]
			for _, v := range array {
				if !matchesFun(v) {
					return false
				}
			}
			return true
		default:
			return false
		}
	})

	govalidator.TagMap["matches-BusinessShow"] = func(str string) bool {
		return utf8.RuneCountInString(str) <= 32
	}

	govalidator.TagMap["matches-BusinessOrder"] = func(str string) bool {
		value64, _ := govalidator.ToInt(str)
		return value64 > 0 && value64 <= 100
	}

	govalidator.TagMap["matches-ServantName"] = func(str string) bool {
		match, _ := regexp.MatchString(`^[a-zA-Z1-9]{1,64}Obj\z`, str)
		return match
	}

	govalidator.TagMap["matches-ServantPort"] = func(str string) bool {
		value64, _ := govalidator.ToInt(str)
		return value64 > 0 && value64 <= 65535 && value64 != 19385
	}

	govalidator.TagMap["matches-ServantHostPort"] = func(str string) bool {
		value64, _ := govalidator.ToInt(str)
		return value64 >= 0 && value64 <= 65535
	}

	govalidator.TagMap["matches-ServantThreads"] = func(str string) bool {
		value64, _ := govalidator.ToInt(str)
		return value64 > 0 && value64 <= 50
	}

	govalidator.TagMap["matches-ServantConnections"] = func(str string) bool {
		value64, _ := govalidator.ToInt(str)
		return value64 > 0 && value64 <= 30000000
	}

	govalidator.TagMap["matches-ServantCapacity"] = func(str string) bool {
		value64, _ := govalidator.ToInt(str)
		return value64 > 0 && value64 <= 3000000
	}

	govalidator.TagMap["matches-ServantTimeout"] = func(str string) bool {
		value64, _ := govalidator.ToInt(str)
		return value64 > 0 && value64 <= 1200000
	}

	govalidator.CustomTypeTagMap.Set("matches-ServerServant", func(i interface{}, o interface{}) bool {
		switch i.(type) {
		case ServerServant:
			break
		default:
			return false
		}

		s := i.(ServerServant)
		var port = make(map[int]interface{}, len(s))
		var name = make(map[string]interface{}, len(s))

		for k, v := range s {
			if k != strings.ToLower(v.Name) {
				return false
			}
			if ok, _ := govalidator.ValidateStruct(v); !ok {
				return false
			}
			if _, ok := port[v.Port]; ok {
				return false
			}

			if _, ok := name[v.Name]; ok {
				return false
			}
			name[v.Name] = nil
			port[v.Port] = nil
		}

		return true
	})

	govalidator.CustomTypeTagMap.Set("matches-ServerOption", func(i interface{}, o interface{}) bool {
		switch i.(type) {
		case ServerOption:
			break
		default:
			return false
		}

		opt := i.(ServerOption)
		if opt.ServerImportant < 0 || opt.ServerImportant > 5 {
			return false
		}
		if opt.AsyncThread < 0 || opt.ServerImportant > 12 {
			return false
		}

		return true
	})

	govalidator.TagMap["matches-ServerImportant"] = func(str string) bool {
		value64, _ := govalidator.ToInt(str)
		return value64 >= 0 && value64 <= 20
	}

	govalidator.TagMap["matches-ServerStartScript"] = func(str string) bool {
		return true
	}

	govalidator.TagMap["matches-ServerStopScript"] = func(str string) bool {
		return true
	}

	govalidator.TagMap["matches-ServerMonitorScript"] = func(str string) bool {
		return true
	}

	govalidator.CustomTypeTagMap.Set("matches-ServerK8S", func(i interface{}, o interface{}) bool {

		var serverK8S *ServerK8S

		switch i.(type) {
		case ServerK8S:
			serverK8S_ := i.(ServerK8S)
			serverK8S = &serverK8S_
			break
		case *ServerK8S:
			serverK8S = i.(*ServerK8S)
		default:
			return false
		}

		if serverK8S.NodeSelector.Kind != NodeBind && serverK8S.NodeSelector.Kind != AbilityPool && serverK8S.NodeSelector.Kind != PublicPool {
			return false
		}

		if serverK8S.NodeSelector.Kind == NodeBind {
			checkNodeSelectorValueFun, _ := govalidator.CustomTypeTagMap.Get("each-matches-NodeName")
			if !checkNodeSelectorValueFun(serverK8S.NodeSelector.Value, nil) {
				return false
			}

			if serverK8S.NotStacked && int(serverK8S.Replicas) > len(serverK8S.NodeSelector.Value) {
				return false
			}
		}

		if serverK8S.NodeSelector.Kind != NodeBind {

			if serverK8S.HostIpc || serverK8S.HostNetwork {
				return false
			}

			if serverK8S.NotStacked && int(serverK8S.Replicas) > len(serverK8S.NodeSelector.Value) {
				return false
			}
		}

		checkImageFun, _ := govalidator.TagMap["matches-ServiceImage"]
		if !checkImageFun(serverK8S.Image) {
			return false
		}

		return true
	})

	govalidator.TagMap["matches-ServerAsyncThread"] = func(str string) bool {
		value64, _ := govalidator.ToInt(str)
		return value64 >= 0 && value64 <= 20
	}

	govalidator.TagMap["matches-NodeName"] = func(str string) bool {
		match, _ := regexp.MatchString(`[0-9a-zA-Z.-]{1,256}\z`, str)
		return match
	}

	govalidator.CustomTypeTagMap.Set("each-matches-NodeName", func(i interface{}, o interface{}) bool {
		switch i.(type) {
		case []string:
			var array = i.([]string)
			var matchesFun = govalidator.TagMap["matches-NodeName"]
			for _, v := range array {
				if !matchesFun(v) {
					return false
				}
			}
			return true
		default:
			return false
		}
	})

	govalidator.TagMap["matches-TemplateId"] = func(str string) bool {
		value64, _ := govalidator.ToInt(str)
		return value64 > 0 && value64 <= math.MaxInt32
	}

	govalidator.TagMap["matches-TemplateName"] = func(str string) bool {
		match, _ := regexp.MatchString(`[0-9a-zA-Z.-]{5,64}\z`, str)
		return match
	}

	govalidator.TagMap["matches-AdapterId"] = func(str string) bool {
		value64, _ := govalidator.ToInt(str)
		return value64 > 0 && value64 <= math.MaxInt32
	}

	govalidator.TagMap["matches-Domain"] = func(str string) bool {
		match, _ := regexp.MatchString(`^([0-9a-zA-Z-*]{1,63}\.)?([0-9a-zA-Z-]{1,63}\.){0,10}([a-zA-Z]{2,3})$`, str)
		return match
	}

	govalidator.TagMap["matches-ConfigId"] = func(str string) bool {
		value64, _ := govalidator.ToInt(str)
		return value64 >= 1 && value64 <= math.MaxInt32
	}

	govalidator.TagMap["matches-ConfigName"] = func(str string) bool {
		match, _ := regexp.MatchString(`^[0-9a-zA-Z.-]{1,64}\z`, str)
		return match
	}

	govalidator.TagMap["matches-ConfigPodSeq"] = func(str string) bool {
		value64, _ := govalidator.ToInt(str)
		return value64 >= -1 && value64 <= 32
	}

	govalidator.TagMap["matches-HistoryConfigId"] = func(str string) bool {
		value64, _ := govalidator.ToInt(str)
		return value64 >= 1 && value64 <= math.MaxInt32
	}

	govalidator.CustomTypeTagMap.Set("each-matches-IPv4", func(i interface{}, o interface{}) bool {
		switch i.(type) {
		case []string:
			var array = i.([]string)
			for _, v := range array {
				if !govalidator.IsIPv4(v) {
					return false
				}
			}
			return true
		default:
			return false
		}
	})
}
