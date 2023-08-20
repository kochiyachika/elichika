package serverdb

import (
	"elichika/generic"
	"elichika/model"
	"elichika/utils"
)

func (session *Session) GetUserPlayList() generic.ObjectByObjectIDWrite[*model.UserPlayListItem] {
	playlist := generic.ObjectByObjectIDWrite[*model.UserPlayListItem]{}
	err := Engine.Table("s_user_play_list").Where("user_id = ?", session.UserStatus.UserID).
		Find(&playlist.Objects)
	utils.CheckErr(err)
	return playlist
}

func (session *Session) UpdateUserPlayList(item model.UserPlayListItem) {
	exists, err := Engine.Table("s_user_play_list").
		Where("user_id = ? AND user_play_list_id = ?", session.UserStatus.UserID, item.UserPlayListID).
		AllCols().Update(&item)
	utils.CheckErr(err)
	if exists == 0 {
		exists, err := Engine.Table("s_user_play_list").Insert(item)
		utils.CheckErrMustExist(err, exists != 0)
	}

}
