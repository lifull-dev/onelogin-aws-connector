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
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/lifull-dev/onelogin-aws-connector/aws/configuration"
	"github.com/lifull-dev/onelogin-aws-connector/cmd/config"
	"github.com/lifull-dev/onelogin-aws-connector/cmd/login"
	"github.com/lifull-dev/onelogin-aws-connector/onelogin"
	"github.com/lifull-dev/onelogin-aws-connector/onelogin/samlassertion"
)

var region string
var force bool

type LoginEvent struct {
	reader *bufio.Reader
}

func NewLoginEvent(reader *bufio.Reader) *LoginEvent {
	return &LoginEvent{
		reader: reader,
	}
}

func (m *LoginEvent) ChooseDeviceIndex(devices []samlassertion.GenerateResponseFactorDevice) (int, error) {
	if debug {
		fmt.Println("")
		log.Println("MFA Devices:")
		for _, device := range devices {
			log.Printf("  %v:\t\t%v\n", device.DeviceID, device.DeviceType)
		}
	}
	length := len(devices)
	selected := length
	for {
		fmt.Println("--------")
		for i, device := range devices {
			fmt.Printf("%d : %s\n", i, device.DeviceType)
		}
		fmt.Println("--------")
		fmt.Print("Select your MFA device: ")
		tmp, err := m.reader.ReadString('\n')
		if err != nil {
			return 0, err
		}
		tmp = strings.Trim(tmp, "\n")
		if tmp == "" {
			continue
		}
		selected, err = strconv.Atoi(tmp)
		if err != nil {
			return 0, err
		}
		if selected < length && selected >= 0 {
			break
		}
	}
	return selected, nil
}

