import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface RRDQuery extends DataQuery {
  xport: string;
}

export interface RRDVariableQuery {
  glob: string;
}

export interface RRDDataSourceOptions extends DataSourceJsonData {
  /* nothing */
}

export interface RRDSecureJsonData {
  /* nothing */
}
