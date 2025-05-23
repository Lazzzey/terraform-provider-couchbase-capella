package resources

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/jollytachome/terraform-provider-couchbase-capella/internal/api"
	internalerrors "github.com/jollytachome/terraform-provider-couchbase-capella/internal/errors"
	providerschema "github.com/jollytachome/terraform-provider-couchbase-capella/internal/schema"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ApiKey{}
	_ resource.ResourceWithConfigure   = &ApiKey{}
	_ resource.ResourceWithImportState = &ApiKey{}
)

const errorMessageAfterApiKeyCreation = "Api Key creation is successful, but encountered an error while checking the current" +
	" state of the api key. Please run `terraform plan` after 1-2 minutes to know the" +
	" current api key state. Additionally, run `terraform apply --refresh-only` to update" +
	" the state from remote, unexpected error: "

const errorMessageWhileApiKeyCreation = "There is an error during api key creation. Please check in Capella to see if any hanging resources" +
	" have been created, unexpected error: "

// ApiKey is the ApiKey resource implementation.
type ApiKey struct {
	*providerschema.Data
}

func NewApiKey() resource.Resource {
	return &ApiKey{}
}

// Metadata returns the apiKey resource type name.
func (r *ApiKey) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apikey"
}

// Schema defines the schema for the apiKey resource.
func (r *ApiKey) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ApiKeySchema()
}

// Configure adds the provider configured client to the apiKey resource.
func (r *ApiKey) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerschema.Data)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ProviderSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.Data = data
}

// Create creates a new apiKey.
func (a *ApiKey) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan providerschema.ApiKey
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := a.validateCreateApiKeyRequest(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ApiKey",
			"Could not create ApiKey, unexpected error:"+err.Error(),
		)
		return
	}

	var organizationId = plan.OrganizationId.ValueString()

	apiKeyRequest := api.CreateApiKeyRequest{
		Name:              plan.Name.ValueString(),
		OrganizationRoles: a.convertOrganizationRoles(plan.OrganizationRoles),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		apiKeyRequest.Description = plan.Description.ValueStringPointer()
	}

	if !plan.Expiry.IsNull() && !plan.Expiry.IsUnknown() {
		expiry := float32(*plan.Expiry.ValueFloat64Pointer())
		apiKeyRequest.Expiry = &expiry
	}

	convertedResources, err := a.convertResources(plan.Resources)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ApiKey",
			"Could not create ApiKey, unexpected error:"+err.Error(),
		)
		return
	}
	apiKeyRequest.Resources = &convertedResources

	if !plan.AllowedCIDRs.IsNull() && !plan.AllowedCIDRs.IsUnknown() {
		convertedAllowedCidr, err := a.convertAllowedCidrs(ctx, plan.AllowedCIDRs)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating ApiKey",
				"Could not create ApiKey, unexpected error:"+err.Error(),
			)
			return
		}
		apiKeyRequest.AllowedCIDRs = &convertedAllowedCidr
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/apikeys", a.HostURL, organizationId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		apiKeyRequest,
		a.Token,
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ApiKey Here",
			errorMessageWhileApiKeyCreation+api.ParseError(err),
		)
		return
	}

	apiKeyResponse := api.CreateApiKeyResponse{}
	err = json.Unmarshal(response.Body, &apiKeyResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ApiKey",
			errorMessageWhileApiKeyCreation+"error during unmarshalling: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, initializeApiKeyWithPlanAndId(plan, apiKeyResponse.Id))
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	refreshedState, err := a.retrieveApiKey(ctx, organizationId, apiKeyResponse.Id)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error creating ApiKey",
			errorMessageAfterApiKeyCreation+api.ParseError(err),
		)
		return
	}

	refreshedState.Token = types.StringValue(apiKeyResponse.Token)
	refreshedState = a.retainResourcesIfOrgOwner(&plan, refreshedState)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, refreshedState)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads ApiKey information.
