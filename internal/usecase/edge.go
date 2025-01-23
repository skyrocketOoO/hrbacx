package usecase

type NebulaUsecase struct{}

func NewNebulaUsecase() *NebulaUsecase {
	return &NebulaUsecase{}
}

func (u *NebulaUsecase) AddLeader(leaderID string, roleID string) error {
	return nil
}

func (u *NebulaUsecase) AssignPermission(objectID, permissionType, roleID string) error {
	return nil
}

func (u *NebulaUsecase) AssignRole(userID, roleID string) error {
	return nil
}

func (u *NebulaUsecase) CheckPermission(userID, permissionType, objectID string) (
	ok bool, err error,
) {
	return false, nil
}

func (u *NebulaUsecase) ClearAll() error {
	return nil
}
