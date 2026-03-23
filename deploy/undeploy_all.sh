#!/bin/bash

echo "Undeploy everything..."

DEPLOY_PATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

if oc version; then
    KUBE_CLI=oc
else
    KUBE_CLI=kubectl
fi

$KUBE_CLI delete -f $DEPLOY_PATH/crds
$KUBE_CLI delete -f $DEPLOY_PATH/service_account.yaml
$KUBE_CLI delete -f $DEPLOY_PATH/role.yaml
$KUBE_CLI delete -f $DEPLOY_PATH/role_binding.yaml
$KUBE_CLI delete -f $DEPLOY_PATH/cluster_role.yaml
$KUBE_CLI delete -f $DEPLOY_PATH/cluster_role_binding.yaml
$KUBE_CLI delete -f $DEPLOY_PATH/election_role.yaml
$KUBE_CLI delete -f $DEPLOY_PATH/election_role_binding.yaml
$KUBE_CLI delete -f $DEPLOY_PATH/network_policy.yaml
$KUBE_CLI delete -f $DEPLOY_PATH/operator.yaml
