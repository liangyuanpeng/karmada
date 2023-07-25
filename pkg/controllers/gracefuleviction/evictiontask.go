package gracefuleviction

import (
	"math"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	workv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
)

type assessmentOption struct {
	timeout        time.Duration
	scheduleResult []workv1alpha2.TargetCluster
	observedStatus []workv1alpha2.AggregatedStatusItem
	replicas       int32
}

// assessEvictionTasks assesses each task according to graceful eviction rules and
// returns the tasks that should be kept.
func assessEvictionTasks(bindingSpec workv1alpha2.ResourceBindingSpec,
	observedStatus []workv1alpha2.AggregatedStatusItem,
	timeout time.Duration,
	now metav1.Time,
) ([]workv1alpha2.GracefulEvictionTask, []string) {
	var keptTasks []workv1alpha2.GracefulEvictionTask
	var evictedClusters []string

	for _, task := range bindingSpec.GracefulEvictionTasks {
		// set creation timestamp for new task
		if task.CreationTimestamp.IsZero() {
			task.CreationTimestamp = now
			keptTasks = append(keptTasks, task)
			klog.Info("lan.dev.assessEvictionTasks.new task:", task.CreationTimestamp)
			continue
		}
		klog.Info("lan.dev.assessEvictionTasks task:", task.FromCluster, task.CreationTimestamp)

		// assess task according to observed status
		kt := assessSingleTask(task, assessmentOption{
			scheduleResult: bindingSpec.Clusters,
			timeout:        timeout,
			observedStatus: observedStatus,
			replicas:       bindingSpec.Replicas,
		})
		if kt != nil {
			klog.Info("lan.dev.assessEvictionTasks task1:", task.CreationTimestamp)
			keptTasks = append(keptTasks, *kt)
		} else {
			klog.Info("lan.dev.assessEvictionTasks task2:", task.CreationTimestamp)
			evictedClusters = append(evictedClusters, task.FromCluster)
		}
	}
	return keptTasks, evictedClusters
}

func assessSingleTask(task workv1alpha2.GracefulEvictionTask, opt assessmentOption) *workv1alpha2.GracefulEvictionTask {
	if task.SuppressDeletion != nil {
		if *task.SuppressDeletion {
			return &task
		}
		// If *task.SuppressDeletion is equal to false,
		// it means users have confirmed that they want to delete the redundant copy.
		// In that case, we will delete the task immediately.
		return nil
	}

	timeout := opt.timeout
	if task.GracePeriodSeconds != nil {
		timeout = time.Duration(*task.GracePeriodSeconds) * time.Second
	}
	klog.Info("lan.dev.assessSingleTask:", timeout)
	// task exceeds timeout
	if metav1.Now().After(task.CreationTimestamp.Add(timeout)) {
		klog.Info("lan.dev.assessSingleTask2:", timeout)
		return nil
	}

	if allScheduledResourceInHealthyState(opt) {
		klog.Info("lan.dev.assessSingleTask3:")
		return nil
	}

	klog.Info("lan.dev.assessSingleTask4:")
	return &task
}

func allScheduledResourceInHealthyState(opt assessmentOption) bool {
	var scheduleReplaces int32 = 0
	for _, targetCluster := range opt.scheduleResult {
		scheduleReplaces += targetCluster.DeepCopy().Replicas
		var statusItem *workv1alpha2.AggregatedStatusItem

		// find the observed status of targetCluster
		for index, aggregatedStatus := range opt.observedStatus {
			if aggregatedStatus.ClusterName == targetCluster.Name {
				statusItem = &opt.observedStatus[index]
				break
			}
		}

		// no observed status found, maybe the resource hasn't been applied
		if statusItem == nil {
			return false
		}
		klog.Info("lan.dev.allScheduledResourceInHealthyState.aggregatedStatus:", statusItem.ClusterName, statusItem.Health, statusItem.Health)
		klog.Info("lan.dev.allScheduledResourceInHealthyState.aggregatedStatus2:", statusItem.Status)

		// resource not in healthy state
		if statusItem.Health != workv1alpha2.ResourceHealthy {
			return false
		}

	}
	if scheduleReplaces != opt.replicas {
		klog.Infof("lan.dev.allScheduledResourceInHealthyState.replicas is not equal:%d,%d", opt.replicas, scheduleReplaces)
		// return false
	}
	klog.Infof("lan.dev.allScheduledResourceInHealthyState.return true:%d,%d", opt.replicas, scheduleReplaces)

	return true
}

func nextRetry(tasks []workv1alpha2.GracefulEvictionTask, timeout time.Duration, timeNow time.Time) time.Duration {
	if len(tasks) == 0 {
		return 0
	}

	retryInterval := time.Duration(math.MaxInt64)

	// Skip tasks whose type is SuppressDeletion because they are manually controlled by users.
	// We currently take the minimum value of the timeout of all GracefulEvictionTasks besides above.
	// When the application on the new cluster becomes healthy, a new event will be queued
	// because the controller can watch the changes of binding status.
	for i := range tasks {
		if tasks[i].SuppressDeletion != nil {
			continue
		}
		if tasks[i].GracePeriodSeconds != nil {
			timeout = time.Duration(*tasks[i].GracePeriodSeconds) * time.Second
		}
		next := tasks[i].CreationTimestamp.Add(timeout).Sub(timeNow)
		if next < retryInterval {
			retryInterval = next
		}
	}

	// When there are only tasks whose type is SuppressDeletion, we do not need to retry.
	if retryInterval == time.Duration(math.MaxInt64) {
		retryInterval = 0
	}
	return retryInterval
}
