# Abuse Cases

## State Exfiltration

An attacker with state backend access reads sensitive generated values. Control: encrypted backend, least privilege, audit logging, and the state leak runbook.

## Token Misuse

A provider token with broad scopes is reused in CI. Control: short-lived tokens, separate CI identity, and scope review.

## Rate-Limit Bypass

A partner requests `report` mode or excessive burst limits. Control: review Terraform changes and keep product defaults in platform policy.

## Manual Remote Change

An operator changes a resource through the portal. Control: drift detection during plan and the drift runbook.

## Secret Printed by Output

A module outputs a secret without `sensitive = true`. Control: module review and Terraform policy checks outside this provider.
