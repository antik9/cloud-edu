#!/bin/bash
instanceid="$(ec2metadata --instance-id)"
region=$(
    curl -s http://169.254.169.254/latest/dynamic/instance-identity/document \
        | grep region \
        | awk -F\" '{print $4}'
)
username="$(
    aws --region "$region" ec2 describe-tags \
        --filters "Name=resource-id,Values=${instanceid}" "Name=key,Values=student" \
        --output text \
    | cut -f5
)"

cd ..
sed -i "s/student/$username/g" hello.go
nohup go run hello.go &
