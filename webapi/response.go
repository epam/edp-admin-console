package webapi

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func OKJsonResponse(ctx context.Context, w http.ResponseWriter, data interface{}) {
	logger := LoggerFromContext(ctx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("encode json response failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		if _, wErr := w.Write([]byte("build response failed")); wErr != nil {
			logger.Error("write error response failed", zap.Error(wErr))
		}
	}
	return
}

func ErrNotFoundResponse(ctx context.Context, w http.ResponseWriter, msg string) {
	logger := LoggerFromContext(ctx)
	w.WriteHeader(http.StatusNotFound)
	_, wErr := w.Write([]byte(msg))
	if wErr != nil {
		logger.Error("write error response failed", zap.Error(wErr))
	}
}
