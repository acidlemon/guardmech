package main

import (
	"fmt"
	"log"
	"reflect"

	db "github.com/acidlemon/guardmech/backend/persistence/db"
	"github.com/acidlemon/seacle"
)

//go:generate go run github.com/acidlemon/guardmech/backend/gen

type genset struct {
	tp      reflect.Type
	pkg     string
	table   string
	outfile string
}

func prepareGenset() []genset {
	var set = []struct {
		tp    reflect.Type
		table string
	}{
		{reflect.TypeOf(db.PrincipalRow{}), "principal"},
		{reflect.TypeOf(db.AuthOIDCRow{}), "auth_oidc"},
		{reflect.TypeOf(db.AuthAPIKeyRow{}), "auth_apikey"},
		{reflect.TypeOf(db.GroupRow{}), "group_info"},
		{reflect.TypeOf(db.PermissionRow{}), "permission"},
		{reflect.TypeOf(db.RoleRow{}), "role_info"},
		{reflect.TypeOf(db.MappingRuleRow{}), "mapping_rule"},

		{reflect.TypeOf(db.PrincipalGroupMapRow{}), "principal_group_map"},
		{reflect.TypeOf(db.PrincipalRoleMapRow{}), "principal_role_map"},
		{reflect.TypeOf(db.GroupRoleMapRow{}), "group_role_map"},
		{reflect.TypeOf(db.RolePermissionMapRow{}), "role_permission_map"},
	}

	result := []genset{}
	for _, v := range set {
		result = append(result, genset{
			tp:      v.tp,
			pkg:     "db",
			table:   v.table,
			outfile: fmt.Sprintf("../persistence/db/%s_gen.go", v.table),
		})
	}

	return result
}

func main() {

	log.Println("go!! generate!!")

	gen := seacle.Generator{
		Tag: "db",
	}

	sets := prepareGenset()
	for _, s := range sets {
		err := gen.Generate(s.tp, s.pkg, s.table, s.outfile)
		if err != nil {
			log.Printf("failed to generate file: table=%s, outfile=%s, err=%s", s.table, s.outfile, err)
		}
	}
}
