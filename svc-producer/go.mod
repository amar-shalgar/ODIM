module github.com/ODIM-Project/ODIM/svc-producer

go 1.13

require (
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20201201072448-9772421f1b55
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/kataras/iris/v12 v12.1.9-0.20200616210209-a85c83b70ad0
	github.com/sirupsen/logrus v1.4.2
)

replace (
	github.com/ODIM-Project/ODIM/lib-dmtf => ../lib-dmtf
	github.com/ODIM-Project/ODIM/lib-messagebus => ../lib-messagebus
	github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager
	github.com/ODIM-Project/ODIM/lib-rest-client => ../lib-rest-client
	github.com/ODIM-Project/ODIM/lib-utilities => ../lib-utilities
)
