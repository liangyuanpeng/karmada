#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${REPO_ROOT}/hack/util.sh"

helm push _output/charts/karmada-chart-$1.tgz oci://docker.io/gcslyp
helm push _output/charts/karmada-operator-chart-$1.tgz oci://docker.io/gcslyp

util::signImage docker.io/gcslyp/karmada:$1
util::signImage docker.io/gcslyp/karmada-operator:$1
