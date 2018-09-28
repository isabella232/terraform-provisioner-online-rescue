package rescue

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/src-d/terraform-provider-online-net/online"
	"github.com/src-d/terraform-provider-online-net/online/mock"
)

var onlineMock *mock.OnlineClientMock

func init() {
	onlineMock = new(mock.OnlineClientMock)
	provisionerOnlineClient = onlineMock
}

func testConfig(t *testing.T, c map[string]interface{}) *terraform.ResourceConfig {
	os.Setenv("ONLINE_TOKEN", "dummy")

	r, err := config.NewRawConfig(c)
	if err != nil {
		t.Fatalf("config error: %s", err)
	}

	return terraform.NewResourceConfig(r)
}

func TestProvisioner(t *testing.T) {
	if err := Provisioner().(*schema.Provisioner).InternalValidate(); err != nil {
		t.Fatalf("error: %s", err)
	}
}

func TestProvisioner_validate_good(t *testing.T) {
	c := testConfig(t, map[string]interface{}{
		"enabled":         "true",
		"server":          "123",
		"credentials_dir": "/tmp",
	})

	warn, errs := Provisioner().Validate(c)
	if len(warn) > 0 {
		t.Fatalf("Warnings: %v", warn)
	}
	if len(errs) > 0 {
		t.Fatalf("Errors: %v", errs)
	}
}

func TestProvisioner_validate_bad_server(t *testing.T) {
	c := testConfig(t, map[string]interface{}{
		"enabled":         "true",
		"server":          "terraform",
		"credentials_dir": "/tmp",
	})

	warn, errs := Provisioner().Validate(c)
	if len(warn) > 0 {
		t.Fatalf("Warnings: %v", warn)
	}
	if len(errs) != 1 {
		t.Fatalf("Should have one error")
	}
}

func TestProvisioner_validate_no_server(t *testing.T) {
	c := testConfig(t, map[string]interface{}{
		"enabled":         "true",
		"credentials_dir": "/tmp",
	})

	warn, errs := Provisioner().Validate(c)
	if len(warn) > 0 {
		t.Fatalf("Warnings: %v", warn)
	}
	if len(errs) != 1 {
		t.Fatalf("Should have one error")
	}
}

func TestProvisioner_validate_enabled_no_image(t *testing.T) {
	tmp, _ := ioutil.TempDir("/tmp", "unittestprovisioner")
	defer os.RemoveAll(tmp)

	c := testConfig(t, map[string]interface{}{
		"enabled":         "true",
		"server":          "123",
		"credentials_dir": tmp,
	})

	err := Provisioner().Apply(nil, nil, c)
	if err == nil {
		t.Fatalf("Expected error")
	}
}

func TestProvisioner_validate_enabled_no_path(t *testing.T) {
	tmp, _ := ioutil.TempDir("/tmp", "unittestprovisioner")
	defer os.RemoveAll(tmp)

	c := testConfig(t, map[string]interface{}{
		"enabled": "true",
		"server":  "123",
		"image":   "ubuntu-18.04_amd64",
	})

	err := Provisioner().Apply(nil, nil, c)
	if err == nil {
		t.Fatalf("Expected error")
	}
}

func TestProvisioner_apply_disabled(t *testing.T) {

	onlineMock.On("BootNormalMode", 123).Return(nil)
	c := testConfig(t, map[string]interface{}{
		"enabled": "false",
		"server":  "123",
	})

	err := Provisioner().Apply(nil, nil, c)
	if err != nil {
		t.Fatalf(err.Error())
	}

	onlineMock.AssertCalled(t, "BootNormalMode", 123)
}

func TestProvisioner_apply_enabled(t *testing.T) {
	// credentials to test
	ip := "127.0.0.1"
	username := "root"
	password := "gophers"
	protocol := "ssh"

	tmp, _ := ioutil.TempDir("/tmp", "unittestprovisioner")
	defer os.RemoveAll(tmp)

	onlineMock.On("BootRescueMode", 123, "ubuntu-18.04_amd64").Return(&online.RescueCredentials{
		IP:       ip,
		Login:    username,
		Password: password,
		Protocol: protocol,
	}, nil)

	c := testConfig(t, map[string]interface{}{
		"enabled":         "true",
		"server":          "123",
		"credentials_dir": tmp,
		"image":           "ubuntu-18.04_amd64",
	})

	err := Provisioner().Apply(nil, nil, c)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if readFile(tmp, "ip") != ip {
		t.Fatalf("ip file reads %s but needs %s", readFile(tmp, "ip"), ip)
	}

	if readFile(tmp, "username") != username {
		t.Fatalf("username file reads %s but needs %s", readFile(tmp, "username"), username)
	}

	if readFile(tmp, "password") != password {
		t.Fatalf("password file reads %s but needs %s", readFile(tmp, "password"), password)
	}

	if readFile(tmp, "protocol") != protocol {
		t.Fatalf("protocol file reads %s but needs %s", readFile(tmp, "protocol"), protocol)
	}
}

func readFile(dir, file string) string {
	bytes, _ := ioutil.ReadFile(path.Join(dir, file))

	return string(bytes)
}
