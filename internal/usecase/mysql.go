package usecase

import (
	"fmt"

	"github.com/skyrocketOoO/hrbacx/internal/global"
	"github.com/skyrocketOoO/hrbacx/internal/model"
	"gorm.io/gorm"
)

type MysqlUsecase struct {
	db *gorm.DB
}

func NewMysqlUsecase(db *gorm.DB) *MysqlUsecase {
	return &MysqlUsecase{db}
}

func (u *MysqlUsecase) AddLeader(leaderID string, roleID string) error {
	return u.db.Create(&model.Edge{
		FromV:    "role_" + leaderID,
		Relation: global.LearderOf,
		ToV:      "role_" + roleID,
	}).Error
}

func (u *MysqlUsecase) AssignPermission(objectID, permissionType, roleID string) error {
	return u.db.Create(&model.Edge{
		FromV:    "role_" + roleID,
		Relation: permissionType,
		ToV:      "obj_" + objectID,
	}).Error
}

func (u *MysqlUsecase) AssignRole(userID, roleID string) error {
	return u.db.Create(&model.Edge{
		FromV:    "user_" + userID,
		Relation: global.BelongsTo,
		ToV:      "role_" + roleID,
	}).Error
}

func (u *MysqlUsecase) CheckPermission(userID, permissionType, objectID string) (bool, error) {
	userID = "user_" + userID
	objectID = "obj_" + objectID
	// Initialize queue and visited map
	queue := []string{}
	visited := make(map[string]bool)
	var current string

	var edges []model.Edge
	if err := u.db.
		Where("from_v = ? AND relation = ?", userID, "belongs_to").
		Find(&edges).Error; err != nil {
		return false, fmt.Errorf("failed to query belongs_to relations: %v", err)
	}

	for _, edge := range edges {
		queue = append(queue, edge.ToV)
	}

	for len(queue) > 0 {
		current = queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}

		// Mark the node as visited
		visited[current] = true

		// 2. Check if the permission exists for this node
		var permissionExists int64
		if err := u.db.Model(&model.Edge{}).
			Where("from_v = ? AND relation = ? AND to_v = ?", current, permissionType, objectID).
			Limit(1).
			Count(&permissionExists).Error; err != nil {
			return false, fmt.Errorf("failed to check permission: %v", err)
		}

		if permissionExists > 0 {
			return true, nil
		}

		// 3. Enqueue neighbors with 'leader_of' relation that haven't been visited
		if err := u.db.Where("from_v = ? AND relation = ?", current, "leader_of").Find(&edges).Error; err != nil {
			return false, fmt.Errorf("failed to query leader_of relations: %v", err)
		}

		for _, edge := range edges {
			if !visited[edge.ToV] {
				queue = append(queue, edge.ToV)
			}
		}
	}

	// If no permission found, return false
	return false, nil
}

func (u *MysqlUsecase) ClearAll() error {
	return u.db.Unscoped().Where("1=1").Delete(&model.Edge{}).Error
}
