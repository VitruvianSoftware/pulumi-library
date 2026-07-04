# compute_instance

A Pulumi component for creating GCP Compute Instances from an existing instance template.

**Upstream Reference:** [terraform-google-modules/vm/google//modules/compute_instance](https://registry.terraform.io/modules/terraform-google-modules/vm/google)

## Installation

> This module is versioned and tagged independently. Git tags use the Go module-path format `go/pkg/compute_instance/vX.Y.Z`, so pin a version with `go get github.com/VitruvianSoftware/pulumi-library/go/pkg/compute_instance@vX.Y.Z`.

```bash
go get github.com/VitruvianSoftware/pulumi-library/go/pkg/compute_instance
```

```go
import "github.com/VitruvianSoftware/pulumi-library/go/pkg/compute_instance"
```

## Overview

Creates one or more compute instances from a template with:
- Configurable instance count
- Deletion protection
- Automatic sequential naming

## API Reference

### ComputeInstanceArgs

| Field | Type | Required | Description |
|:--|:--|:--|:--|
| `Project` | `pulumi.StringInput` | ✅ | GCP project ID |
| `Zone` | `string` | ✅ | GCP zone |
| `Hostname` | `string` | ✅ | Base hostname for instances |
| `InstanceTemplate` | `pulumi.StringInput` | ✅ | Self-link of the instance template |
| `NumInstances` | `int` | | Number of instances (default: 1) |
| `DeletionProtection` | `bool` | | Enable deletion protection |

### Outputs

| Field | Type | Description |
|:--|:--|:--|
| `Instances` | `[]*compute.InstanceFromTemplate` | The created instances |

## Usage

```go
instances, err := compute_instance.NewComputeInstance(ctx, "web-pool", &compute_instance.ComputeInstanceArgs{
    Project:          pulumi.String("my-project"),
    Zone:             "us-central1-a",
    Hostname:         "web-server",
    InstanceTemplate: tmpl.Template.SelfLinkUnique,
    NumInstances:     3,
})
```
