package auth

import (
	"net/http"
	"testing"
)

const (
	_pathPrivateKey = "keys/private-key"
	_pathPublicKey  = "keys/public-key.pub"
)

func TestJWT_GenerateToken(t *testing.T) {
	jwt := NewJWT(_pathPrivateKey, _pathPublicKey)
	token, err := jwt.GenerateToken("+5511999999999", 1)
	if err != nil {
		t.Error(err)
	}

	t.Log(token)
}

func TestJWT_ExtractToken(t *testing.T) {
	jwt := NewJWT(_pathPrivateKey, _pathPublicKey)
	token, err := jwt.GenerateToken("+5511999999999", 1)
	if err != nil {
		t.Fatal(err)
	}

	req, err := makeRequest(token)
	if err != nil {
		t.Fatal(err)
	}

	extractToken, err := jwt.ExtractToken(req)
	if err != nil {
		t.Error(err)
	}

	if extractToken == token {
		t.Log(extractToken)
	} else {
		t.Log("Invalid token")
	}
}

func TestJWT_GetDataToken(t *testing.T) {
	jwt := NewJWT(_pathPrivateKey, _pathPublicKey)
	token, err := jwt.GenerateToken("+5511999999999", 1)
	if err != nil {
		t.Fatal(err)
	}

	req, err := makeRequest(token)
	if err != nil {
		t.Fatal(err)
	}

	d, err := jwt.GetDataToken(req, "id")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(d)
}

func makeRequest(token string) (*http.Request, error) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token)

	return req, nil
}
