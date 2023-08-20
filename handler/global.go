package handler

import (
	"elichika/config"
	"elichika/encrypt"
	"elichika/model"
	"elichika/utils"

	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"xorm.io/xorm"
)

var (
	IsGlobal      = false
	MasterVersion = "b66ec2295e9a00aa"
	StartUpKey    = "5f7IZY1QrAX0D49g"

	MainEng *xorm.Engine

	presetDataPath = "assets/preset/"
	userDataPath   = "assets/userdata/"
)

func init() {
	MainEng = config.MainEng

	os.Mkdir(userDataPath, 0755)
}

func SignResp(ep, body, key string) (resp string) {
	signBody := fmt.Sprintf("%d,\"%s\",0,%s", time.Now().UnixMilli(), MasterVersion, body)
	sign := encrypt.HMAC_SHA1_Encrypt([]byte(ep+" "+signBody), []byte(key))
	// fmt.Println(sign)

	resp = fmt.Sprintf("[%s,\"%s\"]", signBody, sign)
	return
}

func GetData(fileName string) string {
	presetDataFile := presetDataPath + fileName
	if !utils.PathExists(presetDataFile) {
		panic("File not exists")
	}

	return utils.ReadAllText(presetDataFile)
}

func GetUserData(fileName string) string {
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

func GetUserAccessoryData() string {
	if IsGlobal {
		return GetData("userAccessory_gl.json")
	}
	return GetData("userAccessory.json")
}

func GetPartyInfoByRoleIds(roleIds []int) (partyIcon int, partyName string) {
	// 脑残逻辑部分
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == j {
				continue
			}
			exists, err := MainEng.Table("m_live_party_name").
				Where("role_1 = ? AND role_2 = ? AND role_3 = ?", roleIds[i], roleIds[j], roleIds[3-i-j]).
				Cols("name,live_party_icon").Get(&partyName, &partyIcon)
			utils.CheckErr(err)
			if exists {
				return
			}
		}
	}
	panic("not found")
	return
}

func GetRealPartyName(partyName string) (realPartyName string) {
	_, err := MainEng.Table("m_dictionary").Where("id = ?", strings.ReplaceAll(partyName, "k.", "")).
		Cols("message").Get(&realPartyName)
	utils.CheckErr(err)
	return
}

func GetMemberMasterIdByCardMasterId(cardMasterId int) (memberMasterId int) {
	_, err := MainEng.Table("m_card").Where("id = ?", cardMasterId).
		Cols("member_m_id").Get(&memberMasterId)
	utils.CheckErr(err)
	return
}

func GetMemberDefaultSuitByCardMasterId(cardMasterId int) int {
	suitMasterId, err := strconv.Atoi(fmt.Sprintf("10%03d1001", GetMemberMasterIdByCardMasterId(cardMasterId)))
	utils.CheckErr(err)

	return suitMasterId
}

func GetMemberInfoByCardMasterId(cardMasterId int) (memberInfo model.UserMemberInfo) {
	key := fmt.Sprintf("user_member_by_member_id.#(member_master_id==%d)", GetMemberMasterIdByCardMasterId(cardMasterId))
	memberData := GetUserData("memberSettings.json")
	if err := json.Unmarshal([]byte(gjson.Parse(memberData).Get(key).String()), &memberInfo); err != nil {
		panic(err)
	}
	return
}

func GetMemberInfo(memberMasterId int) (memberInfo model.UserMemberInfo) {
	key := fmt.Sprintf("user_member_by_member_id.#(member_master_id==%d)", memberMasterId)
	memberData := GetUserData("memberSettings.json")
	if err := json.Unmarshal([]byte(gjson.Parse(memberData).Get(key).String()), &memberInfo); err != nil {
		panic(err)
	}
	return
}
