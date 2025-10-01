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
