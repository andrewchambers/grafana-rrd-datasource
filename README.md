# Grafana RRD Datasource

A grafana datasource for reading '.rrd' files via [RRDTool](https://oss.oetiker.ch/rrdtool/) and 
[RRDsrv](https://github.com/andrewchambers/rrdsrv).

With this datasource you will be able to create grafana dashboards and use grafana alerts with RRD file metrics.

[Demo video](https://www.youtube.com/watch?v=BuoPcyJik38).

## Usage

Install rrdtool and https://github.com/andrewchambers/rrdsrv.

From grafana add a new rrd datasource and point it at your rrdsrv instance.

## Thanks

Thanks to Tobias Oetiker and all other [contributors to RRDTool](https://oss.oetiker.ch/rrdtool/cast.en.html)