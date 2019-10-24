#!/bin/bash

oc get dc gateway > /dev/null 2>&1

if [ $ret -e 0 ]; then
  oc apply -f gateway/app.yaml
  oc apply -f dest1/app.yaml
  oc apply -f dest2/app.yaml
  oc apply -f dest3/app.yaml
else
  oc rollout latest gateway
  oc rollout latest dest1
  oc rollout latest dest2
  oc rollout latest dest3
fi
