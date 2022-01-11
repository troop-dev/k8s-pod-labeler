# K8s pod labeler

This repo implelents a [MutatingAdmissionWebhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#mutatingadmissionwebhook) that can add missing labels and annotations to pods. The webhook only adds the given label or mutation if it is not already present on the pod definition.

### Deployment

TODO: talk about helm

### Local Dev

TODO: talk about earthly, etc

### Motivation

Troop uses the Linkerd service mesh in GKE, which requires a special annotation to inform the cluster autoscaler that it can evict pods despite them having a local volume.

### Resources

- https://github.com/slackhq/simple-kubernetes-webhook
- https://github.com/operator-framework/operator-sdk
- https://github.com/openshift/generic-admission-server
- https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#webhookclientconfig-v1-apiextensions-k8s-io
- https://github.com/morvencao/kube-mutating-webhook-tutorial



Troop settings
```yaml
imagePullSecrets:
- name: troop-dev-github-docker-registry

appInputs:
  labels:
    "cluster-autoscaler.kubernetes.io/safe-to-evict": "true"
  annotations:
    "cluster-autoscaler.kubernetes.io/safe-to-evict": "true"

```
