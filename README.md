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
| **Bootstrap** | `pkg/bootstrap` | Core foundation seed project, KMS keys/rings, encrypted state buckets, and base organization policies | [README](./pkg/bootstrap/README.md) |
| **Project** | `pkg/project` | Project factory: creates GCP projects with API enablement, billing association, and automatic default-VPC suppression | [README](./pkg/project/README.md) |
| **Group** | `pkg/group` | Google Workspace / Cloud Identity group provisioning with structured ownership and dynamic typing | [README](./pkg/group/README.md) |
| **IAM** | `pkg/iam` | Scope-isolated IAM bindings (additive + authoritative) with dedicated constructors per GCP scope: organization, folder, project, service account, and billing account | [README](./pkg/iam/README.md) |
| **Policy** | `pkg/policy` | Organization policy constraint enforcement (boolean + list) using the v2 Org Policy API | [README](./pkg/policy/README.md) |
| **Logging** | `pkg/logging` | Centralized log export infrastructure with org/folder sinks, internal project sinks, and destinations | [README](./pkg/logging/README.md) |
| **Networking** | `pkg/networking` | VPC networks with subnets (secondary ranges, flow logs, Private Google Access), and optional Private Service Access | [README](./pkg/networking/README.md) |
| **App** | `pkg/app` | Cloud Run v2 service deployment with environment variables, custom service accounts, and ingress control | [README](./pkg/app/README.md) |
| **Data** | `pkg/data` | BigQuery data platform with raw + curated datasets | [README](./pkg/data/README.md) |
| **CI/CD** | `pkg/cicd` | Workload Identity Federation (WIF) integrations for external pipelines (GitHub Actions, GitLab CI) and Cloud Build infrastructure (Source Repos, Artifact Registry, Triggers) | [README](./pkg/cicd/README.md) |
| **Storage** | `pkg/storage` | Hardened Google Cloud Storage buckets with enforced public access prevention and optional KMS | [README](./pkg/storage/README.md) |

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

### 1. Scope-Isolated Constructors

When a package operates across multiple GCP scopes (e.g., IAM at organization, folder, project levels), each scope gets a **dedicated constructor** with a scope-specific `Args` struct. This replaces the earlier strategy-pattern approach that used a `ParentType` string for runtime dispatch.

**Why:** Scope isolation provides compile-time safety (no magic strings), independent Pulumi component types (blast radius isolation), and Args structs that contain only scope-relevant fields.

```go
// ✅ Correct: scope-specific constructor with typed Args
iam.NewOrganizationIAMMember(ctx, "org-admin", &iam.OrganizationIAMMemberArgs{
    OrgID:  pulumi.String(orgID),
    Role:   pulumi.String("roles/viewer"),
    Member: sa.Email,
})

// ❌ Deprecated: unified constructor with magic string dispatch
iam.NewIAMMember(ctx, "org-admin", &iam.IAMMemberArgs{
    ParentType: "organization",  // magic string, runtime error on typo
    ParentID:   pulumi.String(orgID),
    ...
})
```

### 2. Plan-Time Values for Dispatch Fields

Fields that control *which* GCP resource type to create (like `ProjectArgs.ActivateApis`) use **plain Go types** (`string`, `[]string`) rather than `pulumi.StringInput`.

**Why:** This ensures resources are registered directly in the Pulumi state graph with proper dependency ordering and error propagation — not inside `ApplyT` callbacks where errors are silently swallowed and resources are invisible to the engine.

### 3. Pulumi Inputs for GCP Resource Fields

Fields that map directly to GCP resource arguments (like `OrgID`, `Role`, `Member`) remain `pulumi.StringInput` so they can accept outputs from other resources, enabling proper dependency chains.

### 4. Sensible Security Defaults

- `AutoCreateNetwork` defaults to `false` (the default GCP network has overly permissive firewall rules)
- `PrivateIpGoogleAccess` is always enabled on subnets
- `DisableOnDestroy` is `false` for API services (prevents orphaned APIs)
- Subnets enforce flow logging when `FlowLogs: true` is set

### 5. File-Per-Scope Within Packages

Simple packages use a single Go file. Packages that span multiple GCP scopes (like `pkg/iam`) use one file per scope (e.g., `organization.go`, `project.go`, `billing.go`) to keep each scope's logic independent and reviewable.

## Usage Examples

### Bootstrap the Foundation

```go
import "github.com/VitruvianSoftware/pulumi-library/pkg/bootstrap"

seed, err := bootstrap.NewBootstrap(ctx, "foundation-seed", &bootstrap.BootstrapArgs{
    OrgID:          "1234567890",
    BillingAccount: "XXXXXX-XXXXXX-XXXXXX",
    ProjectPrefix:  "prj",
    DefaultRegion:  "us-central1",
    StateBucketIAMMembers: []pulumi.StringInput{
        pulumi.String("group:gcp-organization-admins@example.com"),
    },
})
// seed.SeedProject is the underlying project component
// seed.StateBucketName is the generated KMS-encrypted GCS bucket
```

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
iam.NewOrganizationIAMMember(ctx, "org-admin", &iam.OrganizationIAMMemberArgs{
    OrgID:  pulumi.String(orgID),
    Role:   pulumi.String("roles/resourcemanager.organizationAdmin"),
    Member: pulumi.Sprintf("serviceAccount:%s", sa.Email),
})

