#!/bin/bash

MAX=$1

oc delete dc gateway
oc delete svc gateway-svc
oc delete bc gateway
oc delete is gateway
oc delete configmap gateway
oc delete gateway tracing-gateway
oc delete virtualservice tracing-app
oc delete destinationrule gateway

for i in `seq 1 $MAX`
do
  oc delete dc "dest$i"
  oc delete svc "dest$i"
  oc delete configmap "dest$i-cfg"
done
