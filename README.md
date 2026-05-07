# writeonly Terraform Provider

This is a helper Terraform provider to work around the inability to detect changes to the contents of an ephemeral resource.

## Example Usage

```hcl
variable "some_secret" {
    type = string
    sensitive = true
    ephemeral = true
}

resource "writeonly_hash" "some_secret" {
  input_wo = var.some_secret
}

resource "time_static" "secret_update_timestamp" {
  triggers = {
    value = writeonly_hash.some_secret.output
  }
}

resource "aws_secretsmanager_secret_version" "some_secret" {
  secret_id                = "some_secret"
  secret_string_wo         = var.some_secret
  secret_string_wo_version = time_static.secret_update_timestamp.unix
}
```
