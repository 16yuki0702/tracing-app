#/bin/bash

oc adm policy add-scc-to-user privileged -z tekton-triggers-admin
oc adm policy add-role-to-user edit -z tekton-triggers-admin
