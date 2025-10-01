package awskms_test

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/go-xlan/go-aws-kms/awskms"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
)

var envOptions *awskms.EnvOptions

// TestMain sets up test environment and validates AWS credentials access
// Configures environment variable options with custom session token name
// Skips each test when AWS credentials are not configured in CI environments
// Prints environment options in debugging when credentials are accessible
//
// TestMain 设置测试环境并验证 AWS 凭证访问
// 配置环境变量选项，带有自定义会话令牌名称
// 在 CI 环境中当 AWS 凭证未配置时跳过每个测试
// 当凭证可访问时打印环境选项以供调试
func TestMain(m *testing.M) {
	envOptions = awskms.NewEnvOptions()
	envOptions.WithSessionToken("AWS_SESSION_TOKEN") // change the environment variable key name // 更改环境变量键名

	// Check if AWS credentials are configured // 检查 AWS 凭证是否配置
	if os.Getenv(envOptions.RegionID) == "" {
		println("Skipping tests: AWS KMS credentials not configured")
		println("Required environment variables:")
		println(neatjsons.S(envOptions))
		os.Exit(0) // Exit with success to skip tests // 成功退出，跳过测试
	}

	println(neatjsons.S(envOptions)) // print the environment options // 打印环境选项
	m.Run()
}

// TestAwsKms_Encrypt tests bytes encryption and decryption operations
// Creates AWS KMS client with static credentials and SlogLogger in operation logging
// Verifies encrypt/decrypt process maintains data intact in bytes messages
// Uses must.Nice in environment variable validation and require in test assertions
//
// TestAwsKms_Encrypt 测试字节加密和解密操作
// 使用静态凭证和 SlogLogger 创建 AWS KMS 客户端以记录操作
// 验证加密/解密过程保持字节消息的数据完整
// 使用 must.Nice 验证环境变量和 require 进行测试断言
func TestAwsKms_Encrypt(t *testing.T) {
	region := must.Nice(os.Getenv(envOptions.RegionID))

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			must.Nice(os.Getenv(envOptions.AccessKeyID)),
			must.Nice(os.Getenv(envOptions.SecretAccessKey)),
			os.Getenv(envOptions.SessionToken),
		)),
		config.WithLogger(awskms.NewSlogLogger()),
		config.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody),
	)
	require.NoError(t, err)

	kmsClient := kms.NewFromConfig(cfg)
	encryptKeyID := must.Nice(os.Getenv(envOptions.EncryptKeyID))

	awsKms := awskms.NewAwsKms(kmsClient, encryptKeyID)

	t.Run("TestAwsKms_Encrypt", func(t *testing.T) {
		msg := []byte("test message")

		ciphertext, err := awsKms.Encrypt(msg)
		require.NoError(t, err)
		t.Logf("Ciphertext: %v", ciphertext)

		plaintext, err := awsKms.Decrypt(ciphertext)
		require.NoError(t, err)
		t.Logf("Plaintext: %s", plaintext)

		require.Equal(t, msg, plaintext)
	})
}

// TestAwsKms_Encrypts tests string encryption with base64 encoding operations
// Creates AWS KMS client with static credentials and ZapLogger in structured logging
// Verifies Encrypts/Decrypts process with base64 encoding maintains string data intact
// Demonstrates alternate implementation using uber/zap structured logging
//
// TestAwsKms_Encrypts 测试带 base64 编码的字符串加密操作
// 使用静态凭证和 ZapLogger 创建 AWS KMS 客户端以进行结构化日志记录
// 验证带 base64 编码的 Encrypts/Decrypts 过程保持字符串数据完整
// 演示使用 uber/zap 结构化日志记录的替代实现
func TestAwsKms_Encrypts(t *testing.T) {
	region := must.Nice(os.Getenv(envOptions.RegionID))
	accessKey := must.Nice(os.Getenv(envOptions.AccessKeyID))
	secretKey := must.Nice(os.Getenv(envOptions.SecretAccessKey))
	sessionToken := os.Getenv(envOptions.SessionToken) // allow empty session token // 允许空会话令牌
	encryptKeyID := must.Nice(os.Getenv(envOptions.EncryptKeyID))

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			sessionToken,
		)),
		config.WithLogger(awskms.NewZapLogger()),
		config.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody),
	)
	require.NoError(t, err)

	kmsClient := kms.NewFromConfig(cfg)
	awsKms := awskms.NewAwsKms(kmsClient, encryptKeyID)

	t.Run("TestAwsKms_Encrypts", func(t *testing.T) {
		msg := "test message"

		ciphertext, err := awsKms.Encrypts(msg)
		require.NoError(t, err)
		t.Log(ciphertext)

		plaintext, err := awsKms.Decrypts(ciphertext)
		require.NoError(t, err)
		t.Log(plaintext)

		require.Equal(t, msg, plaintext)
	})
}

// TestNewAwsKmsFromEnv tests environment-based AwsKms instance creation
// Uses NewAwsKmsFromEnv function in simplified client setup from environment variables
// Validates complete encryption workflow using convenient environment configuration
// Demonstrates the recommended pattern in production environment setup
//
// TestNewAwsKmsFromEnv 测试基于环境变量的 AwsKms 实例创建
// 使用 NewAwsKmsFromEnv 函数从环境变量简化客户端设置
// 使用便捷的环境配置验证完整的加密工作流程
// 演示生产环境设置的推荐模式
func TestNewAwsKmsFromEnv(t *testing.T) {
	awsKms := rese.P1(awskms.NewAwsKmsFromEnv(envOptions))

	t.Run("TestNewAwsKmsFromEnv", func(t *testing.T) {
		msg := "test message from env"

		ciphertext, err := awsKms.Encrypts(msg)
		require.NoError(t, err)
		t.Log("Encrypted:", ciphertext)

		plaintext, err := awsKms.Decrypts(ciphertext)
		require.NoError(t, err)
		t.Log("Decrypted:", plaintext)

		require.Equal(t, msg, plaintext)
	})
}
