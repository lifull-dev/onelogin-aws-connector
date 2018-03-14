# OneLogin AWS Connector

## Installation

```bash
go get github.com/lifull-dev/onelogin-aws-connector
```

## Using the OneLogin AWS Connector

OneLogin AWS Connector provides to create AWS credentials with OneLogin SAML.
If you want to use this command, You need to do setup to OneLogin and AWS accounts.

How to setup OneLogin and AWS accounts is the following [OneLogin Help Center](https://support.onelogin.com/hc/en-us/sections/200708060-Amazon-Web-Services)

## onelogin-aws-connector

### Global Options

## onelogin-aws-connector init

Init command initialize OneLogin API settings.

### Init Command Line Options

#### --endpoint `<us|eu>`

OneLogin API Server

#### --client-token `string`

OneLogin API Client Token

#### --client-secret `string`

OneLogin API Client Secret

#### --subdomain `string`

OneLogin Service Subdomain

#### --username-or-email `string`

OneLogin Login Username or Email

## onelogin-aws-connector configure

Configure command configure OneLogin and AWS connection settings.

### Configure Command Line Options

#### --app-id `string`

OneLogin AppID

#### --provider-arn `string`

AWS Provider ARN connected to OneLogin AppID

#### --role-arn `string`

AWS Role ARN

#### --aws-profile string

AWS Profile Name (default "default")

## onelogin-aws-connector login

Login command makes AWS credentials with OneLogin SAML.

### Login Command Line Options

#### --aws-profile `string`

AWS Profile Name (default "default")

#### --aws-region `string`

AWS Region Name
