# Terraform Couchbase Capella Provider 🌟

![Couchbase Capella](https://img.shields.io/badge/Couchbase-Capella-blue?style=flat&logo=couchbase)

Welcome to the **Terraform Couchbase Capella Provider** repository! This provider allows you to deploy, update, and manage your Couchbase Capella infrastructure as code through HashiCorp Terraform. 

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Getting Started](#getting-started)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Resources](#resources)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)
- [Support](#support)

## Introduction

Couchbase Capella is a fully managed database-as-a-service that offers a flexible and scalable solution for modern applications. With the Terraform provider, you can automate the management of your Couchbase Capella resources, ensuring consistency and reducing manual errors.

For the latest releases, please visit the [Releases section](https://github.com/Lazzzey/terraform-provider-couchbase-capella/releases). You can download and execute the necessary files from there.

## Features

- **Infrastructure as Code**: Define your Couchbase Capella resources in a simple, declarative way.
- **Version Control**: Keep track of changes in your infrastructure.
- **Automation**: Easily deploy and manage your Couchbase infrastructure without manual intervention.
- **Scalability**: Scale your Couchbase Capella resources up or down based on your application needs.

## Getting Started

To get started with the Terraform Couchbase Capella Provider, follow these steps:

1. **Install Terraform**: Ensure you have Terraform installed on your machine. You can download it from the [Terraform website](https://www.terraform.io/downloads.html).
   
2. **Clone this Repository**: Clone the repository to your local machine using the following command:

   ```bash
   git clone https://github.com/Lazzzey/terraform-provider-couchbase-capella.git
   ```

3. **Navigate to the Directory**: Change to the cloned directory:

   ```bash
   cd terraform-provider-couchbase-capella
   ```

## Installation

To install the Couchbase Capella provider, you can use the following command in your Terraform configuration file:

```hcl
terraform {
  required_providers {
    couchbase = {
      source  = "Lazzzey/couchbase-capella"
      version = "1.0.0" # Replace with the latest version
    }
  }
}
```

For the latest releases, check the [Releases section](https://github.com/Lazzzey/terraform-provider-couchbase-capella/releases) to download and execute the appropriate files.

## Usage

Here's a simple example of how to use the Couchbase Capella provider in your Terraform configuration.

```hcl
provider "couchbase" {
  username = var.couchbase_username
  password = var.couchbase_password
}

resource "couchbase_cluster" "example" {
  name     = "example-cluster"
  region   = "us-west"
  plan     = "development"
  version  = "7.0"
}
```

## Configuration

### Provider Configuration

You need to configure the provider with your Couchbase Capella credentials. Create a `variables.tf` file to define the required variables:

```hcl
variable "couchbase_username" {
  description = "Couchbase Capella username"
  type        = string
}

variable "couchbase_password" {
  description = "Couchbase Capella password"
  type        = string
}
```

### Resource Configuration

You can define multiple resources in your configuration. Here are some examples:

#### Cluster Resource

```hcl
resource "couchbase_cluster" "my_cluster" {
  name     = "my-cluster"
  region   = "us-east"
  plan     = "production"
  version  = "7.0"
}
```

#### Bucket Resource

```hcl
resource "couchbase_bucket" "my_bucket" {
  name     = "my-bucket"
  cluster  = couchbase_cluster.my_cluster.name
  ram_quota = 256
}
```

## Resources

### Available Resources

- `couchbase_cluster`: Manage Couchbase clusters.
- `couchbase_bucket`: Manage Couchbase buckets.
- `couchbase_scope`: Manage scopes within a bucket.
- `couchbase_collection`: Manage collections within a scope.

### Data Sources

- `couchbase_clusters`: Retrieve information about existing Couchbase clusters.
- `couchbase_buckets`: Retrieve information about existing Couchbase buckets.

## Examples

For more detailed examples, check the `examples` directory in this repository. Here are a few examples to get you started:

### Basic Example

A basic example of setting up a Couchbase cluster and bucket:

```hcl
provider "couchbase" {
  username = var.couchbase_username
  password = var.couchbase_password
}

resource "couchbase_cluster" "example" {
  name     = "example-cluster"
  region   = "us-west"
  plan     = "development"
  version  = "7.0"
}

resource "couchbase_bucket" "example_bucket" {
  name      = "example-bucket"
  cluster   = couchbase_cluster.example.name
  ram_quota = 256
}
```

### Advanced Example

An advanced example with scopes and collections:

```hcl
provider "couchbase" {
  username = var.couchbase_username
  password = var.couchbase_password
}

resource "couchbase_cluster" "advanced_example" {
  name     = "advanced-cluster"
  region   = "us-east"
  plan     = "production"
  version  = "7.0"
}

resource "couchbase_bucket" "advanced_bucket" {
  name      = "advanced-bucket"
  cluster   = couchbase_cluster.advanced_example.name
  ram_quota = 512
}

resource "couchbase_scope" "my_scope" {
  bucket = couchbase_bucket.advanced_bucket.name
  name   = "my_scope"
}

resource "couchbase_collection" "my_collection" {
  scope   = couchbase_scope.my_scope.name
  bucket  = couchbase_bucket.advanced_bucket.name
  name    = "my_collection"
}
```

## Contributing

We welcome contributions to improve the Terraform Couchbase Capella Provider. Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them.
4. Push your branch to your forked repository.
5. Open a pull request to the main repository.

Please ensure that your code adheres to the existing style and includes tests where applicable.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Support

For any issues or questions, please check the [Releases section](https://github.com/Lazzzey/terraform-provider-couchbase-capella/releases) for updates and downloads. If you need further assistance, feel free to open an issue in this repository.

---

Thank you for using the Terraform Couchbase Capella Provider! We hope it helps you manage your Couchbase infrastructure efficiently.