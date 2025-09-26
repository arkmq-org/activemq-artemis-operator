#!/bin/bash

version=$1
dir=$2
YQ=${3:-yq}

# version
$YQ -i ".version=\"${version}\"" ${dir}/Chart.yaml
$YQ -i ".appVersion=\"${version}\"" ${dir}/Chart.yaml

# cluster scoped, namespaces
sed -i 's~kind: Role~kind: {{ if .Values.clusterScoped }}Cluster{{ end }}Role~' ${dir}/templates/operator-rbac.yaml
$YQ -i ".clusterScoped=true" ${dir}/values.yaml
$YQ -i ".watchNamespaces=[]" ${dir}/values.yaml
sed -i -z 's~WATCH_NAMESPACE\n          valueFrom:\n            fieldRef:\n              fieldPath: metadata.namespace~WATCH_NAMESPACE\n          value: {{ if .Values.clusterScoped }}{{ join "," .Values.watchNamespaces | quote }}{{ else }}{{ .Release.Namespace | quote }}{{ end }}~' \
${dir}/templates/deployment.yaml

# related image prefix
$YQ -i '.controllerManager.manager.relatedImages.activemqArtemisBrokerInitRepository="quay.io/arkmq-org/activemq-artemis-broker-init"' ${dir}/values.yaml
$YQ -i '.controllerManager.manager.relatedImages.activemqArtemisBrokerKubernetesRepository="quay.io/arkmq-org/activemq-artemis-broker-kubernetes"' ${dir}/values.yaml
$YQ -i '.controllerManager.manager.relatedImages += (.controllerManager.manager.env | with_entries(select(.key | test("^relatedImageA.+"))) | with_entries(.key |= sub("relatedImageA","a")) | with_entries(.value |= {"digest": sub("quay.io.*@","")}))' ${dir}/values.yaml
$YQ -i 'del(.controllerManager.manager.env)' ${dir}/values.yaml
sed -i -E 's~\.Values\.controllerManager\.manager\.env\.relatedImageA(ctivemqArtemisBrokerInit.+)~(printf "%s%s" .Values.controllerManager.manager.relatedImages.activemqArtemisBrokerInitRepository .Values.controllerManager.manager.relatedImages.a\1.digest)~' ${dir}/templates/deployment.yaml
sed -i -E 's~\.Values\.controllerManager\.manager\.env\.relatedImageA(ctivemqArtemisBrokerKubernetes.+)~(printf "%s%s" .Values.controllerManager.manager.relatedImages.activemqArtemisBrokerKubernetesRepository .Values.controllerManager.manager.relatedImages.a\1.digest)~' ${dir}/templates/deployment.yaml

# pack CRDs as template
crds="${dir}/templates/crds.yaml"
echo "{{- if .Values.crds.apply }}" > $crds
for file in ${dir}/crds/*.yaml; do
  $YQ -i '.metadata.annotations."helm.sh/resource-policy"="{{ if .Values.crds.keep }}keep{{ else }}delete{{ end }}"' $file
  cat $file >> $crds
  echo "---" >> $crds
done
echo "{{- end }}" >> $crds
rm -R ${dir}/crds/
$YQ -i ".crds.apply=true" ${dir}/values.yaml
$YQ -i ".crds.keep=true" ${dir}/values.yaml

# resources
$YQ -i '.controllerManager.manager.resources=null' ${dir}/values.yaml
sed -i -z 's~name: manager~name: manager\n        resources: {{- toYaml .Values.controllerManager.manager.resources | nindent 10 }}~' ${dir}/templates/deployment.yaml

# ENABLE_WEBHOOKS always false
$YQ -i 'del(.controllerManager.manager.env.enableWebhooks)' ${dir}/values.yaml
sed -i "s~.Values.controllerManager.manager.env.enableWebhooks~false~" ${dir}/templates/deployment.yaml
