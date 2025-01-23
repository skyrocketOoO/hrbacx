package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/skyrocketOoO/hrbacx/internal/common"
)

type Usecase interface {
	AddLeader(leaderID string, roleID string) error
	AssignPermission(objectID, permissionType, roleID string) error
	AssignRole(userID, roleID string) error
	CheckPermission(userID, permissionType, objectID string) (ok bool, err error)
}

type Handler struct {
	Usecase Usecase
}

func NewHandler(u Usecase) *Handler {
	return &Handler{
		Usecase: u,
	}
}

func (h *Handler) AddLeader(c *gin.Context) {
	type Req struct {
		LeaderID string `json:"leaderID"`
		RoleID   string `json:"roleID"`
	}

	var req Req
	if ok := common.BindAndValidate(c, &req); !ok {
		return
	}

	if err := h.Usecase.AddLeader(req.LeaderID, req.RoleID); err != nil {
		c.Status(400)
	}
	c.Status(200)
}

func (h *Handler) AssignPermission(c *gin.Context) {
	type Req struct {
		ObjectID       string `json:"objectID"`
		PermissionType string `json:"permissionType"`
		RoleID         string `json:"roleID"`
	}

	var req Req
	if ok := common.BindAndValidate(c, &req); !ok {
		return
	}

	if err := h.Usecase.AssignPermission(req.ObjectID, req.PermissionType, req.RoleID); err != nil {
		c.Status(400)
	}
	c.Status(200)
}

func (h *Handler) AssignRole(c *gin.Context) {
	type Req struct {
		UserID string `json:"userID"`
		RoleID string `json:"roleID"`
	}

	var req Req
	if ok := common.BindAndValidate(c, &req); !ok {
		return
	}

	if err := h.Usecase.AssignRole(req.UserID, req.RoleID); err != nil {
		c.Status(400)
	}
	c.Status(200)
}

func (h *Handler) CheckPermission(c *gin.Context) {
	type Req struct {
		ObjectID       string `json:"objectID"`
		PermissionType string `json:"permissionType"`
		UserID         string `json:"userID"`
	}

	var req Req
	if ok := common.BindAndValidate(c, &req); !ok {
		return
	}

	ok, err := h.Usecase.CheckPermission(req.UserID, req.PermissionType, req.ObjectID)
	if err != nil {
		c.Status(500)
	}
	if !ok {
		c.Status(400)
	}
	c.Status(200)
}
