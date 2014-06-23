package libcontainer

import (
	"encoding/json"
	"os"
	"testing"
)

// Checks whether the expected capability is specified in the capabilities.
func contains(expected string, values []string) bool {
	for _, v := range values {
		if v == expected {
			return true
		}
	}
	return false
}

func TestContainerJsonFormat(t *testing.T) {
	f, err := os.Open("sample_configs/attach_to_bridge.json")
	if err != nil {
		t.Fatal("Unable to open container.json")
	}
	defer f.Close()

	var container *Container
	if err := json.NewDecoder(f).Decode(&container); err != nil {
		t.Fatalf("failed to decode container config: %s", err)
	}
	if container.Hostname != "koye" {
		t.Log("hostname is not set")
		t.Fail()
	}

	if !container.Tty {
		t.Log("tty should be set to true")
		t.Fail()
	}

	if !container.Namespaces["NEWNET"] {
		t.Log("namespaces should contain NEWNET")
		t.Fail()
	}

	if container.Namespaces["NEWUSER"] {
		t.Log("namespaces should not contain NEWUSER")
		t.Fail()
	}

	if contains("SYS_ADMIN", container.Capabilities) {
		t.Log("SYS_ADMIN should not be enabled in capabilities mask")
		t.Fail()
	}

	if !contains("MKNOD", container.Capabilities) {
		t.Log("MKNOD should be enabled in capabilities mask")
		t.Fail()
	}

	if !contains("SYS_CHROOT", container.Capabilities) {
		t.Log("capabilities mask should contain SYS_CHROOT")
		t.Fail()
	}

	for _, n := range container.Networks {
		if n.Type == "veth" {
			if n.Bridge != "docker0" {
				t.Logf("veth bridge should be docker0 but received %q", n.Bridge)
				t.Fail()
			}

			if n.Address != "172.17.0.101/16" {
				t.Logf("veth address should be 172.17.0.101/61 but received %q", n.Address)
				t.Fail()
			}

			if n.VethPrefix != "veth" {
				t.Logf("veth prefix should be veth but received %q", n.VethPrefix)
				t.Fail()
			}

			if n.Gateway != "172.17.42.1" {
				t.Logf("veth gateway should be 172.17.42.1 but received %q", n.Gateway)
				t.Fail()
			}

			if n.Mtu != 1500 {
				t.Logf("veth mtu should be 1500 but received %d", n.Mtu)
				t.Fail()
			}

			break
		}
	}
}
