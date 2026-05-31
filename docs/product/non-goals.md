# Non-Goals

- No real production BankPort API is implemented in this repository.
- No CDKTF constructs are generated or used as the core provider implementation.
- No Terraform Registry publishing flow is included.
- No local daemon, database, or control plane is run by the provider.
- No attempt is made to make generated secrets disappear from Terraform state; the provider marks them sensitive and documents state controls.
