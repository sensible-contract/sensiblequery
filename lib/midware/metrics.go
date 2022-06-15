package midware

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	kAll        int = iota // 总值下标
	kThisSecond     = iota // 当前秒下标
	kLastSecond     = iota // 前一秒下标
	kThisMinute     = iota // 当前分钟下标
	kLastMinute     = iota // 前一分钟下标
	kNum            = iota

	mapNumMax int64 = 100000
)

var programStartTime = time.Now()

var enabledMetricsDuration []int = []int{
	// kAll,
	kThisSecond, kThisMinute,
}

var enabledMetricsDurationWithAll []int = []int{
	kAll,
	kThisSecond, kThisMinute,
}

var enabledMetricsDurationDump []int = []int{
	kAll,
	kThisSecond, kThisMinute,
	kLastSecond, kLastMinute,
}

var modifiedDuration [kNum]time.Duration = [kNum]time.Duration{
	kAll:        time.Hour * 876000, // 100 years, int64 is 292 years
	kThisSecond: time.Second,
	kLastSecond: time.Second * 2,
	kThisMinute: time.Minute,
	kLastMinute: time.Minute * 2,
}

var metricsKeyName [kNum]string = [kNum]string{
	kAll:        "total",
	kThisSecond: "second",
	kLastSecond: "second",
	kThisMinute: "minute",
	kLastMinute: "minute",
}

var metricsStageName [kNum]string = [kNum]string{
	kAll:        "past",
	kThisSecond: "now",
	kLastSecond: "past",
	kThisMinute: "now",
	kLastMinute: "past",
}

type metricsItem struct {
	v [kNum]metricsValue
	m sync.RWMutex
}

type metricsValue struct {
	expiredTime time.Time
	value       int64 // value
	latency     int64 // latency
}

type serviceKey struct {
	Name string
	Type string
}

type serviceMetrics struct {
	metricsItemIdx   map[serviceKey]int64
	metricsKeyCache  [mapNumMax]serviceKey
	metricsItemCache [mapNumMax]metricsItem
	m                sync.RWMutex
	n                *int64
}

var ServiceMetrics *serviceMetrics = &serviceMetrics{
	metricsItemIdx: make(map[serviceKey]int64),
	n:              new(int64),
}

func metricFormat(strBuilder *strings.Builder, key, sName, sType, stage string, value int64) {
	fmt.Fprintf(strBuilder, `%s{name="%s",type="%s",stage="%s"} %d%s`, key, sName, sType, stage, value, "\n")
}

func (sm *serviceMetrics) Dump(keyFilter, typeFilter, stageFilter string) (metrics string) {
	metricsCount := atomic.LoadInt64(sm.n)

	var metricsBuilder strings.Builder

	metricFormat(&metricsBuilder, "metrics_num", "debug", "runtime", "now", metricsCount)

	// goroutine number
	metricFormat(&metricsBuilder, "goroutine_num", "debug", "runtime", "now", int64(runtime.NumGoroutine()))
	metricFormat(&metricsBuilder, "uptime", "debug", "runtime", "now", int64(time.Since(programStartTime)/time.Second))

	// memory
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	metricFormat(&metricsBuilder, "memory_alloc", "debug", "runtime", "now", int64(ms.Alloc))
	metricFormat(&metricsBuilder, "memory_heap_idle", "debug", "runtime", "now", int64(ms.HeapIdle))
	metricFormat(&metricsBuilder, "memory_heap_released", "debug", "runtime", "now", int64(ms.HeapReleased))

	for idx := range make([]struct{}, metricsCount) {
		itemIdx := int64(idx)
		sName := sm.metricsKeyCache[itemIdx].Name
		sType := sm.metricsKeyCache[itemIdx].Type
		if typeFilter != "-" && typeFilter != sType { // match type
			continue
		}

		sm.flushData(itemIdx)
		data := &(sm.metricsItemCache[itemIdx])
		data.m.RLock()

		for _, durIdx := range enabledMetricsDurationDump {
			stage := metricsStageName[durIdx]
			if stageFilter != "-" && stageFilter != stage { // match stage
				continue
			}

			if data.v[durIdx].value == 0 {
				continue
			}

			keyName := metricsKeyName[durIdx]

			// 输出counter
			if data.v[durIdx].latency == 0 {
				if time.Now().After(data.v[durIdx].expiredTime) {
					data.v[durIdx].value = 0 // counter类型数据清理
					break
				}
				metricFormat(&metricsBuilder, "counter_"+keyName, sName, sType, stage, data.v[durIdx].value)
				break
			}

			if keyFilter != "-" && keyFilter != keyName { // match key
				continue
			}
			metricFormat(&metricsBuilder, "queries_"+keyName, sName, sType, stage, data.v[durIdx].value)
			metricFormat(&metricsBuilder, "latency_"+keyName, sName, sType, stage, data.v[durIdx].latency)
			metricFormat(&metricsBuilder, "latency_avg_"+keyName, sName, sType, stage, data.v[durIdx].latency/data.v[durIdx].value)
		}

		data.m.RUnlock()
	}

	metrics = metricsBuilder.String()
	return metrics
}

