// Code generated by ogen, DO NOT EDIT.

package client

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// AddGroupMemberV1 implements addGroupMemberV1 operation.
	//
	// Add Group Member.
	//
	// PUT /api/admin/v1/groups/{id}/members/{email}
	AddGroupMemberV1(ctx context.Context, params AddGroupMemberV1Params) error
	// BulkDeleteAttributesFormIdentities implements bulkDeleteAttributesFormIdentities operation.
	//
	// Delete attributes from multiple identities.
	//
	// DELETE /api/v2/bulk/identities/attributes
	BulkDeleteAttributesFormIdentities(ctx context.Context, req []IdentityAttributeKeysModel) (*IdentitiesAttributesResponseModel, error)
	// BulkUpsertAttributesForIdentities implements bulkUpsertAttributesForIdentities operation.
	//
	// Adds or Updates attributes to multiple identities.
	//
	// PUT /api/v2/bulk/identities/attributes
	BulkUpsertAttributesForIdentities(ctx context.Context, req []IdentityAttributeModel) (*IdentitiesAttributesResponseModel, error)
	// CreateAccessBundle implements createAccessBundle operation.
	//
	// Create Access Bundle.
	//
	// POST /api/v1/access-bundles
	CreateAccessBundle(ctx context.Context, req *UpsertAccessBundleV1) (*AccessBundleV1, error)
	// CreateAccessFlowV1 implements createAccessFlowV1 operation.
	//
	// Create Access Flow.
	//
	// POST /api/v1/access-flows
	CreateAccessFlowV1(ctx context.Context, req *UpsertAccessFlowV1) (*AccessFlowV1, error)
	// CreateAccessFlowV2 implements createAccessFlowV2 operation.
	//
	// Create access flow.
	//
	// POST /api/admin/v2/access-flows
	CreateAccessFlowV2(ctx context.Context, req *AccessFlowUpsertPublicV2Model) (*AccessFlowPublicV2Model, error)
	// CreateAccessRequest implements createAccessRequest operation.
	//
	// Create access request.
	//
	// POST /api/v3/access-requests
	CreateAccessRequest(ctx context.Context, req *CreateAccessRequest) (*AccessRequest, error)
	// CreateAccessRequestV4 implements createAccessRequestV4 operation.
	//
	// Create access request.
	//
	// POST /api/user/v4/access-requests
	CreateAccessRequestV4(ctx context.Context, req *CreateAccessRequestV4) ([]AccessRequestV4, error)
	// CreateAccessScopesV1 implements createAccessScopesV1 operation.
	//
	// Create Access Scope.
	//
	// POST /api/admin/v1/access-scopes
	CreateAccessScopesV1(ctx context.Context, req *UpsertAccessScopeV1) (*AccessScopeV1, error)
	// CreateBundleV2 implements createBundleV2 operation.
	//
	// Create Bundle.
	//
	// POST /api/admin/v2/bundles
	CreateBundleV2(ctx context.Context, req *UpsertBundlePublicV2Model) (*BundlePublicV2Model, error)
	// CreateGroupV1 implements createGroupV1 operation.
	//
	// Create Group.
	//
	// POST /api/admin/v1/groups
	CreateGroupV1(ctx context.Context, req *CreateGroupV1) (*GroupV1, error)
	// CreateIntegrationV2 implements createIntegrationV2 operation.
	//
	// Create integration.
	//
	// POST /api/v2/integrations
	CreateIntegrationV2(ctx context.Context, req *CreateIntegration) (*Integration, error)
	// CreateIntegrationV4 implements createIntegrationV4 operation.
	//
	// Create integration.
	//
	// POST /api/admin/v4/integrations
	CreateIntegrationV4(ctx context.Context, req *CreateIntegrationV4) (*IntegrationV4, error)
	// DeleteAccessBundle implements deleteAccessBundle operation.
	//
	// Delete Access Bundle.
	//
	// DELETE /api/v1/access-bundles/{id}
	DeleteAccessBundle(ctx context.Context, params DeleteAccessBundleParams) (*MessageResponse, error)
	// DeleteAccessFlowV1 implements deleteAccessFlowV1 operation.
	//
	// Delete Access Flow.
	//
	// DELETE /api/v1/access-flows/{id}
	DeleteAccessFlowV1(ctx context.Context, params DeleteAccessFlowV1Params) (*MessageResponse, error)
	// DeleteAccessFlowV2 implements deleteAccessFlowV2 operation.
	//
	// Delete access flow.
	//
	// DELETE /api/admin/v2/access-flows/{id}
	DeleteAccessFlowV2(ctx context.Context, params DeleteAccessFlowV2Params) error
	// DeleteAccessScopesV1 implements deleteAccessScopesV1 operation.
	//
	// Delete Access Scope.
	//
	// DELETE /api/admin/v1/access-scopes/{id}
	DeleteAccessScopesV1(ctx context.Context, params DeleteAccessScopesV1Params) error
	// DeleteBundleV2 implements deleteBundleV2 operation.
	//
	// Delete Bundle.
	//
	// DELETE /api/admin/v2/bundles/{id}
	DeleteBundleV2(ctx context.Context, params DeleteBundleV2Params) error
	// DeleteConnectorV3 implements deleteConnectorV3 operation.
	//
	// Delete Connector.
	//
	// DELETE /api/admin/v3/connectors/{id}
	DeleteConnectorV3(ctx context.Context, params DeleteConnectorV3Params) error
	// DeleteGroupV1 implements deleteGroupV1 operation.
	//
	// Delete Group.
	//
	// DELETE /api/admin/v1/groups/{id}
	DeleteGroupV1(ctx context.Context, params DeleteGroupV1Params) error
	// DeleteIntegrationV2 implements deleteIntegrationV2 operation.
	//
	// Delete integration.
	//
	// DELETE /api/v2/integrations/{id}
	DeleteIntegrationV2(ctx context.Context, params DeleteIntegrationV2Params) (*MessageResponse, error)
	// DeleteIntegrationV4 implements deleteIntegrationV4 operation.
	//
	// Delete integration.
	//
	// DELETE /api/admin/v4/integrations/{id}
	DeleteIntegrationV4(ctx context.Context, params DeleteIntegrationV4Params) error
	// GetAccessBundle implements getAccessBundle operation.
	//
	// Get Access Bundle.
	//
	// GET /api/v1/access-bundles/{id}
	GetAccessBundle(ctx context.Context, params GetAccessBundleParams) (*AccessBundleV1, error)
	// GetAccessFlowV1 implements getAccessFlowV1 operation.
	//
	// Get Access Flow.
	//
	// GET /api/v1/access-flows/{id}
	GetAccessFlowV1(ctx context.Context, params GetAccessFlowV1Params) (*AccessFlowV1, error)
	// GetAccessFlowV2 implements getAccessFlowV2 operation.
	//
	// Get access flow.
	//
	// GET /api/admin/v2/access-flows/{id}
	GetAccessFlowV2(ctx context.Context, params GetAccessFlowV2Params) (*AccessFlowPublicV2Model, error)
	// GetAccessRequest implements getAccessRequest operation.
	//
	// Get access request.
	//
	// GET /api/v3/access-requests/{id}
	GetAccessRequest(ctx context.Context, params GetAccessRequestParams) (*AccessRequest, error)
	// GetAccessRequestDetails implements getAccessRequestDetails operation.
	//
	// Get access request access details.
	//
	// GET /api/v3/access-requests/{id}/access-details
	GetAccessRequestDetails(ctx context.Context, params GetAccessRequestDetailsParams) (*ConnectDetailsResponse, error)
	// GetAccessRequestEntitlementsV4 implements getAccessRequestEntitlementsV4 operation.
	//
	// Get access request entitlements.
	//
	// GET /api/user/v4/access-requests/{id}/entitlements
	GetAccessRequestEntitlementsV4(ctx context.Context, params GetAccessRequestEntitlementsV4Params) (*PublicApiListResponseAccessRequestEntitlementPublicV4Model, error)
	// GetAccessRequestsV4 implements getAccessRequestsV4 operation.
	//
	// Get access request.
	//
	// GET /api/user/v4/access-requests/{id}
	GetAccessRequestsV4(ctx context.Context, params GetAccessRequestsV4Params) (*AccessRequestV4, error)
	// GetAccessScopesV1 implements getAccessScopesV1 operation.
	//
	// Get Access Scope.
	//
	// GET /api/admin/v1/access-scopes/{id}
	GetAccessScopesV1(ctx context.Context, params GetAccessScopesV1Params) (*AccessScopeV1, error)
	// GetAccessSessionAccessDetailsV1 implements getAccessSessionAccessDetailsV1 operation.
	//
	// Get session access details.
	//
	// GET /api/user/v1/access-sessions/{id}/access-details
	GetAccessSessionAccessDetailsV1(ctx context.Context, params GetAccessSessionAccessDetailsV1Params) (*AccessSessionDetailsV1, error)
	// GetAccessSessionV1 implements getAccessSessionV1 operation.
	//
	// Get access session.
	//
	// GET /api/user/v1/access-sessions/{id}
	GetAccessSessionV1(ctx context.Context, params GetAccessSessionV1Params) (*AccessSessionV1, error)
	// GetBundleV2 implements getBundleV2 operation.
	//
	// Get Bundle.
	//
	// GET /api/admin/v2/bundles/{id}
	GetBundleV2(ctx context.Context, params GetBundleV2Params) (*BundlePublicV2Model, error)
	// GetConnectorV3 implements getConnectorV3 operation.
	//
	// Get Connector.
	//
	// GET /api/admin/v3/connectors/{id}
	GetConnectorV3(ctx context.Context, params GetConnectorV3Params) (*ConnectorV3, error)
	// GetGroupV1 implements getGroupV1 operation.
	//
	// Get Group.
	//
	// GET /api/admin/v1/groups/{id}
	GetGroupV1(ctx context.Context, params GetGroupV1Params) (*GroupV1, error)
	// GetIntegrationConfig implements getIntegrationConfig operation.
	//
	// Get integration config.
	//
	// GET /api/v2/integrations-catalog/{type}
	GetIntegrationConfig(ctx context.Context, params GetIntegrationConfigParams) (*IntegrationConfig, error)
	// GetIntegrationPermissions implements getIntegrationPermissions operation.
	//
	// Get integration permissions for the entire tenant.
	//
	// GET /api/v3/integrations/{id}/permissions
	GetIntegrationPermissions(ctx context.Context, params GetIntegrationPermissionsParams) (*PaginatedResponsePermissionV3Response, error)
	// GetIntegrationResources implements getIntegrationResources operation.
	//
	// Get integration resources for the entire tenant.
	//
	// GET /api/v3/integrations/{id}/resources
	GetIntegrationResources(ctx context.Context, params GetIntegrationResourcesParams) (*PaginatedResponseResourceV3Response, error)
	// GetIntegrationV2 implements getIntegrationV2 operation.
	//
	// Get integration.
	//
	// GET /api/v2/integrations/{id}
	GetIntegrationV2(ctx context.Context, params GetIntegrationV2Params) (*Integration, error)
	// GetIntegrationsByIdV4 implements getIntegrationsByIdV4 operation.
	//
	// Get integration by id.
	//
	// GET /api/admin/v4/integrations/{id}
	GetIntegrationsByIdV4(ctx context.Context, params GetIntegrationsByIdV4Params) (*IntegrationV4, error)
	// GetResourceUserTags implements getResourceUserTags operation.
	//
	// Get user tags of a resource.
	//
	// GET /api/v3/integrations/resources/{resource_id}/user-tags
	GetResourceUserTags(ctx context.Context, params GetResourceUserTagsParams) (*ResourceUserTagsResponse, error)
	// GetSelectableIntegrations implements getSelectableIntegrations operation.
	//
	// Get selectable integrations.
	//
	// GET /api/v3/selectable-integrations
	GetSelectableIntegrations(ctx context.Context, params GetSelectableIntegrationsParams) (*PaginatedResponseSelectableIntegrationV3, error)
	// GetSelectablePermissions implements getSelectablePermissions operation.
	//
	// Get selectable permissions.
	//
	// GET /api/v3/selectable-integrations/{integration_id}/{resource_type}/permissions
	GetSelectablePermissions(ctx context.Context, params GetSelectablePermissionsParams) (*SelectablePermissionsResponse, error)
	// GetSelectableResourceTypes implements getSelectableResourceTypes operation.
	//
	// Get selectable resource types.
	//
	// GET /api/v3/selectable-integrations/{integration_id}/resource-types
	GetSelectableResourceTypes(ctx context.Context, params GetSelectableResourceTypesParams) (*PaginatedResponseSelectableResourceTypeV3, error)
	// GetSelectableResources implements getSelectableResources operation.
	//
	// Get selectable resources.
	//
	// GET /api/v3/selectable-integrations/{integration_id}/{resource_type}/resources
	GetSelectableResources(ctx context.Context, params GetSelectableResourcesParams) (*PaginatedResponseSelectableResourceV3, error)
	// GetUser implements getUser operation.
	//
	// Get user by Id or Email.
	//
	// GET /api/v2/users/{id}
	GetUser(ctx context.Context, params GetUserParams) (*UserModel, error)
	// ListAccessBundles implements listAccessBundles operation.
	//
	// List Access Bundles.
	//
	// GET /api/v1/access-bundles
	ListAccessBundles(ctx context.Context) (*PaginatedResponseAccessBundleV1Model, error)
	// ListAccessFlowsV1 implements listAccessFlowsV1 operation.
	//
	// List Access Flows.
	//
	// GET /api/v1/access-flows
	ListAccessFlowsV1(ctx context.Context) (*PaginatedResponseAccessFlowV1Model, error)
	// ListAccessFlowsV2 implements listAccessFlowsV2 operation.
	//
	// List access flows.
	//
	// GET /api/admin/v2/access-flows
	ListAccessFlowsV2(ctx context.Context, params ListAccessFlowsV2Params) (*PublicApiListResponseAccessFlowPublicV2Model, error)
	// ListAccessRequests implements listAccessRequests operation.
	//
	// List access requests.
	//
	// GET /api/v3/access-requests
	ListAccessRequests(ctx context.Context, params ListAccessRequestsParams) (*PaginatedResponseAccessRequestV3, error)
	// ListAccessRequestsV4 implements listAccessRequestsV4 operation.
	//
	// List access requests.
	//
	// GET /api/user/v4/access-requests
	ListAccessRequestsV4(ctx context.Context, params ListAccessRequestsV4Params) (*PublicApiListResponseAccessRequestV4PublicModel, error)
	// ListAccessScopesV1 implements listAccessScopesV1 operation.
	//
	// List Access Scopes.
	//
	// GET /api/admin/v1/access-scopes
	ListAccessScopesV1(ctx context.Context, params ListAccessScopesV1Params) (*PublicApiListResponseAccessScopePublicV1Model, error)
	// ListAccessSessionsV1 implements listAccessSessionsV1 operation.
	//
	// List access sessions.
	//
	// GET /api/user/v1/access-sessions
	ListAccessSessionsV1(ctx context.Context, params ListAccessSessionsV1Params) (*PublicApiListResponseAccessSessionPublicV1Model, error)
	// ListActivity implements listActivity operation.
	//
	// List Activity.
	//
	// GET /api/v3/activity
	ListActivity(ctx context.Context, params ListActivityParams) (*PaginatedResponseActivityReportJsonExportModel, error)
	// ListAttributesForIdentities implements listAttributesForIdentities operation.
	//
	// List attributes for multiple identities.
	//
	// GET /api/v2/bulk/identities/attributes
	ListAttributesForIdentities(ctx context.Context, params ListAttributesForIdentitiesParams) (*PaginatedResponseIdentityAttributeModel, error)
	// ListAvailableBundlesV1 implements listAvailableBundlesV1 operation.
	//
	// List available bundles.
	//
	// GET /api/user/v1/available-access/bundles
	ListAvailableBundlesV1(ctx context.Context, params ListAvailableBundlesV1Params) (*PublicApiListResponseAvailableBundlePublicV1Model, error)
	// ListAvailableEntitlementsV1 implements listAvailableEntitlementsV1 operation.
	//
	// List available entitlements.
	//
	// GET /api/user/v1/available-access/entitlements
	ListAvailableEntitlementsV1(ctx context.Context, params ListAvailableEntitlementsV1Params) (*PublicApiListResponseAvailableEntitlementPublicV1Model, error)
	// ListBundlesV2 implements listBundlesV2 operation.
	//
	// List Bundles.
	//
	// GET /api/admin/v2/bundles
	ListBundlesV2(ctx context.Context, params ListBundlesV2Params) (*PublicApiListResponseBundlePublicV2Model, error)
	// ListConnectors implements listConnectors operation.
	//
	// List connectors.
	//
	// GET /api/v2/connectors
	ListConnectors(ctx context.Context) ([]Connector, error)
	// ListConnectorsV3 implements listConnectorsV3 operation.
	//
	// List Connectors.
	//
	// GET /api/admin/v3/connectors
	ListConnectorsV3(ctx context.Context, params ListConnectorsV3Params) (*PublicApiListResponseConnectorPublicV3Model, error)
	// ListGroupMembersV1 implements listGroupMembersV1 operation.
	//
	// Get Group Members.
	//
	// GET /api/admin/v1/groups/{id}/members
	ListGroupMembersV1(ctx context.Context, params ListGroupMembersV1Params) (*PublicApiListResponseGroupMemberPublicV1Model, error)
	// ListGroupsV1 implements listGroupsV1 operation.
	//
	// List Groups.
	//
	// GET /api/admin/v1/groups
	ListGroupsV1(ctx context.Context, params ListGroupsV1Params) (*PublicApiListResponseGroupPublicV1Model, error)
	// ListIdentities implements listIdentities operation.
	//
	// List identities, grantees and approvers.
	//
	// GET /api/v2/identities
	ListIdentities(ctx context.Context) (*PaginatedResponseIdentityModelV2, error)
	// ListIntegrationConfigs implements listIntegrationConfigs operation.
	//
	// List integration configs.
	//
	// GET /api/v2/integrations-catalog
	ListIntegrationConfigs(ctx context.Context) (*PaginatedResponseIntegrationConfigPublicModel, error)
	// ListIntegrationsV2 implements listIntegrationsV2 operation.
	//
	// List integrations.
	//
	// GET /api/v2/integrations
	ListIntegrationsV2(ctx context.Context) (*PaginatedResponseIntegrationModel, error)
	// ListIntegrationsV4 implements listIntegrationsV4 operation.
	//
	// List integrations.
	//
	// GET /api/admin/v4/integrations
	ListIntegrationsV4(ctx context.Context, params ListIntegrationsV4Params) (*PublicApiListResponseIntegrationPublicV4Model, error)
	// ListUsers implements listUsers operation.
	//
	// List users.
	//
	// GET /api/v2/users
	ListUsers(ctx context.Context) (*PaginatedResponseUserModel, error)
	// RefreshIntegrationV2 implements refreshIntegrationV2 operation.
	//
	// Refresh integration.
	//
	// POST /api/v2/integrations/{id}/refresh
	RefreshIntegrationV2(ctx context.Context, params RefreshIntegrationV2Params) (*MessageResponse, error)
	// RemoveGroupMemberV1 implements removeGroupMemberV1 operation.
	//
	// Remove Group Member.
	//
	// DELETE /api/admin/v1/groups/{id}/members/{email}
	RemoveGroupMemberV1(ctx context.Context, params RemoveGroupMemberV1Params) error
	// RequestAccessAgainV4 implements requestAccessAgainV4 operation.
	//
	// Request access again.
	//
	// POST /api/user/v4/access-requests/{id}/request-again
	RequestAccessAgainV4(ctx context.Context, req *RequestAgainV4, params RequestAccessAgainV4Params) ([]AccessRequestV4, error)
	// ResetAccessRequestCredentials implements resetAccessRequestCredentials operation.
	//
	// Reset access request credentials.
	//
	// POST /api/v3/access-requests/{id}/reset
	ResetAccessRequestCredentials(ctx context.Context, params ResetAccessRequestCredentialsParams) (*MessageResponse, error)
	// ResetAccessSessionCredentialsV1 implements resetAccessSessionCredentialsV1 operation.
	//
	// Reset session credentials.
	//
	// POST /api/user/v1/access-sessions/{id}/reset-credentials
	ResetAccessSessionCredentialsV1(ctx context.Context, params ResetAccessSessionCredentialsV1Params) (*PublicApiMessageResponse, error)
	// RevokeAccessRequestV4 implements revokeAccessRequestV4 operation.
	//
	// Revoke access request.
	//
	// POST /api/user/v4/access-requests/{id}/revoke
	RevokeAccessRequestV4(ctx context.Context, params RevokeAccessRequestV4Params) (*PublicApiMessageResponse, error)
	// RevokeAccessRequests implements revokeAccessRequests operation.
	//
	// Revoke multiple access requests.
	//
	// POST /api/v3/access-requests-bulk/revoke
	RevokeAccessRequests(ctx context.Context, req *AccessRequestsBulkRevokeRequestV3) (*MessageResponse, error)
	// UpdateAccessBundle implements updateAccessBundle operation.
	//
	// Update Access Bundle.
	//
	// PATCH /api/v1/access-bundles/{id}
	UpdateAccessBundle(ctx context.Context, req *UpdateAccessBundleV1, params UpdateAccessBundleParams) (*AccessBundleV1, error)
	// UpdateAccessFlowV1 implements updateAccessFlowV1 operation.
	//
	// Update Access Flow.
	//
	// PATCH /api/v1/access-flows/{id}
	UpdateAccessFlowV1(ctx context.Context, req *UpdateAccessFlowV1, params UpdateAccessFlowV1Params) (*AccessFlowV1, error)
	// UpdateAccessFlowV2 implements updateAccessFlowV2 operation.
	//
	// Update access flow.
	//
	// PUT /api/admin/v2/access-flows/{id}
	UpdateAccessFlowV2(ctx context.Context, req *AccessFlowUpsertPublicV2Model, params UpdateAccessFlowV2Params) (*AccessFlowPublicV2Model, error)
	// UpdateAccessScopesV1 implements updateAccessScopesV1 operation.
	//
	// Update Access Scope.
	//
	// PUT /api/admin/v1/access-scopes/{id}
	UpdateAccessScopesV1(ctx context.Context, req *UpsertAccessScopeV1, params UpdateAccessScopesV1Params) (*AccessScopeV1, error)
	// UpdateBundleV2 implements updateBundleV2 operation.
	//
	// Update Bundle.
	//
	// PUT /api/admin/v2/bundles/{id}
	UpdateBundleV2(ctx context.Context, req *UpsertBundlePublicV2Model, params UpdateBundleV2Params) (*BundlePublicV2Model, error)
	// UpdateConnectorV3 implements updateConnectorV3 operation.
	//
	// Update Connector.
	//
	// PUT /api/admin/v3/connectors/{id}
	UpdateConnectorV3(ctx context.Context, req *UpsertConnectorV3, params UpdateConnectorV3Params) (*ConnectorV3, error)
	// UpdateGroupMembersV1 implements updateGroupMembersV1 operation.
	//
	// Update Group Members.
	//
	// PUT /api/admin/v1/groups/{id}/members
	UpdateGroupMembersV1(ctx context.Context, req *UpdateGroupMembersV1, params UpdateGroupMembersV1Params) error
	// UpdateGroupV1 implements updateGroupV1 operation.
	//
	// Update Group.
	//
	// PUT /api/admin/v1/groups/{id}/name
	UpdateGroupV1(ctx context.Context, req *UpdateGroupV1, params UpdateGroupV1Params) (*GroupV1, error)
	// UpdateIntegrationV2 implements updateIntegrationV2 operation.
	//
	// Update integration.
	//
	// PUT /api/v2/integrations/{id}
	UpdateIntegrationV2(ctx context.Context, req *UpdateIntegration, params UpdateIntegrationV2Params) (*Integration, error)
	// UpdateIntegrationV4 implements updateIntegrationV4 operation.
	//
	// Update integration.
	//
	// PUT /api/admin/v4/integrations/{id}
	UpdateIntegrationV4(ctx context.Context, req *UpdateIntegrationV4, params UpdateIntegrationV4Params) (*IntegrationV4, error)
	// UpdateResourceUserTags implements updateResourceUserTags operation.
	//
	// Update user tags of a resource.
	//
	// PUT /api/v3/integrations/resources/{resource_id}/user-tags
	UpdateResourceUserTags(ctx context.Context, req *UpdateResourceUserTagsRequest, params UpdateResourceUserTagsParams) (*MessageResponse, error)
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h   Handler
	sec SecurityHandler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, sec SecurityHandler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		sec:        sec,
		baseServer: s,
	}, nil
}
