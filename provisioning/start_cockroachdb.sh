#!/bin/bash
region=$(
    curl -s http://169.254.169.254/latest/dynamic/instance-identity/document \
        | grep region \
        | awk -F\" '{print $4}'
)
asginstanceips=()

for asginstanceid in $(
    aws --region "$region" autoscaling describe-auto-scaling-groups \
        --auto-scaling-group-name rteregulov-ec2-asg \
        | grep -i instanceid  \
        | awk '{ print $2}' \
        | cut -d',' -f1 \
        | sed -e 's/"//g'
)
do
    asginstanceips+=("$(
        aws --region "$region" ec2 describe-instances --instance-ids "$asginstanceid" \
            | grep -i PrivateIpAddress \
            | awk '{ print $2 }' \
            | head -1 \
            | cut -d"," -f1 \
            | sed -e 's/"//g'
    )")
done

asginstances=$( IFS=$','; echo "${asginstanceips[*]}" )
instanceip=$(ec2metadata --local-ipv4)

cockroach start --insecure --background --advertise-addr="$instanceip" --join="$asginstances"
