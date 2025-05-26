package awskms

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

// NewEnvCredentials creates a new credentials provider that retrieves
// cp from https://github.com/aws/aws-sdk-go/blob/163aada692ed32951f979aacf452ded4c03b8a7c/aws/credentials/env_provider.go#L36
func NewEnvCredentials(options *EnvOptions) *credentials.Credentials {
	return credentials.NewCredentials(NewEnvProvider(options))
}

type EnvProvider struct {
	success bool
	options *EnvOptions
}

func NewEnvProvider(options *EnvOptions) *EnvProvider {
	return &EnvProvider{
		success: false,
		options: options,
	}
}

// Retrieve retrieves the keys from the environment.
func (e *EnvProvider) Retrieve() (credentials.Value, error) {
	e.success = false

	access := os.Getenv(e.options.AccessKeyID)
	if access == "" {
		return credentials.Value{ProviderName: credentials.EnvProviderName}, credentials.ErrAccessKeyIDNotFound
	}
	secret := os.Getenv(e.options.SecretAccessKey)
	if secret == "" {
		return credentials.Value{ProviderName: credentials.EnvProviderName}, credentials.ErrSecretAccessKeyNotFound
	}

	e.success = true
	return credentials.Value{
		AccessKeyID:     access,
		SecretAccessKey: secret,
		SessionToken:    os.Getenv(e.options.SessionToken), // allow empty session token
		ProviderName:    credentials.EnvProviderName,
	}, nil
}

func (e *EnvProvider) IsExpired() bool {
	return !e.success
}

type EnvOptions struct {
	RegionID        string // Environment variable name for AWS region ID
	AccessKeyID     string // Environment variable name for AWS access key ID
	SecretAccessKey string // Environment variable name for AWS secret access key
	SessionToken    string // Environment variable name for AWS session token
	EncryptKeyID    string // Environment variable name for AWS KMS encrypt key ID
}

func NewEnvOptions() *EnvOptions {
	return &EnvOptions{
		RegionID:        "AWS_KMS_REGION_ID",
		AccessKeyID:     "AWS_KMS_ACCESS_KEY",
		SecretAccessKey: "AWS_KMS_SECRET_KEY",
		SessionToken:    "AWS_KMS_SESSION_TOKEN",
		EncryptKeyID:    "AWS_KMS_ENCRYPT_KEY_ID",
	}
}

func (op *EnvOptions) WithRegionID(keyName string) *EnvOptions {
	op.RegionID = keyName
	return op
}

func (op *EnvOptions) WithAccessKeyID(keyName string) *EnvOptions {
	op.AccessKeyID = keyName
	return op
}

func (op *EnvOptions) WithSecretAccessKey(keyName string) *EnvOptions {
	op.SecretAccessKey = keyName
	return op
}

func (op *EnvOptions) WithSessionToken(keyName string) *EnvOptions {
	op.SessionToken = keyName
	return op
}

func (op *EnvOptions) WithEncryptKeyID(keyName string) *EnvOptions {
	op.EncryptKeyID = keyName
	return op
}
