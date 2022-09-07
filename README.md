# terraform-provider-gcp-ext

## google_compute_resource_bind
insert resource policy in an exist google compute instance
```azure
terraform {
  required_providers {
    gcp-ext = {
      source = "frankzye/gcp-ext"
      version = "1.0.3"
    }
  }
}

provider "gcp-ext" {
  # Configuration options
}

resource "google_compute_resource_bind" "example" {
    provider = gcp-ext
    project  = "xxx"
    zone = "europe-west1"
    instance = "vm-name"
    policy = "projects/xxx/regions/europe-west1/resourcePolicies/policy-example"
}
```