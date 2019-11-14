#!/bin/bash

ret=`oc get dc gateway 2>&1`
if [[ $ret =~ NAME ]]; then
  oc rollout latest gateway
  oc rollout latest dest1
  oc rollout latest dest2
  oc rollout latest dest3
else
  oc apply -f gateway/app.yaml
  oc apply -f dest1/app.yaml
  oc apply -f dest2/app.yaml
  oc apply -f dest3/app.yaml
fi
