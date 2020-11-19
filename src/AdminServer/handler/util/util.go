package util

import (
	"encoding/json"
	"fmt"
	"strings"
	"tarsadmin/openapi/models"
)

type SelectParams struct {
	Filter       *models.SelectRequestFilter  	`json:"filter"`
	Order        *models.SelectRequestOrder   	`json:"order"`
	Limiter      *models.SelectRequestLimiter 	`json:"limiter"`
	Confirmation bool `json:"confirmation"`
}

func ParseSelectQuery(Filter, Limiter, Order *string) (*SelectParams, error) {
	selectParams := &SelectParams{}

	if Filter != nil {
		var filter models.SelectRequestFilter
		if err := filter.UnmarshalBinary([]byte(*Filter)); err != nil {
			return nil, err
		}
		selectParams.Filter = &filter
	}

	if Limiter != nil {
		var limiter models.SelectRequestLimiter
		if err := limiter.UnmarshalBinary([]byte(*Limiter)); err != nil {
			return nil, err
		}
		selectParams.Limiter = &limiter
	}

	if Order != nil {
		var order models.SelectRequestOrder
		if err := json.Unmarshal([]byte(*Order), &order); err != nil {
			return nil, err
		}
		selectParams.Order = &order
	}

	return selectParams, nil
}

func GetServerId(serverApp, serverName string) string {
	return fmt.Sprintf("%s.%s", serverApp, serverName)
}

func GetTServerName(serverId string) string {
	return strings.ToLower(strings.ReplaceAll(serverId, ".", "-"))
}

