package vm

import (
	"strings"

	"github.com/KevinZonda/GoX/pkg/stringx"
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coodinator/models/openapi"
	"github.com/kigland/OpenHPC/coodinator/utils"
)

func token(c *gin.Context) {
	req := utils.BodyAsF[openapi.VmTokenReq](c)

	provider := MustGetProvider(c, req.Provider)
	if provider == nil {
		return
	}

	summary, ok := provider.TryGetContainer(req.Id)
	if !ok {
		c.JSON(400, gin.H{
			"message": "container not found",
		})
		return
	}

	inspect, err := provider.ContainerInspect(summary.ID)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "failed to fetch container",
		})
		return
	}

	env := inspect.Config.Env
	tokens := filterToken(env)

	c.JSON(200, openapi.VmTokenResp{
		Token: tokens,
	})
}

func filterToken(env []string) []string {
	var tokens []string
	for _, e := range env {
		if strings.Contains(stringx.TrimLower(e), "token") {
			tokens = append(tokens, e)
		}
	}
	return tokens
}
