package serverdb

import (
	"elichika/enum"
	"elichika/klab"
	"elichika/model"

	"fmt"
	"xorm.io/xorm"
)

func (session *Session) GetMember(memberMasterID int) model.UserMemberInfo {
	member, exist := session.UserMemberDiffs[memberMasterID]
	if exist {
		return member
	}
	exists, err := Engine.Table("s_user_member").
		Where("user_id = ? AND member_master_id = ?", session.UserStatus.UserID, memberMasterID).Get(&member)
	// inserted at login if not exist
	if err != nil {
		panic(err)
	}
	if !exists {
		panic("member not found")
	}
	return member
}

func (session *Session) GetAllMembers() []model.UserMemberInfo {
	members := []model.UserMemberInfo{}
	err := Engine.Table("s_user_member").Where("user_id = ?", session.UserStatus.UserID).Find(&members)
	if err != nil {
		panic(err)
	}
	return members
}

func (session *Session) UpdateMember(member model.UserMemberInfo) {
	session.UserMemberDiffs[member.MemberMasterID] = member
}

func (session *Session) InsertMembers(members []model.UserMemberInfo) {
	affected, err := Engine.Table("s_user_member").Insert(&members)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted ", affected, " members")
}

func (session *Session) FinalizeUserMemberDiffs(dbSession *xorm.Session) []any {
	userMemberByMemberID := []any{}
	for memberMasterID, member := range session.UserMemberDiffs {
		userMemberByMemberID = append(userMemberByMemberID, memberMasterID)
		userMemberByMemberID = append(userMemberByMemberID, member)
		affected, err := dbSession.Table("s_user_member").
			Where("user_id = ? AND member_master_id = ?", session.UserStatus.UserID, memberMasterID).AllCols().Update(member)
		if (err != nil) || (affected != 1) {
			panic(err)
		}
	}
	return userMemberByMemberID
}

// add love point and return the love point added (in case maxed out)
func (session *Session) AddLovePoint(memberID, point int) int {
	member := session.GetMember(memberID)
	if point > member.LovePointLimit-member.LovePoint {
		point = member.LovePointLimit - member.LovePoint
	}
	member.LovePoint += point

	oldLoveLevel := member.LoveLevel
	member.LoveLevel = klab.BondLevelFromBondValue(member.LovePoint)
	// unlock bond stories, unlock bond board
	if oldLoveLevel < member.LoveLevel {
		db := session.Ctx.MustGet("masterdata.db").(*xorm.Engine)

		rewards := []model.RewardByContent{}
		err := db.Table("m_member_love_level_reward").Where("member_m_id = ? AND love_level > ? and love_level <= ?",
			memberID, oldLoveLevel, member.LoveLevel).Find(&rewards)
		if err != nil {
			panic(err)
		}
		for i, _ := range rewards {
			session.AddRewardContent(rewards[i])
		}

		latestLovePanelLevel := klab.MaxLovePanelLevelFromLoveLevel(member.LoveLevel)
		currentLovePanel := session.GetMemberLovePanel(memberID)
		if (currentLovePanel.LovePanelLevel < latestLovePanelLevel) && (len(currentLovePanel.LovePanelLastLevelCellIDs) == 5) {
			currentLovePanel.LevelUp()
			session.AddTriggerBasic(0, &model.TriggerBasic{
				InfoTriggerType: enum.InfoTriggerTypeUnlockBondBoard,
				LimitAt:         nil,
				Description:     nil,
				ParamInt:        currentLovePanel.LovePanelLevel*1000 + currentLovePanel.MemberID})

			session.UpdateMemberLovePanel(currentLovePanel)
		}
		session.AddTriggerMemberLoveLevelUp(0,
			&model.TriggerMemberLoveLevelUp{
				TriggerID:       0,
				MemberMasterID:  memberID,
				BeforeLoveLevel: member.LoveLevel - 1})

	}
	session.UpdateMember(member)
	return point
}
