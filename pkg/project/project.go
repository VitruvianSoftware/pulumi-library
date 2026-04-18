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

// ProjectArgs configures the Project component.
// ActivateApis is a plain []string (not a Pulumi Input) because API names are
// always known at plan time. This ensures each projects.Service resource is
// properly registered in the Pulumi state graph — NOT created inside an
// ApplyT callback where errors are silently swallowed and resources are
// invisible to the engine.
type ProjectArgs struct {
	ProjectID         pulumi.StringInput
	Name              pulumi.StringInput
	FolderID          pulumi.StringInput
	BillingAccount    pulumi.StringInput
	ActivateApis      []string // plain Go slice — always known at plan time
	AutoCreateNetwork pulumi.BoolPtrInput
	Labels            pulumi.StringMapInput
	DeletionPolicy    pulumi.StringPtrInput
}

type Project struct {
	pulumi.ResourceState
	Project  *organizations.Project
	Services []*projects.Service
}

func NewProject(ctx *pulumi.Context, name string, args *ProjectArgs, opts ...pulumi.ResourceOption) (*Project, error) {
	component := &Project{}
	err := ctx.RegisterComponentResource("pkg:index:Project", name, component, opts...)
	if err != nil {
		return nil, err
	}

	// Default to false for autoCreateNetwork — security best practice:
	// the default VPC has overly permissive firewall rules.
	autoCreateNetwork := args.AutoCreateNetwork
	if autoCreateNetwork == nil {
		autoCreateNetwork = pulumi.Bool(false)
	}

	// 1. Create the Project
	pArgs := &organizations.ProjectArgs{
		ProjectId:         args.ProjectID,
		Name:              args.Name,
		FolderId:          args.FolderID,
		BillingAccount:    args.BillingAccount,
		AutoCreateNetwork: autoCreateNetwork,
		Labels:            args.Labels,
	}

	if args.DeletionPolicy != nil {
		pArgs.DeletionPolicy = args.DeletionPolicy
	}

	p, err := organizations.NewProject(ctx, name, pArgs, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}
	component.Project = p

	// 2. Enable APIs — each Service is a first-class Pulumi resource,
	// properly tracked in state with correct dependency ordering.
	for _, api := range args.ActivateApis {
		svc, err := projects.NewService(ctx, fmt.Sprintf("%s-%s", name, api), &projects.ServiceArgs{
			Project:                  p.ProjectId,
			Service:                  pulumi.String(api),
			DisableOnDestroy:         pulumi.Bool(false),
			DisableDependentServices: pulumi.Bool(false),
		}, pulumi.Parent(p))
		if err != nil {
			return nil, err
		}
		component.Services = append(component.Services, svc)
	}

	ctx.RegisterResourceOutputs(component, pulumi.Map{
		"projectId": p.ProjectId,
	})

	return component, nil
}
