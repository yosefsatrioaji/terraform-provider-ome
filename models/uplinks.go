package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type UplinkDataSource struct {
	ID          types.String `tfsdk:"id"`
	FabricID    types.String `tfsdk:"fabric_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	MediaType   types.String `tfsdk:"media_type"`
	NativeVLAN  types.Int64  `tfsdk:"native_vlan"`
	Ports       types.List   `tfsdk:"ports"`
	Networks    types.List   `tfsdk:"networks"`
	UfdEnable   types.String `tfsdk:"ufd_enable"`
}

// type NetworkUplink struct {
// 	NetworkID   types.Int64  `tfsdk:"network_id"`
// 	NetworkName types.String `tfsdk:"network_name"`
// }

type PortUplink struct {
	PortID   types.String `tfsdk:"port_id"`
	PortName types.String `tfsdk:"port_name"`
}

//uplink get
type OMEUplink struct {
	ID                     string `json:"Id"`
	Name                   string `json:"Name"`
	Description            string `json:"Description"`
	MediaType              string `json:"MediaType"`
	NativeVLAN             int64  `json:"NativeVLAN"`
	PortsNavigationLink    string `json:"Ports@odata.navigationLink"`
	NetworksNavigationLink string `json:"Networks@odata.navigationLink"`
	UfdEnable              string `json:"UfdEnable"`
}

type OMEUplinkPort struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

type OMEUplinkPorts struct {
	Ports []OMEUplinkPort `json:"value"`
}

// type OMEUplinkNetwork struct {
// 	ID   string `json:"Id"`
// 	Name string `json:"Name"`
// }

type OMEUplinkNetworks struct {
	Networks []VLanNetworks `json:"value"`
}

type UplinkUpdate struct {
	ID          types.String `tfsdk:"id"`
	FabricID    types.String `tfsdk:"fabric_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	MediaType   types.String `tfsdk:"media_type"`
	NativeVLAN  types.Int64  `tfsdk:"native_vlan"`
	Ports       types.List   `tfsdk:"ports"`
	Networks    types.List   `tfsdk:"networks"`
	UfdEnable   types.String `tfsdk:"ufd_enable"`
}

type PortUplinkUpdate struct {
	ID types.String `tfsdk:"id"`
}

type NetworkUplinkUpdate struct {
	ID types.Int64 `tfsdk:"id"`
}

type OMEUplinkUpdate struct {
	ID          string                   `json:"Id"`
	Name        string                   `json:"Name"`
	Description string                   `json:"Description"`
	MediaType   string                   `json:"MediaType"`
	NativeVLAN  int64                    `json:"NativeVLAN"`
	Ports       []OMEPortUplinkUpdate    `json:"Ports"`
	Networks    []OMENetworkUplinkUpdate `json:"Networks"`
	UfdEnable   string                   `json:"UfdEnable"`
}

type OMEPortUplinkUpdate struct {
	ID string `json:"Id"`
}

type OMENetworkUplinkUpdate struct {
	ID int64 `json:"Id"`
}
