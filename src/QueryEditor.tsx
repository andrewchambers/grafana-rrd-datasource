import { QueryEditorProps } from '@grafana/data';
import { CodeEditor, InlineFieldRow, InlineLabel } from '@grafana/ui';

import React, { ComponentType } from 'react';
import AutoSizer from 'react-virtualized-auto-sizer';
import { DataSource } from './datasource';

import { RRDDataSourceOptions, RRDQuery } from './types';

type Props = QueryEditorProps<DataSource, RRDQuery, RRDDataSourceOptions>;

interface LastQuery {
  xport: string;
}

export const QueryEditor: ComponentType<Props> = ({ datasource, onChange, onRunQuery, query }) => {
  const [xport, setXport] = React.useState(query.xport ?? '');
  const [lastQuery, setLastQuery] = React.useState<LastQuery | null>(null);

  React.useEffect(() => {
    if (lastQuery !== null && xport === lastQuery.xport) {
      return;
    }

    setLastQuery({ xport });

    onChange({ ...query, xport });

    onRunQuery();
  }, [xport, lastQuery, onChange, onRunQuery, query]);

  return (
    <>
      <InlineFieldRow>
        <AutoSizer disableHeight>
          {({ width }) => (
            <div style={{ width: width + 'px' }}>
              <InlineLabel>RRDTool (e)xport</InlineLabel>
              <CodeEditor
                width="100%"
                height="100px"
                language=""
                showLineNumbers={true}
                showMiniMap={xport.length > 100}
                value={xport}
                onBlur={value => setXport(value)}
              />
            </div>
          )}
        </AutoSizer>
      </InlineFieldRow>
    </>
  );
};
