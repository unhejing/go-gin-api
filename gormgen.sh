#!/bin/bash

# ./gormgen.sh 127.0.0.1:3306 root 123456 go_gin_api config
shellExit()
{
if [ $1 -eq 1 ]; then
    printf "\nfailed!!!\n\n"
    exit 1
fi
}

printf "\nRegenerating file\n\n"
time go run -v ./generate/mysqlmd/main.go  -addr $1 -user $2 -pass $3 -name $4 -table $5
shellExit $?

printf "\ncreate curd code : \n"
time go build -o gormgen ./generate/gormgen/main.go
shellExit $?

if [ ! -d $GOPATH/bin ];then
   mkdir -p $GOPATH/bin
fi

#mv gormgen $GOPATH/bin
mv gormgen $GOPATH/bin/gormgen.exe
shellExit $?

#go generate ./...
printf "\ngenerate repostory...\n"
printf ./model/$5_model/
printf "\n"
go generate ./model/$5_model/
shellExit $?

# generate serive handler dto
printf "\ngenerate service,handler,dto\n"
go run -v ./generate/servicegen/main.go -service $5

printf "\nFormatting code\n\n"
time go run -v ./generate/mfmt/main.go
shellExit $?

printf "\nDone.\n\n"
