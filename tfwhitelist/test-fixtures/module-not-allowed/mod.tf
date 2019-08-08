module "bulk-scan-orchestrator" {
  source = "git@github.com:hmcts/cnp-module-webapp?ref=master"
}

module "bulk-scan-custom" {
  source = "git@github.com:hmcts/cnp-module-webapp?ref=my-branch"
}
