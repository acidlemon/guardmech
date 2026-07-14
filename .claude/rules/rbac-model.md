---
paths:
  - "backend/app/logic/membership/*.go"
---

# RBAC モデルの型からは見えない規約

## 背景 / 罠の説明

`backend/app/logic/membership` パッケージの RBAC には、struct 定義だけでは分からない
以下の規約がある。

- **有効ロールの算出**: `Principal.Roles()` は「直接付与ロール（`AttachedRoles()`）」と
  「所属グループ経由のロール（各 `Group.Roles()`）」の和集合を、`RoleID` で重複排除して返す。
- **有効パーミッションの算出**: `Principal.HavingPermissions()` は `Roles()` が返す
  全ロールの `Permissions()` を集約し、`PermissionID` で重複排除して返す。
  `Group.HavingPermissions()` も同様にグループ配下のロールから集約する。
- **well-known エンティティの固定 UUID**: 名前で判定して固定 UUID を割り当てるエンティティがある。
  - Role: `RoleOwnerName = "_GuardmechOwnerRole"` に対して `RoleOwnerID`（`role.go` の `newRole`）。
  - Permission: `PermissionOwnerName`/`PermissionWriteName`/`PermissionReadOnlyName`
    （`"_GUARDMECH_OWNER"`/`"_GUARDMECH_WRITE"`/`"_GUARDMECH_READONLY"`）に対して、
    それぞれ `PermissionOwnerID`/`PermissionWriteID`/`PermissionReadOnlyID`（`permission.go` の `newPermission`）。
  これらの固定 UUID は frontend の `AuthorityStatus`（`OWNER`/`WRITABLE`/`READONLY`）と対応する
  （→ `.claude/rules/frontend-authority.md`）。
- **MappingRule（ログイン時のロール/グループ自動付与）**: `mapping_rule.go` の `MappingType` は
  `1=MappingSpecificDomain`（`"@"+Detail` サフィックス一致）/
  `2=MappingWholeDomain`（`Detail` サフィックス一致）/
  `3=MappingGroupMember`（GSuite の `GroupInquirer` に問い合わせ。`inquirer` が `nil` なら
  そのルールは**黙ってスキップ**される）/
  `4=MappingSpecificAddress`（完全一致）。
  `MappingRuleManager`（`mapping_rule_manager.go`）は生成時に `Priority` **昇順**の
  `sort.Stable` で評価順を確定し、`FindMatchedRules()` はマッチした**全ルール**を返す。
  マッチが 0 件の場合はエラーを返す（認証エラーとして扱われる）。
- **Attach/Detach 系は idempotent**: `Principal.AttachRole/AttachGroup`、
  `Role.AttachPermission`、`Group.AttachRole` は既に付与済みなら何もせず `nil` を返す。
  `DetachRole/DetachGroup/DetachPermission` は対象が存在しなくても `nil` を返す
  （エラーにしない設計）。
