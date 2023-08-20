package handler

import (
	"elichika/config"
	"elichika/serverdb"
	"elichika/utils"

	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	// "github.com/tidwall/sjson"
)

func FetchProfile(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	type FetchProfileReq struct {
		UserID int `json:"user_id"`
	}
	req := FetchProfileReq{}
	if err := json.Unmarshal([]byte(reqBody), &req); err != nil {
		panic(err)
	}

	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)
	profile := session.FetchProfile(req.UserID)

	signBody, err := json.Marshal(profile)
	if err != nil {
		panic(err)
	}

	resp := SignResp(ctx.GetString("ep"), string(signBody), config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
	// fmt.Println(resp)
}

func SetProfile(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")
	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)
	// fmt.Println(reqBody)

	req := gjson.Parse(reqBody).Array()[0]
	if req.Get("name").String() != "" {
		session.UserStatus.Name.DotUnderText = gjson.Parse(reqBody).Array()[0].Get("name").String()
	} else if req.Get("nickname").String() != "" {
		session.UserStatus.Nickname.DotUnderText = gjson.Parse(reqBody).Array()[0].Get("nickname").String()
	} else if req.Get("message").String() != "" {
		session.UserStatus.Message.DotUnderText = gjson.Parse(reqBody).Array()[0].Get("message").String()
	}

	signBody := session.Finalize(GetData("setProfile.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func SetRecommendCard(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")
	// fmt.Println(reqBody)
	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)
	cardMasterId := int(gjson.Parse(reqBody).Array()[0].Get("card_master_id").Int())
	session.UserStatus.RecommendCardMasterID = cardMasterId

	signBody := session.Finalize(GetData("setRecommendCard.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func SetLivePartner(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	type SetLivePartnerReq struct {
		LivePartnerCategoryID int `json:"live_partner_category_id"`
		CardMasterID          int `json:"card_master_id"`
	}
	req := SetLivePartnerReq{}
	if err := json.Unmarshal([]byte(reqBody), &req); err != nil {
		panic(err)
	}

	// set the bit on the correct card
	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)
	newCard := session.GetUserCard(req.CardMasterID)
	newCard.LivePartnerCategories |= (1 << req.LivePartnerCategoryID)
	session.UpdateUserCard(newCard)

	// remove the bit on the other cards
	partnerCards := serverdb.FetchPartnerCards(UserID)
	for _, card := range partnerCards {
		if card.CardMasterID == req.CardMasterID {
			continue
		}
		if (card.LivePartnerCategories & (1 << req.LivePartnerCategoryID)) != 0 {
			card.LivePartnerCategories ^= (1 << req.LivePartnerCategoryID)
			session.UpdateUserCard(card)
		}
	}

	session.Finalize("{}", "")
	// this is correct, the server send {}
	resp := SignResp(ctx.GetString("ep"), "{}", config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func SetScoreOrComboLive(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	type SetScoreOrComboReq struct {
		LiveDifficultyMasterID int `json:"live_difficulty_master_id"`
	}
	req := SetScoreOrComboReq{}
	err := json.Unmarshal([]byte(reqBody), &req)
	utils.CheckErr(err)

	userID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, userID)
	customSetProfile := session.GetUserCustomSetProfile()
	if ctx.Request.URL.Path == "/userProfile/setScoreLive" {
		customSetProfile.VoltageLiveDifficultyID = req.LiveDifficultyMasterID
	} else {
		customSetProfile.ComboLiveDifficultyID = req.LiveDifficultyMasterID
	}
	session.SetUserCustomSetProfile(customSetProfile)
	resp := SignResp(ctx.GetString("ep"), reqBody, config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
