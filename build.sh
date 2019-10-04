#!/bin/bash

oc new-build --strategy=docker --binary=true --name=gateway
oc start-build gateway --from-dir=. --from-file= --follow
