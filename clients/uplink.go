package clients

import (
	"fmt"
	"terraform-provider-ome/models"
)

func (c *Client) GetUplinkByName(fabricID string, name string) (models.OMEUplink, error) {
	omeUplinkResponse := []models.OMEUplink{}
	// err := c.GetPaginatedDataWithQueryParam(fmt.Sprintf(UplinkAPI, fabricID), nil, &omeUplinkResponse)
	err := c.GetPaginatedData(fmt.Sprintf(UplinkAPI, fabricID), &omeUplinkResponse)
	// response, err := c.Get(fmt.Sprintf(UplinkAPI, fabricID), nil, nil)
	if err != nil {
		return models.OMEUplink{}, err
	}
	if len(omeUplinkResponse) == 0 {
		return models.OMEUplink{}, nil
	}
	for _, u := range omeUplinkResponse {
		if u.Name == name {
			return u, nil
		}
	}
	return models.OMEUplink{}, nil
}

func (c *Client) GetUplinkPorts(fabricID string, uplinkID string) (models.OMEUplinkPorts, error) {
	omeUplinkPortResponse := models.OMEUplinkPorts{}
	resp, err := c.Get(fmt.Sprintf(UplinkAPI+"('%s')/Ports", fabricID, uplinkID), nil, nil)
	if err != nil {
		return models.OMEUplinkPorts{}, err
	}
	respBody, errorBody := c.GetBodyData(resp.Body)
	if errorBody != nil {
		return models.OMEUplinkPorts{}, errorBody
	}
	err = c.JSONUnMarshal(respBody, &omeUplinkPortResponse)
	if err != nil {
		return models.OMEUplinkPorts{}, err
	}
	return omeUplinkPortResponse, nil
}

func (c *Client) GetUplinkNetworks(fabricID string, uplinkID string) (models.OMEUplinkNetworks, error) {
	omeUplinkNetworkResponse := models.OMEUplinkNetworks{}
	resp, err := c.Get(fmt.Sprintf(UplinkAPI+"('%s')/Networks", fabricID, uplinkID), nil, nil)
	if err != nil {
		return models.OMEUplinkNetworks{}, err
	}
	respBody, errorBody := c.GetBodyData(resp.Body)
	if errorBody != nil {
		return models.OMEUplinkNetworks{}, errorBody
	}
	err = c.JSONUnMarshal(respBody, &omeUplinkNetworkResponse)
	if err != nil {
		return models.OMEUplinkNetworks{}, err
	}
	return omeUplinkNetworkResponse, nil
}

func (c *Client) UpdateUplinkNetwork(fabricID string, uplink models.OMEUplinkUpdate) error {
	data, errMarshal := c.JSONMarshal(uplink)
	if errMarshal != nil {
		return errMarshal
	}
	fullPath := fmt.Sprintf(UplinkAPI+"('%s')", fabricID, uplink.ID)
	_, err := c.Put(fullPath, nil, data)
	return err
}
