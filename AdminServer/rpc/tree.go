package rpc

const _RequestKindServerTree = "ServerTree"

func selectServerTree(rpcRequest *Request) (Result, Error) {

	type ServerElem struct {
		ServerId   int    `json:"ServerId"`
		ServerName string `json:"ServerName"`
	}

	type AppElem struct {
		AppName string       `json:"AppName"`
		Server  []ServerElem `json:"Server"`
	}

	type Business struct {
		BusinessName string    `json:"BusinessName"`
		BusinessShow string    `json:"BusinessShow"`
		App          []AppElem `json:"App"`
	}

	var allApp = make(map[string][]ServerElem, 100)

	var lastApp string

	var posServer []ServerElem

	var targetBusiness *Business

	var lastBusiness string

	const sql1 = "select f_server_id ,f_server_app ,f_server_name from t_server order by f_server_app"
	serverRows, err := tafDb.Query(sql1)
	if err != nil {
		return nil, Error{"内部错误", -1}
	}

	defer func() {
		if serverRows != nil {
			_ = serverRows.Close()
		}
	}()

	for serverRows.Next() {
		var serverId int
		var serverApp string
		var serverName string
		if err := serverRows.Scan(&serverId, &serverApp, &serverName); err != nil {
			return nil, Error{"内部错误", -1}
		}

		server := ServerElem{
			ServerId:   serverId,
			ServerName: serverName,
		}

		if serverApp != lastApp {
			if posServer != nil {
				allApp[lastApp] = posServer
			}
			posServer = make([]ServerElem, 0, 10)
			lastApp = serverApp
		}
		posServer = append(posServer, server)
	}

	if posServer != nil {
		allApp[lastApp] = posServer
	}

	const sql2 = "select f_business_order bo, f_business_name bn, f_business_show, ifnull(f_app_name, '') an from t_business a left join t_app b using (f_business_name) union select -1, '', '', f_app_name an from t_app where f_business_name is null order by bo desc, bn, an"
	appRows, err := tafDb.Query(sql2)
	if err != nil {
		return nil, Error{"内部错误", -1}
	}

	defer func() {
		if appRows != nil {
			_ = appRows.Close()
		}
	}()

	result := make([]Business, 0, 20)

	emptyBusiness := Business{
		BusinessName: "",
		BusinessShow: "",
		App:          make([]AppElem, 0, 10),
	}

	for appRows.Next() {
		var businessOrder int
		var businessName string
		var businessShow string
		var appName string
		if err := appRows.Scan(&businessOrder, &businessName, &businessShow, &appName); err != nil {
			return nil, Error{"内部错误", -1}
		}

		for {
			if appName == "" {
				result = append(result, Business{
					BusinessName: businessName,
					BusinessShow: businessShow,
					App:          make([]AppElem, 0, 10),
				})
				break
			}

			if businessName == "" {
				if server, ok := allApp[appName]; ok {
					emptyBusiness.App = append(emptyBusiness.App, AppElem{
						AppName: appName,
						Server:  server,
					})
				} else {
					emptyBusiness.App = append(emptyBusiness.App, AppElem{
						AppName: appName,
						Server:  make([]ServerElem, 0),
					})
				}
				break
			}

			if lastBusiness != businessName {
				if targetBusiness != nil {
					result = append(result, *targetBusiness)
				}
				lastBusiness = businessName
				targetBusiness = new(Business)
				targetBusiness.BusinessName = businessName
				targetBusiness.BusinessShow = businessShow
				targetBusiness.App = make([]AppElem, 0, 10)
			}

			if serverSlice, ok := allApp[appName]; ok {
				targetBusiness.App = append(targetBusiness.App, AppElem{
					AppName: appName,
					Server:  serverSlice,
				})
			} else {
				targetBusiness.App = append(targetBusiness.App, AppElem{
					AppName: appName,
					Server:  make([]ServerElem, 0),
				})
			}
			break
		}
	}

	if targetBusiness != nil {
		result = append(result, *targetBusiness)
	}

	if len(emptyBusiness.App) != 0 {
		result = append(result, emptyBusiness)
	}

	return result, Error{"", Success}
}

func init() {
	registryHandle(RequestMethodSelect+_RequestKindServerTree, selectServerTree)
}
