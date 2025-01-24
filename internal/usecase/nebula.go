package usecase

import (
	"fmt"

	nebulaservice "github.com/skyrocketOoO/hrbacx/internal/service/exter/nebula"
	nebula "github.com/vesoft-inc/nebula-go/v3"
)

type NebulaUsecase struct {
	SessionPool *nebula.SessionPool
}

func NewNebulaUsecase(sessionPool *nebula.SessionPool) *NebulaUsecase {
	return &NebulaUsecase{
		SessionPool: sessionPool,
	}
}

func (u *NebulaUsecase) AddLeader(leaderID string, roleID string) error {
	if err := u.addVertex("role()", "role_"+leaderID); err != nil {
		return err
	}
	if err := u.addVertex("role()", "role_"+roleID); err != nil {
		return err
	}

	sql := fmt.Sprintf(`INSERT EDGE IF NOT EXISTS leader_of() VALUES "%s"->"%s":();`,
		"role_"+leaderID,
		"role_"+roleID)

	resp, err := u.SessionPool.Execute(sql)
	if err != nil {
		return fmt.Errorf("query execution failed: %v", err)
	}

	if !resp.IsSucceed() {
		return fmt.Errorf("query failed: %s", resp.GetErrorMsg())
	}
	return nil
}

func (u *NebulaUsecase) AssignPermission(objectID, permissionType, roleID string) error {
	if err := u.addVertex("role()", "role_"+roleID); err != nil {
		return err
	}
	if err := u.addVertex("obj()", "obj_"+objectID); err != nil {
		return err
	}
	sql := fmt.Sprintf(
		`INSERT EDGE IF NOT EXISTS has_permission(type) VALUES "%s"->"%s":("%s");`,
		"role_"+roleID, "obj_"+objectID, permissionType,
	)

	resp, err := u.SessionPool.Execute(sql)
	if err != nil {
		return fmt.Errorf("query execution failed: %v", err)
	}

	if !resp.IsSucceed() {
		return fmt.Errorf("query failed: %s", resp.GetErrorMsg())
	}
	return nil
}

func (u *NebulaUsecase) AssignRole(userID, roleID string) error {
	if err := u.addVertex("user()", "user_"+userID); err != nil {
		return err
	}
	if err := u.addVertex("role()", "role_"+roleID); err != nil {
		return err
	}
	sql := fmt.Sprintf(`INSERT EDGE IF NOT EXISTS belongs_to() VALUES "%s"->"%s":();`,
		"user_"+userID,
		"role_"+roleID,
	)

	resp, err := u.SessionPool.Execute(sql)
	if err != nil {
		return fmt.Errorf("query execution failed: %v", err)
	}

	if !resp.IsSucceed() {
		return fmt.Errorf("query failed: %s", resp.GetErrorMsg())
	}
	return nil
}

func (u *NebulaUsecase) CheckPermission(userID, permissionType, objectID string) (
	ok bool, err error,
) {
	sql := fmt.Sprintf(
		`MATCH (v)-[belongs_to]->(:role)-[leader_of*0..]->(:role)`+
			`-[has_permission{type:'%s'}]->(d:obj) WHERE id(v) == '%s' AND id(d) == '%s' RETURN d LIMIT 1;`,
		permissionType,
		"user_"+userID,
		"obj_"+objectID,
	)

	resp, err := u.SessionPool.Execute(sql)
	if err != nil {
		return false, fmt.Errorf("query execution failed: %v", err)
	}

	if !resp.IsSucceed() {
		return false, fmt.Errorf("query failed: %s", resp.GetErrorMsg())
	}

	if resp.GetRowSize() > 0 {
		return true, nil
	}

	return false, nil
}

func (u *NebulaUsecase) ClearAll() error {
	sql := fmt.Sprintf(`DROP SPACE IF EXISTS %s;`, nebulaservice.SPACE)
	resp, err := u.SessionPool.Execute(sql)
	if err != nil {
		return fmt.Errorf("query execution failed: %v", err)
	}

	if !resp.IsSucceed() {
		return fmt.Errorf("query failed: %s", resp.GetErrorMsg())
	}
	return nil
}

func (u *NebulaUsecase) addVertex(tag, id string) error {
	resp, err := u.SessionPool.Execute(
		fmt.Sprintf(`INSERT VERTEX IF NOT EXISTS %s VALUES "%s":();`, tag, id))
	if err != nil {
		return err
	}
	if !resp.IsSucceed() {
		return fmt.Errorf("leader vertex insertion failed: %s", resp.GetErrorMsg())
	}
	return nil
}
