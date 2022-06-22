#!/usr/bin/env bash

set -eo pipefail

go get -u github.com/bufbuild/buf/cmd/buf@v1.0.0-rc10
go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u google.golang.org/protobuf/cmd/protoc-gen-go
go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


protoc_gen_gocosmos() {
  if ! grep "github.com/gogo/protobuf => github.com/regen-network/protobuf" go.mod &>/dev/null ; then
    echo -e "\tPlease run this command from somewhere inside the ibc-go folder."
    return 1
  fi

  go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@latest 2>/dev/null
}

protoc_gen_gocosmos

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  buf protoc \
  -I "proto" \
  -I "third_party/proto" \
  --gocosmos_out=plugins=interfacetype+grpc,\
Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
  --grpc-gateway_out=logtostderr=true:. \
  $(find "${dir}" -maxdepth 1 -name '*.proto')

done

# command to generate docs using protoc-gen-doc
# buf protoc \
#     -I "proto" \
#     -I "third_party/proto" \
#     --doc_out=./docs/ibc \
#     --doc_opt=./docs/protodoc-markdown.tmpl,proto-docs.md \
#     $(find "$(pwd)/proto" -maxdepth 7 -name '*.proto')
go mod tidy -compat=1.17

# move proto files to the right places
cp -r github.com/cosmos/ibc-go/v*/modules/* modules/
rm -rf github.com
