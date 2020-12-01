# lvm_exporter
[Prometheus](https://prometheus.io) exporter for LVM metrics.

Currently calls `pvs`, `vgs`, and `lvs` to get data and exposes them as a Prometheus exporter.
Future plans include making direct calls to the `liblvm`.

## Metrics

Currently only gets metrics for volume size and use.

```
# HELP lvm_lv_bytes_size Shows total size of the LV in bytes
# TYPE lvm_lv_bytes_size gauge
lvm_lv_bytes_size{lv_name="home",lv_uuid="m2sdfN-T234-Dvsd-HfsJ-7fsd-6dfd-IDyreZ",vg_name="centos"} 4.48744390656e+11
lvm_lv_bytes_size{lv_name="root",lv_uuid="xHGsdc-4sdf-edfG-osdZ-vaPI-ogsd-ISgsdL",vg_name="centos"} 5.36870912e+10
lvm_lv_bytes_size{lv_name="swap",lv_uuid="nFsdfb-8sdZ-edsE-gsdc-Xg7r-gsgR-Z6r3ot",vg_name="centos"} 8.388608e+09
# HELP lvm_pv_bytes_free Shows free space of the PV in bytes
# TYPE lvm_pv_bytes_free gauge
lvm_pv_bytes_free{pv_name="/dev/sda3",pv_uuid="Sasdfa-oasQ-tdg1-r34p-jj43-S43g-gK343Y",vg_name="centos"} 4.194304e+06
# HELP lvm_pv_bytes_size Shows total size of the PV in bytes
# TYPE lvm_pv_bytes_size gauge
lvm_pv_bytes_size{pv_name="/dev/sda3",pv_uuid="Sasdfa-oasQ-tdg1-r34p-jj43-S43g-gK343Y",vg_name="centos"} 5.1082428416e+11
# HELP lvm_pv_bytes_used Shows used space of the PV in bytes
# TYPE lvm_pv_bytes_used gauge
lvm_pv_bytes_used{pv_name="/dev/sda3",pv_uuid="Sasdfa-oasQ-tdg1-r34p-jj43-S43g-gK343Y",vg_name="centos"} 5.10820089856e+11
# HELP lvm_vg_bytes_free Shows free space of the VG in bytes
# TYPE lvm_vg_bytes_free gauge
lvm_vg_bytes_free{vg_name="centos",vg_uuid="23rFt7-ddgQ-q356-c5fh-tw34-7bsi-2y3sgm"} 4.194304e+06
# HELP lvm_vg_bytes_size Shows total size of the VG in bytes
# TYPE lvm_vg_bytes_size gauge
lvm_vg_bytes_size{vg_name="centos",vg_uuid="23rFt7-ddgQ-q356-c5fh-tw34-7bsi-2y3sgm"} 5.1082428416e+11
```

