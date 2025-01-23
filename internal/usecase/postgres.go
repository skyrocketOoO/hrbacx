package usecase

import (
	"fmt"

	"github.com/skyrocketOoO/hrbacx/internal/global"
	"github.com/skyrocketOoO/hrbacx/internal/model"
	"gorm.io/gorm"
)

type PgUsecase struct {
	db *gorm.DB
}

func NewPgUsecase(db *gorm.DB) *PgUsecase {
	return &PgUsecase{db}
}

func (u *PgUsecase) AddLeader(leaderID string, roleID string) error {
	return u.db.Create(&model.Edge{
		FromV:    "role_" + leaderID,
		Relation: global.LearderOf,
		ToV:      "role_" + roleID,
	}).Error
}

func (u *PgUsecase) AssignPermission(objectID, permissionType, roleID string) error {
	return u.db.Create(&model.Edge{
		FromV:    "role_" + roleID,
		Relation: permissionType,
		ToV:      "obj_" + objectID,
	}).Error
}

func (u *PgUsecase) AssignRole(userID, roleID string) error {
	return u.db.Create(&model.Edge{
		FromV:    "user_" + userID,
		Relation: global.BelongsTo,
		ToV:      "role_" + roleID,
	}).Error
}

func (u *PgUsecase) CheckPermission(userID, permissionType, objectID string) (
	ok bool, err error,
) {
	fmt.Println("CheckPermission", userID, permissionType, objectID)
	var result bool
	query := `SELECT check_permission($1, $2, $3)`
	if err := u.db.Raw(query, userID, permissionType, objectID).Scan(&result).Error; err != nil {
		return false, err
	}
	return result, nil
}

func (u *PgUsecase) ClearAll() error {
	return u.db.Unscoped().Where("1=1").Delete(&model.Edge{}).Error
}
