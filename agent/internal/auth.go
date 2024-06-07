package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/reconmap/shared-lib/pkg/logging"

	"github.com/golang-jwt/jwt"
)

var logger = logging.GetLoggerInstance()

func GetPublicKeys() string {
	var body map[string]string
	keycloakHostname, _ := os.LookupEnv("RMAP_KEYCLOAK_HOSTNAME")
	uri := keycloakHostname + "/realms/reconmap"
	resp, _ := http.Get(uri)
	err := json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		logger.Error(err)
	}

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
			err := fmt.Errorf("validate: parse key: %w", err)
			return err
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
			return errors.New("unable to parse claims")
		}
	}

	return nil
}
