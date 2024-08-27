package repositories

import "weex_admin/internal/domain/entities"

type ResourceRepository interface {
	FindByGroup(group string) ([]*entities.Resource, error)
	FindByPath(path string) (*entities.Resource, error)
	List() ([]*entities.Resource, error)
}
