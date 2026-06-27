package dal

import (
	"fmt"
	"go-web-demo/internal/common/response"
	"go-web-demo/internal/dto"
	"go-web-demo/internal/model"
	"strings"

	"gorm.io/gorm"
)

type UserDao struct{}

var userDao *UserDao

func init() {
	userDao = &UserDao{}
}

func GetUserDao() *UserDao {
	return userDao
}

// Insert 插入一条记录
func (*UserDao) Insert(db *gorm.DB, user *model.User) error {
	result := db.Create(user)
	return result.Error
}

// InsertBatch 批量插入
func (*UserDao) InsertBatch(db *gorm.DB, userList []*model.User) error {
	result := db.Create(userList)
	return result.Error
}

// SelectById 根据 ID 查询记录
func (*UserDao) SelectById(db *gorm.DB, id int64) (*model.User, error) {
	var user model.User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdateSelective 选择性更新（只更新非零值字段）
func (*UserDao) UpdateSelective(db *gorm.DB, user *model.User, updateFields ...string) error {
	// 如果指定了字段，则只更新这些字段
	if len(updateFields) > 0 {
		result := db.Select(updateFields).Updates(user)
		return result.Error
	}
	// 否则，只更新非零值字段
	result := db.Updates(user)
	return result.Error
}

// Update 更新（支持空值更新）
func (*UserDao) Update(db *gorm.DB, id int64, userMap map[string]interface{}) error {
	result := db.Table("user").Where("id = ?", id).Updates(userMap)
	return result.Error
}

// Delete 根据 ID 删除记录（物理删除）
func (*UserDao) Delete(db *gorm.DB, id int64) error {
	result := db.Unscoped().Delete(&model.User{}, id)
	return result.Error
}

// DeleteBatch 批量删除用户
func (*UserDao) DeleteBatch(db *gorm.DB, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	result := db.Unscoped().Where("id IN ?", ids).Delete(&model.User{})
	return result.Error
}

// DeleteSoft 软删除（设置 deleted_at 字段）
func (*UserDao) DeleteSoft(db *gorm.DB, id int64) error {
	result := db.Delete(&model.User{}, id)
	return result.Error
}

// SelectList 列表查询
func (*UserDao) SelectList(db *gorm.DB, query *dto.UserPageQueryRequest) ([]model.User, error) {
	var userList []model.User
	baseDB := buildWhereByUser(db, query)
	if err := baseDB.Find(&userList).Error; err != nil {
		return nil, err
	}
	return userList, nil
}

// SelectListByPage 查询分页列表
func (*UserDao) SelectListByPage(db *gorm.DB, query *dto.UserPageQueryRequest) (*response.PageResult[model.User], error) {
	var userList []model.User
	var total int64

	// 构建基础查询（带所有 WHERE 条件）
	baseDB := buildWhereByUser(db, query)

	// 第一步：获取总数（COUNT）
	if err := baseDB.Count(&total).Error; err != nil {
		return nil, err
	}

	// 第二步：获取分页数据（LIST）
	page := query.Page
	if page == 0 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 构建排序
	order := buildOrderByUser(query)

	if err := baseDB.Offset(int(offset)).Limit(int(pageSize)).Order(order).Find(&userList).Error; err != nil {
		return nil, err
	}

	// 计算总页数
	totalPages := int32((total + int64(pageSize) - 1) / int64(pageSize))

	return &response.PageResult[model.User]{
		List:       userList,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// buildOrderByUser 自定义排序
func buildOrderByUser(req *dto.UserPageQueryRequest) string {
	// 1. 定义白名单：只允许这些字段被排序。  键是前端传来的字段名，值是数据库中真实的列名（可以做映射，也可以直接校验）
	allowedSortFields := map[string]bool{
		"id":         true,
		"created_at": true,
		"updated_at": true,
		// 如果有其他需要排序的字段，在这里添加
	}

	// 默认排序
	order := "id DESC"

	if req.SortField != nil {
		field := strings.TrimSpace(*req.SortField)

		// 2. 核心校验：检查字段是否在白名单中
		// 注意：这里建议转为小写比较，防止大小写绕过 (如 "Id", "ID")
		if allowedSortFields[strings.ToLower(field)] {
			direction := "ASC"
			if req.SortOrder != nil && strings.ToLower(strings.TrimSpace(*req.SortOrder)) == "desc" {
				direction = "DESC"
			}
			// 安全拼接：因为 field 已经过白名单校验，direction 也是固定的枚举值
			order = fmt.Sprintf("%s %s", field, direction)
		} else {
			// 可选：记录日志，表示有人尝试非法排序
			// log.Warnf("Invalid sort field attempted: %s", field)

			// 策略 A: 忽略非法字段，使用默认排序 (推荐)
			order = "id DESC"

			// 策略 B: 直接返回错误 (更严格)
			// return "" // 并在调用处处理 error
		}
	}

	return order
}

// buildWhere 构建where查询
func buildWhereByUser(db *gorm.DB, req *dto.UserPageQueryRequest) *gorm.DB {
	db = db.Table("user")
	if req.Nickname != nil {
		db = db.Where("nickname = ?", *req.Nickname)
	}
	if req.Email != nil {
		db = db.Where("email = ?", *req.Email)
	}
	if req.Phone != nil {
		db = db.Where("phone = ?", *req.Phone)
	}
	if req.Gender != nil {
		db = db.Where("gender = ?", *req.Gender)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	db = db.Where("is_deleted = 0")
	return db
}
