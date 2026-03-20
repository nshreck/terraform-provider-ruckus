package provider

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WLANGroupDS struct{ client *APIClient }

type WLANGroupDSItem struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Members     types.List   `tfsdk:"members"`
}

type WLANGroupDSModel struct {
	ZoneID types.String `tfsdk:"zone_id"`
	Groups types.List   `tfsdk:"groups"`
}

func NewWLANGroupDataSource() datasource.DataSource { return &WLANGroupDS{} }

func (d *WLANGroupDS) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "ruckus_wlan_groups"
}

func (d *WLANGroupDS) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a list of WLAN Groups in a specified zone.",
		Attributes: map[string]schema.Attribute{
			"zone_id": schema.StringAttribute{
				Description: "ID of the zone to list WLAN Groups from.",
				Required:    true,
			},
			"groups": schema.ListNestedAttribute{
				Description: "List of WLAN Groups.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the WLAN Group.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the WLAN Group.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the WLAN Group.",
							Computed:    true,
						},
						"members": schema.ListAttribute{
							Description: "List of WLAN IDs that are members of this group.",
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *WLANGroupDS) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		d.client = req.ProviderData.(*APIClient)
	}
}

type wlanGroupListResp struct {
	List []wlanGroupListItem `json:"list"`
}

type wlanGroupListItem struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Members     []wlanGroupMember `json:"members"`
}

func (d *WLANGroupDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("not configured", "missing API client")
		return
	}
	var state WLANGroupDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GET {base}/wsg/api/public/{ver}/rkszones/{zoneID}/wlangroups?serviceTicket=...
	q := url.Values{}
	q.Set("serviceTicket", d.client.ServiceTicket)

	endpoint := fmt.Sprintf("%s/wsg/api/public/%s/rkszones/%s/wlangroups?%s",
		d.client.BaseURL, d.client.APIVersion, state.ZoneID.ValueString(), q.Encode())

	var wglr wlanGroupListResp
	if err := doGET(ctx, d.client.HTTP, endpoint, &wglr); err != nil {
		resp.Diagnostics.AddError("read WLAN groups failed", err.Error())
		return
	}

	groups := make([]WLANGroupDSItem, 0, len(wglr.List))
	for _, wg := range wglr.List {
		item := WLANGroupDSItem{
			ID:   types.StringValue(wg.ID),
			Name: types.StringValue(wg.Name),
		}
		if wg.Description != "" {
			item.Description = types.StringValue(wg.Description)
		} else {
			item.Description = types.StringNull()
		}
		members := make([]string, 0, len(wg.Members))
		for _, m := range wg.Members {
			members = append(members, m.ID)
		}
		item.Members, _ = types.ListValueFrom(ctx, types.StringType, members)
		groups = append(groups, item)
	}
	state.Groups, _ = types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":          types.StringType,
			"name":        types.StringType,
			"description": types.StringType,
			"members":     types.ListType{ElemType: types.StringType},
		},
	}, groups)

	resp.State.Set(ctx, &state)
}
