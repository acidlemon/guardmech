module github.com/acidlemon/guardmech/backend

go 1.16

require (
	github.com/acidlemon/seacle v0.3.0
	github.com/coreos/go-oidc/v3 v3.1.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/k0kubun/pp/v3 v3.1.0
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.13.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/serenize/snaker v0.0.0-20201027110005-a7ad2135616e // indirect
	golang.org/x/crypto v0.0.0-20220314234724-5d542ad81a58
	golang.org/x/oauth2 v0.0.0-20220309155454-6242fa91716a
	golang.org/x/tools v0.1.9 // indirect
	google.golang.org/api v0.72.0
	google.golang.org/genproto v0.0.0-20220314164441-57ef72a4c106 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)

//replace github.com/acidlemon/seacle => ../../seacle
