package config

import (
	"fmt"
	"log"

	"github.com/karmada-io/karmada/pkg/karmadactl/cmdinit/kubernetes"
	"github.com/karmada-io/karmada/pkg/karmadactl/util"
	"github.com/karmada-io/karmada/pkg/version"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"k8s.io/kubectl/pkg/util/templates"
)

var (
	initLong = templates.LongDesc(`
		hello initLong.`)

	initExamples = templates.Examples(`
		# hello initExamples
		`)
)

// NewCmdInit install Karmada on Kubernetes
func NewCmdConfig(parentCommand string) *cobra.Command {
	opts := kubernetes.CommandInitOption{}
	cmd := &cobra.Command{
		Use:                   "config",
		Short:                 "karmadactl config",
		Long:                  initLong,
		Example:               initExample(parentCommand),
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := opts.Validate(parentCommand); err != nil {
				return err
			}
			log.Println("aa:", kubernetes.DefaultKarmadaAggregatedAPIServerImage)
			if err := opts.Complete(); err != nil {
				return err
			}
			if err := opts.RunInit(parentCommand); err != nil {
				return err
			}
			return nil
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
		Annotations: map[string]string{
			util.TagCommandGroup: util.GroupClusterRegistration,
		},
	}
	// flags := cmd.Flags()
	// flags.StringVarP(&opts.ImageRegistry, "private-image-registry", "", "", "Private image registry where pull images from. If set, all required images will be downloaded from it, it would be useful in offline installation scenarios.  In addition, you still can use --kube-image-registry to specify the registry for Kubernetes's images.")
	// flags.StringSliceVar(&opts.PullSecrets, "image-pull-secrets", nil, "Image pull secrets are used to pull images from the private registry, could be secret list separated by comma (e.g '--image-pull-secrets PullSecret1,PullSecret2', the secrets should be pre-settled in the namespace declared by '--namespace')")

	return cmd
}

func initExample(parentCommand string) string {
	releaseVer, err := version.ParseGitVersion(version.Get().GitVersion)
	if err != nil {
		klog.Infof("No default release version found. build version: %s", version.Get().String())
		releaseVer = &version.ReleaseVersion{}
	}
	return fmt.Sprintf(initExamples, parentCommand, releaseVer.FirstMinorRelease())
}
