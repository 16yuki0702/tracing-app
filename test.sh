#!/bin/bash

while true
do
  curl http://istio-ingressgateway-istio-system.apps.openshift-1.16yuki0702dev1.mobi/productpage
  curl http://istio-ingressgateway-istio-system.apps.openshift-1.16yuki0702dev1.mobi
  curl http://istio-ingressgateway-istio-system.apps.openshift-1.16yuki0702dev1.mobi/propagate1
  curl http://istio-ingressgateway-istio-system.apps.openshift-1.16yuki0702dev1.mobi/propagate2
  curl http://istio-ingressgateway-istio-system.apps.openshift-1.16yuki0702dev1.mobi/propagate3

  sleep 30
done
