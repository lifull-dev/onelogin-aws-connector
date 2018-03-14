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
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	awsProfile string
	debug      bool
)

var (
	configFile string
	cacheFile  string
	cacheDir   string
	awsDir     string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "onelogin-aws-connector",
	Short: "Generate AWS Credentials with OneLogin SAML",
	Long: `This is a CLI command to generate AWS credentials with OneLogin SAML
This command write to credentials to ~/.aws/config and ~/.aws/credentials.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		errorExit(err)
	}
}

func init() {
	home, err := homedir.Dir()
	if err != nil {
		errorExit(err)
	}
	dir := path.Join(home, ".onelogin-aws-connector")
	if err := os.Mkdir(dir, 0700); err != nil {
		if !os.IsExist(err) {
			errorExit(err)
		}
	}
	awsDir = path.Join(home, ".aws")
	if err := os.Mkdir(dir, 0700); err != nil {
		if !os.IsExist(err) {
			errorExit(err)
		}
	}
	cacheDir = path.Join(dir, "cache")
	if err := os.Mkdir(cacheDir, 0700); err != nil {
		if !os.IsExist(err) {
			errorExit(err)
		}
	}
	configFile = path.Join(dir, "config.toml")
	cacheFile = path.Join(cacheDir, "response.cache")
	awsProfile = os.Getenv("AWS_PROFILE")
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "debug mode")
}
