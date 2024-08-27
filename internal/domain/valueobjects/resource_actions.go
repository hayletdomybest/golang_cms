package valueobjects

import "weex_admin/internal/shared/errors"

type ResourceAction string

const (
	Read      ResourceAction = "read"
	Write     ResourceAction = "write"
	ReadWrite ResourceAction = "read_write"
)

func (a ResourceAction) GetActions() []ResourceAction {
	if a == ReadWrite {
		return []ResourceAction{Read, Write}
	}
	return []ResourceAction{a}
}

func Wrap(s string) ([]ResourceAction, error) {
	switch s {
	case "read":
		return []ResourceAction{Read}, nil
	case "write":
		return []ResourceAction{Write}, nil
	case "read_write":
		return []ResourceAction{Read, Write}, nil
	default:
		return nil, errors.ErrPanic.WithMessagef("invalid resource action: %s", s)
	}
}
