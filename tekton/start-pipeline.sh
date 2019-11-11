#!/bin/bash

oc apply -f tasks.yaml
oc apply -f pipeline-resources.yaml
oc apply -f pipeline.yaml
oc apply -f pipeline-run.yaml
