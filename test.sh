#!/bin/bash

while true
do
  curl http://istio-ingressgateway-istio-system.apps.openshift-1.16yuki0702dev1.mobi
  curl http://istio-ingressgateway-istio-system.apps.openshift-1.16yuki0702dev1.mobi/todest1
  curl http://istio-ingressgateway-istio-system.apps.openshift-1.16yuki0702dev1.mobi/todest2
  curl http://istio-ingressgateway-istio-system.apps.openshift-1.16yuki0702dev1.mobi/todest3

  sleep 5
done
