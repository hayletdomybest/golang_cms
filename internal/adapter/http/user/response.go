package user

import "weex_admin/internal/domain/entities"

type GetUserItem struct {
	ID       uint64 `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	RoleName string `json:"role_name"`
}

func getUserItemFromEntity(entity *entities.User) *GetUserItem {
	return getUserItemFromEntities(entity)[0]
}

func getUserItemFromEntities(entities ...*entities.User) []*GetUserItem {
	var res []*GetUserItem
	for _, entity := range entities {
		res = append(res, &GetUserItem{
			ID:       entity.ID,
			UserName: entity.Name,
			Email:    entity.Email,
			RoleName: entity.Role.Name,
		})
	}
	return res
}
