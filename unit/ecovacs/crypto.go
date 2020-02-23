package ecovacs

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

const (
	// https://github.com/wpietri/sucks/blob/master/sucks/__init__.py
	CLIENT_KEY = "eJUWrzRv34qFSaYk"
	SECRET     = "Cyu5jcR4zyK6QEPn1hdIGXB5QIDAQABMA0GC"
	PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDb8V0OYUGP3Fs63E1gJzJh+7iq
eymjFUKJUqSD60nhWReZ+Fg3tZvKKqgNcgl7EGXp1yNifJKUNC/SedFG1IJRh5hB
eDMGq0m0RQYDpf9l0umqYURpJ5fmfvH/gjfHe3Eg/NTLm7QEa0a0Il2t3Cyu5jcR
4zyK6QEPn1hdIGXB5QIDAQAB
-----END PUBLIC KEY-----`
)

func DecodePublicKey() (*rsa.PublicKey, error) {
	if block, _ := pem.Decode([]byte(PUBLIC_KEY)); block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	} else if pub, err := x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		return nil, err
	} else {
		return pub.(*rsa.PublicKey), nil
	}
}

func Encrypt(key *rsa.PublicKey, data string) (string, error) {
	if data, err := rsa.EncryptPKCS1v15(rand.Reader, key, []byte(data)); err != nil {
		return "", fmt.Errorf("Encrypt: %w", err)
	} else {
		return base64.StdEncoding.EncodeToString(data), nil
	}
}

func MD5String(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}
