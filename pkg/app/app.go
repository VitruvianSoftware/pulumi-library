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

package app

import (
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/cloudrun"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CloudRunAppArgs struct {
	ProjectID pulumi.StringInput
	Name      pulumi.StringInput
	Image     pulumi.StringInput
	Region    pulumi.StringInput
}

type CloudRunApp struct {
	pulumi.ResourceState
	Service *cloudrun.Service
}

func NewCloudRunApp(ctx *pulumi.Context, name string, args *CloudRunAppArgs, opts ...pulumi.ResourceOption) (*CloudRunApp, error) {
	component := &CloudRunApp{}
	err := ctx.RegisterComponentResource("pkg:index:CloudRunApp", name, component, opts...)
	if err != nil {
		return nil, err
	}

	svc, err := cloudrun.NewService(ctx, name, &cloudrun.ServiceArgs{
		Name:     args.Name,
		Project:  args.ProjectID,
		Location: args.Region,
		Template: &cloudrun.ServiceTemplateArgs{
			Spec: &cloudrun.ServiceTemplateSpecArgs{
				Containers: cloudrun.ServiceTemplateSpecContainerArray{
					&cloudrun.ServiceTemplateSpecContainerArgs{
						Image: args.Image,
					},
				},
			},
		},
	}, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}

	component.Service = svc
	return component, nil
}
