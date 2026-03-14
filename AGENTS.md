# AI Coding Agent Guidelines for terraform-provider-ruckus

## Architecture Overview
This is a Terraform provider for Ruckus SmartZone controllers, built on HashiCorp's Terraform Plugin Framework. It manages WLAN configurations via REST API calls to `/wsg/api/public/{api_version}/` endpoints. Authentication uses `serviceTicket` obtained via login and passed in query parameters.

Key components:
- `internal/provider/provider.go`: Defines provider schema (host, username, password, api_version, etc.) and configures APIClient with authenticated HTTP client.
- `internal/provider/api_client.go`: APIClient struct and login logic for serviceTicket.
- `internal/provider/http.go`: HTTP utilities (`doJSON`, `doGET`) for JSON API interactions.
- Resources in `internal/provider/res_*.go`: Implement CRUD for WLANs, WLAN groups.
- Data sources in `internal/provider/ds_*.go`: Read-only access, e.g., zones by name.

## Developer Workflows
- **Build & Install**: Run `make install` (builds with `go build -v ./...`, installs with `go install -v ./...`).
- **Lint**: `make lint` runs `golangci-lint run`.
- **Format**: `make fmt` applies `gofmt -s -w -e .`.
- **Test**: `make test` for unit tests (`go test -v -cover -timeout=120s -parallel=10 ./...`); `make testacc` for acceptance tests with `TF_ACC=1` and 120m timeout.
- **Generate**: `make generate` runs `go generate ./...` in `tools/` directory (if present).

## Project-Specific Patterns
- **API Endpoints**: Construct URLs like `fmt.Sprintf("%s/wsg/api/public/%s/rkszones/%s/wlans?serviceTicket=%s", baseURL, apiVersion, zoneID, ticket)`. Always include `serviceTicket` in query params for auth.
- **Nested Blocks**: Use `schema.SingleNestedBlock` for complex attributes (e.g., `security`, `vlan` in WLAN resource). Map Terraform models to API structs with conditional field setting (e.g., `if !plan.Security.Mode.IsNull() { sec.Mode = plan.Security.Mode.ValueString() }`).
- **API Payloads**: Define Go structs with `json` tags matching Ruckus API (e.g., `type wlanSecurity struct { Mode string `json:"method,omitempty"` }`). Use pointers for optional fields in API structs.
- **Error Handling**: Check HTTP status in 200-299 range; drain response body on errors. Use `resp.Diagnostics.AddError()` for Terraform errors.
- **Defaults**: API version defaults to `"v13_1"` for SmartZone 7.1.1; timeout to 30s. Domain optional in login.
- **Examples**: Reference `examples/basic/main.tf` for provider config and resource usage with nested blocks.

## Integration Points
- **External API**: Ruckus SmartZone REST API. Verify payloads against controller's OpenAPI at `https://{host}:8443/wsg/apiDoc/openapi`.
- **Dependencies**: Relies on `github.com/hashicorp/terraform-plugin-framework` for provider framework; no external libraries beyond standard and HashiCorp's.
- **Cross-Component Communication**: APIClient shared between resources/data sources via `req.ProviderData`. Resources reference zone IDs from data sources (e.g., `zone_id = data.ruckus_zone.hq.id`).</content>
<parameter name="filePath">C:\Users\shrec\GolandProjects\terraform-provider-ruckus\AGENTS.md
