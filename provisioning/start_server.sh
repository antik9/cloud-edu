#!/bin/bash
username="$(
    aws --region "$region" ec2 describe-tags \
        --filters "Name=resource-id,Values=${instanceid}" "Name=key,Values=student" \
        --output text \
    | cut -f5
)"

sed -i "s/student/$username/g" hello.go
go run hello.go
