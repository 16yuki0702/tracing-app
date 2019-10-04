#!/bin/bash

MAX=$1

oc delete dc gateway
oc delete svc gateway-svc
oc delete bc gateway
oc delete is gateway
oc delete configmap gateway

for i in `seq 1 $MAX`
do
  oc delete dc "dest$i"
  oc delete svc "dest$i-svc"
  oc delete configmap "dest$i-cfg"
done
