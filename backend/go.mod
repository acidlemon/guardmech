module github.com/acidlemon/guardmech/backend

go 1.16

require (
	github.com/acidlemon/seacle v0.3.0
	github.com/coreos/go-oidc/v3 v3.0.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.7.3
	github.com/k0kubun/pp/v3 v3.0.7
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.13.0 // indirect
	github.com/pkg/errors v0.8.1
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e
	golang.org/x/oauth2 v0.0.0-20210628180205-a41e5a781914
	google.golang.org/api v0.52.0
)

//replace github.com/acidlemon/seacle => ../../seacle
