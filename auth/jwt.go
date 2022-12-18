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
	_privateKey, err := privateKey(pathPrivateKey)
	if err != nil {
		panic(err)
	}

	_publicKey, err := publicKey(pathPublicKey)
	if err != nil {
		panic(err)
	}

	return &_jwt{
		privateKey: _privateKey,
		publicKey:  _publicKey,
	}
}

// GenerateToken generates a JWT token and returns the token in string format.
func (j *_jwt) GenerateToken(payload map[string]interface{}, exp int) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	mapClaims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * time.Duration(exp)).Unix(),
		"iat": time.Now().Unix(),
	}

	if len(payload) == 0 {
		return "", fmt.Errorf("invalid payload")
	}
	for k, v := range payload {
		mapClaims[k] = v
	}

	token.Claims = mapClaims
	tokenString, err := token.SignedString(j.privateKey)
	if err != nil {
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

// GetDataToken receives a key per parameter and extracts the corresponding token value.
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
	// token, err := request.ParseFromRequest(r, request.OAuth2Extractor, j.parseKeyFunc())
	extractor := &request.MultiExtractor{
		request.OAuth2Extractor,
		request.ArgumentExtractor{"authorization"},
	}
	token, err := request.ParseFromRequest(r, extractor, j.parseKeyFunc())
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func (j *_jwt) parseKeyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return j.publicKey, nil
	}
}

func privateKey(pathPrivateKey string) (*rsa.PrivateKey, error) {
	privateKeyFile, err := os.Open(pathPrivateKey)
	if err != nil {
		return nil, err
	}

	pemFileInfo, err := privateKeyFile.Stat()
	if err != nil {
		return nil, err
	}

	size := pemFileInfo.Size()
	pemBytes := make([]byte, size)
	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pemBytes)
	if err != nil {
		return nil, err
	}

	data, _ := pem.Decode(pemBytes)
	if err = privateKeyFile.Close(); err != nil {
		return nil, err
	}

	privatePKCS1Key, errPKCS1 := x509.ParsePKCS1PrivateKey(data.Bytes)
	if errPKCS1 == nil {
		return privatePKCS1Key, nil
	}

	privatePKCS8Key, errPKCS8 := x509.ParsePKCS8PrivateKey(data.Bytes)
	if errPKCS8 == nil {
		privatePKCS8RsaKey, ok := privatePKCS8Key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("PKCS8 contained non-RSA key. Expected RSA key")
		}
		return privatePKCS8RsaKey, nil
	}

	return nil, fmt.Errorf("failed to parse private key as PKCS#1 or PKCS#8. (%s). (%s)",
		errPKCS1, errPKCS8)
}

func publicKey(pathPublicKey string) (*rsa.PublicKey, error) {
	publicKeyFile, err := os.Open(pathPublicKey)
	if err != nil {
		return nil, err
	}

	pemFileInfo, err := publicKeyFile.Stat()
	if err != nil {
		return nil, err
	}

	size := pemFileInfo.Size()
	pemBytes := make([]byte, size)
	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pemBytes)
	if err != nil {
		return nil, err
	}

	data, _ := pem.Decode(pemBytes)
	if err = publicKeyFile.Close(); err != nil {
		return nil, err
	}

	publicPKIXKey, errPKIX := x509.ParsePKIXPublicKey(data.Bytes)
	if errPKIX == nil {
		rsaPub, ok := publicPKIXKey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("PKIX contained non-RSA key. Expected RSA key")
		}
		return rsaPub, nil
	}

	publicPKCS1Key, errPKCS1 := x509.ParsePKCS1PublicKey(data.Bytes)
	if errPKCS1 == nil {
		return publicPKCS1Key, nil
	}

	return nil, fmt.Errorf("failed to parse public key as PKIX or PKCS#1. (%s). (%s)",
		errPKIX, errPKCS1)
}
