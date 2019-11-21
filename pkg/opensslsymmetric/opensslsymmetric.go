// Based on https://dequeue.blogspot.com/2014/11/decrypting-something-encrypted-with.html
// and https://play.golang.org/p/r3VObSIB4o

package opensslsymmetric

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"fmt"
)

var openSSLSaltHeader string = "Salted_" // OpenSSL salt is always this string + 8 bytes of actual salt

type OpenSSLCreds struct {
	key []byte
	iv  []byte
}

// Decrypt that was encrypted using OpenSSL and AES-256-CBC
func Decrypt(passphrase string, data []byte) ([]byte, error) {
	saltHeader := data[:aes.BlockSize]
	if string(saltHeader[:7]) != openSSLSaltHeader {
		return nil, fmt.Errorf("does not appear to have been encrypted with OpenSSL, salt header missing")
	}
	salt := saltHeader[8:]
	creds, err := extractOpenSSLCreds([]byte(passphrase), salt)
	if err != nil {
		return nil, err
	}
	return decrypt(creds.key, creds.iv, data)
}

func decrypt(key, iv, data []byte) ([]byte, error) {
	if len(data) == 0 || len(data)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("bad blocksize(%v), aes.BlockSize = %v", len(data), aes.BlockSize)
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cbc := cipher.NewCBCDecrypter(c, iv)
	cbc.CryptBlocks(data[aes.BlockSize:], data[aes.BlockSize:])
	out, err := pkcs7Unpad(data[aes.BlockSize:], aes.BlockSize)
	if out == nil {
		return nil, err
	}
	return out, nil
}

// openSSLEvpBytesToKey follows the OpenSSL (undocumented?) convention for extracting the key and IV from passphrase.
// It uses the EVP_BytesToKey() method which is basically:
// D_i = HASH^count(D_(i-1) || password || salt) where || denotes concatentaion, until there are sufficient bytes available
// 48 bytes since we're expecting to handle AES-256, 32bytes for a key and 16bytes for the IV
func extractOpenSSLCreds(password, salt []byte) (OpenSSLCreds, error) {
	m := make([]byte, 48)
	prev := []byte{}
	for i := 0; i < 3; i++ {
		prev = hash(prev, password, salt)
		copy(m[i*16:], prev)
	}
	return OpenSSLCreds{key: m[:32], iv: m[32:]}, nil
}

func hash(prev, password, salt []byte) []byte {
	a := make([]byte, len(prev)+len(password)+len(salt))
	copy(a, prev)
	copy(a[len(prev):], password)
	copy(a[len(prev)+len(password):], salt)
	return md5sum(a)
}

func md5sum(data []byte) []byte {
	h := md5.New()
	h.Write(data)
	return h.Sum(nil)
}

// pkcs7Unpad returns slice of the original data without padding.
func pkcs7Unpad(data []byte, blocklen int) ([]byte, error) {
	if blocklen <= 0 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}
	if len(data)%blocklen != 0 || len(data) == 0 {
		return nil, fmt.Errorf("invalid data len %d", len(data))
	}
	padlen := int(data[len(data)-1])
	if padlen > blocklen || padlen == 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	pad := data[len(data)-padlen:]
	for i := 0; i < padlen; i++ {
		if pad[i] != byte(padlen) {
			return nil, fmt.Errorf("invalid padding")
		}
	}
	return data[:len(data)-padlen], nil
}
