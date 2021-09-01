# Grafana RRD Datasource

A grafana datasource for reading '.rrd' files via [RRDTool](https://oss.oetiker.ch/rrdtool/) and 
[RRDsrv](https://github.com/andrewchambers/rrdsrv).

With this datasource you will be able to create grafana dashboards and use grafana alerts with RRD file metrics.

[Demo video](https://www.youtube.com/watch?v=BuoPcyJik38).

## Getting started

A data source backend plugin consists of both frontend and backend components.

### Frontend

1. Install dependencies

   ```bash
   yarn install
   ```

2. Build plugin in development mode or run in watch mode

   ```bash
   yarn dev
   ```

   or

   ```bash
   yarn watch
   ```

3. Build plugin in production mode

   ```bash
   yarn build
   ```

### Backend

1. Update [Grafana plugin SDK for Go](https://grafana.com/docs/grafana/latest/developers/plugins/backend/grafana-plugin-sdk-for-go/) dependency to the latest minor version:

   ```bash
   go get -u github.com/grafana/grafana-plugin-sdk-go
   go mod tidy
   ```

2. Build backend plugin binaries for Linux, Windows and Darwin:

   ```bash
   mage -v
   ```

3. List all available Mage targets for additional commands:

   ```bash
   mage -l
   ```

### Usage

Install https://github.com/andrewchambers/rrdsrv.

In a terminal start rrdsrv:
```
$ rrdsrv
listening on localhost:9191
```

From grafana add a new datasource and point it to http://localhost:9191.

## Thanks

Thanks to Tobias Oetiker and all other [contributors to RRDTool](https://oss.oetiker.ch/rrdtool/cast.en.html)