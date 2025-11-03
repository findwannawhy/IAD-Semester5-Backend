package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const prefix = "Bearer"
const guestSessionCookie = "guest_session_id"
const guestSessionTTL = 20 * time.Minute

func (h *Handler) ModeratorMiddleware(allowedRole bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractTokenFromHeader(c.Request)

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, nil
			}
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		blacklisted, err := h.Repository.IsTokenBlacklisted(context.Background(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error"})
			return
		}
		if blacklisted {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		jwtIsModerator, ok := claims["is_moderator"].(bool)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}

		if allowedRole && !jwtIsModerator {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "forbidden"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func (h *Handler) WithOptionalAuthCheck() func(ctx *gin.Context) {
	return func(c *gin.Context) {
			tokenString := c.GetHeader("Authorization")
			prefix := "Bearer "

			if tokenString == "" || !strings.HasPrefix(tokenString, prefix) {
					c.Set("user_id", "") 
					c.Next()
					return
			}

			tokenString = strings.TrimPrefix(tokenString, prefix)

	blacklisted, err := h.Repository.IsTokenBlacklisted(context.Background(), tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "error"})
		return
	}
	if blacklisted {
		c.Set("user_id", "")
					c.Next()
					return
	}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
							return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return []byte(os.Getenv("JWT_KEY")), nil
			})

			if err != nil || token == nil || !token.Valid {
					c.Set("user_id", "")
					c.Next()
					return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
					c.Set("user_id", "")
					c.Next()
					return
			}

			userIDValue, exists := claims["user_id"]
			if !exists || userIDValue == nil {
					c.Set("user_id", "")
					c.Next()
					return
			}

			var userID string
			switch v := userIDValue.(type) {
			case string:
					userID = v
			case float64:
					userID = fmt.Sprintf("%.0f", v) 
			default:
					userID = fmt.Sprintf("%v", v)
			}

			c.Set("user_id", userID)
			c.Next()
	}
}

func extractTokenFromHeader(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")

	if bearerToken == "" {
		return ""
	}

	if strings.Split(bearerToken, " ")[0] != prefix {
		return ""
	}

	return strings.Split(bearerToken, " ")[1]
}

// GuestSessionMiddleware создает или обновляет гостевую сессию
func (h *Handler) GuestSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie(guestSessionCookie)
		
		needNewSession := false
		if err != nil || sessionID == "" {
			needNewSession = true
		} else {
			// Проверяем существование и время создания сессии в Redis
			createdAt, err := h.Repository.GetGuestSessionCreatedAt(context.Background(), sessionID)
			if err != nil || createdAt.IsZero() {
				needNewSession = true
			} else {
				// Проверяем, не истек ли срок действия сессии
				if time.Since(createdAt) > guestSessionTTL {
					// Удаляем старую сессию
					h.Repository.DeleteGuestSession(context.Background(), sessionID)
					needNewSession = true
				}
			}
		}
		
		if needNewSession {
			// Создаем новую сессию
			sessionID = uuid.New().String()
			err := h.Repository.CreateGuestSession(context.Background(), sessionID, guestSessionTTL)
			if err != nil {
				c.Next()
				return
			}
			
			// Устанавливаем cookie
			c.SetCookie(
				guestSessionCookie,
				sessionID,
				int(guestSessionTTL.Seconds()),
				"/",
				"",
				false,
				true,
			)
		}
		
		c.Set("guest_session_id", sessionID)
		c.Next()
	}
}