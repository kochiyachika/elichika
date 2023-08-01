package handler

import (
	"elichika/config"
	"elichika/model"
	"elichika/serverdb"
	"fmt"
	"net/http"
	"strings"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func ExecuteLesson(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")
	// fmt.Println(reqBody)

	req := gjson.Parse(reqBody).Array()[0]
	deckId := req.Get("selected_deck_id").Int()

	var deckInfo string
	var actionList []model.LessonMenuAction
	gjson.Parse(GetUserData("lessonDeck.json")).Get("user_lesson_deck_by_id").ForEach(func(key, value gjson.Result) bool {
		if value.IsObject() && value.Get("user_lesson_deck_id").Int() == deckId {
			deckInfo = value.String()
			// fmt.Println("Deck Info:", deckInfo)

			gjson.Parse(deckInfo).ForEach(func(kk, vv gjson.Result) bool {
				// fmt.Printf("kk: %s, vv: %s\n", kk.String(), vv.String())
				if strings.Contains(kk.String(), "card_master_id") {
					actionList = append(actionList, model.LessonMenuAction{
						CardMasterID:                  vv.Int(),
						Position:                      0,
						IsAddedPassiveSkill:           true,
						IsAddedSpecialPassiveSkill:    true,
						IsRankupedPassiveSkill:        true,
						IsRankupedSpecialPassiveSkill: true,
						IsPromotedSkill:               true,
						MaxRarity:                     4,
						UpCount:                       1,
					})
				}
				return true
			})
			return false
		}
		return true
	})
	// fmt.Println(actionList)

	SetUserData("userStatus.json", "main_lesson_deck_id", deckId)

	signBody := GetData("executeLesson.json")
	signBody, _ = sjson.Set(signBody, "lesson_menu_actions.1", actionList)
	signBody, _ = sjson.Set(signBody, "lesson_menu_actions.3", actionList)
	signBody, _ = sjson.Set(signBody, "lesson_menu_actions.5", actionList)
	signBody, _ = sjson.Set(signBody, "lesson_menu_actions.7", actionList)
	signBody, _ = sjson.Set(signBody, "user_model_diff.user_status", GetUserStatus())
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	// fmt.Println(resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func ResultLesson(ctx *gin.Context) {
	userData := GetUserStatus()
	signBody, _ := sjson.Set(GetData("resultLesson.json"),
		"user_model_diff.user_status", userData)
	signBody, _ = sjson.Set(signBody, "selected_deck_id", userData["main_lesson_deck_id"])
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	// fmt.Println(resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func SkillEditResult(ctx *gin.Context) {
	reqBody := ctx.GetString("reqBody")
	// fmt.Println(reqBody)

	req := gjson.Parse(reqBody).Array()[0]

	session := serverdb.GetSession(UserID)
	skillList := req.Get("selected_skill_ids")
	skillList.ForEach(func(key, cardId gjson.Result) bool {
		if key.Int()%2 == 0 {
			cardInfo := session.GetCard(int(cardId.Int()))
			skills := skillList.Get(fmt.Sprintf("%d", key.Int()+1))
			cardJsonBytes, _ := json.Marshal(cardInfo)
			cardJson := string(cardJsonBytes)
			skills.ForEach(func(kk, vv gjson.Result) bool {
				skillIdKey := fmt.Sprintf("additional_passive_skill_%d_id", kk.Int()+1)
				cardJson, _ = sjson.Set(cardJson, skillIdKey, vv.Int())
				return true
			})
			if err := json.Unmarshal([]byte(cardJson), &cardInfo); err != nil {
				panic(err)
			}
			session.UpdateCard(cardInfo)
		}
		return true
	})
	signBody := session.Finalize(GetData("skillEditResult.json"), "user_model")
	signBody, _ = sjson.Set(signBody, "user_model.user_status", GetUserStatus())
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	// fmt.Println(resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}

func SaveDeckLesson(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0]
	fmt.Println(reqBody)
	type SaveDeckReq struct {
		DeckID int `json:"deck_id"`
		CardMasterIDs []int `json:"card_master_ids"`
	}
	req := SaveDeckReq{}
	if err := json.Unmarshal([]byte(reqBody.String()), &req); err != nil {
		panic(err)
	}

	session := serverdb.GetSession(UserID)
	userLessonDeck := session.GetLessonDeck(req.DeckID)
	deckByte, _ := json.Marshal(userLessonDeck)
	deckInfo := string(deckByte)
	for i := 0; i < len(req.CardMasterIDs); i+=2 {
		deckInfo, _ = sjson.Set(deckInfo, fmt.Sprintf("card_master_id_%d", req.CardMasterIDs[i]), req.CardMasterIDs[i+1])
	}
	if err := json.Unmarshal([]byte(deckInfo), &userLessonDeck); err != nil {
		panic(err)
	}
	session.UpdateLessonDeck(userLessonDeck)
	signBody := session.Finalize(GetUserData("userModel.json"), "user_model")
	signBody, _ = sjson.Set(signBody, "user_model.user_status", GetUserStatus())
	resp := SignResp(ctx.GetString("ep"), signBody, config.SessionKey)
	// fmt.Println(resp)

	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
}
