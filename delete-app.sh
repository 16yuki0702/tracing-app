#!/bin/bash

MAX=$1
USER=$2

oc delete dc "$USER-gateway"
oc delete svc "$USER-gateway"
#oc delete bc gateway
#oc delete is gateway
oc delete configmap "$USER-gateway-cfg"
oc delete gateway "$USER-gateway"
oc delete virtualservice "$USER-gateway"
oc delete destinationrule "$USER-gateway"

# delete gitlab runner
oc delete dc gitlab-runner
oc delete is gitlab-runner
oc delete sa gitlab-runner-user
oc delete configmap gitlab-runner-scripts
oc delete rolebinding gitlab-runner_edit

for i in `seq 1 $MAX`
do
  oc delete dc "dest$i"
  oc delete svc "dest$i"
  oc delete configmap "dest$i-cfg"
done
