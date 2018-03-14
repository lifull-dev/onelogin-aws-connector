// Copyright © 2017 LIFULL Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/lifull-dev/onelogin-aws-connector/cmd/config"
)

var appID string
var roleArn string
var principalArn string

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Add config to login to onelogin api",
	Long:  `Configure is add config to login to onelogin api.`,
	Run: func(cmd *cobra.Command, args []string) {
		if awsProfile == "" {
			awsProfile = "default"
		}
		if err := initAppConfig(configFile, awsProfile); err != nil {
			errorExit(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(configureCmd)
	configureCmd.Flags().StringVarP(&appID, "app-id", "", "", "OneLogin AppID")
	configureCmd.Flags().StringVarP(&roleArn, "role-arn", "", "", "Login Target AWS Role ARN")
	configureCmd.Flags().StringVarP(&principalArn, "principal-arn", "", "", "AWS Provider ARN connected to OneLogin AppID")
	configureCmd.Flags().StringVarP(&awsProfile, "aws-profile", "", awsProfile, "aws profile name")
}

func initAppConfig(file string, profile string) error {
	c, err := config.Load(file)
	if err != nil {
		return err
	}
	appConfig, ok := c.App[profile]
	if !ok {
		appConfig = &config.AppConfig{}
	}
	if appID != "" {
		appConfig.AppID = appID
	}
	if roleArn != "" {
		appConfig.RoleArn = roleArn
	}
	if principalArn != "" {
		appConfig.PrincipalArn = principalArn
	}
	serviceProfile := "default"
	if _, ok := c.Service[serviceProfile]; !ok {
		return errors.Errorf("There is no initialized service. Please run `onelogin-aws-connector init`")
	}
	c.App[profile] = appConfig
	if err := c.Save(); err != nil {
		return err
	}
	if debug {
		log.Printf("AppConfig: %#v\n", appConfig)
	}
	return nil
}
