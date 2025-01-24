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
		CREATE OR REPLACE FUNCTION check_permission(user_id TEXT, permission_type TEXT, object_id TEXT)
		RETURNS BOOLEAN AS $$
		DECLARE
		    queue TEXT[] :=ARRAY[]::TEXT[];
		    visited TEXT[] := '{}';
		    current TEXT;
		BEGIN
			queue := queue || ARRAY(
			    SELECT to_v
			    FROM "edges"
			    WHERE from_v = user_id
			      AND relation = 'belongs_to'
			);

		    WHILE array_length(queue, 1) > 0 LOOP
		        current := queue[1];
		        queue := queue[2:array_length(queue, 1)];

		        IF current = ANY(visited) THEN
		            CONTINUE;
		        END IF;

		        visited := array_append(visited, current);

		        PERFORM 1
				FROM "edges"
				WHERE from_v = current
				  AND relation = permission_type
				  AND to_v = object_id;

				IF FOUND THEN
				    RETURN TRUE;
				END IF;

				queue := queue || ARRAY(
				    SELECT to_v
				    FROM "edges"
				    WHERE from_v = current
				      AND relation = 'leader_of'
				      AND NOT (to_v = ANY(visited))
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
