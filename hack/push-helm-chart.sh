#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail



function cosignImage(){

echo "====================begin cosign for :"$1

echo "github.sha:" ${{ github.sha }}
echo "github.run_id:" ${{ github.run_id }}
echo "github.run_attempt:" ${{ github.run_attempt }}

cosign sign --yes \
            -a sha=${{ github.sha }} \
            -a run_id=${{ github.run_id }} \
            -a run_attempt=${{ github.run_attempt }} \
            $1
}


helm push _output/charts/karmada-chart-${VERSION}.tgz oci://$REGISTRY
cosignImage $REGISTRY:$VERSION
