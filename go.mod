module server_module

go 1.15

require (
	cloud.google.com/go/firestore v1.3.0
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.6.3
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
	google.golang.org/api v0.30.0
	server_module/auth v0.0.0-00010101000000-000000000000
	server_module/convert v0.0.0-00010101000000-000000000000 // indirect
	server_module/query v0.0.0-00010101000000-000000000000
	server_module/setting v0.0.0-00010101000000-000000000000
	server_module/status_code v0.0.0-00010101000000-000000000000 // indirect
)

replace (
	server_module/auth => ./auth
	server_module/convert => ./convert
	server_module/query => ./query
	server_module/setting => ./setting
	server_module/status_code => ./status_code
)
