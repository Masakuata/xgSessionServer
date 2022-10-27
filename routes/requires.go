package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xgss/session"
)

type Role struct {
	Role string `json:"role,omitempty" validate:"required" binding:"required"`
}

func RequireRoutes(group *gin.Engine) {
	const OFFSET = len("BEARER ")

	routerGroup := group.RouterGroup.Group("/requires_role")

	routerGroup.PUT("", func(context *gin.Context) {
		var role Role
		if err := context.ShouldBindJSON(&role); err != nil {
			context.Status(http.StatusNotAcceptable)
			return
		}

		token := context.GetHeader("Authorization")
		if len(token) > 0 {
			token = token[OFFSET:]

			if session.Exists(token) {
				if session.IsRole(token, role.Role) {
					context.Status(http.StatusOK)
				} else {
					context.Status(http.StatusForbidden)
				}
			}
		} else {
			context.Status(http.StatusUnauthorized)
		}
		return
	})
}
