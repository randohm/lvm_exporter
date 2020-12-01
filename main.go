package main

import (
    "flag"
    "net/http"
    "log"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/prometheus/client_golang/prometheus"
    "os/exec"
    "encoding/json"
    "strconv"
    "strings"
    "fmt"
    "os"
)



const (
    swVersion = "0.0.1"
    swName = "lvm_exporter"
    defaultListenAddress = "0.0.0.0:9777"
    defaultMetricsPath = "/metrics"
    metricNamePrefix = "lvm_"
)



var (
    pvSizeMetric = prometheus.NewDesc("lvm_pv_bytes_size", "Shows total size of the PV in bytes", []string{"pv_name","pv_uuid","vg_name"}, nil)
    pvFreeMetric = prometheus.NewDesc("lvm_pv_bytes_free", "Shows free space of the PV in bytes", []string{"pv_name","pv_uuid","vg_name"}, nil)
    pvUsedMetric = prometheus.NewDesc("lvm_pv_bytes_used", "Shows used space of the PV in bytes", []string{"pv_name","pv_uuid","vg_name"}, nil)

    vgSizeMetric = prometheus.NewDesc("lvm_vg_bytes_size", "Shows total size of the VG in bytes", []string{"vg_name","vg_uuid"}, nil)
    vgFreeMetric = prometheus.NewDesc("lvm_vg_bytes_free", "Shows free space of the VG in bytes", []string{"vg_name","vg_uuid"}, nil)

    lvSizeMetric = prometheus.NewDesc("lvm_lv_bytes_size", "Shows total size of the LV in bytes", []string{"lv_name","lv_uuid","vg_name"}, nil)
)



type LvmExporter struct {
    pvs     []map[string]string
    vgs     []map[string]string
    lvs     []map[string]string
}



func NewLvmExporter() (*LvmExporter, error) {
    return &LvmExporter{}, nil
}



func (e *LvmExporter) Describe(ch chan<- *prometheus.Desc) {
    ch <- pvSizeMetric
    ch <- pvFreeMetric
    ch <- pvUsedMetric
    ch <- vgSizeMetric
    ch <- vgFreeMetric
    ch <- lvSizeMetric
}



func (e *LvmExporter) Collect(ch chan<- prometheus.Metric) {
    // Get the names of the VGs first, then get full reports for each VG
    var vgNames []string
    var vgs map[string][]map[string][]map[string]string
    vgs_json, err := exec.Command("/usr/sbin/vgs", "--verbose", "--units", "b", "--reportformat", "json").Output()
    if err != nil {
        log.Print(err)
        return
    }
    err = json.Unmarshal(vgs_json, &vgs)
    if err != nil {
        log.Print(err)
        return
    }

    for _, v := range vgs["report"][0]["vg"] {
        vgNames = append(vgNames, v["vg_name"])
    }

    for _, vgName := range vgNames {
        var report map[string][]map[string][]map[string]string
        report_json, err := exec.Command("/usr/sbin/lvm", "fullreport", "--units", "b", "--reportformat", "json", vgName).Output()
        if err != nil {
            log.Print(err)
            return
        }
        err = json.Unmarshal(report_json, &report)
        if err != nil {
            log.Print(err)
            return
        }

        pvCollect(ch, report["report"][0]["pv"], vgName)
        vgCollect(ch, report["report"][0]["vg"])
        lvCollect(ch, report["report"][0]["lv"], vgName)
    }
}



// Collects metrics from PVs
func pvCollect(ch chan<- prometheus.Metric, pvs []map[string]string, vgName string) {
    for _, pv := range  pvs {
        pvSizeF, err := strconv.ParseFloat(strings.Trim(pv["pv_size"], "B"), 64)
        if err != nil {
            log.Print(err)
            return
        }
        ch <- prometheus.MustNewConstMetric(pvSizeMetric, prometheus.GaugeValue, pvSizeF, pv["pv_name"], pv["pv_uuid"], vgName)

        pvFreeF, err := strconv.ParseFloat(strings.Trim(pv["pv_free"], "B"), 64)
        if err != nil {
            log.Print(err)
            return
        }
        ch <- prometheus.MustNewConstMetric(pvFreeMetric, prometheus.GaugeValue, pvFreeF, pv["pv_name"], pv["pv_uuid"], vgName)

        pvUsedF, err := strconv.ParseFloat(strings.Trim(pv["pv_used"], "B"), 64)
        if err != nil {
            log.Print(err)
            return
        }
        ch <- prometheus.MustNewConstMetric(pvUsedMetric, prometheus.GaugeValue, pvUsedF, pv["pv_name"], pv["pv_uuid"], vgName)
    }
}



// Collects metrics from VGs
func vgCollect(ch chan<- prometheus.Metric, vgs []map[string]string) {
    for _, vg := range  vgs {
        vgSizeF, err := strconv.ParseFloat(strings.Trim(vg["vg_size"], "B"), 64)
        if err != nil {
            log.Print(err)
            return
        }
        ch <- prometheus.MustNewConstMetric(vgSizeMetric, prometheus.GaugeValue, vgSizeF, vg["vg_name"], vg["vg_uuid"])

        vgFreeF, err := strconv.ParseFloat(strings.Trim(vg["vg_free"], "B"), 64)
        if err != nil {
            log.Print(err)
            return
        }
        ch <- prometheus.MustNewConstMetric(vgFreeMetric, prometheus.GaugeValue, vgFreeF, vg["vg_name"], vg["vg_uuid"])
    }
}



// Collects metrics from LVs
func lvCollect(ch chan<- prometheus.Metric, lvs []map[string]string, vgName string) {
    for _, lv := range  lvs {
        lvSizeF, err := strconv.ParseFloat(strings.Trim(lv["lv_size"], "B"), 64)
        if err != nil {
            log.Print(err)
            return
        }
        ch <- prometheus.MustNewConstMetric(lvSizeMetric, prometheus.GaugeValue, lvSizeF, lv["lv_name"], lv["lv_uuid"], vgName)
    }
}



func init() {
    log.SetFlags(log.Ldate|log.Ltime|log.Lshortfile)
}



func main() {
    listenAddress := flag.String("web.listen-address", defaultListenAddress, "Listen address for HTTP requests")
    metricsPath := flag.String("web.telemetry-path", defaultMetricsPath, "Path under which to expose metrics")
    showVersion := flag.Bool("version", false, "Show version and exit")
    flag.Parse()

    if *showVersion {
        fmt.Printf("%s v%s\n", swName, swVersion)
        os.Exit(0)
    }


    exporter, err := NewLvmExporter()
    if err != nil {
        log.Fatal(err)
    }
    prometheus.MustRegister(exporter)

    http.Handle(*metricsPath, promhttp.Handler())
    log.Printf("Listening on %s", *listenAddress)
    http.ListenAndServe(*listenAddress, nil)
}
