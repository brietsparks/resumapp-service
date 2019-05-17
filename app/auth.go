package app

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/auth0-community/go-auth0"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"io/ioutil"
	"net/http"
)

func loadPublicKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	// Try to load SubjectPublicKeyInfo
	pub, err0 := x509.ParsePKIXPublicKey(input)
	if err0 == nil {
		return pub, nil
	}

	cert, err1 := x509.ParseCertificate(input)
	if err1 == nil {
		return cert.PublicKey, nil
	}

	return nil, fmt.Errorf("square/go-jose: parse error, got '%s' and '%s'", err0, err1)
}

type ValidateToken func (r *http.Request) (*jwt.JSONWebToken, *jwt.Claims, error)

func NewValidateToken(auth0CertPath string, auth0Audience string, issuer string) ValidateToken  {
	spew.Dump(auth0CertPath, auth0Audience, issuer)
	p, err := ioutil.ReadFile(auth0CertPath)
	if err != nil {
		panic(err)
	}

	secret, err := loadPublicKey(p)
	if err != nil {
		panic(err)
	}

	secretProvider := auth0.NewKeyProvider(secret)

	configuration := auth0.NewConfiguration(
		secretProvider,
		[]string{auth0Audience},
		issuer,
		jose.RS256,
	)
	validator := auth0.NewValidator(configuration, nil)

	return func(r *http.Request) (*jwt.JSONWebToken, *jwt.Claims, error) {
		token, err := validator.ValidateRequest(r)
		if err != nil {
			return nil, nil, err
		}

		claims := &jwt.Claims{}
		err = validator.Claims(r, token, claims)
		if err != nil {
			return token, nil, err
		}

		return token, claims, nil
	}
}
