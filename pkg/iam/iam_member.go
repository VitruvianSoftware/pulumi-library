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
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// IAMMemberArgs is the legacy unified args struct.
//
// Deprecated: Use the scope-specific constructors instead:
//   - iam.NewOrganizationIAMMember
//   - iam.NewFolderIAMMember
//   - iam.NewProjectIAMMember
//   - iam.NewServiceAccountIAMMember
//   - iam.NewBillingIAMMember
type IAMMemberArgs struct {
	ParentID   pulumi.StringInput
	ParentType string // "organization", "folder", "project", "serviceAccount", "billing"
	Role       pulumi.StringInput
	Member     pulumi.StringInput
}

// IAMMember is the legacy unified component.
//
// Deprecated: Use the scope-specific types instead.
type IAMMember struct {
	pulumi.ResourceState
}

// NewIAMMember creates an additive IAM member binding at the specified scope.
//
// Deprecated: Use the scope-specific constructors instead:
//   - iam.NewOrganizationIAMMember
//   - iam.NewFolderIAMMember
//   - iam.NewProjectIAMMember
//   - iam.NewServiceAccountIAMMember
//   - iam.NewBillingIAMMember
func NewIAMMember(ctx *pulumi.Context, name string, args *IAMMemberArgs, opts ...pulumi.ResourceOption) (*IAMMember, error) {
	component := &IAMMember{}
	err := ctx.RegisterComponentResource("pkg:index:IAMMember", name, component, opts...)
	if err != nil {
		return nil, err
	}

	switch args.ParentType {
	case "organization":
		_, err = NewOrganizationIAMMember(ctx, name+"-delegated", &OrganizationIAMMemberArgs{
			OrgID:  args.ParentID,
			Role:   args.Role,
			Member: args.Member,
		}, pulumi.Parent(component))
	case "folder":
		_, err = NewFolderIAMMember(ctx, name+"-delegated", &FolderIAMMemberArgs{
			FolderID: args.ParentID,
			Role:     args.Role,
			Member:   args.Member,
		}, pulumi.Parent(component))
	case "project":
		_, err = NewProjectIAMMember(ctx, name+"-delegated", &ProjectIAMMemberArgs{
			ProjectID: args.ParentID,
			Role:      args.Role,
			Member:    args.Member,
		}, pulumi.Parent(component))
	case "serviceAccount":
		_, err = NewServiceAccountIAMMember(ctx, name+"-delegated", &ServiceAccountIAMMemberArgs{
			ServiceAccountID: args.ParentID,
			Role:             args.Role,
			Member:           args.Member,
		}, pulumi.Parent(component))
	case "billing":
		_, err = NewBillingIAMMember(ctx, name+"-delegated", &BillingIAMMemberArgs{
			BillingAccountID: args.ParentID,
			Role:             args.Role,
			Member:           args.Member,
		}, pulumi.Parent(component))
	default:
		return nil, fmt.Errorf("unsupported IAM parent type: %q (expected organization, folder, project, serviceAccount, or billing)", args.ParentType)
	}

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(component, pulumi.Map{})
	return component, nil
}

// IAMBindingArgs is the legacy unified args struct for authoritative bindings.
//
// Deprecated: Use the scope-specific constructors instead:
//   - iam.NewOrganizationIAMBinding
//   - iam.NewFolderIAMBinding
//   - iam.NewProjectIAMBinding
//   - iam.NewServiceAccountIAMBinding
//   - iam.NewBillingIAMBinding
type IAMBindingArgs struct {
	ParentID   pulumi.StringInput
	ParentType string // "organization", "folder", "project", "serviceAccount", "billing"
	Role       pulumi.StringInput
	Members    pulumi.StringArrayInput
}

// IAMBinding is the legacy unified component for authoritative bindings.
//
// Deprecated: Use the scope-specific types instead.
type IAMBinding struct {
	pulumi.ResourceState
}

// NewIAMBinding creates an authoritative IAM binding for a specific role at
// the target scope. It will remove any members assigned to this role outside
// of this binding — use with caution.
//
// Deprecated: Use the scope-specific constructors instead:
//   - iam.NewOrganizationIAMBinding
//   - iam.NewFolderIAMBinding
//   - iam.NewProjectIAMBinding
//   - iam.NewServiceAccountIAMBinding
//   - iam.NewBillingIAMBinding
func NewIAMBinding(ctx *pulumi.Context, name string, args *IAMBindingArgs, opts ...pulumi.ResourceOption) (*IAMBinding, error) {
	component := &IAMBinding{}
	err := ctx.RegisterComponentResource("pkg:index:IAMBinding", name, component, opts...)
	if err != nil {
		return nil, err
	}

	switch args.ParentType {
	case "organization":
		_, err = NewOrganizationIAMBinding(ctx, name+"-delegated", &OrganizationIAMBindingArgs{
			OrgID:   args.ParentID,
			Role:    args.Role,
			Members: args.Members,
		}, pulumi.Parent(component))
	case "folder":
		_, err = NewFolderIAMBinding(ctx, name+"-delegated", &FolderIAMBindingArgs{
			FolderID: args.ParentID,
			Role:     args.Role,
			Members:  args.Members,
		}, pulumi.Parent(component))
	case "project":
		_, err = NewProjectIAMBinding(ctx, name+"-delegated", &ProjectIAMBindingArgs{
			ProjectID: args.ParentID,
			Role:      args.Role,
			Members:   args.Members,
		}, pulumi.Parent(component))
	case "serviceAccount":
		_, err = NewServiceAccountIAMBinding(ctx, name+"-delegated", &ServiceAccountIAMBindingArgs{
			ServiceAccountID: args.ParentID,
			Role:             args.Role,
			Members:          args.Members,
		}, pulumi.Parent(component))
	case "billing":
		_, err = NewBillingIAMBinding(ctx, name+"-delegated", &BillingIAMBindingArgs{
			BillingAccountID: args.ParentID,
			Role:             args.Role,
			Members:          args.Members,
		}, pulumi.Parent(component))
	default:
		return nil, fmt.Errorf("unsupported IAM parent type: %q", args.ParentType)
	}

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(component, pulumi.Map{})
	return component, nil
}
