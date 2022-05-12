<!-- markdownlint-disable MD024 -->

# Use Terraform AzAPI Provider to test API and provide examples

## Overview

This document introduced a Terraform AzAPI Provider based practice to check API runtime correctness and provide a better
example format.

The [AzAPI provider](https://docs.microsoft.com/en-us/azure/developer/terraform/overview-azapi-provider) is a thin layer
on top of the Azure ARM REST APIs, and it enables you to manage any Azure resource type using any API version. 

The idea is to use AzAPI provider to manage the Azure resource and see whether there are Terraform development blockers.

And examples in Terraform configuration can describe complex dependency references and deployment orders, because
Terraform will figure out the dependency graph of the resources and deploy them in a correct oder.

## Prerequisites

- Install [Terraform](https://www.terraform.io/)
- Install [Terraform VSCode Extension](https://marketplace.visualstudio.com/items?itemName=HashiCorp.terraform)
- Install [AzApi VSCode Extension](https://marketplace.visualstudio.com/items?itemName=azapi-vscode.azapi)
- Learn Terraform basics:
  - Terraform's primary function is to create, modify, and destroy infrastructure resources to match the desired state
    described in a [Terraform configuration](https://www.terraform.io/language).
  - The `terraform init` command is used to initialize a working directory containing Terraform configuration files.
  - The `terraform plan` command creates an execution plan, which lets you preview the changes that Terraform plans 
  to make to your infrastructure.
  - The `terraform apply` command performs a plan just like terraform plan does, but then actually carries out the 
    planned changes to each resource using the relevant infrastructure provider's API.
  - Step-by-step tutorials: [Get Started - Azure](https://learn.hashicorp.com/collections/terraform/azure-get-started)

## How does it work?

1. Compose a Terraform configuration, see the below example, which uses `azapi` provider to manage the testing resources, and `azurerm`
provider to manage other dependencies. The resources supported in `azurerm` provider can be found [here](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs).
   Each resource has a valid example configuration can be used directly.
2. Run `terraform init` commands to download the providers, only needs be executed in a new working directory.
3. Run `terraform apply` commands to deploy them. 
4. Run `terraform plan` command to check whether the resources actually deployed are matched with the configurations, 
   if not, there must be some requirements are not met, for example, the property's value isn't consistent between 
   PUT and GET requests, or incorrect casing.

In next sections, we'll introduce tools to help you compose and test the Terraform configuration. 
```hcl
##### Provider declaration: tells Terraform which provider will be used
terraform {
  required_providers {
    azapi = {
      source = "azure/azapi"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "azapi" {
}

#### azurerm_resource_group manages a resource group
resource "azurerm_resource_group" "test" {
  name     = "myResourceGroup"
  location = "westus"
}

#### azapi_resource manages the testing resource
resource "azapi_resource" "account" {
  type = "Microsoft.LabServices/labaccounts@2018-10-15"
  name = "myAccount"
  // refers to its parent resource's id, terraform will know this resource depends on azurerm_resource_group.test
  parent_id = azurerm_resource_group.test.id
  
  // refers to resource group's location
  location = azurerm_resource_group.test.location
  
  // request payload, jsonencode converts a hcl object to json
  body = jsonencode({
    properties = {
      enabledRegionSelection = false
    }
  })
  
  /*
  # you can pass JSON as well.
  body = <<BODY
  {
    "properties" : {
      "enabledRegionSelection": false
    }
  }
  BODY  
  */
}
```

## How to compose the configurations?

### Automatically generate from swagger example

This section introduces a tool to help users automatically generate a full Terraform configuration from a swagger example
file.

1. Download [azurerm-restapi-testing-tool](https://github.com/ms-henglu/azurerm-restapi-testing-tool/releases)  
2. Run azure cli commands to login: `az login`. Terraform supports a number of different methods for [authenticating to Azure](https://registry.terraform.io/providers/Azure/azapi/latest/docs#authenticating-to-azure). 
3. Generate terraform files and Test
   1. Generate testing files by running `azurerm-restapi-testing-tool generate -path path_to_swagger_example`.
       Here's an example:
       ```bash
       azurerm-restapi-testing-tool generate -path ./CreateComputeInstanceMinimal.json
       ```

       Then `dependency.tf` and `testing.tf` will be generated.
   2. Run API tests by running `azurerm-restapi-testing-tool test`. This command will set up dependencies and test the ARM resource API.
   3. There's an `auto` command, it can generate testing files, then run the tests and remove all resources if test is passed. Example:

      ```bash
      azurerm-restapi-testing-tool auto -path ./CreateComputeInstanceMinimal.json
      ```

Here are some [examples](https://github.com/ms-henglu/azurerm-restapi-testing-tool/tree/main/examples).

### Manually compose the configuration

[AzApi VSCode Extension](https://marketplace.visualstudio.com/items?itemName=azapi-vscode.azapi) provides a rich 
authoring experience, with the following benefits:
- Intellisense
- Code auto-completion
- Hints
- Syntax validation
- Quick info

> [!IMPORTANT]
> Note: AzApi Extension uses an embedded schema which is generated from `azure-rest-api-specs` to provide the language
> features. Because of its release cadence, the schema may not be the latest, and it may have one month latency.

The extension suggests to compose `body` with `jsonencode` function, which is more concise and easier to read.

## Summary

This practice can detect Terraform development blockers during API development phase, it also provides rich examples
which can serve as valuable E2E testcases, and the examples can be used by customers directly if they want to have Terraform
support on the preview features on day 0.
