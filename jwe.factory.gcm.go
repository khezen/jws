package jwc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
)

func newJWEA128GCM(protectedHeaders *JOSEHeaders, pubKey *rsa.PublicKey, plaintext []byte) (*JWE, error) {
	cek, cipherCEK, cipherCEKB64, err := GenerateCEK(16, protectedHeaders.Algorithm, pubKey)
	if err != nil {
		return nil, err
	}
	return newGCM(cek, cipherCEK, cipherCEKB64, protectedHeaders, pubKey, plaintext)
}

func newJWEA256GCM(protectedHeaders *JOSEHeaders, pubKey *rsa.PublicKey, plaintext []byte) (*JWE, error) {
	cek, cipherCEK, cipherCEKB64, err := GenerateCEK(32, protectedHeaders.Algorithm, pubKey)
	if err != nil {
		return nil, err
	}
	return newGCM(cek, cipherCEK, cipherCEKB64, protectedHeaders, pubKey, plaintext)
}

func newJWEA512GCM(protectedHeaders *JOSEHeaders, pubKey *rsa.PublicKey, plaintext []byte) (*JWE, error) {
	cek, cipherCEK, cipherCEKB64, err := GenerateCEK(64, protectedHeaders.Algorithm, pubKey)
	if err != nil {
		return nil, err
	}
	return newGCM(cek, cipherCEK, cipherCEKB64, protectedHeaders, pubKey, plaintext)
}

func newGCM(
	cek []byte,
	cipherCEK []byte,
	cipherCEKB64 string,
	protectedHeaders *JOSEHeaders,
	pubKey *rsa.PublicKey,
	plaintext []byte,
) (*JWE, error) {
	headersBytes, err := json.Marshal(protectedHeaders)
	if err != nil {
		return nil, err
	}
	headersB64 := base64.RawURLEncoding.EncodeToString(headersBytes)
	iv, ivB64, err := GenerateInitVector(aes.BlockSize)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(cek)
	if err != nil {
		return nil, err
	}
	mode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	authTag := []byte(String2ASCII(headersB64))
	authTagB64 := base64.RawURLEncoding.EncodeToString(authTag)
	ciphertext := mode.Seal(nil, iv, plaintext, nil)
	ciphertextB64 := base64.RawURLEncoding.EncodeToString(ciphertext)
	jwe := JWE{
		ProtectedB64:  headersB64,
		CipherCEKB64:  cipherCEKB64,
		InitVectorB64: ivB64,
		CiphertextB64: ciphertextB64,
		TagB64:        authTagB64,
	}
	return &jwe, nil
}
