package database

import (
	"weex_admin/internal/domain/entities"
	"weex_admin/internal/domain/repositories"
	"weex_admin/internal/shared/constants"
	"weex_admin/internal/shared/errors"
)

type ResourceRepository struct {
	resources []*entities.Resource

	byPath  map[string]*entities.Resource
	byGroup map[string][]*entities.Resource
}

func (r *ResourceRepository) init() {
	r.resources = []*entities.Resource{
		{
			Group: constants.ROLE_GROUP,
			Path:  constants.ROLE_CREATE_PATH,
		},
		{
			Group: constants.ROLE_GROUP,
			Path:  constants.ROLE_GET_PATH,
		},
	}

	for _, resource := range r.resources {
		if _, ok := r.byPath[resource.Path]; ok {
			panic("resource repository encounter error: duplicated resource path")
		}
		r.byPath[resource.Path] = resource
		r.byGroup[resource.Group] = append(r.byGroup[resource.Group], resource)
	}
}

// FindByGroup finds resources by group.
func (r *ResourceRepository) FindByGroup(group string) ([]*entities.Resource, error) {
	resources, ok := r.byGroup[group]
	if !ok {
		return nil, errors.ErrNotFound.WithMessagef("can not find group %s", group)
	}

	var result []*entities.Resource
	for _, res := range resources {
		var clone entities.Resource = *res
		result = append(result, &clone)
	}
	return result, nil
}

// FindByPath finds a resource by its path.
func (r *ResourceRepository) FindByPath(path string) (*entities.Resource, error) {
	resource, ok := r.byPath[path]
	if !ok {
		return nil, errors.ErrNotFound.WithMessagef("can not find path %s", path)
	}

	var res = *resource
	return &res, nil
}

// List lists all resources.
func (r *ResourceRepository) List() ([]*entities.Resource, error) {
	var result []*entities.Resource
	for _, res := range r.resources {
		var clone entities.Resource = *res
		result = append(result, &clone)
	}
	return result, nil
}

var _ repositories.ResourceRepository = (*ResourceRepository)(nil)

var resourceRepository *ResourceRepository = NewResourceRepository()

func NewResourceRepository() *ResourceRepository {
	repo := &ResourceRepository{
		byPath:  make(map[string]*entities.Resource),
		byGroup: make(map[string][]*entities.Resource),
	}
	repo.init()

	return repo
}
