#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
source "${REPO_ROOT}/hack/util.sh"

helm push _output/charts/karmada-chart-$1.tgz oci://ghcr.io/liangyuanpeng
helm push _output/charts/karmada-operator-chart-$1.tgz oci://ghcr.io/liangyuanpeng

util::signImage ghcr.io/liangyuanpeng/karmada:$1
util::signImage ghcr.io/liangyuanpeng/karmada-operator:$1
