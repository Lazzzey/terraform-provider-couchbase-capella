package schema

import (
	"testing"

	"github.com/jollytachome/terraform-provider-couchbase-capella/internal/errors"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

func TestNetworkPeerSchemaValidate(t *testing.T) {
	type test struct {
		expectedErr            error
		name                   string
		expectedProjectId      string
		expectedOrganizationId string
		expectedClusterId      string
		expectedPeerId         string
		input                  NetworkPeer
	}

	tests := []test{
		{
			name: "[POSITIVE] project ID, organization ID, and cluster ID are passed via terraform apply",
			input: NetworkPeer{
				Id:             basetypes.NewStringValue("100"),
				ClusterId:      basetypes.NewStringValue("200"),
				ProjectId:      basetypes.NewStringValue("300"),
				OrganizationId: basetypes.NewStringValue("400"),
			},
			expectedPeerId:         "100",
			expectedClusterId:      "200",
			expectedProjectId:      "300",
			expectedOrganizationId: "400",
		},
		{
			name: "[POSITIVE] IDs are passed via terraform import",
			input: NetworkPeer{
				Id: basetypes.NewStringValue("id=100,cluster_id=200,project_id=300,organization_id=400"),
			},
			expectedPeerId:         "100",
			expectedClusterId:      "200",
			expectedProjectId:      "300",
			expectedOrganizationId: "400",
		},
		{
			name: "[NEGATIVE] only Peer ID is passed via terraform apply",
			input: NetworkPeer{
				Id: basetypes.NewStringValue("100"),
			},
			expectedErr: errors.ErrInvalidImport,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			IDs, err := test.input.Validate()

			if test.expectedErr != nil {
				assert.ErrorContains(t, err, test.expectedErr.Error())
				return
			}

			assert.Equal(t, test.expectedPeerId, IDs[Id])
			assert.Equal(t, test.expectedClusterId, IDs[ClusterId])
			assert.Equal(t, test.expectedProjectId, IDs[ProjectId])
			assert.Equal(t, test.expectedOrganizationId, IDs[OrganizationId])
		})
	}
}
