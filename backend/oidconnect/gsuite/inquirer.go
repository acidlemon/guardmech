package gsuite

import (
	"context"
	"log"

	"github.com/acidlemon/guardmech/backend/oidconnect"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

type inquirer struct {
	svc *admin.Service
}

func NewGroupInquirer(ctx context.Context) (oidconnect.GroupInquirer, error) {
	svc, err := adminService(ctx)
	if err != nil {
		return nil, err
	}

	gi := &inquirer{
		svc: svc,
	}
	return gi, nil
}

func adminService(ctx context.Context) (*admin.Service, error) {
	client, err := google.DefaultClient(ctx, admin.AdminDirectoryGroupReadonlyScope, admin.AdminDirectoryGroupMemberReadonlyScope)
	if err != nil {
		log.Println("Failed to prepare Client", err)
		return nil, err
	}

	svc, err := admin.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Println("Failed to retrieve directory Client", err)
		return nil, err
	}

	return svc, nil
}
func (gi *inquirer) IsMember(ctx context.Context, email, group string) (bool, error) {
	result, err := gi.svc.Members.HasMember(group, email).Do()
	if err != nil {
		log.Println("admin.Members.HasMember returned error:", err)
		return false, err
	}

	return result.IsMember, nil
}
