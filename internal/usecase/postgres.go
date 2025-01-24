package usecase

import (
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
	// fmt.Println("CheckPermission", userID, permissionType, objectID)
	var result bool
	query := `SELECT check_permission($1, $2, $3)`
	if err := u.db.Raw(query, "user_"+userID, permissionType, "obj_"+objectID).Scan(&result).Error; err != nil {
		return false, err
	}
	return result, nil

	/*
		CREATE EXTENSION IF NOT EXISTS hstore;

		CREATE OR REPLACE FUNCTION check_permission(user_id TEXT, permission_type TEXT, object_id TEXT)
		RETURNS BOOLEAN AS $$
		DECLARE
		    queue TEXT[] := ARRAY[]::TEXT[];
		    visited hstore := hstore('');  -- Use hstore to track visited nodes
		    current TEXT;
		BEGIN
		    -- Initialize queue with nodes from 'belongs_to' relation
		    queue := queue || ARRAY(
		        SELECT to_v
		        FROM "edges"
		        WHERE from_v = user_id
		          AND relation = 'belongs_to'
		    );

		    -- BFS traversal
		    WHILE array_length(queue, 1) > 0 LOOP
		        -- Dequeue the first element
		        current := queue[1];
		        queue := queue[2:array_length(queue, 1)];

		        -- Skip if already visited
		        IF visited -> current IS NOT NULL THEN
		            CONTINUE;
		        END IF;

		        -- Mark the node as visited
		        visited := visited || hstore(current, 'visited');

		        -- Check if the permission exists
		        PERFORM 1
		        FROM "edges"
		        WHERE from_v = current
		          AND relation = permission_type
		          AND to_v = object_id;

		        IF FOUND THEN
		            RETURN TRUE;
		        END IF;

		        -- Enqueue neighbors (leader_of relation) if not visited
		        queue := queue || ARRAY(
		            SELECT to_v
		            FROM "edges"
		            WHERE from_v = current
		              AND relation = 'leader_of'
		              AND visited -> to_v IS NULL  -- Only enqueue if not visited
		        );
		    END LOOP;

		    RETURN FALSE;
		END;
		$$ LANGUAGE plpgsql;

	*/
}

func (u *PgUsecase) ClearAll() error {
	return u.db.Unscoped().Where("1=1").Delete(&model.Edge{}).Error
}
