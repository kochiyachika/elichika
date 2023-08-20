package serverdb

import (
	"elichika/config"
	"elichika/model"
	"elichika/utils"

	"fmt"

	"xorm.io/xorm"
)

var (
	Engine *xorm.Engine
)

func InitTable(tableName string, structure interface{}, overwrite bool) {
	exist, err := Engine.Table(tableName).IsTableExist(tableName)
	utils.CheckErr(err)

	if !exist {
		fmt.Println("Creating new table:", tableName)
		err = Engine.Table(tableName).CreateTable(structure)
		utils.CheckErr(err)
	} else {
		if !overwrite {
			fmt.Println("Skip existing table:", tableName)
			return
		}
		fmt.Println("Overwrite existing table:", tableName)
		err := Engine.DropTables(tableName)
		utils.CheckErr(err)
		err = Engine.Table(tableName).CreateTable(structure)
		utils.CheckErr(err)
	}
}

func InitTables(overwrite bool) {
	type DbUser struct {
		model.UserStatus           `xorm:"extends"`
		model.UserProfileLiveStats `xorm:"extends"`
	}

	InitTable("s_user_info", DbUser{}, overwrite)
	InitTable("s_user_custom_set_profile", model.UserCustomSetProfile{}, overwrite)
	InitTable("s_user_card", model.UserCard{}, overwrite)
	InitTable("s_user_suit", model.UserSuit{}, overwrite)
	InitTable("s_user_training_tree_cell", model.TrainingTreeCell{}, overwrite)

	type DbMember struct {
		model.UserMemberInfo      `xorm:"extends"`
		LovePanelLevel            int   `xorm:"'love_panel_level' default 1"`
		LovePanelLastLevelCellIds []int `xorm:"'love_panel_last_level_cell_ids' default '[]'"`
	}
	InitTable("s_user_member", DbMember{}, overwrite)
	InitTable("s_user_lesson_deck", model.UserLessonDeck{}, overwrite)
	InitTable("s_user_live_deck", model.UserLiveDeck{}, overwrite)
	InitTable("s_user_live_party", model.UserLiveParty{}, overwrite)
	InitTable("s_user_live_state", model.LiveState{}, overwrite)
	InitTable("s_user_play_list", model.UserPlayListItem{}, overwrite)
	type DbLiveRecord struct {
		model.UserLiveDifficultyRecord `xorm:"extends"`
		Voltage                        int   `xorm:"'last_clear_voltage'" json:"voltage"`
		IsCleared                      bool  `xorm:"'last_clear_is_cleared'" json:"is_cleared"`
		RecordedAt                     int64 `xorm:"'last_clear_recorded_at'" json:"recorded_at"`
		CardWithSuitDict               []int `xorm:"'last_clear_cards_and_suits'" json:"card_with_suit_dict"`
		SquadDict                      []any `xorm:"'squad_dict'" json:"squad_dict"`
	}
	InitTable("s_user_live_record", DbLiveRecord{}, overwrite)
	InitTable("s_user_trigger_basic", model.TriggerBasic{}, overwrite)
	InitTable("s_user_trigger_card_grade_up", model.TriggerCardGradeUp{}, overwrite)
	InitTable("s_user_trigger_member_love_level_up", model.TriggerMemberLoveLevelUp{}, overwrite)

	InitTable("s_gacha_appeal", model.GachaAppeal{}, overwrite)
	InitTable("s_gacha_draw", model.GachaDraw{}, overwrite)
	InitTable("s_gacha", model.Gacha{}, overwrite)
	InitTable("s_gacha_group", model.GachaGroup{}, overwrite)
	InitTable("s_gacha_card", model.GachaCard{}, overwrite)
	InitTable("s_gacha_guarantee", model.GachaGuarantee{}, overwrite)
}

func init() {
	var err error
	Engine, err = xorm.NewEngine("sqlite", config.ServerdataDb)
	utils.CheckErr(err)
	Engine.SetMaxOpenConns(50)
	Engine.SetMaxIdleConns(10)
}
