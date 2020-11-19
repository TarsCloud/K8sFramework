// This file is safe to edit. Once it exists it will not be overwritten

package handler

import (
	"net/http"
	"tarsadmin/handler/compatible"
	"tarsadmin/handler/k8s"
	"tarsadmin/handler/mysql"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"tarsadmin/openapi/restapi/operations"
)

//go:generate swagger generate server --target ../../openapi --name TarsadminOpenapi --spec ../../doc/Admin.yaml --principal interface{} --exclude-main

type tarsAdminHandler struct {}

func (h *tarsAdminHandler) ConfigureAPI(api *operations.TarsadminOpenapiAPI) http.Handler {
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

	// -----------------------------

	api.ApplicationsCreateAppHandler = &k8s.CreateAppHandler{}

	api.ApprovalCreateApprovalHandler = &k8s.CreateApprovalHandler{}

	api.BusinessCreateBusinessHandler = &k8s.CreateBusinessHandler{}

	api.ServerServantCreateServerAdapterHandler = &k8s.CreateServerAdapterHandler{}

	api.ConfigCreateServerConfigHandler = &k8s.CreateServerConfigHandler{}

	api.DeployCreateDeployHandler = &k8s.CreateDeployHandler{}

	api.TemplateCreateTemplateHandler = &k8s.CreateTemplateHandler{}

	api.ReleaseCreateServicePoolHandler = &k8s.CreateServicePoolHandler{}

	// -----------------------------

	api.ApplicationsDeleteAppHandler = &k8s.DeleteAppHandler{}

	api.BusinessDeleteBusinessHandler = &k8s.DeleteBusinessHandler{}

	api.ServerDeleteServerHandler = &k8s.DeleteServerHandler{}

	api.ServerServantDeleteServerAdapterHandler = &k8s.DeleteServerAdapterHandler{}

	api.ConfigDeleteServerConfigHandler = &k8s.DeleteServerConfigHandler{}

	api.ConfigDeleteServerConfigHistoryHandler = &k8s.DeleteServerConfigHistoryHandler{}

	api.DeployDeleteDeployHandler = &k8s.DeleteDeployHandler{}

	api.TemplateDeleteTemplateHandler = &k8s.DeleteTemplateHandler{}

	// ------------------------------------

	api.ConfigDoActiveHistoryConfigHandler = &k8s.DoActiveHistoryConfigHandler{}

	api.BusinessDoAddBusinessAppHandler = &k8s.DoAddBusinessAppHandler{}

	api.ReleaseDoEnableServiceHandler = &k8s.DoEnableServiceHandler{}

	api.BusinessDoListBusinessAppHandler = &k8s.DoListBusinessAppHandler{}

	api.ConfigDoPreviewConfigContentHandler = &k8s.DoPreviewConfigContentHandler{}

	api.ServerOptionDoPreviewTemplateContentHandler = &k8s.DoPreviewTemplateContentHandler{}

	// ------------------------------------

	api.ApplicationsSelectAppHandler = &k8s.SelectAppHandler{}

	api.ApprovalSelectApprovalHandler = &k8s.SelectApprovalHandler{}

	api.BusinessSelectBusinessHandler = &k8s.SelectBusinessHandler{}

	api.DefaultOperationsSelectDefaultValueHandler = &k8s.SelectDefaultValueHandler{}

	api.ServerK8sSelectK8SHandler = &k8s.SelectServerK8SHandler{}

	api.ServerSelectServerHandler = &k8s.SelectServerHandler{}

	api.ServerServantSelectServerAdapterHandler = &k8s.SelectServerAdapterHandler{}

	api.ConfigSelectServerConfigHandler = &k8s.SelectServerConfigHandler{}

	api.ConfigSelectServerConfigHistoryHandler = &k8s.SelectServerConfigHistoryHandler{}

	api.ServerOptionSelectServerOptionHandler = &k8s.SelectServerOptionHandler{}

	api.TreeSelectServerTreeHandler = &k8s.SelectServerTreeHandler{}

	api.ReleaseSelectServiceEnabledHandler = &k8s.SelectServiceEnabledHandler{}

	api.ReleaseSelectServicePoolHandler = &k8s.SelectServicePoolHandler{}

	api.DeploySelectDeployHandler = &k8s.SelectDeployHandler{}

	api.TemplateSelectTemplateHandler = &k8s.SelectTemplateHandler{}

	api.NotifySelectNotifyHandler = &mysql.SelectNotifyHandler{}

	api.ServerPodSelectPodAliveHandler = &k8s.SelectPodAliveHandler{}

	api.ServerPodSelectPodPerishedHandler = &k8s.SelectPodPerishedHandler{}

	api.AgentSelectAvailHostPortHandler = &k8s.SelectAvailHostPortHandler{}

	// ------------------------------------

	api.ApplicationsUpdateAppHandler = &k8s.UpdateAppHandler{}

	api.BusinessUpdateBusinessHandler = &k8s.UpdateBusinessHandler{}

	api.ConfigUpdateServerConfigHandler = &k8s.UpdateServerConfigHandler{}

	api.ServerK8sUpdateK8SHandler = &k8s.UpdateServerK8SHandler{}

	api.ServerUpdateServerHandler = &k8s.UpdateServerHandler{}

	api.ServerServantUpdateServerAdapterHandler = &k8s.UpdateServerAdapterHandler{}

	api.ServerOptionUpdateServerOptionHandler = &k8s.UpdateServerOptionHandler{}

	api.DeployUpdateDeployHandler = &k8s.UpdateDeployHandler{}

	api.TemplateUpdateTemplateHandler = &k8s.UpdateTemplateHandler{}

	// -----------------------------

	api.NodeSelectNodeHandler = &compatible.SelectNodeHandler{}

	api.NodeDoListClusterNodeHandler = &compatible.DoListClusterNodeHandler{}

	api.NodeDoSetPublicNodeHandler = &compatible.DoSetPublicNodeHandler{}

	api.NodeDoDeletePublicNodeHandler = &compatible.DoDeletePublicNodeHandler{}

	// -----------------------------

	api.AffinityDoAddNodeEnableServerHandler = &compatible.DoAddNodeEnableServerHandler{}

	api.AffinityDoAddServerEnableNodeHandler = &compatible.DoAddServerEnableNodeHandler{}

	api.AffinityDoDeleteNodeEnableServerHandler = &compatible.DoDeleteNodeEnableServerHandler{}

	api.AffinityDoDeleteServerEnableNodeHandler = &compatible.DoDeleteServerEnableNodeHandler{}

	api.AffinityDoListAffinityGroupByNodeHandler = &compatible.DoListAffinityGroupByNodeHandler{}

	api.AffinityDoListAffinityGroupByAbilityHandler = &compatible.DoListAffinityGroupByAbilityHandler{}

	// -----------------------------

	api.ShellSSHPodShellHandler = &k8s.SSHPodShellHandler{}

	// -----------------------------

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return h.setupGlobalMiddleware(api.Serve(h.setupMiddlewares))
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func (h *tarsAdminHandler) setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func (h *tarsAdminHandler) setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
