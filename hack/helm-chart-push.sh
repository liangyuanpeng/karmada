#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${REPO_ROOT}/hack/util.sh"

helm push _output/charts/karmada-chart-$1.tgz oci://docker.io/lypgcs
helm push _output/charts/karmada-operator-chart-$1.tgz oci://docker.io/lypgcs

signImage docker.io/lypgcs/karmada:$1
signImage docker.io/lypgcs/karmada-operator:$1