// Project-level (additive)
iam.NewProjectIAMMember(ctx, "project-editor", &iam.ProjectIAMMemberArgs{
    ProjectID: p.Project.ProjectId,
    Role:      pulumi.String("roles/editor"),
    Member:    pulumi.String("user:admin@example.com"),
})

// Billing-level (additive)
iam.NewBillingIAMMember(ctx, "billing-user", &iam.BillingIAMMemberArgs{
    BillingAccountID: pulumi.String("XXXXXX-XXXXXX-XXXXXX"),
    Role:             pulumi.String("roles/billing.user"),
    Member:           pulumi.Sprintf("serviceAccount:%s", sa.Email),
})

// Authoritative — REMOVES members not in this list
iam.NewProjectIAMBinding(ctx, "project-viewers", &iam.ProjectIAMBindingArgs{
    ProjectID: p.Project.ProjectId,
    Role:      pulumi.String("roles/viewer"),
    Members: pulumi.StringArray{
        pulumi.String("user:alice@example.com"),
        pulumi.String("user:bob@example.com"),
    },
})
```

### Provision Google Workspace Groups

```go
import "github.com/VitruvianSoftware/pulumi-library/pkg/group"

g, err := group.NewGroup(ctx, "org-admins", &group.GroupArgs{
    ID:          "gcp-org-admins@example.com",
    DisplayName: "GCP Organization Admins",
    CustomerID:  pulumi.String("C01234abc"),
    Types:       []string{"default", "security"},
    Owners:      []string{"admin-owner@example.com"},
    Managers:    []string{"admin-manager@example.com"},
    Members:     []string{"admin-user@example.com"},
})
// g.GroupResource is the underlying Cloud Identity group
// g.GroupEmail is the managed group email
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

### Export Centralized Logs

```go
import "github.com/VitruvianSoftware/pulumi-library/pkg/logging"

cl, err := logging.NewCentralizedLogging(ctx, "logs-export", &logging.CentralizedLoggingArgs{
    Resources:                   map[string]string{"resource": "organizations/1234567890"},
    ResourceType:                "organization",
    LoggingDestinationProjectID: pulumi.String("my-audit-project"),
    StorageOptions: &logging.StorageOptions{
        LoggingSinkName:   "sk-c-logging-bkt",
        LoggingSinkFilter: "logName: /logs/cloudaudit.googleapis.com%2Factivity",
        Location:          "US",
    },
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

### Configure CI/CD Workload Identity (GitHub/GitLab)

```go
import "github.com/VitruvianSoftware/pulumi-library/pkg/cicd"

// Configure GitHub Actions OIDC
gh, err := cicd.NewGitHubOIDC(ctx, "gh-oidc", &cicd.GitHubOIDCArgs{
    ProjectID:          pulumi.String("my-cicd-project"),
    PoolID:             pulumi.String("foundation-pool"),
    ProviderID:         pulumi.String("foundation-gh-provider"),
    AttributeCondition: pulumi.String("assertion.repository_owner=='my-org'"),
    SAMapping: map[string]cicd.SAMappingEntry{
        "bootstrap": {
            SAName:    pulumi.String("projects/my-cicd-project/serviceAccounts/bootstrap@..."),
            Attribute: pulumi.String("attribute.repository/my-org/gcp-bootstrap"),
        },
    },
})
```

### Create a Secure Cloud Storage Bucket

```go
import "github.com/VitruvianSoftware/pulumi-library/pkg/storage"

enabled := true
bucket, err := storage.NewSimpleBucket(ctx, "state-bucket", &storage.SimpleBucketArgs{
    Name:         pulumi.String("my-encrypted-state-bucket"),
    ProjectID:    pulumi.String("my-seed-project"),
    Location:     pulumi.String("us-central1"),
    ForceDestroy: pulumi.Bool(false),
    Versioning:   &enabled,
})
```

### Provision Cloud Build Infrastructure

```go
import "github.com/VitruvianSoftware/pulumi-library/pkg/cicd"

cb, err := cicd.NewCloudBuild(ctx, "pipeline", &cicd.CloudBuildArgs{
    ProjectID:  pulumi.String("my-cicd-project"),
    Region:     pulumi.String("us-central1"),
    SourceType: cicd.CloudBuildSourceGitHub, // Default; CSR is deprecated
    Triggers: map[string]cicd.CloudBuildTriggerConfig{
        "bootstrap": {
            RepoName:       "gcp-bootstrap",
            RepoOwner:      "my-org",
            ServiceAccount: pulumi.String("projects/my-cicd-project/serviceAccounts/sa@..."),
        },
    },
})
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
