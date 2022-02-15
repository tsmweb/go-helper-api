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

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
)

// JWT uses JSON Web Token for generate and extract token.
type JWT interface {
	GenerateToken(payload map[string]interface{}, exp int) (string, error)
	ExtractToken(r *http.Request) (string, error)
	GetDataToken(r *http.Request, key string) (interface{}, error)
}

type _jwt struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewJWT create new instance of JWT.
func NewJWT(pathPrivateKey, pathPublicKey string) JWT {
	return &_jwt{
		privateKey: privateKey(pathPrivateKey),
		publicKey:  publicKey(pathPublicKey),
	}
}

// GenerateToken generates a JWT token and returns the token in string format.
func (j *_jwt) GenerateToken(payload map[string]interface{}, exp int) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	mapClaims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * time.Duration(exp)).Unix(),
		"iat": time.Now().Unix(),
	}

	if payload == nil || len(payload) == 0 {
		return "", fmt.Errorf("invalid payload")
	}
	for k, v := range payload {
		mapClaims[k] = v
	}

	token.Claims = mapClaims
	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
		panic(err)
		return "", err
	}

	return tokenString, nil
}

// ExtractToken extracts the token from an HTTP Request.
func (j *_jwt) ExtractToken(r *http.Request) (string, error) {
	t, err := j.token(r)
	if err != nil {
		return "", err
	}

	if !t.Valid {
		return "", fmt.Errorf("invalid token")
	}

	return t.Raw, nil
}

// MapClaims receives a key per parameter and extracts the corresponding token value.
func (j *_jwt) GetDataToken(r *http.Request, key string) (interface{}, error) {
	token, err := j.token(r)
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	if claims != nil {
		return claims[key], nil
	}

	return nil, fmt.Errorf("key not found")
}

// token extracts the token from an HTTP Request.
func (j *_jwt) token(r *http.Request) (*jwt.Token, error) {
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, j.parseKeyfunc())
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func (j *_jwt) parseKeyfunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return j.publicKey, nil
	}
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
