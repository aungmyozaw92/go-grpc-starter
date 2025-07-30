#!/bin/bash

echo "Generating proto files..."

# Remove old proto files
rm -rf proto/userpb/*.pb.go

# Generate new proto files
protoc --go_out=proto/userpb \
       --go_opt=module=github.com/aungmyozaw92/go-grpc-starter \
       --go-grpc_out=proto/userpb \
       --go-grpc_opt=module=github.com/aungmyozaw92/go-grpc-starter \
       proto/user.proto

# Move files to correct location if nested structure was created
if [ -d "proto/userpb/proto/userpb" ]; then
    mv proto/userpb/proto/userpb/* proto/userpb/
    rm -rf proto/userpb/proto
fi

echo "✅ Proto files generated successfully in proto/userpb/"
ls -la proto/userpb/*.pb.go 