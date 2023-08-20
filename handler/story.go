package handler

import (
	"elichika/config"
	"elichika/serverdb"

	"net/http"

	"github.com/gin-gonic/gin"
)

func FinishStory(ctx *gin.Context) {
	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)
	signBody := session.Finalize(GetData("finishStory.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func FinishStoryMain(ctx *gin.Context) {
	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)
	signBody := session.Finalize(GetData("finishUserStoryMain.json"), "user_model_diff")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func FinishStoryLinkage(ctx *gin.Context) {
	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)
	signBody := session.Finalize(GetData("finishStoryLinkage.json"), "user_model_diff")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
