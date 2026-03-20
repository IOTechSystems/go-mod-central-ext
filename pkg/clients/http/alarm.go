// Copyright (C) 2026 IOTech Ltd

package http

import (
	"context"
	"fmt"
	"net/url"

	pkgCommon "github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/http/utils"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

const defaultMaxLimit = "100"

// AlarmClient encapsulates HTTP operations against the support-alarm service.
type AlarmClient struct {
	baseUrl               string
	authInjector          interfaces.AuthenticationInjector
	enableNameFieldEscape bool
}

// NewAlarmClient creates a new AlarmClient for the given base URL.
func NewAlarmClient(baseUrl string, authInjector interfaces.AuthenticationInjector, enableNameFieldEscape bool) *AlarmClient {
	return &AlarmClient{
		baseUrl:               baseUrl,
		authInjector:          authInjector,
		enableNameFieldEscape: enableNameFieldEscape,
	}
}

// AlarmConfigByName checks whether an alarm config with the given name already exists.
func (c *AlarmClient) AlarmConfigByName(ctx context.Context, name string) (map[string]any, errors.EdgeX) {
	requestPath := common.NewPathBuilder().EnableNameFieldEscape(c.enableNameFieldEscape).
		SetPath(pkgCommon.AlarmConfigAPIRoute).SetNameFieldPath(name).BuildPath()
	var res map[string]any
	err := utils.GetRequest(ctx, &res, c.baseUrl, requestPath, nil, c.authInjector)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.Kind(err), "fail to query alarm config", err)
	}
	return res, nil
}

// AddAlarmConfig creates an alarm config with the given name and JSON data.
func (c *AlarmClient) AddAlarmConfig(ctx context.Context, name string, data []byte) errors.EdgeX {
	requestPath := common.NewPathBuilder().EnableNameFieldEscape(c.enableNameFieldEscape).
		SetPath(pkgCommon.AlarmConfigAPIRoute).SetNameFieldPath(name).BuildPath()
	var res map[string]any
	err := utils.PostRequest(ctx, &res, c.baseUrl, requestPath, data, common.ContentTypeJSON, c.authInjector)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), "fail to add alarm config", err)
	}
	return nil
}

// Associations query associations with the given query parameters.
func (c *AlarmClient) Associations(ctx context.Context, queryParams map[string]string) (models.AlarmMultiAssociationResponse, errors.EdgeX) {
	params := url.Values{}
	for k, v := range queryParams {
		params.Set(k, v)
	}
	var res models.AlarmMultiAssociationResponse
	err := utils.GetRequest(ctx, &res, c.baseUrl, pkgCommon.AssociationQueryAPIRoute, params, c.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("fail to query association with query parameters '%v'", queryParams), err)
	}
	return res, nil
}

// AddDeviceAssociation creates an edgexDevice association.
func (c *AlarmClient) AddDeviceAssociation(ctx context.Context, deviceName, resourceName, configName string) errors.EdgeX {
	requestPath := common.NewPathBuilder().EnableNameFieldEscape(c.enableNameFieldEscape).
		SetPath(pkgCommon.AssociationAPIRoute).SetPath(pkgCommon.AlarmAssociationEdgex).SetPath("device").SetNameFieldPath(deviceName).
		SetPath(pkgCommon.AlarmAssociationResource).SetNameFieldPath(resourceName).
		SetPath(pkgCommon.AlarmAssociationConfigName).SetNameFieldPath(configName).BuildPath()
	return c.postAssociation(ctx, requestPath)
}

// AddProfileAssociation creates an edgexProfile association.
func (c *AlarmClient) AddProfileAssociation(ctx context.Context, profileName, resourceName, configName string) errors.EdgeX {
	requestPath := common.NewPathBuilder().EnableNameFieldEscape(c.enableNameFieldEscape).
		SetPath(pkgCommon.AssociationAPIRoute).SetPath(pkgCommon.AlarmAssociationEdgex).SetPath("profile").SetNameFieldPath(profileName).
		SetPath(pkgCommon.AlarmAssociationResource).SetNameFieldPath(resourceName).
		SetPath(pkgCommon.AlarmAssociationConfigName).SetNameFieldPath(configName).BuildPath()
	return c.postAssociation(ctx, requestPath)
}

