package testfilter

import (
	"context"

	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	workv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
	"github.com/karmada-io/karmada/pkg/scheduler/framework"
)

const (
	// Name is the name of the plugin used in the plugin registry and configurations.
	Name = "Schedulerb"
)

type TestFilter struct{}

var _ framework.FilterPlugin = &TestFilter{}
var _ framework.ScorePlugin = &TestFilter{}

// New instantiates the TestFilter plugin.
func New() (framework.Plugin, error) {
	return &TestFilter{}, nil
}

// Name returns the plugin name.
func (p *TestFilter) Name() string {
	return Name
}

// Filter implements the filtering logic of the TestFilter plugin.
func (p *TestFilter) Filter(ctx context.Context,
	bindingSpec *workv1alpha2.ResourceBindingSpec, bindingStatus *workv1alpha2.ResourceBindingStatus, cluster *clusterv1alpha1.Cluster) *framework.Result {

	// implementation

	return framework.NewResult(framework.Success)
}

func (p *TestFilter) Score(ctx context.Context, spec *workv1alpha2.ResourceBindingSpec, cluster *clusterv1alpha1.Cluster) (int64, *framework.Result) {
	return 100, framework.NewResult(framework.Success)
}
func (p *TestFilter) ScoreExtensions() framework.ScoreExtensions {
	// 返回 实现了ScoreExtensions接口的实现 或 nil
	return nil
}
