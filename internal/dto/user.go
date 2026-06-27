package dto

type UserPageQueryRequest struct {
	Id       *uint64 `json:"id"`        // 用户ID
	Username *string `json:"username"`  // 用户名
	Nickname *string `json:"nickname"`  // 昵称
	Email    *string `json:"email"`     // 邮箱
	Phone    *string `json:"phone"`     // 手机号
	Gender   *uint8  `json:"gender"`    // 性别：0-未知，1-男，2-女
	Status   *uint8  `json:"status"`    // 状态：0-禁用，1-启用
	Page     int32   `json:"page"`      // 页码
	PageSize int32   `json:"page_size"` // 每页数量

	SortField *string `json:"sort_field"` // 排序字段
	SortOrder *string `json:"sort_order"` // 排序顺序：asc 或 desc
}

type UserCreateRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Gender   uint8  `json:"gender"`
	Status   uint8  `json:"status"`
}

type UserUpdateRequest struct {
	ID       uint   `json:"id" binding:"required"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Gender   uint8  `json:"gender"`
	Status   uint8  `json:"status"`
}

type UserDeleteRequest struct {
	ID uint `json:"id" binding:"required"`
}