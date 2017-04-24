#!/usr/bin/env sh

# Ensure this script fails if anything errors
set -e

# Copy the private key for access to deploy
mkdir -p /root/.ssh/ && chmod 0600 /root/.ssh
echo "${DEPLOYMENT_PRIVATE_KEY}" > /root/.ssh/id_rsa && chmod 0600 /root/.ssh/id_rsa

# Untar the automation release
echo "Using the following automation artifact version: `cat automation-source-code/version`"
tar -xf automation-source-code/source.tar.gz

# Change to plays directory
cd hellofresh-janus-automation-*/plays

# Deploy to staging using ansible
export ANSIBLE_FORCE_COLOR=true
ansible-playbook \
    -i ../${DEPLOYMENT_ENVIRONMENT}.ini \
    -u policy \
    -vvvv \
    --skip-tags=provision \
    -e deployment_force=true \
    ${GROUP_NAME}.yml
