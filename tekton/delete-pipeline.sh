#!/bin/bash

oc delete pipeline tracing-pipeline
oc delete pipelinerun tracing-pipeline-run
oc delete task oc
oc delete task build-and-push
oc delete pipelineresource pipeline-source
oc delete pipelineresource pipeline-image
