// Copyright (C) 2026 IOTech Ltd

package alarm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/clients/http"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/common"
	"github.com/IOTechSystems/go-mod-central-ext/v4/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/errors"
)

type Adapter struct {
	lc     logger.LoggingClient
	client *http.AlarmClient
}

func NewAdapter(lc logger.LoggingClient, client *http.AlarmClient) *Adapter {
	return &Adapter{lc: lc, client: client}
}

func (a *Adapter) AlarmConfigExists(ctx context.Context, data []byte) (bool, errors.EdgeX) {
	var config models.AlarmSetting
	if marshalErr := json.Unmarshal(data, &config); marshalErr != nil {
		return false, errors.NewCommonEdgeX(errors.Kind(marshalErr), "fail to unmarshal alarm config", marshalErr)
	}
	if len(strings.TrimSpace(config.Name)) == 0 {
		return false, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("alarm config '%s' missing or empty 'name' field", data), nil)
	}
	_, err := a.client.AlarmConfigByName(ctx, config.Name)
	if err == nil {
		a.lc.Debugf("alarm config '%s' already exists, skip adding", config.Name)
		return true, nil
	}
	return false, nil
}

func (a *Adapter) AddAlarmConfig(ctx context.Context, data []byte) errors.EdgeX {
	var config models.AlarmSetting
	if marshalErr := json.Unmarshal(data, &config); marshalErr != nil {
		return errors.NewCommonEdgeX(errors.Kind(marshalErr), "fail to unmarshal alarm config from %s", marshalErr)
	}
	if err := a.client.AddAlarmConfig(ctx, config.Name, data); err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("fail to add alarm config '%s'", config.Name), err)
	}
	return nil
}

func (a *Adapter) AssociationsExists(ctx context.Context, data []byte) (bool, errors.EdgeX) {
	var assoc models.AlarmAssociation
	if marshalErr := json.Unmarshal(data, &assoc); marshalErr != nil {
		return false, errors.NewCommonEdgeX(errors.Kind(marshalErr), "fail to unmarshal alarm association from %s", marshalErr)
	}

	queryParams, err := buildAssociationQueryParams(assoc)
	if err != nil {
		return false, errors.NewCommonEdgeXWrapper(err)
	}
	ass, err := a.client.Associations(ctx, queryParams)
	if err != nil {
		return false, errors.NewCommonEdgeXWrapper(err)
	}
	if ass.Metadata.Count > 0 {
		a.lc.Debugf("Association (sourceType=%s, configName=%s) already exists, skip adding", assoc.SourceType, assoc.ConfigName)
		return true, nil
	}
	return false, nil
}

func (a *Adapter) AddAssociations(ctx context.Context, data []byte) errors.EdgeX {
	var assoc models.AlarmAssociation
	if marshalErr := json.Unmarshal(data, &assoc); marshalErr != nil {
		return errors.NewCommonEdgeX(errors.Kind(marshalErr), "fail to unmarshal alarm association from %s", marshalErr)
	}
	switch assoc.SourceType {
	case common.AlarmSourceTypeDevice:
		return a.client.AddDeviceAssociation(ctx, assoc.DeviceName, assoc.ResourceName, assoc.ConfigName)
	case common.AlarmSourceTypeProfile:
		return a.client.AddProfileAssociation(ctx, assoc.ProfileName, assoc.ResourceName, assoc.ConfigName)
	case common.AlarmSourceTypeMessageBus:
		return a.client.AddMessageBusAssociation(ctx, assoc.MessageBusSourceName, assoc.ConfigName)
	case common.AlarmSourceTypeSparkplug:
		return a.client.AddSparkplugAssociation(ctx, assoc.SparkplugNodeId, assoc.SparkplugDeviceName, assoc.SparkplugMetricName, assoc.ConfigName)
	default:
		return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("unknown association sourceType: %s", assoc.SourceType), nil)
	}
}

func buildAssociationQueryParams(assoc models.AlarmAssociation) (map[string]string, errors.EdgeX) {
	params := map[string]string{
		"alarmConfigName": assoc.ConfigName,
		"sourceType":      assoc.SourceType,
	}

	switch assoc.SourceType {
	case common.AlarmSourceTypeDevice:
		params[common.AlarmSourceTypeDevice] = assoc.DeviceName
		params[common.AlarmAssociationResource] = assoc.ResourceName
	case common.AlarmSourceTypeProfile:
		params[common.AlarmSourceTypeProfile] = assoc.ProfileName
		params[common.AlarmAssociationResource] = assoc.ResourceName
	case common.AlarmSourceTypeMessageBus:
		params["messageBusSourceName"] = assoc.MessageBusSourceName
	case common.AlarmSourceTypeSparkplug:
		params["sparkplugNodeId"] = assoc.SparkplugNodeId
		params["sparkplugDeviceName"] = assoc.SparkplugDeviceName
		params["sparkplugMetricName"] = assoc.SparkplugMetricName
	default:
		return nil, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("unknown association sourceType: %s", assoc.SourceType), nil)
	}

	return params, nil
}

