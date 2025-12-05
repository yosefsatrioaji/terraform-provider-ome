package ome

import (
	"context"
	"fmt"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource = &uplinkUpdateResource{}
)

func NewUplinkUpdateResource() resource.Resource {
	return &uplinkUpdateResource{}
}

type uplinkUpdateResource struct {
	p *omeProvider
}

func (u *uplinkUpdateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	u.p = req.ProviderData.(*omeProvider)
}

func (u *uplinkUpdateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "uplink_update"
}

func (u *uplinkUpdateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This terraform resource is used to Update Uplink on OME." +
			"We can Update Uplink using this resource.",
		Version: 1,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of Uplink.",
				Description:         "ID of Uplink.",
			},
			"fabric_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the fabric.",
				Description:         "ID of the fabric.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the uplink.",
				Description:         "Name of the uplink.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the uplink.",
				Description:         "Description of the uplink.",
			},
			"media_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Media Type of the uplink.",
				Description:         "Media Type of the uplink.",
			},
			"native_vlan": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Native VLAN",
				Description:         "Native VLAN",
			},
			"ufd_enable": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "UFD Enable",
				Description:         "UFD Enable",
			},
			"ports": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "Ports",
				Description:         "Ports",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id": types.StringType,
					},
				},
			},
			"networks": schema.ListAttribute{
				Optional:            true,
				MarkdownDescription: "Networks",
				Description:         "Networks",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id": types.Int64Type,
					},
				},
			},
		},
	}
}

func (u *uplinkUpdateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "resource_uplink_update create: started")
	var plan models.UplinkUpdate
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	omeClient, d := u.p.createOMESession(ctx, "resource_uplink_update Create")
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer omeClient.RemoveSession()
	up := getUplinkPayload(ctx, &plan)
	tflog.Trace(ctx, "resource_uplink_update Create create Updating Uplink")
	tflog.Debug(ctx, "resource_uplink_update Create create Updating Uplink", map[string]interface{}{
		"Create Update Uplink": up,
	})
	tflog.Debug(ctx, fmt.Sprintf("%+v", up))
	err := omeClient.UpdateUplinkNetwork(plan.FabricID.ValueString(), up)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrUpdateUplink, err.Error(),
		)
		return
	}
	tflog.Trace(ctx, "resource_uplink_update create: create Finished Updating Uplink")
	tflog.Trace(ctx, "resource_uplink_update create: updating state finished, saving ...")

	OMEUplinkData, err := omeClient.GetUplinkByName(plan.FabricID.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"error reading the uplink", err.Error(),
		)
		return
	}
	omePortsUplinkData, err := omeClient.GetUplinkPorts(plan.FabricID.ValueString(), OMEUplinkData.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to refresh ports of the uplink:",
			err.Error(),
		)
	}
	omeNetworksUplinkData, err := omeClient.GetUplinkNetworks(plan.FabricID.ValueString(), OMEUplinkData.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to refresh networks of the uplink:",
			err.Error(),
		)
	}
	updateUplinkState(&plan, &OMEUplinkData, omePortsUplinkData, omeNetworksUplinkData)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (u *uplinkUpdateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "resource_uplink_update read: started")
	var uplink models.UplinkUpdate
	diags := req.State.Get(ctx, &uplink)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	omeClient, d := u.p.createOMESession(ctx, "resource_uplink_update Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()
	uplinkName := uplink.Name.ValueString()
	fabricId := uplink.FabricID.ValueString()
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
	updateUplinkState(&uplink, &omeUplinkData, omePortsUplinkData, omeNetworksUplinkData)

	diags = resp.State.Set(ctx, uplink)
	resp.Diagnostics.Append(diags...)
}

