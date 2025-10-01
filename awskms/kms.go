// Package awskms: AWS KMS encryption and decryption package with convenient operations
// Provides simple interfaces to encrypt and decrypt data using AWS KMS Service
// Supports both bytes and string operations with automatic base64 encoding
// Built on AWS SDK Go v2 in modern cloud-native applications
//
// awskms: AWS KMS 加密和解密包，提供便捷操作
// 提供简单接口来加密和解密 AWS KMS 服务的数据
// 支持字节和字符串操作，带有自动 base64 编码
// 基于 AWS SDK Go v2 构建，适用于现代云原生应用
package awskms

import (
	"context"
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
)

// AwsKms provides encryption and decryption operations using AWS KMS
// Wraps AWS KMS client to enable convenient encrypt/decrypt operations
// Maintains encryption ID and client instance in seamless cryptographic operations
//
// AwsKms 使用 AWS KMS 提供加密和解密操作
// 封装 AWS KMS 客户端以提供便捷的加解密操作
// 维护加密 ID 和客户端实例以实现无缝加密操作
type AwsKms struct {
	client       *kms.Client // AWS KMS client to handle API operations // AWS KMS 客户端用于处理 API 操作
	encryptKeyID string      // KMS ID used in encryption // 用于加密的 KMS ID
}

// NewAwsKms creates a new AwsKms instance with given KMS client and encryption ID
// Validates client address using must.Full and ID using must.Nice to ensure input safe
// Returns configured AwsKms instance suited in encryption and decryption operations
//
// NewAwsKms 使用给定的 KMS 客户端和加密 ID 创建新的 AwsKms 实例
// 使用 must.Full 验证客户端地址和 must.Nice 验证 ID 以确保输入安全
// 返回配置好的 AwsKms 实例，适用于加密和解密操作
func NewAwsKms(client *kms.Client, encryptKeyID string) *AwsKms {
	return &AwsKms{
		client:       must.Full(client),
		encryptKeyID: must.Nice(encryptKeyID),
	}
}

// Encrypt encrypts plaintext bytes using AWS KMS with configured encryption ID
// Calls AWS KMS Encrypt API and wraps exception with erero in enhanced context
// Returns encrypted ciphertext blob suitable in storage and transmission
//
// Encrypt 使用配置的加密 ID 通过 AWS KMS 加密明文字节
// 调用 AWS KMS Encrypt API 并使用 erero 包装异常以增强上下文
// 返回适合存储和传输的加密密文块
func (a *AwsKms) Encrypt(plaintext []byte) ([]byte, error) {
	res, err := a.client.Encrypt(context.Background(), &kms.EncryptInput{
		KeyId:     &a.encryptKeyID,
		Plaintext: plaintext,
	})
	if err != nil {
		return nil, erero.Wro(err)
	}
	return res.CiphertextBlob, nil
}

// Decrypt decrypts ciphertext blob using AWS KMS with automatic ID detection
// KMS auto-detects the encryption ID from ciphertext metadata in decryption
// Returns decrypted plaintext bytes and wraps exception with erero in enhanced context
//
// Decrypt 使用 AWS KMS 解密密文块，自动检测 ID
// KMS 从密文元数据中自动检测用于解密的加密 ID
// 返回解密的明文字节并使用 erero 包装异常以增强上下文
func (a *AwsKms) Decrypt(ciphertextBlob []byte) ([]byte, error) {
	res, err := a.client.Decrypt(context.Background(), &kms.DecryptInput{
		CiphertextBlob: ciphertextBlob,
	})
	if err != nil {
		return nil, erero.Wro(err)
	}
	return res.Plaintext, nil
}

// Encrypts encrypts plaintext string using AWS KMS and returns base64 encoded outcome
// Converts string to bytes, encrypts with KMS, and encodes outcome in base64 format
// Returns base64 encoded ciphertext string suited in text-based storage and transmission
//
// Encrypts 使用 AWS KMS 加密明文字符串并返回 base64 编码结果
// 将字符串转换为字节，使用 KMS 加密，并将结果编码为 base64 格式
// 返回 base64 编码的密文字符串，适用于基于文本的存储和传输
func (a *AwsKms) Encrypts(plaintext string) (string, error) {
	ciphertextBlob, err := a.Encrypt([]byte(plaintext))
	if err != nil {
		return "", erero.Wro(err)
	}
	return base64.StdEncoding.EncodeToString(ciphertextBlob), nil
}

// Decrypts decrypts base64 encoded ciphertext string using AWS KMS
// Decodes base64 string to bytes, decrypts with KMS, and converts outcome to string
// Returns decrypted plaintext string and wraps exception with erero in enhanced context
//
// Decrypts 使用 AWS KMS 解密 base64 编码的密文字符串
// 将 base64 字符串解码为字节，使用 KMS 解密，并将结果转换为字符串
// 返回解密的明文字符串并使用 erero 包装异常以增强上下文
func (a *AwsKms) Decrypts(cipherText string) (string, error) {
	ciphertextBlob, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", erero.Wro(err)
	}
	plaintext, err := a.Decrypt(ciphertextBlob)
	if err != nil {
		return "", erero.Wro(err)
	}
	return string(plaintext), nil
}
