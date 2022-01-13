package resource_test

import (
	"log"
	"testing"

	"github.com/ms-henglu/azurerm-rest-api-testing-tool/resource"
	"github.com/ms-henglu/azurerm-rest-api-testing-tool/types"
)

func Test_NewResourceFromExample(t *testing.T) {
	r, err := resource.NewResourceFromExample("./testdata/example.json")
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("expect valid resource, but got nil")
	}
	if r.ApiVersion != "2020-06-01" {
		t.Fatalf("expect api-version 2020-06-01, but got %s", r.ApiVersion)
	}

	expectExampleId := "/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.MachineLearningServices/workspaces/workspaces123/computes/compute123"
	if r.ExampleId != expectExampleId {
		t.Fatalf("expect exampleId %s, but got %s", expectExampleId, r.ExampleId)
	}

	if len(r.PropertyDependencyMappings) != 7 {
		t.Fatalf("expect PropertyDependencyMappings length %v, but got %v", 7, len(r.PropertyDependencyMappings))
	}
}

func Test_GetDependencyHcl(t *testing.T) {
	r, err := resource.NewResourceFromExample("./testdata/example.json")
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("expect valid resource, but got nil")
	}

	deps := make([]types.Dependency, 0)
	deps = append(deps, types.Dependency{
		Pattern:              "/subscriptions/resourceGroups/providers/Microsoft.MachineLearningServices/workspaces",
		ExampleConfiguration: "\r\nprovider \"azurerm\" {\r\n  features {}\r\n}\r\n\r\ndata \"azurerm_client_config\" \"current\" {}\r\n\r\nresource \"azurerm_resource_group\" \"example\" {\r\n  name     = \"example-resources\"\r\n  location = \"West Europe\"\r\n}\r\n\r\nresource \"azurerm_application_insights\" \"example\" {\r\n  name                = \"workspace-example-ai\"\r\n  location            = azurerm_resource_group.example.location\r\n  resource_group_name = azurerm_resource_group.example.name\r\n  application_type    = \"web\"\r\n}\r\n\r\nresource \"azurerm_key_vault\" \"example\" {\r\n  name                = \"workspaceexamplekeyvault\"\r\n  location            = azurerm_resource_group.example.location\r\n  resource_group_name = azurerm_resource_group.example.name\r\n  tenant_id           = data.azurerm_client_config.current.tenant_id\r\n  sku_name            = \"premium\"\r\n}\r\n\r\nresource \"azurerm_storage_account\" \"example\" {\r\n  name                     = \"workspacestorageaccount\"\r\n  location                 = azurerm_resource_group.example.location\r\n  resource_group_name      = azurerm_resource_group.example.name\r\n  account_tier             = \"Standard\"\r\n  account_replication_type = \"GRS\"\r\n}\r\n\r\nresource \"azurerm_machine_learning_workspace\" \"example\" {\r\n  name                    = \"example-workspace\"\r\n  location                = azurerm_resource_group.example.location\r\n  resource_group_name     = azurerm_resource_group.example.name\r\n  application_insights_id = azurerm_application_insights.example.id\r\n  key_vault_id            = azurerm_key_vault.example.id\r\n  storage_account_id      = azurerm_storage_account.example.id\r\n\r\n  identity {\r\n    type = \"SystemAssigned\"\r\n  }\r\n}\r\n",
		ResourceType:         "azurerm_machine_learning_workspace",
		ReferredProperty:     "id",
	})
	output := r.GetDependencyHcl(deps)
	log.Printf("Test_GetDependencyHcl output: %s", output)
	if len(output) == 0 {
		t.Fatal("expect valid config, but got empty string")
	}
}

