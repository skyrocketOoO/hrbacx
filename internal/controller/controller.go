package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/skyrocketOoO/hrbacx/internal/common"
)

type Usecase interface {
	AddLeader(leaderID string, roleID string) error
	AssignPermission(objectID, permissionType, roleID string) error
	AssignRole(userID, roleID string) error
	CheckPermission(userID, permissionType, objectID string) (ok bool, err error)
	ClearAll() error
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
		log.Error().Msg(err.Error())
		c.Status(400)
		return
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
		log.Error().Msg(err.Error())
		c.Status(400)
		return
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
		log.Error().Msg(err.Error())
		c.Status(400)
		return
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
		log.Error().Msg(err.Error())
		c.Status(500)
		return
	}
	if !ok {
		c.Status(400)
		return
	}
	c.Status(200)
}

func (h *Handler) ClearAll(c *gin.Context) {
	if err := h.Usecase.ClearAll(); err != nil {
		log.Error().Msg(err.Error())
		c.Status(500)
		return
	}
	c.Status(200)
}
