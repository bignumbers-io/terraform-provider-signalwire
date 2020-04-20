package provider

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSignalwireSipEndpoint_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSignalwireSipEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSignalwireSipEndpointConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSignalwireSipEndpointExists("signalwire_sip_endpoint.test_endpoint"),
				),
			},
		},
	})
}

func testAccCheckSignalwireSipEndpointDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "signalwire_sip_endpoint" {
			continue
		}

		resp, err := client.Req("GET", os.Getenv("SIGNALWIRE_SPACE"), "endpoints/sip/", nil)
		if err != nil {
			return err
		}

		endpoints := resp["data"].([]interface{})

		if len(endpoints) > 0 {
			return errors.New("Endpoints still exists")
		}
		return nil
	}
	return nil
}

func testAccCheckSignalwireSipEndpointExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No key ID is set")
		}

		client := testAccProvider.Meta().(*Client)

		resp, err := client.Req("GET", os.Getenv("SIGNALWIRE_SPACE"), "endpoints/sip/"+rs.Primary.ID, nil)

		if err != nil {
			return err
		}

		if resp["id"] != rs.Primary.ID {
			return errors.New("SIP Endpoint does not match")
		}

		return nil
	}
}

var testAccSignalwireSipEndpointConfig = fmt.Sprintf(`
	resource "signalwire_sip_endpoint" "test_endpoint" {
        space = "%[1]s"
        username = "c3p0"
        password = "password"
        caller_id = "C-3P0"
        ciphers = [
			"AEAD_AES_256_GCM_8",
			"AES_256_CM_HMAC_SHA1_80",
			"AES_CM_128_HMAC_SHA1_80",
			"AES_256_CM_HMAC_SHA1_32",
			"AES_CM_128_HMAC_SHA1_32"
		]
  		codecs = [
			"OPUS",
			"G722",
			"PCMU",
			"PCMA",
			"VP8",
			"H264"
  		]
  		encryption = "optional"
	}
`, testSignalwireSpace)
