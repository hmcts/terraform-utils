resource "azurerm_key_vault_secret" "test_secret_for_tests" {
  name  = "test-secret-for-tests"
}

resource "azurerm_resource_group" "test" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_dns_zone" "test" {
  name                = "mydomain.com"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_key_vault_secret" "s2s_secret_for_tests" {
  name  = "s2s-secret-for-tests"
}

module "bulk-scan-orchestrator" {
  source = "git@github.com:hmcts/cnp-module-webapp?ref=master"
}
