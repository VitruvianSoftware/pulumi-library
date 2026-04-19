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

	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/billing"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/folder"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// IAMMemberArgs configures an additive IAM member binding.
// ParentType is a plain string (not a Pulumi Input) because the scope type is
// always known at plan time. This avoids the ApplyT anti-pattern where
// resources created inside callbacks are invisible to the Pulumi engine.
type IAMMemberArgs struct {
	ParentID   pulumi.StringInput
	ParentType string // "organization", "folder", "project", "serviceAccount", "billing"
	Role       pulumi.StringInput
	Member     pulumi.StringInput
}

type IAMMember struct {
	pulumi.ResourceState
}

// NewIAMMember creates an additive IAM member binding at the specified scope.
// The underlying GCP resource is created directly in the resource graph,
// ensuring proper state tracking, dependency ordering, and error propagation.
func NewIAMMember(ctx *pulumi.Context, name string, args *IAMMemberArgs, opts ...pulumi.ResourceOption) (*IAMMember, error) {
	component := &IAMMember{}
	err := ctx.RegisterComponentResource("pkg:index:IAMMember", name, component, opts...)
	if err != nil {
		return nil, err
	}

	childOpts := append(opts, pulumi.Parent(component))

	switch args.ParentType {
	case "organization":
		_, err = organizations.NewIAMMember(ctx, name+"-member", &organizations.IAMMemberArgs{
			OrgId:  args.ParentID,
			Role:   args.Role,
			Member: args.Member,
		}, childOpts...)
	case "folder":
		_, err = folder.NewIAMMember(ctx, name+"-member", &folder.IAMMemberArgs{
			Folder: args.ParentID,
			Role:   args.Role,
			Member: args.Member,
		}, childOpts...)
	case "project":
		_, err = projects.NewIAMMember(ctx, name+"-member", &projects.IAMMemberArgs{
			Project: args.ParentID,
			Role:    args.Role,
			Member:  args.Member,
		}, childOpts...)
	case "serviceAccount":
		_, err = serviceaccount.NewIAMMember(ctx, name+"-member", &serviceaccount.IAMMemberArgs{
			ServiceAccountId: args.ParentID,
			Role:             args.Role,
			Member:           args.Member,
		}, childOpts...)
	case "billing":
		_, err = billing.NewAccountIamMember(ctx, name+"-member", &billing.AccountIamMemberArgs{
			BillingAccountId: args.ParentID,
			Role:             args.Role,
			Member:           args.Member,
		}, childOpts...)
	default:
		return nil, fmt.Errorf("unsupported IAM parent type: %q (expected organization, folder, project, serviceAccount, or billing)", args.ParentType)
	}

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(component, pulumi.Map{})
	return component, nil
}

// IAMBindingArgs configures an authoritative IAM binding for a specific role.
type IAMBindingArgs struct {
	ParentID   pulumi.StringInput
	ParentType string // "organization", "folder", "project", "serviceAccount", "billing"
	Role       pulumi.StringInput
	Members    pulumi.StringArrayInput
}

type IAMBinding struct {
	pulumi.ResourceState
}

// NewIAMBinding creates an authoritative IAM binding for a specific role at
// the target scope. It will remove any members assigned to this role outside
// of this binding — use with caution.
func NewIAMBinding(ctx *pulumi.Context, name string, args *IAMBindingArgs, opts ...pulumi.ResourceOption) (*IAMBinding, error) {
	component := &IAMBinding{}
	err := ctx.RegisterComponentResource("pkg:index:IAMBinding", name, component, opts...)
	if err != nil {
		return nil, err
	}

	childOpts := append(opts, pulumi.Parent(component))

	switch args.ParentType {
	case "organization":
		_, err = organizations.NewIAMBinding(ctx, name+"-binding", &organizations.IAMBindingArgs{
			OrgId:   args.ParentID,
			Role:    args.Role,
			Members: args.Members,
		}, childOpts...)
	case "folder":
		_, err = folder.NewIAMBinding(ctx, name+"-binding", &folder.IAMBindingArgs{
			Folder:  args.ParentID,
			Role:    args.Role,
			Members: args.Members,
		}, childOpts...)
	case "project":
		_, err = projects.NewIAMBinding(ctx, name+"-binding", &projects.IAMBindingArgs{
			Project: args.ParentID,
			Role:    args.Role,
			Members: args.Members,
		}, childOpts...)
	case "serviceAccount":
		_, err = serviceaccount.NewIAMBinding(ctx, name+"-binding", &serviceaccount.IAMBindingArgs{
			ServiceAccountId: args.ParentID,
			Role:             args.Role,
			Members:          args.Members,
		}, childOpts...)
	case "billing":
		_, err = billing.NewAccountIamBinding(ctx, name+"-binding", &billing.AccountIamBindingArgs{
			BillingAccountId: args.ParentID,
			Role:             args.Role,
			Members:          args.Members,
		}, childOpts...)
	default:
		return nil, fmt.Errorf("unsupported IAM parent type: %q", args.ParentType)
	}

	if err != nil {
		return nil, err
	}

	ctx.RegisterResourceOutputs(component, pulumi.Map{})
	return component, nil
}
