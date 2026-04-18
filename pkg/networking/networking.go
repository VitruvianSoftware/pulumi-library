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

package networking

import (
	"fmt"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/servicenetworking"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type SubnetArgs struct {
	Name   string
	Region string
	CIDR   string
}

type NetworkingArgs struct {
	ProjectID   pulumi.StringInput
	VPCName     pulumi.StringInput
	Subnets     []SubnetArgs
	EnablePSA   bool // Private Service Access
}

type Networking struct {
	pulumi.ResourceState
	VPC     *compute.Network
	Subnets map[string]*compute.Subnetwork
}

func NewNetworking(ctx *pulumi.Context, name string, args *NetworkingArgs, opts ...pulumi.ResourceOption) (*Networking, error) {
	component := &Networking{
		Subnets: make(map[string]*compute.Subnetwork),
	}
	err := ctx.RegisterComponentResource("pkg:index:Networking", name, component, opts...)
	if err != nil {
		return nil, err
	}

	// 1. VPC
	vpc, err := compute.NewNetwork(ctx, name+"-vpc", &compute.NetworkArgs{
		Project:               args.ProjectID,
		Name:                  args.VPCName,
		AutoCreateSubnetworks: pulumi.Bool(false),
	}, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}
	component.VPC = vpc

	// 2. Subnets
	for _, s := range args.Subnets {
		sub, err := compute.NewSubnetwork(ctx, name+"-"+s.Name, &compute.SubnetworkArgs{
			Project:     args.ProjectID,
			Name:        pulumi.String(s.Name),
			Region:      pulumi.String(s.Region),
			Network:     vpc.ID(),
			IpCidrRange: pulumi.String(s.CIDR),
		}, pulumi.Parent(vpc))
		if err != nil {
			return nil, err
		}
		component.Subnets[s.Name] = sub
	}

	// 3. Private Service Access
	if args.EnablePSA {
		reservedIP, err := compute.NewGlobalAddress(ctx, name+"-psa-ip", &compute.GlobalAddressArgs{
			Project:      args.ProjectID,
			Name:         pulumi.String(name + "-psa-range"),
			Purpose:      pulumi.String("VPC_PEERING"),
			AddressType:  pulumi.String("INTERNAL"),
			PrefixLength: pulumi.Int(16),
			Network:      vpc.ID(),
		}, pulumi.Parent(vpc))
		if err != nil {
			return nil, err
		}

		_, err = servicenetworking.NewConnection(ctx, name+"-psa-conn", &servicenetworking.ConnectionArgs{
			Network:                vpc.ID(),
			Service:                pulumi.String("servicenetworking.googleapis.com"),
			ReservedPeeringRanges: pulumi.StringArray{reservedIP.Name},
		}, pulumi.Parent(vpc))
		if err != nil {
			return nil, err
		}
	}

	ctx.RegisterResourceOutputs(component, pulumi.Map{
		"vpcId": vpc.ID(),
	})

	return component, nil
}
