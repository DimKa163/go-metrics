package runtime

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func ReadMemoryStats(values map[string]float64) error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	values["Alloc"] = float64(memStats.Alloc)
	values["BuckHashSys"] = float64(memStats.BuckHashSys)
	values["Frees"] = float64(memStats.Frees)
	values["GCCPUFraction"] = float64(memStats.GCCPUFraction)
	values["GCSys"] = float64(memStats.GCSys)
	values["HeapAlloc"] = float64(memStats.HeapAlloc)
	values["HeapIdle"] = float64(memStats.HeapIdle)
	values["HeapInuse"] = float64(memStats.HeapInuse)
	values["HeapObjects"] = float64(memStats.HeapObjects)
	values["HeapReleased"] = float64(memStats.HeapReleased)
	values["HeapSys"] = float64(memStats.HeapSys)
	values["LastGC"] = float64(memStats.LastGC)
	values["MCacheInuse"] = float64(memStats.MCacheInuse)
	values["MCacheSys"] = float64(memStats.MCacheSys)
	values["MSpanSys"] = float64(memStats.MSpanSys)
	values["Mallocs"] = float64(memStats.Mallocs)
	values["NextGC"] = float64(memStats.NextGC)

	values["NumForcedGC"] = float64(memStats.NumForcedGC)
	values["NumGC"] = float64(memStats.NumGC)
	values["OtherSys"] = float64(memStats.OtherSys)

	values["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	values["StackInuse"] = float64(memStats.StackInuse)
	values["StackSys"] = float64(memStats.StackSys)
	values["Sys"] = float64(memStats.Sys)

	values["TotalAlloc"] = float64(memStats.TotalAlloc)
	values["MSpanInuse"] = float64(memStats.MSpanInuse)
	values["Lookups"] = float64(memStats.Lookups)
	values["RandomValue"] = rand.Float64()

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	values["TotalMemory"] = float64(vmStat.Total / 1024 / 1024)
	values["FreeMemory"] = float64(vmStat.Free / 1024 / 1024)
	return nil
}

func ReadCPUStats(values map[string]float64) error {
	cp, err := cpu.Percent(time.Second, true)
	if err != nil {
		return err
	}
	for i, c := range cp {
		values[fmt.Sprintf("CPUutilization%d", i+1)] = c
	}
	return nil
}
