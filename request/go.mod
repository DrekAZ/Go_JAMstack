module request_module

go 1.15

require (
	cloud.google.com/go/firestore v1.3.0
	github.com/gin-gonic/gin v1.6.3
	server_module/convert v0.0.0-00010101000000-000000000000
	server_module/query v0.0.0-00010101000000-000000000000
)

replace (
	server_module/convert => ../convert
	server_module/query => ../query
)
