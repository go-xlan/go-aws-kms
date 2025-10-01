package awskms

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
)

// EnvOptions defines environment variable names in AWS configuration
// Contains mappings to specify region, credentials, and encryption ID with flexible customization
// Supports chain configuration pattern to set custom environment variable names
//
// EnvOptions 定义 AWS 配置的环境变量名称
// 包含区域、凭证和加密 ID 的映射，支持灵活定制
// 支持链式配置模式以设置自定义环境变量名称
type EnvOptions struct {
	RegionID        string // Environment variable name to specify AWS region ID // AWS 区域 ID 的环境变量名
	AccessKeyID     string // Environment variable name to specify AWS access ID // AWS 访问 ID 的环境变量名
	SecretAccessKey string // Environment variable name to specify AWS secret access code // AWS 密钥访问码的环境变量名
	SessionToken    string // Environment variable name to specify AWS session token // AWS 会话令牌的环境变量名
	EncryptKeyID    string // Environment variable name to specify AWS KMS encrypt ID // AWS KMS 加密 ID 的环境变量名
}

// NewEnvOptions creates EnvOptions with default environment variable names
// Default names follow AWS_KMS_* pattern in quick identification
//
// NewEnvOptions 创建具有默认环境变量名称的 EnvOptions
// 默认名称遵循 AWS_KMS_* 模式以便于识别
func NewEnvOptions() *EnvOptions {
	return &EnvOptions{
		RegionID:        "AWS_KMS_REGION_ID",
		AccessKeyID:     "AWS_KMS_ACCESS_KEY",
		SecretAccessKey: "AWS_KMS_SECRET_KEY",
		SessionToken:    "AWS_KMS_SESSION_TOKEN",
		EncryptKeyID:    "AWS_KMS_ENCRYPT_KEY_ID",
	}
}

// WithRegionID sets custom environment variable name to specify region ID
// Returns self in method chaining
//
// WithRegionID 设置区域 ID 的自定义环境变量名
// 返回自身以支持链式调用
func (op *EnvOptions) WithRegionID(keyName string) *EnvOptions {
	op.RegionID = keyName
	return op
}

// WithAccessKeyID sets custom environment variable name to specify access ID
// Returns self in method chaining
//
// WithAccessKeyID 设置访问 ID 的自定义环境变量名
// 返回自身以支持链式调用
func (op *EnvOptions) WithAccessKeyID(keyName string) *EnvOptions {
	op.AccessKeyID = keyName
	return op
}

// WithSecretAccessKey sets custom environment variable name to specify secret access code
// Returns self in method chaining
//
// WithSecretAccessKey 设置密钥访问码的自定义环境变量名
// 返回自身以支持链式调用
func (op *EnvOptions) WithSecretAccessKey(keyName string) *EnvOptions {
	op.SecretAccessKey = keyName
	return op
}

// WithSessionToken sets custom environment variable name to specify session token
// Returns self in method chaining
//
// WithSessionToken 设置会话令牌的自定义环境变量名
// 返回自身以支持链式调用
func (op *EnvOptions) WithSessionToken(keyName string) *EnvOptions {
	op.SessionToken = keyName
	return op
}

// WithEncryptKeyID sets custom environment variable name to specify encryption ID
// Returns self in method chaining
//
// WithEncryptKeyID 设置加密 ID 的自定义环境变量名
// 返回自身以支持链式调用
func (op *EnvOptions) WithEncryptKeyID(keyName string) *EnvOptions {
	op.EncryptKeyID = keyName
	return op
}

// NewAwsKmsFromEnv creates AwsKms instance from environment variables using provided options
// Reads AWS credentials, region, and encryption ID from environment variables specified in options
// Validates required environment variables using must.Nice and returns exception when missing encryption ID
// Configures AWS SDK v2 client with static credentials and creates ready-to-use AwsKms instance
//
// NewAwsKmsFromEnv 使用提供的选项从环境变量创建 AwsKms 实例
// 从选项中指定的环境变量读取 AWS 凭证、区域和加密 ID
// 使用 must.Nice 验证必需环境变量，对缺失的加密 ID 返回异常
// 使用静态凭证配置 AWS SDK v2 客户端并创建可用的 AwsKms 实例
func NewAwsKmsFromEnv(options *EnvOptions) (*AwsKms, error) {
	region := must.Nice(os.Getenv(options.RegionID))
	accessKey := must.Nice(os.Getenv(options.AccessKeyID))
	secretKey := must.Nice(os.Getenv(options.SecretAccessKey))
	sessionToken := os.Getenv(options.SessionToken) // allow empty session token // 允许空会话令牌
	encryptKeyID := os.Getenv(options.EncryptKeyID)
	if encryptKeyID == "" {
		return nil, erero.New("encrypt key ID environment variable is none")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			sessionToken,
		)),
	)
	if err != nil {
		return nil, erero.Wro(err)
	}

	return NewAwsKms(kms.NewFromConfig(cfg), encryptKeyID), nil
}
