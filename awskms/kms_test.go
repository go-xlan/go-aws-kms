package awskms_test

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/go-xlan/go-aws-kms/awskms"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/done"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
)

var envOptions *awskms.EnvOptions

func TestMain(m *testing.M) {
	envOptions = awskms.NewEnvOptions()
	envOptions.WithSessionToken("AWS_SESSION_TOKEN") // change the environment variable key name

	println(neatjsons.S(envOptions)) // print the environment options
	m.Run()
}

func TestAwsKms_Encrypt(t *testing.T) {
	region := aws.String(must.Nice(os.Getenv(envOptions.RegionID)))

	oneKms := kms.New(rese.P1(session.NewSession(&aws.Config{
		Credentials: awskms.NewEnvCredentials(envOptions),
		Region:      region,
		Logger:      awskms.NewSlogLogger(),
	})))
	encryptKeyID := must.Nice(os.Getenv(envOptions.EncryptKeyID))

	awsKms := awskms.NewAwsKms(oneKms, encryptKeyID)

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

func TestAwsKms_Encrypts(t *testing.T) {
	// Avoid using static string access/secret in GitHub. Still use env variables.
	regionString := must.Nice(os.Getenv(envOptions.RegionID))
	accessString := must.Nice(os.Getenv(envOptions.AccessKeyID))
	secretString := must.Nice(os.Getenv(envOptions.SecretAccessKey))
	xTokenString := os.Getenv(envOptions.SessionToken) // allow empty session token
	encryptKeyID := must.Nice(os.Getenv(envOptions.EncryptKeyID))

	region := aws.String(regionString)
	oneKms := kms.New(done.VCE(session.NewSession(&aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentials(accessString, secretString, xTokenString),
	})).Nice())
	awsKms := awskms.NewAwsKms(oneKms, encryptKeyID)

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
