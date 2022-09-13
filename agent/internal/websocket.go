package internal

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// UpgradeRequest converts http connection to a websocket one.
func UpgradeRequest(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithError(err).Error("Unable to upgrade connection")
		return nil, err
	}
	return conn, nil
}
