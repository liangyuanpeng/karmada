package images

import (
	"fmt"

	"github.com/karmada-io/karmada/pkg/version"
	"k8s.io/klog/v2"
)

var (
	imageRepositories = map[string]string{
		"global": "registry.k8s.io",
		"cn":     "registry.cn-hangzhou.aliyuncs.com/google_containers",
	}

	defaultEtcdImage = "etcd:3.5.3-0"

	// DefaultInitImage etcd init container image
	DefaultInitImage string
	// DefaultKarmadaSchedulerImage Karmada scheduler image
	DefaultKarmadaSchedulerImage string
	// DefaultKarmadaControllerManagerImage Karmada controller manager image
	DefaultKarmadaControllerManagerImage string
	// DefualtKarmadaWebhookImage Karmada webhook image
	DefualtKarmadaWebhookImage string
	// DefaultKarmadaAggregatedAPIServerImage Karmada aggregated apiserver image
	DefaultKarmadaAggregatedAPIServerImage string

	karmadaRelease string
)

func init() {
	releaseVer, err := version.ParseGitVersion(version.Get().GitVersion)
	if err != nil {
		klog.Infof("No default release version found. build version: %s", version.Get().String())
		releaseVer = &version.ReleaseVersion{} // initialize to avoid panic
	}
	karmadaRelease = releaseVer.PatchRelease()

	DefaultInitImage = "docker.io/alpine:3.15.1"
	DefaultKarmadaSchedulerImage = fmt.Sprintf("docker.io/karmada/karmada-scheduler:%s", releaseVer.PatchRelease())
	DefaultKarmadaControllerManagerImage = fmt.Sprintf("docker.io/karmada/karmada-controller-manager:%s", releaseVer.PatchRelease())
	DefualtKarmadaWebhookImage = fmt.Sprintf("docker.io/karmada/karmada-webhook:%s", releaseVer.PatchRelease())
	DefaultKarmadaAggregatedAPIServerImage = fmt.Sprintf("docker.io/karmada/karmada-aggregated-apiserver:%s", releaseVer.PatchRelease())
}

func GetImageRepositories() map[string]string {
	return imageRepositories
}

// get kube components registry
func kubeRegistry() string {
	registry := "registry.k8s.io"
	// registry := i.KubeImageRegistry
	// mirrorCountry := strings.ToLower(i.KubeImageMirrorCountry)
	mirrorCountry := ""

	if registry != "" {
		return registry
	}

	if mirrorCountry != "" {
		value, ok := imageRepositories[mirrorCountry]
		if ok {
			return value
		}
	}

	// if i.ImageRegistry != "" {
	// 	return i.ImageRegistry
	// }
	return imageRepositories["global"]
}

// get kube-apiserver image
func GetkubeAPIServerImage(KarmadaAPIServerImage, kubeImageTag string) string {
	if KarmadaAPIServerImage != "" {
		return KarmadaAPIServerImage
	}

	return kubeRegistry() + "/kube-apiserver:" + kubeImageTag
}

// get kube-controller-manager image
func GetkubeControllerManagerImage(kubeControllerManagerImage, kubeImageTag string) string {
	if kubeControllerManagerImage != "" {
		return kubeControllerManagerImage
	}

	return kubeRegistry() + "/kube-controller-manager:" + kubeImageTag
}

// get etcd-init image
func GetetcdInitImage(imageRegistry, etcdInitImage string) string {
	if imageRegistry != "" && etcdInitImage == DefaultInitImage {
		return imageRegistry + "/alpine:3.15.1"
	}
	return etcdInitImage
}

// get etcd image
func GetetcdImage(etcdImage string) string {
	if etcdImage != "" {
		return etcdImage
	}
	return kubeRegistry() + "/" + defaultEtcdImage
}

// get karmada-scheduler image
func GetkarmadaSchedulerImage(imageRegistry, karmadaSchedulerImage string) string {
	if imageRegistry != "" && karmadaSchedulerImage == DefaultKarmadaSchedulerImage {
		return imageRegistry + "/karmada-scheduler:" + karmadaRelease
	}
	return karmadaSchedulerImage
}

// get karmada-controller-manager
func GetkarmadaControllerManagerImage(imageRegistry, KarmadaControllerManagerImage string) string {
	if imageRegistry != "" && KarmadaControllerManagerImage == DefaultKarmadaControllerManagerImage {
		return imageRegistry + "/karmada-controller-manager:" + karmadaRelease
	}
	return KarmadaControllerManagerImage
}

// get karmada-webhook image
func GetkarmadaWebhookImage(imageRegistry, KarmadaWebhookImage string) string {
	if imageRegistry != "" && KarmadaWebhookImage == DefualtKarmadaWebhookImage {
		return imageRegistry + "/karmada-webhook:" + karmadaRelease
	}
	return KarmadaWebhookImage
}

// get karmada-aggregated-apiserver image
func GetkarmadaAggregatedAPIServerImage(imageRegistry, KarmadaAggregatedAPIServerImage string) string {
	if imageRegistry != "" && KarmadaAggregatedAPIServerImage == DefaultKarmadaAggregatedAPIServerImage {
		return imageRegistry + "/karmada-aggregated-apiserver:" + karmadaRelease
	}
	return KarmadaAggregatedAPIServerImage
}