func Test_GetParentReference(t *testing.T) {
	r, err := resource.NewResourceFromExample("./testdata/example.json")
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("expect valid resource, but got nil")
	}

	deps := make([]types.Dependency, 0)
	deps = append(deps, types.Dependency{
		Pattern:              "/subscriptions/resourceGroups/providers/Microsoft.MachineLearningServices/workspaces",
		ExampleConfiguration: "\r\nprovider \"azurerm\" {\r\n  features {}\r\n}\r\n\r\ndata \"azurerm_client_config\" \"current\" {}\r\n\r\nresource \"azurerm_resource_group\" \"example\" {\r\n  name     = \"example-resources\"\r\n  location = \"West Europe\"\r\n}\r\n\r\nresource \"azurerm_application_insights\" \"example\" {\r\n  name                = \"workspace-example-ai\"\r\n  location            = azurerm_resource_group.example.location\r\n  resource_group_name = azurerm_resource_group.example.name\r\n  application_type    = \"web\"\r\n}\r\n\r\nresource \"azurerm_key_vault\" \"example\" {\r\n  name                = \"workspaceexamplekeyvault\"\r\n  location            = azurerm_resource_group.example.location\r\n  resource_group_name = azurerm_resource_group.example.name\r\n  tenant_id           = data.azurerm_client_config.current.tenant_id\r\n  sku_name            = \"premium\"\r\n}\r\n\r\nresource \"azurerm_storage_account\" \"example\" {\r\n  name                     = \"workspacestorageaccount\"\r\n  location                 = azurerm_resource_group.example.location\r\n  resource_group_name      = azurerm_resource_group.example.name\r\n  account_tier             = \"Standard\"\r\n  account_replication_type = \"GRS\"\r\n}\r\n\r\nresource \"azurerm_machine_learning_workspace\" \"example\" {\r\n  name                    = \"example-workspace\"\r\n  location                = azurerm_resource_group.example.location\r\n  resource_group_name     = azurerm_resource_group.example.name\r\n  application_insights_id = azurerm_application_insights.example.id\r\n  key_vault_id            = azurerm_key_vault.example.id\r\n  storage_account_id      = azurerm_storage_account.example.id\r\n\r\n  identity {\r\n    type = \"SystemAssigned\"\r\n  }\r\n}\r\n",
		ResourceType:         "azurerm_machine_learning_workspace",
		ReferredProperty:     "id",
	})
	depHcl := r.GetDependencyHcl(deps)
	output := r.GetParentReference(depHcl)
	expect := "azurerm_machine_learning_workspace.test.id"
	if output != expect {
		t.Fatalf("expect %s but got %s", expect, output)
	}
}

func Test_GetHcl(t *testing.T) {
	r, err := resource.NewResourceFromExample("./testdata/example.json")
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("expect valid resource, but got nil")
	}

	deps := make([]types.Dependency, 0)
	deps = append(deps, types.Dependency{
		Pattern:              "/subscriptions/resourceGroups/providers/Microsoft.MachineLearningServices/workspaces",
		ExampleConfiguration: "\r\nprovider \"azurerm\" {\r\n  features {}\r\n}\r\n\r\ndata \"azurerm_client_config\" \"current\" {}\r\n\r\nresource \"azurerm_resource_group\" \"example\" {\r\n  name     = \"example-resources\"\r\n  location = \"West Europe\"\r\n}\r\n\r\nresource \"azurerm_application_insights\" \"example\" {\r\n  name                = \"workspace-example-ai\"\r\n  location            = azurerm_resource_group.example.location\r\n  resource_group_name = azurerm_resource_group.example.name\r\n  application_type    = \"web\"\r\n}\r\n\r\nresource \"azurerm_key_vault\" \"example\" {\r\n  name                = \"workspaceexamplekeyvault\"\r\n  location            = azurerm_resource_group.example.location\r\n  resource_group_name = azurerm_resource_group.example.name\r\n  tenant_id           = data.azurerm_client_config.current.tenant_id\r\n  sku_name            = \"premium\"\r\n}\r\n\r\nresource \"azurerm_storage_account\" \"example\" {\r\n  name                     = \"workspacestorageaccount\"\r\n  location                 = azurerm_resource_group.example.location\r\n  resource_group_name      = azurerm_resource_group.example.name\r\n  account_tier             = \"Standard\"\r\n  account_replication_type = \"GRS\"\r\n}\r\n\r\nresource \"azurerm_machine_learning_workspace\" \"example\" {\r\n  name                    = \"example-workspace\"\r\n  location                = azurerm_resource_group.example.location\r\n  resource_group_name     = azurerm_resource_group.example.name\r\n  application_insights_id = azurerm_application_insights.example.id\r\n  key_vault_id            = azurerm_key_vault.example.id\r\n  storage_account_id      = azurerm_storage_account.example.id\r\n\r\n  identity {\r\n    type = \"SystemAssigned\"\r\n  }\r\n}\r\n",
		ResourceType:         "azurerm_machine_learning_workspace",
		ReferredProperty:     "id",
	})
	depHcl := r.GetDependencyHcl(deps)
	output := r.GetHcl(depHcl)
	log.Printf("Test_GetHcl output: %s", output)
	if len(output) == 0 {
		t.Fatal("expect valid config, but got empty string")
	}
}