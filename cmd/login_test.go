package cmd

import "testing"

func TestLoginCmdFetchConfigConfigVars(t *testing.T) {
	_, app, err := fetchConfig("fixtures/fullfilled.toml", "other")
	if err != nil {
		t.Errorf("%#v", err)
	}
	if app.AppID != "other-app-id" {
		t.Errorf("%s is not equal %s", app.AppID, "other-app-id")
	}
	if app.PrincipalArn != "other-provider-arn" {
		t.Errorf("%s is not equal %s", app.PrincipalArn, "other-provider-arn")
	}
	if app.RoleArn != "other-role-arn" {
		t.Errorf("%s is not equal %s", app.RoleArn, "other-role-arn")
	}
}

func TestLoginCmdFetchConfigNoProfile(t *testing.T) {
	var err error
	_, _, err = fetchConfig("fixtures/fullfilled.toml", "none")
	if err.Error() != "none profile is not exists" {
		t.Error(err.Error())
	}
}
