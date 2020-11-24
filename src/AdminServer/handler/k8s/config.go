package k8s

import (
	"database/sql"
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
	"tarsadmin/handler/mysql"
	"tarsadmin/handler/util"
	"tarsadmin/openapi/models"
	"tarsadmin/openapi/restapi/operations/config"
	"tarsadmin/openapi/restapi/operations/server_pod"
	"tarsadmin/openapi/restapi/operations/template"

	"github.com/go-openapi/runtime/middleware"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	crdv1alpha1 "k8s.tars.io/api/crd/v1alpha1"

	tafConf "github.com/TarsCloud/TarsGo/tars/util/conf"
)

var configHistoryColumnsSqlColumnsMap = mysql.RequestColumnSqlColumnMap{
	"HistoryId": mysql.SqlColumn{
		ColumnName: "f_history_id",
		ColumnType: "int",
	},
	"ConfigId": mysql.SqlColumn{
		ColumnName: "f_config_id",
		ColumnType: "int",
	},
	"ConfigName": mysql.SqlColumn{
		ColumnName: "f_config_name",
		ColumnType: "string",
	},
	"ConfigVersion": mysql.SqlColumn{
		ColumnName: "f_config_version",
		ColumnType: "int",
	},
	"ConfigContent": mysql.SqlColumn{
		ColumnName: "f_config_content",
		ColumnType: "string",
	},
	"CreateTime": mysql.SqlColumn{
		ColumnName: "f_create_time",
		ColumnType: "string",
	},
	"CreatePerson": mysql.SqlColumn{
		ColumnName: "f_create_person",
		ColumnType: "string",
	},
	"ConfigMark": mysql.SqlColumn{
		ColumnName: "f_config_mark",
		ColumnType: "string",
	},
	"ServerApps": mysql.SqlColumn{
		ColumnName: "f_app_server",
		ColumnType: "string",
	},
}

type CreateServerConfigHandler struct{}

