module server_module

go 1.15

require (
	auth_module/codes v0.0.0-00010101000000-000000000000 // indirect
	cloud.google.com/go/firestore v1.3.0
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.6.3
	github.com/joho/godotenv v1.3.0
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
	google.golang.org/api v0.30.0
	server_module/auth v0.0.0-00010101000000-000000000000
	server_module/query v0.0.0-00010101000000-000000000000
)

replace (
	auth_module/codes => ./auth/codes
	server_module/auth => ./auth
	server_module/query => ./query
)
