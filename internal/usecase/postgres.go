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

	/*
		CREATE OR REPLACE FUNCTION check_permission(user_id TEXT, permission_type TEXT, object_id TEXT)
		RETURNS BOOLEAN AS $$
		DECLARE
		    queue TEXT[];            -- Queue for BFS traversal
		    current TEXT;            -- Current node being processed
		    edge RECORD;             -- Holds edges fetched during traversal
		    permission_found BOOLEAN := FALSE; -- Flag to indicate if permission is found
		BEGIN
		    queue := ARRAY['user_' || user_id];

		    WHILE array_length(queue, 1) > 0 LOOP
		        current := queue[1];
		        queue := queue[2:array_length(queue, 1)];

		        SELECT TRUE
		        INTO permission_found
		        FROM "edges"
		        WHERE from_v = current
		          AND relation = permission_type
		          AND to_v = 'obj_' || object_id;

		        IF permission_found THEN
		            RETURN TRUE;
		        END IF;

		        FOR edge IN
		            SELECT to_v
		            FROM "edges"
		            WHERE from_v = current
		              AND (relation = 'leader_of' or relation = 'belongs_to')

		        LOOP
		            queue := array_append(queue, edge.to_v);
		        END LOOP;
		    END LOOP;

		    RETURN FALSE;
		END;
		$$ LANGUAGE plpgsql;
	*/
}

func (u *PgUsecase) ClearAll() error {
	return u.db.Unscoped().Where("1=1").Delete(&model.Edge{}).Error
}
