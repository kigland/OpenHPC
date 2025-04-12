package scheduler

import (
	"time"
)

type Scheduler struct {
	Interval time.Duration
}

func (s *Scheduler) StartGCDaemon() {
	for {
		time.Sleep(s.Interval)
		// s.GC()
	}
}

// func (s *Scheduler) GC() {
// 	containers, err := shared.DockerHelper.AllKHSContainers()
// 	if err != nil {
// 		log.Println("Error getting all containers:", err)
// 		return
// 	}
// 	dh := shared.DockerHelper
// 	for n, c := range containers {
// 		if c.State != "running" {
// 			continue
// 		}
// 		timeStr, err := shared.GCKVStore.Get(n)
// 		if err != nil {
// 			log.Println("Error getting TO:", err)
// 			continue
// 		}
// 		timeStr = strings.TrimSpace(timeStr)
// 		if timeStr == "" {
// 			continue
// 		}
// 		ddl, err := time.Parse(time.RFC3339, timeStr)
// 		if err != nil {
// 			log.Println("Error parsing TO:", err)
// 			continue
// 		}
// 		if ddl.Before(time.Now()) {
// 			continue
// 		}
// 		err = dh.StopContainer(n)
// 		if err != nil {
// 			log.Println("Error stopping container:", err)
// 			continue
// 		}

// 		err = dh.RemoveContainer(n)
// 		if err != nil {
// 			log.Println("Error removing container:", err)
// 			continue
// 		}
// 		err = shared.GCKVStore.Delete(n)
// 		if err != nil {
// 			log.Println("Error removing container:", err)
// 			continue
// 		}
// 	}
// }
