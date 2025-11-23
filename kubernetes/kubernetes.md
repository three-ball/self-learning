# Kubernetes

## Overview

- Kubernetes is a platform to orchestrate the deployment, scaling, and management of container-based applications.

### Cluster Architecture

![alt text](images/architecture.png)

#### Cluster

| A cluster is a collection of hosts (nodes) that run containerized applications.

**Node**: A node is a single host in the cluster, which can be a physical or virtual machine. Its job is to run pods. Each node contains the services necessary to run pods and is managed by the master components.
- `kubelet`: An agent that runs on each node in the cluster. It ensures that containers are running in a pod. The kubelet takes a set of PodSpecs that are provided through various mechanisms and ensures that the containers described in those PodSpecs are running and healthy.
- `kube-proxy` (optional): A network proxy that runs on each node in the cluster. It maintains network rules on nodes, allowing network communication to your pods from network sessions inside or outside of your cluster. We can also use other networking solutions instead of kube-proxy, such as Cilium or Calico.
- `Container Runtime`: The container runtime is the software that is responsible for running containers. Kubernetes supports several container runtimes, including Docker, containerd, and CRI-O.

**Pod**: Smallest deployable unit—one or more containers sharing network and storage. Each pod contains one or more containers that share the same network namespace and storage.

#### Control Plane

| The brain of the cluster that makes scheduling decisions and maintains desired state. 

The control plane is responsible for **managing the state of the cluster**, including scheduling, scaling, and updating applications. It consists of several components, including the API server, scheduler, controller manager, and `etcd` (a key-value store for cluster data).
- `kube-apiserver`: The API server is a component of the kubernetes control plane that exposes the Kubernetes API. It is the front-end for the Kubernetes control plane. Node components communicate with the API server to manage the state of the cluster.
- `etcd`: Consistent and highly-available key value store used as Kubernetes' backing store for all cluster data.
- `kube-scheduler`: The scheduler is a component of the Kubernetes control plane that is responsible for scheduling pods onto nodes. It watches for newly created pods that have no node assigned, and selects a node for them to run on. Scheduler just decides on which Node new Pod should be scheduled.
`kube-controller-manager`: The controller manager is a component of the Kubernetes control plane that runs controller processes. Each controller is a separate process, but to reduce complexity, they are all compiled into a single binary and run in a single process.

### Some definitions

![alt text](images/overviews_defi.png)

#### Service

| Stable network endpoint to access a set of pods.

Services are used to expose some functionality to users or other services. You can have services that provide access to external resources, or pods you control directly at the virtual IP level. Native Kubernetes services are exposed through convenient endpoints. Note that services operate at layer 3 (TCP/UDP).
There are different types of services:
- `ClusterIP` (default): Internal clients send requests to a stable internal IP address.
- `NodePort`: Clients send requests to the IP address of a node on one or more nodePort values that are specified by the Service.
- `LoadBalancer`: Clients send requests to the IP address of a network load balancer.
- `ExternalName`: Internal clients use the DNS name of a Service as an alias for an external DNS name.
- `Headless`: You can use a headless service when you want a Pod grouping, but don't need a stable IP address.

#### Volume

| Storage that persists beyond pod lifetime and can be shared by containers.

Local storage used by the pod is ephemeral and goes away with the pod in most cases. Sometimes that's all you need, if the goal is just to exchange data between containers of the node, but sometimes it's important for the data to outlive the pod, or it's necessary to share data between pods. The volume concept supports that need.

Originally, Kubernetes directly supported many volume types, but the modern approach for extending Kubernetes with volume types is through the Container Storage Interface (CSI).

K8s doesn't manage data persistence.

#### PersistentVolume (PV) & PersistentVolumeClaim (PVC)

| PV: Cluster storage resource; PVC: Request for storage by a user.

![alt text](images/pv_pvc.png)

The `PersistentVolume` subsystem provides an API for users and administrators that abstracts details of how storage is provided from how it is consumed. To do this, we introduce two new API resources: `PersistentVolume` and `PersistentVolumeClaim`.
**A PersistentVolume (PV) is a piece of storage in the cluster that has been provisioned by an administrator or dynamically provisioned using Storage Classes.** It is a resource in the cluster just like a node is a cluster resource. PVs are volume plugins like Volumes, but have **a lifecycle independent of any individual Pod** that uses the PV.
**A PersistentVolumeClaim (PVC) is a request for storage by a user.** It is similar to a Pod. Pods consume node resources and PVCs consume PV resources.

#### ReplicaSet

| Ensures a specified number of identical pod replicas are running.

A ReplicaSet's purpose is to maintain a stable set of replica Pods running at any given time. Usually, you define a Deployment and let that Deployment manage ReplicaSets automatically.

#### StatefulSet

| Like ReplicaSet but for stateful apps needing stable identity and persistent storage.

A StatefulSet runs a group of Pods, and maintains a sticky identity for each of those Pods. This is useful for managing applications that need persistent storage or a stable, unique network identity. **Each pod gets a predictable name (pod-0, pod-1) and its own PVC. Used for databases, message queues—anything requiring stable identity.**

StatefulSet is the workload API object used to manage stateful applications.

#### Secret

| Stores sensitive data (passwords, tokens, keys) in base64 encoding.

