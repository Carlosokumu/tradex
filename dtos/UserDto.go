package dtos

type UserDto struct {
	UserName string `form:"username" binding:"required"`
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type UserLoginDto struct {
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}
