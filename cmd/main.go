package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/logging/logadmin"
	"github.com/olekukonko/tablewriter"
	flag "github.com/spf13/pflag"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/cloud/audit"
)

func main() {
	projectID := flag.String("projectID", "", "GCP Project ID")
	clusterName := flag.String("clusterName", "", "GKE Cluster Name")
	sinceDays := flag.Int64("sinceDays", 2, "How many days back to query for?")
	flag.Parse()

	if *projectID == "" || *clusterName == "" {
		log.Fatal("required parameter is missing")
	}

	ctx := context.Background()
	client, err := logadmin.NewClient(ctx, *projectID)
	if err != nil {
		log.Fatalf("failed creating log client: %v", err)
	}

	fmt.Printf("Showing GKE Master and worker node upgrades history of cluster [%v] at project ID [%v] for the last [%v] days. Please wait...\n", *clusterName, *projectID, *sinceDays)
	fmt.Println("--- GKE Master Upgrades ---")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetColWidth(60)
	table.SetHeader([]string{"Time", "Previous Master Version", "Current Master Version"})

	durationSince := time.Duration(*sinceDays) * 24 * time.Hour
	sinceTs := time.Now().Add(-durationSince).Format(time.RFC3339)
	logFilter := fmt.Sprintf("resource.labels.cluster_name=\"%s\" AND resource.type=\"gke_cluster\" AND\nlog_id(\"cloudaudit.googleapis.com/activity\") \nAND timestamp > \"%s\"", *clusterName, sinceTs)

	it := client.Entries(ctx, logadmin.Filter(logFilter))
	for {
		entry, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("failed using iterator: %v", err)
		}

		auditLogEntry := entry.Payload.(*audit.AuditLog)
		if auditLogEntry.MethodName == "google.container.v1.ClusterManager.UpdateCluster" && auditLogEntry.GetMetadata().GetFields() != nil {
			currentMasterVersion := auditLogEntry.Metadata.Fields["currentMasterVersion"]
			previousMasterVersion := auditLogEntry.Metadata.Fields["previousMasterVersion"]

			table.Append([]string{entry.Timestamp.String(), previousMasterVersion.GetStringValue(), currentMasterVersion.GetStringValue()})
		}

		if auditLogEntry.MethodName == "google.container.v1.ClusterManager.UpdateNodePool" && auditLogEntry.GetRequest().GetFields() != nil {
			nodeVersion := auditLogEntry.Request.Fields["nodeVersion"]
			fmt.Printf("Found GKE node pool upgrade. time: %v, currentVersion: %v\n", entry.Timestamp.String(), nodeVersion.GetStringValue())
		}
	}

	table.Render()

	fmt.Println("--- GKE Node Pool Upgrades ---")

	table = tablewriter.NewWriter(os.Stdout)
	table.SetColWidth(60)
	table.SetHeader([]string{"Time", "Current Node Version"})

	logFilter = fmt.Sprintf("resource.labels.cluster_name=\"%s\" AND log_id(\"cloudaudit.googleapis.com/activity\") AND\nresource.type=\"gke_nodepool\" \nAND timestamp > \"%s\"", *clusterName, sinceTs)
	it = client.Entries(ctx, logadmin.Filter(logFilter))
	for {
		entry, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("failed using iterator: %v", err)
		}

		auditLogEntry := entry.Payload.(*audit.AuditLog)
		if auditLogEntry.MethodName == "google.container.v1.ClusterManager.UpdateNodePool" && auditLogEntry.GetRequest().GetFields() != nil {
			nodeVersion := auditLogEntry.Request.Fields["nodeVersion"]

			table.Append([]string{entry.Timestamp.String(), nodeVersion.GetStringValue()})

		}
	}

	table.Render()
}
