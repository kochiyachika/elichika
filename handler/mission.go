package handler

import (
	"elichika/config"
	"elichika/serverdb"

	"net/http"

	"github.com/gin-gonic/gin"
)

func FetchMission(ctx *gin.Context) {
	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)
	signBody := session.Finalize(GetData("fetchMission.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func ClearMissionBadge(ctx *gin.Context) {
	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)
	signBody := session.Finalize(GetData("clearMissionBadge.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
