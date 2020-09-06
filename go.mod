module server_module

go 1.15

require (
	cloud.google.com/go/firestore v1.3.0
	github.com/gin-gonic/gin v1.6.3
	github.com/stretchr/testify v1.6.1 // indirect
	google.golang.org/api v0.30.0
	server_module/request v0.0.0-00010101000000-000000000000
	server_module/setting v0.0.0-00010101000000-000000000000
)

replace (
	server_module/convert => ./convert
	server_module/query => ./query
	server_module/request => ./request
	server_module/setting => ./setting
)
