# Vitruvian Software Pulumi Library

Reusable [ComponentResource](https://www.pulumi.com/docs/concepts/resources/components/) building blocks for Google Cloud Platform infrastructure, written in Go.

This library provides enterprise-grade, opinionated components that enforce Google Cloud security best practices. Each component wraps one or more GCP resources into a single, well-tested abstraction with sensible defaults.

## Quick Start

Add to your Go module:

```bash
go get github.com/VitruvianSoftware/pulumi-library
```

Then import the package you need:

```go
import "github.com/VitruvianSoftware/pulumi-library/pkg/project"
```

## Packages

| Package | Import Path | Description | Docs |
|---------|-------------|-------------|------|
| **Project** | `pkg/project` | Project factory: creates GCP projects with API enablement, billing association, and automatic default-VPC suppression | [README](./pkg/project/README.md) |
| **IAM** | `pkg/iam` | Multi-scope IAM bindings (additive + authoritative) at organization, folder, project, service account, and billing account scopes | [README](./pkg/iam/README.md) |
| **Policy** | `pkg/policy` | Organization policy constraint enforcement (boolean + list) using the v2 Org Policy API | [README](./pkg/policy/README.md) |
| **Networking** | `pkg/networking` | VPC networks with subnets (secondary ranges, flow logs, Private Google Access), and optional Private Service Access | [README](./pkg/networking/README.md) |
| **App** | `pkg/app` | Cloud Run v2 service deployment with environment variables, custom service accounts, and ingress control | [README](./pkg/app/README.md) |
| **Data** | `pkg/data` | BigQuery data platform with raw + curated datasets | [README](./pkg/data/README.md) |

## Architecture

Each package follows the same pattern:

```
pkg/<name>/
  └── <name>.go           # Single-file component with Args struct, Component struct, and constructor
  └── README.md            # Package documentation with API reference and examples
```

Every component is a Pulumi [ComponentResource](https://www.pulumi.com/docs/concepts/resources/components/). This means:

- Child resources appear grouped under the component in `pulumi stack`
- The component can be composed into larger abstractions
- Standard Pulumi resource options (`pulumi.Parent`, `pulumi.DependsOn`, `pulumi.Protect`) work on all components

### Resource Graph

```
pkg:index:Project ("seed-project")
├── gcp:organizations:Project ("seed-project")
├── gcp:projects:Service ("seed-project-compute.googleapis.com")
├── gcp:projects:Service ("seed-project-iam.googleapis.com")
└── gcp:projects:Service ("seed-project-cloudkms.googleapis.com")
```

## Design Principles

### 1. Plan-Time Values for Dispatch Fields

Fields that control *which* GCP resource type to create (like `IAMMemberArgs.ParentType` or `ProjectArgs.ActivateApis`) use **plain Go types** (`string`, `[]string`) rather than `pulumi.StringInput`.

**Why:** This ensures resources are registered directly in the Pulumi state graph with proper dependency ordering and error propagation — not inside `ApplyT` callbacks where errors are silently swallowed and resources are invisible to the engine.

```go
// ✅ Correct: ParentType is a plain string
iam.NewIAMMember(ctx, "binding", &iam.IAMMemberArgs{
    ParentType: "organization",              // plain Go string
    ParentID:   pulumi.String(orgID),         // Pulumi Input
    Role:       pulumi.String("roles/viewer"),
    Member:     sa.Email,                     // Pulumi Output
})

// ❌ Wrong: would require ApplyT to resolve the type at runtime
iam.NewIAMMember(ctx, "binding", &iam.IAMMemberArgs{
    ParentType: pulumi.String("organization"), // DON'T do this
    ...
})
```

### 2. Pulumi Inputs for GCP Resource Fields

Fields that map directly to GCP resource arguments (like `ParentID`, `Role`, `Member`) remain `pulumi.StringInput` so they can accept outputs from other resources, enabling proper dependency chains.

### 3. Sensible Security Defaults

- `AutoCreateNetwork` defaults to `false` (the default GCP network has overly permissive firewall rules)
- `PrivateIpGoogleAccess` is always enabled on subnets
- `DisableOnDestroy` is `false` for API services (prevents orphaned APIs)
- Subnets enforce flow logging when `FlowLogs: true` is set

### 4. Single-File Components

Each package is a single Go file. This is intentional — each component is small enough that splitting would add unnecessary navigation cost. As components grow, they can be split while maintaining the same public API.

## Usage Examples

### Create a Project with APIs

```go
import "github.com/VitruvianSoftware/pulumi-library/pkg/project"

p, err := project.NewProject(ctx, "my-project", &project.ProjectArgs{
    ProjectID:      pulumi.String("my-project-id"),
    Name:           pulumi.String("My Project"),
    FolderID:       folderID,  // pulumi.StringOutput from another resource
    BillingAccount: pulumi.String("XXXXXX-XXXXXX-XXXXXX"),
    ActivateApis: []string{
        "compute.googleapis.com",
        "iam.googleapis.com",
        "cloudkms.googleapis.com",
    },
})
// p.Project is the underlying *organizations.Project
// p.Services is a []*projects.Service slice
```

### Bind IAM at Multiple Scopes

```go
import "github.com/VitruvianSoftware/pulumi-library/pkg/iam"

// Organization-level (additive)
iam.NewIAMMember(ctx, "org-admin", &iam.IAMMemberArgs{
    ParentID:   pulumi.String(orgID),
    ParentType: "organization",
    Role:       pulumi.String("roles/resourcemanager.organizationAdmin"),
    Member:     pulumi.Sprintf("serviceAccount:%s", sa.Email),
})

// Project-level (additive)
iam.NewIAMMember(ctx, "project-editor", &iam.IAMMemberArgs{
    ParentID:   p.Project.ProjectId,
    ParentType: "project",
    Role:       pulumi.String("roles/editor"),
    Member:     pulumi.String("user:admin@example.com"),
})

// Billing-level (additive)
iam.NewIAMMember(ctx, "billing-user", &iam.IAMMemberArgs{
    ParentID:   pulumi.String("XXXXXX-XXXXXX-XXXXXX"),
    ParentType: "billing",
    Role:       pulumi.String("roles/billing.user"),
    Member:     pulumi.Sprintf("serviceAccount:%s", sa.Email),
})

// Authoritative — REMOVES members not in this list
iam.NewIAMBinding(ctx, "project-viewers", &iam.IAMBindingArgs{
    ParentID:   p.Project.ProjectId,
    ParentType: "project",
    Role:       pulumi.String("roles/viewer"),
    Members: pulumi.StringArray{
        pulumi.String("user:alice@example.com"),
        pulumi.String("user:bob@example.com"),
    },
})
```

### Enforce Organization Policies

```go
import "github.com/VitruvianSoftware/pulumi-library/pkg/policy"

// Boolean constraint — enforce
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

### Create a VPC with Subnets and PSA

```go
import "github.com/VitruvianSoftware/pulumi-library/pkg/networking"

net, err := networking.NewNetworking(ctx, "my-vpc", &networking.NetworkingArgs{
    ProjectID: p.Project.ProjectId,
    VPCName:   pulumi.String("vpc-shared-base"),
    EnablePSA: true,
    Subnets: []networking.SubnetArgs{
        {
            Name:     "sb-us-central1",
            Region:   "us-central1",
            CIDR:     "10.0.0.0/21",
            FlowLogs: true,
            SecondaryRanges: []networking.SecondaryRangeArgs{
                {RangeName: "gke-pods", CIDR: "100.64.0.0/21"},
                {RangeName: "gke-svcs", CIDR: "100.64.8.0/21"},
            },
        },
    },
})
// net.VPC is the underlying *compute.Network
// net.Subnets is a map[string]*compute.Subnetwork
```

### Deploy a Cloud Run Service

```go
import "github.com/VitruvianSoftware/pulumi-library/pkg/app"

cloudRun, err := app.NewCloudRunApp(ctx, "my-api", &app.CloudRunAppArgs{
    ProjectID: p.Project.ProjectId,
    Name:      pulumi.String("my-api"),
    Image:     pulumi.String("us-docker.pkg.dev/my-project/my-repo/my-image:latest"),
    Region:    pulumi.String("us-central1"),
    EnvVars: map[string]pulumi.StringInput{
        "LOG_LEVEL":  pulumi.String("info"),
        "PROJECT_ID": p.Project.ProjectId,
    },
})
// cloudRun.Service is the underlying *cloudrunv2.Service
```

### Create a BigQuery Data Platform

```go
import "github.com/VitruvianSoftware/pulumi-library/pkg/data"

dp, err := data.NewDataPlatform(ctx, "analytics", &data.DataPlatformArgs{
    ProjectID:        p.Project.ProjectId,
    Location:         pulumi.String("US"),
    RawDatasetID:     "raw_events",     // optional, defaults to "raw_data"
    CuratedDatasetID: "curated_events", // optional, defaults to "curated_data"
})
// dp.RawDataset is the underlying *bigquery.Dataset
// dp.CuratedDataset is the underlying *bigquery.Dataset
```

## Compatibility

| Dependency | Version |
|------------|---------|
| Go | 1.21+ |
| Pulumi SDK | v3.231.0+ |
| Pulumi GCP Provider | v9.20.0+ |

## Development

```bash
# Build all packages
make build

# Run tests
make test

# Lint and format
make lint

# Tidy module dependencies
make tidy
```

## Related

- [Pulumi Example Foundation](https://github.com/VitruvianSoftware/pulumi-example-foundation) — Enterprise GCP foundation that consumes this library
- [Google Terraform Example Foundation](https://github.com/terraform-google-modules/terraform-example-foundation) — The upstream reference architecture
- [Pulumi ComponentResources](https://www.pulumi.com/docs/concepts/resources/components/) — Pulumi's documentation on component resources

## License

Apache License 2.0 — see [LICENSE](./LICENSE) for details.
