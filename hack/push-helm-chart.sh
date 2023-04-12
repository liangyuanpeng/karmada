#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

helm push _output/charts/karmada-chart-${VERSION}.tgz oci://$REGISTRY