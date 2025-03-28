package internal

import (
	"net/http"

	"go.uber.org/zap"
)

func handleNotifications(w http.ResponseWriter, r *http.Request) {
	logger.Debug("handling notification request")

	err := CheckRequestToken(r)
	if err != nil {
		logger.Error("cannot check request token", zap.Error(err))
		return
	}

	conn, err := UpgradeRequest(w, r)
	if err != nil {
		logger.Error("cannot upgrade request", zap.Error(err))
		return
	}

	registerClient(conn)
}
