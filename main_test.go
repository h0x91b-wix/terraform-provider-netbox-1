package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"testing"
)

// pluginVersionExpected provides the expected plugin version that this release
// of the plugin is built with.
const pluginVersionExpected = 5

// TestPluginVersion tests the plugin API version of the vendored Terraform
// libraries.
func TestPluginVersion(t *testing.T) {
	if plugin.Handshake.ProtocolVersion != pluginVersionExpected {
		t.Fatalf("Expected vendored plugin version to be %d, got %d", pluginVersionExpected, plugin.Handshake.ProtocolVersion)
	}
}
