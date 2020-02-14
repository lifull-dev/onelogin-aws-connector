package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestInitCmdWithoutConfigFile(t *testing.T) {
	file := path.Join(os.TempDir(), "example.toml")
	defer os.Remove(file)

	resetInitFlags()
	endpoint = "api-server"
	clientToken = "client-token"
	clientSecret = "client-secret"
	subdomain = "subdomain"
	usernameOrEmail = "username-or-email"

	if err := initServiceConfig(file, "default"); err != nil {
		t.Errorf("%#v", err)
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		t.Errorf("%#v", err)
	}
	actual := string(data)
	expected := `[service]
  [service.default]
    endpoint = "api-server"
    client_token = "client-token"
    client_secret = "client-secret"
    subdomain = "subdomain"
    username_or_email = "username-or-email"

[app]
`
	if actual != expected {
		t.Errorf("'%v' is not equal '%v'", actual, expected)
	}
}

func TestInitCmdWithConfigFile(t *testing.T) {
	source, err := os.Open("fixtures/serviceconfig.toml")
	if err != nil {
		t.Errorf("%#v", err)
	}
	dist, err := ioutil.TempFile("", "onelogin-aws-connector")
	if err != nil {
		t.Errorf("%#v", err)
	}
	file := dist.Name()
	defer os.Remove(file)
	_, err = io.Copy(dist, source)
	if err != nil {
		t.Errorf("%#v", err)
	}

	resetInitFlags()
	endpoint = "new-api-server"
	clientToken = "new-client-token"
	clientSecret = "new-client-secret"
	subdomain = "new-subdomain"
	usernameOrEmail = "new-username-or-email"
	err = initServiceConfig(file, "default")
	if err != nil {
		t.Errorf("%#v", err)
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		t.Errorf("%#v", err)
	}
	actual := string(data)
	expected := `[service]
  [service.default]
    endpoint = "new-api-server"
    client_token = "new-client-token"
    client_secret = "new-client-secret"
    subdomain = "new-subdomain"
    username_or_email = "new-username-or-email"

[app]
`
	if actual != expected {
		t.Errorf("'%v' is not equal '%v'", actual, expected)
	}
}

func resetInitFlags() {
	endpoint = ""
	clientToken = ""
	clientSecret = ""
	subdomain = ""
	usernameOrEmail = ""
}
