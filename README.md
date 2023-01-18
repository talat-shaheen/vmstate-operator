# vmstate-operator
This Operator will manage the state of cloud resources from kubernetes environment. Right now this only supports VMs on AWS. In future, other cloud services can be included and for various other cloud platforms like GCP, Azure etc.

## Description
There are two CRs that need to created.
First CR will be the manager pod that will keep on watching the state of the cloud resources and take corrective action if that does not match. 
Second CR will create the cloud resource that is supposed to be managed. Upon deletion of the CR, the resource is suposed to be deleted.

![image](https://user-images.githubusercontent.com/36874355/212598021-01716c9f-ea1a-4f11-b106-777781de06f0.png)

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.

**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

Install operator-sdk & golang.

Need to login to any image registry and replace registry in the command below & create a secret in the operator namespace with AWS environment variables

```
git clone <your-repo>
cd <your-repo>
git branch -b <branch name>
go mod init github.com/<your-repo>
go mod tidy
operator-sdk init --domain <your-name>.com --repo github.com/<your-repo>
operator-sdk edit --multigroup=true
operator-sdk create api     --group=Azure     --version=v1     --kind=<your-name>AzureAVM
operator-sdk create api     --group=gcp     --version=v1     --kind=<your-name>GCPGCE
operator-sdk create api     --group=aws     --version=v1     --kind=<your-name>AWSEC2
git add *
git commit -m"...."
git push origin <branch name>
```



### Running on the cluster
1. Build and push your image:
	
```sh
make generate;make manifests;
make docker-build;
sudo docker push quay.io/<your-registry>/vmstate-operator:latest;
```
OR

```
make generate;make manifests;
make docker-build docker-push IMG="quay.io/<your-registry>/vmstate-operator:latest"
make deploy IMG="quay.io/<your-registry>/vmstate-operator:latest"
```
	
2. Deploy the controller to the cluster:

```sh
make deploy IMG="quay.io/<your-registry>/vmstate-operator:tag"
```

3. Apply Custom Resources & create secret:

```sh
kubectl create secret generic aws-secret --from-literal=region=us-east-1 --from-literal=aws-secret-access-key=<secret access key> --from-literal=aws-access-key-id=<secret access key id>
kubectl apply -f config/samples/awsmanager_v1_awsmanager.yaml -n vmstate-operator-system;
kubectl apply -f config/samples/aws_v1_awsec2.yaml -n vmstate-operator-system;
```
3. Check jobs & AWSEC2/AWSManager resources:

```
kubectl get jobs -n vmstate-operator-system;
kubectl get awsec2 -n vmstate-operator-system;
kubectl get awsmanager -n vmstate-operator-system;
```

4. Delete CR:

```
kubectl delete awsec2 <cr-name> -n vmstate-operator-system;
kubectl delete awsmanager <cr-name> -n vmstate-operator-system;
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

