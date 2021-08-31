import React, { ComponentType } from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { RRDDataSourceOptions } from './types';
import { DataSourceHttpSettings } from '@grafana/ui';

type Props = DataSourcePluginOptionsEditorProps<RRDDataSourceOptions>;

export const ConfigEditor: ComponentType<Props> = ({ options, onOptionsChange }) => {
  const defaultUrl = 'http://localhost:9191';

  return (
    <>
      <DataSourceHttpSettings
        defaultUrl={defaultUrl}
        dataSourceConfig={options}
        showAccessOptions={false}
        onChange={onOptionsChange}
      />
    </>
  );
};
