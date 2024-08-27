package validators

import "weex_admin/internal/domain/valueobjects"

type CheckPermissionRequest struct {
	UserID   *uint64
	UserName *string
	Policy   valueobjects.Policy

	Tx interface{}
}

type AuthValidator interface {
	CheckPermission(req *CheckPermissionRequest) (bool, error)
}
