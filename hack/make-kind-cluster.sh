#! /usr/bin/env bash

set -o pipefail
set -o nounset
set -o errexit

readonly PROGNAME=$(basename $0)

readonly KIND=${KIND:-kind}
readonly KUBECTL=${KUBECTL:-kubectl}
readonly KUSTOMIZE=${KUSTOMIZE:-kustomize}
readonly DOCKER=${DOCKER:-docker}

readonly HERE=$(cd $(dirname "$0") && pwd)
readonly REPO=$(cd "${HERE}/.." && pwd)

usage() {
    echo usage: $PROGNAME [CLUSTER]
    exit 64 # EX_USAGE
}

case $# in
0) ;;
1) CLUSTER="$1";;
*) usage;;
esac

readonly CLUSTER=${CLUSTER:-"envoy-controller"}

kind::cluster::list() {
    ${KIND} get clusters
}

# Emit a Kind config that maps the envoy listener ports to the host.
kind::cluster::config() {
    cat <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    listenAddress: "0.0.0.0"
  - containerPort: 443
    hostPort: 443
    listenAddress: "0.0.0.0"
EOF
}

kind::cluster::create() {
    ${KIND} create cluster \
        --config <(kind::cluster::config) \
        --name ${CLUSTER}
}

kind::cluster::load() {
    local -r img="$1"

    ${DOCKER} pull "$img"
    ${KIND} load docker-image \
        --name "${CLUSTER}" \
        "$img"
}

kind::cluster::exists() {
    ${KIND} get clusters | grep -q "$1"
}


# Create a fresh kind cluster.
if ! kind::cluster::exists "$CLUSTER" ; then
    kind::cluster::create
fi

kind::cluster::load "docker.io/envoyproxy/envoy:v1.15.0"
kind::cluster::load "docker.io/agervais/ingress-conformance-echo:latest"

${KUBECTL} apply -f <(${KUSTOMIZE} build ${REPO}/config/crd)
