package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"log"

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

	//Make the cipher text a byte array of size BlockSize + the length of the message
	cipherText := make([]byte, aes.BlockSize+len(plainText))

	//iv is the ciphertext up to the blocksize (16)
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	//Encrypt the data:
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	//Return string encoded in base64
	return base64.RawStdEncoding.EncodeToString(cipherText), err
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
	log.Println(cipherKey)

	//Remove base64 encoding:
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)
	if err != nil {
		return "", err
	}

	//Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return "", err
	}

	//IF the length of the cipherText is less than 16 Bytes:
	if len(cipherText) < aes.BlockSize {
		err = errors.New("Ciphertext block size is too short!")
		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	//Decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), err
}
