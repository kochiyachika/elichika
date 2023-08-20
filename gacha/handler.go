package gacha

import (
	"elichika/model"
	"elichika/serverdb"
	"elichika/utils"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

var (
	// handlers take the context
	// the context have the session and the list of available cards in the gacha builtin to it
	// handlers should choose a random card from the card pool if it match, or return 0 if there is no card that fit
	GuaranteeHandlers map[string]func(*gin.Context, []string) int
)

// take no params, return a new card
func GuaranteeNewCard(ctx *gin.Context, params []string) int {
	cardPool := ctx.MustGet("gacha_card_pool").([]model.GachaCard)
	userCards := []int{}
	err := serverdb.Engine.Table("s_user_card").Where("user_id = ?", ctx.GetInt("user_id")).Cols("card_master_id").
		Find(&userCards)
	utils.CheckErr(err)
	cardSet := map[int]bool{}
	for _, id := range userCards {
		cardSet[id] = true
	}
	newCards := []model.GachaCard{}
	for _, card := range cardPool {
		_, have := cardSet[card.CardMasterID]
		if !have {
			newCards = append(newCards, card)
		}
	}
	return ChooseRandomCard(newCards)
}

// card_rarity_type = 10 for R
// card_rarity_type = 20 for SR
// card_rarity_type = 30 for UR
func GuaranteeCardInSet(ctx *gin.Context, params []string) int {
	if len(params) == 0 {
		return 0
	}
	cardPool := ctx.MustGet("gacha_card_pool").([]model.GachaCard)
	// use SQL to get card from the set
	db := ctx.MustGet("masterdata.db").(*xorm.Engine)
	cardIDs := []int{}
	err := db.Table("m_card").Where(params[0]).Cols("id").Find(&cardIDs)
	utils.CheckErr(err)
	cardSet := map[int]bool{}
	for _, id := range cardIDs {
		cardSet[id] = true
	}
	availableCards := []model.GachaCard{}
	for _, card := range cardPool {
		_, have := cardSet[card.CardMasterID]
		if have {
			availableCards = append(availableCards, card)
		}
	}

	return ChooseRandomCard(availableCards)
}

func init() {
	GuaranteeHandlers = make(map[string]func(*gin.Context, []string) int)
	GuaranteeHandlers["guarantee_new_card"] = GuaranteeNewCard
	GuaranteeHandlers["guarantee_card_in_set"] = GuaranteeCardInSet
}