func (a *ApiKey) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state providerschema.ApiKey
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading api key",
			"Could not read api key id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		apiKeyId       = resourceIDs[providerschema.Id]
	)

	// Get refreshed api key value from Capella
	refreshedState, err := a.retrieveApiKey(ctx, organizationId, apiKeyId)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading api key",
			"Could not read api key id "+state.Id.String()+": "+errString,
		)
		return
	}

	refreshedState.Token = state.Token
	refreshedState.Rotate = state.Rotate
	refreshedState.Secret = state.Secret
	refreshedState = a.retainResourcesIfOrgOwner(&state, refreshedState)

	// Set refreshed state
	diags = resp.State.Set(ctx, &refreshedState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update rotates the ApiKey.
func (a *ApiKey) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan, state providerschema.ApiKey
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error rotate api key",
			"Could not rotate api key id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		apiKeyId       = resourceIDs[providerschema.Id]
	)

	if plan.Rotate.IsNull() || plan.Rotate.IsUnknown() {
		resp.Diagnostics.AddError(
			"Error rotating api key",
			"Could not rotate api key id "+state.Id.String()+": rotate value is not set",
		)
		return
	}

	if !state.Rotate.IsNull() && !state.Rotate.IsUnknown() {
		planRotate := *plan.Rotate.ValueBigFloat()
		stateRotate := *state.Rotate.ValueBigFloat()
		if planRotate.Cmp(&stateRotate) != 1 {
			resp.Diagnostics.AddError(
				"Error rotating api key",
				"Could not rotate api key id "+state.Id.String()+": plan rotate value is not greater than state rotate value",
			)
			return
		}
	}

	var rotateApiRequest api.RotateApiKeyRequest
	if !plan.Secret.IsNull() || !plan.Secret.IsUnknown() {
		rotateApiRequest = api.RotateApiKeyRequest{
			Secret: plan.Secret.ValueStringPointer(),
		}
	}

	url := fmt.Sprintf("%s/v4/organizations/%s/apikeys/%s/rotate", a.HostURL, organizationId, apiKeyId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodPost, SuccessStatus: http.StatusCreated}
	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		rotateApiRequest,
		a.Token,
		nil,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error rotating api key",
			"Could not rotate api key id "+state.Id.String()+": "+errString,
		)
		return
	}

	rotateApiKeyResponse := api.RotateApiKeyResponse{}
	err = json.Unmarshal(response.Body, &rotateApiKeyResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error rotating api key",
			"Could not rotate api key id "+state.Id.String()+": "+err.Error(),
		)
		return
	}

	currentState, err := a.retrieveApiKey(ctx, organizationId, apiKeyId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error rotating api key",
			"Could not rotate api key id "+state.Id.String()+": "+api.ParseError(err),
		)
		return
	}

	currentState.Secret = types.StringValue(rotateApiKeyResponse.SecretKey)
	if !currentState.Id.IsNull() && !currentState.Id.IsUnknown() && !currentState.Secret.IsNull() && !currentState.Secret.IsUnknown() {
		currentState.Token = types.StringValue(base64.StdEncoding.EncodeToString([]byte(currentState.Id.ValueString() + ":" + currentState.Secret.ValueString())))
	}
	currentState.Rotate = plan.Rotate
	currentState = a.retainResourcesIfOrgOwner(&plan, currentState)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, currentState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the ApiKey.
func (a *ApiKey) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state providerschema.ApiKey
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceIDs, err := state.Validate()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting api key",
			"Could not delete api key id "+state.Id.String()+" unexpected error: "+err.Error(),
		)
		return
	}

	var (
		organizationId = resourceIDs[providerschema.OrganizationId]
		apiKeyId       = resourceIDs[providerschema.Id]
	)

	// Delete existing api key
	url := fmt.Sprintf("%s/v4/organizations/%s/apikeys/%s", a.HostURL, organizationId, apiKeyId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodDelete, SuccessStatus: http.StatusNoContent}
	_, err = a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		resourceNotFound, errString := api.CheckResourceNotFoundError(err)
		if resourceNotFound {
			tflog.Info(ctx, "resource doesn't exist in remote server removing resource from state file")
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting api key",
			"Could not delete api key id "+state.Id.String()+" unexpected error: "+errString,
		)
		return
	}
}

func (a *ApiKey) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// retrieveApiKey retrieves apikey information for a specified organization and apiKeyId.
func (a *ApiKey) retrieveApiKey(ctx context.Context, organizationId, apiKeyId string) (*providerschema.ApiKey, error) {
	url := fmt.Sprintf("%s/v4/organizations/%s/apikeys/%s", a.HostURL, organizationId, apiKeyId)
	cfg := api.EndpointCfg{Url: url, Method: http.MethodGet, SuccessStatus: http.StatusOK}
	response, err := a.Client.ExecuteWithRetry(
		ctx,
		cfg,
		nil,
		a.Token,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", internalerrors.ErrExecutingRequest, err)
	}

	apiKeyResp := api.GetApiKeyResponse{}
	err = json.Unmarshal(response.Body, &apiKeyResp)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", internalerrors.ErrUnmarshallingResponse, err)
	}

	audit := providerschema.NewCouchbaseAuditData(apiKeyResp.Audit)

	auditObj, diags := types.ObjectValueFrom(ctx, audit.AttributeTypes(), audit)
	if diags.HasError() {
		return nil, fmt.Errorf("error while audit conversion")
	}

	refreshedState, err := providerschema.NewApiKey(&apiKeyResp, organizationId, auditObj)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", internalerrors.ErrRefreshingState, err)
	}
	return refreshedState, nil
}

