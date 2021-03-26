#!/bin/bash

# There are 4 arguments passed in.
# $1 is the name of the micro service
# $2 is the image id
# $3 is the GITHUB OAUTH Key
# This script will first get the latest raw of the docker-compose file for impulse and copy it into a local file
# It will then replace the entire name of the image for the passed in ($1) micro service with the image name that is built
# Finally it will commit the new docker-compose file to impulse with a message

set -x
# echo running script
#
# echo getting file from GITHUB
curl -H 'Authorization: token '"$3"'' -H 'Accept: application/vnd.github.v3.raw' -O -L https://api.github.com/repos/motiv-labs/impulse/contents/docker-tools/compose-template/config/astro.conf

# echo search and replace
sed -i "s+motivlabs/$1:.*+$2+g" astro.conf
# echo pushing file to GITHUB

curl -X PUT -H 'Authorization: token '"$3"'' -d '{"message":"updated janus image","content":"'"$(base64 -w0 astro.conf)"'","sha":'"$(curl -s -X GET -H 'Authorization: token '"$3"'' https://api.github.com/repos/motiv-labs/impulse/contents/docker-tools/compose-template/config/astro.conf | jq .sha)"'}' -L https://api.github.com/repos/motiv-labs/impulse/contents/docker-tools/compose-template/config/astro.conf

# update googlecloud image
curl -H 'Authorization: token '"$3"'' -H 'Accept: application/vnd.github.v3.raw' -O -L https://api.github.com/repos/motiv-labs/impulse-googlecloud/contents/impulse/janus-deployment.yaml
sed -i "s+motivlabs/$1:.*+$2+g" janus-deployment.yaml
curl -X PUT -H 'Authorization: token '"$3"'' -d '{"message":"updated janus image","content":"'"$(base64 -w0 janus-deployment.yaml)"'","sha":'"$(curl -s -X GET -H 'Authorization: token '"$3"'' https://api.github.com/repos/motiv-labs/impulse-googlecloud/contents/impulse/janus-deployment.yaml | jq .sha)"'}' -L https://api.github.com/repos/motiv-labs/impulse-googlecloud/contents/impulse/janus-deployment.yaml

# echo finished script
