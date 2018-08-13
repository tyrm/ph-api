go build -o puphaus-api.exe .

SET CGO_ENABLED=0
SET GOOS=linux
go build -a -installsuffix cgo -o puphaus-api .