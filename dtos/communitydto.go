package dtos

type CommunityDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Creator     string `json:"creator" binding:"required"`
}

type JoinCommunityDto struct {
	Username      string `json:"username" binding:"required"`
	CommunityName string `json:"community_name" binding:"required"`
}

type PostToCommunityDto struct {
	Username      string `form:"username" binding:"required"`
	CommunityName string `form:"community_name" binding:"required"`
	Content       string `form:"content" binding:"required,max=1000"`
}
