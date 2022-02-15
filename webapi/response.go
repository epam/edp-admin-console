package webapi

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
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

func BadRequestResponse(ctx context.Context, w http.ResponseWriter, msg string) {
	logger := LoggerFromContext(ctx)
	w.WriteHeader(http.StatusBadRequest)
	_, wErr := w.Write([]byte(msg))
	if wErr != nil {
		logger.Error("write error response failed", zap.Error(wErr))
	}
}

func InternalErrorResponse(ctx context.Context, w http.ResponseWriter, msg string) {
	logger := LoggerFromContext(ctx)
	w.WriteHeader(http.StatusInternalServerError)
	_, wErr := w.Write([]byte(msg))
	if wErr != nil {
		logger.Error("write error response failed", zap.Error(wErr))
	}
}

func OkHTMLResponse(ctx context.Context, writer http.ResponseWriter, data interface{}, tmpl *template.Template, mainFileName string, templatePaths ...string) {
	logger := LoggerFromContext(ctx)
	tmplInternal, err := tmpl.ParseFiles(templatePaths...)
	if err != nil {
		logger.Error("cant find template page", zap.Error(err), zap.Strings("template_paths", templatePaths))
		InternalErrorResponse(ctx, writer, "cant find template page")
		return
	}
	var buf bytes.Buffer
	err = tmplInternal.ExecuteTemplate(&buf, mainFileName, data)
	if err != nil {
		logger.Error("cant exec page", zap.Error(err), zap.Strings("template_paths", templatePaths), zap.String("main_template", mainFileName))
		InternalErrorResponse(ctx, writer, "cant exec page")
		return
	}
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(buf.Bytes())
	if err != nil {
		logger.Error("write error response failed", zap.Error(err))
	}
}
