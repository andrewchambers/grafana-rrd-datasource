import { DataSourceInstanceSettings, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { RRDDataSourceOptions, RRDQuery, RRDVariableQuery } from './types';

export class DataSource extends DataSourceWithBackend<RRDQuery, RRDDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<RRDDataSourceOptions>) {
    super(instanceSettings);
  }

  applyTemplateVariables(query: RRDQuery, scopedVars: ScopedVars): RRDQuery {
    query.xport = getTemplateSrv().replace(query.xport);
    return query;
  }

  async metricFindQuery(query: RRDVariableQuery, options?: any) {
    const params = query.glob ? null : { glob: query.glob };
    const response: string[] = await this.getResource('/api/v1/list_metrics', params);
    return response.map((rrd: string) => {
      return { text: rrd };
    });
  }
}
