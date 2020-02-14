package config

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestLoadNoFile(t *testing.T) {
	file := path.Join(os.TempDir(), "notexists.toml")
	defer os.Remove(file)

	c, err := Load(file)
	if err != nil {
		t.Errorf("%#v", err)
	}
	if len(c.Service) > 0 {
		t.Errorf("ServiceConfigs is not empty: %#v", c.Service)
	}
	if len(c.App) > 0 {
		t.Errorf("AppConfigs is not empty: %#v", c.App)
	}
	err = c.Save()
	if err != nil {
		t.Errorf("%#v", err)
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		t.Errorf("%#v", err)
	}
	actual := string(data)
	expected := `[service]

[app]
`
	if actual != expected {
		t.Errorf("%s is not equal %s", actual, expected)
	}
}

func TestLoadServiceOnlyFile(t *testing.T) {
	source, err := os.Open("../fixtures/serviceconfig.toml")
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

	c, err := Load(file)
	if err != nil {
		t.Errorf("%#v", err)
	}
	if len(c.Service) != 1 {
		t.Errorf("ServiceConfigs is not empty: %#v", c.Service)
	}
	if len(c.App) > 0 {
		t.Errorf("AppConfigs is not empty: %#v", c.App)
	}
	err = c.Save()
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
    endpoint = "api-server"
    client_token = "client-token"
    client_secret = "client-secret"
    subdomain = "subdomain"
    username_or_email = "username-or-email"

[app]
`
	if actual != expected {
		t.Errorf("%v is not equal %v", actual, expected)
	}

	service := &ServiceConfig{
		Endpoint:        "new-api-server",
		ClientToken:     "new-client-token",
		ClientSecret:    "new-client-secret",
		Subdomain:       "new-subdomain",
		UsernameOrEmail: "new-username-or-email",
	}
	c.Service["default"] = service
	err = c.Save()
	if err != nil {
		t.Errorf("%#v", err)
	}
	data, err = ioutil.ReadFile(file)
	if err != nil {
		t.Errorf("%#v", err)
	}
	actual = string(data)
	expected = `[service]
  [service.default]
    endpoint = "new-api-server"
    client_token = "new-client-token"
    client_secret = "new-client-secret"
    subdomain = "new-subdomain"
    username_or_email = "new-username-or-email"

[app]
`
	if actual != expected {
		t.Errorf("%v is not equal %v", actual, expected)
	}
}

func TestLoadNormalFile(t *testing.T) {

	source, err := os.Open("../fixtures/fullfilled.toml")
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

	c, err := Load(file)
	if err != nil {
		t.Errorf("%#v", err)
	}
	if len(c.Service) != 1 {
		t.Errorf("ServiceConfigs is not empty: %#v", c.Service)
	}
	if len(c.App) != 2 {
		t.Errorf("AppConfigs is not empty: %#v", c.App)
	}
	err = c.Save()
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
    duration_seconds = 0
  [app.other]
    app_id = "other-app-id"
    role_arn = "other-role-arn"
    principal_arn = "other-provider-arn"
    duration_seconds = 0
`
	if actual != expected {
		t.Errorf("%s is not equal %s", actual, expected)
	}

	service := &ServiceConfig{
		Endpoint:        "new-api-server",
		ClientToken:     "new-client-token",
		ClientSecret:    "new-client-secret",
		Subdomain:       "new-subdomain",
		UsernameOrEmail: "new-username-or-email",
	}
	app := &AppConfig{
		AppID:        "new-app-id",
		RoleArn:      "new-role-arn",
		PrincipalArn: "new-principal-arn",
	}
	c.Service["default"] = service
	c.App["other"] = app
	err = c.Save()
	if err != nil {
		t.Errorf("%#v", err)
	}
	data, err = ioutil.ReadFile(file)
	if err != nil {
		t.Errorf("%#v", err)
	}
	actual = string(data)
	expected = `[service]
  [service.default]
    endpoint = "new-api-server"
    client_token = "new-client-token"
    client_secret = "new-client-secret"
    subdomain = "new-subdomain"
    username_or_email = "new-username-or-email"

[app]
  [app.default]
    app_id = "app-id"
    role_arn = "role-arn"
    principal_arn = "provider-arn"
    duration_seconds = 0
  [app.other]
    app_id = "new-app-id"
    role_arn = "new-role-arn"
    principal_arn = "new-principal-arn"
    duration_seconds = 0
`
	if actual != expected {
		t.Errorf("%v is not equal %v", actual, expected)
	}
}