func (a *Adapter) TemplateExists(ctx context.Context, data []byte) (bool, errors.EdgeX) {
	return a.itemExistsFromList(ctx, data, a.client.TemplateByName)
}

func (a *Adapter) ConditionExists(ctx context.Context, data []byte) (bool, errors.EdgeX) {
	return a.itemExistsFromList(ctx, data, a.client.ConditionByName)
}

func (a *Adapter) ActionExists(ctx context.Context, data []byte) (bool, errors.EdgeX) {
	return a.itemExistsFromList(ctx, data, a.client.ActionByName)
}

func (a *Adapter) RouteExists(ctx context.Context, data []byte) (bool, errors.EdgeX) {
	return a.itemExistsFromList(ctx, data, a.client.RouteByName)
}

// listByNameFunc returns the result of a query by filtering the name,
// the name is a unique key in the schema, and the result will return only one item in the list
type listByNameFunc func(ctx context.Context, name string) (models.AlarmMultiResponse, errors.EdgeX)

func (a *Adapter) itemExistsFromList(ctx context.Context, data []byte, listByName listByNameFunc) (bool, errors.EdgeX) {
	var setting models.AlarmSetting
	if marshalErr := json.Unmarshal(data, &setting); marshalErr != nil {
		return false, errors.NewCommonEdgeX(errors.Kind(marshalErr), "fail to unmarshal alarm setting from %s", marshalErr)
	}
	if len(strings.TrimSpace(setting.Name)) == 0 {
		return false, errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("alarm setting '%s' missing or empty 'name' field", data), nil)
	}
	var res models.AlarmMultiResponse
	res, err := listByName(ctx, setting.Name)
	if err != nil {
		return false, errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("fail to query '%s'", setting.Name), err)
	}
	if res.Metadata.Count > 0 {
		a.lc.Debugf("The '%s' already exists", setting.Name)
		return true, nil
	}
	return false, nil
}

func (a *Adapter) AddTemplate(ctx context.Context, data []byte) errors.EdgeX {
	return a.client.AddTemplate(ctx, data)
}

func (a *Adapter) AddCondition(ctx context.Context, data []byte) errors.EdgeX {
	return a.client.AddCondition(ctx, data)
}

func (a *Adapter) AddAction(ctx context.Context, data []byte) errors.EdgeX {
	var action struct {
		Name         string
		TemplateName string
	}
	if err := json.Unmarshal(data, &action); err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("fail to unmarshal action data from %s", data), err)
	}

	// prepare the http post body, convert templateName to templateId
	var actionMap map[string]any
	if err := json.Unmarshal(data, &actionMap); err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("fail to unmarshal action data from %s", data), err)
	}
	if len(strings.TrimSpace(action.TemplateName)) > 0 { // templateName is optional, but if it is provided, it must be unique
		templateRes, err := a.client.TemplateByName(ctx, action.TemplateName)
		if err != nil {
			return errors.NewCommonEdgeXWrapper(err)
		}
		if templateRes.Metadata.Count == 0 {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("template '%s' not found", action.TemplateName), nil)
		}
		actionMap["templateId"] = templateRes.Templates[0].Id
		delete(actionMap, "templateName")
	}

	postData, marshalErr := json.Marshal(actionMap)
	if marshalErr != nil {
		return errors.NewCommonEdgeX(errors.Kind(marshalErr), fmt.Sprintf("fail to marshal action '%s'", action.Name), marshalErr)
	}

	if err := a.client.AddAction(ctx, postData); err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	a.lc.Debugf("successfully added action '%s'", action.Name)
	return nil
}

// ExportItem represents an exported alarm setting item.
type ExportItem struct {
	Name string // file name (without extension), e.g. "my-template"
	Data []byte // JSON content, can be written directly to a file
}

