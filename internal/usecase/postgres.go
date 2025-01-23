package usecase

import (
	"github.com/skyrocketOoO/gox/Collection/queue"
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

func (u *PgUsecase) CheckPermission(userID, permissionType, objectID string) (ok bool, err error) {
	q := queue.NewQueue[string]()
	q.Push("user_" + userID)

	for !q.IsEmpty() {
		n := q.Len()
		for i := 0; i < n; i++ {
			v, _ := q.Pop()
			edges := []model.Edge{}
			if err := u.db.Debug().Where(`from_v = ?`, v).Find(&edges).Error; err != nil {
				return false, err
			}

			for _, edge := range edges {
				switch edge.Relation {
				case global.BelongsTo:
					q.Push(edge.ToV)
				case global.LearderOf:
					q.Push(edge.ToV)
				case permissionType:
					if edge.ToV == "obj_"+objectID {
						return true, nil
					}
				}
			}
		}
	}

	return false, nil
}
