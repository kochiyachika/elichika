package db

import (
	"elichika/config"
	"elichika/handler"
	"elichika/klab"
	"elichika/model"
	"elichika/serverdb"
	"elichika/utils"

	"encoding/json"
	"fmt"
	"strings"
	"time"

	"xorm.io/xorm"

	"github.com/tidwall/gjson"
)

var (
	presetDataPath = "assets/preset/"
	userDataPath   = "assets/userdata/"
	UserID         = 588296696
	IsGlobal       = true
)

func GetJsonData(fileName string) string {
	userDataFile := userDataPath + fileName
	if utils.PathExists(userDataFile) {
		return utils.ReadAllText(userDataFile)
	}

	presetDataFile := presetDataPath + fileName
	if !utils.PathExists(presetDataFile) {
		panic("File not exists")
	}

	userData := utils.ReadAllText(presetDataFile)
	utils.WriteAllText(userDataFile, userData)

	return userData
}

func GetLiveDeckData() string {
	if IsGlobal {
		return GetJsonData("liveDeck_gl.json")
	}
	return GetJsonData("liveDeck.json")
}

func ImportJsonUser() {
	fmt.Println("Insert new user: ", UserID)
	data := GetJsonData("userStatus.json")
	status := model.UserStatus{}
	if err := json.Unmarshal([]byte(data), &status); err != nil {
		panic(err)
	}
	status.UserID = UserID
	status.LivePointBroken = 1000 // give 1000 LP
	// insert into the db
	_, err := serverdb.Engine.Table("s_user_info").AllCols().Insert(&status)
	if err != nil {
		panic(err)
	}
}

func LoadMemberFromJson() []model.UserMemberInfo {
	members := []model.UserMemberInfo{}
	userMemberInfo := model.UserMemberInfo{}
	memberData := gjson.Parse(GetJsonData("memberSettings.json"))
	memberData.Get("user_member_by_member_id").ForEach(func(key, value gjson.Result) bool {
		if value.IsObject() {
			if err := json.Unmarshal([]byte(value.String()), &userMemberInfo); err != nil {
				panic(err)
			}
			userMemberInfo.UserID = UserID
			members = append(members, userMemberInfo)
		}
		return true
	})
	return members
}

func LoadLiveDeckAndLivePartyFromJson() ([]model.UserLiveDeck, []model.UserLiveParty) {
	fmt.Println("importing live deck data to db")
	liveDecks := []model.UserLiveDeck{}
	liveDeckInfo := model.UserLiveDeck{}
	liveDeckData := gjson.Parse(GetLiveDeckData())
	liveDeckData.Get("user_live_deck_by_id").ForEach(func(key, value gjson.Result) bool {
		if value.IsObject() {
			if err := json.Unmarshal([]byte(value.String()), &liveDeckInfo); err != nil {
				panic(err)
			}
			liveDeckInfo.UserID = UserID
			liveDecks = append(liveDecks, liveDeckInfo)
		}
		return true
	})

	liveParties := []model.UserLiveParty{}

	fmt.Println("importing live party data to db")
	livePartyInfo := model.UserLiveParty{}
	var livePartyData []json.RawMessage
	decoder := json.NewDecoder(strings.NewReader(liveDeckData.Get("user_live_party_by_id").String()))
	decoder.UseNumber()
	err := decoder.Decode(&livePartyData)
	if err != nil {
		panic(err)
	}
	for i := 1; i < len(livePartyData); i += 2 {
		err := json.Unmarshal(livePartyData[i], &livePartyInfo)
		if err != nil {
			panic(err)
		}
		livePartyInfo.UserID = UserID
		liveParties = append(liveParties, livePartyInfo)
	}
	return liveDecks, liveParties
}

func LoadLessonDeckFromJson() []model.UserLessonDeck {
	lessonDecks := []model.UserLessonDeck{}
	fmt.Println("importing lesson deck data to db")
	lessonData := gjson.Parse(GetJsonData("lessonDeck.json"))
	userLessonDeck := model.UserLessonDeck{}
	lessonData.Get("user_lesson_deck_by_id").ForEach(func(key, value gjson.Result) bool {
		if value.IsObject() {
			if err := json.Unmarshal([]byte(value.String()), &userLessonDeck); err != nil {
				panic(err)
			}
			userLessonDeck.UserID = UserID
			lessonDecks = append(lessonDecks, userLessonDeck)
		}
		return true
	})
	return lessonDecks
}

