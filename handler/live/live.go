package live

import (
	"elichika/config"
	"elichika/handler"
	"elichika/model"
	"elichika/serverdb"
	"elichika/utils"

	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	// "xorm.io/xorm"
)

func FetchLiveMusicSelect(ctx *gin.Context) {
	now := time.Now()
	year, month, day := now.Year(), now.Month(), now.Day()
	tomorrow := time.Date(year, month, day+1, 0, 0, 0, 0, now.Location()).Unix()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	liveDailyList := []model.LiveDaily{}
	err := handler.MainEng.Table("m_live_daily").Where("weekday = ?", weekday).Cols("id,live_id").Find(&liveDailyList)
	utils.CheckErr(err)
	for k := range liveDailyList {
		liveDailyList[k].EndAt = int(tomorrow)
		liveDailyList[k].RemainingPlayCount = 5
		liveDailyList[k].RemainingRecoveryCount = 10
	}

	signBody := handler.GetData("fetchLiveMusicSelect.json")
	signBody, _ = sjson.Set(signBody, "weekday_state.weekday", weekday)
	signBody, _ = sjson.Set(signBody, "weekday_state.next_weekday_at", tomorrow)
	signBody, _ = sjson.Set(signBody, "live_daily_list", liveDailyList)
	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)
	signBody = session.Finalize(signBody, "user_model_diff")
	resp := handler.SignResp(ctx.GetString("ep"), signBody, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func FetchLivePartners(ctx *gin.Context) {
	// a set of partners player (i.e. friends and others), then fetch the card for them
	// this set include the current user, so we can use our own cards.
	// currently only have current user
	// note that all card are available, but we need to use the filter functionality to actually get them to show up.
	partnerIDs := []int{}
	UserID := ctx.GetInt("user_id")
	partnerIDs = append(partnerIDs, UserID)
	livePartners := []model.LiveStartLivePartner{}
	for _, partnerID := range partnerIDs {
		partner := model.LiveStartLivePartner{}
		partner.IsFriend = true
		serverdb.FetchDBProfile(partnerID, &partner)
		partnerCards := serverdb.FetchPartnerCards(partnerID) // model.UserCard
		if len(partnerCards) == 0 {
			continue
		}
		for _, card := range partnerCards {
			for i := 1; i <= 7; i++ {
				if (card.LivePartnerCategories & (1 << i)) != 0 {
					partnerCardInfo := serverdb.GetPartnerCardFromUserCard(card)
					partner.CardByCategory = append(partner.CardByCategory, i)
					partner.CardByCategory = append(partner.CardByCategory, partnerCardInfo)
				}
			}
		}
		livePartners = append(livePartners, partner)
	}

	signBody := "{}"
	signBody, _ = sjson.Set(signBody, "partner_select_state.live_partners", livePartners)
	signBody, _ = sjson.Set(signBody, "partner_select_state.friend_count", len(livePartners))
	resp := handler.SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	// fmt.Println(resp)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func LiveStart(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	req := model.LiveStartReq{}
	if err := json.Unmarshal([]byte(reqBody), &req); err != nil {
		panic(err)
	}
	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)

	session.UserStatus.LastLiveDifficultyID = req.LiveDifficultyID
	session.UserStatus.LatestLiveDeckID = req.DeckID

	// 保存请求包因为 /live/finish 接口的响应包里有部分字段不在该接口的请求包里
	// live is stored in db
	liveState := model.LiveState{}
	liveState.UserID = UserID
	liveState.PartnerUserID = req.PartnerUserID
	liveState.LiveID = time.Now().UnixNano()
	liveState.LiveType = 1 // not sure what this is
	liveState.IsPartnerFriend = true
	liveState.DeckID = req.DeckID
	liveState.CellID = req.CellID // cell id send player to the correct place after playing, normal live don't have cell id.

	liveNotes := utils.ReadAllText(fmt.Sprintf("assets/stages/%d.json", req.LiveDifficultyID))
	if liveNotes == "" {
		panic("歌曲情报信息不存在！(song doesn't exist)")
	}

	if err := json.Unmarshal([]byte(liveNotes), &liveState.LiveStage); err != nil {
		panic(err)
	}

	if req.IsAutoPlay {
		for k := range liveState.LiveStage.LiveNotes {
			liveState.LiveStage.LiveNotes[k].AutoJudgeType = 30
		}
	}

	if req.PartnerUserID != 0 {
		liveState.LivePartnerCard = serverdb.GetPartnerCardFromUserCard(
			serverdb.GetOtherUserCard(req.PartnerUserID, req.PartnerCardMasterID))
	}

	liveStartResp := session.Finalize(handler.GetData("userModelDiff.json"), "user_model_diff")
	liveStartResp, _ = sjson.Set(liveStartResp, "live", liveState)
	if req.PartnerUserID == 0 {
		liveStartResp, _ = sjson.Set(liveStartResp, "live.live_partner_card", nil)
	}
	serverdb.SaveLiveState(liveState)
	resp := handler.SignResp(ctx.GetString("ep"), liveStartResp, config.SessionKey)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
