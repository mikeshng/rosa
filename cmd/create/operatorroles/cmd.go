/*
Copyright (c) 2021 Red Hat, Inc.

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

package operatorroles

import (
	"os"
	"strings"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/spf13/cobra"

	"github.com/openshift/rosa/pkg/aws"
	"github.com/openshift/rosa/pkg/interactive"
	"github.com/openshift/rosa/pkg/interactive/confirm"
	"github.com/openshift/rosa/pkg/ocm"
	"github.com/openshift/rosa/pkg/rosa"
)

const (
	PrefixFlag           = "prefix"
	HostedCpFlag         = "hosted-cp"
	OidcEndpointUrlFlag  = "oidc-endpoint-url"
	InstallerRoleArnFlag = "installer-role-arn"
)

var args struct {
	prefix              string
	hostedCp            bool
	installerRoleArn    string
	permissionsBoundary string
	forcePolicyCreation bool
	oidcEndpointUrl     string
}

var Cmd = &cobra.Command{
	Use:     "operator-roles",
	Aliases: []string{"operatorroles"},
	Short:   "Create operator IAM roles for a cluster.",
	Long:    "Create cluster-specific operator IAM roles based on your cluster configuration.",
	Example: `  # Create default operator roles for cluster named "mycluster"
  rosa create operator-roles --cluster=mycluster

  # Create operator roles with a specific permissions boundary
  rosa create operator-roles -c mycluster --permissions-boundary arn:aws:iam::123456789012:policy/perm-boundary`,
	RunE: run,
}

func init() {
	flags := Cmd.Flags()

	ocm.AddOptionalClusterFlag(Cmd)

	flags.StringVar(
		&args.prefix,
		PrefixFlag,
		"",
		"User-defined prefix for generated AWS operator policies. Not to be used alongside --cluster flag.",
	)

	flags.StringVar(
		&args.oidcEndpointUrl,
		OidcEndpointUrlFlag,
		"",
		"Oidc endpoint URL to add as the trusted relationship to the operator roles.",
	)

	flags.StringVar(
		&args.installerRoleArn,
		InstallerRoleArnFlag,
		"",
		"Installer role ARN supplied to retrieve operator policy prefix and path.",
	)

	flags.BoolVar(
		&args.hostedCp,
		HostedCpFlag,
		false,
		"Indicates whether to create the hosted control planes operator roles when using --prefix option.",
	)

	flags.StringVar(
		&args.permissionsBoundary,
		"permissions-boundary",
		"",
		"The ARN of the policy that is used to set the permissions boundary for the operator roles.",
	)

	flags.BoolVarP(
		&args.forcePolicyCreation,
		"force-policy-creation",
		"f",
		false,
		"Forces creation of policies skipping compatibility check",
	)

	aws.AddModeFlag(Cmd)
	confirm.AddFlag(flags)
	interactive.AddFlag(flags)
}

func isByoOidcSet(cluster *cmv1.Cluster) bool {
	return cluster != nil && !strings.Contains(cluster.AWS().STS().OIDCEndpointURL(), cluster.ID())
}

func run(cmd *cobra.Command, argv []string) error {
	r := rosa.NewRuntime().WithAWS().WithOCM()
	defer r.Cleanup()

	// Allow the command to be called programmatically
	isProgmaticallyCalled := false
	if len(argv) == 3 && !cmd.Flag("cluster").Changed {
		ocm.SetClusterKey(argv[0])
		aws.SetModeKey(argv[1])
		args.permissionsBoundary = argv[2]

		// if mode is empty skip interactive is true
		if argv[1] != "" {
			isProgmaticallyCalled = true
		}
	}

	env, err := ocm.GetEnv()
	if err != nil {
		r.Reporter.Errorf("Failed to determine OCM environment: %v", err)
		os.Exit(1)
	}

	mode, err := aws.GetMode()
	if err != nil {
		r.Reporter.Errorf("%s", err)
		os.Exit(1)
	}

	// Determine if interactive mode is needed
	if !interactive.Enabled() && !cmd.Flags().Changed("mode") && !isProgmaticallyCalled {
		interactive.Enable()
	}

	if !cmd.Flag("cluster").Changed && !cmd.Flag(PrefixFlag).Changed && !isProgmaticallyCalled {
		r.Reporter.Errorf("Either a cluster key for STS cluster or an operator roles prefix must be specified.")
		os.Exit(1)
	}

	var cluster *cmv1.Cluster
	if args.prefix == "" {
		cluster = r.FetchCluster()
	}

	if args.forcePolicyCreation && mode != aws.ModeAuto {
		r.Reporter.Warnf("Forcing creation of policies only works in auto mode")
		os.Exit(1)
	}

	if interactive.Enabled() && !isProgmaticallyCalled {
		mode, err = interactive.GetOption(interactive.Input{
			Question: "Role creation mode",
			Help:     cmd.Flags().Lookup("mode").Usage,
			Default:  aws.ModeAuto,
			Options:  aws.Modes,
			Required: true,
		})
		if err != nil {
			r.Reporter.Errorf("Expected a valid role creation mode: %s", err)
			os.Exit(1)
		}
	}

	if cluster == nil && interactive.Enabled() && !isProgmaticallyCalled {
		handleOperatorRolesPrefixOptions(r, cmd)
	}

	permissionsBoundary := args.permissionsBoundary
	if interactive.Enabled() && !isProgmaticallyCalled {
		permissionsBoundary, err = interactive.GetString(interactive.Input{
			Question: "Permissions boundary ARN",
			Help:     cmd.Flags().Lookup("permissions-boundary").Usage,
			Default:  permissionsBoundary,
			Validators: []interactive.Validator{
				aws.ARNValidator,
			},
		})
		if err != nil {
			r.Reporter.Errorf("Expected a valid policy ARN for permissions boundary: %s", err)
			os.Exit(1)
		}
	}

	if permissionsBoundary != "" {
		err = aws.ARNValidator(permissionsBoundary)
		if err != nil {
			r.Reporter.Errorf("Expected a valid policy ARN for permissions boundary: %s", err)
			os.Exit(1)
		}
	}

	policies, err := r.OCMClient.GetPolicies("OperatorRole")
	if err != nil {
		r.Reporter.Errorf("Expected a valid role creation mode: %s", err)
		os.Exit(1)
	}

	defaultPolicyVersion, err := r.OCMClient.GetDefaultVersion()
	if err != nil {
		r.Reporter.Errorf("Error getting latest default version: %s", err)
		os.Exit(1)
	}

	if args.prefix != "" {
		if args.oidcEndpointUrl == "" {
			r.Reporter.Errorf("%s is mandatory for %s param flow.", OidcEndpointUrlFlag, PrefixFlag)
			os.Exit(1)
		}

		if args.installerRoleArn == "" {
			r.Reporter.Errorf("%s is mandatory for %s param flow.", InstallerRoleArnFlag, PrefixFlag)
			os.Exit(1)
		}
		return handleOperatorRoleCreationByPrefix(r, env, permissionsBoundary,
			mode, policies, defaultPolicyVersion)
	}
	return handleOperatorRoleCreationByClusterKey(r, env, permissionsBoundary,
		mode, policies, defaultPolicyVersion)
}
