/*
 * Copyright 2026 Vitruvian Software
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package iam

import (
	"testing"

	"github.com/VitruvianSoftware/pulumi-library/internal/testutil"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==================== Scope-Specific IAM Member Tests ====================

func TestNewOrganizationIAMMember(t *testing.T) {
	tracker := testutil.NewTracker()
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		_, err := NewOrganizationIAMMember(ctx, "test-org-iam", &OrganizationIAMMemberArgs{
			OrgID:  pulumi.String("123456"),
			Role:   pulumi.String("roles/viewer"),
			Member: pulumi.String("user:test@example.com"),
		})
		return err
	}, pulumi.WithMocks("test-project", "test-stack", tracker))
	require.NoError(t, err)

	members := tracker.RequireType(t, "gcp:organizations/iAMMember:IAMMember", 1)
	assert.Equal(t, "123456", members[0].Inputs["orgId"].StringValue())
	assert.Equal(t, "roles/viewer", members[0].Inputs["role"].StringValue())
	assert.Equal(t, "user:test@example.com", members[0].Inputs["member"].StringValue())
}

func TestNewFolderIAMMember(t *testing.T) {
	tracker := testutil.NewTracker()
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		_, err := NewFolderIAMMember(ctx, "test-folder-iam", &FolderIAMMemberArgs{
			FolderID: pulumi.String("folders/789"),
			Role:     pulumi.String("roles/editor"),
			Member:   pulumi.String("serviceAccount:sa@p.iam.gserviceaccount.com"),
		})
		return err
	}, pulumi.WithMocks("test-project", "test-stack", tracker))
	require.NoError(t, err)

	members := tracker.RequireType(t, "gcp:folder/iAMMember:IAMMember", 1)
	assert.Equal(t, "folders/789", members[0].Inputs["folder"].StringValue())
}

func TestNewProjectIAMMember(t *testing.T) {
	tracker := testutil.NewTracker()
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		_, err := NewProjectIAMMember(ctx, "test-proj-iam", &ProjectIAMMemberArgs{
			ProjectID: pulumi.String("my-project-id"),
			Role:      pulumi.String("roles/storage.admin"),
			Member:    pulumi.String("group:admins@example.com"),
		})
		return err
	}, pulumi.WithMocks("test-project", "test-stack", tracker))
	require.NoError(t, err)

	members := tracker.RequireType(t, "gcp:projects/iAMMember:IAMMember", 1)
	assert.Equal(t, "my-project-id", members[0].Inputs["project"].StringValue())
}

func TestNewServiceAccountIAMMember(t *testing.T) {
	tracker := testutil.NewTracker()
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		_, err := NewServiceAccountIAMMember(ctx, "test-sa-iam", &ServiceAccountIAMMemberArgs{
			ServiceAccountID: pulumi.String("projects/p/serviceAccounts/sa@p.iam.gserviceaccount.com"),
			Role:             pulumi.String("roles/iam.serviceAccountTokenCreator"),
			Member:           pulumi.String("user:dev@example.com"),
		})
		return err
	}, pulumi.WithMocks("test-project", "test-stack", tracker))
	require.NoError(t, err)

	tracker.RequireType(t, "gcp:serviceaccount/iAMMember:IAMMember", 1)
}

func TestNewBillingIAMMember(t *testing.T) {
	tracker := testutil.NewTracker()
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		_, err := NewBillingIAMMember(ctx, "test-bill-iam", &BillingIAMMemberArgs{
			BillingAccountID: pulumi.String("AAAAAA-BBBBBB-CCCCCC"),
			Role:             pulumi.String("roles/billing.user"),
			Member:           pulumi.String("serviceAccount:tf@p.iam.gserviceaccount.com"),
		})
		return err
	}, pulumi.WithMocks("test-project", "test-stack", tracker))
	require.NoError(t, err)

	members := tracker.RequireType(t, "gcp:billing/accountIamMember:AccountIamMember", 1)
	assert.Equal(t, "AAAAAA-BBBBBB-CCCCCC", members[0].Inputs["billingAccountId"].StringValue())
}

// ==================== Scope-Specific IAM Binding Tests ====================

func TestNewOrganizationIAMBinding(t *testing.T) {
	tracker := testutil.NewTracker()
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		_, err := NewOrganizationIAMBinding(ctx, "test-org-binding", &OrganizationIAMBindingArgs{
			OrgID: pulumi.String("123456"),
			Role:  pulumi.String("roles/viewer"),
			Members: pulumi.StringArray{
				pulumi.String("user:a@example.com"),
				pulumi.String("user:b@example.com"),
			},
		})
		return err
	}, pulumi.WithMocks("test-project", "test-stack", tracker))
	require.NoError(t, err)

	bindings := tracker.RequireType(t, "gcp:organizations/iAMBinding:IAMBinding", 1)
	assert.Equal(t, "123456", bindings[0].Inputs["orgId"].StringValue())
}

func TestNewFolderIAMBinding(t *testing.T) {
	tracker := testutil.NewTracker()
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		_, err := NewFolderIAMBinding(ctx, "test-folder-binding", &FolderIAMBindingArgs{
			FolderID: pulumi.String("folders/999"),
			Role:     pulumi.String("roles/editor"),
			Members:  pulumi.StringArray{pulumi.String("user:c@example.com")},
		})
		return err
	}, pulumi.WithMocks("test-project", "test-stack", tracker))
	require.NoError(t, err)

	tracker.RequireType(t, "gcp:folder/iAMBinding:IAMBinding", 1)
}

func TestNewProjectIAMBinding(t *testing.T) {
	tracker := testutil.NewTracker()
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		_, err := NewProjectIAMBinding(ctx, "test-proj-binding", &ProjectIAMBindingArgs{
			ProjectID: pulumi.String("my-project"),
			Role:      pulumi.String("roles/compute.admin"),
			Members:   pulumi.StringArray{pulumi.String("group:sre@example.com")},
		})
		return err
	}, pulumi.WithMocks("test-project", "test-stack", tracker))
	require.NoError(t, err)

	tracker.RequireType(t, "gcp:projects/iAMBinding:IAMBinding", 1)
}

func TestNewServiceAccountIAMBinding(t *testing.T) {
	tracker := testutil.NewTracker()
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		_, err := NewServiceAccountIAMBinding(ctx, "test-sa-binding", &ServiceAccountIAMBindingArgs{
			ServiceAccountID: pulumi.String("projects/p/serviceAccounts/sa@p.iam.gserviceaccount.com"),
			Role:             pulumi.String("roles/iam.serviceAccountUser"),
			Members:          pulumi.StringArray{pulumi.String("user:dev@example.com")},
		})
		return err
	}, pulumi.WithMocks("test-project", "test-stack", tracker))
	require.NoError(t, err)

	tracker.RequireType(t, "gcp:serviceaccount/iAMBinding:IAMBinding", 1)
}

func TestNewBillingIAMBinding(t *testing.T) {
	tracker := testutil.NewTracker()
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		_, err := NewBillingIAMBinding(ctx, "test-bill-binding", &BillingIAMBindingArgs{
			BillingAccountID: pulumi.String("AAAAAA-BBBBBB-CCCCCC"),
			Role:             pulumi.String("roles/billing.viewer"),
			Members:          pulumi.StringArray{pulumi.String("user:fin@example.com")},
		})
		return err
	}, pulumi.WithMocks("test-project", "test-stack", tracker))
	require.NoError(t, err)

	tracker.RequireType(t, "gcp:billing/accountIamBinding:AccountIamBinding", 1)
}

// ==================== Legacy Deprecated Constructor Tests ====================

func TestNewIAMMember_Legacy_Organization(t *testing.T) {
	tracker := testutil.NewTracker()
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		_, err := NewIAMMember(ctx, "test-legacy-org", &IAMMemberArgs{
			ParentID:   pulumi.String("123456"),
			ParentType: "organization",
			Role:       pulumi.String("roles/viewer"),
			Member:     pulumi.String("user:test@example.com"),
		})
		return err
	}, pulumi.WithMocks("test-project", "test-stack", tracker))
	require.NoError(t, err)

	tracker.RequireType(t, "gcp:organizations/iAMMember:IAMMember", 1)
}

func TestNewIAMMember_Legacy_UnsupportedType(t *testing.T) {
	tracker := testutil.NewTracker()
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		_, err := NewIAMMember(ctx, "test-bad-iam", &IAMMemberArgs{
			ParentID:   pulumi.String("something"),
			ParentType: "unsupported",
			Role:       pulumi.String("roles/viewer"),
			Member:     pulumi.String("user:test@example.com"),
		})
		return err
	}, pulumi.WithMocks("test-project", "test-stack", tracker))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported IAM parent type")
	assert.Contains(t, err.Error(), `"unsupported"`)
}

func TestNewIAMBinding_Legacy_UnsupportedType(t *testing.T) {
	tracker := testutil.NewTracker()
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		_, err := NewIAMBinding(ctx, "test-bad-binding", &IAMBindingArgs{
			ParentID:   pulumi.String("something"),
			ParentType: "bucket",
			Role:       pulumi.String("roles/viewer"),
			Members:    pulumi.StringArray{pulumi.String("user:x@x.com")},
		})
		return err
	}, pulumi.WithMocks("test-project", "test-stack", tracker))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported IAM parent type")
}
