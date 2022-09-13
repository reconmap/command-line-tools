package internal

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

var rsakeys map[string]*rsa.PublicKey

func GetPublicKeys() string {
	rsakeys = make(map[string]*rsa.PublicKey)
	var body map[string]string
	keycloakHostname, _ := os.LookupEnv("RMAP_KEYCLOAK_HOSTNAME")
	uri := keycloakHostname + "/realms/reconmap"
	resp, _ := http.Get(uri)
	json.NewDecoder(resp.Body).Decode(&body)
	return body["public_key"]
}

func CheckRequestToken(r *http.Request) error {
	params := r.URL.Query()

	if !params.Has("token") {
		return errors.New("missing \"token\" parameter")
	} else {
		tokenParam := params.Get("token")
		pubkey := "-----BEGIN PUBLIC KEY-----\n" + GetPublicKeys() + "\n-----END PUBLIC KEY-----"
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubkey))
		if err != nil {
			fmt.Errorf("validate: parse key: %w", err)
		}

		token, err := jwt.Parse(tokenParam, func(jwtToken *jwt.Token) (interface{}, error) {
			if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
			}

			return key, nil
		})
		if !token.Valid {
			return err
		}

		if _, ok := token.Claims.(jwt.MapClaims); !ok {
			return errors.New("Unable to parse claims")
		}
	}

	return nil
}
