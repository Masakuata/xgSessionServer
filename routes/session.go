package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	session2 "test/session"
	"time"
)

type Login struct {
	Email    string `json:"email,omitempty" validate:"required" binding:"required"`
	Password string `json:"password,omitempty" validate:"required" binding:"required"`
}

type SessionData struct {
	Data map[string]any `json:"data,omitempty" validate:"required" binding:"required"`
}

func SessionRoutes(group *gin.Engine) {
	const OFFSET = len("BEARER ")
	session := group.RouterGroup.Group("/session")

	session.POST("", func(context *gin.Context) {
		var user Login
		err := context.ShouldBindJSON(&user)
		if err != nil {
			context.Status(http.StatusNotAcceptable)
			return
		}

		var timestamp = time.Now().Unix()

		token, tokenError := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email":     user.Email,
			"password":  user.Password,
			"timestamp": timestamp,
		}).SigningString()

		if tokenError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"token_error": tokenError.Error()})
			return
		}

		if !session2.Exists(token) {
			session2.NewSession(
				token, user.Email, user.Password, timestamp)

			context.JSON(http.StatusCreated, gin.H{"token": token})
		} else {
			context.Status(http.StatusConflict)
		}
		return
	})

	session.PATCH("", func(context *gin.Context) {
		var data SessionData

		err := context.ShouldBindJSON(&data)
		if err != nil {
			context.Status(http.StatusNotAcceptable)
			return
		}

		token := context.GetHeader("Authorization")
		if len(token) > 0 {
			token = token[OFFSET:]

			if session2.Exists(token) {
				session2.AddData(token, data.Data)
				context.Status(http.StatusOK)
			} else {
				context.Status(http.StatusNotFound)
			}
		} else {
			context.Status(http.StatusUnauthorized)
		}
		return
	})

	session.GET("", func(context *gin.Context) {
		token := context.GetHeader("Authorization")
		if len(token) > 0 {
			token = token[OFFSET:]

			if session2.Exists(token) {
				data := session2.GetData(token)
				if data != nil {
					context.JSON(http.StatusOK, data)
				} else {
					context.Status(http.StatusNoContent)
				}
			} else {
				context.Status(http.StatusNotFound)
			}
		} else {
			context.Status(http.StatusUnauthorized)
		}
		return
	})

	session.PUT("", func(context *gin.Context) {
		token := context.GetHeader("Authorization")
		if len(token) > 0 {
			token = token[OFFSET:]

			if session2.Exists(token) {
				session2.UpdateLifetime(token)
				context.Status(http.StatusOK)
			} else {
				context.Status(http.StatusNotFound)
			}
		} else {
			context.Status(http.StatusUnauthorized)
		}
		return
	})

	session.DELETE("", func(context *gin.Context) {
		token := context.GetHeader("Authorization")
		if len(token) > 0 {
			token = token[OFFSET:]
			if session2.Exists(token) {
				session2.Delete(token)
				context.Status(http.StatusOK)
			} else {
				context.Status(http.StatusNotFound)
			}
		} else {
			context.Status(http.StatusUnauthorized)
		}
		return
	})
}