func (s *CreateServerConfigHandler) Handle(params config.CreateServerConfigParams) middleware.Responder {
	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tConfig, err := buildTConfig(namespace, metadata)
	if err != nil {
		return config.NewCreateServerConfigInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	tConfigInterface := K8sOption.CrdClientSet.CrdV1alpha1().TConfigs(namespace)
	if _, err := tConfigInterface.Create(context.TODO(), tConfig, k8sMetaV1.CreateOptions{}); err != nil && !errors.IsAlreadyExists(err) {
		return config.NewCreateServerConfigInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return config.NewCreateServerConfigOK().WithPayload(&config.CreateServerConfigOKBody{Result: 0})
}

type SelectServerConfigHandler struct{}

func (s *SelectServerConfigHandler) Handle(params config.SelectServerConfigParams) middleware.Responder {
	// fetch list
	namespace := K8sOption.Namespace

	// parse query
	selectParams, err := util.ParseSelectQuery(params.Filter, params.Limiter, nil)
	if err != nil {
		return config.NewSelectServerConfigInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	listAll := true
	if selectParams.Filter != nil && selectParams.Filter.Eq != nil {
		_, appOk := selectParams.Filter.Eq["ServerApp"]
		_, nameOk := selectParams.Filter.Eq["ServerName"]
		if appOk || nameOk {
			listAll = false
		}
	}

	tConfigInterface := K8sOption.CrdClientSet.CrdV1alpha1().TConfigs(namespace)

	allItems := make([]crdv1alpha1.TConfig, 0, 10)
	if listAll {
		return config.NewSelectServerConfigInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Invalid Select Query Request."})
	} else {
		requirements := BuildDoubleEqualSelector(selectParams.Filter, KeyLabel)

		// 查询某个configName的节点配置
		pattern, ok := selectParams.Filter.Eq["ConfigName"]
		if ok && pattern != "" {
			requirement, _ := labels.NewRequirement(TConfigNameLabel, selection.DoubleEquals, []string{pattern.(string)})
			requirements = append(requirements, *requirement)
			requirement, _ = labels.NewRequirement(TConfigPodSeqLabel, selection.NotEquals, []string{"m"})
			requirements = append(requirements, *requirement)
		} else {
			requirement, _ := labels.NewRequirement(TConfigPodSeqLabel, selection.DoubleEquals, []string{"m"})
			requirements = append(requirements, *requirement)
		}
		list, err := tConfigInterface.List(context.TODO(), k8sMetaV1.ListOptions{LabelSelector: labels.NewSelector().Add(requirements...).String()})
		if err != nil {
			return server_pod.NewSelectPodAliveInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}
		allItems = list.Items
	}

	// filter
	filterItems := allItems
	// Admin临时版本，Template特化已有web实现
	if selectParams.Filter != nil {
		filterItems = make([]crdv1alpha1.TConfig, 0, len(allItems))
		for _, elem := range allItems {
			if selectParams.Filter.Ne != nil {
				return template.NewSelectTemplateInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Ne Is Not Supported."})
			}
			if selectParams.Filter.Like != nil {
				return template.NewSelectTemplateInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Like Is Not Supported."})
			}
			filterItems = append(filterItems, elem)
		}
	}

	// limiter
	if selectParams.Limiter != nil {
		start, stop := PageList(len(filterItems), selectParams.Limiter)
		filterItems = filterItems[start:stop]
	}

	// Count填充
	result := &models.SelectResult{}
	result.Count = make(models.MapInt)
	result.Count["AllCount"] = int32(len(filterItems))
	result.Count["FilterCount"] = int32(len(filterItems))

	// Data填充
	result.Data = make(models.ArrayMapInterface, 0, len(filterItems))
	for _, item := range filterItems {
		elem := make(map[string]interface{})
		elem["ConfigId"] = item.Name
		elem["ConfigVersion"] = "10000"

		if item.AppConfig != nil {
			elem["ServerId"] = item.AppConfig.App
			elem["ConfigName"] = item.AppConfig.ConfigName
			elem["ConfigContent"] = item.AppConfig.ConfigContent
			elem["CreateTime"] = item.AppConfig.UpdateTime
			elem["CreatePerson"] = item.AppConfig.UpdatePerson
			elem["ConfigMark"] = item.AppConfig.UpdateReason
		} else if item.ServerConfig != nil {
			elem["ServerId"] = util.GetServerId(item.ServerConfig.App, item.ServerConfig.Server)
			if item.ServerConfig.PodSeq != nil {
				elem["PodSeq"] = item.ServerConfig.PodSeq
			} else {
				elem["PodSeq"] = ""
			}
			elem["ConfigName"] = item.ServerConfig.ConfigName
			elem["ConfigContent"] = item.ServerConfig.ConfigContent
			elem["CreateTime"] = item.ServerConfig.UpdateTime
			elem["CreatePerson"] = item.ServerConfig.UpdatePerson
			elem["ConfigMark"] = item.ServerConfig.UpdateReason
		}

		result.Data = append(result.Data, elem)
	}

	return config.NewSelectServerConfigOK().WithPayload(result)
}

type UpdateServerConfigHandler struct{}

func (s *UpdateServerConfigHandler) Handle(params config.UpdateServerConfigParams) middleware.Responder {
	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tConfigInterface := K8sOption.CrdClientSet.CrdV1alpha1().TConfigs(namespace)

	tConfig, err := tConfigInterface.Get(context.TODO(), *metadata.ConfigID, k8sMetaV1.GetOptions{})
	if err != nil {
		return config.NewUpdateServerConfigInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	err = validTConfig(tConfig)
	if err != nil {
		return config.NewUpdateServerConfigInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	target := params.Params.Target

	conf := tafConf.New()
	if err := conf.InitFromString(target.ConfigContent); err != nil {
		return config.NewUpdateServerConfigInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	tConfigCopy := tConfig.DeepCopy()
	if tConfigCopy.AppConfig != nil {
		tConfigCopy.AppConfig.ConfigContent = target.ConfigContent
		tConfigCopy.AppConfig.UpdateReason = target.ConfigMark
		tConfigCopy.AppConfig.UpdateTime = k8sMetaV1.Now()
	} else if tConfigCopy.ServerConfig != nil {
		tConfigCopy.ServerConfig.ConfigContent = target.ConfigContent
		tConfigCopy.ServerConfig.UpdateReason = target.ConfigMark
		tConfigCopy.ServerConfig.UpdateTime = k8sMetaV1.Now()
	}

	_, err = tConfigInterface.Update(context.TODO(), tConfigCopy, k8sMetaV1.UpdateOptions{})
	if err != nil {
		return config.NewUpdateServerConfigInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return config.NewUpdateServerConfigOK().WithPayload(&config.UpdateServerConfigOKBody{Result: 0})
}

type DeleteServerConfigHandler struct{}

func (s *DeleteServerConfigHandler) Handle(params config.DeleteServerConfigParams) middleware.Responder {

	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tConfigInterface := K8sOption.CrdClientSet.CrdV1alpha1().TConfigs(namespace)

	// 由validating验证是否删除成功：比如app和pod配置不能脱离server配置而单独删除
	err := tConfigInterface.Delete(context.TODO(), *metadata.ConfigID, k8sMetaV1.DeleteOptions{})
	if err != nil {
		return config.NewDeleteServerConfigInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return config.NewDeleteServerConfigOK().WithPayload(&config.DeleteServerConfigOKBody{Result: 0})
}

type DoPreviewConfigContentHandler struct{}

func (s *DoPreviewConfigContentHandler) Handle(params config.DoPreviewConfigContentParams) middleware.Responder {

	namespace := K8sOption.Namespace
	metadata := params

	tConfigInterface := K8sOption.CrdClientSet.CrdV1alpha1().TConfigs(namespace)

	serverConfigId := getConfId(*metadata.ServerApp, *metadata.ServerName, *metadata.ConfigName)
	tServerConfig, err := tConfigInterface.Get(context.TODO(), serverConfigId, k8sMetaV1.GetOptions{})
	if err != nil {
		return config.NewDoPreviewConfigContentInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	podConfigId := getConfId(*metadata.ServerApp, *metadata.ServerName, *metadata.ConfigName, *metadata.PodSeq)
	tPodConfig, err := tConfigInterface.Get(context.TODO(), podConfigId, k8sMetaV1.GetOptions{})
	if err != nil {
		return config.NewDoPreviewConfigContentInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	result := fmt.Sprintf("%s\r\n\r\n%s", tServerConfig.ServerConfig.ConfigContent, tPodConfig.ServerConfig.ConfigContent)
	return config.NewDoPreviewConfigContentOK().WithPayload(&config.DoPreviewConfigContentOKBody{Result: result})
}

type SelectServerConfigHistoryHandler struct{}

func (s *SelectServerConfigHistoryHandler) Handle(params config.SelectServerConfigHistoryParams) middleware.Responder {
	/*
		const from = "t_config_history"

		result, err := mysql.SelectQueryResult(from, params.Filter, params.Limiter, params.Order, configHistoryColumnsSqlColumnsMap)
		if err != nil {
			return config.NewSelectServerConfigHistoryInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}

		return config.NewSelectServerConfigHistoryOK().WithPayload(result)
	*/
	return config.NewSelectServerConfigHistoryInternalServerError().WithPayload(&models.Error{Code: -1, Message: "History Config Has Not Been Supported."})
}

type DeleteServerConfigHistoryHandler struct {
	tafDb *sql.DB
}

func (s *DeleteServerConfigHistoryHandler) Handle(params config.DeleteServerConfigHistoryParams) middleware.Responder {
	/*
		metadata := params.Params.Metadata

		DeleteConfigHistoryResourceSql1 := "delete from t_config_history where f_history_id=?"
		if _, err := mysql.TafDb.Exec(DeleteConfigHistoryResourceSql1, metadata.HistoryID); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))
			return config.NewDeleteServerConfigHistoryInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}

		return config.NewDeleteServerConfigHistoryOK().WithPayload(&config.DeleteServerConfigHistoryOKBody{Result: 0})
	*/
	return config.NewDeleteServerConfigHistoryInternalServerError().WithPayload(&models.Error{Code: -1, Message: "History Config Has Not Been Supported."})
}

type DoActiveHistoryConfigHandler struct{}

func (s *DoActiveHistoryConfigHandler) Handle(params config.DoActiveHistoryConfigParams) middleware.Responder {
	/*
		metadata := params.Params.Metadata

		updateConfigSql1 := "insert into t_config_history (f_app_server, f_config_name, f_config_version, f_config_content,f_config_id,f_create_person,f_create_time, f_config_mark,f_pod_seq) select f_app_server,f_config_name,f_config_version,f_config_content,f_config_id,f_create_person,f_create_time,f_config_mark,f_pod_seq from t_config where f_config_id = (select f_config_id from t_config_history where f_history_id=?) ON DUPLICATE KEY UPDATE f_history_id=f_history_id"
		if _, err := mysql.TafDb.Exec(updateConfigSql1, metadata.HistoryID); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))
			return config.NewDoActiveHistoryConfigInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}

		updateConfigSql2 := "update t_config a inner join t_config_history b using (f_config_id) set a.f_config_version=b.f_config_version,a.f_config_content=b.f_config_content,a.f_create_person=b.f_create_person,a.f_config_content=b.f_config_content, a.f_create_time=b.f_create_time,a.f_config_mark=b.f_config_mark,a.f_pod_seq=b.f_pod_seq where b.f_history_id=?"
		if _, err := mysql.TafDb.Exec(updateConfigSql2, metadata.HistoryID); err != nil {
			_, file, line, _ := runtime.Caller(0)
			fmt.Println(fmt.Sprintf("file %s , Line: %d , Err: %s ", file, line, err.Error()))
			return config.NewDoActiveHistoryConfigInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
		}

		return config.NewDoActiveHistoryConfigOK().WithPayload(&config.DoActiveHistoryConfigOKBody{Result: 0})
	*/
	return config.NewDoPreviewConfigContentInternalServerError().WithPayload(&models.Error{Code: -1, Message: "History Config Has Not Been Supported."})
}

func getConfId(param ...string) string {
	var configId string

	n := len(param)
	for i := 0; i < n; i++ {
		// 特殊处理：默认第3个字段是configName，需要hash处理避免无效字符串
		if i == 2 {
			crc := crc32.ChecksumIEEE([]byte(param[i]))
			configId += strconv.FormatUint(uint64(crc), 10)
		} else {
			configId += param[i]
		}
		if i < (n - 1) {
			configId += "-"
		}
	}

	return strings.ToLower(configId)
}

func buildTConfig(namespace string, metadata *models.ConfigMeta) (*crdv1alpha1.TConfig, error) {

	tConfig := &crdv1alpha1.TConfig{
		ObjectMeta: k8sMetaV1.ObjectMeta{
			Namespace: namespace,
		},
	}

	conf := tafConf.New()
	if err := conf.InitFromString(*metadata.ConfigContent); err != nil {
		return nil, fmt.Errorf("Bad Schema : Bad Params.Metadata.ConfigContent Value. ")
	}

	if metadata.PodSeq != "" {
		// 节点配置
		if metadata.ServerApp == "" || metadata.ServerName == "" {
			return nil, fmt.Errorf("Bad Schema : Bad Params.Metadata.PodSeq Value. ")
		}
		tConfig.Name = getConfId(metadata.ServerApp, metadata.ServerName, *metadata.ConfigName, metadata.PodSeq)
		tConfig.ServerConfig = &crdv1alpha1.TConfigServer{
			App:           metadata.ServerApp,
			Server:        metadata.ServerName,
			ConfigName:    *metadata.ConfigName,
			ConfigContent: *metadata.ConfigContent,
			PodSeq:        &metadata.PodSeq,
			UpdateTime:    k8sMetaV1.Now(),
			UpdatePerson:  "TafAdmin",
			UpdateReason:  "Create",
		}
	} else if metadata.ServerName != "" {
		// 服务配置
		if metadata.ServerApp == "" {
			return nil, fmt.Errorf("Bad Schema : Bad Params.Metadata.ServerName Value. ")
		}
		tConfig.Name = getConfId(metadata.ServerApp, metadata.ServerName, *metadata.ConfigName)
		tConfig.ServerConfig = &crdv1alpha1.TConfigServer{
			App:           metadata.ServerApp,
			Server:        metadata.ServerName,
			ConfigName:    *metadata.ConfigName,
			ConfigContent: *metadata.ConfigContent,
			PodSeq:        nil,
			UpdateTime:    k8sMetaV1.Now(),
			UpdatePerson:  "TafAdmin",
			UpdateReason:  "Create",
		}
	} else if metadata.ServerApp != "" {
		// 应用配置
		tConfig.Name = getConfId(metadata.ServerApp, *metadata.ConfigName)
		tConfig.AppConfig = &crdv1alpha1.TConfigApp{
			App:           metadata.ServerApp,
			ConfigName:    *metadata.ConfigName,
			ConfigContent: *metadata.ConfigContent,
			UpdateTime:    k8sMetaV1.Now(),
			UpdatePerson:  "TafAdmin",
			UpdateReason:  "Create",
		}
	} else {
		if metadata.ServerApp == "" {
			return nil, fmt.Errorf("Bad Schema : Bad Params.Metadata.ServerApp Value. ")
		}
	}

	return tConfig, nil
}

func validTConfig(newTConfig *crdv1alpha1.TConfig) error {
	if newTConfig.AppConfig != nil {
		if newTConfig.Name != getConfId(newTConfig.AppConfig.App, newTConfig.AppConfig.ConfigName) {
			return fmt.Errorf("unexpected app config resources name")
		}
	}

	if newTConfig.ServerConfig != nil {
		if newTConfig.ServerConfig.PodSeq == nil {
			if newTConfig.Name != getConfId(newTConfig.ServerConfig.App, newTConfig.ServerConfig.Server, newTConfig.ServerConfig.ConfigName) {
				return fmt.Errorf("unexpected server config resources name")
			}
			return nil
		}

		if newTConfig.Name != getConfId(newTConfig.ServerConfig.App, newTConfig.ServerConfig.Server, newTConfig.ServerConfig.ConfigName, *newTConfig.ServerConfig.PodSeq) {
			return fmt.Errorf("unexpected server-podseq resources name")
		}

		masterConfigName := getConfId(newTConfig.ServerConfig.App, newTConfig.ServerConfig.Server, newTConfig.ServerConfig.ConfigName)
		_, err := K8sOption.CrdClientSet.CrdV1alpha1().TConfigs(K8sOption.Namespace).Get(context.TODO(), masterConfigName, k8sMetaV1.GetOptions{})
		if err != nil {
			return fmt.Errorf("get %s:%s in %s err: %s", "tconfig", masterConfigName, K8sOption.Namespace, err)
		}
	}

	return nil
}
