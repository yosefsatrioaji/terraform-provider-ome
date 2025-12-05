package ome

import (
	"context"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &uplinkDataSource{}
	_ datasource.DataSourceWithConfigure = &uplinkDataSource{}
)

func NewUplinkDataSource() datasource.DataSource {
	return &uplinkDataSource{}
}

type uplinkDataSource struct {
	p *omeProvider
}

func (u *uplinkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	u.p = req.ProviderData.(*omeProvider)
}

// Metadata implements datasource.DataSource
func (*uplinkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "uplink_info"
}

func (u uplinkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This Terraform DataSource is used to query uplink information from OME." +
			" The information fetched from this data source can be used for getting the details / for further processing in resource block.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the uplink data source.",
				Description:         "ID of the uplik data source.",
				Computed:            true,
			},
			"fabric_id": schema.StringAttribute{
				MarkdownDescription: "Fabric ID of the uplink.",
				Description:         "Fabric ID of the uplink.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the uplink.",
				Description:         "Name of the uplink.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the uplink.",
				Description:         "Description of the uplink.",
				Computed:            true,
			},
			"media_type": schema.StringAttribute{
				MarkdownDescription: "Media type of the uplink.",
				Description:         "Media type of the uplink.",
				Computed:            true,
			},
			"native_vlan": schema.Int64Attribute{
				MarkdownDescription: "Native VLAN of the uplink",
				Description:         "Native VLAN of the uplink",
				Computed:            true,
			},
			"ufd_enable": schema.StringAttribute{
				MarkdownDescription: "UFD enable of the uplink",
				Description:         "UFD enable of the uplink",
				Computed:            true,
			},
			"ports": schema.ListAttribute{
				MarkdownDescription: "Uplink ports",
				Description:         "Uplink ports",
				Computed:            true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"port_id":   types.StringType,
						"port_name": types.StringType,
					},
				},
			},
			"networks": schema.ListAttribute{
				MarkdownDescription: "Uplink networks",
				Description:         "Uplink networks",
				Computed:            true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"vlan_id": types.Int64Type,
						"name":    types.StringType,
					},
				},
			},
		},
	}
}

func (u uplinkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var uplink models.UplinkDataSource
	diags := req.Config.Get(ctx, &uplink)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	uplinkName := uplink.Name.ValueString()
	fabricId := uplink.FabricID.ValueString()

	omeClient, d := u.p.createOMESession(ctx, "datasource_uplink read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	omeUplinkData, err := omeClient.GetUplinkByName(fabricId, uplinkName)
	if err != nil {
		resp.Diagnostics.AddError(
			"error reading the uplink", err.Error(),
		)
		return
	}
	omePortsUplinkData, err := omeClient.GetUplinkPorts(fabricId, omeUplinkData.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to refresh ports of the uplink:",
			err.Error(),
		)
	}
	omeNetworksUplinkData, err := omeClient.GetUplinkNetworks(fabricId, omeUplinkData.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to refresh networks of the uplink:",
			err.Error(),
		)
	}
	updateUplinkDataSource(&uplink, &omeUplinkData, omePortsUplinkData, omeNetworksUplinkData)

	diags = resp.State.Set(ctx, uplink)
	resp.Diagnostics.Append(diags...)
}

func updateUplinkDataSource(uplink *models.UplinkDataSource, omeUplinkData *models.OMEUplink, omePorts models.OMEUplinkPorts, omeNetworks models.OMEUplinkNetworks) {
	uplink.ID = types.StringValue(omeUplinkData.ID)
	uplink.FabricID = types.StringValue(uplink.FabricID.ValueString())
	uplink.Name = types.StringValue(omeUplinkData.Name)
	uplink.Description = types.StringValue(omeUplinkData.Description)
	uplink.MediaType = types.StringValue(omeUplinkData.MediaType)
	uplink.NativeVLAN = types.Int64Value(int64(omeUplinkData.NativeVLAN))
	uplink.UfdEnable = types.StringValue(omeUplinkData.UfdEnable)
	portObjects := []attr.Value{}

	for _, port := range omePorts.Ports {
		portDetails := map[string]attr.Value{}
		portDetails["port_id"] = types.StringValue(port.ID)
		portDetails["port_name"] = types.StringValue(port.Name)
		portObject, _ := types.ObjectValue(
			map[string]attr.Type{
				"port_id":   types.StringType,
				"port_name": types.StringType,
			}, portDetails,
		)
		portObjects = append(portObjects, portObject)
	}
	portTfsdk, _ := types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"port_id":   types.StringType,
				"port_name": types.StringType,
			},
		},
		portObjects,
	)
	uplink.Ports = portTfsdk

	networkObjects := []attr.Value{}
	for _, network := range omeNetworks.Networks {
		networkDetails := map[string]attr.Value{}
		networkDetails["vlan_id"] = types.Int64Value(int64(network.ID))
		networkDetails["name"] = types.StringValue(network.Name)
		networkObject, _ := types.ObjectValue(
			map[string]attr.Type{
				"vlan_id": types.Int64Type,
				"name":    types.StringType,
			}, networkDetails,
		)
		networkObjects = append(networkObjects, networkObject)
	}
	networkTfsdk, _ := types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"vlan_id": types.Int64Type,
				"name":    types.StringType,
			},
		},
		networkObjects,
	)
	uplink.Networks = networkTfsdk
}
