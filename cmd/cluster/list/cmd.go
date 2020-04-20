/*
Copyright (c) 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package list

import (
	"fmt"
	"os"
	"text/tabwriter"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/spf13/cobra"

	"gitlab.cee.redhat.com/service/moactl/pkg/aws"
	"gitlab.cee.redhat.com/service/moactl/pkg/logging"
	"gitlab.cee.redhat.com/service/moactl/pkg/ocm"
	"gitlab.cee.redhat.com/service/moactl/pkg/ocm/properties"
	rprtr "gitlab.cee.redhat.com/service/moactl/pkg/reporter"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List clusters",
	Long:  "List clusters.",
	Run:   run,
}

func run(_ *cobra.Command, argv []string) {
	// Create the reporter:
	reporter, err := rprtr.New().
		Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't create reporter: %v\n", err)
		os.Exit(1)
	}

	// Create the logger:
	logger, err := logging.NewLogger().Build()
	if err != nil {
		reporter.Errorf("Can't create logger: %v", err)
		os.Exit(1)
	}

	// Check command line arguments:
	if len(argv) != 0 {
		reporter.Errorf("Expected exactly zero command line parameters")
		os.Exit(1)
	}

	// Create the AWS client:
	awsClient, err := aws.NewClient().
		Logger(logger).
		Build()
	if err != nil {
		reporter.Errorf("Can't create AWS client: %v", err)
		os.Exit(1)
	}

	awsCreator, err := awsClient.GetCreator()
	if err != nil {
		reporter.Errorf("Can't get AWS creator: %v", err)
		os.Exit(1)
	}

	// Create the client for the OCM API:
	ocmConnection, err := ocm.NewConnection().
		Logger(logger).
		Build()
	if err != nil {
		reporter.Errorf("Can't create OCM connection: %v", err)
		os.Exit(1)
	}
	defer func() {
		err = ocmConnection.Close()
		if err != nil {
			reporter.Errorf("Can't close OCM connection: %v", err)
		}
	}()

	// Create the writer that will be used to print the tabulated results:
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(writer, "ID\tNAME\tSTATE\n")

	// Retrieve the list of clusters:
	ocmQuery := fmt.Sprintf("properties.%s = '%s'", properties.CreatorARN, awsCreator.ARN)
	ocmRequest := ocmConnection.ClustersMgmt().V1().Clusters().List().
		Search(ocmQuery)
	page := 1
	size := 100
	for {
		ocmResponse, err := ocmRequest.Page(page).Size(size).Send()
		if err != nil {
			reporter.Errorf("Can't retrieve clusters: %v", err)
			os.Exit(1)
		}
		ocmResponse.Items().Each(func(ocmCluster *cmv1.Cluster) bool {
			fmt.Fprintf(
				writer,
				"%s\t%s\t%s\n",
				ocmCluster.ID(),
				ocmCluster.Name(),
				ocmCluster.State(),
			)
			return true
		})
		writer.Flush()
		if ocmResponse.Size() != size {
			break
		}
		page++
	}
}
