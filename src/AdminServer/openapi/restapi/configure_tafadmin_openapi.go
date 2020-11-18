// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"tafadmin/openapi/restapi/operations"
	"tafadmin/openapi/restapi/operations/affinity"
	"tafadmin/openapi/restapi/operations/agent"
	"tafadmin/openapi/restapi/operations/applications"
	"tafadmin/openapi/restapi/operations/approval"
	"tafadmin/openapi/restapi/operations/business"
	"tafadmin/openapi/restapi/operations/config"
	"tafadmin/openapi/restapi/operations/default_operations"
	"tafadmin/openapi/restapi/operations/deploy"
	"tafadmin/openapi/restapi/operations/node"
	"tafadmin/openapi/restapi/operations/notify"
	"tafadmin/openapi/restapi/operations/release"
	"tafadmin/openapi/restapi/operations/server"
	"tafadmin/openapi/restapi/operations/server_k8s"
	"tafadmin/openapi/restapi/operations/server_option"
	"tafadmin/openapi/restapi/operations/server_pod"
	"tafadmin/openapi/restapi/operations/server_servant"
	"tafadmin/openapi/restapi/operations/shell"
	"tafadmin/openapi/restapi/operations/template"
	"tafadmin/openapi/restapi/operations/tree"
)

//go:generate swagger generate server --target ../../openapi --name TafadminOpenapi --spec ../../doc/Admin.yaml --principal interface{} --exclude-main

