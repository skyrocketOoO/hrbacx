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

	query := `
		WITH RECURSIVE permission_cte AS (
		    -- Base case: Start with the user's belongs_to relations
		    SELECT to_v AS node
		    FROM edges
		    WHERE from_v = ? AND relation = 'belongs_to'

		    UNION ALL

		    -- Recursive case: Find neighbors through leader_of relation
		    SELECT e.to_v AS node
		    FROM edges e
		    INNER JOIN permission_cte pc ON e.from_v = pc.node
		    WHERE e.relation = 'leader_of'
		)
		-- Check if any node in the CTE has the required permission
		SELECT 1
		FROM edges
		WHERE from_v IN (SELECT node FROM permission_cte)
		  AND relation = ?
		  AND to_v = ?
		LIMIT 1;
	`

	var result int
	if err := u.db.Raw(query, userID, permissionType, objectID).Scan(&result).Error; err != nil {
		return false, fmt.Errorf("failed to execute permission query: %v", err)
	}

	return result > 0, nil
}

func (u *MysqlUsecase) ClearAll() error {
	return u.db.Unscoped().Where("1=1").Delete(&model.Edge{}).Error
}
