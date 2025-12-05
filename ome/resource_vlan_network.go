package ome

import (
	"context"
	"reflect"
	"terraform-provider-ome/clients"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource = &vlanNetworkResource{}
)

func NewVlanNetworkResource() resource.Resource {
	return &vlanNetworkResource{}
}

type vlanNetworkResource struct {
	p *omeProvider
}

func (r *vlanNetworkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.p = req.ProviderData.(*omeProvider)
}

func (r *vlanNetworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "network_vlan"
}

func (r *vlanNetworkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This terraform resource is used to manage VLAN Network on OME." +
			"We can Create, Update and Delete OME VLAN Network using this resource.",
		Version: 1,
		Attributes: map[string]schema.Attribute{
			"vlan_id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "ID of the template resource.",
				Description:         "ID of the template resource.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the template resource.",
				Description:         "Name of the template resource.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the template resource.",
				Description:         "Description of the template resource.",
			},
			"vlan_maximum": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "VLAN ID Maximum",
				Description:         "VLAN ID Maximum",
			},
			"vlan_minimum": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "VLAN ID Minimum",
				Description:         "VLAN ID Minimum",
			},
			"type": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "VLAN Type",
				Description:         "VLAN Type",
			},
			"internal_ref_nwuu_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *vlanNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Trace(ctx, "resource_network_vlan create: started")
	var plan, state models.VLanNetworksTfsdk
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	omeClient, d := r.p.createOMESession(ctx, "resource_network_vlan Create")
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	vp := getVlanNetworkPayload(ctx, &plan)

	tflog.Trace(ctx, "resource_vlan_network create Creating VLAN Network")
	tflog.Debug(ctx, "resource_vlan_network create Creating VLAN Network", map[string]interface{}{
		"Create VLAN Network": vp,
	})

	cVlan, err := omeClient.CreateVlanNetwork(vp)
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrCreateVlanNetwork, err.Error(),
		)
		return
	}
	tflog.Trace(ctx, "resource_network_vlan create: create Finished creating VLAN Network")
	tflog.Trace(ctx, "resource_network_vlan create: updating state finished, saving ...")
	// Save into State
	state = saveVlanNetworkState(cVlan)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_network_vlan create: finish")
}

func (r *vlanNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Trace(ctx, "resource_network_vlan read: started")
	var state models.VLanNetworksTfsdk
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	omeClient, d := r.p.createOMESession(ctx, "resource_network_vlan Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()
	vlan, err := omeClient.GetVlanNetwork(state.VlanID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrReadVlanNetwork, err.Error(),
		)
		return
	}

	tflog.Trace(ctx, "resource_network_vlan read: read Finished reading VLAN Network")
	state = saveVlanNetworkState(vlan)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_network_vlan read: finish")
}

func (r *vlanNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Trace(ctx, "resource_network_vlan update: started")
	var state, plan models.VLanNetworksTfsdk
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	omeClient, d := r.p.createOMESession(ctx, "resource_network_vlan Update")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()

	if !reflect.DeepEqual(state, plan) {
		updatePayload := models.UpdateVlanNetwork{
			ID:          state.VlanID.ValueInt64(),
			Name:        plan.Name.ValueString(),
			Description: plan.Description.ValueString(),
			VLANMaximum: plan.VLANMaximum.ValueInt64(),
			VLANMinimum: plan.VLANMinimum.ValueInt64(),
			Type:        plan.Type.ValueInt64(),
		}
		vlan, err := omeClient.UpdateVlanNetwork(updatePayload)
		if err != nil {
			resp.Diagnostics.AddError(
				clients.ErrGnrUpdateVlanNetwork, err.Error(),
			)
			return
		}
		tflog.Trace(ctx, "resource_network_vlan update: update Finished updating VLAN Network")
		state = saveVlanNetworkState(vlan)
	}
	tflog.Trace(ctx, "resource_network_vlan update: updating state finished, saving ...")
	// Save into State
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "resource_network_vlan update: finish")
}

func (r *vlanNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Trace(ctx, "resource_network_vlan delete: started")
	var state models.VLanNetworksTfsdk
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	omeClient, d := r.p.createOMESession(ctx, "resource_network_vlan Delete")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()
	vlan, err := omeClient.DeleteVlanNetwork(state.VlanID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			clients.ErrGnrDeleteVlanNetwork, err.Error(),
		)
		return
	}
	resp.State.RemoveResource(ctx)
	tflog.Trace(ctx, "resource_network_vlan delete: delete Finished deleting VLAN Network")
	tflog.Trace(ctx, "resource_network_vlan delete: finished "+vlan)
}

func getVlanNetworkPayload(ctx context.Context, plan *models.VLanNetworksTfsdk) models.CreateVlanNetwork {
	vlan := models.CreateVlanNetwork{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		VLANMaximum: plan.VLANMaximum.ValueInt64(),
		VLANMinimum: plan.VLANMinimum.ValueInt64(),
		Type:        plan.Type.ValueInt64(),
	}
	return vlan
}

func saveVlanNetworkState(resp models.VLanNetworks) (state models.VLanNetworksTfsdk) {
	state.VlanID = types.Int64Value(resp.ID)
	state.Name = types.StringValue(resp.Name)
	state.Description = types.StringValue(resp.Description)
	state.VLANMaximum = types.Int64Value(resp.VLANMaximum)
	state.VLANMinimum = types.Int64Value(resp.VLANMinimum)
	state.Type = types.Int64Value(int64(resp.Type))
	state.InternalRefNWUUID = types.StringValue(resp.InternalRefNWUUID)
	return
}
