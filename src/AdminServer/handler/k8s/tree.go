package k8s

import (
	"github.com/go-openapi/runtime/middleware"
	"k8s.io/apimachinery/pkg/labels"
	"sort"
	"tafadmin/handler/util"
	"tafadmin/openapi/models"
	"tafadmin/openapi/restapi/operations/tree"
)

type SelectServerTreeHandler struct {}

func (s *SelectServerTreeHandler) Handle(params tree.SelectServerTreeParams) middleware.Responder {

	namespace := K8sOption.Namespace

	var allApp = make(map[string][]*models.ServerElem, 100)

	var targetBusiness *models.Business

	var lastBusiness string

	requirements := BuildSubTypeTafSelector()
	list, err := K8sWatcher.tServerLister.TServers(namespace).List(labels.NewSelector().Add(requirements ...))
	if err != nil {
		return tree.NewSelectServerTreeInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	for _, serverRow := range list {
		serverApp := serverRow.Spec.App
		serverName := serverRow.Spec.Server

		server := &models.ServerElem{
			ServerID:   util.GetServerId(serverApp, serverName),
			ServerName: serverName,
		}

		_, ok := allApp[serverApp]
		if !ok {
			allApp[serverApp] = make([]*models.ServerElem, 0, 10)
		}
		allApp[serverApp] = append(allApp[serverApp], server)
	}

	result := make([]*models.Business,0,10)

	emptyBusiness := &models.Business{
		BusinessName: "",
		BusinessShow: "",
		App:          make([]*models.AppElem, 0, 10),
	}

	tTree, err := K8sWatcher.tTreeLister.TTrees(namespace).Get(TafTreeName)
	if err != nil {
		return tree.NewSelectServerTreeInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	bussinessMap := make(map[string]models.Business, len(tTree.Businesses))
	for _, bn := range tTree.Businesses {
		bussinessMap[bn.Name] = models.Business{BusinessName: bn.Name, BusinessShow: bn.Show, App: make([]*models.AppElem, 0, len(tTree.Apps))}
	}

	for _, app := range tTree.Apps {
		appName := app.Name
		businessName := app.BusinessRef
		var businessShow string
		if buz, ok := bussinessMap[businessName]; ok {
			businessShow = buz.BusinessShow
		}

		for {
			if appName == "" {
				result = append(result, &models.Business{
					BusinessName: businessName,
					BusinessShow: businessShow,
					App:          make([]*models.AppElem, 0, 10),
				})
				break
			}

			if businessName == "" {
				if server, ok := allApp[appName]; ok {
					emptyBusiness.App = append(emptyBusiness.App, &models.AppElem{
						AppName: appName,
						Server:  server,
					})
				} else {
					emptyBusiness.App = append(emptyBusiness.App, &models.AppElem{
						AppName: appName,
						Server:  make([]*models.ServerElem, 0, 1),
					})
				}
				break
			}

			if lastBusiness != businessName {
				if targetBusiness != nil {
					result = append(result, targetBusiness)
				}
				lastBusiness = businessName
				targetBusiness = new(models.Business)
				targetBusiness.BusinessName = businessName
				targetBusiness.BusinessShow = businessShow
				targetBusiness.App = make([]*models.AppElem, 0, 10)
			}

			if serverSlice, ok := allApp[appName]; ok {
				targetBusiness.App = append(targetBusiness.App, &models.AppElem{
					AppName: appName,
					Server:  serverSlice,
				})
			} else {
				targetBusiness.App = append(targetBusiness.App, &models.AppElem{
					AppName: appName,
					Server:  make([]*models.ServerElem, 0, 1),
				})
			}
			break
		}
	}

	if targetBusiness != nil {
		result = append(result, targetBusiness)
	}

	if len(emptyBusiness.App) != 0 {
		result = append(result, emptyBusiness)
	}

	// 排序输出
	for i, business := range result {
		for j, app := range business.App {
			serverWrapper := ServerElemWrapper{Server: app.Server, By: func(e1, e2 *models.ServerElem) bool {
				return e1.ServerName < e2.ServerName
			}}
			sort.Sort(serverWrapper)
			business.App[j].Server = serverWrapper.Server
		}

		appWrapper := AppWrapper{App: business.App, By: func(e1, e2 *models.AppElem) bool {
			return e1.AppName < e2.AppName
		}}
		sort.Sort(appWrapper)
		result[i].App = appWrapper.App
	}

	return tree.NewSelectServerTreeOK().WithPayload(&tree.SelectServerTreeOKBody{Result: result})
}