Secrets are small objects that contain sensitive info such as credentials and tokens. They are stored by default as plaintext in etcd, accessible by the Kubernetes API server, and can be mounted as files into pods (using dedicated secret volumes that piggyback on regular data volumes) that need access to them. The same secret can be mounted into multiple pods.

#### ConfigMap

| Stores non-sensitive configuration data as key-value pairs.

Separates config from container images. Mount as files or inject as env vars. Changes don't auto-reload in running pods. **Do not put credentials in ConfigMap**.

#### Namespace

| Virtual cluster for resource isolation and multi-tenancy.

In Kubernetes, namespaces provide a mechanism for isolating groups of resources within a single cluster. Names of resources need to be unique within a namespace, but not across namespaces. Namespace-based scoping is applicable only for namespaced objects (e.g. Deployments, Services, etc.) and not for cluster-wide objects (e.g. StorageClass, Nodes, PersistentVolumes, etc.).
By default, a Kubernetes cluster has three initial namespaces:
- `default`: The default namespace for objects with no other namespace.
- `kube-system`: The namespace for objects created by the Kubernetes system.
- `kube-public`: This namespace is created automatically and is readable by all users (including those not authenticated). It is mostly reserved for cluster usage, in case that some resources should be visible and readable publicly throughout the cluster.
If you don't specify a namespace, objects are created in the `default` namespace.

#### Ingress

| Manages external access (HTTP/HTTPS) to services in a cluster, typically via load balancer.

```mermaid
graph LR;
  client([client])-. Ingress-managed <br> load balancer .->ingress[Ingress];
  ingress-->|routing rule|service[Service];
  subgraph cluster
  ingress;
  service-->pod1[Pod];
  service-->pod2[Pod];
  end
  classDef plain fill:#ddd,stroke:#fff,stroke-width:4px,color:#000;
  classDef k8s fill:#326ce5,stroke:#fff,stroke-width:4px,color:#fff;
  classDef cluster fill:#fff,stroke:#bbb,stroke-width:2px,color:#326ce5;
  class ingress,service,pod1,pod2 k8s;
  class client plain;
  class cluster cluster;
```

Make your HTTP (or HTTPS) network service available using a protocol-aware configuration mechanism, that understands web concepts like URIs, hostnames, paths, and more. The Ingress concept lets you map traffic to different backends based on rules you define via the Kubernetes API.

#### Helm

| Package manager for Kubernetes, simplifying app deployment and management.

Helm uses a packaging format called charts. A chart is a collection of files that describe a related set of Kubernetes resources. Helm charts help you define, install, and upgrade even the most complex Kubernetes applications. Examples: You can define a set of yaml files that describe the way to deploy an Elasticsearch stack, but with Helm, you can package those files into a chart and share it with others or pull charts created by the community.
Example charts repository: [Artifact Hub](https://artifacthub.io/), we use Helm to install applications like Prometheus, Grafana, Nginx Ingress Controller:

```bash
# add prometheus community repo and install prometheus
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
# update repo to get latest charts
helm repo update
# install prometheus
helm install prometheus prometheus-community/prometheus
```

### Configuration File

![alt text](images/conf_file.png)

First line is **apiVersion**, we can use `kubectl api-resources` to find the correct version for any resource.

| apiVersion | Resource Type |
|------------|----------------|
| v1         | Pod, Service, ConfigMap, Secret, Namespace, PersistentVolume, PersistentVolumeClaim |
| apps/v1    | Deployment, StatefulSet, DaemonSet, ReplicaSet |
| batch/v1   | Job, CronJob |
| networking.k8s.io/v1 | Ingress, NetworkPolicy |
| rbac.authorization.k8s.io/v1 | Role, ClusterRole, RoleBinding, ClusterRoleBinding |
| autoscaling/v1 | HorizontalPodAutoscaler |
| storage.k8s.io/v1 | StorageClass, VolumeAttachment |
| policy/v1 | PodDisruptionBudget |

Second line usually is what we want to create; Example: `kind: Deployment`
The remain parts is include 03 parts:
- Metadata
- Specification
    - `template`: template also has it's own metadata and spec section because this template configuration will apply to pod; so pod should have it's own metadata and spec. **This will be the "blueprint" for a pod.**
    - `selector`: Selectors use `labels` to connect Services to Pods dynamically. See: [Labels and Selectors](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/)
- Status (usually omitted in manifest files, auto generated by kubernetes)

```yaml
apiVersion: apps/v1
kind: Deployment
# Metadata for the Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx # ← [1] Label for the DEPLOYMENT object itself
# Specification of the Deployment    
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx # ← [2] Deployment manages pods with THIS label
  template:
    metadata:
      labels:
        app: nginx # ← [3] Stamps this label on every pod it creates
    spec:
      containers:
      - name: nginx
        image: nginx:1.16
        ports:
        - containerPort: 8080
# Can be status section, but usually omitted in manifest files

apiVersion: v1
kind: Service
# Metadata for the Service
metadata:
  name: nginx-service
# Specification of the Service
spec:
  selector:
    app: nginx # ← Service finds ALL pods with this label
  ports:
    - protocol: TCP
      port: 80 # ← Service listens on port 80
      targetPort: 8080 # ← Forwards to pod's port 8080
# Status section, usually omitted in manifest files

```