func LoadCardFromJson() []model.UserCard {
	cards := []model.UserCard{}
	fmt.Println("importing json card data to db")
	cardData := gjson.Parse(GetJsonData("userCard.json"))
	userCard := model.UserCard{}
	cardData.Get("user_card_by_card_id").ForEach(func(key, value gjson.Result) bool {
		if value.IsObject() {
			if err := json.Unmarshal([]byte(value.String()), &userCard); err != nil {
				panic(err)
			}
			userCard.UserID = UserID
			cards = append(cards, userCard)
		}
		return true
	})
	return cards
}

func LoadSuitFromJson() []model.UserSuit {
	anys := []any{}
	suits := []model.UserSuit{}
	fmt.Println("importing json suit data to db")
	err := json.Unmarshal([]byte(gjson.Get(GetJsonData("login.json"), "user_model.user_suit_by_suit_id").String()), &anys)
	if err != nil {
		panic(err)
	}
	n := len(anys)
	for i := 0; i < n; i += 2 {
		suits = append(suits, model.UserSuit{
			UserID:       UserID,
			SuitMasterID: int(anys[i].(float64)),
			IsNew:        false})
	}
	return suits
}

func InsertAccount(session *serverdb.Session, members []model.UserMemberInfo,
	liveDecks []model.UserLiveDeck,
	liveParties []model.UserLiveParty,
	lessonDecks []model.UserLessonDeck,
	cards []model.UserCard,
	suits []model.UserSuit,
	lovePanels []model.UserMemberLovePanel,
	trainingTreeCells []model.TrainingTreeCell) {

	session.InsertMembers(members)
	session.InsertLiveDecks(liveDecks)
	session.InsertLiveParties(liveParties)
	session.InsertLessonDecks(lessonDecks)
	session.InsertCards(cards)
	serverdb.InsertUserSuits(suits, nil)
	for _, lovePanel := range lovePanels {
		affected, err := serverdb.Engine.Table("s_user_member").AllCols().
			Where("user_id = ? AND member_master_id = ?", lovePanel.UserID, lovePanel.MemberID).
			Update(&lovePanel)
		utils.CheckErr(err)
		if affected == 0 {
			panic("affected = 0")
		}
	}
	session.InsertTrainingCells(&trainingTreeCells)
}

func ImportFromJson() {
	fmt.Println("Importing account data from json")
	ImportJsonUser()
	session := serverdb.GetSession(nil, UserID)

	// member data
	members := LoadMemberFromJson()
	lovePanels := []model.UserMemberLovePanel{}
	for _, member := range members {
		lovePanel := model.UserMemberLovePanel{
			UserID:                    member.UserID,
			MemberID:                  member.MemberMasterID,
			LovePanelLevel:            32,
			LovePanelLastLevelCellIDs: []int{}}
		for i := 1; i <= 5; i++ {
			lovePanel.LovePanelLastLevelCellIDs = append(lovePanel.LovePanelLastLevelCellIDs,
				310000+i*1000+member.MemberMasterID)
		}
		lovePanels = append(lovePanels, lovePanel)
	}

	liveDecks, liveParties := LoadLiveDeckAndLivePartyFromJson()
	lessonDecks := LoadLessonDeckFromJson()
	cards := LoadCardFromJson()
	suits := LoadSuitFromJson()
	trainingTreeCells := []model.TrainingTreeCell{}
	for _, card := range cards {
		for cellID := 1; cellID <= 89; cellID++ {
			// this is a bit lazy, could have used the correct number instead
			trainingTreeCells = append(trainingTreeCells,
				model.TrainingTreeCell{
					UserID:       card.UserID,
					CardMasterID: card.CardMasterID,
					CellID:       cellID,
					ActivatedAt:  int64(1688094000 + cellID)})
		}
	}
	InsertAccount(&session, members, liveDecks, liveParties, lessonDecks, cards, suits, lovePanels, trainingTreeCells)
}

