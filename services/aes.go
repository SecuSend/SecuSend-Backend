package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

func generateAESKeyFromSecret(secret string) ([]byte, error) {
	salt := []byte("my-secret-salt") // Use a byte slice for the salt
	key := pbkdf2.Key([]byte(secret), salt, 10000, 32, sha256.New)
	return key, nil
}

/*
 *	FUNCTION		: encrypt
 *	DESCRIPTION		:
 *		This function takes a string and a cipher key and uses AES to encrypt the message
 *
 *	PARAMETERS		:
 *		string secret	: String secret used to encrypt
 *		string message	: String containing the message to encrypt
 *
 *	RETURNS			:
 *		string encoded	: String containing the encoded user input
 *		error err	: Error message
 */
func Encrypt(secret string, message string) (encoded string, err error) {
	// Use PBKDF2 to derive a 256-bit key from the secret
	cipherKey, err := generateAESKeyFromSecret(secret)
	if err != nil {
		return "", err
	}

	//Create byte array from the input string
	plainText := []byte(message)

	//Create a new AES cipher using the key
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	//Create a new GCM cipher
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//Generate a random nonce
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	//Encrypt the data:
	cipherText := aesgcm.Seal(nil, nonce, plainText, nil)

	//Return string encoded in base64
	return base64.RawStdEncoding.EncodeToString(append(nonce, cipherText...)), err
}

/*
 *	FUNCTION		: decrypt
 *	DESCRIPTION		:
 *		This function takes a string and a key and uses AES to decrypt the string into plain text
 *
 *	PARAMETERS		:
 *		string secret	: String secret used to encrypt
 *		string secure	: String containing an encrypted message
 *
 *	RETURNS			:
 *		string decoded	: String containing the decrypted equivalent of secure
 *		error err	: Error message
 */
func Decrypt(secret string, secure string) (decoded string, err error) {
	// Use PBKDF2 to derive a 256-bit key from the secret
	cipherKey, err := generateAESKeyFromSecret(secret)
	if err != nil {
		return "", err
	}

	//Remove base64 encoding:
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)
	if err != nil {
		return "", err
	}

	//Create a new AES-GCM cipher with the key and encrypted message
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < gcm.NonceSize() {
		return "", errors.New("Ciphertext block size is too short!")
	}

	nonce := cipherText[:gcm.NonceSize()]
	cipherText = cipherText[gcm.NonceSize():]

	//Decrypt the message
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), err
}
