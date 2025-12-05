package ome

import (
	"context"
	"terraform-provider-ome/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &fabricDataSource{}
	_ datasource.DataSourceWithConfigure = &fabricDataSource{}
)

func NewFabricDataSource() datasource.DataSource {
	return &fabricDataSource{}
}

type fabricDataSource struct {
	p *omeProvider
}

func (f *fabricDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	f.p = req.ProviderData.(*omeProvider)
}

func (*fabricDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "fabric_info"
}

func (f fabricDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This Terraform DataSource is used to query fabric from OME." +
			" The information fetched from this data source can be used for getting the details / for further processing in resource block.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the fabric data source.",
				Description:         "ID of the fabric data source.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the fabric.",
				Description:         "Name of the fabric.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description for the fabric.",
				Description:         "Description for the fabric.",
				Computed:            true,
			},
		},
	}
}

func (f fabricDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var fabric models.FabricDataSource
	diags := req.Config.Get(ctx, &fabric)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	fabricName := fabric.Name.ValueString()
	omeClient, d := f.p.createOMESession(ctx, "datasource_fabric Read")
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}
	defer omeClient.RemoveSession()
	fabricData, err := omeClient.GetFabricByName(fabricName)
	if err != nil {
		resp.Diagnostics.AddError(
			"error reading the fabric", err.Error(),
		)
		return
	}
	updateFabricDataSourceState(&fabric, &fabricData)
	diags = resp.State.Set(ctx, fabric)
	resp.Diagnostics.Append(diags...)
}

func updateFabricDataSourceState(fabric *models.FabricDataSource, omeFabric *models.OMEFabric) {
	fabric.ID = types.StringValue(omeFabric.ID)
	fabric.Name = types.StringValue(omeFabric.Name)
	fabric.Description = types.StringValue(omeFabric.Description)
}
