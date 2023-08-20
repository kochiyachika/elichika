package serverdb

import (
	"elichika/model"
	"elichika/utils"
	// "fmt"
)

func (session *Session) GetGachaList() []model.Gacha {
	gachaList := []model.Gacha{}
	err := Engine.Table("s_gacha").Find(&gachaList)
	utils.CheckErr(err)
	for i, _ := range gachaList {
		err = Engine.Table("s_gacha_appeal").In("gacha_appeal_master_id", gachaList[i].DbGachaAppeals).
			Find(&gachaList[i].GachaAppeals)
		utils.CheckErr(err)
		err = Engine.Table("s_gacha_draw").In("gacha_draw_master_id", gachaList[i].DbGachaDraws).
			Find(&gachaList[i].GachaDraws)
		utils.CheckErr(err)
	}
	return gachaList
}

func (session *Session) GetGacha(gachaMasterID int) model.Gacha {
	gacha := model.Gacha{}
	exists, err := Engine.Table("s_gacha").Where("gacha_master_id = ?", gachaMasterID).Get(&gacha)
	utils.CheckErrMustExist(err, exists)
	// work on the state of the gacha if
	err = Engine.Table("s_gacha_appeal").In("gacha_appeal_master_id", gacha.DbGachaAppeals).
		Find(&gacha.GachaAppeals)
	utils.CheckErr(err)
	err = Engine.Table("s_gacha_draw").In("gacha_draw_master_id", gacha.DbGachaDraws).
		Find(&gacha.GachaDraws)
	utils.CheckErr(err)
	return gacha
}
