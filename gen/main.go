package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/acidlemon/guardmech/infra"
	"github.com/acidlemon/seacle"
)

//go:generate go run github.com/acidlemon/guardmech/gen

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
		{reflect.TypeOf(infra.Principal{}), "principal"},
		{reflect.TypeOf(infra.Auth{}), "auth"},
		{reflect.TypeOf(infra.APIKey{}), "api_key"},
		{reflect.TypeOf(infra.Group{}), "group_info"},
		{reflect.TypeOf(infra.Permission{}), "permission"},
		{reflect.TypeOf(infra.Role{}), "role_info"},
	}

	result := []genset{}
	for _, v := range set {
		result = append(result, genset{
			tp:      v.tp,
			pkg:     "infra",
			table:   v.table,
			outfile: fmt.Sprintf("../infra/%s_gen.go", v.table),
		})
	}

	return result
}

func main() {

	log.Println("go!! generate!!")

	gen := seacle.Generator{
		Tag: "json",
	}

	// files := []genset{
	// 	{reflect.TypeOf(infra.Principal{}), "infra", "principal", "../infra/principal_gen.go"},
	// 	{reflect.TypeOf(infra.Auth{}), "infra", "auth", "../infra/auth_gen.go"},
	// 	{reflect.TypeOf(infra.APIKey{}), "infra", "api_key", "../infra/api_key_gen.go"},
	// 	{reflect.TypeOf(infra.Group{}), "infra", "group", "../infra/group_gen.go"},
	// 	{reflect.TypeOf(infra.Permission{}), "infra", "permission", "../infra/permission_gen.go"},
	// }

	sets := prepareGenset()

	for _, s := range sets {
		err := gen.Generate(s.tp, s.pkg, s.table, s.outfile)
		if err != nil {
			log.Printf("failed to generate file: table=%s, outfile=%s, err=%s", s.table, s.outfile, err)
		}
	}
}
