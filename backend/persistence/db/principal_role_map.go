package db

type PrincipalRoleMapRow struct {
	PrincipalSeqID int64 `db:"principal_seq_id,primary"`
	RoleSeqID      int64 `db:"role_seq_id,primary"`
}

func (row *PrincipalRoleMapRow) TargetSeqID() int64 {
	return row.RoleSeqID
}

type PrincipalGroupMapRow struct {
	PrincipalSeqID int64 `db:"principal_seq_id,primary"`
	GroupSeqID     int64 `db:"group_seq_id,primary"`
}

func (row *PrincipalGroupMapRow) TargetSeqID() int64 {
	return row.GroupSeqID
}

type GroupRoleMapRow struct {
	GroupSeqID int64 `db:"group_seq_id,primary"`
	RoleSeqID  int64 `db:"role_seq_id,primary"`
}

func (row *GroupRoleMapRow) TargetSeqID() int64 {
	return row.RoleSeqID
}

type RolePermissionMapRow struct {
	RoleSeqID       int64 `db:"role_seq_id,primary"`
	PermissionSeqID int64 `db:"permission_seq_id,primary"`
}

func (row *RolePermissionMapRow) TargetSeqID() int64 {
	return row.PermissionSeqID
}
