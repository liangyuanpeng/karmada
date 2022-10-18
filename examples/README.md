# Examples

## Resource Interpreter

This example implements a new CustomResourceDefinition(CRD), `Workload`, and creates a resource interpreter webhook.

### Install

### Prerequisites

For karmada deploy using `hack/local-up-karmada.sh`, there are `karmada-host`, `karmada-apiserver` and three member clusters named `member1`, `member2` and `member3`.

Then you need to deploy MetalLB as a Load Balancer to expose the webhook.

```bash
kubectl --context="karmada-host" get configmap kube-proxy -n kube-system -o yaml | \
  sed -e "s/strictARP: false/strictARP: true/" | \
  kubectl --context="karmada-host" apply -n kube-system -f -

curl https://raw.githubusercontent.com/metallb/metallb/v0.13.5/config/manifests/metallb-native.yaml -k | \
  sed '0,/args:/s//args:\n        - --webhook-mode=disabled/' | \
  sed '/apiVersion: admissionregistration/,$d' | \
  kubectl --context="karmada-host" apply -f -

export interpreter_webhook_example_service_external_ip_address=$(kubectl config view --template='{{range $_, $value := .clusters }}{{if eq $value.name "karmada-apiserver"}}{{$value.cluster.server}}{{end}}{{end}}' | \
  awk -F/ '{print $3}' | \
  sed 's/:.*//' | \
  awk -F. '{printf "%s.%s.%s.8",$1,$2,$3}')

cat <<EOF | kubectl --context="karmada-host" apply -f -
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: metallb-config
  namespace: metallb-system
spec:
  addresses:
  - ${interpreter_webhook_example_service_external_ip_address}-${interpreter_webhook_example_service_external_ip_address}
---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: metallb-advertisement
  namespace: metallb-system
EOF
```

#### Step1: Install `Workload` CRD in `karmada-apiserver` and member clusters

Install CRD in `karmada-apiserver` by running the following command:

```bash
kubectl --kubeconfig $HOME/.kube/karmada.config --context karmada-apiserver apply -f examples/customresourceinterpreter/apis/workload.example.io_workloads.yaml
```

Create `ClusterPropagationPolicy` object to propagate CRD to member clusters:

workload-crd-cpp.yaml:

<details>

<summary>unfold me to see the yaml</summary>

```yaml
apiVersion: policy.karmada.io/v1alpha1
kind: ClusterPropagationPolicy
metadata:
  name: workload-crd-cpp
spec:
  resourceSelectors:
    - apiVersion: apiextensions.k8s.io/v1
      kind: CustomResourceDefinition
      name: workloads.workload.example.io
  placement:
    clusterAffinity:
      clusterNames:
        - member1
        - member2
        - member3
```
</details>

```bash
kubectl --kubeconfig $HOME/.kube/karmada.config --context karmada-apiserver apply -f workload-crd-cpp.yaml
```

#### Step2: Deploy webhook configuration in karmada-apiserver

Execute below script:

webhook-configuration.sh

<details>

<summary>unfold me to see the script</summary>

```bash
#!/usr/bin/env bash

export ca_string=$(cat ${HOME}/.karmada/ca.crt | base64 | tr "\n" " "|sed s/[[:space:]]//g)
export temp_path=$(mktemp -d)
export interpreter_webhook_example_service_external_ip_address=$(kubectl config view --template='{{range $_, $value := .clusters }}{{if eq $value.name "karmada-apiserver"}}{{$value.cluster.server}}{{end}}{{end}}' | \
  awk -F/ '{print $3}' | \
  sed 's/:.*//' | \
  awk -F. '{printf "%s.%s.%s.8",$1,$2,$3}')

cp -rf "examples/customresourceinterpreter/webhook-configuration.yaml" "${temp_path}/temp.yaml"
sed -i'' -e "s/{{caBundle}}/${ca_string}/g" -e "s/{{karmada-interpreter-webhook-example-svc-address}}/${interpreter_webhook_example_service_external_ip_address}/g" "${temp_path}/temp.yaml"
kubectl --kubeconfig $HOME/.kube/karmada.config --context karmada-apiserver apply -f "${temp_path}/temp.yaml"
rm -rf "${temp_path}"
```

