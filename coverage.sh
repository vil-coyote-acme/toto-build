#!/usr/bin/env bash


usage ()
{
  echo 'Usage : Script <package>'
  exit
}

if [ $# -eq 0 ]
  then
    usage
fi
# todo permit to use atomic covermode for parallels treatments testing
go test -covermode=count -coverprofile=coverage.out toto-build-agent/"$1"
go tool cover -html=coverage.out