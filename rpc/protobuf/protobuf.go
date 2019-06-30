/*
	Mutablehome: Home Automation in Golang
	(c) Copyright David Thorpe 2019
	All Rights Reserved

    https://github.com/djthorpe/mutablehome/
	For Licensing and Usage information, please see LICENSE
*/

package protobuf

//go:generate protoc mutablehome/mutablehome.proto --go_out=plugins=grpc:.

/*
	This folder contains all the protocol buffer definitions including
	the RPC Service definitions. You generate golang code by running:

	go generate -x github.com/djthorpe/mutablehome/rpc/protobuf

	where you have installed the protoc compiler and the GRPC plugin for
	golang. In order to do that on a Mac:

	mac# brew install protobuf
	mac# go get -u github.com/golang/protobuf/protoc-gen-go

	On Debian Linux (including Raspian Linux) use the following commands
	instead:

	rpi# sudo apt install protobuf-compiler
	rpi# go get -u github.com/golang/protobuf/protoc-gen-go
*/
