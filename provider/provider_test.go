package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SIGNALWIRE_AUTH_TOKEN"); v == "" {
		t.Fatal("SIGNALWIRE_AUTH_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("SIGNALWIRE_PROJECT_ID"); v == "" {
		t.Fatal("SIGNALWIRE_PROJECT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("SIGNALWIRE_PROJECT_ID"); v == "" {
		t.Fatal("SIGNALWIRE_PROJECT_ID must be set for acceptance tests")
	}
	if v := os.Getenv("SIGNALWIRE_SPACE"); v == "" {
		t.Fatal("SIGNALWIRE_SPACE must be set for acceptance tests")
	}
}

var testSignalwireSpace = os.Getenv("SIGNALWIRE_SPACE")

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"signalwire": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderImpl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}
