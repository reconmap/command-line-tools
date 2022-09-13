package internal

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func handleNotifications(w http.ResponseWriter, r *http.Request) {
	log.Debug("handling notification request")

	err := CheckRequestToken(r)
	if err != nil {
		log.Error(err)
		return
	}

	conn, err := UpgradeRequest(w, r)
	if err != nil {
		log.Error(err)
		return
	}

	registerClient(conn)
}