func configureFlags(api *operations.TafadminOpenapiAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.TafadminOpenapiAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.NotifySelectNotifyHandler == nil {
		api.NotifySelectNotifyHandler = notify.SelectNotifyHandlerFunc(func(params notify.SelectNotifyParams) middleware.Responder {
			return middleware.NotImplemented("operation notify.SelectNotify has not yet been implemented")
		})
	}
	if api.ServerPodSelectPodAliveHandler == nil {
		api.ServerPodSelectPodAliveHandler = server_pod.SelectPodAliveHandlerFunc(func(params server_pod.SelectPodAliveParams) middleware.Responder {
			return middleware.NotImplemented("operation server_pod.SelectPodAlive has not yet been implemented")
		})
	}
	if api.ServerPodSelectPodPerishedHandler == nil {
		api.ServerPodSelectPodPerishedHandler = server_pod.SelectPodPerishedHandlerFunc(func(params server_pod.SelectPodPerishedParams) middleware.Responder {
			return middleware.NotImplemented("operation server_pod.SelectPodPerished has not yet been implemented")
		})
	}
	if api.ApplicationsCreateAppHandler == nil {
		api.ApplicationsCreateAppHandler = applications.CreateAppHandlerFunc(func(params applications.CreateAppParams) middleware.Responder {
			return middleware.NotImplemented("operation applications.CreateApp has not yet been implemented")
		})
	}
	if api.ApprovalCreateApprovalHandler == nil {
		api.ApprovalCreateApprovalHandler = approval.CreateApprovalHandlerFunc(func(params approval.CreateApprovalParams) middleware.Responder {
			return middleware.NotImplemented("operation approval.CreateApproval has not yet been implemented")
		})
	}
	if api.BusinessCreateBusinessHandler == nil {
		api.BusinessCreateBusinessHandler = business.CreateBusinessHandlerFunc(func(params business.CreateBusinessParams) middleware.Responder {
			return middleware.NotImplemented("operation business.CreateBusiness has not yet been implemented")
		})
	}
	if api.DeployCreateDeployHandler == nil {
		api.DeployCreateDeployHandler = deploy.CreateDeployHandlerFunc(func(params deploy.CreateDeployParams) middleware.Responder {
			return middleware.NotImplemented("operation deploy.CreateDeploy has not yet been implemented")
		})
	}
	if api.ServerServantCreateServerAdapterHandler == nil {
		api.ServerServantCreateServerAdapterHandler = server_servant.CreateServerAdapterHandlerFunc(func(params server_servant.CreateServerAdapterParams) middleware.Responder {
			return middleware.NotImplemented("operation server_servant.CreateServerAdapter has not yet been implemented")
		})
	}
	if api.ConfigCreateServerConfigHandler == nil {
		api.ConfigCreateServerConfigHandler = config.CreateServerConfigHandlerFunc(func(params config.CreateServerConfigParams) middleware.Responder {
			return middleware.NotImplemented("operation config.CreateServerConfig has not yet been implemented")
		})
	}
	if api.ReleaseCreateServicePoolHandler == nil {
		api.ReleaseCreateServicePoolHandler = release.CreateServicePoolHandlerFunc(func(params release.CreateServicePoolParams) middleware.Responder {
			return middleware.NotImplemented("operation release.CreateServicePool has not yet been implemented")
		})
	}
	if api.TemplateCreateTemplateHandler == nil {
		api.TemplateCreateTemplateHandler = template.CreateTemplateHandlerFunc(func(params template.CreateTemplateParams) middleware.Responder {
			return middleware.NotImplemented("operation template.CreateTemplate has not yet been implemented")
		})
	}
	if api.ApplicationsDeleteAppHandler == nil {
		api.ApplicationsDeleteAppHandler = applications.DeleteAppHandlerFunc(func(params applications.DeleteAppParams) middleware.Responder {
			return middleware.NotImplemented("operation applications.DeleteApp has not yet been implemented")
		})
	}
	if api.BusinessDeleteBusinessHandler == nil {
		api.BusinessDeleteBusinessHandler = business.DeleteBusinessHandlerFunc(func(params business.DeleteBusinessParams) middleware.Responder {
			return middleware.NotImplemented("operation business.DeleteBusiness has not yet been implemented")
		})
	}
	if api.DeployDeleteDeployHandler == nil {
		api.DeployDeleteDeployHandler = deploy.DeleteDeployHandlerFunc(func(params deploy.DeleteDeployParams) middleware.Responder {
			return middleware.NotImplemented("operation deploy.DeleteDeploy has not yet been implemented")
		})
	}
	if api.ServerDeleteServerHandler == nil {
		api.ServerDeleteServerHandler = server.DeleteServerHandlerFunc(func(params server.DeleteServerParams) middleware.Responder {
			return middleware.NotImplemented("operation server.DeleteServer has not yet been implemented")
		})
	}
	if api.ServerServantDeleteServerAdapterHandler == nil {
		api.ServerServantDeleteServerAdapterHandler = server_servant.DeleteServerAdapterHandlerFunc(func(params server_servant.DeleteServerAdapterParams) middleware.Responder {
			return middleware.NotImplemented("operation server_servant.DeleteServerAdapter has not yet been implemented")
		})
	}
	if api.ConfigDeleteServerConfigHandler == nil {
		api.ConfigDeleteServerConfigHandler = config.DeleteServerConfigHandlerFunc(func(params config.DeleteServerConfigParams) middleware.Responder {
			return middleware.NotImplemented("operation config.DeleteServerConfig has not yet been implemented")
		})
	}
	if api.ConfigDeleteServerConfigHistoryHandler == nil {
		api.ConfigDeleteServerConfigHistoryHandler = config.DeleteServerConfigHistoryHandlerFunc(func(params config.DeleteServerConfigHistoryParams) middleware.Responder {
			return middleware.NotImplemented("operation config.DeleteServerConfigHistory has not yet been implemented")
		})
	}
	if api.TemplateDeleteTemplateHandler == nil {
		api.TemplateDeleteTemplateHandler = template.DeleteTemplateHandlerFunc(func(params template.DeleteTemplateParams) middleware.Responder {
			return middleware.NotImplemented("operation template.DeleteTemplate has not yet been implemented")
		})
	}
	if api.ConfigDoActiveHistoryConfigHandler == nil {
		api.ConfigDoActiveHistoryConfigHandler = config.DoActiveHistoryConfigHandlerFunc(func(params config.DoActiveHistoryConfigParams) middleware.Responder {
			return middleware.NotImplemented("operation config.DoActiveHistoryConfig has not yet been implemented")
		})
	}
	if api.BusinessDoAddBusinessAppHandler == nil {
		api.BusinessDoAddBusinessAppHandler = business.DoAddBusinessAppHandlerFunc(func(params business.DoAddBusinessAppParams) middleware.Responder {
			return middleware.NotImplemented("operation business.DoAddBusinessApp has not yet been implemented")
		})
	}
	if api.AffinityDoAddNodeEnableServerHandler == nil {
		api.AffinityDoAddNodeEnableServerHandler = affinity.DoAddNodeEnableServerHandlerFunc(func(params affinity.DoAddNodeEnableServerParams) middleware.Responder {
			return middleware.NotImplemented("operation affinity.DoAddNodeEnableServer has not yet been implemented")
		})
	}
	if api.AffinityDoAddServerEnableNodeHandler == nil {
		api.AffinityDoAddServerEnableNodeHandler = affinity.DoAddServerEnableNodeHandlerFunc(func(params affinity.DoAddServerEnableNodeParams) middleware.Responder {
			return middleware.NotImplemented("operation affinity.DoAddServerEnableNode has not yet been implemented")
		})
	}
	if api.AffinityDoDeleteNodeEnableServerHandler == nil {
		api.AffinityDoDeleteNodeEnableServerHandler = affinity.DoDeleteNodeEnableServerHandlerFunc(func(params affinity.DoDeleteNodeEnableServerParams) middleware.Responder {
			return middleware.NotImplemented("operation affinity.DoDeleteNodeEnableServer has not yet been implemented")
		})
	}
	if api.NodeDoDeletePublicNodeHandler == nil {
		api.NodeDoDeletePublicNodeHandler = node.DoDeletePublicNodeHandlerFunc(func(params node.DoDeletePublicNodeParams) middleware.Responder {
			return middleware.NotImplemented("operation node.DoDeletePublicNode has not yet been implemented")
		})
	}
	if api.AffinityDoDeleteServerEnableNodeHandler == nil {
		api.AffinityDoDeleteServerEnableNodeHandler = affinity.DoDeleteServerEnableNodeHandlerFunc(func(params affinity.DoDeleteServerEnableNodeParams) middleware.Responder {
			return middleware.NotImplemented("operation affinity.DoDeleteServerEnableNode has not yet been implemented")
		})
	}
	if api.ReleaseDoEnableServiceHandler == nil {
		api.ReleaseDoEnableServiceHandler = release.DoEnableServiceHandlerFunc(func(params release.DoEnableServiceParams) middleware.Responder {
			return middleware.NotImplemented("operation release.DoEnableService has not yet been implemented")
		})
	}
	if api.AffinityDoListAffinityGroupByAbilityHandler == nil {
		api.AffinityDoListAffinityGroupByAbilityHandler = affinity.DoListAffinityGroupByAbilityHandlerFunc(func(params affinity.DoListAffinityGroupByAbilityParams) middleware.Responder {
			return middleware.NotImplemented("operation affinity.DoListAffinityGroupByAbility has not yet been implemented")
		})
	}
	if api.AffinityDoListAffinityGroupByNodeHandler == nil {
		api.AffinityDoListAffinityGroupByNodeHandler = affinity.DoListAffinityGroupByNodeHandlerFunc(func(params affinity.DoListAffinityGroupByNodeParams) middleware.Responder {
			return middleware.NotImplemented("operation affinity.DoListAffinityGroupByNode has not yet been implemented")
		})
	}
	if api.BusinessDoListBusinessAppHandler == nil {
		api.BusinessDoListBusinessAppHandler = business.DoListBusinessAppHandlerFunc(func(params business.DoListBusinessAppParams) middleware.Responder {
			return middleware.NotImplemented("operation business.DoListBusinessApp has not yet been implemented")
		})
	}
	if api.NodeDoListClusterNodeHandler == nil {
		api.NodeDoListClusterNodeHandler = node.DoListClusterNodeHandlerFunc(func(params node.DoListClusterNodeParams) middleware.Responder {
			return middleware.NotImplemented("operation node.DoListClusterNode has not yet been implemented")
		})
	}
	if api.ConfigDoPreviewConfigContentHandler == nil {
		api.ConfigDoPreviewConfigContentHandler = config.DoPreviewConfigContentHandlerFunc(func(params config.DoPreviewConfigContentParams) middleware.Responder {
			return middleware.NotImplemented("operation config.DoPreviewConfigContent has not yet been implemented")
		})
	}
	if api.ServerOptionDoPreviewTemplateContentHandler == nil {
		api.ServerOptionDoPreviewTemplateContentHandler = server_option.DoPreviewTemplateContentHandlerFunc(func(params server_option.DoPreviewTemplateContentParams) middleware.Responder {
			return middleware.NotImplemented("operation server_option.DoPreviewTemplateContent has not yet been implemented")
		})
	}
	if api.NodeDoSetPublicNodeHandler == nil {
		api.NodeDoSetPublicNodeHandler = node.DoSetPublicNodeHandlerFunc(func(params node.DoSetPublicNodeParams) middleware.Responder {
			return middleware.NotImplemented("operation node.DoSetPublicNode has not yet been implemented")
		})
	}
	if api.ApplicationsSelectAppHandler == nil {
		api.ApplicationsSelectAppHandler = applications.SelectAppHandlerFunc(func(params applications.SelectAppParams) middleware.Responder {
			return middleware.NotImplemented("operation applications.SelectApp has not yet been implemented")
		})
	}
	if api.ApprovalSelectApprovalHandler == nil {
		api.ApprovalSelectApprovalHandler = approval.SelectApprovalHandlerFunc(func(params approval.SelectApprovalParams) middleware.Responder {
			return middleware.NotImplemented("operation approval.SelectApproval has not yet been implemented")
		})
	}
	if api.AgentSelectAvailHostPortHandler == nil {
		api.AgentSelectAvailHostPortHandler = agent.SelectAvailHostPortHandlerFunc(func(params agent.SelectAvailHostPortParams) middleware.Responder {
			return middleware.NotImplemented("operation agent.SelectAvailHostPort has not yet been implemented")
		})
	}
	if api.BusinessSelectBusinessHandler == nil {
		api.BusinessSelectBusinessHandler = business.SelectBusinessHandlerFunc(func(params business.SelectBusinessParams) middleware.Responder {
			return middleware.NotImplemented("operation business.SelectBusiness has not yet been implemented")
		})
	}
	if api.DefaultOperationsSelectDefaultValueHandler == nil {
		api.DefaultOperationsSelectDefaultValueHandler = default_operations.SelectDefaultValueHandlerFunc(func(params default_operations.SelectDefaultValueParams) middleware.Responder {
			return middleware.NotImplemented("operation default_operations.SelectDefaultValue has not yet been implemented")
		})
	}
	if api.DeploySelectDeployHandler == nil {
		api.DeploySelectDeployHandler = deploy.SelectDeployHandlerFunc(func(params deploy.SelectDeployParams) middleware.Responder {
			return middleware.NotImplemented("operation deploy.SelectDeploy has not yet been implemented")
		})
	}
	if api.ServerK8sSelectK8SHandler == nil {
		api.ServerK8sSelectK8SHandler = server_k8s.SelectK8SHandlerFunc(func(params server_k8s.SelectK8SParams) middleware.Responder {
			return middleware.NotImplemented("operation server_k8s.SelectK8S has not yet been implemented")
		})
	}
	if api.NodeSelectNodeHandler == nil {
		api.NodeSelectNodeHandler = node.SelectNodeHandlerFunc(func(params node.SelectNodeParams) middleware.Responder {
			return middleware.NotImplemented("operation node.SelectNode has not yet been implemented")
		})
	}
	if api.ServerSelectServerHandler == nil {
		api.ServerSelectServerHandler = server.SelectServerHandlerFunc(func(params server.SelectServerParams) middleware.Responder {
			return middleware.NotImplemented("operation server.SelectServer has not yet been implemented")
		})
	}
	if api.ServerServantSelectServerAdapterHandler == nil {
		api.ServerServantSelectServerAdapterHandler = server_servant.SelectServerAdapterHandlerFunc(func(params server_servant.SelectServerAdapterParams) middleware.Responder {
			return middleware.NotImplemented("operation server_servant.SelectServerAdapter has not yet been implemented")
		})
	}
	if api.ConfigSelectServerConfigHandler == nil {
		api.ConfigSelectServerConfigHandler = config.SelectServerConfigHandlerFunc(func(params config.SelectServerConfigParams) middleware.Responder {
			return middleware.NotImplemented("operation config.SelectServerConfig has not yet been implemented")
		})
	}
	if api.ConfigSelectServerConfigHistoryHandler == nil {
		api.ConfigSelectServerConfigHistoryHandler = config.SelectServerConfigHistoryHandlerFunc(func(params config.SelectServerConfigHistoryParams) middleware.Responder {
			return middleware.NotImplemented("operation config.SelectServerConfigHistory has not yet been implemented")
		})
	}
	if api.ServerOptionSelectServerOptionHandler == nil {
		api.ServerOptionSelectServerOptionHandler = server_option.SelectServerOptionHandlerFunc(func(params server_option.SelectServerOptionParams) middleware.Responder {
			return middleware.NotImplemented("operation server_option.SelectServerOption has not yet been implemented")
		})
	}
	if api.TreeSelectServerTreeHandler == nil {
		api.TreeSelectServerTreeHandler = tree.SelectServerTreeHandlerFunc(func(params tree.SelectServerTreeParams) middleware.Responder {
			return middleware.NotImplemented("operation tree.SelectServerTree has not yet been implemented")
		})
	}
	if api.ReleaseSelectServiceEnabledHandler == nil {
		api.ReleaseSelectServiceEnabledHandler = release.SelectServiceEnabledHandlerFunc(func(params release.SelectServiceEnabledParams) middleware.Responder {
			return middleware.NotImplemented("operation release.SelectServiceEnabled has not yet been implemented")
		})
	}
	if api.ReleaseSelectServicePoolHandler == nil {
		api.ReleaseSelectServicePoolHandler = release.SelectServicePoolHandlerFunc(func(params release.SelectServicePoolParams) middleware.Responder {
			return middleware.NotImplemented("operation release.SelectServicePool has not yet been implemented")
		})
	}
	if api.TemplateSelectTemplateHandler == nil {
		api.TemplateSelectTemplateHandler = template.SelectTemplateHandlerFunc(func(params template.SelectTemplateParams) middleware.Responder {
			return middleware.NotImplemented("operation template.SelectTemplate has not yet been implemented")
		})
	}
	if api.ShellSSHPodShellHandler == nil {
		api.ShellSSHPodShellHandler = shell.SSHPodShellHandlerFunc(func(params shell.SSHPodShellParams) middleware.Responder {
			return middleware.NotImplemented("operation shell.SSHPodShell has not yet been implemented")
		})
	}
	if api.ApplicationsUpdateAppHandler == nil {
		api.ApplicationsUpdateAppHandler = applications.UpdateAppHandlerFunc(func(params applications.UpdateAppParams) middleware.Responder {
			return middleware.NotImplemented("operation applications.UpdateApp has not yet been implemented")
		})
	}
	if api.BusinessUpdateBusinessHandler == nil {
		api.BusinessUpdateBusinessHandler = business.UpdateBusinessHandlerFunc(func(params business.UpdateBusinessParams) middleware.Responder {
			return middleware.NotImplemented("operation business.UpdateBusiness has not yet been implemented")
		})
	}
	if api.DeployUpdateDeployHandler == nil {
		api.DeployUpdateDeployHandler = deploy.UpdateDeployHandlerFunc(func(params deploy.UpdateDeployParams) middleware.Responder {
			return middleware.NotImplemented("operation deploy.UpdateDeploy has not yet been implemented")
		})
	}
	if api.ServerK8sUpdateK8SHandler == nil {
		api.ServerK8sUpdateK8SHandler = server_k8s.UpdateK8SHandlerFunc(func(params server_k8s.UpdateK8SParams) middleware.Responder {
			return middleware.NotImplemented("operation server_k8s.UpdateK8S has not yet been implemented")
		})
	}
	if api.ServerUpdateServerHandler == nil {
		api.ServerUpdateServerHandler = server.UpdateServerHandlerFunc(func(params server.UpdateServerParams) middleware.Responder {
			return middleware.NotImplemented("operation server.UpdateServer has not yet been implemented")
		})
	}
	if api.ServerServantUpdateServerAdapterHandler == nil {
		api.ServerServantUpdateServerAdapterHandler = server_servant.UpdateServerAdapterHandlerFunc(func(params server_servant.UpdateServerAdapterParams) middleware.Responder {
			return middleware.NotImplemented("operation server_servant.UpdateServerAdapter has not yet been implemented")
		})
	}
	if api.ConfigUpdateServerConfigHandler == nil {
		api.ConfigUpdateServerConfigHandler = config.UpdateServerConfigHandlerFunc(func(params config.UpdateServerConfigParams) middleware.Responder {
			return middleware.NotImplemented("operation config.UpdateServerConfig has not yet been implemented")
		})
	}
	if api.ServerOptionUpdateServerOptionHandler == nil {
		api.ServerOptionUpdateServerOptionHandler = server_option.UpdateServerOptionHandlerFunc(func(params server_option.UpdateServerOptionParams) middleware.Responder {
			return middleware.NotImplemented("operation server_option.UpdateServerOption has not yet been implemented")
		})
	}
	if api.TemplateUpdateTemplateHandler == nil {
		api.TemplateUpdateTemplateHandler = template.UpdateTemplateHandlerFunc(func(params template.UpdateTemplateParams) middleware.Responder {
			return middleware.NotImplemented("operation template.UpdateTemplate has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