</details>

```bash
chmod +x webhook-configuration.sh

./webhook-configuration.sh
```

#### Step3: Deploy interpreter webhook example in karmada-host

```bash
kubectl --kubeconfig $HOME/.kube/karmada.config --context karmada-host apply -f examples/customresourceinterpreter/karmada-interpreter-webhook-example.yaml
```

### Usage

Create a `Workload` resource and propagate it to the member clusters:

workload-interpret-test.yaml:

<details>

<summary>unfold me to see the yaml</summary>

```yaml
apiVersion: workload.example.io/v1alpha1
kind: Workload
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  replicas: 3
  paused: false
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: nginx
        name: nginx
---
apiVersion: policy.karmada.io/v1alpha1
kind: PropagationPolicy
metadata:
  name: nginx-workload-propagation
spec:
  resourceSelectors:
    - apiVersion: workload.example.io/v1alpha1
      kind: Workload
      name: nginx
  placement:
    clusterAffinity:
      clusterNames:
        - member1
        - member2
        - member3
    replicaScheduling:
      replicaDivisionPreference: Weighted
      replicaSchedulingType: Divided
      weightPreference:
        staticWeightList:
          - targetCluster:
              clusterNames:
                - member1
            weight: 1
          - targetCluster:
              clusterNames:
                - member2
            weight: 1
          - targetCluster:
              clusterNames:
                - member3
            weight: 1
```

</details>

```bash
kubectl --kubeconfig $HOME/.kube/karmada.config --context karmada-apiserver apply -f workload-interpret-test.yaml
```

#### InterpretReplica

You can get `ResourceBinding` to check if the `replicas` field is interpreted successfully.

```bash
kubectl get rb nginx-workload -o yaml
```

#### ReviseReplica

You can check if the replicas field of `Workload` object is revised to 1 in all member clusters.

```bash
kubectl --kubeconfig $HOME/.kube/members.config --context member1 get workload nginx --template={{.spec.replicas}}
```

#### Retain

Update `spec.paused` of `Workload` object in member1 cluster to `true`.

```bash
kubectl --kubeconfig $HOME/.kube/members.config --context member1 patch workload nginx --type='json' -p='[{"op": "replace", "path": "/spec/paused", "value":true}]'
```

Check if it is retained successfully.
```bash
kubectl --kubeconfig $HOME/.kube/members.config --context member1 get workload nginx --template={{.spec.paused}}
```

#### InterpretStatus

There is no `Workload` controller deployed on member clusters, so in order to simulate the `Workload` CR handling, 
we will manually update `status.readyReplicas` of `Workload` object in member1 cluster to 1. 

```bash
kubectl proxy --port=8001 &
curl  http://127.0.0.1:8001/apis/workload.example.io/v1alpha1/namespaces/default/workloads/nginx/status  -XPATCH -d'{"status":{"readyReplicas": 1}}' -H "Content-Type: application/merge-patch+json
```

Then you can get `ResourceBinding` to check if the `status.aggregatedStatus[x].status` field is interpreted successfully.

```bash
kubectl get rb nginx-workload --kubeconfig $HOME/.kube/karmada.config --context karmada-apiserver -o yaml
```

You can also check the `status.manifestStatuses[x].status` field of Karmada `Work` object in namespace karmada-es-member1.

#### InterpretHealth

You can get `ResourceBinding` to check if the `status.aggregatedStatus[x].health` field is interpreted successfully.

```bash
kubectl get rb nginx-workload --kubeconfig $HOME/.kube/karmada.config --context karmada-apiserver -o yaml
```

You can also check the `status.manifestStatuses[x].health` field of Karmada `Work` object in namespace karmada-es-member1.

#### AggregateStatus

You can check if the `status` field of `Workload` object is aggregated correctly.

```bash
kubectl get workload nginx --kubeconfig $HOME/.kube/karmada.config --context karmada-apiserver -o yaml
```
 

> Note: If you want to use `Retain`/`InterpretStatus`/`InterpretHealth` function in pull mode cluster, you need to deploy interpreter webhook example in this member cluster.