package db

type PrincipalRoleMapRow struct {
	PrincipalSeqID int64 `db:"principal_seq_id,primary"`
	RoleSeqID      int64 `db:"role_seq_id,primary"`
}

type PrincipalGroupMapRow struct {
	PrincipalSeqID int64 `db:"principal_seq_id,primary"`
	GroupSeqID     int64 `db:"group_seq_id,primary"`
}

type GroupRoleMapRow struct {
	GroupSeqID int64 `db:"group_seq_id,primary"`
	RoleSeqID  int64 `db:"role_seq_id,primary"`
}

type RolePermissionMapRow struct {
	RoleSeqID       int64 `db:"role_seq_id,primary"`
	PermissionSeqID int64 `db:"permission_seq_id,primary"`
}
