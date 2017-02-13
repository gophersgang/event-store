package main

import (
	"fmt"
	"time"
	"golang.org/x/oauth2/google"
	"crypto/rsa"
	"encoding/pem"
	"golang.org/x/oauth2/jws"
	"crypto/x509"
	"net/url"
	"net/http"
	"log"
	"errors"
	"encoding/json"
)

type GoogleIdTokenResponse struct {
	IDToken string `json:"id_token"`
}

func GenerateJWT(scope string) {
	conf, err := google.JWTConfigFromJSON([]byte(VENDASTA_LOCAL_JSON_KEY), "")
	if err != nil {
		log.Fatalf("Could not parse service account JSON: %v", err)
	}
	rsaKey, err := parseKey(conf.PrivateKey)
	if err != nil {
		log.Fatalf("Could not get RSA key: %v", err)
	}

	iat := time.Now()
	exp := iat.Add(time.Hour)

	jwt := &jws.ClaimSet{
		Iss:   "vendasta-local@repcore-prod.iam.gserviceaccount.com",
		Aud:   "https://www.googleapis.com/oauth2/v4/token",
		Scope: scope,
		Iat:   iat.Unix(),
		Exp:   exp.Unix(),
	}
	jwsHeader := &jws.Header{
		Algorithm: "RS256",
		Typ:       "JWT",
	}

	msg, err := jws.Encode(jwsHeader, jwt, rsaKey)
	if err != nil {
		log.Fatalf("Could not encode JWT: %v", err)
	}
	resp, err := http.PostForm("https://www.googleapis.com/oauth2/v4/token", url.Values{"grant_type": {"urn:ietf:params:oauth:grant-type:jwt-bearer"}, "assertion": {string(msg[:])}})
	defer resp.Body.Close()
	var idToken GoogleIdTokenResponse
	json.NewDecoder(resp.Body).Decode(&idToken)
	fmt.Println(idToken.IDToken)

}

func parseKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block != nil {
		key = block.Bytes
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		parsedKey, err = x509.ParsePKCS1PrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("private key should be a PEM or plain PKSC1 or PKCS8; parse error: %v", err)
		}
	}
	parsed, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is invalid")
	}
	return parsed, nil
}