// flushData 轮转数据
func (sm *serviceMetrics) flushData(itemIdx int64) {
	data := &(sm.metricsItemCache[itemIdx])
	data.m.Lock()
	defer data.m.Unlock()

	// Total
	if data.v[kAll].latency > int64(modifiedDuration[kAll]) {
		data.v[kAll].value = 0
		data.v[kAll].latency = 0
	}
	now := time.Now()
	for _, durIdx := range enabledMetricsDuration {
		if data.v[durIdx+1].value == 0 && data.v[durIdx].value == 0 {
			continue
		}

		if now.After(data.v[durIdx].expiredTime) {
			data.v[durIdx+1].value = data.v[durIdx].value
			data.v[durIdx+1].latency = data.v[durIdx].latency
			data.v[durIdx].expiredTime = now.Add(modifiedDuration[durIdx])
			data.v[durIdx].value = 0
			data.v[durIdx].latency = 0
		}
	}
}

func (sm *serviceMetrics) addOrSet(serviceName, serviceType string, isAdd bool, value int64, expire time.Duration) {
	serviceKeyPair := serviceKey{
		Name: serviceName,
		Type: serviceType,
	}

	sm.m.RLock()
	itemIdx, ok := sm.metricsItemIdx[serviceKeyPair]
	sm.m.RUnlock()
	if !ok {
		// 新增监控项
		nextItemIdx := atomic.LoadInt64(sm.n)
		for {
			if atomic.CompareAndSwapInt64(sm.n, nextItemIdx, (nextItemIdx + 1)) {
				break
			}
			nextItemIdx = atomic.LoadInt64(sm.n)
		}
		// 限最多监控条目
		if nextItemIdx >= mapNumMax {
			return
		}
		itemIdx = nextItemIdx
		sm.metricsKeyCache[itemIdx].Name = serviceName
		sm.metricsKeyCache[itemIdx].Type = serviceType

		sm.m.Lock()
		sm.metricsItemIdx[serviceKeyPair] = itemIdx
		sm.m.Unlock()
	}

	data := &(sm.metricsItemCache[itemIdx])
	data.m.Lock()
	defer data.m.Unlock()

	data.v[kAll].expiredTime = time.Now().Add(expire)
	if isAdd {
		data.v[kAll].value += value
	} else {
		data.v[kAll].value = value
	}
}

// Set 累加器设置，带过期时长
func (sm *serviceMetrics) Set(serviceName, serviceType string, value int64, expire time.Duration) {
	sm.addOrSet(serviceName, serviceType, false, value, expire)
}

// Inc 累加器增加，默认不过期
func (sm *serviceMetrics) Inc(serviceName, serviceType string) {
	sm.addOrSet(serviceName, serviceType, true, 1, modifiedDuration[kAll])
}

// Dec 累加器减少，默认不过期
func (sm *serviceMetrics) Dec(serviceName, serviceType string) {
	sm.addOrSet(serviceName, serviceType, true, -1, modifiedDuration[kAll])
}

// Observe 记录一次访问请求，并统计延迟
func (sm *serviceMetrics) Observe(serviceName, serviceType string, latencyDuration time.Duration) {
	sm.Update(serviceName, serviceType, latencyDuration)
}

func (sm *serviceMetrics) Update(serviceName, serviceType string, latencyDuration time.Duration) {
	latency := int64(latencyDuration)
	if latency <= 0 {
		latency = 1
	}
	serviceKeyPair := serviceKey{
		Name: serviceName,
		Type: serviceType,
	}

	sm.m.RLock()
	itemIdx, ok := sm.metricsItemIdx[serviceKeyPair]
	sm.m.RUnlock()
	if !ok {
		// 新增监控项
		nextItemIdx := atomic.LoadInt64(sm.n)
		for {
			if atomic.CompareAndSwapInt64(sm.n, nextItemIdx, (nextItemIdx + 1)) {
				break
			}
			nextItemIdx = atomic.LoadInt64(sm.n)
		}
		// 限最多监控条目
		if nextItemIdx >= mapNumMax {
			return
		}
		itemIdx = nextItemIdx
		//初始化监控时间
		for _, durIdx := range enabledMetricsDuration {
			sm.metricsItemCache[itemIdx].v[durIdx].expiredTime = time.Now().Add(modifiedDuration[durIdx])
		}
		sm.metricsKeyCache[itemIdx].Name = serviceName
		sm.metricsKeyCache[itemIdx].Type = serviceType

		sm.m.Lock()
		sm.metricsItemIdx[serviceKeyPair] = itemIdx
		sm.m.Unlock()
	} else {
		//先尝试轮转数据
		sm.flushData(itemIdx)
	}

	data := &(sm.metricsItemCache[itemIdx])
	data.m.Lock()
	defer data.m.Unlock()

	for _, durIdx := range enabledMetricsDurationWithAll {
		// QPS
		data.v[durIdx].value += 1
		// Response Time
		data.v[durIdx].latency += latency
	}
}
