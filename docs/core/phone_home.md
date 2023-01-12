---
layout: default
title: Phone Home
nav_order: 15
parent: Core
---

# Phone Home

We're committed to being more aware of our products, how they are used and how they behave in different conditions and hardware.

To do so, we actually need to expand how the Phone Home worked in NS6-NS7 by adding additional telemetry to be sent and processed for analysis.

Starting from NS8, the Phone Home will be activated by default, to disable  this behaviour, simply disable the `phonehome.timer` on leader node of the cluster:

```bash
systemctl disable phonehome.timer
```

To find what data is sent, you can run:
```bash
api-cli run get-facts
```
