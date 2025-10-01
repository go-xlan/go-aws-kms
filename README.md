[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/go-xlan/go-aws-kms/release.yml?branch=main&label=BUILD)](https://github.com/go-xlan/go-aws-kms/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-xlan/go-aws-kms)](https://pkg.go.dev/github.com/go-xlan/go-aws-kms)
[![Coverage Status](https://img.shields.io/coveralls/github/go-xlan/go-aws-kms/main.svg)](https://coveralls.io/github/go-xlan/go-aws-kms?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.22--1.25-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/go-xlan/go-aws-kms.svg)](https://github.com/go-xlan/go-aws-kms/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-xlan/go-aws-kms)](https://goreportcard.com/report/github.com/go-xlan/go-aws-kms)

# go-aws-kms

AWS KMS encryption and decryption package with convenient operations.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->
## CHINESE README

[中文说明](README.zh.md)
<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## Main Features

🔐 **Simple Encryption**: Easy-to-use AWS KMS encryption and decryption operations
⚡ **Dual Modes**: Supports both bytes and base64-encoded string operations
🔧 **Environment Config**: Convenient environment variable-based configuration
📝 **Dual Logging**: Includes both slog and zap logging implementations
🎯 **Context Support**: Built on AWS SDK v2 with context-based API

## Installation

```bash
go get github.com/go-xlan/go-aws-kms
```

## Usage

### Bytes Encryption

This example shows AWS KMS encryption/decryption with raw bytes using slog logging.

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

⬆️ **Source:** [Source](internal/demos/demo1x/main.go)

### String Encryption with Base64

This example shows AWS KMS string encryption/decryption with base64 encoding using environment-based configuration.

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

⬆️ **Source:** [Source](internal/demos/demo2x/main.go)

## API Reference

### Core Functions

- `NewAwsKms(client, keyID)` - Create AwsKms instance with KMS client and encryption ID
- `Encrypt(plaintext)` - Encrypt bytes, returns encrypted bytes
- `Decrypt(ciphertext)` - Decrypt bytes, returns plaintext bytes
- `Encrypts(plaintext)` - Encrypt string, returns base64-encoded string
- `Decrypts(ciphertext)` - Decrypt base64 string, returns plaintext string

### Environment Functions

- `NewEnvOptions()` - Create environment options with default variable names
- `NewAwsKmsFromEnv(options)` - Create AwsKms instance from environment variables

### Logger Functions

- `NewSlogLogger()` - Create slog-based logger in AWS SDK operations
- `NewZapLogger()` - Create zap-based logger in AWS SDK operations

## Examples

### Environment-Based Configuration

**Using default environment variable names:**
```go
envOptions := awskms.NewEnvOptions()
awsKms, _ := awskms.NewAwsKmsFromEnv(envOptions)
```

**Custom environment variable names:**
```go
envOptions := awskms.NewEnvOptions()
envOptions.
	WithRegionID("CUSTOM_REGION").
	WithAccessKeyID("CUSTOM_ACCESS_KEY").
	WithSecretAccessKey("CUSTOM_SECRET_KEY").
	WithEncryptKeyID("CUSTOM_KMS_KEY_ID")

awsKms, _ := awskms.NewAwsKmsFromEnv(envOptions)
```

### Logging Configuration

**Using slog:**
```go
cfg, _ := config.LoadDefaultConfig(
	context.Background(),
	config.WithLogger(awskms.NewSlogLogger()),
	config.WithClientLogMode(aws.LogRequestWithBody | aws.LogResponseWithBody),
)
```

**Using zap:**
```go
cfg, _ := config.LoadDefaultConfig(
	context.Background(),
	config.WithLogger(awskms.NewZapLogger()),
	config.WithClientLogMode(aws.LogRequestWithBody | aws.LogResponseWithBody),
)
```

### Quick String Encryption

**Encrypt and decrypt strings:**
```go
encrypted, _ := awsKms.Encrypts("secret message")
decrypted, _ := awsKms.Decrypts(encrypted)
```

### Direct Bytes Operations

**Working with raw bytes:**
```go
plaintext := []byte("sensitive data")
ciphertext, _ := awsKms.Encrypt(plaintext)
decrypted, _ := awsKms.Decrypt(ciphertext)
```

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 License

MIT License. See [LICENSE](LICENSE).

---

## 🤝 Contributing

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Found a mistake?** Open an issue on GitHub with reproduction steps
- 💡 **Have a feature idea?** Create an issue to discuss the suggestion
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share the use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize through reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo to get new releases and features
- 🌟 **Success stories?** Share how this package improved the workflow
- 💬 **Feedback?** We welcome suggestions and comments

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage UI).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement the changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation to support client-facing changes and use significant commit messages
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a merge request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project via submitting merge requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉🎉🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/go-xlan/go-aws-kms.svg?variant=adaptive)](https://starchart.cc/go-xlan/go-aws-kms)
