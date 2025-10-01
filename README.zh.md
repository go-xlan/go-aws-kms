[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/go-xlan/go-aws-kms/release.yml?branch=main&label=BUILD)](https://github.com/go-xlan/go-aws-kms/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-xlan/go-aws-kms)](https://pkg.go.dev/github.com/go-xlan/go-aws-kms)
[![Coverage Status](https://img.shields.io/coveralls/github/go-xlan/go-aws-kms/main.svg)](https://coveralls.io/github/go-xlan/go-aws-kms?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.22--1.25-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/go-xlan/go-aws-kms.svg)](https://github.com/go-xlan/go-aws-kms/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-xlan/go-aws-kms)](https://goreportcard.com/report/github.com/go-xlan/go-aws-kms)

# go-aws-kms

AWS KMS 加密和解密包，提供便捷操作。

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->
## 英文文档

[ENGLISH README](README.md)
<!-- TEMPLATE (ZH) END: LANGUAGE NAVIGATION -->

## 核心特性

🔐 **简单加密**: 易于使用的 AWS KMS 加密和解密操作
⚡ **双重模式**: 支持字节和 base64 编码的字符串操作
🔧 **环境配置**: 便捷的基于环境变量的配置方式
📝 **双重日志**: 包含 slog 和 zap 两种日志实现
🎯 **上下文支持**: 基于 AWS SDK v2 构建，带有上下文 API

## 安装

```bash
go get github.com/go-xlan/go-aws-kms
```

## 使用方法

### 字节加密

此示例展示了使用 slog 日志进行 AWS KMS 原始字节的加密/解密操作。

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/go-xlan/go-aws-kms/awskms"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
)

func main() {
	// Set AWS region
	region := must.Nice(os.Getenv("AWS_KMS_REGION_ID"))

	// Load AWS configuration
	cfg := rese.V1(config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_KMS_ACCESS_KEY"),
			os.Getenv("AWS_KMS_SECRET_KEY"),
			os.Getenv("AWS_SESSION_TOKEN"),
		)),
		config.WithLogger(awskms.NewSlogLogger()),
		config.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody),
	))

	// Create KMS client
	kmsClient := kms.NewFromConfig(cfg)

	// Get encryption key ID from environment
	encryptKeyID := must.Nice(os.Getenv("AWS_KMS_ENCRYPT_KEY_ID"))

	// Create AwsKms instance
	awsKms := awskms.NewAwsKms(kmsClient, encryptKeyID)

	// Encrypt and decrypt bytes
	fmt.Println("=== Bytes Encryption Example ===")
	plaintext := []byte("sensitive data")
	fmt.Printf("Message: %s\n", plaintext)

	ciphertext := rese.A1(awsKms.Encrypt(plaintext))
	fmt.Printf("Encrypted (bytes): %d bytes\n", len(ciphertext))

	decrypted := rese.A1(awsKms.Decrypt(ciphertext))
	fmt.Printf("Decrypted: %s\n", decrypted)

	fmt.Println("\n✅ Success!")
}
```

⬆️ **源码:** [源码](internal/demos/demo1x/main.go)

### Base64 字符串加密

此示例展示了使用基于环境变量配置的 AWS KMS 字符串加密/解密和 base64 编码操作。

```go
package main

import (
	"fmt"

	"github.com/go-xlan/go-aws-kms/awskms"
	"github.com/yyle88/rese"
)

func main() {
	// Create AwsKms instance from environment variables
	envOptions := awskms.NewEnvOptions()
	awsKms := rese.P1(awskms.NewAwsKmsFromEnv(envOptions))

	// Encrypt and decrypt string with base64
	fmt.Println("=== String Encryption Example ===")
	message := "secret message"
	fmt.Printf("Message: %s\n", message)

	encrypted := rese.C1(awsKms.Encrypts(message))
	fmt.Printf("Encrypted (base64): %s\n", encrypted)

	decryptedStr := rese.C1(awsKms.Decrypts(encrypted))
	fmt.Printf("Decrypted: %s\n", decryptedStr)

	fmt.Println("\n✅ Success!")
}
```

⬆️ **源码:** [源码](internal/demos/demo2x/main.go)

## API 参考

### 核心函数

- `NewAwsKms(client, keyID)` - 使用 KMS 客户端和加密 ID 创建 AwsKms 实例
- `Encrypt(plaintext)` - 加密字节，返回加密字节
- `Decrypt(ciphertext)` - 解密字节，返回明文字节
- `Encrypts(plaintext)` - 加密字符串，返回 base64 编码字符串
- `Decrypts(ciphertext)` - 解密 base64 字符串，返回明文字符串

### 环境函数

- `NewEnvOptions()` - 创建带有默认变量名的环境选项
- `NewAwsKmsFromEnv(options)` - 从环境变量创建 AwsKms 实例

### 日志函数

- `NewSlogLogger()` - 在 AWS SDK 操作中创建基于 slog 的日志记录器
- `NewZapLogger()` - 在 AWS SDK 操作中创建基于 zap 的日志记录器

## 示例

### 环境变量配置

**使用默认环境变量名:**
```go
envOptions := awskms.NewEnvOptions()
awsKms, _ := awskms.NewAwsKmsFromEnv(envOptions)
```

**自定义环境变量名:**
```go
envOptions := awskms.NewEnvOptions()
envOptions.
	WithRegionID("CUSTOM_REGION").
	WithAccessKeyID("CUSTOM_ACCESS_KEY").
	WithSecretAccessKey("CUSTOM_SECRET_KEY").
	WithEncryptKeyID("CUSTOM_KMS_KEY_ID")

awsKms, _ := awskms.NewAwsKmsFromEnv(envOptions)
```

### 日志配置

**使用 slog:**
```go
cfg, _ := config.LoadDefaultConfig(
	context.Background(),
	config.WithLogger(awskms.NewSlogLogger()),
	config.WithClientLogMode(aws.LogRequestWithBody | aws.LogResponseWithBody),
)
```

**使用 zap:**
```go
cfg, _ := config.LoadDefaultConfig(
	context.Background(),
	config.WithLogger(awskms.NewZapLogger()),
	config.WithClientLogMode(aws.LogRequestWithBody | aws.LogResponseWithBody),
)
```

### 快速字符串加密

**加密和解密字符串:**
```go
encrypted, _ := awsKms.Encrypts("secret message")
decrypted, _ := awsKms.Decrypts(encrypted)
```

### 直接字节操作

**处理原始字节:**
```go
plaintext := []byte("sensitive data")
ciphertext, _ := awsKms.Encrypt(plaintext)
decrypted, _ := awsKms.Decrypt(ciphertext)
```

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 许可证类型

MIT 许可证。详见 [LICENSE](LICENSE)。

---

## 🤝 项目贡献

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **发现问题？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **功能建议？** 创建 issue 讨论您的想法
- 📖 **文档疑惑？** 报告问题，帮助我们改进文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，帮助我们优化性能
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **反馈意见？** 欢迎提出建议和意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：为面向用户的更改更新文档，并使用有意义的提交消息
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Merge Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Merge Request 和报告问题来为此项目做出贡献。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**祝你用这个包编程愉快！** 🎉🎉🎉

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## GitHub 标星点赞

[![Stargazers](https://starchart.cc/go-xlan/go-aws-kms.svg?variant=adaptive)](https://starchart.cc/go-xlan/go-aws-kms)
