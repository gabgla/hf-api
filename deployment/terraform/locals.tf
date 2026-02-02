locals {
  statefile_bucket = "tf-state-${var.namespace}-${var.environment}"
  statefile_key    = "hfapi.tfstate"

  resource_prefix = (
    var.environment == "prod" || var.environment == "live"
    ? "${var.namespace}-"
    : "${var.namespace}-${var.environment}-"
  )

  function_name = "${local.resource_prefix}hfapi-function"
}
