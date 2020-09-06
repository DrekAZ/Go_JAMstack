module convert_module

go 1.15

replace server_module/setting => ../setting

require (
	github.com/gin-gonic/gin v1.6.3
	server_module/setting v0.0.0-00010101000000-000000000000
)
