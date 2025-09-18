#!/bin/bash

version=$1
dir=$2
YQ=${3:-yq}

$YQ -i ".version=\"${version}\"" ${dir}/Chart.yaml
$YQ -i ".appVersion=\"${version}\"" ${dir}/Chart.yaml

sed -i 's~kind: Role~kind: {{ if .Values.clusterWide }}Cluster{{ end }}Role~' ${dir}/templates/operator-rbac.yaml


$YQ -i ".clusterWide=false" ${dir}/values.yaml
$YQ -i ".watchNamespace=null" ${dir}/values.yaml

sed -i -z 's~WATCH_NAMESPACE\n          valueFrom:\n            fieldRef:\n              fieldPath: metadata.namespace~WATCH_NAMESPACE\n          value: {{ if .Values.clusterWide }}"*"{{ else }}{{ .Values.watchNamespace | default (print .Release.Namespace) | quote }}{{ end }}~' \
${dir}/templates/deployment.yaml