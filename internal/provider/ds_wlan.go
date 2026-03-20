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

type WLANDS struct{ client *APIClient }

type WLANDSItem struct {
	ID                  types.String                  `tfsdk:"id"`
	Name                types.String                  `tfsdk:"name"`
	SSID                types.String                  `tfsdk:"ssid"`
	Description         types.String                  `tfsdk:"description"`
	GroupID             types.String                  `tfsdk:"group_id"`
	Encryption          *WLANEncryptionModel          `tfsdk:"encryption"`
	VLAN                *WLANVLANModel                `tfsdk:"vlan"`
	AccessTunnelProfile *WLANAccessTunnelProfileModel `tfsdk:"access_tunnel_profile"`
}

type WLANDSModel struct {
	ZoneID types.String `tfsdk:"zone_id"`
	WLANs  types.List   `tfsdk:"wlans"`
}

func NewWLANDataSource() datasource.DataSource { return &WLANDS{} }

func (d *WLANDS) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "ruckus_wlans"
}

func (d *WLANDS) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves a list of WLANs in a specified zone.",
		Attributes: map[string]schema.Attribute{
			"zone_id": schema.StringAttribute{
				Description: "ID of the zone to list WLANs from.",
				Required:    true,
			},
			"wlans": schema.ListNestedAttribute{
				Description: "List of WLANs.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the WLAN.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the WLAN.",
							Computed:    true,
						},
						"ssid": schema.StringAttribute{
							Description: "SSID of the WLAN.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the WLAN.",
							Computed:    true,
						},
						"group_id": schema.StringAttribute{
							Description: "ID of the WLAN Group this WLAN belongs to.",
							Computed:    true,
						},
						"encryption": schema.SingleNestedAttribute{
							Description: "Encryption configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"mode": schema.StringAttribute{
									Description: "Encryption mode.",
									Computed:    true,
								},
								"passphrase": schema.StringAttribute{
									Description: "Encryption passphrase.",
									Computed:    true,
									Sensitive:   true,
								},
								"algorithm": schema.StringAttribute{
									Description: "Encryption algorithm.",
									Computed:    true,
								},
							},
						},
						"vlan": schema.SingleNestedAttribute{
							Description: "VLAN configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"access_vlan": schema.Int64Attribute{
									Description: "Access VLAN.",
									Computed:    true,
								},
							},
						},
						"access_tunnel_profile": schema.SingleNestedAttribute{
							Description: "Access tunnel profile configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Description: "Access tunnel profile name.",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *WLANDS) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		d.client = req.ProviderData.(*APIClient)
	}
}

type wlanListResp struct {
	List []wlanListItem `json:"list"`
}

type wlanListItem struct {
	ID                  string                   `json:"id"`
	Name                string                   `json:"name"`
	SSID                string                   `json:"ssid"`
	Description         string                   `json:"description,omitempty"`
	GroupID             string                   `json:"groupId,omitempty"`
	Encryption          *wlanEncryption          `json:"encryption,omitempty"`
	VLAN                *wlanVLAN                `json:"vlan,omitempty"`
	AccessTunnelProfile *wlanAccessTunnelProfile `json:"accessTunnelProfile,omitempty"`
}

func (d *WLANDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("not configured", "missing API client")
		return
	}
	var state WLANDSModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// GET {base}/wsg/api/public/{ver}/rkszones/{zoneID}/wlans?serviceTicket=...
	q := url.Values{}
	q.Set("serviceTicket", d.client.ServiceTicket)

	endpoint := fmt.Sprintf("%s/wsg/api/public/%s/rkszones/%s/wlans?%s",
		d.client.BaseURL, d.client.APIVersion, state.ZoneID.ValueString(), q.Encode())

	var wlr wlanListResp
	if err := doGET(ctx, d.client.HTTP, endpoint, &wlr); err != nil {
		resp.Diagnostics.AddError("read WLANs failed", err.Error())
		return
	}

	wlans := make([]WLANDSItem, 0, len(wlr.List))
	for _, w := range wlr.List {
		item := WLANDSItem{
			ID:   types.StringValue(w.ID),
			Name: types.StringValue(w.Name),
			SSID: types.StringValue(w.SSID),
		}
		if w.Description != "" {
			item.Description = types.StringValue(w.Description)
		} else {
			item.Description = types.StringNull()
		}
		if w.GroupID != "" {
			item.GroupID = types.StringValue(w.GroupID)
		} else {
			item.GroupID = types.StringNull()
		}
		if w.Encryption != nil {
			item.Encryption = &WLANEncryptionModel{}
			if w.Encryption.Mode != "" {
				item.Encryption.Mode = types.StringValue(w.Encryption.Mode)
			} else {
				item.Encryption.Mode = types.StringNull()
			}
			if w.Encryption.Passphrase != "" {
				item.Encryption.Passphrase = types.StringValue(w.Encryption.Passphrase)
			} else {
				item.Encryption.Passphrase = types.StringNull()
			}
			if w.Encryption.Algorithm != "" {
				item.Encryption.Algorithm = types.StringValue(w.Encryption.Algorithm)
			} else {
				item.Encryption.Algorithm = types.StringNull()
			}
		}
		if w.VLAN != nil && w.VLAN.AccessVLAN != nil {
			item.VLAN = &WLANVLANModel{
				AccessVLAN: types.Int64Value(int64(*w.VLAN.AccessVLAN)),
			}
		}
		if w.AccessTunnelProfile != nil {
			item.AccessTunnelProfile = &WLANAccessTunnelProfileModel{
				Name: types.StringValue(w.AccessTunnelProfile.Name),
			}
		}
		wlans = append(wlans, item)
	}

	// Define the object type for the list
	objType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":          types.StringType,
			"name":        types.StringType,
			"ssid":        types.StringType,
			"description": types.StringType,
			"group_id":    types.StringType,
			"encryption": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"mode":       types.StringType,
					"passphrase": types.StringType,
					"algorithm":  types.StringType,
				},
			},
			"vlan": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"access_vlan": types.Int64Type,
				},
			},
			"access_tunnel_profile": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"name": types.StringType,
				},
			},
		},
	}

	state.WLANs, _ = types.ListValueFrom(ctx, objType, wlans)

	resp.State.Set(ctx, &state)
}
