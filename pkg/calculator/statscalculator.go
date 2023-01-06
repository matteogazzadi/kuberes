package calculator

import (
	"sync"

	domain "github.com/matteogazzadi/kuberes/pkg/domain"
	helper "github.com/matteogazzadi/kuberes/pkg/stats-helper"
	v1 "k8s.io/api/core/v1"
)

// Calculate Resources
func CalculateResources(groupByNamespace bool, pods *[]v1.Pod, resources *[]domain.K8sStats) {

	// Statistics Channel
	k8sStatsChan := make(chan domain.K8sStats, len(*pods))

	// Wait Group
	var wg sync.WaitGroup

	for _, pod := range *pods {
		wg.Add(1)
		go helper.GetPodStats(pod, k8sStatsChan, &wg)
	}

	wg.Wait()
	close(k8sStatsChan)

	if groupByNamespace {
		resByNs := make(map[string]*domain.K8sStats)

		for stats := range k8sStatsChan {

			curStats, ok := resByNs[stats.Namespace]

			if ok {
				// Update current stats
				curStats.Cpu.Limit += stats.Cpu.Limit
				curStats.Cpu.Request += stats.Cpu.Request
				curStats.Memory.Limit += stats.Memory.Limit
				curStats.Memory.Request += stats.Memory.Request
			} else {
				// Create new stats
				var newStats domain.K8sStats

				newStats.Namespace = stats.Namespace
				newStats.Cpu = stats.Cpu
				newStats.Memory = stats.Memory

				resByNs[stats.Namespace] = &newStats
			}
		}

		for _, stats := range resByNs {
			*resources = append(*resources, *stats)
		}

	} else {
		for stats := range k8sStatsChan {
			*resources = append(*resources, stats)
		}
	}
}
