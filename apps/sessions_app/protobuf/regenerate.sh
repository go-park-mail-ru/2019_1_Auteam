#!/bin/bash
protoc -I ./ sessions_route.proto --go_out=plugins=grpc:./
