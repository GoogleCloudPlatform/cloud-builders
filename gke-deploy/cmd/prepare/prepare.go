// Package prepare contains the logic for `gke-deploy prepare` subcommand.
package prepare

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/spf13/cobra"

	"github.com/GoogleCloudPlatform/cloud-builders/gke-deploy/cmd/common"
)

const (
	short = "Execute prepare phase and skip apply phase"
	long  = `Prepare to deploy to GKE by generating expanded Kubernetes configuration files. Skip apply.

- Expand Kubernetes configuration files to:
  - Set the digest of images that match the [--image|-i] flag, if provided.
  - Add app.kubernetes.io/name=[--app|-a] label, if provided.
  - Add app.kubernetes.io/version=[--version|-v] label, if provided.
`
	example = `  # Prepare only.
  gke-deploy prepare -f configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o expanded -n my-namespace

  # Execute prepare and apply, with an intermediary step in between (e.g., manually check expanded YAMLs)
  gke-deploy prepare -f configs -i gcr.io/my-project/my-app:1.0.0 -a my-app -v 1.0.0 -o expanded -n my-namespace
  cat expanded/*
  gke-deploy apply -f expanded -n my-namespace -c my-cluster -l us-east1-b  # Pass expanded directory to -f

  # Pipe output from another templating engine to gke-deploy prepare.
  kustomize build overlays/staging | gke-deploy prepare -f - -a my-app
  helm template charts/prometheus | gke-deploy prepare -f - -a prometheus`
)

type options struct {
	appName             string
	appVersion          string
	filename            string
	image               string
	labels              []string
	annotations         []string
	namespace           string
	output              string
	exposePort          int
	createApplicationCR bool
	applicationLinks    []string
	verbose             bool
	recursive           bool
}

// NewPrepareCommand creates the `gke-deploy prepare` subcommand.
func NewPrepareCommand() *cobra.Command {
	options := &options{}

	cmd := &cobra.Command{
		Use:     "prepare",
		Aliases: []string{"p"},
		Short:   short,
		Long:    long,
		Example: example,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return prepare(cmd, options)
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVarP(&options.appName, "app", "a", "", "Application name of the Kubernetes deployment.")
	cmd.Flags().StringVarP(&options.appVersion, "version", "v", "", "Version of the Kubernetes deployment.")
	cmd.Flags().StringVarP(&options.filename, "filename", "f", "", "Local or GCS path to configuration file or directory of configuration files to use to create Kubernetes objects (file or files in directory must end in \".yml\" or \".yaml\"). Prefix this value with \"gs://\" to indicate a GCS path. If this field is not provided, a Deployment (with image provided by --image) and a HorizontalPodAutoscaler are created as suggested based configs. The application's name is inferred from the image name's suffix.")
	cmd.Flags().StringVarP(&options.image, "image", "i", "", "Image to be deployed.")
	cmd.Flags().StringSliceVarP(&options.labels, "label", "L", nil, "Label(s) to add to Kubernetes objects (k1=v1). Labels can be set comma-delimited or as separate flags. If two or more labels with the same key are listed, the last one is used.")
	cmd.Flags().StringVarP(&options.namespace, "namespace", "n", "", "Namespace of GKE cluster to deploy to. Creates a namespace Kubernetes configuration file to reflect this and updates the namespace field of each supplied Kubernetes configuration file.")
	cmd.Flags().StringSliceVarP(&options.annotations, "annotation", "A", nil, "Annotation(s) to add to Kubernetes configuration files (k1=v1). Annotations can be set comma-delimited or as separate flags. If two or more annotations with the same key are listed, the last one is used.")
	cmd.Flags().StringVarP(&options.output, "output", "o", "./output", "Target directory or GCS path to store suggested and expanded Kubernetes configuration files. Prefix this value with \"gs://\" to indicate a GCS path. Suggested files will be stored in \"<output>/suggested\" and expanded files will be stored in \"<output>/expanded\".")
	cmd.Flags().IntVarP(&options.exposePort, "expose", "x", 0, "Creates a Service object that connects to a deployed workload object using a selector that matches the label with key as 'app.kubernetes.io/name' and value specified by --app. The port provided will be used to expose the deployed workload object (i.e., port and targetPort will be set to the value provided in this flag).")
	cmd.Flags().BoolVarP(&options.verbose, "verbose", "V", false, "Prints underlying commands being called to stdout.")
	cmd.Flags().BoolVarP(&options.recursive, "recursive", "R", false, "Recursively search through the provided path in --filename for all YAML files.")
	cmd.Flags().BoolVar(&options.createApplicationCR, "create-application-cr", false, "Creates an Application CR object with the name provided by --app and connects to deployed objects using a selector that matches the label with key as 'app.kubernetes.io/name' and value specified by --app.")
	cmd.Flags().StringSliceVar(&options.applicationLinks, "links", nil, "Links(s) to add to the spec.descriptor.links field of an Application CR generated with the --create-application-cr flag or provided via the --filename flag (description=URL). Links can be set comma-delimited or as separate flags.")

	return cmd
}

func prepare(_ *cobra.Command, options *options) error {
	ctx := context.Background()

	var im name.Reference
	if options.image != "" {
		ref, err := name.ParseReference(options.image)
		if err != nil {
			return err
		}
		im = ref
	}

	if options.filename == "" && options.image == "" {
		return fmt.Errorf("omitting -f|--filename requires -i|--image to be set")
	}
	if options.output == "" {
		return fmt.Errorf("value of -o|--output cannot be empty")
	}

	if options.exposePort < 0 {
		return fmt.Errorf("value of -x|--expose must be > 0")
	}
	if options.exposePort > 0 && options.appName == "" {
		return fmt.Errorf("exposing a deployed workload object requires -a|--app to be set")
	}

	if options.createApplicationCR && options.appName == "" {
		return fmt.Errorf("creating an Application CR requires -a|--app to be set")
	}

	labelsMap, err := common.CreateMapFromEqualDelimitedStrings(options.labels)
	if err != nil {
		return err
	}
	annotationsMap, err := common.CreateMapFromEqualDelimitedStrings(options.annotations)
	if err != nil {
		return err
	}
	applicationLinks, err := common.CreateApplicationLinksListFromEqualDelimitedStrings(options.applicationLinks)
	if err != nil {
		return err
	}
	d, err := common.CreateDeployer(ctx, false /* useGcloud */, options.verbose, false /* serverDryRun */)
	if err != nil {
		return err
	}

	if err := d.Prepare(ctx, im, options.appName, options.appVersion, options.filename, common.SuggestedOutputPath(options.output), common.ExpandedOutputPath(options.output), options.namespace, labelsMap, annotationsMap, options.exposePort, options.recursive, options.createApplicationCR, applicationLinks); err != nil {
		return fmt.Errorf("failed to prepare deployment: %v", err)
	}

	return nil
}
