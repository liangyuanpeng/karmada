package images

import (
	"fmt"
	"log"
	"strings"

	"github.com/karmada-io/karmada/pkg/version"
	"k8s.io/klog/v2"
)

var (
	imageRepositories = map[string]string{
		"global": "registry.k8s.io",
		"cn":     "registry.cn-hangzhou.aliyuncs.com/google_containers",
	}

	DefaultEtcdImage = "etcd:3.5.3-0"

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

	log.Println("================================version:", version.Get().GitVersion, version.Get().GitCommit, releaseVer.PatchRelease(), karmadaRelease)

	DefaultInitImage = "docker.io/alpine:3.15.1"
	DefaultKarmadaSchedulerImage = fmt.Sprintf("docker.io/karmada/karmada-scheduler:%s", releaseVer.PatchRelease())
	DefaultKarmadaControllerManagerImage = fmt.Sprintf("docker.io/karmada/karmada-controller-manager:%s", releaseVer.PatchRelease())
	DefualtKarmadaWebhookImage = fmt.Sprintf("docker.io/karmada/karmada-webhook:%s", releaseVer.PatchRelease())
	DefaultKarmadaAggregatedAPIServerImage = fmt.Sprintf("docker.io/karmada/karmada-aggregated-apiserver:%s", releaseVer.PatchRelease())
}

func GetKarmadaRelease() string {
	return karmadaRelease
}

func GetImageRepositories() map[string]string {
	return imageRepositories
}

// get kube components registry
func KubeRegistry(kubeImageRegistry, kubeImageMirrorCountry string) string {
	registry := kubeImageRegistry
	mirrorCountry := strings.ToLower(kubeImageMirrorCountry)

	if registry != "" {
		return registry
	}

	if mirrorCountry != "" {
		value, ok := imageRepositories[mirrorCountry]
		if ok {
			return value
		}
	}

	if kubeImageRegistry != "" {
		return kubeImageRegistry
	}
	return imageRepositories["global"]
}

// get kube-apiserver image
func GetkubeAPIServerImage(kubeImageRegistry, kubeImageMirrorCountry, KarmadaAPIServerImage, kubeImageTag string) string {
	if KarmadaAPIServerImage != "" {
		return KarmadaAPIServerImage
	}

	return KubeRegistry(kubeImageRegistry, kubeImageMirrorCountry) + "/kube-apiserver:" + kubeImageTag
}

// get kube-controller-manager image
func GetkubeControllerManagerImage(kubeImageRegistry, kubeImageMirrorCountry, kubeControllerManagerImage, kubeImageTag string) string {
	if kubeControllerManagerImage != "" {
		return kubeControllerManagerImage
	}

	return KubeRegistry(kubeImageRegistry, kubeImageMirrorCountry) + "/kube-controller-manager:" + kubeImageTag
}

// get etcd-init image
func GetetcdInitImage(imageRegistry, etcdInitImage string) string {
	if imageRegistry != "" && etcdInitImage == DefaultInitImage {
		return imageRegistry + "/alpine:3.15.1"
	}
	return etcdInitImage
}

// get etcd image
func GetetcdImage(kubeImageRegistry, kubeImageMirrorCountry, etcdImage string) string {
	if etcdImage != "" {
		return etcdImage
	}
	return KubeRegistry(kubeImageRegistry, kubeImageMirrorCountry) + "/" + DefaultEtcdImage
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
