export GO111MODULE=on
#export GONOPROXY=
export GOPROXY=https://goproxy.cn
go mod tidy
go build -gcflags "-N -l" -o "escript" main.go