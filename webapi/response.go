package webapi

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"edp-admin-console/internal/applog"
)

func OKJsonResponse(ctx context.Context, w http.ResponseWriter, data interface{}) {
	logger := applog.LoggerFromContext(ctx)
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

func NotFoundResponse(ctx context.Context, w http.ResponseWriter, msg string) {
	logger := applog.LoggerFromContext(ctx)
	w.WriteHeader(http.StatusNotFound)
	_, wErr := w.Write([]byte(msg))
	if wErr != nil {
		logger.Error("write error response failed", zap.Error(wErr))
	}
}

func BadRequestResponse(ctx context.Context, w http.ResponseWriter, msg string) {
	logger := applog.LoggerFromContext(ctx)
	w.WriteHeader(http.StatusBadRequest)
	_, wErr := w.Write([]byte(msg))
	if wErr != nil {
		logger.Error("write error response failed", zap.Error(wErr))
	}
}

func InternalErrorResponse(ctx context.Context, w http.ResponseWriter, msg string) {
	logger := applog.LoggerFromContext(ctx)
	w.WriteHeader(http.StatusInternalServerError)
	_, wErr := w.Write([]byte(msg))
	if wErr != nil {
		logger.Error("write error response failed", zap.Error(wErr))
	}
}
