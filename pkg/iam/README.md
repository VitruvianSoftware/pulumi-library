# pkg/iam — Multi-Scope IAM Bindings

Additive and authoritative IAM bindings across 5 GCP scopes: organization, folder, project, service account, and billing account.

## Overview

The `IAMMember` component creates **additive** IAM bindings (adds a member without affecting others). The `IAMBinding` component creates **authoritative** bindings for a specific role (removes any members not in the list).

Both components use a `ParentType` dispatch pattern with a plain Go `string` to determine which underlying GCP IAM resource to create. This ensures the resource is registered directly in Pulumi's state graph, not inside an `ApplyT` callback.

### Additive vs. Authoritative

| Mode | Component | Behavior | When to Use |
|------|-----------|----------|-------------|
| **Additive** | `IAMMember` | Adds a single member to a role; does not affect other members | Most environments — safe to use alongside manual IAM |
| **Authoritative** | `IAMBinding` | Sets the complete list of members for a role; **removes unlisted members** | Strict environments where IAM must be fully managed by code |

> ⚠️ **Warning:** `IAMBinding` will remove any members assigned to the specified role that are not included in the `Members` list. Use with extreme caution in production.

## Supported Scopes

| `ParentType` | `ParentID` Expected | Underlying GCP Resource |
|--------------|--------------------|-----------------------|
| `"organization"` | Organization ID (numeric) | `organizations.IAMMember` / `IAMBinding` |
| `"folder"` | Folder ID (numeric or `folders/ID`) | `organizations.FolderIamMember` / `FolderIamBinding` |
| `"project"` | Project ID (string) | `projects.IAMMember` / `IAMBinding` |
| `"serviceAccount"` | Service account ID (email or resource name) | `serviceaccount.IAMMember` / `IAMBinding` |
| `"billing"` | Billing account ID (`XXXXXX-XXXXXX-XXXXXX`) | `billing.AccountIamMember` / `AccountIamBinding` |

## API Reference

### `IAMMemberArgs` (Additive)

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| `ParentID` | `pulumi.StringInput` | ✅ | The ID of the parent resource |
| `ParentType` | `string` (plain) | ✅ | One of: `organization`, `folder`, `project`, `serviceAccount`, `billing` |
| `Role` | `pulumi.StringInput` | ✅ | The IAM role to grant |
| `Member` | `pulumi.StringInput` | ✅ | The member identity (e.g., `user:alice@example.com`, `serviceAccount:sa@project.iam.gserviceaccount.com`) |

### `IAMBindingArgs` (Authoritative)

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| `ParentID` | `pulumi.StringInput` | ✅ | The ID of the parent resource |
| `ParentType` | `string` (plain) | ✅ | One of: `organization`, `folder`, `project`, `serviceAccount`, `billing` |
| `Role` | `pulumi.StringInput` | ✅ | The IAM role to bind |
| `Members` | `pulumi.StringArrayInput` | ✅ | The complete list of members for this role |

### Constructors

```go
func NewIAMMember(ctx *pulumi.Context, name string, args *IAMMemberArgs, opts ...pulumi.ResourceOption) (*IAMMember, error)
func NewIAMBinding(ctx *pulumi.Context, name string, args *IAMBindingArgs, opts ...pulumi.ResourceOption) (*IAMBinding, error)
```

## Examples

### Grant a Service Account Org-Level Access

```go
iam.NewIAMMember(ctx, "sa-org-admin", &iam.IAMMemberArgs{
    ParentID:   pulumi.String("123456789"),
    ParentType: "organization",
    Role:       pulumi.String("roles/resourcemanager.organizationAdmin"),
    Member:     pulumi.Sprintf("serviceAccount:%s", sa.Email),
})
```

### Grant Project-Level Access Using Outputs

```go
iam.NewIAMMember(ctx, "project-viewer", &iam.IAMMemberArgs{
    ParentID:   project.Project.ProjectId,  // pulumi.StringOutput
    ParentType: "project",
    Role:       pulumi.String("roles/viewer"),
    Member:     pulumi.String("group:developers@example.com"),
})
```

### Authoritative Binding on a Folder

```go
iam.NewIAMBinding(ctx, "folder-admins", &iam.IAMBindingArgs{
    ParentID:   folder.ID(),
    ParentType: "folder",
    Role:       pulumi.String("roles/resourcemanager.folderAdmin"),
    Members: pulumi.StringArray{
        pulumi.String("user:admin@example.com"),
        pulumi.Sprintf("serviceAccount:%s", bootstrapSA.Email),
    },
})
```

### Billing Account Binding

```go
iam.NewIAMMember(ctx, "billing-user", &iam.IAMMemberArgs{
    ParentID:   pulumi.String("XXXXXX-XXXXXX-XXXXXX"),
    ParentType: "billing",
    Role:       pulumi.String("roles/billing.user"),
    Member:     pulumi.Sprintf("serviceAccount:%s", sa.Email),
})
```

## Error Handling

Passing an unsupported `ParentType` returns an error:

```text
unsupported IAM parent type: "unknown" (expected organization, folder, project, serviceAccount, or billing)
```
