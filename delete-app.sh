#!/bin/bash

oc delete dc gateway
oc delete svc gateway-svc
oc delete bc gateway
oc delete is gateway

oc delete dc dest1
oc delete svc dest1
oc delete bc dest1
oc delete is dest1

oc delete dc dest2
oc delete svc dest2
oc delete bc dest2
oc delete is dest2
