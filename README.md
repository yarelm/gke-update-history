# gke-update-history
A quick glance to the update history of a GKE cluster master and worker nodes

# Installation
Fetch the latest [release](https://github.com/yarelm/gke-update-history/releases), extract it, and run it!

Usage:

```
./gke-update-history -h
      --clusterName string   GKE Cluster Name
      --projectID string     GCP Project ID
      --sinceDays int        How many days back to query for? (default 2)
```