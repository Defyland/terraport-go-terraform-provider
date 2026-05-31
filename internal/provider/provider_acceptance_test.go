package provider

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"terraport": providerserver.NewProtocol6WithError(New("test")()),
}

func TestAccAPIProductDataSource(t *testing.T) {
	api := newFakeBankPortAPI(t)
	defer api.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProviderConfig(api, api.Token(), 5000, 4) + `
data "terraport_bankport_api_product" "bankport" {
  product_code = "bankport"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.terraport_bankport_api_product.bankport", "name", "BankPort Partner API"),
					resource.TestCheckResourceAttr("data.terraport_bankport_api_product.bankport", "category", "payments"),
					resource.TestCheckResourceAttr("data.terraport_bankport_api_product.bankport", "beta", "false"),
				),
			},
		},
	})
}

func TestAccPartnerAppLifecycleImportDrift(t *testing.T) {
	api := newFakeBankPortAPI(t)
	defer api.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkFakePartnerAppsDestroyed(api),
		Steps: []resource.TestStep{
			{
				Config: testProviderConfig(api, api.Token(), 5000, 4) + partnerAppConfig("ledger-connect", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("terraport_bankport_partner_app.main", "name", "ledger-connect"),
					resource.TestCheckResourceAttr("terraport_bankport_partner_app.main", "status", "active"),
					resource.TestCheckResourceAttrSet("terraport_bankport_partner_app.main", "client_id"),
					resource.TestCheckResourceAttrSet("terraport_bankport_partner_app.main", "client_secret"),
				),
			},
			{
				Config: testProviderConfig(api, api.Token(), 5000, 4) + partnerAppConfig("ledger-connect-updated", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("terraport_bankport_partner_app.main", "name", "ledger-connect-updated"),
					resource.TestCheckResourceAttr("terraport_bankport_partner_app.main", "client_secret_version", "2"),
				),
			},
			{
				ResourceName:            "terraport_bankport_partner_app.main",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"client_secret_version"},
			},
			{
				PreConfig: func() {
					api.MutateFirstPartnerAppName("drifted-outside-terraform")
				},
				Config:             testProviderConfig(api, api.Token(), 5000, 4) + partnerAppConfig("ledger-connect-updated", 2),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccWebhookEndpointLifecycleImport(t *testing.T) {
	api := newFakeBankPortAPI(t)
	defer api.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProviderConfig(api, api.Token(), 5000, 4) + webhookConfig(true, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("terraport_bankport_webhook_endpoint.main", "enabled", "true"),
					resource.TestCheckResourceAttrSet("terraport_bankport_webhook_endpoint.main", "signing_secret"),
				),
			},
			{
				Config: testProviderConfig(api, api.Token(), 5000, 4) + webhookConfig(false, 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("terraport_bankport_webhook_endpoint.main", "enabled", "false"),
					resource.TestCheckResourceAttr("terraport_bankport_webhook_endpoint.main", "signing_secret_version", "2"),
				),
			},
			{
				ResourceName:            "terraport_bankport_webhook_endpoint.main",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"signing_secret_version"},
			},
		},
	})
}

func TestAccRateLimitPolicyLifecycleImport(t *testing.T) {
	api := newFakeBankPortAPI(t)
	defer api.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProviderConfig(api, api.Token(), 5000, 4) + rateLimitPolicyConfig(600, 60, "enforce"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("terraport_bankport_rate_limit_policy.main", "requests_per_minute", "600"),
					resource.TestCheckResourceAttr("terraport_bankport_rate_limit_policy.main", "burst_limit", "60"),
				),
			},
			{
				Config: testProviderConfig(api, api.Token(), 5000, 4) + rateLimitPolicyConfig(1200, 120, "report"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("terraport_bankport_rate_limit_policy.main", "requests_per_minute", "1200"),
					resource.TestCheckResourceAttr("terraport_bankport_rate_limit_policy.main", "mode", "report"),
				),
			},
			{
				ResourceName:      "terraport_bankport_rate_limit_policy.main",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSandboxEnvironmentLifecycleImport(t *testing.T) {
	api := newFakeBankPortAPI(t)
	defer api.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProviderConfig(api, api.Token(), 5000, 4) + sandboxEnvironmentConfig("platform-sandbox", "sa-east-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("terraport_bankport_sandbox_environment.main", "region", "sa-east-1"),
					resource.TestCheckResourceAttrSet("terraport_bankport_sandbox_environment.main", "api_key_token"),
				),
			},
			{
				Config: testProviderConfig(api, api.Token(), 5000, 4) + sandboxEnvironmentConfig("platform-sandbox-updated", "us-east-1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("terraport_bankport_sandbox_environment.main", "name", "platform-sandbox-updated"),
					resource.TestCheckResourceAttr("terraport_bankport_sandbox_environment.main", "region", "us-east-1"),
				),
			},
			{
				ResourceName:      "terraport_bankport_sandbox_environment.main",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAuthFailure(t *testing.T) {
	api := newFakeBankPortAPI(t)
	defer api.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testProviderConfig(api, "wrong-token", 5000, 1) + partnerAppConfig("auth-failure", 1),
				ExpectError: regexp.MustCompile(`status=401`),
			},
		},
	})
}

func TestAccApplyTimeout(t *testing.T) {
	api := newFakeBankPortAPI(t)
	defer api.Close()
	api.Delay("POST /v1/partner-apps", 100*time.Millisecond)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testProviderConfig(api, api.Token(), 5, 1) + partnerAppConfig("timeout-app", 1),
				ExpectError: regexp.MustCompile(`timeout|deadline|Client.Timeout`),
			},
		},
	})
}

func TestAccRateLimitRetry(t *testing.T) {
	api := newFakeBankPortAPI(t)
	defer api.Close()
	api.Fail429("POST /v1/partner-apps", 2)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkFakePartnerAppsDestroyed(api),
		Steps: []resource.TestStep{
			{
				Config: testProviderConfig(api, api.Token(), 5000, 4) + partnerAppConfig("rate-limited-app", 1),
				Check: func(_ *terraform.State) error {
					if got := api.RequestCount("POST /v1/partner-apps"); got != 3 {
						return fmt.Errorf("expected 3 create attempts after rate-limit retries, got %d", got)
					}
					return nil
				},
			},
		},
	})
}

func TestAccPlanOnlyAvoidsRemoteResourceCalls(t *testing.T) {
	api := newFakeBankPortAPI(t)
	defer api.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             testProviderConfig(api, api.Token(), 5000, 4) + partnerAppConfig("plan-only-app", 1),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
	if got := api.TotalRequests(); got != 0 {
		t.Fatalf("plan-only resource config should not call the fake API, got %d requests", got)
	}
}

func TestAccHundredPartnerAppsApply(t *testing.T) {
	api := newFakeBankPortAPI(t)
	defer api.Close()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkFakePartnerAppsDestroyed(api),
		Steps: []resource.TestStep{
			{
				Config: testProviderConfig(api, api.Token(), 10000, 4) + hundredPartnerAppsConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					func(_ *terraform.State) error {
						if got := api.PartnerAppCount(); got != 100 {
							return fmt.Errorf("expected 100 partner apps in fake API, got %d", got)
						}
						return nil
					},
					resource.TestCheckResourceAttr("terraport_bankport_partner_app.app_099", "name", "Load App 099"),
				),
			},
		},
	})
}

func testProviderConfig(api *fakeBankPortAPI, token string, timeoutMS int, retryAttempts int) string {
	return fmt.Sprintf(`
provider "terraport" {
  endpoint             = %q
  token                = %q
  timeout_ms           = %d
  retry_max_attempts   = %d
  retry_min_delay_ms   = 1
}
`, api.URL(), token, timeoutMS, retryAttempts)
}

func partnerAppConfig(name string, secretVersion int) string {
	return fmt.Sprintf(`
resource "terraport_bankport_partner_app" "main" {
  name                  = %q
  product_code          = "bankport"
  redirect_uris         = ["https://partner.example.test/callback"]
  scopes                = ["accounts:read", "payments:write"]
  client_secret_version = %d
}
`, name, secretVersion)
}

func webhookConfig(enabled bool, secretVersion int) string {
	return fmt.Sprintf(`
resource "terraport_bankport_partner_app" "main" {
  name                  = "webhook-owner"
  product_code          = "bankport"
  redirect_uris         = ["https://partner.example.test/callback"]
  scopes                = ["webhooks:write"]
  client_secret_version = 1
}

resource "terraport_bankport_webhook_endpoint" "main" {
  partner_app_id         = terraport_bankport_partner_app.main.id
  url                    = "https://partner.example.test/webhooks/bankport"
  event_types            = ["partner_app.created", "payment.settled"]
  enabled                = %t
  signing_secret_version = %d
}
`, enabled, secretVersion)
}

func rateLimitPolicyConfig(rpm int, burst int, mode string) string {
	return fmt.Sprintf(`
resource "terraport_bankport_rate_limit_policy" "main" {
  product_code        = "bankport"
  subject_type        = "partner_app"
  subject_id          = "app-demo"
  requests_per_minute = %d
  burst_limit         = %d
  mode                = %q
}
`, rpm, burst, mode)
}

func sandboxEnvironmentConfig(name, region string) string {
	return fmt.Sprintf(`
resource "terraport_bankport_sandbox_environment" "main" {
  name     = %q
  products = ["bankport", "pixguard", "settleflow"]
  region   = %q
}
`, name, region)
}

func hundredPartnerAppsConfig() string {
	var b strings.Builder
	for i := range 100 {
		fmt.Fprintf(&b, `
resource "terraport_bankport_partner_app" "app_%03d" {
  name                  = "Load App %03d"
  product_code          = "bankport"
  redirect_uris         = ["https://partner-%03d.example.test/callback"]
  scopes                = ["accounts:read"]
  client_secret_version = 1
}
`, i, i, i)
	}
	return b.String()
}

func checkFakePartnerAppsDestroyed(api *fakeBankPortAPI) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if got := api.PartnerAppCount(); got != 0 {
			return fmt.Errorf("expected fake partner apps to be destroyed, got %d", got)
		}
		return nil
	}
}
