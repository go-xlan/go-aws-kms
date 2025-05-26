package awskms

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
)

type AwsKms struct {
	kms          *kms.KMS
	encryptKeyID string
}

func NewAwsKms(kms *kms.KMS, encryptKeyID string) *AwsKms {
	return &AwsKms{
		kms:          must.Nice(kms),
		encryptKeyID: must.Nice(encryptKeyID),
	}
}

func (AWS *AwsKms) Encrypt(plaintext []byte) ([]byte, error) {
	res, err := AWS.kms.Encrypt(&kms.EncryptInput{KeyId: aws.String(AWS.encryptKeyID), Plaintext: plaintext})
	if err != nil {
		return nil, erero.Wro(err)
	}
	return res.CiphertextBlob, nil
}

func (AWS *AwsKms) Decrypt(ciphertextBlob []byte) ([]byte, error) {
	res, err := AWS.kms.Decrypt(&kms.DecryptInput{CiphertextBlob: ciphertextBlob})
	if err != nil {
		return nil, erero.Wro(err)
	}
	return res.Plaintext, nil
}

func (AWS *AwsKms) Encrypts(plaintext string) (string, error) {
	ciphertextBlob, err := AWS.Encrypt([]byte(plaintext))
	if err != nil {
		return "", erero.Wro(err)
	}
	return base64.StdEncoding.EncodeToString(ciphertextBlob), nil
}

func (AWS *AwsKms) Decrypts(cipherText string) (string, error) {
	ciphertextBlob, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", erero.Wro(err)
	}
	plaintext, err := AWS.Decrypt(ciphertextBlob)
	if err != nil {
		return "", erero.Wro(err)
	}
	return string(plaintext), nil
}
