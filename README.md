# OneLogin AWS Connector

## Installation

```bash
go install github.com/lifull-dev/onelogin-aws-connector@latest
```

## Using the OneLogin AWS Connector

OneLogin AWS Connector provides to create AWS credentials with OneLogin SAML.
If you want to use this command, You need to do setup to OneLogin and AWS accounts.

How to setup OneLogin and AWS accounts is the following [OneLogin Help Center](https://support.onelogin.com/hc/en-us/sections/200708060-Amazon-Web-Services)

## onelogin-aws-connector

### Global Options

## onelogin-aws-connector init

Init command initialize OneLogin API settings.

```bash
onelogin-aws-connector init \
    --endpoint us \
    --client-token [TOKEN] \
    --client-secret [SECRET] \
    --subdomain [SUBDOMAIN] \
    --username-or-email [USERNAME_OR_EMAIL]
```

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

```bash
onelogin-aws-connector configure \
    --app-id [APP_ID] \
    --role-arn [AWS_ROLE_ARN] \
    --provider-arn [AWS_SAML_PROVIDER_ARN] \
    --aws-profile [AWS_PROFILE_NAME]
```

### Configure Command Line Options

#### --app-id `string`

OneLogin AppID

#### --provider-arn `string`

AWS Provider ARN connected to OneLogin AppID

#### --role-arn `string`

AWS Role ARN

#### --duration `int`

The value can range from 900 seconds (15 minutes) to maximum session duration setting (default 3600 seconds (1 hour)).

#### --aws-profile string

AWS Profile Name (default "default")

## onelogin-aws-connector login

Login command makes AWS credentials with OneLogin SAML.

### Login Command Line Options

```bash
onelogin-aws-connector login \
    --aws-profile [AWS_PROFILE_NAME] \
    --aws-region [AWS_REGION_NAME]
```

#### --aws-profile `string`

AWS Profile Name (default "default")

#### --aws-region `string`

AWS Region Name