// AddMessageBusAssociation creates a messageBus association.
func (c *AlarmClient) AddMessageBusAssociation(ctx context.Context, messageBusSourceName, configName string) errors.EdgeX {
	requestPath := common.NewPathBuilder().EnableNameFieldEscape(c.enableNameFieldEscape).
		SetPath(pkgCommon.AssociationAPIRoute).SetPath("messagebus").SetPath("name").SetNameFieldPath(messageBusSourceName).
		SetPath(pkgCommon.AlarmAssociationConfigName).SetNameFieldPath(configName).BuildPath()
	return c.postAssociation(ctx, requestPath)
}

// AddSparkplugAssociation creates a sparkplug association.
func (c *AlarmClient) AddSparkplugAssociation(ctx context.Context, nodeId, deviceName, metricName, configName string) errors.EdgeX {
	requestPath := common.NewPathBuilder().EnableNameFieldEscape(c.enableNameFieldEscape).
		SetPath(pkgCommon.AssociationAPIRoute).SetPath("sparkplug").SetPath("node").SetNameFieldPath(nodeId).
		SetPath("device").SetNameFieldPath(deviceName).
		SetPath("metric").SetNameFieldPath(metricName).
		SetPath(pkgCommon.AlarmAssociationConfigName).SetNameFieldPath(configName).BuildPath()
	return c.postAssociation(ctx, requestPath)
}

func (c *AlarmClient) postAssociation(ctx context.Context, requestPath string) errors.EdgeX {
	var res map[string]any
	err := utils.PostRequest(ctx, &res, c.baseUrl, requestPath, nil, "", c.authInjector)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), "fail to create association", err)
	}
	return nil
}

// TemplateByName queries a template by name.
func (c *AlarmClient) TemplateByName(ctx context.Context, name string) (models.AlarmMultiResponse, errors.EdgeX) {
	return c.queryByName(ctx, pkgCommon.AlarmTemplateAPIRoute, name)
}

// ConditionByName queries a condition by name.
func (c *AlarmClient) ConditionByName(ctx context.Context, name string) (models.AlarmMultiResponse, errors.EdgeX) {
	return c.queryByName(ctx, pkgCommon.AlarmConditionAPIRoute, name)
}

// ActionByName queries an action by name.
func (c *AlarmClient) ActionByName(ctx context.Context, name string) (models.AlarmMultiResponse, errors.EdgeX) {
	return c.queryByName(ctx, pkgCommon.AlarmActionAPIRoute, name)
}

// RouteByName queries a route by name.
func (c *AlarmClient) RouteByName(ctx context.Context, name string) (models.AlarmMultiResponse, errors.EdgeX) {
	return c.queryByName(ctx, pkgCommon.AlarmRouteAPIRoute, name)
}

// AddTemplate creates a template resource.
func (c *AlarmClient) AddTemplate(ctx context.Context, data []byte) errors.EdgeX {
	var res map[string]any
	err := utils.PostRequest(ctx, &res, c.baseUrl, pkgCommon.AlarmTemplateAPIRoute, data, common.ContentTypeJSON, c.authInjector)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("fail to create template at %s", pkgCommon.AlarmTemplateAPIRoute), err)
	}
	return nil
}

// AddCondition creates a condition resource.
func (c *AlarmClient) AddCondition(ctx context.Context, data []byte) errors.EdgeX {
	var res map[string]any
	err := utils.PostRequest(ctx, &res, c.baseUrl, pkgCommon.AlarmConditionAPIRoute, data, common.ContentTypeJSON, c.authInjector)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("fail to create condition at %s", pkgCommon.AlarmConditionAPIRoute), err)
	}
	return nil
}

// AddAction creates an action resource.
func (c *AlarmClient) AddAction(ctx context.Context, data []byte) errors.EdgeX {
	var res map[string]any
	err := utils.PostRequest(ctx, &res, c.baseUrl, pkgCommon.AlarmActionAPIRoute, data, common.ContentTypeJSON, c.authInjector)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("fail to create action at %s", pkgCommon.AlarmActionAPIRoute), err)
	}
	return nil
}

