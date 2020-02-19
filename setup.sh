#!/bin/bash

USER=$1

sed -i -e "s/sample-tracing/$USER-tracing/g" gateway/app.yaml
sed -i -e "s/sample-gateway/$USER-gateway/g" gateway/app.yaml
sed -i -e "s/sample-tracing/$USER-tracing/g" dest1/app.yaml
sed -i -e "s/sample-tracing/$USER-tracing/g" dest1/canary.yaml
sed -i -e "s/sample-tracing/$USER-tracing/g" dest2/app.yaml
sed -i -e "s/sample-tracing/$USER-tracing/g" dest3/app.yaml

oc project "$USER-tracing"

HOSTNAME=`oc get route istio-ingressgateway -n "$USER-smcp" | awk 'NR==2{print $2}'`

sed -i -e "s/gateway/$HOSTNAME/g" test.sh
