#!/bin/bash

oc delete route gateway
oc delete dc gateway
oc delete svc gateway
