package serverdb

import (
	"elichika/model"
	"elichika/utils"

	// "encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	// "github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"xorm.io/xorm"
)

// A session is a complete transation between server and client
// so 1 session per request
// A session fetch the data needs to be modified, and sometime modify the data if the code is shared between handlers.
// session can use ctx to get things like user id / master db, but it should not make any network operation

type Session struct {
	Ctx                           *gin.Context
	UserStatus                    model.UserStatus
	CardDiffs                     map[int]model.UserCard
	UserMemberDiffs               map[int]model.UserMemberInfo
	UserLessonDeckDiffs           map[int]model.UserLessonDeck
	UserLiveDeckDiffs             map[int]model.UserLiveDeck
	UserLivePartyDiffs            map[int]model.UserLiveParty
	UserMemberLovePanelDiffs      map[int]model.UserMemberLovePanel
	UserLiveDifficultyRecordDiffs map[int]model.UserLiveDifficultyRecord
	UserSuitDiffs                 []model.UserSuit
	TriggerCardGradeUps           []any
	TriggerBasics                 []any
	TriggerMemberLoveLevelUps     []any
}

// Push update into the db and create the diff
// The actual response depend on the API, but they often contain the diff somewhere
// The mainKey is the key to the diff
func (session *Session) Finalize(jsonBody string, mainKey string) string {
	dbSession := Engine.NewSession()
	err := dbSession.Begin()
	utils.CheckErr(err)
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_member_by_member_id", session.FinalizeUserMemberDiffs(dbSession))
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_card_by_card_id", session.FinalizeCardDiffs(dbSession))
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_lesson_deck_by_id", session.FinalizeUserLessonDeckDiffs(dbSession))
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_live_deck_by_id", session.FinalizeUserLiveDeckDiffs(dbSession))
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_live_party_by_id", session.FinalizeUserLivePartyDiffs(dbSession))
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_suit_by_suit_id", session.FinalizeUserSuitDiffs(dbSession))
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_live_difficulty_by_difficulty_id", session.FinalizeLiveDifficultyRecords(dbSession))
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_status", session.FinalizeUserInfo(dbSession))
	memberLovePanels := session.FinalizeMemberLovePanelDiffs(dbSession)
	if len(memberLovePanels) != 0 {
		jsonBody, _ = sjson.Set(jsonBody, "member_love_panels", memberLovePanels)
	}
	dbSession.Commit()
	dbSession.Close()
	// this could be a delta patch, but we can just send the whole thing
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_play_list_by_id", session.GetUserPlayList())
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_info_trigger_card_grade_up_by_trigger_id", session.TriggerCardGradeUps)
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_info_trigger_basic_by_trigger_id", session.TriggerBasics)
	jsonBody, _ = sjson.Set(jsonBody, mainKey+".user_info_trigger_member_love_level_up_by_trigger_id", session.TriggerMemberLoveLevelUps)

	return jsonBody
}

// fetch the user, this is always sent back to client
func (session *Session) InitUser(userID int) {
	session.UserStatus.UserID = userID
	exists, err := Engine.Table("s_user_info").Where("user_id = ?", userID).Get(&session.UserStatus)
	if err != nil {
		panic(err)
	}
	if !exists {
		// insert user and stuff for now
		panic(fmt.Sprintf("user doesn't exist %d\nNote: use \"elichika make [json/new] [jp/gl]\" to init the db", userID))
	}
}

func (session *Session) FinalizeUserInfo(dbSession *xorm.Session) model.UserStatus {
	_, err := dbSession.Table("s_user_info").Where("user_id = ?", session.UserStatus.UserID).Update(&session.UserStatus)
	if err != nil {
		panic(err)
	}
	return session.UserStatus
}

func GetSession(ctx *gin.Context, userId int) Session {
	s := Session{}
	s.Ctx = ctx
	s.CardDiffs = make(map[int]model.UserCard)
	s.UserMemberDiffs = make(map[int]model.UserMemberInfo)
	s.UserLessonDeckDiffs = make(map[int]model.UserLessonDeck)
	s.UserLiveDeckDiffs = make(map[int]model.UserLiveDeck)
	s.UserLivePartyDiffs = make(map[int]model.UserLiveParty)
	s.UserMemberLovePanelDiffs = make(map[int]model.UserMemberLovePanel)
	s.UserLiveDifficultyRecordDiffs = make(map[int]model.UserLiveDifficultyRecord)
	s.TriggerCardGradeUps = make([]any, 0)
	s.TriggerBasics = make([]any, 0)
	s.TriggerMemberLoveLevelUps = make([]any, 0)
	s.UserSuitDiffs = make([]model.UserSuit, 0)
	s.InitUser(userId)
	return s
}
