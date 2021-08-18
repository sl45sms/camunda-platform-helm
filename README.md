[![Community Extension](https://img.shields.io/badge/Community%20Extension-An%20open%20source%20community%20maintained%20project-FF4700)](https://github.com/camunda-community-hub/community)[![Lifecycle: Incubating](https://img.shields.io/badge/Lifecycle-Incubating-blue)](https://github.com/Camunda-Community-Hub/community/blob/main/extension-lifecycle.md#incubating-)[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# camunda-cloud-helm

The camunda-cloud helm mono-repo, contains and host all camunda-cloud related helm charts.

The charts can be accessed by adding the following HELM repo to your HELM setup:
```
> helm repo add camundacloud https://helm.camunda.io
> helm repo update
```

There are three main charts which are represented in the following image:
![HELM CHARTS](imgs/HelmChartImage.png)

You can consume each individual chart, or use the `zeebe-full-helm` chart which will install all the components, including an NGINX Ingress Controller.

Currently we host the following charts:
- [zeebe-cluster-helm](charts/zeebe-cluster-helm)
  - Depends on: [ElasticSearch](https://github.com/elastic/helm-charts/tree/master/elasticsearch), [Kibana](https://github.com/elastic/helm-charts/tree/master/kibana), [Prometheus Operator](https://github.com/helm/charts/tree/master/stable/prometheus-operator)
- [zeebe-operate-helm](charts/zeebe-operate-helm)
  - Can be configured to point to a Zeebe Cluster
- [zeebe-full-helm](charts/zeebe-full-helm)
  - Depends on: zeebe-cluster, zeebe-operate and [nginx-ingress](https://github.com/helm/charts/tree/master/stable/nginx-ingress)
- [zeebe-tasklist-helm](charts/zeebe-tasklist-helm) **(Experimental)**
- [zeebe-zeeqs-helm](charts/zeebe-zeeqs-helm) **(Experimental)**
  - Depends on: [Hazelcast](https://github.com/hazelcast/charts)

Follow [the instructions in the Zeebe docs](https://docs.zeebe.io/kubernetes/installing-helm.html) to install Zeebe to a K8s cluster using these charts.

Each Chart contains it's own configurations and parameters, you can visit each chart README for more information. 

**Note check the [Zeebe Helm Profiles](https://github.com/zeebe-io/zeebe-helm-profiles) repository for different configurations for your clusters, such as Dev, HA, etc. Feel free to contribute with your own profiles if you want to**

## Installing Charts

You can install these Helm Charts by running:
```
helm install --name <YOUR HELM RELEASE NAME> zeebe/zeebe-full-helm
```
This command will install all Zeebe Components that are provided by the `zeebe-full-helm` chart.


## Uninstalling Charts

You can remove these charts by running:
```
helm delete <YOUR HELM RELEASE NAME> --no-hooks --purge
```

> Notice that all the services and pods will be deleted, but not the Persistence Volume Claims which are used to hold the storage for the data generated by the cluster and ElasticSearch. In order to free up the storage you need to manually delete all the Persistent Volume Claims. You can do this by running:
```
kubectl get pvc
```
Then delete the ones that you don't want to keep:
```
kubectl delete pvc <PVC ids here>
```

Or delete the related kubernetes namespace, which contains the resources.

## Issues

Please create [new issues](https://github.com/zeebe-io/zeebe-helm/issues) if you find problems with these charts. This repository is hosted using GitHub Pages and the source code repository can be found here: [https://github.com/zeebe-io/zeebe-helm/](https://github.com/zeebe-io/zeebe-helm/)

## Releasing these Charts

These charts are build and linted on every push to the main branch. If the chart version changes a new github release with the the corresponding packaged helm chart is created. The charts are hosted via github pages and use the release artifacts. We use the [chart-releaser-action](https://github.com/helm/chart-releaser-action) to release the charts.