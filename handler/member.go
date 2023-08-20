package handler

import (
	"elichika/config"
	"elichika/enum"
	"elichika/klab"
	"elichika/model"
	"elichika/serverdb"
	"elichika/utils"

	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	// "github.com/tidwall/sjson"
)

func OpenMemberLovePanel(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	type OpenMemberLovePanelReq struct {
		MemberID               int   `json:"member_id"`
		MemberLovePanelID      int   `json:"member_love_panel_id"`
		MemberLovePanelCellIDs []int `json:"member_love_panel_cell_ids"`
	}
	req := OpenMemberLovePanelReq{}
	err := json.Unmarshal([]byte(reqBody), &req)
	utils.CheckErr(err)
	UserID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, UserID)
	panel := session.GetMemberLovePanel(req.MemberID)

	panel.LovePanelLastLevelCellIDs = append(panel.LovePanelLastLevelCellIDs, req.MemberLovePanelCellIDs...)

	// if is full panel, then we have to send a basic info trigger to actually open up the next panel
	if len(panel.LovePanelLastLevelCellIDs) == 5 {
		member := session.GetMember(panel.MemberID)
		maxPanelLevel := klab.MaxLovePanelLevelFromLoveLevel(member.LoveLevel)
		if panel.LovePanelLevel < maxPanelLevel {
			// unlock the next board if available
			// otherwise it will be unlocked when bond level reach the value
			panel.LevelUp()
			session.AddTriggerBasic(0, &model.TriggerBasic{
				InfoTriggerType: enum.InfoTriggerTypeUnlockBondBoard,
				LimitAt:         nil,
				Description:     nil,
				ParamInt:        panel.LovePanelLevel*1000 + panel.MemberID})
		}
	}
	session.UpdateMemberLovePanel(panel)

	signBody := session.Finalize(GetData("userModel.json"), "user_model")
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	// fmt.Println(resp)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
