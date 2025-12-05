package clients

import (
	"fmt"
	"terraform-provider-ome/models"
)

func (c *Client) GetFabricByName(name string) (models.OMEFabric, error) {
	omeFabricResponse := []models.OMEFabric{}
	err := c.GetPaginatedDataWithQueryParam(FabricAPI, map[string]string{"$filter": fmt.Sprintf("%s eq '%s'", "Name", name)}, &omeFabricResponse)
	if err != nil {
		return models.OMEFabric{}, err
	}
	if len(omeFabricResponse) == 0 {
		return models.OMEFabric{}, nil
	}
	for _, omeFabric := range omeFabricResponse {
		if name == omeFabric.Name {
			return omeFabric, nil
		}
	}
	return models.OMEFabric{}, nil
}
