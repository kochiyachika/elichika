package handler

import (
	"elichika/config"
	"elichika/serverdb"

	"encoding/json"
	"net/http"
	// "fmt"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	// "github.com/tidwall/sjson"
)

func SaveUserNaviVoice(ctx *gin.Context) {
	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)
	signBody := session.Finalize(GetData("saveUserNaviVoice.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func TapLovePoint(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	type TapLovePointReq struct {
		MemberMasterID int `json:"member_master_id"`
	}

	req := TapLovePointReq{}
	if err := json.Unmarshal([]byte(reqBody), &req); err != nil {
		panic(err)
	}
	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)
	session.AddLovePoint(req.MemberMasterID, 20)
	signBody := session.Finalize(GetData("saveUserNaviVoice.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
