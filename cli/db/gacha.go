package db

import (
	"elichika/config"
	"elichika/model"
	"elichika/serverdb"
	"elichika/utils"

	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

func InitGacha(args []string) {
	// insert some relevant gacha group, gacha card, and gacha guarantee
	masterdata := config.MasterdataEngGl

	// 9 groups for now:
	// (R, SR, UR) * (muse, aqours, niji)
	dbSession := serverdb.Engine.NewSession()
	err := dbSession.Begin()
	utils.CheckErr(err)
	defer dbSession.Close()
	weight := make(map[int]int64)
	weight[10] = 85
	weight[20] = 10
	weight[30] = 5
	for rarity := 10; rarity <= 30; rarity += 10 {
		for school := 0; school < 3; school++ {
			groupMasterID := rarity*10 + school
			cardMasterIDs := []int{}
			err := masterdata.Table("m_card").Where("card_rarity_type = ? AND member_m_id / 100 == ?", rarity, school).
				Cols("id").Find(&cardMasterIDs)
			utils.CheckErr(err)
			for _, cardMasterID := range cardMasterIDs {
				_, err := dbSession.Table("s_gacha_card").Insert(model.GachaCard{
					GroupMasterID: groupMasterID,
					CardMasterID:  cardMasterID,
				})
				utils.CheckErr(err)
			}
			dbSession.Table("s_gacha_group").Insert(model.GachaGroup{
				GroupMasterID: groupMasterID,
				GroupWeight:   weight[rarity],
			})
		}
	}

	// gacha guarantee: new card
	dbSession.Table("s_gacha_guarantee").Insert(model.GachaGuarantee{
		GachaGuaranteeMasterID: 0,
		GuaranteeHandler:       "guarantee_new_card",
		GuaranteeParams:        []string{},
	})
	// gacha guarantee: UR card
	dbSession.Table("s_gacha_guarantee").Insert(model.GachaGuarantee{
		GachaGuaranteeMasterID: 1,
		GuaranteeHandler:       "guarantee_card_in_set",
		GuaranteeParams:        []string{"card_rarity_type = 30"},
	})
	// gacha guarantee: SR+ card
	dbSession.Table("s_gacha_guarantee").Insert(model.GachaGuarantee{
		GachaGuaranteeMasterID: 2,
		GuaranteeHandler:       "guarantee_card_in_set",
		GuaranteeParams:        []string{"card_rarity_type >= 20"},
	})
	// gacha guarantee: festival / party card
	dbSession.Table("s_gacha_guarantee").Insert(model.GachaGuarantee{
		GachaGuaranteeMasterID: 3,
		GuaranteeHandler:       "guarantee_card_in_set",
		GuaranteeParams:        []string{"passive_skill_slot == 2"},
	})
	dbSession.Commit()
}

func InsertGacha(args []string) {
	if len(args) == 0 {
		fmt.Println("Invalid params:", args)
		return
	}
	// insert gacha from json format, with some exceptions.
	file := args[0]
	gachas := []model.Gacha{}
	gachaJsons := utils.ReadAllText(file)

	err := json.Unmarshal([]byte(gachaJsons), &gachas)
	utils.CheckErr(err)
	for pos, gacha := range gachas {
		for i, appeal := range gacha.GachaAppeals {
			appeal.GachaAppealMasterID = gacha.GachaMasterID*10 + i
			gacha.DbGachaAppeals = append(gacha.DbGachaAppeals, appeal.GachaAppealMasterID)
			_, err := serverdb.Engine.Table("s_gacha_appeal").Insert(appeal)
			utils.CheckErr(err)
		}
		for i, draw := range gacha.GachaDraws {
			draw.GachaMasterID = gacha.GachaMasterID
			gacha.DbGachaDraws = append(gacha.DbGachaDraws, draw.GachaDrawMasterID)
			gjson.Get(gachaJsons, fmt.Sprintf("%d.gacha_draws.%d.guarantees", pos, i)).ForEach(
				func(_, value gjson.Result) bool {
					draw.Guarantees = append(draw.Guarantees, int(value.Int()))
					return true
				})
			_, err := serverdb.Engine.Table("s_gacha_draw").Insert(draw)
			utils.CheckErr(err)
		}
		gjson.Get(gachaJsons, fmt.Sprintf("%d.gacha_groups", pos)).ForEach(
			func(_, value gjson.Result) bool {
				gacha.DbGachaGroups = append(gacha.DbGachaGroups, int(value.Int()))
				return true
			})

		_, err := serverdb.Engine.Table("s_gacha").Insert(gacha)
		utils.CheckErr(err)
	}
}

func Gacha(args []string) {
	if len(args) == 0 {
		fmt.Println("Invalid params:", args)
		return
	}
	switch args[0] {
	case "init":
		InitGacha(args[1:])
	case "insert":
		InsertGacha(args[1:])
	}
}
