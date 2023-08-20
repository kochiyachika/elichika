package handler

import (
	"elichika/config"
	"elichika/serverdb"

	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

func FetchBootstrap(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	type BootstrapReq struct {
		BootstrapFetchTypes []int  `json:"bootstrap_fetch_types"`
		DeviceToken         string `json:"device_token"`
		DeviceName          string `json:"device_name"`
	}
	// fmt.Println(reqBody)
	req := BootstrapReq{}
	if err := json.Unmarshal([]byte(reqBody), &req); err != nil {
		panic(err)
	}

	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)
	session.UserStatus.BootstrapSifidCheckAt = time.Now().UnixMilli()
	session.UserStatus.DeviceToken = req.DeviceToken
	signBody := session.Finalize(GetData("fetchBootstrap.json"), "user_model_diff")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
	// fmt.Println(resp)
}

func GetClearedPlatformAchievement(ctx *gin.Context) {
	signBody := GetData("getClearedPlatformAchievement.json")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
