# K8s pod labeler

This repo implelents a [MutatingAdmissionWebhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#mutatingadmissionwebhook) that can add missing labels and annotations to pods. The webhook only adds the given label or mutation if it is not already present on the pod definition.

TODO:
 - add unit tests!
 - support external cert generation like cert-manager
### Configuration

Configure the desired labels/annotations for your pods by setting the `appInputs.labels` and `appInputs.annotations` values.

```yaml
appInputs:
  labels:
    "foo": "bar"
  annotations:
    "bar": "baz"
```

By default, this uses namespace selectors to enable the labeler in a given namespace by adding a namespace label `troop.com/k8s-pod-labeler: enabled`.

For all configurable settings, see the [values.yaml](./helm/k8s-pod-labeler/values.yaml) file.
### Deployment

Helm install...
```bash
# this enables pulling helm charts from github packages
export HELM_EXPERIMENTAL_OCI=1
# test your local config at local-values.yml
helm template \
  -f ./local-values.yml \
  --version 0.3.6 \
  oci://ghcr.io/troop-dev/helm/k8s-pod-labeler
# install helm from a local-values.yml config
helm install \
  -f ./local-values.yml \
  pod-labeler \
  --version 0.3.6 \
  oci://ghcr.io/troop-dev/helm/k8s-pod-labeler
```

### Local Dev

This repo uses [Earthly](https://earthly.dev/) to build the golang application.

```bash
# install earthly on osx with homebrew
brew install earthly/earthly/earthly && earthly bootstrap
# build the docker image
earthly +build
```

### Motivation

Troop uses the Linkerd service mesh in GKE, which requires a special annotation to inform the cluster autoscaler that it can evict pods despite them having a local volume.

### Resources

- https://github.com/slackhq/simple-kubernetes-webhook
- https://github.com/operator-framework/operator-sdk
- https://github.com/openshift/generic-admission-server
- https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#webhookclientconfig-v1-apiextensions-k8s-io
- https://github.com/morvencao/kube-mutating-webhook-tutorial
