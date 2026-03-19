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
- **Build & Install**: Run `make install` (builds with `go build -v ./...`, installs with `go install -v ./...`). Provider address is `registry.terraform.io/nshreck/ruckus`.
- **Lint**: `make lint` runs `golangci-lint run`.
- **Format**: `make fmt` applies `gofmt -s -w -e .`.
- **Test**: `make test` for unit tests (`go test -v -cover -timeout=120s -parallel=10 ./...`); `make testacc` for acceptance tests with `TF_ACC=1` and 120m timeout.
- **Generate**: `make generate` runs `go generate ./...` in the `tools/` subdirectory.

## Project-Specific Patterns
- **API Endpoints**: Construct URLs like `fmt.Sprintf("%s/wsg/api/public/%s/rkszones/%s/wlans?serviceTicket=%s", baseURL, apiVersion, zoneID, ticket)`. Always include `serviceTicket` in query params for auth.
- **Nested Blocks**: Use `schema.SingleNestedBlock` for complex attributes (e.g., `encryption`, `vlan` in WLAN resource). Map Terraform models to API structs with conditional field setting (e.g., `if !plan.Encryption.Mode.IsNull() { enc.Mode = plan.Encryption.Mode.ValueString() }`).
- **HTTP Utilities**: Use `doJSON()` for request/response cycles with automatic status checking (200-299) and body draining. For simple GET requests with JSON, use `doGET()`. For lower-level control, use `doRequest()` and manually handle response with `closeWith()` and `drainBody()` for connection reuse.
- **API Payloads**: Define Go structs with `json` tags matching Ruckus API (e.g., `type wlanEncryption struct { Mode string `json:"method,omitempty"` }`). Use pointers for optional fields in API structs.
- **Error Handling**: `doJSON()` checks HTTP status in 200-299 range and returns errors as `fmt.Errorf("http status %d", resp.StatusCode)`. Always use `resp.Diagnostics.AddError()` for Terraform-level errors. Drain response bodies even on error for keep-alive connection reuse.
- **Defaults**: API version defaults to `"v13_1"` for SmartZone 7.1.1; timeout to 30s. Domain optional in login.
- **Examples**: Reference `examples/basic/main.tf` for provider config and resource usage with nested blocks (encryption, vlan).
- **Schema Definition**: Mark sensitive fields with `Sensitive: true` (e.g., passwords, passphrases). Use validators like `stringvalidator.OneOf` for enum values.

## Integration Points
- **External API**: Ruckus SmartZone REST API. Verify payloads against controller's OpenAPI at `https://{host}:8443/wsg/apiDoc/openapi`.
- **Dependencies**: Relies on `github.com/hashicorp/terraform-plugin-framework` for provider framework; no external libraries beyond standard and HashiCorp's.
- **Cross-Component Communication**: APIClient shared between resources/data sources via `req.ProviderData`. Resources reference zone IDs from data sources (e.g., `zone_id = data.ruckus_zone.hq.id`).</content>
<parameter name="filePath">C:\Users\shrec\GolandProjects\terraform-provider-ruckus\AGENTS.md
