package k8s

import (
	"fmt"
	"tarsadmin/handler/util"
	"tarsadmin/openapi/models"
	"tarsadmin/openapi/restapi/operations/template"

	"github.com/go-openapi/runtime/middleware"
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	crdv1alpha1 "k8s.tars.io/api/crd/v1alpha1"
)

type CreateTemplateHandler struct{}

func (s *CreateTemplateHandler) Handle(params template.CreateTemplateParams) middleware.Responder {
	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	_, err := K8sWatcher.tTemplateLister.TTemplates(namespace).Get(*metadata.TemplateName)
	if err == nil {
		return template.NewCreateTemplateInternalServerError().WithPayload(&models.Error{Code: -1, Message: fmt.Sprintf("%s Already Existed.", metadata.TemplateName)})
	}

	if !errors.IsNotFound(err) {
		return template.NewCreateTemplateInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	tTemplate := buildTTemplate(namespace, *metadata.TemplateName, *metadata.TemplateParent, *metadata.TemplateContent)
	tTemplateInterface := K8sOption.CrdClientSet.CrdV1alpha1().TTemplates(namespace)
	if _, err = tTemplateInterface.Create(context.TODO(), tTemplate, k8sMetaV1.CreateOptions{}); err != nil && !errors.IsAlreadyExists(err) {
		return template.NewCreateTemplateInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return template.NewCreateTemplateOK().WithPayload(&template.CreateTemplateOKBody{Result: 0})
}

type SelectTemplateHandler struct{}

func (s *SelectTemplateHandler) Handle(params template.SelectTemplateParams) middleware.Responder {

	namespace := K8sOption.Namespace

	tTemplateInterface := K8sOption.CrdClientSet.CrdV1alpha1().TTemplates(namespace)
	list, err := tTemplateInterface.List(context.TODO(), k8sMetaV1.ListOptions{})
	if err != nil {
		return template.NewSelectTemplateInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}
	allItems := list.Items

	// parse query
	selectParams, err := util.ParseSelectQuery(params.Filter, params.Limiter, nil)
	if err != nil {
		return template.NewSelectTemplateInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	// filter
	filterItems := allItems
	// Admin临时版本，Template特化已有web实现
	if selectParams.Filter != nil {
		filterItems = make([]crdv1alpha1.TTemplate, 0, len(allItems))
		for _, elem := range allItems {
			if selectParams.Filter.Eq != nil {
				return template.NewSelectTemplateInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Eq Is Not Supported."})
			}
			if selectParams.Filter.Ne != nil {
				return template.NewSelectTemplateInternalServerError().WithPayload(&models.Error{Code: -1, Message: "Ne Is Not Supported."})
			}
			if selectParams.Filter.Like != nil {
				match, err := LikeMatch("TemplateName", selectParams.Filter.Like, elem.ObjectMeta.Name)
				if err != nil {
					return template.NewSelectTemplateInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
				}
				if !match {
					continue
				}
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
	result.Count["AllCount"] = int32(len(allItems))
	result.Count["FilterCount"] = int32(len(filterItems))

	// Data填充
	result.Data = make(models.ArrayMapInterface, 0, len(filterItems))
	for _, item := range filterItems {
		elem := make(map[string]interface{})
		elem["TemplateId"] = item.ObjectMeta.Name
		elem["TemplateName"] = item.ObjectMeta.Name
		elem["TemplateParent"] = item.Spec.Parent
		elem["TemplateContent"] = item.Spec.Content
		result.Data = append(result.Data, elem)
	}

	return template.NewSelectTemplateOK().WithPayload(result)
}

type UpdateTemplateHandler struct{}

func (s *UpdateTemplateHandler) Handle(params template.UpdateTemplateParams) middleware.Responder {

	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	tTemplate, err := K8sWatcher.tTemplateLister.TTemplates(namespace).Get(*metadata.TemplateID)
	if err != nil {
		return template.NewUpdateTemplateInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	target := params.Params.Target
	if equalTTemplate(tTemplate, target.TemplateParent, target.TemplateContent) {
		return template.NewUpdateTemplateInternalServerError().WithPayload(&models.Error{Code: -1, Message: fmt.Sprintf("%s Is Same.", *metadata.TemplateID)})
	}

	tTemplateCopy := tTemplate.DeepCopy()
	tTemplateCopy.Spec.Parent = target.TemplateParent
	tTemplateCopy.Spec.Content = target.TemplateContent

	tTemplateInterface := K8sOption.CrdClientSet.CrdV1alpha1().TTemplates(namespace)
	if _, err = tTemplateInterface.Update(context.TODO(), tTemplateCopy, k8sMetaV1.UpdateOptions{}); err != nil {
		return template.NewUpdateTemplateInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return template.NewUpdateTemplateOK().WithPayload(&template.UpdateTemplateOKBody{Result: 0})
}

type DeleteTemplateHandler struct{}

func (s *DeleteTemplateHandler) Handle(params template.DeleteTemplateParams) middleware.Responder {

	namespace := K8sOption.Namespace
	metadata := params.Params.Metadata

	_, err := K8sWatcher.tTemplateLister.TTemplates(namespace).Get(*metadata.TemplateID)
	if err != nil {
		return template.NewDeleteTemplateInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	tTemplateInterface := K8sOption.CrdClientSet.CrdV1alpha1().TTemplates(namespace)
	if err = tTemplateInterface.Delete(context.TODO(), *metadata.TemplateID, k8sMetaV1.DeleteOptions{}); err != nil {
		return template.NewDeleteTemplateInternalServerError().WithPayload(&models.Error{Code: -1, Message: err.Error()})
	}

	return template.NewDeleteTemplateOK().WithPayload(&template.DeleteTemplateOKBody{Result: 0})
}

func equalTTemplate(oldTemplate *crdv1alpha1.TTemplate, targetParent, targetContent string) bool {
	if oldTemplate.Spec.Parent != targetParent {
		return false
	}
	if oldTemplate.Spec.Content != targetContent {
		return false
	}
	return true
}

func buildTTemplate(namespace, templateName, templateParent, templateContent string) *crdv1alpha1.TTemplate {
	tTemplate := &crdv1alpha1.TTemplate{
		ObjectMeta: k8sMetaV1.ObjectMeta{
			Name:      templateName,
			Namespace: namespace,
		},
		Spec: crdv1alpha1.TTemplateSpec{
			Content: templateContent,
			Parent:  templateParent,
		},
	}
	return tTemplate
}
