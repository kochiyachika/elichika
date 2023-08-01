package model

// UserCommunicationMemberDetailBadgeByID ...
type UserCommunicationMemberDetailBadgeByID struct {
	MemberMasterID     int  `json:"member_master_id"`
	IsStoryMemberBadge bool `json:"is_story_member_badge"`
	IsStorySideBadge   bool `json:"is_story_side_badge"`
	IsVoiceBadge       bool `json:"is_voice_badge"`
	IsThemeBadge       bool `json:"is_theme_badge"`
	IsCardBadge        bool `json:"is_card_badge"`
	IsMusicBadge       bool `json:"is_music_badge"`
}

// UserMemberInfo ...
type UserMemberInfo struct {
	UserID                   int  `xorm:"pk 'user_id'" json:"-"`
	MemberMasterID           int  `xorm:"pk 'member_master_id'" json:"member_master_id"`
	CustomBackgroundMasterID int  `xorm:"'custom_background_master_id'" json:"custom_background_master_id"`
	SuitMasterID             int  `xorm:"'suit_master_id'" json:"suit_master_id"`
	LovePoint                int  `json:"love_point"`
	LovePointLimit           int  `json:"love_point_limit"`
	LoveLevel                int  `json:"love_level"`
	ViewStatus               int  `json:"view_status"`
	IsNew                    bool `json:"is_new"`
}

// Bond board tile
type UserMemberLovePanel struct {
	UserID                int `xorm:"pk 'user_id'"`
	MemberID              int `xorm:"'member_id'"`
	MemberLovePanelCellID int `xorm:"pk 'member_love_panel_cell_id'"` // level * 10000 + tile_id * 1000 + member_id
	// can store user_id, member_id, level * 32 + 5bits instead, but that's not necessary for now
}

// SuitInfo ...
type SuitInfo struct {
	SuitMasterID int  `json:"suit_master_id"`
	IsNew        bool `json:"is_new"`
}
