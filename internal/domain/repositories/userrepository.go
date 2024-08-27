package repositories

import "weex_admin/internal/domain/entities"

type GetUserRequest struct {
	UserID   *uint64
	UserName *string

	JoinRole *bool
}

type GetUserOption = TransactionOption
type CreateUserOption = TransactionOption

type UserRepository interface {
	GetUser(req *GetUserRequest, option *GetUserOption) (*entities.User, error)
	GetUsers(reqs []*GetUserRequest, option *GetUserOption) ([]*entities.User, error)
	CreateUser(user *entities.User, option *CreateUserOption) (*entities.User, error)
	CreateUsers(users []*entities.User, option *CreateUserOption) ([]*entities.User, error)
}
