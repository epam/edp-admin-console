package webapi

import (
	"net/http"
	"path"
	"time"

	"edp-admin-console/internal/config"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

const (
	SessionIDName         = "sessionID"
	AuthSessionIDName     = "authSessionID"
	SessionExpirationTime = 24 * time.Hour
)

func WithLoggerMw(logger *zap.Logger) func(next http.Handler) http.Handler {
	mw := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			r = r.WithContext(ContextWithLogger(ctx, logger))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return mw
}

func WithAuth(tokenMap map[string]oauth2.TokenSource, urlMap map[string]string, stateMap map[string]string, oauth2config *config.AuthController) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			log := LoggerFromContext(r.Context())
			sessionID, ok := GetCookieByName(r, SessionIDName)
			if !ok {
				startAuth(w, r, stateMap, urlMap, oauth2config)
				return
			}
			tc, ok := tokenMap[sessionID]
			if !ok {
				log.Info("cant find token source", zap.String("key", "token_source"))
				startAuth(w, r, stateMap, urlMap, oauth2config)
				return
			}

			token, err := tc.Token()
			if err != nil {
				delete(tokenMap, sessionID)
				log.Error("cant convert token source to token", zap.Any("token source", tc), zap.Error(err))
				startAuth(w, r, stateMap, urlMap, oauth2config)
				return
			}
			if !token.Valid() {
				log.Info("failed to validate token", zap.Any("token", token))
				startAuth(w, r, stateMap, urlMap, oauth2config)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func startAuth(w http.ResponseWriter, r *http.Request, stateMap map[string]string, urlMap map[string]string, oauth2config *config.AuthController) {
	log := LoggerFromContext(r.Context())
	state := uuid.New().String()
	authSessionID := uuid.New().String()

	log.Info("start new auth session")
	cookie := &http.Cookie{
		Name:    AuthSessionIDName,
		Value:   authSessionID,
		Expires: time.Now().Add(oauth2config.AuthSessionTTL),
		Path:    "/",
	}

	stateMap[authSessionID] = state
	log.Info("State has been generated, saved in the session and added in the auth request",
		zap.String("key", authSessionID), zap.String("state", state))
	if r.Method == http.MethodGet {
		urlMap[authSessionID] = path.Join(r.URL.Host, r.URL.Path)
	}
	redirectURL := oauth2config.Oauth2Service.AuthCodeURL(state)
	log.Info("redirect url", zap.String("url", redirectURL))
	http.SetCookie(w, cookie)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func GetCookieByName(r *http.Request, name string) (string, bool) {
	log := LoggerFromContext(r.Context())
	cookie, err := r.Cookie(name)
	if err != nil {
		log.Info("cant find cookie by name")
		return "", false
	}
	return cookie.Value, true
}
