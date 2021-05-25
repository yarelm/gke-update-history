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


Example output:
```
./gke-update-history --projectID my-proj --clusterName cluster-1
Showing GKE Master and worker node upgrades history of cluster [cluster-1] at project ID [my-proj] for the last [2] days. Please wait...
--- GKE Master Upgrades ---
+-----------------------------------------+-------------------------+------------------------+
|                  TIME                   | PREVIOUS MASTER VERSION | CURRENT MASTER VERSION |
+-----------------------------------------+-------------------------+------------------------+
| 2021-05-25 08:58:20.213613341 +0000 UTC | 1.17.17-gke.4900        | 1.18.17-gke.700        |
| 2021-05-25 12:38:29.825690238 +0000 UTC | 1.18.17-gke.700         | 1.18.17-gke.1200       |
+-----------------------------------------+-------------------------+------------------------+
--- GKE Node Pool Upgrades ---
+-----------------------------------------+----------------------+
|                  TIME                   | CURRENT NODE VERSION |
+-----------------------------------------+----------------------+
| 2021-05-25 09:12:50.541920125 +0000 UTC | 1.18.17-gke.700      |
| 2021-05-25 12:40:25.564877517 +0000 UTC | 1.18.17-gke.1200     |
+-----------------------------------------+----------------------+
```