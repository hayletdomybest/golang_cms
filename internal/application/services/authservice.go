package services

import (
	"weex_admin/internal/domain/entities"
	"weex_admin/internal/domain/repositories"
	"weex_admin/internal/domain/validators"
	"weex_admin/internal/domain/valueobjects"
	"weex_admin/internal/infrastructure/jwt"
	"weex_admin/internal/shared/crypto"
	"weex_admin/internal/shared/errors"
	"weex_admin/internal/shared/pointer"
)

type CreateUserRequest struct {
	Name     string
	Email    string
	Password string
	RoleID   *uint
	RoleName *string
	Tx       interface{}
}

type CreateRoleRequest struct {
	RoleName string
	Policies []valueobjects.Policy
	Tx       interface{}
}

func (s *AuthService) CheckPermission(req *validators.CheckPermissionRequest) (bool, error) {
	pass, err := s.authReposiotry.CheckPermission(&repositories.CheckPermissionRequest{
		UserID:   req.UserID,
		UserName: req.UserName,
		Policy:   req.Policy,
	}, &repositories.CheckPermissionOption{
		Tx: req.Tx,
	})

	if err != nil {
		return false, err
	}
	return pass, nil
}

func (s *AuthService) Login(name, password string) (string, error) {
	user, err := s.userRepository.GetUser(&repositories.GetUserRequest{
		UserName: &name,
		JoinRole: pointer.Bool(true),
	}, nil)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.ErrNotFound.
			WithMessagef("user %s not found", name)
	}

	if !crypto.CompareHashAndPassword(user.Password, password) {
		return "", errors.ErrUnauthorized.
			WithMessagef("invalid username or password")
	}
	return s.jWTManager.GenerateToken(&jwt.ClaimsPayload{
		UserID:   user.ID,
		UserName: user.Name,
		RoleName: user.Role.Name,
	})
}

func (s *AuthService) CreateUser(
	req *CreateUserRequest) (*entities.User, error) {
	var err error
	tx := req.Tx
	shouldCommit := false

	rollback := func() {
		if shouldCommit {
			_ = s.dbTransactioner.Rollback(tx)
		}
	}

	commit := func() {
		if shouldCommit {
			_ = s.dbTransactioner.Commit(tx)
		}
	}

	if tx == nil {
		shouldCommit = true
		tx, err = s.dbTransactioner.Begin()
		if err != nil {
			return nil, err
		}
	}

	hash, err := crypto.HashPassword(req.Password)
	if err != nil {
		rollback()
		return nil, err
	}

	role, err := s.roleRepository.GetRole(&repositories.GetRoleRequest{
		RoleID:   req.RoleID,
		RoleName: req.RoleName,
	}, &repositories.GetRoleOption{
		Tx: tx,
	})
	if err != nil {
		rollback()
		return nil, err
	}

	if role == nil {
		rollback()
		return nil, errors.ErrNotFound.
			WithMessagef("role id:%d role name:%s not found", *req.RoleID, *req.RoleName)
	}

	user, err := s.userRepository.CreateUser(&entities.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hash,
		RoleID:   role.ID,
	}, &repositories.CreateUserOption{
		Tx: tx,
	})
	if err != nil {
		rollback()
		return nil, err
	}

	if err := s.authReposiotry.AddRoleForUser(&repositories.AddRoleForUserRequest{
		UserName: req.Name,
		RoleName: role.Name,
	}, &repositories.AddRoleForUserOption{
		Tx: tx,
	}); err != nil {
		rollback()
		return nil, err
	}
	commit()

	user.Role = role

	return user, nil

}

func (s *AuthService) CreateRole(req *CreateRoleRequest) error {
	shouldCommit := false
	tx := req.Tx
	if tx == nil {
		shouldCommit = true
		var err error
		tx, err = s.dbTransactioner.Begin()
		if err != nil {
			return err
		}
	}
	if _, err := s.roleRepository.AddRole(&repositories.AddRoleRequest{
		RoleName: req.RoleName,
	}, &repositories.AddRoleOption{
		Tx: tx,
	}); err != nil {
		if !errors.ErrAlreadyExist.Is(err) {
			return err
		}
	}

	if len(req.Policies) > 0 {
		if err := s.authReposiotry.AddPoliciesForRole(&repositories.AddPoliciesForRoleRequest{
			RoleName: req.RoleName,
			Policies: req.Policies,
		}, nil); err != nil {
			return err
		}
	}

	if shouldCommit {
		if err := s.dbTransactioner.Commit(tx); err != nil {
			return err
		}
	}
	return nil
}

func (s *AuthService) GetUserByID(tx interface{}, id uint64) (*entities.User, error) {
	user, err := s.userRepository.GetUser(&repositories.GetUserRequest{
		UserID:   &id,
		JoinRole: pointer.Bool(true),
	}, &repositories.GetUserOption{
		Tx: tx,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

type AuthService struct {
	dbTransactioner repositories.DBTransactioner
	userRepository  repositories.UserRepository
	roleRepository  repositories.RoleRepository
	authReposiotry  repositories.AuthRepository
	jWTManager      *jwt.JWTManager
}

func NewAuthService(
	dbTransactioner repositories.DBTransactioner,
	userRepository repositories.UserRepository,
	roleRepository repositories.RoleRepository,
	authReposiotry repositories.AuthRepository,
	jWTManager *jwt.JWTManager,
) *AuthService {
	return &AuthService{
		dbTransactioner: dbTransactioner,
		userRepository:  userRepository,
		roleRepository:  roleRepository,
		authReposiotry:  authReposiotry,
		jWTManager:      jWTManager,
	}
}

var _ validators.AuthValidator = (*AuthService)(nil)
