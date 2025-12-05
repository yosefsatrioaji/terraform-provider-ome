/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://mozilla.org/MPL/2.0/
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package clients

import (
	"fmt"
	"terraform-provider-ome/models"
)

// GetAllVlanNetworks returns the vlan data from OME
func (c *Client) GetAllVlanNetworks() ([]models.VLanNetworks, error) {
	vlanData := []models.VLanNetworks{}
	err := c.GetPaginatedData(VlanNetworksAPI, &vlanData)
	if err != nil {
		return []models.VLanNetworks{}, err
	}
	return vlanData, nil
}

func (c *Client) GetVlanNetwork(id int64) (models.VLanNetworks, error) {
	vlanData := models.VLanNetworks{}
	fullPath := fmt.Sprintf(VlanNetworksAPI+"(%d)", id)
	response, err := c.Get(fullPath, nil, nil)
	if err != nil {
		return vlanData, err
	}
	respData, getBodyError := c.GetBodyData(response.Body)
	if getBodyError != nil {
		return vlanData, getBodyError
	}
	err = c.JSONUnMarshal(respData, &vlanData)
	return vlanData, err
}

func (c *Client) CreateVlanNetwork(vlan models.CreateVlanNetwork) (models.VLanNetworks, error) {
	data, errMarshal := c.JSONMarshal(vlan)
	if errMarshal != nil {
		return models.VLanNetworks{}, errMarshal
	}
	response, err := c.Post(VlanNetworksAPI, nil, data)
	if err != nil {
		return models.VLanNetworks{}, err
	}
	respData, getBodyError := c.GetBodyData(response.Body)
	if getBodyError != nil {
		return models.VLanNetworks{}, getBodyError
	}
	vlanData := models.VLanNetworks{}
	err = c.JSONUnMarshal(respData, &vlanData)
	if err != nil {
		return models.VLanNetworks{}, err
	}
	return vlanData, nil
}

func (c *Client) UpdateVlanNetwork(vlan models.UpdateVlanNetwork) (models.VLanNetworks, error) {
	omeVlan := models.VLanNetworks{}
	data, errMarshal := c.JSONMarshal(vlan)
	if errMarshal != nil {
		return omeVlan, errMarshal
	}
	fullPath := fmt.Sprintf(VlanNetworksAPI+"(%d)", vlan.ID)
	response, err := c.Put(fullPath, nil, data)
	if err != nil {
		return omeVlan, err
	}
	respData, getBodyError := c.GetBodyData(response.Body)
	if getBodyError != nil {
		return omeVlan, getBodyError
	}
	err = c.JSONUnMarshal(respData, &omeVlan)
	if err != nil {
		return omeVlan, err
	}
	return omeVlan, nil
}

func (c *Client) DeleteVlanNetwork(id int64) (string, error) {
	fullPath := fmt.Sprintf(VlanNetworksAPI+"(%d)", id)
	response, err := c.Delete(fullPath, nil, nil)
	if err != nil {
		return "", err
	}
	return response.Status, nil
}
