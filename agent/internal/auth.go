package internal

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Nerzal/gocloak/v11"
	"github.com/reconmap/shared-lib/pkg/logging"

	"github.com/golang-jwt/jwt"
)

var logger = logging.GetLoggerInstance()

func GetPublicKeys() string {
	keycloakHostname, _ := os.LookupEnv("RMAP_KEYCLOAK_HOSTNAME")
	realm := "reconmap"

	client := gocloak.NewClient(keycloakHostname, gocloak.SetAuthAdminRealms("admin/realms"), gocloak.SetAuthRealms("realms"))
	restyClient := client.RestyClient()
	restyClient.SetDebug(true)
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	// this goes to host:port/realms/name
	issuerResponse, err := client.GetIssuer(context.Background(), realm)
	if err != nil {
		logger.Error("error retrieving issue", err)
	}

	return *issuerResponse.PublicKey
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
