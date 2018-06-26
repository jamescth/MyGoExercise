package parse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Hosts struct {
	BaseFields

	// The following is the json fmt for da_alarms_show output
	HostSummary []struct {
		DataIP                string  `json:"dataIp"`
		DataReduction         float64 `json:"dataReduction"`
		DNSName               string  `json:"dnsName"`
		FlashMode             string  `json:"flashMode"`
		Headroom              float64 `json:"headroom"`
		HostHeadroom          float64 `json:"hostHeadroom"`
		HostOptimalCacheSize  float64 `json:"hostOptimalCacheSize"`
		HostPhysicalCacheSize float64 `json:"hostPhysicalCacheSize"`
		ID                    string  `json:"id"`
		Iops                  float64 `json:"iops"`
		Name                  string  `json:"name"`
		OptimalCacheSize      float64 `json:"optimalCacheSize"`
		OsName                string  `json:"osName"`
		PhysicalCacheSize     float64 `json:"physicalCacheSize"`
		PrivateCacheUsage     float64 `json:"privateCacheUsage"`
		QueueDepth            float64 `json:"queueDepth"`
		SharedCacheUsage      float64 `json:"sharedCacheUsage"`
		Status                struct {
			Description        string        `json:"description"`
			Events             []interface{} `json:"events"`
			ExtendedAttributes []interface{} `json:"extendedAttributes"`
			ID                 string        `json:"id"`
			Kind               string        `json:"kind"`
			Status             int           `json:"status"`
		} `json:"status"`
		Throughput float64 `json:"throughput"`
		Vcenter    string  `json:"vcenter"`
		WriteIops  float64 `json:"writeIops"`
	} `json:"hostSummary"`
	SsdSummary []interface{} `json:"ssdSummary"`
	VMSummary  []interface{} `json:"vmSummary"`
}

func NewHosts(p string) (DaOut, error) {
	content, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("NewHosts"))
	}
	var h Hosts
	if err := json.Unmarshal(content, &h); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("NewHosts %s", p))
	}

	h.FileName = p

	base := path.Base(p)
	base = strings.TrimPrefix(base, "da_hosts_show.")
	base = strings.TrimRight(base, ".out")
	sTime, err := time.Parse("2006-01-02T15_04_05", base)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("NewHosts %s parse time %s", p, base))
	}

	h.StartTime = sTime.UnixNano()
	h.EndTime = sTime.UnixNano()

	var da DaOut = &h
	return da, nil
}

// func (h *Hosts) ListIssue(start, end int64, c, s string, w io.Writer) error {
func (h *Hosts) ListIssue(req *Request) (*bytes.Buffer, error) {
	var b bytes.Buffer

	//b.WriteString(fmt.Sprintf("file: %s, time: %s\n", h.Name(), daNanosecond(h.End())))
	b.WriteString(fmt.Sprintf("%30s %20s %12s %10s %30s %10s %1s\n", "DNS Name", "Name", "FlashMode", "OSName", "ID", "Kind", "S"))
	b.WriteString(fmt.Sprintf("============================== ==================== ============ ========== ============================== ========== =\n"))
	//fmt.Fprintf(w, "file: %s, time: %s\n", h.Name(), daNanosecond(h.End()))
	//fmt.Fprintf(w, "%30s %20s %12s %10s %30s %10s %1s\n", "DNS Name", "Name", "FlashMode", "OSName", "ID", "Kind", "S")
	//fmt.Fprintf(w, "============================== ==================== ============ ========== ============================== ========== =\n")
	for _, host := range h.HostSummary {
		//fmt.Fprintf(w, "%30s %20s %12s %10s %30s %10s %d\n",
		b.WriteString(fmt.Sprintf("%30s %20s %12s %10s %30s %10s %d\n",
			host.DNSName,
			host.Name,
			host.FlashMode,
			host.OsName,
			host.Status.ID,
			host.Status.Kind,
			host.Status.Status))
	}

	return &b, nil
}
