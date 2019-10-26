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

echo "CREATE USER IF NOT EXISTS $username;
CREATE DATABASE IF NOT EXISTS ${username}db;
GRANT ALL ON DATABASE ${username}db TO $username;
USE ${username}db;
CREATE TABLE IF NOT EXISTS views (
  view_id SERIAL PRIMARY KEY,
  client_ip INET,
  view_date TIMESTAMP
);" |  cockroach sql --insecure