func (u *uplinkUpdateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "resource_uplink_update update: started")
	var plan models.UplinkUpdate
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	omeClient, d := u.p.createOMESession(ctx, "resource_uplink_update Update")
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer omeClient.RemoveSession()
	up := getUplinkPayload(ctx, &plan)
	tflog.Trace(ctx, "resource_uplink_update Create update Updating Uplink")
	tflog.Debug(ctx, "resource_uplink_update Create update Updating Uplink", map[string]interface{}{
		"Create Update Uplink": up,
	})
	err := omeClient.UpdateUplinkNetwork(plan.FabricID.ValueString(), up)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrUpdateUplink, err.Error(),
		)
		return
	}
	tflog.Trace(ctx, "resource_uplink_update create: update Finished Updating Uplink")
	tflog.Trace(ctx, "resource_uplink_update create: updating state finished, saving ...")

	OMEUplinkData, err := omeClient.GetUplinkByName(plan.FabricID.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"error reading the uplink", err.Error(),
		)
		return
	}
	omePortsUplinkData, err := omeClient.GetUplinkPorts(plan.FabricID.ValueString(), OMEUplinkData.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to refresh ports of the uplink:",
			err.Error(),
		)
	}
	omeNetworksUplinkData, err := omeClient.GetUplinkNetworks(plan.FabricID.ValueString(), OMEUplinkData.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to refresh networks of the uplink:",
			err.Error(),
		)
	}
	updateUplinkState(&plan, &OMEUplinkData, omePortsUplinkData, omeNetworksUplinkData)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (u *uplinkUpdateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "resource_uplink_update delete: started")
	tflog.Trace(ctx, "resource_template delete: finished")
}

func getUplinkPayload(ctx context.Context, plan *models.UplinkUpdate) models.OMEUplinkUpdate {
	planPortsObjects := []types.Object{}
	planNetworksObjects := []types.Object{}
	planPorts := []models.PortUplinkUpdate{}
	planNetworks := []models.NetworkUplinkUpdate{}
	OMEPort := []models.OMEPortUplinkUpdate{}
	OMENetwork := []models.OMENetworkUplinkUpdate{}
	plan.Ports.ElementsAs(ctx, &planPortsObjects, true)
	plan.Networks.ElementsAs(ctx, &planNetworksObjects, true)

	for _, planPortsObject := range planPortsObjects {
		port := models.PortUplinkUpdate{}
		planPortsObject.As(ctx, &port, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})
		planPorts = append(planPorts, port)
	}

	for _, planPort := range planPorts {
		updatePort := models.OMEPortUplinkUpdate{
			ID: planPort.ID.ValueString(),
		}
		OMEPort = append(OMEPort, updatePort)
	}

	for _, planNetworksObject := range planNetworksObjects {
		network := models.NetworkUplinkUpdate{}
		planNetworksObject.As(ctx, &network, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true})
		planNetworks = append(planNetworks, network)
	}

	for _, planNetwork := range planNetworks {
		updateNetwork := models.OMENetworkUplinkUpdate{
			ID: planNetwork.ID.ValueInt64(),
		}
		OMENetwork = append(OMENetwork, updateNetwork)
	}

	up := models.OMEUplinkUpdate{
		ID:          plan.ID.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		MediaType:   plan.MediaType.ValueString(),
		NativeVLAN:  plan.NativeVLAN.ValueInt64(),
		Ports:       OMEPort,
		Networks:    OMENetwork,
		UfdEnable:   plan.UfdEnable.ValueString(),
	}
	return up
}

func updateUplinkState(uplink *models.UplinkUpdate, omeUplinkData *models.OMEUplink, omePorts models.OMEUplinkPorts, omeNetworks models.OMEUplinkNetworks) {
	uplink.ID = types.StringValue(omeUplinkData.ID)
	uplink.FabricID = types.StringValue(uplink.FabricID.ValueString())
	uplink.Name = types.StringValue(omeUplinkData.Name)
	uplink.Description = types.StringValue(omeUplinkData.Description)
	uplink.MediaType = types.StringValue(omeUplinkData.MediaType)
	uplink.NativeVLAN = types.Int64Value(omeUplinkData.NativeVLAN)
	uplink.UfdEnable = types.StringValue(omeUplinkData.UfdEnable)

	portObjects := []attr.Value{}
	for _, port := range omePorts.Ports {
		portDetails := map[string]attr.Value{}
		portDetails["id"] = types.StringValue(port.ID)
		portObject, _ := types.ObjectValue(
			map[string]attr.Type{
				"id": types.StringType,
			}, portDetails,
		)
		portObjects = append(portObjects, portObject)
	}
	portTfsdk, _ := types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id": types.StringType,
			},
		},
		portObjects,
	)
	uplink.Ports = portTfsdk

	networkObjects := []attr.Value{}
	for _, network := range omeNetworks.Networks {
		networkDetails := map[string]attr.Value{}
		networkDetails["id"] = types.Int64Value(network.ID)
		networkObject, _ := types.ObjectValue(
			map[string]attr.Type{
				"id": types.Int64Type,
			}, networkDetails,
		)
		networkObjects = append(networkObjects, networkObject)
	}
	networkTfsdk, _ := types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id": types.Int64Type,
			},
		},
		networkObjects,
	)
	uplink.Networks = networkTfsdk
}
