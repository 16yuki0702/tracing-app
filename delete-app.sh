#!/bin/bash

MAX=$1

oc delete dc gateway
oc delete svc gateway
oc delete bc gateway
oc delete is gateway
oc delete configmap gateway-cfg
oc delete gateway tracing-gateway
oc delete virtualservice tracing-app
oc delete destinationrule gateway

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
