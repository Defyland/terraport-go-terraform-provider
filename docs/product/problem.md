# Product Problem

Partner onboarding for BankPort-style APIs has several control-plane resources that are usually created outside version control: applications, webhook endpoints, rate-limit exceptions, and sandbox environments. Those manual changes are hard to review, hard to reproduce, and easy to forget during incident response.

Terraport makes those resources declarative. The product value is not that Terraform can call HTTP. The value is that onboarding changes become reviewable, importable, drift-detectable, and repeatable across partner environments.

## Core Workflow

1. A platform engineer opens a Terraform change for a partner app and sandbox.
2. A reviewer checks scopes, webhook destination, and rate-limit policy.
3. Terraform apply creates or updates the remote BankPort resources.
4. Future plans refresh remote state and show drift when a platform operator changes the resource outside Terraform.

## Business Value

- Reduces manual setup time for partner integrations.
- Makes sensitive generated values visible only through Terraform-sensitive outputs and state policy.
- Gives support engineers import and drift runbooks instead of ad hoc portal recovery.
