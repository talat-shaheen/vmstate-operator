# vmstate-operator
// TODO(user): Add simple overview of use/purpose

## Description
// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.

**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

**Pre-requisite** Need to login to any image registry and replace registry in the command below & create a secret in the operator namespace with AWS environment variables
```
kubectl create secret generic aws-secret --from-literal=region=us-east-1 --from-literal=aws-secret-access-key=<secret access key> --from-literal=aws-access-key-id=<secret access key id>
```
### Running on the cluster
1. Build and push your image:
	
```sh
make generate;make manifests;
make docker-build;
sudo docker push quay.io/talat_shaheen0/vmstate-operator:latest;
```
OR

```
make generate;make manifests;
make docker-build docker-push IMG="quay.io/talat_shaheen0/vmstate-operator:latest"
make deploy IMG="quay.io/talat_shaheen0/vmstate-operator:latest"
```
	
2. Deploy the controller to the cluster:

```sh
make deploy quay.io/talat_shaheen0/vmstate-operator:tag
```

3. Apply Custom Resources:

```sh
kubectl apply -f config/samples/awsmanager_v1_awsmanager.yaml -n vmstate-operator-system;
kubectl apply -f config/samples/aws_v1_awsec2.yaml -n vmstate-operator-system;
```
3. Check jobs & AWSEC2:

```
kubectl get jobs -n vmstate-operator-system;
kubectl get awsec2 -n vmstate-operator-system;
```

4. Delete CR:

```
kubectl delete awsec2 <xyz> -n vmstate-operator-system;
```
### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller to the cluster:

```sh
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) 
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster 

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

