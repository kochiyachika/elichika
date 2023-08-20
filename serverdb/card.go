package serverdb

import (
	"elichika/model"
	"elichika/utils"

	"fmt"
	"time"

	"xorm.io/xorm"
)

// fetch a card, use the value in diff is present, otherwise fetch from db
func (session *Session) GetUserCard(cardMasterID int) model.UserCard {
	card, exist := session.CardDiffs[cardMasterID]
	if exist {
		return card
	}
	card = model.UserCard{}
	exists, err := Engine.Table("s_user_card").
		Where("user_id = ? AND card_master_id = ?", session.UserStatus.UserID, cardMasterID).Get(&card)
	if err != nil {
		panic(err)
	}

	if !exists {
		card = model.UserCard{
			UserID:                     session.UserStatus.UserID,
			CardMasterID:               cardMasterID,
			Level:                      1,
			Exp:                        0,
			LovePoint:                  0,
			IsFavorite:                 false,
			IsAwakening:                false,
			IsAwakeningImage:           false,
			IsAllTrainingActivated:     false,
			TrainingActivatedCellCount: 0,
			MaxFreePassiveSkill:        -1,
			Grade:                      -1, // check this for new card
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
			IsNew:                      true,
			LivePartnerCategories:      0,
			LiveJoinCount:              0,
			ActiveSkillPlayCount:       0,
		}
		exists, err := session.Ctx.MustGet("masterdata.db").(*xorm.Engine).Table("m_card").Where("id = ?", cardMasterID).
			Cols("passive_skill_slot").Get(&card.MaxFreePassiveSkill)
		utils.CheckErrMustExist(err, exists)
	}
	return card
}

func (session *Session) GetAllCards() []model.UserCard {
	var cards []model.UserCard
	err := Engine.Table("s_user_card").Where("user_id = ?", session.UserStatus.UserID).Find(&cards)
	if err != nil {
		panic(err)
	}
	return cards
}

func (session *Session) UpdateUserCard(card model.UserCard) {
	session.CardDiffs[card.CardMasterID] = card
}

func (session *Session) FinalizeCardDiffs(dbSession *xorm.Session) []any {
	userCardByCardID := []any{}
	for cardMasterID, card := range session.CardDiffs {
		userCardByCardID = append(userCardByCardID, cardMasterID)
		userCardByCardID = append(userCardByCardID, card)
		affected, err := dbSession.Table("s_user_card").
			Where("user_id = ? AND card_master_id = ?", session.UserStatus.UserID, cardMasterID).AllCols().Update(card)
		utils.CheckErr(err)
		if affected == 0 {
			_, err := dbSession.Table("s_user_card").AllCols().Insert(card)
			utils.CheckErr(err)
		}
	}
	return userCardByCardID
}

// insert all the cards
func (session *Session) InsertCards(cards []model.UserCard) {
	affected, err := Engine.Table("s_user_card").Insert(&cards)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted ", affected, " cards")
}
