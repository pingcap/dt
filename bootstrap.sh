#!/bin/bash
echo "install go dependencies"

go get -u "github.com/ngaut/log"
go get -u "github.com/juju/errors"
go get -u "gopkg.in/check.v1"
go get -u "github.com/gorilla/mux"
go get -u "github.com/BurntSushi/toml"

