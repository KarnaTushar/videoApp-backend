package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/KarnaTushar/videoApp-backend/pkg/models"
	"github.com/KarnaTushar/videoApp-backend/utils"

	"github.com/spf13/viper"
)

type contextKey struct {
	name string
}

var userContextKey = &contextKey{"user"}

// AuthHandler is a middleware for authentication
func AuthHandler(db *models.Database, logger *utils.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			if !viper.GetBool("ENABLE_OAUTH") {
				next.ServeHTTP(w, r)
				return
			}

			header := r.Header.Get("Authorization")

			if header == "" {
				logger.Debug().Msg("No Token Provided")
			} else {
				splitToken := strings.Split(header, "Bearer ")
				token := splitToken[1]

				var tokenData models.Token
				var user models.User

				// Fetch the token
				err := db.Get(&tokenData, "SELECT token_id, user_id FROM tokens WHERE token_id= ?", token)
				if err != nil {
					logger.Debug().Str("token", token).Msg("Passed Invalid token")
					next.ServeHTTP(w, r)
					return
				}

				err = db.Get(&user, "SELECT id, identifier, user_name, email FROM users WHERE id= ?", tokenData.UserID)
				if err != nil {
					logger.Error().Int64("id", tokenData.UserID).Str("token", token).Msg("User does not exist for the provided token")
					next.ServeHTTP(w, r)
					return
				}

				logger.Info().Str("token", token).Interface("user", user).Msg("Successfull")
				ctx := context.WithValue(r.Context(), userContextKey, &user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// TODO: See if there's a better way to structure this code
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext fetches the user from the context
func GetUserFromContext(ctx context.Context) (*models.UserAccount, error) {
	userObject := ctx.Value(userContextKey)
	if userObject != nil {
		return userObject.(*models.UserAccount), nil
	}

	return nil, errors.New("No such user")
}
