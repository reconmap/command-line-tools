package internal

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Nerzal/gocloak/v13"
	"go.uber.org/zap"

	"github.com/golang-jwt/jwt"
)

const realm = "reconmap"

func NewGocloakClient() *gocloak.GoCloak {
	keycloakHostname, _ := os.LookupEnv("RMAP_KEYCLOAK_HOSTNAME")
	keycloakDebug, _ := os.LookupEnv("RMAP_KEYCLOAK_DEBUG")
	keycloakSkipVerify, _ := os.LookupEnv("RMAP_KEYCLOAK_SKIP_TLS_VERIFY")

	client := gocloak.NewClient(keycloakHostname, gocloak.SetAuthAdminRealms("admin/realms"), gocloak.SetAuthRealms("realms"))

	restyClient := client.RestyClient()
	restyClient.SetDebug(keycloakDebug == "true")
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: keycloakSkipVerify == "true"})

	return client
}

func GetAccessToken(app *App) (string, error) {
	clientID, _ := os.LookupEnv("RMAP_AGENT_CLIENT_ID")
	clientSecret, _ := os.LookupEnv("RMAP_AGENT_CLIENT_SECRET")

	client := NewGocloakClient()

	ctx := context.Background()
	token, err := client.LoginClient(ctx, clientID, clientSecret, realm)
	if err != nil {
		return "", err
	}

	tokenInfo, err := client.RetrospectToken(ctx, token.AccessToken, clientID, clientSecret, realm)
	if err != nil {
		app.Logger.Error("unable to inspect token", zap.Error(err))
		panic(err)
	}

	if !*tokenInfo.Active {
		app.Logger.Error("token is not active")
		panic("token is not active")
	}

	return token.AccessToken, nil
}

func GetPublicKeys() string {
	client := NewGocloakClient()

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
