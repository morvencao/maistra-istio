#!/bin/bash

# Copyright Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set +e

source common.sh

log "deploy v2 ratings system into mesh1 and mesh2"

oc1 logs -n mesh1-bookinfo deploy/ratings-v2 -f &
oc2 logs -n mesh2-bookinfo deploy/ratings-v2 -f &

log "manual steps to test:
  1. Open http://istio-ingressgateway-mesh2-system.apps.kiali-qez-49.maistra.upshift.redhat.com/productpage
  2. Refresh the page several times and observe requests hitting either the mesh1 or the mesh2 cluster.
"