// listItems is a generic helper that lists all items from a given alarm API,
// extracts items by the responseKey, and returns them as ExportItem slice.
func (a *Adapter) listItems(ctx context.Context,
	listFn func(context.Context) (map[string]any, errors.EdgeX),
	responseKey string) ([]ExportItem, errors.EdgeX) {

	res, err := listFn(ctx)
	if err != nil {
		return nil, err
	}

	items, ok := res[responseKey].([]any)
	if !ok {
		a.lc.Debugf("no %s found to export", responseKey)
		return nil, nil
	}

	var result []ExportItem
	for _, item := range items {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		name, _ := itemMap[common.Name].(string)
		if name == "" {
			a.lc.Warnf("%s missing 'name', skipping export the data '%v'", responseKey, item)
			continue
		}

		data, marshalErr := json.MarshalIndent(itemMap, "", "  ")
		if marshalErr != nil {
			a.lc.Warnf("failed to marshal %s '%s': %v", responseKey, name, marshalErr)
			continue
		}
		result = append(result, ExportItem{Name: name, Data: data})
	}
	return result, nil
}

func (a *Adapter) ExportAlarmConfigs(ctx context.Context) ([]ExportItem, errors.EdgeX) {
	return a.listItems(ctx, a.client.AllAlarmConfigs, common.AlarmJsonKeyAlarmConfigs)
}

func (a *Adapter) ExportTemplates(ctx context.Context) ([]ExportItem, errors.EdgeX) {
	return a.listItems(ctx, a.client.AllTemplates, common.AlarmJsonKeyTemplates)
}

func (a *Adapter) ExportConditions(ctx context.Context) ([]ExportItem, errors.EdgeX) {
	return a.listItems(ctx, a.client.AllConditions, common.AlarmJsonKeyConditions)
}

func (a *Adapter) ExportActions(ctx context.Context) ([]ExportItem, errors.EdgeX) {
	res, err := a.client.AllActions(ctx)
	if err != nil {
		return nil, err
	}

	items, ok := res[common.AlarmJsonKeyActions].([]any)
	if !ok {
		a.lc.Debug("no actions found to export")
		return nil, nil
	}

	var result []ExportItem
	for _, item := range items {
		actionMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		name, _ := actionMap[common.Name].(string)
		if name == "" {
			continue
		}

		// convert templateId → templateName
		if templateId, ok := actionMap["templateId"].(string); ok && templateId != "" {
			templateRes, err := a.client.TemplateById(ctx, templateId)
			if err != nil {
				a.lc.Warnf("failed to resolve templateId '%s' for action '%s': %v", templateId, name, err)
			} else {
				if templateName, ok := templateRes[common.Name].(string); ok {
					actionMap["templateName"] = templateName
				}
			}
			delete(actionMap, "templateId")
		}

		data, marshalErr := json.MarshalIndent(actionMap, "", "  ")
		if marshalErr != nil {
			a.lc.Warnf("failed to marshal action '%s': %v", name, marshalErr)
			continue
		}
		result = append(result, ExportItem{Name: name, Data: data})
	}
	return result, nil
}

func (a *Adapter) ExportRoutes(ctx context.Context) ([]ExportItem, errors.EdgeX) {
	res, err := a.client.AllRoutes(ctx)
	if err != nil {
		return nil, err
	}

	items, ok := res[common.AlarmJsonKeyRoutes].([]any)
	if !ok {
		a.lc.Debug("no routes found to export")
		return nil, nil
	}

	var result []ExportItem
	for _, item := range items {
		routeMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		name, _ := routeMap[common.Name].(string)
		if name == "" {
			continue
		}

		// convert conditionId → conditionName
		if conditionId, ok := routeMap["conditionId"].(string); ok && conditionId != "" {
			conditionRes, err := a.client.ConditionById(ctx, conditionId)
			if err != nil {
				a.lc.Warnf("failed to resolve conditionId '%s' for route '%s': %v", conditionId, name, err)
			} else {
				if conditionName, ok := conditionRes[common.Name].(string); ok {
					routeMap["conditionName"] = conditionName
				}
			}
			delete(routeMap, "conditionId")
		}

		// convert actions (ID array) → actionNames (name array)
		if actionIds, ok := routeMap[common.AlarmJsonKeyActions].([]any); ok {
			var actionNames []string
			for _, aid := range actionIds {
				actionId, ok := aid.(string)
				if !ok {
					continue
				}
				actionRes, err := a.client.ActionById(ctx, actionId)
				if err != nil {
					a.lc.Warnf("failed to resolve actionId '%s' for route '%s': %v", actionId, name, err)
					continue
				}
				if actionName, ok := actionRes[common.Name].(string); ok {
					actionNames = append(actionNames, actionName)
				}
			}
			routeMap["actionNames"] = actionNames
			delete(routeMap, common.AlarmJsonKeyActions)
		}

		data, marshalErr := json.MarshalIndent(routeMap, "", "  ")
		if marshalErr != nil {
			a.lc.Warnf("failed to marshal route '%s': %v", name, marshalErr)
			continue
		}
		result = append(result, ExportItem{Name: name, Data: data})
	}
	return result, nil
}

