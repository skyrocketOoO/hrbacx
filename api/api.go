package api

import (
	"github.com/gin-gonic/gin"
	"github.com/skyrocketOoO/hrbacx/internal/controller"
)

func Bind(r *gin.Engine, h *controller.Handler) {
	r.POST("/addLeader", h.AddLeader)
	r.POST("/assignPermission", h.AssignPermission)
	r.POST("/assignRole", h.AssignRole)
	r.POST("/checkPermission", h.CheckPermission)
	r.POST("/clearAll", h.ClearAll)
}
