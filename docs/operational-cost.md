# Operational Cost

## Infrastructure Cost

The provider is a local binary. There is no service, database, or broker to host. CI runs Go tests, Terraform Plugin Testing, benchmarks, OpenAPI lint, security scan, and Docker build validation.

## Debugging Cost

Most incidents happen at the Terraform/API boundary: auth failure, rate limits, timeouts, drift, or import mismatch. Runbooks reduce support time by naming exact checks.

## Deploy Cost

Publishing is deferred. A real release process would add signed binaries, changelog generation, provider docs, and registry publishing.

## Backup and Retention Cost

Terraform state backend retention matters because state contains generated secrets. Retention should be short enough to limit blast radius but long enough for recovery.

## Monitoring Burden

The provider does not emit long-lived metrics. Teams should monitor CI failures, Terraform apply failures, and BankPort API rate-limit/auth metrics.

## Vendor Lock-in

The provider depends on Terraform Plugin Framework and Terraform state semantics. The rejected alternative, CDKTF, would add another runtime and generated abstraction without improving provider lifecycle evidence.

## Simpler Alternatives

Manual portal changes are cheaper initially but lose review, import, drift detection, and repeatability. A custom CLI would be easier to write but would not integrate with Terraform state and plans.
