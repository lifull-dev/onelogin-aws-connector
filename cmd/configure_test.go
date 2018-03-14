package cmd

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestConfigureCmdWithoutInit(t *testing.T) {
	file := path.Join(os.TempDir(), "example.toml")
	defer os.Remove(file)

	resetConfigureFlags()
	appID = "app-id"
	roleArn = "role-arn"
	principalArn = "provider-arn"
	err := initAppConfig(file, "default")
	if err == nil {
		t.Error("It need to return none initialized error.")
	}
	errorMessage := "There is no initialized service. Please run `onelogin-aws-connector init`"
	if err.Error() != errorMessage {
		t.Errorf("%s is not equal to %s", err.Error(), errorMessage)
	}
}

func TestConfigureCmdWithService(t *testing.T) {
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
	io.Copy(dist, source)

	resetConfigureFlags()
	appID = "app-id"
	roleArn = "role-arn"
	principalArn = "provider-arn"
	if err := initAppConfig(file, "default"); err != nil {
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
  [app.default]
    app_id = "app-id"
    role_arn = "role-arn"
    principal_arn = "provider-arn"
`
	if actual != expected {
		t.Errorf("'%v' is not equal '%v'", actual, expected)
	}
}

func TestConfigureCmdWithDefault(t *testing.T) {
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
	io.Copy(dist, source)

	resetConfigureFlags()
	appID = "new-app-id"
	roleArn = "new-role-arn"
	principalArn = "new-provider-arn"
	if err := initAppConfig(file, "default"); err != nil {
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
  [app.default]
    app_id = "new-app-id"
    role_arn = "new-role-arn"
    principal_arn = "new-provider-arn"
`
	if actual != expected {
		t.Errorf("'%v' is not equal '%v'", actual, expected)
	}
}

func resetConfigureFlags() {
	appID = ""
	roleArn = ""
	principalArn = ""
}