func (a *Adapter) ExportAssociations(ctx context.Context) ([]ExportItem, errors.EdgeX) {
	res, err := a.client.AllAssociations(ctx)
	if err != nil {
		return nil, err
	}

	items, ok := res["associations"].([]any)
	if !ok {
		a.lc.Debug("no associations found to export")
		return nil, nil
	}

	var associations []models.AlarmAssociation
	for _, item := range items {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		sourceType, ok := itemMap["sourceType"].(string)
		if !ok {
			a.lc.Warn("association item missing 'sourceType', skipping")
			continue
		}
		source, ok := itemMap["source"].(map[string]any)
		if !ok {
			a.lc.Warnf("association item missing 'source' for sourceType '%s', skipping", sourceType)
			continue
		}
		configNames, ok := itemMap["alarmConfigNames"].([]any)
		if !ok {
			a.lc.Warnf("association item missing 'alarmConfigNames' for sourceType '%s', skipping", sourceType)
			continue
		}

		for _, cn := range configNames {
			configName, _ := cn.(string)
			if configName == "" {
				continue
			}
			assoc := models.AlarmAssociation{
				SourceType: sourceType,
				ConfigName: configName,
			}
			switch sourceType {
			case common.AlarmSourceTypeDevice:
				assoc.DeviceName, _ = source[common.AlarmSourceTypeDevice].(string)
				assoc.ResourceName, _ = source[common.AlarmAssociationResource].(string)
			case common.AlarmSourceTypeProfile:
				assoc.ProfileName, _ = source[common.AlarmSourceTypeProfile].(string)
				assoc.ResourceName, _ = source[common.AlarmAssociationResource].(string)
			case common.AlarmSourceTypeMessageBus:
				assoc.MessageBusSourceName, _ = source["messageBusSourceName"].(string)
			case common.AlarmSourceTypeSparkplug:
				assoc.SparkplugNodeId, _ = source["sparkplugNodeId"].(string)
				assoc.SparkplugDeviceName, _ = source["sparkplugDeviceName"].(string)
				assoc.SparkplugMetricName, _ = source["sparkplugMetricName"].(string)
			}
			associations = append(associations, assoc)
		}
	}

	if len(associations) == 0 {
		return nil, nil
	}

	data, marshalErr := json.MarshalIndent(associations, "", "  ")
	if marshalErr != nil {
		return nil, errors.NewCommonEdgeXWrapper(marshalErr)
	}
	return []ExportItem{{Name: "associations", Data: data}}, nil
}

func (a *Adapter) AddRoute(ctx context.Context, data []byte) errors.EdgeX {
	var route struct {
		Name          string
		ConditionName string
		ActionNames   []string
	}
	if err := json.Unmarshal(data, &route); err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("fail to unmarshal route data from %s", data), err)
	}

	// convert conditionName to conditionId
	res, err := a.client.ConditionByName(ctx, route.ConditionName) // name is the primary key of Condition, but the query result is multiple response
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	if res.Metadata.Count == 0 {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("condition '%s' not found", route.ConditionName), nil)
	}
	conditionId := res.Conditions[0].Id

	if len(route.ActionNames) == 0 {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("route %s missing or empty 'actionNames' field", data), nil)
	}
	// retrieve action IDs from names
	var actionIds []string
	for _, an := range route.ActionNames {
		if len(strings.TrimSpace(an)) == 0 {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("route %s has empty actionName entry", data), nil)
		}
		multiActionRes, err := a.client.ActionByName(ctx, an) // name is the primary key of Action, but the query result is multiple response
		if err != nil {
			return errors.NewCommonEdgeXWrapper(err)
		}
		if multiActionRes.Metadata.Count == 0 {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, fmt.Sprintf("action '%s' not found", an), nil)
		}
		actionIds = append(actionIds, multiActionRes.Actions[0].Id)
	}

	var routeMap map[string]any
	if err := json.Unmarshal(data, &routeMap); err != nil {
		return errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("fail to unmarshal route data from %s", data), err)
	}
	// ID required for REST API
	routeMap["conditionId"] = conditionId
	routeMap[common.AlarmJsonKeyActions] = actionIds
	postData, marshalErr := json.Marshal(routeMap)
	if marshalErr != nil {
		return errors.NewCommonEdgeX(errors.Kind(marshalErr), fmt.Sprintf("fail to marshal routeMap '%s'", route.Name), marshalErr)
	}

	if err := a.client.AddRoute(ctx, postData); err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	a.lc.Debugf("Successfully added route '%s'", route.Name)
	return nil
}
