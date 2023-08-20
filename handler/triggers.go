package handler

import (
	"elichika/config"
	"elichika/model"
	"elichika/serverdb"

	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func TriggerReadCardGradeUp(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)
	req := model.TriggerReadReq{}
	if err := json.Unmarshal([]byte(reqBody), &req); err != nil {
		panic(err)
	}

	session.AddTriggerCardGradeUp(req.TriggerID, nil)
	resp := session.Finalize(GetData("userModel.json"), "user_model")
	resp = SignResp(ctx.GetString("ep"), resp, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
	// fmt.Println(resp)
}

func TriggerRead(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)
	req := model.TriggerReadReq{}
	if err := json.Unmarshal([]byte(reqBody), &req); err != nil {
		panic(err)
	}

	session.AddTriggerBasic(req.TriggerID, nil)
	resp := session.Finalize(GetData("userModel.json"), "user_model")
	resp = SignResp(ctx.GetString("ep"), resp, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
	// fmt.Println(resp)
}

func TriggerReadMemberLoveLevelUp(ctx *gin.Context) {
	// req is null, so we need to pull the triggers from db here
	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)

	triggers := session.GetAllTriggerMemberLoveLevelUps()
	for _, trigger := range triggers.Objects {
		session.AddTriggerMemberLoveLevelUp(trigger.TriggerID, nil)
	}

	resp := session.Finalize(GetData("userModel.json"), "user_model")
	resp = SignResp(ctx.GetString("ep"), resp, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
	// fmt.Println(resp)
}
