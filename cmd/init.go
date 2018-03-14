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
	"fmt"
	"log"

	"github.com/lifull-dev/onelogin-aws-connector/cmd/config"
	"github.com/spf13/cobra"
)

var endpoint string
var clientToken string
var clientSecret string
var subdomain string
var usernameOrEmail string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialze settings for call to onelogin api ",
	Long:  `Init is initializing settings for onelogin api.`,
	Run: func(cmd *cobra.Command, args []string) {
		if endpoint != "" {
			endpoint = fmt.Sprintf("api.%s.onelogin.com", endpoint)
		}
		if err := initServiceConfig(configFile, "default"); err != nil {
			errorExit(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(&endpoint, "endpoint", "", "", "OneLogin API Server")
	initCmd.Flags().StringVarP(&clientToken, "client-token", "", "", "OneLogin API Client Token")
	initCmd.Flags().StringVarP(&clientSecret, "client-secret", "", "", "OneLogin API Client Secret")
	initCmd.Flags().StringVarP(&subdomain, "subdomain", "", "", "OneLogin Service Subdomain")
	initCmd.Flags().StringVarP(&usernameOrEmail, "username-or-email", "", "", "OneLogin Login Username or Email")
}

func initServiceConfig(file string, profile string) error {
	c, err := config.Load(file)
	if err != nil {
		return err
	}
	serviceConfig, ok := c.Service[profile]
	if !ok {
		serviceConfig = &config.ServiceConfig{}
	}
	if endpoint != "" {
		serviceConfig.Endpoint = endpoint
	}
	if clientToken != "" {
		serviceConfig.ClientToken = clientToken
	}
	if clientSecret != "" {
		serviceConfig.ClientSecret = clientSecret
	}
	if subdomain != "" {
		serviceConfig.Subdomain = subdomain
	}
	if usernameOrEmail != "" {
		serviceConfig.UsernameOrEmail = usernameOrEmail
	}
	c.Service["default"] = serviceConfig
	if err := c.Save(); err != nil {
		return err
	}
	if debug {
		log.Printf("ServiceConfig: %#v\n", serviceConfig)
	}
	return nil
}
