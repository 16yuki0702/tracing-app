#!/bin/bash

while true
do
  curl http://istio-ingressgateway-istio-system.apps.openshift-1.16yuki0702dev1.mobi/productpage
  curl http://istio-ingressgateway-istio-system.apps.openshift-1.16yuki0702dev1.mobi
  curl http://istio-ingressgateway-istio-system.apps.openshift-1.16yuki0702dev1.mobi/trace1
  curl http://istio-ingressgateway-istio-system.apps.openshift-1.16yuki0702dev1.mobi/trace2
  curl http://istio-ingressgateway-istio-system.apps.openshift-1.16yuki0702dev1.mobi/trace3

  sleep 10
done
