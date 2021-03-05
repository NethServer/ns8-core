# nsdc

This module runs as singleton rootfull instance.

Initialize the Redis DB with

    HSET nsdc/module.env NSDC_IPADDRESS 10.133.0.5 EVENTS_IMAGE ghcr.io/nethserver/cplane-nsdc NSDC_HOSTNAME nsdc1.ad.dp.nethserver.net NSDC_NBDOMAIN AD NSDC_REALM AD.DP.NETHSERVER.NET NSDC_ADMINPASS Nethesis,1234
