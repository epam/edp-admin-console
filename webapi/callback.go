package webapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (h *HandlerAuth) CallBack(writer http.ResponseWriter, request *http.Request) {
	log := LoggerFromContext(request.Context())
	log.Info("Start callback flow...")
	ctx := request.Context()
	authSessionID, ok := GetCookieByName(request, AuthSessionIDName)
	if !ok {
		InternalErrorResponse(ctx, writer, "cant find auth session id")
		return
	}

	queryState := request.URL.Query().Get("state")
	log.Info("State has been retrieved from query param", zap.String("queryState", queryState))
	sessionState, ok := h.StateMap[authSessionID]
	if !ok {
		InternalErrorResponse(ctx, writer, "cant find session state")
		return
	}
	delete(h.StateMap, authSessionID)
	log.Info("State has been retrieved from the session", zap.String("state_key", authSessionID), zap.Any("sessionState", sessionState))
	if queryState != sessionState {
		log.Info("State does not match", zap.String("query state", queryState), zap.String("session state", sessionState))
		BadRequestResponse(ctx, writer, "State does not match")
		return
	}

	authCode := request.URL.Query().Get("code")
	log.Info("Authorization code has been retrieved from query param")
	authConfig := h.AuthController
	token, err := authConfig.Oauth2Service.Exchange(ctx, authCode)
	if err != nil {
		log.Info("Failed to exchange token with code", zap.String("code", authCode))
		InternalErrorResponse(ctx, writer, "Failed to exchange token with code")
		return
	}

	log.Info("Authorization code has been successfully exchanged with token")

	path := h.getRedirectPath(authSessionID)

	ts := authConfig.Oauth2Service.TokenSource(ctx, token)
	sessionID := uuid.New().String()
	h.TokenMap[sessionID] = ts
	log.Info("Token source has been saved to the session")

	cookie := &http.Cookie{
		Name:    SessionIDName,
		Value:   sessionID,
		Expires: time.Now().Add(SessionExpirationTime),
		Path:    "/",
	}
	http.SetCookie(writer, cookie)
	http.Redirect(writer, request, path, http.StatusFound)
}

func (h *HandlerAuth) getRedirectPath(sessionID string) string {
	requestPath := h.UrlMap[sessionID]
	if requestPath == "" {
		return fmt.Sprintf("%s/admin/edp/overview", h.BasePath)
	}
	delete(h.UrlMap, sessionID)
	return requestPath
}
