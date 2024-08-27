package gorm

import (
	"weex_admin/internal/domain/entities"
	"weex_admin/internal/domain/repositories"
	"weex_admin/internal/shared/errors"
	sharedgorm "weex_admin/internal/shared/gorm"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func (u *UserRepository) CreateUser(user *entities.User, option *repositories.CreateUserOption) (*entities.User, error) {
	if users, err := u.CreateUsers([]*entities.User{user}, option); err != nil {
		return nil, err
	} else {
		return users[0], nil
	}
}

func (u *UserRepository) CreateUsers(users []*entities.User, option *repositories.CreateUserOption) ([]*entities.User, error) {
	tx, err := u.getTx(option)
	if err != nil {
		return nil, err
	}

	if err := tx.CreateInBatches(users, len(users)).Error; err != nil {
		return nil, sharedgorm.TryWrapErr(err, "failed to create users")
	}
	return users, nil
}

func (u *UserRepository) GetUser(req *repositories.GetUserRequest, option *repositories.GetUserOption) (*entities.User, error) {
	if reps, err := u.GetUsers([]*repositories.GetUserRequest{req}, option); err != nil {
		return nil, err
	} else {
		if len(reps) > 0 {
			return reps[0], nil
		}
		return nil, nil
	}
}

func (u *UserRepository) GetUsers(req []*repositories.GetUserRequest, option *repositories.GetUserOption) ([]*entities.User, error) {
	query, err := u.getTx(option)
	if err != nil {
		return nil, err
	}

	var userIDs []uint64
	var userNames []string
	joinRole := false

	for _, r := range req {
		if r.JoinRole != nil && *r.JoinRole {
			joinRole = true
		}
		if r.UserID != nil {
			userIDs = append(userIDs, *r.UserID)
		} else if r.UserName != nil {
			userNames = append(userNames, *r.UserName)
		} else {
			return nil, errors.ErrInvalidParams.WithMessage("should pass either user id or user name")
		}
	}

	var users []*entities.User
	if joinRole {
		query = query.Preload("Role")
	}

	if len(userIDs) > 0 {
		query = query.Where("id IN (?)", userIDs)
	}
	if len(userNames) > 0 {
		query = query.Or("name IN (?)", userNames)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

var _ repositories.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) getTx(option *repositories.TransactionOption) (*gorm.DB, error) {
	db := u.db
	if option != nil && option.Tx != nil {
		if tdDB, ok := option.Tx.(*gorm.DB); !ok {
			return nil, errors.ErrPanic.WithMessagef("should pass gorm.DB as transaction")
		} else {
			db = tdDB
		}
	}
	return db.Model(&entities.User{}), nil
}
