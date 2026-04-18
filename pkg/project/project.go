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

package project

import (
	"fmt"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ProjectArgs struct {
	ProjectID      pulumi.StringInput
	Name           pulumi.StringInput
	FolderID       pulumi.StringInput
	BillingAccount pulumi.StringInput
	ActivateApis   pulumi.StringArrayInput
}

type Project struct {
	pulumi.ResourceState
	Project *organizations.Project
}

func NewProject(ctx *pulumi.Context, name string, args *ProjectArgs, opts ...pulumi.ResourceOption) (*Project, error) {
	component := &Project{}
	err := ctx.RegisterComponentResource("pkg:index:Project", name, component, opts...)
	if err != nil {
		return nil, err
	}

	// 1. Create the Project
	p, err := organizations.NewProject(ctx, name, &organizations.ProjectArgs{
		ProjectId:      args.ProjectID,
		Name:           args.Name,
		FolderId:       args.FolderID,
		BillingAccount: args.BillingAccount,
	}, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}
	component.Project = p

	// 2. Enable APIs
	if args.ActivateApis != nil {
		args.ActivateApis.ToStringArrayOutput().ApplyT(func(apis []string) error {
			for _, api := range apis {
				_, err := projects.NewService(ctx, fmt.Sprintf("%s-%s", name, api), &projects.ServiceArgs{
					Project: p.ProjectId,
					Service: pulumi.String(api),
				}, pulumi.Parent(p))
				if err != nil {
					return err
				}
			}
			return nil
		})
	}

	ctx.RegisterResourceOutputs(component, pulumi.Map{
		"projectId": p.ProjectId,
	})

	return component, nil
}
