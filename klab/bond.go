package klab

var (
	LoveLevelRequiredForLovePanel = []int{0, 1, 3, 5, 7, 10, 15, 20, 25, 30, 35, 40, 50, 60, 70, 80, 90, 100, 120, 140, 160, 180, 200, 220, 240, 260, 280, 300, 320, 340, 360, 380, 400, 9999, 9999}
)

func MaxLovePanelLevelFromLoveLevel(loveLevel int) int {
	for i, v := range LoveLevelRequiredForLovePanel {
		if v > loveLevel {
			return i - 1
		}
	}
	panic("love level too high")
}

// bond require to level up to l
func BondRequiredToLevelUp(l int) int {
	res := 30 * l
	if l > 2 {
		res += 10 * (l - 2)
	}
	if l > 6 {
		res += 10 * (l - 6)
	}
	if l > 20 {
		res += 10 * (l - 20)
	}
	if l > 59 {
		res += 10 * (l - 59)
	}
	if l > 153 {
		res += 80 * (l - 153)
	}
	return res
}

func BondRequiredTotal(l int) int {
	// O(1) because I can
	if l == 1 { // [1, 1]: 0
		return 0
	} else if l <= 2 { // [2, 2]: 30l at level l
		return 60 // also 15 * l * l
	} else if l <= 6 { // [3, 6]: 40l - 20 at level l
		return 20*l*l - 20
	} else if l <= 20 { // [7, 20]: 50l - 80 at level l
		return 25*l*l - 55*l + 130
	} else if l <= 59 { // [21, 59]: 60l - 280 at level l
		return 30*l*l - 250*l + 2030
	} else if l <= 153 { // [60, 153]: 70l - 870 at level l
		return 35*l*l - 835*l + 19140
	} else { // [154, inf): 150l - 13110 at level l
		return 75*l*l - 13035*l + 949380
	}

	// SELECT * FROM m_member_love_level WHERE love_level >= 3 AND love_level <= 6 AND 20 * love_level * love_level - 20 != love_point; -> 0
	// SELECT * FROM m_member_love_level WHERE love_level >= 7 AND love_level <= 20 AND 25 * love_level * love_level - 55 * love_level + 130 != love_point; -> 0
	// SELECT * FROM m_member_love_level WHERE love_level >= 21 AND love_level <= 59 AND 30 * love_level * love_level - 250 * love_level + 2030 != love_point; -> 0
	// SELECT * FROM m_member_love_level WHERE love_level >= 60 AND love_level <= 153 AND 35 * love_level * love_level - 835 * love_level + 19140 != love_point; -> 0
	// SELECT * FROM m_member_love_level WHERE love_level > 153 AND 75 * love_level * love_level - 13035 * love_level + 949380 != love_point; -> 0
}

func BondLevelFromBondValue(bond int) int {
	// return the level you would be at
	// highest value of x such that BondRequiredTotal(x) <= bond
	l := 0
	x := 1
	for ; BondRequiredTotal(x*2) <= bond; x *= 2 {
	}
	for ; x > 0; x /= 2 {
		if BondRequiredTotal(l+x) <= bond {
			l += x
		}
	}
	return l
}

func CenterBondGainBasedOnBondGain(gain int) int {
	// https://suyo.be/sifas/wiki/gameplay/live-rewards
	// this doesn't seem to be in DB, so it's probably server sided
	// the amount of bond and LP isn't consistent, because story stage give less bond
	// but the amount of center bond is consistent with the amount of bond per card gain
	// bond for double center is just bond for center / 2, without rounding
	// ignoring the story mode, then the extra center bond has a nice relationship with the lp used:
	// 4 = 0.4 * 10
	// 6 = 0.5 * 12
	// 9 = 0.6 * 15
	// 14 = 0.7 * 20
	if gain < 9 {
		return 0
	} else if gain == 9 {
		return 3
	} else if gain <= 13 {
		return 4
	} else if gain <= 16 {
		return 6
	} else if gain <= 24 {
		return 9
	} else {
		return 14
	}
}