func (m *LoginEvent) InputMFAToken() (string, error) {
	var token string
	var err error
	for {
		fmt.Print("Enter your MFA token: ")
		token, err = m.reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		token = strings.Trim(token, "\n")
		if token != "" {
			break
		}
	}
	return token, nil
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to AWS with OneLogin",
	Long:  `Login is CLI Command to Create AWS Credentials with OneLogin`,
	Run: func(cmd *cobra.Command, args []string) {
		if awsProfile == "" {
			awsProfile = "default"
		}
		err := cached(awsProfile, func() (*sts.Credentials, error) {
			service, app, err := fetchConfig(configFile, awsProfile)
			if err != nil {
				return nil, err
			}
			if debug {
				log.Println("OneLogin Configuration:")
				log.Printf("  Endpoint:\t\t%v\n", service.Endpoint)
				log.Printf("  ClientToken:\t\t%v\n", service.ClientToken)
				log.Printf("  ClientSecret:\t%v\n", service.ClientSecret)
			}

			onelogin.CacheDir = cacheDir
			config := onelogin.NewConfig(service.Endpoint, service.ClientToken, service.ClientSecret)
			if force {
				config.Credentials.Credentials = nil
			}
			if err := config.Save(); err != nil {
				return nil, err
			}
			if debug {
				creds, _ := config.Credentials.Get()
				log.Println("OneLogin Credentials:")
				log.Printf("  AccessToken:\t\t%v\n", creds.AccessToken)
				log.Printf("  RefreshToken:\t%v\n", creds.RefreshToken)
				log.Printf("  CreatedAt:\t\t%v\n", creds.CreatedAt)
				log.Printf("  AccessExpiresAt:\t%v\n", creds.AccessExpiresAt)
				log.Printf("  RefreshExpiresAt:\t%v\n", creds.RefreshExpiresAt)
			}

			fmt.Print("Enter your password: ")
			tmp, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return nil, err
			}
			password := string(tmp)
			fmt.Println("")
			duration := app.DurationSeconds
			if duration == 0 {
				duration = 3600
			}
			if debug {
				fmt.Println("")
				log.Println("Login Parameters:")
				log.Printf("  Subdomain:\t\t%v\n", service.Subdomain)
				log.Printf("  AppID:\t\t%v\n", app.AppID)
				log.Printf("  UsernameOrEmail:\t%v\n", service.UsernameOrEmail)
				log.Printf("  Password:\t\t%v\n", password)
				log.Printf("  PrincipalArn:\t%v\n", app.PrincipalArn)
				log.Printf("  RoleArn:\t\t%v\n", app.RoleArn)
				log.Printf("  DurationSeconds:\t%v\n", duration)
			}
			l := login.New(config, &login.Parameters{
				UsernameOrEmail: service.UsernameOrEmail,
				Password:        password,
				AppID:           app.AppID,
				Subdomain:       service.Subdomain,
				PrincipalArn:    app.PrincipalArn,
				RoleArn:         app.RoleArn,
				DurationSeconds: duration,
			})
			creds, err := l.Login(NewLoginEvent(bufio.NewReader(os.Stdin)))

			if err != nil {
				return nil, err
			}

			if debug {
				log.Println("AWS Credentials:")
				log.Printf("  AccessKeyId:\t%v\n", *creds.AccessKeyId)
				log.Printf("  SecretAccessKey:\t%v\n", *creds.SecretAccessKey)
				log.Printf("  SessionToken:\t%v\n", *creds.SessionToken)
				log.Printf("  Expiration:\t\t%v\n", creds.Expiration)
			}
			options := map[string]string{
				"aws_access_key_id":     *creds.AccessKeyId,
				"aws_secret_access_key": *creds.SecretAccessKey,
				"aws_session_token":     *creds.SessionToken,
			}
			awsCredentials := configuration.NewCredentials(awsDir, awsProfile)
			_ = awsCredentials.Save(options)
			if region != "" {
				awsConfig := configuration.NewConfig(awsDir, awsProfile)
				_ = awsConfig.Save(region)
			}
			return creds, nil
		})
		if err != nil {
			errorExit(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringVarP(&region, "aws-region", "", "", "AWS Region")
	loginCmd.Flags().BoolVarP(&force, "force", "", false, "Force refresh AWS credentials if credentials enabled")
	loginCmd.Flags().StringVarP(&awsProfile, "aws-profile", "", awsProfile, "aws profile name")
}

func fetchConfig(file string, profile string) (config.ServiceConfig, config.AppConfig, error) {
	c, err := config.Load(file)
	if err != nil {
		return config.ServiceConfig{}, config.AppConfig{}, err
	}
	app, ok := c.App[profile]
	if !ok {
		return emptyConfig(fmt.Sprintf("%s profile is not exists", profile))
	}

	service := c.Service["default"]
	if service.Endpoint == "" {
		return emptyConfig("Endpoint is not exists")
	}

	if service.ClientToken == "" {
		return emptyConfig("ClientToken is not exists")
	}

	if service.ClientSecret == "" {
		return emptyConfig("ClientSecret is not exists")
	}

	if service.Subdomain == "" {
		return emptyConfig("Subdomain is not exists")
	}
	return *service, *app, nil
}

func emptyConfig(message string) (config.ServiceConfig, config.AppConfig, error) {
	return config.ServiceConfig{}, config.AppConfig{}, errors.Errorf(message)
}

func cached(profile string, block func() (*sts.Credentials, error)) error {
	file := path.Join(cacheDir, fmt.Sprintf("aws.%s.cache", profile))
	if !force {
		var c *sts.Credentials
		if _, err := toml.DecodeFile(file, &c); err != nil {
			if err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			}
		} else {
			if c.Expiration != nil {
				now := time.Now()
				if now.Before(*c.Expiration) {
					if debug {
						log.Println("use aws credentials cache")
					}
					return nil
				}
			}
		}
	}
	c, err := block()
	if err != nil {
		return err
	}
	fd, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fd.Close()
	encoder := toml.NewEncoder(fd)
	if err := encoder.Encode(c); err != nil {
		return err
	}
	return nil
}