func CreateNewUser() {
	fmt.Println("Insert new user: ", UserID)
	status := model.UserStatus{
		UserID:                                  UserID,
		LastLoginAt:                             time.Now().Unix(),
		Rank:                                    1,
		Exp:                                     0,
		RecommendCardMasterID:                   100011001, // Honoka
		MaxFriendNum:                            99,
		LivePointFullAt:                         time.Now().Unix(),
		LivePointBroken:                         1000, // Start with 1000 LP
		LivePointSubscriptionRecoveryDailyCount: 1,
		LivePointSubscriptionRecoveryDailyResetAt: 1688137200,
		ActivityPointCount:                        3,
		ActivityPointResetAt:                      1688137200,
		ActivityPointPaymentRecoveryDailyCount:    10,
		ActivityPointPaymentRecoveryDailyResetAt:  1688137200,
		GameMoney:                                 1000000000,
		CardExp:                                   1000000000,
		FreeSnsCoin:                               50000,
		AppleSnsCoin:                              25000,
		GoogleSnsCoin:                             25000,
		SubscriptionCoin:                          50000,
		BirthDate:                                 202307,
		BirthMonth:                                7,
		BirthDay:                                  1,
		LatestLiveDeckID:                          1,
		MainLessonDeckID:                          1,
		FavoriteMemberID:                          1,
		LastLiveDifficultyID:                      10001101,
		LpMagnification:                           1,
		EmblemID:                                  10500521, // new player
		DeviceToken:                               "",
		TutorialPhase:                             99,
		TutorialEndAt:                             1622217482,
		LoginDays:                                 1221,
		NaviTapCount:                              0,
		NaviTapRecoverAt:                          1688137200,
		IsAutoMode:                                false,
		MaxScoreLiveDifficultyMasterID:            10001101,
		LiveMaxScore:                              0,
		MaxComboLiveDifficultyMasterID:            10001101,
		LiveMaxCombo:                              0,
		LessonResumeStatus:                        1,
		AccessoryBoxAdditional:                    400,
		TermsOfUseVersion:                         2,
		BootstrapSifidCheckAt:                     1692471111782,
		GdprVersion:                               0,
		MemberGuildMemberMasterID:                 1,
		MemberGuildLastUpdatedAt:                  1659485328,
		Cash:                                      0,
	}
	status.Name.DotUnderText = "Newcomer"
	status.Nickname.DotUnderText = "Newcomer"
	status.Message.DotUnderText = "Hello!"
	// insert into the db
	_, err := serverdb.Engine.Table("s_user_info").AllCols().Insert(&status)
	if err != nil {
		panic(err)
	}
}

