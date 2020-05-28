/*
Package auth provides authentication and authorization methods using JSON Web Tokens.

 */
package auth

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

// JWT uses JSON Web Token for generate and extract token.
type JWT struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewJWT create new instance of JWT.
func NewJWT(pathPrivateKey, pathPublicKey string) *JWT {
	return &JWT{
		privateKey: privateKey(pathPrivateKey),
		publicKey:  publicKey(pathPublicKey),
	}
}

// GenerateToken generates a JWT token and returns the token in string format.
func (a *JWT) GenerateToken(id string, exp int) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims = jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * time.Duration(exp)).Unix(),
		"iat": time.Now().Unix(),
		"sub": id,
	}
	tokenString, err := token.SignedString(a.privateKey)
	if err != nil {
		panic(err)
		return "", err
	}

	return tokenString, nil
}

// Token extracts the token from an HTTP Request.
func (a *JWT) Token(r *http.Request) (*jwt.Token, error) {
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return a.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token invalid")
	}

	return token, nil
}

// MapClaims receives a key per parameter and extracts the corresponding token value.
func (a *JWT) MapClaims(r *http.Request, key string) (interface{}, error) {
	token, err := a.Token(r)
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	if claims != nil {
		return claims[key], nil
	}

	return nil, fmt.Errorf("Key not found")
}

func privateKey(pathPrivateKey string) *rsa.PrivateKey {
	privateKeyFile, err := os.Open(pathPrivateKey)
	if err != nil {
		panic(err)
	}

	pemFileInfo, _ := privateKeyFile.Stat()
	size := pemFileInfo.Size()
	pemBytes := make([]byte, size)
	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pemBytes)
	data, _ := pem.Decode(pemBytes)
	privateKeyFile.Close()

	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		panic(err)
	}

	return privateKeyImported
}

func publicKey(pathPublicKey string) *rsa.PublicKey {
	publicKeyFile, err := os.Open(pathPublicKey)
	if err != nil {
		panic(err)
	}

	pemFileInfo, _ := publicKeyFile.Stat()
	size := pemFileInfo.Size()
	pemBytes := make([]byte, size)
	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pemBytes)
	data, _ := pem.Decode(pemBytes)
	publicKeyFile.Close()

	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)
	if err != nil {
		panic(err)
	}

	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)
	if !ok {
		panic(err)
	}

	return rsaPub
}
