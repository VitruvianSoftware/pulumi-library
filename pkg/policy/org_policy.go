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

package policy

import (
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/orgpolicy"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type OrgPolicyArgs struct {
	ParentID   pulumi.StringInput
	Constraint pulumi.StringInput
	Boolean    pulumi.BoolPtrInput
	List       *OrgPolicyListArgs
}

type OrgPolicyListArgs struct {
	Allow []string
	Deny  []string
}

type OrgPolicy struct {
	pulumi.ResourceState
}

func NewOrgPolicy(ctx *pulumi.Context, name string, args *OrgPolicyArgs, opts ...pulumi.ResourceOption) (*OrgPolicy, error) {
	component := &OrgPolicy{}
	err := ctx.RegisterComponentResource("pkg:index:OrgPolicy", name, component, opts...)
	if err != nil {
		return nil, err
	}

	policyArgs := &orgpolicy.PolicyArgs{
		Parent:     args.ParentID,
		Name:       args.Constraint,
	}

	if args.Boolean != nil {
		policyArgs.Spec = &orgpolicy.PolicySpecArgs{
			Rules: orgpolicy.PolicySpecRuleArray{
				&orgpolicy.PolicySpecRuleArgs{
					Enforce: args.Boolean.ToBoolPtrOutput().ApplyT(func(b *bool) string {
						if b != nil && *b {
							return "TRUE"
						}
						return "FALSE"
					}).(pulumi.StringOutput),
				},
			},
		}
	}

	_, err = orgpolicy.NewPolicy(ctx, name, policyArgs, pulumi.Parent(component))
	return component, err
}