// validateCreateApiKeyRequest validates the required fields in the create request.
func (a *ApiKey) validateCreateApiKeyRequest(plan providerschema.ApiKey) error {
	if plan.OrganizationId.IsNull() {
		return fmt.Errorf("organizationId cannot be empty")
	}
	if plan.Name.IsNull() {
		return fmt.Errorf("name cannot be empty")
	}
	if plan.OrganizationRoles == nil {
		return fmt.Errorf("organizationRoles cannot be empty")
	}
	if !plan.Rotate.IsNull() && !plan.Rotate.IsUnknown() {
		return fmt.Errorf("rotate value should not be set")
	}
	if !plan.Secret.IsNull() && !plan.Secret.IsUnknown() {
		return fmt.Errorf("secret should not be set while create operation")
	}
	return a.validateApiKeyAttributesTrimmed(plan)
}

func (a *ApiKey) validateApiKeyAttributesTrimmed(plan providerschema.ApiKey) error {
	if (!plan.Name.IsNull() && !plan.Name.IsUnknown()) && !providerschema.IsTrimmed(plan.Name.ValueString()) {
		return fmt.Errorf("name %s", internalerrors.ErrNotTrimmed)
	}
	if (!plan.Description.IsNull() && !plan.Description.IsUnknown()) && !providerschema.IsTrimmed(plan.Description.ValueString()) {
		return fmt.Errorf("description %s", internalerrors.ErrNotTrimmed)
	}
	return nil
}

// convertOrganizationRoles is used to convert all roles
// in an array of basetypes.StringValue to strings.
func (a *ApiKey) convertOrganizationRoles(organizationRoles []basetypes.StringValue) []string {
	var convertedRoles []string
	for _, role := range organizationRoles {
		convertedRoles = append(convertedRoles, role.ValueString())
	}
	return convertedRoles
}

// convertResource is used to convert a resource object containing nested fields
// of type basetypes.StringValue to a resource object containing nested fields of go defined type.
func (a *ApiKey) convertResources(resources []providerschema.ApiKeyResourcesItems) ([]api.ResourcesItems, error) {
	var convertedResources []api.ResourcesItems
	for _, resource := range resources {
		id, err := uuid.Parse(resource.Id.ValueString())
		if err != nil {
			return nil, fmt.Errorf("resource id is not valid uuid")
		}
		convertedResource := api.ResourcesItems{
			Id: id,
		}

		var convertedRoles []string
		for _, role := range resource.Roles {
			convertedRoles = append(convertedRoles, role.ValueString())
		}
		convertedResource.Roles = convertedRoles

		if !resource.Type.IsNull() && !resource.Type.IsUnknown() {
			convertedResource.Type = resource.Type.ValueStringPointer()
		}
		convertedResources = append(convertedResources, convertedResource)
	}
	return convertedResources, nil
}

// convertAllowedCidrs is used to convert allowed cidrs in types.List to array of string.
func (a *ApiKey) convertAllowedCidrs(ctx context.Context, allowedCidrs types.Set) ([]string, error) {
	elements := make([]types.String, 0, len(allowedCidrs.Elements()))
	diags := allowedCidrs.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		return nil, fmt.Errorf("error while extracting allowedCidrs elements")
	}

	var convertedAllowedCidrs []string
	for _, allowedCidr := range elements {
		convertedAllowedCidrs = append(convertedAllowedCidrs, allowedCidr.ValueString())
	}
	return convertedAllowedCidrs, nil
}

func (a *ApiKey) retainResourcesIfOrgOwner(apiKeyReq, apiKeyRes *providerschema.ApiKey) *providerschema.ApiKey {
	isOrgOwner := false
	for _, role := range apiKeyRes.OrganizationRoles {
		if role.ValueString() == "organizationOwner" {
			isOrgOwner = true
		}
	}
	if isOrgOwner {
		apiKeyRes.Resources = apiKeyReq.Resources
	}
	return apiKeyRes
}

// initializeApiKeyWithPlanAndId initializes an instance of providerschema.ApiKey
// with the specified plan and ID. It marks all computed fields as null.
func initializeApiKeyWithPlanAndId(plan providerschema.ApiKey, id string) providerschema.ApiKey {
	plan.Id = types.StringValue(id)
	if plan.Secret.IsNull() || plan.Secret.IsUnknown() {
		plan.Secret = types.StringNull()
	}
	if plan.Rotate.IsNull() || plan.Rotate.IsUnknown() {
		plan.Rotate = types.NumberNull()
	}
	plan.Token = types.StringNull()
	plan.Audit = types.ObjectNull(providerschema.CouchbaseAuditData{}.AttributeTypes())
	return plan
}
