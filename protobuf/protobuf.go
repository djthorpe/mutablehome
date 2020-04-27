package protobuf

//go:generate protoc castchannel/castchannel.proto --go_out=plugins=grpc:. --go_opt=paths=source_relative
//go:generate protoc mutablehome/mutablehome.proto --go_out=plugins=grpc:. --go_opt=paths=source_relative

/*
	This folder contains all the protocol buffer definitions including
	the RPC Service definitions. You generate golang code by running:

	go generate -x github.com/djthorpe/mutablehome/protobuf

	where you have installed the protoc compiler and the GRPC plugin for
	golang. In order to do that on a Mac:

	mac# brew install protobuf
	mac# go get -u github.com/golang/protobuf/protoc-gen-go

	On Debian Linux use the following commands instead:

	rpi# sudo apt install protobuf-compiler
	rpi# sudo apt install libprotobuf-dev
	rpi# go get -u github.com/golang/protobuf/protoc-gen-go

*/