// AddRoute creates a route resource.
func (c *AlarmClient) AddRoute(ctx context.Context, data []byte) errors.EdgeX {
	var res map[string]any
	err := utils.PostRequest(ctx, &res, c.baseUrl, pkgCommon.AlarmRouteAPIRoute, data, common.ContentTypeJSON, c.authInjector)
	if err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("fail to create route at %s", pkgCommon.AlarmRouteAPIRoute), err)
	}
	return nil
}

// queryAll sends GET {apiRoute}?offset=0&limit=-1 and returns the full response map.
func (c *AlarmClient) queryAll(ctx context.Context, apiRoute string) (map[string]any, errors.EdgeX) {
	params := url.Values{}
	params.Set(pkgCommon.Offset, "0")
	params.Set(pkgCommon.Limit, defaultMaxLimit) // TODO use -1 if the alarm service supports it
	var res map[string]any
	err := utils.GetRequest(ctx, &res, c.baseUrl, apiRoute, params, c.authInjector)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("fail to query all from %s", apiRoute), err)
	}
	return res, nil
}

// AllAlarmConfigs lists all alarm configs.
func (c *AlarmClient) AllAlarmConfigs(ctx context.Context) (map[string]any, errors.EdgeX) {
	return c.queryAll(ctx, pkgCommon.AlarmConfigsListAPIRoute)
}

// AllAssociations lists all associations.
func (c *AlarmClient) AllAssociations(ctx context.Context) (map[string]any, errors.EdgeX) {
	return c.queryAll(ctx, pkgCommon.AssociationQueryAPIRoute)
}

// AllTemplates lists all templates.
func (c *AlarmClient) AllTemplates(ctx context.Context) (map[string]any, errors.EdgeX) {
	return c.queryAll(ctx, pkgCommon.AlarmTemplateAPIRoute)
}

// AllConditions lists all conditions.
func (c *AlarmClient) AllConditions(ctx context.Context) (map[string]any, errors.EdgeX) {
	return c.queryAll(ctx, pkgCommon.AlarmConditionAPIRoute)
}

// AllActions lists all actions.
func (c *AlarmClient) AllActions(ctx context.Context) (map[string]any, errors.EdgeX) {
	return c.queryAll(ctx, pkgCommon.AlarmActionAPIRoute)
}

// AllRoutes lists all routes.
func (c *AlarmClient) AllRoutes(ctx context.Context) (map[string]any, errors.EdgeX) {
	return c.queryAll(ctx, pkgCommon.AlarmRouteAPIRoute)
}

// queryById sends GET {apiRoute}/{id} and returns the response map.
func (c *AlarmClient) queryById(ctx context.Context, apiRoute, id string) (map[string]any, errors.EdgeX) {
	requestPath := fmt.Sprintf("%s/%s", apiRoute, id)
	var res map[string]any
	err := utils.GetRequest(ctx, &res, c.baseUrl, requestPath, nil, c.authInjector)
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("fail to query %s/%s", apiRoute, id), err)
	}
	return res, nil
}

// TemplateById queries a template by ID.
func (c *AlarmClient) TemplateById(ctx context.Context, id string) (map[string]any, errors.EdgeX) {
	return c.queryById(ctx, pkgCommon.AlarmTemplateByIdRoute, id)
}

// ConditionById queries a condition by ID.
func (c *AlarmClient) ConditionById(ctx context.Context, id string) (map[string]any, errors.EdgeX) {
	return c.queryById(ctx, pkgCommon.AlarmConditionByIdRoute, id)
}

// ActionById queries an action by ID.
func (c *AlarmClient) ActionById(ctx context.Context, id string) (map[string]any, errors.EdgeX) {
	return c.queryById(ctx, pkgCommon.AlarmActionByIdRoute, id)
}

// queryByName sends GET {apiRoute}?name={name}&limit=1 and parses the response.
// Returns (id, true, nil) if found, ("", false, nil) if not found, ("", false, err) on error.
func (c *AlarmClient) queryByName(ctx context.Context, apiRoute, name string) (models.AlarmMultiResponse, errors.EdgeX) {
	params := url.Values{}
	params.Set(pkgCommon.Name, name)
	params.Set(pkgCommon.Limit, "1")

	var res models.AlarmMultiResponse
	err := utils.GetRequest(ctx, &res, c.baseUrl, apiRoute, params, c.authInjector)
	if err != nil {
		return res, errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("fail to query %s '%s'", apiRoute, name), err)
	}

	return res, nil
}
