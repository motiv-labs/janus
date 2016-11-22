#!/usr/bin/env sh

# Ensure this script fails if anything errors
set -e

# Copy the private key for access to deploy
mkdir -p /root/.ssh/ && chmod 0600 /root/.ssh
echo "${DEPLOYMENT_PRIVATE_KEY}" > /root/.ssh/id_rsa && chmod 0600 /root/.ssh/id_rsa

# Untar the automation release
zcat automation-source-code/automation-artifact.tar.gz | tar -xf -

# Change to plays directory
cd automation-artifact/plays

# Temporary fix for the socket
sed -i 's/%%h-%%p-%%r/%%h-%%r/g' ansible.cfg

# Deploy to staging using ansible
export ANSIBLE_FORCE_COLOR=true
ansible-playbook \
    -i ../${DEPLOYMENT_ENVIRONMENT}.ini \
    -u policy \
    -vvvv \
    -t deployment \
    -e deployment_force=true \
    ${GROUP_NAME}.yml
