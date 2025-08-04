package internal

import (
	"encoding/json"
	"net/http"

	"github.com/reconmap/shared-lib/pkg/api"
	"go.uber.org/zap"
)

func handleSystemInfo(w http.ResponseWriter, r *http.Request) {
	logger.Debug("handling notification request")

	err := CheckRequestToken(r)
	if err != nil {
		logger.Error("cannot check request token", zap.Error(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	systemInfo := api.SystemInfo{
		Hostname: "localhost", // Replace with actual hostname retrieval logic
		Os:       "linux",     // Replace with actual OS retrieval logic
		Arch:     "amd64",     // Replace with actual architecture retrieval logic
		CPU:      "4 cores",   // Replace with actual CPU retrieval logic
		Memory:   "8GB",       // Replace with actual memory retrieval logic
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(systemInfo)
}
