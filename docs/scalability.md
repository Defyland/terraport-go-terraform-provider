# Scalability

## Hot Path

Terraform refresh and apply are the hot paths. Each managed resource typically requires one remote read during refresh and one remote write during apply.

## Read-Heavy Paths

- `Read` for resources during refresh.
- `Read` for `bankport_api_product` data sources.

## Write-Heavy Paths

- Creating many partner apps during onboarding.
- Rotating webhook or app secrets across environments.

## Fastest-Growing Data

Terraform state grows with resource count and sensitive generated values. The provider itself stores no persistent data.

## Queue Buildup

No provider-owned queue exists. The analogous bottleneck is Terraform parallel operations waiting on remote API rate limits.

## Hot Partitions

In a real BankPort API, subject IDs such as a large partner organization or product code could become hot keys for rate-limit and webhook operations.

## Horizontal Scaling

Terraform can run independent workspaces separately, but a single apply is bounded by Terraform Core, provider process concurrency, and remote API limits.

## Sharding or Partitioning

Large portfolios should split partner environments by workspace or product boundary. Remote API batch endpoints are preferable to provider-local caching.

## Async Candidates

Sandbox provisioning could become asynchronous if real environment setup takes minutes. The current fake API keeps it synchronous for test speed.

## Must Not Be Eventual

Secret rotation and rate-limit policy updates should be confirmed synchronously because Terraform state must reflect the secret or control value created by the API.
