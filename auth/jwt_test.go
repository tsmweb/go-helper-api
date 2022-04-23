package auth

import (
	"fmt"
	"net/http"
	"testing"
)

const (
	_pathPrivateKey = "keys/localhost-key.pem"
	_pathPublicKey  = "keys/localhost-pub.pem"
)

func TestJWT_GenerateToken(t *testing.T) {
	payload := map[string]interface{}{
		"id": 123456,
	}
	jwt := NewJWT(_pathPrivateKey, _pathPublicKey)
	token, err := jwt.GenerateToken(payload, 1)
	if err != nil {
		t.Error(err)
	}

	t.Log(token)
}

func TestJWT_ExtractToken(t *testing.T) {
	payload := map[string]interface{}{
		"id": 123456,
	}
	jwt := NewJWT(_pathPrivateKey, _pathPublicKey)
	token, err := jwt.GenerateToken(payload, 1)
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
	payload := map[string]interface{}{
		"id":    123456,
		"admin": true,
		"dir":   "user.test",
	}
	jwt := NewJWT(_pathPrivateKey, _pathPublicKey)
	token, err := jwt.GenerateToken(payload, 1)
	if err != nil {
		t.Fatal(err)
	}

	req, err := makeRequest(token)
	if err != nil {
		t.Fatal(err)
	}

	id, err := jwt.GetDataToken(req, "id")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)

	admin, err := jwt.GetDataToken(req, "admin")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(admin)

	dir, err := jwt.GetDataToken(req, "dir")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dir)
}

func makeRequest(token string) (*http.Request, error) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	return req, nil
}
