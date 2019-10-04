#!/bin/bash

oc apply -f gateway/app.yaml
oc apply -f dest1/app.yaml
oc apply -f dest2/app.yaml
oc apply -f dest3/app.yaml
