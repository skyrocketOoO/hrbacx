package usecase

import (
	"github.com/skyrocketOoO/hrbacx/internal/global"
	"github.com/skyrocketOoO/hrbacx/internal/model"
)

func (u *Usecase) AddLeader(leaderID string, roleID string) error {
	return u.db.Create(&model.Edge{
		From:     leaderID,
		Relation: global.LearderOf,
		To:       roleID,
	}).Error
}

func (u *Usecase) AssignPermission(objectID, permissionType, roleID string) error {
	return u.db.Create(&model.Edge{
		From:     roleID,
		Relation: permissionType,
		To:       objectID,
	}).Error
}

func (u *Usecase) AssignRole(userID, roleID string) error {
	return u.db.Create(&model.Edge{
		From:     userID,
		Relation: global.BelongsTo,
		To:       roleID,
	}).Error
}

func (u *Usecase) CheckPermission(userID, permissionType, objectID string) (ok bool, err error) {
	var count int64
	err = u.db.Model(&model.Edge{}).
		Where("from = ? AND relation = ? AND to = ?", userID, permissionType, objectID).
		Count(&count).Error
	return count > 0, err
}
