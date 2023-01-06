package statshelper

import (
	"sync"

	domain "github.com/matteogazzadi/kuberes/pkg/domain"
	v1 "k8s.io/api/core/v1"
)

func GetPodStats(pod v1.Pod, podChan chan<- domain.K8sStats, wg *sync.WaitGroup) {
	defer wg.Done()

	var stats domain.K8sStats

	stats.PodName = pod.Name
	stats.Namespace = pod.Namespace

	// Loop On Containers
	for _, container := range pod.Spec.Containers {

		// CPU
		cpuRequest := container.Resources.Requests.Cpu().MilliValue()
		cpuLimit := container.Resources.Limits.Cpu().MilliValue()

		// Memory
		memRequest, _ := container.Resources.Requests.Memory().AsInt64()
		memLimit, _ := container.Resources.Limits.Memory().AsInt64()

		// Convert MB to Mib
		memRequest = memRequest / 1048576
		memLimit = memLimit / 1048576

		stats.Cpu.Limit += cpuLimit
		stats.Cpu.Request += cpuRequest

		stats.Memory.Limit += memLimit
		stats.Memory.Request += memRequest
	}
	podChan <- stats
}
