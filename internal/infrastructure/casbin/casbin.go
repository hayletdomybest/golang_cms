package casbin

import (
	"weex_admin/internal/shared/errors"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
)

var (
	RBAC_MODEL_CONF = "rbac_model.conf"
)

func NewCasbin(db *gorm.DB) (*casbin.Enforcer, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, errors.ErrInternal.
			Wrap(err).
			WithMessage("failed to create gorm adapter for casbin")
	}

	enforcer, err := casbin.NewEnforcer(RBAC_MODEL_CONF, adapter)
	if err != nil {
		return nil, errors.ErrInternal.
			Wrap(err).
			WithMessage("failed to NewEnforcer for casbin")
	}

	if err := enforcer.LoadPolicy(); err != nil {
		return nil, errors.ErrInternal.
			Wrap(err).
			WithMessage("failed to LoadPolicy for casbin")
	}

	return enforcer, nil
}
