package serverdb

// give reward based on content type, content id, and content amount
// depend on the actual handler, negative amount can lead to a decrease as well
// ContentType: m_content_setting

import (
	"elichika/model"

	"fmt"
)

var (
	RewardHandler map[int]func(*Session, int, int)
)

func (session *Session) AddRewardContent(reward model.RewardByContent) {
	// fmt.Println(reward)
	handler, exists := RewardHandler[reward.ContentType]
	if !exists {
		fmt.Println("TODO: Add reward for content type ", reward.ContentType)
		return
	}
	handler(session, reward.ContentID, reward.ContentAmount)
}

func init() {
	RewardHandler = make(map[int]func(*Session, int, int))
	RewardHandler[1] = SnsCoinRewardHandler
	RewardHandler[7] = SuitRewardHandler
	RewardHandler[12] = TrainingMaterialRewardHandler
}

func SuitRewardHandler(session *Session, suitMasterID, _ int) {
	session.InsertUserSuit(model.UserSuit{
		UserID:       session.UserStatus.UserID,
		SuitMasterID: suitMasterID,
		IsNew:        true})
}

func SnsCoinRewardHandler(session *Session, snsCoinType, amount int) {
	fmt.Println("TODO: Add economy")
}

func TrainingMaterialRewardHandler(session *Session, trainingMaterialType, amount int) {
	fmt.Println("TODO: Add training items")
}
