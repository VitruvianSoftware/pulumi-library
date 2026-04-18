# Vitruvian Software Pulumi Library

Reusable [ComponentResource](https://www.pulumi.com/docs/concepts/resources/components/) building blocks for GCP infrastructure, written in Go.

## Packages

| Package | Description |
|---------|-------------|
| `pkg/project` | Project factory: creates GCP projects with API enablement, billing, and labels |
| `pkg/iam` | IAM bindings at organization, folder, project, service account, and billing scopes |
| `pkg/policy` | Organization policy enforcement (boolean + list constraints) |
| `pkg/networking` | VPC, subnets (with secondary ranges, flow logs), and Private Service Access |
| `pkg/app` | Cloud Run v2 service deployment |
| `pkg/data` | BigQuery data platform (raw + curated datasets) |

## Design Principles

**Plan-time values over Pulumi Inputs for dispatch fields.** Fields that control *which* GCP resource type to create (like `IAMMemberArgs.ParentType` or `ProjectArgs.ActivateApis`) use plain Go types (`string`, `[]string`) rather than `pulumi.StringInput`. This ensures resources are registered directly in the Pulumi state graph with proper dependency ordering and error propagation — not inside `ApplyT` callbacks where errors are silently swallowed.

**Pulumi Inputs for GCP resource fields.** Fields that map to GCP resource arguments (like `ParentID`, `Role`, `Member`) remain `pulumi.StringInput` so they can accept outputs from other resources.

## Usage

Add to your Go module:
```bash
go get github.com/VitruvianSoftware/pulumi-library
```

### IAM (Additive vs. Authoritative)

```go
// Additive — adds a member without affecting others
iam.NewIAMMember(ctx, "my-binding", &iam.IAMMemberArgs{
    ParentID:   pulumi.String(orgID),
    ParentType: "organization",  // plain string, not pulumi.String()
    Role:       pulumi.String("roles/viewer"),
    Member:     pulumi.String("user:alice@example.com"),
})

// Authoritative — removes members NOT in this list
iam.NewIAMBinding(ctx, "my-binding", &iam.IAMBindingArgs{
    ParentID:   pulumi.String(orgID),
    ParentType: "organization",
    Role:       pulumi.String("roles/viewer"),
    Members:    pulumi.StringArray{pulumi.String("user:alice@example.com")},
})
```

### Organization Policy

```go
// Boolean constraint
policy.NewOrgPolicy(ctx, "no-serial-port", &policy.OrgPolicyArgs{
    ParentID:   pulumi.String("organizations/123456"),
    Constraint: pulumi.String("constraints/compute.disableSerialPortAccess"),
    Boolean:    pulumi.Bool(true),
})

// List constraint — deny all
policy.NewOrgPolicy(ctx, "no-external-ip", &policy.OrgPolicyArgs{
    ParentID:   pulumi.String("organizations/123456"),
    Constraint: pulumi.String("constraints/compute.vmExternalIpAccess"),
    DenyAll:    pulumi.Bool(true),
})

// List constraint — allow specific values
policy.NewOrgPolicy(ctx, "restrict-domains", &policy.OrgPolicyArgs{
    ParentID:    pulumi.String("organizations/123456"),
    Constraint:  pulumi.String("constraints/iam.allowedPolicyMemberDomains"),
    AllowValues: pulumi.StringArray{pulumi.String("C0xxxxxxx")},
})
```

## Contributing

This is a public reference repository. Please open a PR for any improvements.
