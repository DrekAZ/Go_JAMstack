module auth_module

go 1.15

require (
	cloud.google.com/go/firestore v1.3.0
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.6.3
	github.com/pquerna/cachecontrol v0.0.0-20200819021114-67c6ae64274f // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	server_module/auth_module/codes v0.0.0-00010101000000-000000000000
)

replace server_module/auth_module/codes => ./codes