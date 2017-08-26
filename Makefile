all:LghSampleMain
	
LghSampleMain:LghSampleMain.go common.go LghSqlExecutor.go LghSqlHelper.go struct.go LghReadConfJson.go func_LUser.go
	go build LghSampleMain.go common.go LghSqlExecutor.go LghSqlHelper.go struct.go LghReadConfJson.go func_LUser.go

install:all
	cp LghSampleMain ./bin

clean:
	rm -f LghSampleMain
