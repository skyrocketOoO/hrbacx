package dm

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/skyrocketOoO/erx/erx"
)

type ErrResp struct {
	ID    string `json:"id"`
	Error string `json:"error"`
}

func RespErr(c *gin.Context, statusCode int, err error) {
	clientMsg, ok := erx.GetClientMsg(err)
	if !ok {
		c.JSON(statusCode, ErrResp{Error: err.Error()})
		return
	}

	c.JSON(statusCode, ErrResp{ID: clientMsg.ID, Error: clientMsg.Err})

	fullMsg, _ := erx.GetFullMsg(err)
	log.Error().Msg(fullMsg.ID)
	log.Error().Msg(fullMsg.Err)
	log.Error().Msg(fullMsg.CallStack)
}
