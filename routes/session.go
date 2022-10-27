package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
	"xgss/model"
	session2 "xgss/session"
)

func SessionRoutes(group *gin.Engine) {
	const OFFSET = len("BEARER ")
	routerGroup := group.RouterGroup.Group("/session")

	routerGroup.POST("", func(context *gin.Context) {
		var user model.User
		err := context.ShouldBindJSON(&user)

		if err != nil {
			context.Status(http.StatusNotAcceptable)
			return
		}

		var timestamp = time.Now().Unix()

		token, tokenError := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email":     user.Email,
			"password":  user.Password,
			"role":      user.Role,
			"timestamp": timestamp,
		}).SigningString()

		if tokenError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"token_error": tokenError.Error()})
			return
		}

		if !session2.Exists(token) {
			session2.NewSession(
				token, user.Email, user.Password, user.Role, timestamp)

			context.JSON(http.StatusCreated, gin.H{"token": token})
		} else {
			context.Status(http.StatusConflict)
		}
		return
	})

	routerGroup.PATCH("", func(context *gin.Context) {
		var data model.SessionData

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

	routerGroup.GET("", func(context *gin.Context) {
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

	routerGroup.PUT("", func(context *gin.Context) {
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

	routerGroup.DELETE("", func(context *gin.Context) {
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