func ImportMinimalAccount() {
	
	CreateNewUser()
	session := serverdb.GetSession(nil, UserID)

	var masterdatadb *xorm.Engine
	if IsGlobal {
		masterdatadb = config.MasterdataEngGl
	} else {
		masterdatadb = config.MasterdataEngJp
	}

	// member and card data
	members := []model.UserMemberInfo{}
	cards := []model.UserCard{}
	type MemberInit struct {
		MemberMasterID           int `xorm:"'member_m_id'"`
		SuitMasterID             int `xorm:"'suit_m_id'"`
		CustomBackgroundMasterID int `xorm:"'custom_background_m_id'"`
		LovePointLimit           int `xorm:"love_point_limit"`
	}

	memberInits := []MemberInit{}
	err := masterdatadb.Table("m_member_init").Find(&memberInits)
	utils.CheckErr(err)
	if len(memberInits) != 30 {
		panic("wrong number of member")
	}
	for _, memberInit := range memberInits {
		members = append(members, model.UserMemberInfo{
			UserID:                   UserID,
			MemberMasterID:           memberInit.MemberMasterID,
			CustomBackgroundMasterID: memberInit.CustomBackgroundMasterID,
			SuitMasterID:             memberInit.SuitMasterID,
			LovePoint:                0,
			LoveLevel:                1,
			LovePointLimit:           memberInit.LovePointLimit,
			ViewStatus:               1,
			IsNew:                    false,
			OwnedCardCount:           1,
			AllTrainingCardCount:     0,
		})
		cards = append(cards, model.UserCard{
			UserID:                     UserID,
			CardMasterID:               memberInit.SuitMasterID,
			Level:                      1,
			Exp:                        0,
			LovePoint:                  0,
			IsFavorite:                 false,
			IsAwakening:                false,
			IsAwakeningImage:           false,
			IsAllTrainingActivated:     false,
			TrainingActivatedCellCount: 0,
			MaxFreePassiveSkill:        1, // all R
			Grade:                      0,
			TrainingLife:               0,
			TrainingAttack:             0,
			TrainingDexterity:          0,
			ActiveSkillLevel:           1,
			PassiveSkillALevel:         1,
			PassiveSkillBLevel:         1,
			PassiveSkillCLevel:         1,
			AdditionalPassiveSkill1ID:  0,
			AdditionalPassiveSkill2ID:  0,
			AdditionalPassiveSkill3ID:  0,
			AdditionalPassiveSkill4ID:  0,
			AcquiredAt:                 time.Now().Unix(),
			IsNew:                      false,
			LivePartnerCategories:      0,
			LiveJoinCount:              0,
			ActiveSkillPlayCount:       0,
		})
	}

	suits := []model.UserSuit{}
	suitMasterIDs := []int{}
	err = masterdatadb.Table("m_suit").Where("suit_release_route == 2").Cols("id").Find(&suitMasterIDs)
	utils.CheckErr(err)
	for _, suitMasterID := range suitMasterIDs {
		suits = append(suits, model.UserSuit{
			UserID:       UserID,
			SuitMasterID: suitMasterID,
			IsNew:        false,
		})
	}
	lovePanels := []model.UserMemberLovePanel{}     // there's no bond board unlocked
	trainingTreeCells := []model.TrainingTreeCell{} // there's no training cell unlocked

	// live deck data
	liveDecks := []model.UserLiveDeck{}
	liveParties := []model.UserLiveParty{}
	for i := 1; i <= 20; i++ {
		cid := [10]int{}
		// this order isn't actually correct to the official server
		for j := 1; j <= 9; j++ {
			cid[j] = klab.DefaultSuitMasterIDFromMemberMasterID(j + 100*((i-1)%3))
		}
		liveDeck := model.UserLiveDeck{
			UserID:         UserID,
			UserLiveDeckID: i,
			CardMasterID1:  cid[1],
			CardMasterID2:  cid[2],
			CardMasterID3:  cid[3],
			CardMasterID4:  cid[4],
			CardMasterID5:  cid[5],
			CardMasterID6:  cid[6],
			CardMasterID7:  cid[7],
			CardMasterID8:  cid[8],
			CardMasterID9:  cid[9],
			SuitMasterID1:  cid[1],
			SuitMasterID2:  cid[2],
			SuitMasterID3:  cid[3],
			SuitMasterID4:  cid[4],
			SuitMasterID5:  cid[5],
			SuitMasterID6:  cid[6],
			SuitMasterID7:  cid[7],
			SuitMasterID8:  cid[8],
			SuitMasterID9:  cid[9],
		}
		liveDeck.Name.DotUnderText = fmt.Sprintf("Show formation %d", i)
		liveDecks = append(liveDecks, liveDeck)
		for j := 1; j <= 3; j++ {
			liveParty := model.UserLiveParty{
				UserID:         UserID,
				PartyID:        i*100 + j,
				UserLiveDeckID: i,
				CardMasterID1:  cid[(j-1)*3+1],
				CardMasterID2:  cid[(j-1)*3+2],
				CardMasterID3:  cid[(j-1)*3+3],
			}
			roles := []int{}
			err = masterdatadb.Table("m_card").
				Where("id IN (?,?,?)", liveParty.CardMasterID1,
					liveParty.CardMasterID2,
					liveParty.CardMasterID3).
				Cols("role").Find(&roles)
			partyIcon, partyName := handler.GetPartyInfoByRoleIds(roles)
			liveParty.IconMasterID = partyIcon
			liveParty.Name.DotUnderText = handler.GetRealPartyName(partyName)
			liveParties = append(liveParties, liveParty)
		}
	}

	lessonDecks := []model.UserLessonDeck{}
	for i := 1; i <= 20; i++ {
		lessonDecks = append(lessonDecks, model.UserLessonDeck{
			UserID:           UserID,
			UserLessonDeckID: i,
			Name:             fmt.Sprintf("Training formation %d", i),
		})
	}
	InsertAccount(&session, members, liveDecks, liveParties, lessonDecks, cards, suits, lovePanels, trainingTreeCells)
}

func Account(args []string) {
	if len(args) != 2 {
		fmt.Println("Invalid params:", args)
		fmt.Println("Required: gl/jp and json/new")
	}
	IsGlobal = (args[0] == "gl")

	if args[1] == "json" { // import from existing jsons
		ImportFromJson()
	} else {
		ImportMinimalAccount() // minimal account
	}
}
