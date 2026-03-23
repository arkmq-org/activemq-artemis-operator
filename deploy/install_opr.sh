#!/bin/bash

echo "Deploying operator to watch single namespace"

DEPLOY_PATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

if oc version; then
    KUBE_CLI=oc
else
    KUBE_CLI=kubectl
fi

$KUBE_CLI create -f $DEPLOY_PATH/crds
$KUBE_CLI create -f $DEPLOY_PATH/service_account.yaml
$KUBE_CLI create -f $DEPLOY_PATH/role.yaml
$KUBE_CLI create -f $DEPLOY_PATH/role_binding.yaml
$KUBE_CLI create -f $DEPLOY_PATH/election_role.yaml
$KUBE_CLI create -f $DEPLOY_PATH/election_role_binding.yaml
$KUBE_CLI create -f $DEPLOY_PATH/network_policy.yaml
$KUBE_CLI create -f $DEPLOY_PATH/operator.yaml
