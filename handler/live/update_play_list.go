package live

import (
	"elichika/config"
	"elichika/handler"
	"elichika/model"
	"elichika/serverdb"
	"elichika/utils"

	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	// "github.com/tidwall/sjson"
	// "xorm.io/xorm"
)

func LiveUpdatePlayList(ctx *gin.Context) {
	reqBody := gjson.Parse(ctx.GetString("reqBody")).Array()[0].String()
	// fmt.Println(reqBody)
	req := model.LiveUpdatePlayListReq{}
	err := json.Unmarshal([]byte(reqBody), &req)
	utils.CheckErr(err)

	userID := ctx.GetInt("user_id")
	session := serverdb.GetSession(ctx, userID)
	mul := 0
	if req.IsSet {
		mul = 1
	}
	session.UpdateUserPlayList(model.UserPlayListItem{
		UserID:         userID,
		UserPlayListID: req.GroupNum + req.LiveMasterID*10,
		GroupNum:       req.GroupNum * mul,
		LiveID:         req.LiveMasterID * mul})

	// not sure what this response is, but "" works
	// this mean that the client is not checking for the response
	// and thus it doesn't send another request when we click on the button again
	// it will also not remove the song when we turn off, the filter will also not work instantly
	// any action that fetch user status will update this

	// not sure if the client can even update after this point, the following and some other like it doesn't seem to work:
	// signBody, _ := json.Marshal(session.GetUserPlayList())
	// maybe if someone can look into how the game handle this function, we can know for sure
	signBody := ""
	resp := handler.SignResp(ctx.GetString("ep"), string(signBody), config.SessionKey)
	ctx.Header("Content-Type", "application/json")
	ctx.String(http.StatusOK, resp)
	// fmt.Println(resp)
}
