# Testing

## Live

Start minikube and run the Kubeflow pipelines api server
```shell
kubectl apply -k config/
```
```shell
export GO_KFP_API_SERVER_ADDRESS=$(minikube service kfp --url -n kubeflow | sed 's/http:\/\///g')
```

Open the Pipelines UI with

```shell
minikube service kfp-ui -n kubeflow
```