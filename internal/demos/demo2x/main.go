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
