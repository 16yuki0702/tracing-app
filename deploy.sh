#!/bin/bash

ret=`oc get dc gateway 2>&1`
if [[ $ret =~ NAME ]]; then
  # for resource limits
  oc delete dc gateway
  oc delete dc dest1
  oc delete dc dest2
  oc delete dc dest3
fi
oc apply -f gateway/app.yaml
oc apply -f dest1/app.yaml
oc apply -f dest2/app.yaml
oc apply -f dest3/app.yaml
