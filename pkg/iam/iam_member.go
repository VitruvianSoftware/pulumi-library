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
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type IAMMemberArgs struct {
	ParentID   pulumi.StringInput
	ParentType pulumi.StringInput // "organization", "folder", "project"
	Role       pulumi.StringInput
	Member     pulumi.StringInput
}

type IAMMember struct {
	pulumi.ResourceState
}

func NewIAMMember(ctx *pulumi.Context, name string, args *IAMMemberArgs, opts ...pulumi.ResourceOption) (*IAMMember, error) {
	component := &IAMMember{}
	err := ctx.RegisterComponentResource("pkg:index:IAMMember", name, component, opts...)
	if err != nil {
		return nil, err
	}

	args.ParentType.ToDetailedStringOutput().ApplyT(func(pType string) error {
		var err error
		switch pType {
		case "organization":
			_, err = organizations.NewIAMMember(ctx, name, &organizations.IAMMemberArgs{
				OrgId:  args.ParentID,
				Role:   args.Role,
				Member: args.Member,
			}, pulumi.Parent(component))
		case "folder":
			_, err = organizations.NewFolderIamMember(ctx, name, &organizations.FolderIamMemberArgs{
				Folder: args.ParentID,
				Role:   args.Role,
				Member: args.Member,
			}, pulumi.Parent(component))
		case "project":
			_, err = projects.NewIAMMember(ctx, name, &projects.IAMMemberArgs{
				Project: args.ParentID,
				Role:    args.Role,
				Member:  args.Member,
			}, pulumi.Parent(component))
		}
		return err
	})

	return component, nil
}
