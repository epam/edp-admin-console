package webapi

import (
	"context"
	"net/http"
	"path"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"edp-admin-console/internal/applog"
	"edp-admin-console/internal/config"
)

const (
	SessionIDName     = "sessionID"
	AuthSessionIDName = "authSessionID"
)

func WithLoggerMw(logger *zap.Logger) func(next http.Handler) http.Handler {
	mw := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			r = r.WithContext(applog.ContextWithLogger(ctx, logger))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return mw
}

func WithAuthN(tokenMap map[string]oauth2.TokenSource, urlMap map[string]string, stateMap map[string]string, oauth2config *config.AuthController) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			log := applog.LoggerFromContext(r.Context())
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
	log := applog.LoggerFromContext(r.Context())
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
	log := applog.LoggerFromContext(r.Context())
	cookie, err := r.Cookie(name)
	if err != nil {
		log.Info("cant find cookie by name")
		return "", false
	}
	return cookie.Value, true
}

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type TokenClaim struct {
	Name        string      `json:"name"`
	RealmAccess RealmAccess `json:"realm_access"`
}

type ConfigRoles struct {
	DevRole   string
	AdminRole string
}

type AuthorisedUser struct {
	TokenClaim  TokenClaim
	ConfigRoles ConfigRoles
}

func NewAuthorisedUser(claim TokenClaim, configRoles ConfigRoles) AuthorisedUser {
	return AuthorisedUser{
		TokenClaim:  claim,
		ConfigRoles: configRoles,
	}
}

func GuestUser() AuthorisedUser {
	return AuthorisedUser{
		TokenClaim: TokenClaim{
			Name: "Guest",
			RealmAccess: RealmAccess{
				Roles: nil,
			},
		},
	}
}

func (user *AuthorisedUser) UserName() string {
	return user.TokenClaim.Name
}

func (user *AuthorisedUser) IsAdmin() bool {
	for _, role := range user.TokenClaim.RealmAccess.Roles {
		if role == user.ConfigRoles.AdminRole {
			return true
		}
	}
	return false
}

func (user *AuthorisedUser) IsDeveloper() bool {
	for _, role := range user.TokenClaim.RealmAccess.Roles {
		if role == user.ConfigRoles.AdminRole {
			return true
		}
	}
	return false
}

type AuthorisedUserKey struct{}

func WithAuthZ(tokenMap map[string]oauth2.TokenSource, authController *config.AuthController) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := applog.LoggerFromContext(ctx)
			sessionID, ok := GetCookieByName(r, SessionIDName)
			if !ok {
				log.Error("cant find session id in cookie")
				InternalErrorResponse(ctx, w, "cant find session id in cookie")
				return
			}

			tc, ok := tokenMap[sessionID]
			if !ok {
				log.Error("cant find token source", zap.String("key", "token_source"))
				InternalErrorResponse(ctx, w, "cant find token source")
				return
			}

			token, err := tc.Token()
			if err != nil {
				log.Error("cant convert token source to token")
				InternalErrorResponse(ctx, w, "cant convert token source")
				return
			}

			tokenID, err := authController.Verifier.Verify(context.Background(), token.AccessToken)
			if err != nil {
				log.Error("cant verify token")
				InternalErrorResponse(ctx, w, "cant verify token")
				return
			}
			var claim TokenClaim
			err = tokenID.Claims(&claim)
			if err != nil {
				log.Error("Error has been occurred during the parsing token", zap.Any("token", token))
				InternalErrorResponse(ctx, w, "cant unmarshall token claims")
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), AuthorisedUserKey{},
				NewAuthorisedUser(claim, ConfigRoles{AdminRole: authController.AdminRoleName, DevRole: authController.DevRoleName})))

			next.ServeHTTP(w, r)
		})
	}
}

func WithLogRequestBoundaries() func(next http.Handler) http.Handler {
	handler := func(next http.Handler) http.Handler {
		mw := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := applog.LoggerFromContext(ctx)
			requestURI := r.RequestURI
			logger.Info("REQUEST_STARTED", zap.String("REQUEST_URI", requestURI))
			next.ServeHTTP(w, r)
			logger.Info("REQUEST_COMPLETED", zap.String("REQUEST_URI", requestURI))
		}
		return http.HandlerFunc(mw)
	}
	return handler
}
