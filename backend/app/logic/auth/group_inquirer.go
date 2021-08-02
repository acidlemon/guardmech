package auth

import (
	"context"

	"github.com/acidlemon/guardmech/backend/oidconnect"
)

type GroupInquirer struct {
	inquirer oidconnect.GroupInquirer
}

func NewGroupInquirer(inquirer oidconnect.GroupInquirer) *GroupInquirer {
	return &GroupInquirer{
		inquirer: inquirer,
	}
}

func (g *GroupInquirer) IsMember(ctx context.Context, email, group string) (bool, error) {
	result, err := g.inquirer.IsMember(ctx, email, group)
	if err != nil {
		return false, err
	}

	return result, nil
}
