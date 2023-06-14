import React, { ChangeEvent } from 'react';
import { Input, Field, Button, TextArea, FieldSet, InlineFieldRow, InlineField } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, MyQuery } from '../types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  const onTestQuery = () => {
    onRunQuery();
  };

  const onDatabaseChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, Db: event.target.value });
  };

  const onCollectionChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, Collection: event.target.value });
  };

  const onAggregationChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    onChange({ ...query, Aggregation: event.target.value });
  };

  const { Db, Collection, Aggregation } = query;

  return (
    // <div className="gf-form" style={{display: "flex", flexDirection: "column"}}>
    <div className="gf-form">
      <FieldSet>
        <InlineFieldRow style={{ flexDirection: 'row-reverse' }}>
          <Button onClick={onTestQuery}>Run query</Button>
        </InlineFieldRow>

        <InlineFieldRow>
          <InlineField label="Database">
            <Input onChange={onDatabaseChange} value={Db || ''} />
          </InlineField>

          <InlineField label="Collection">
            <Input onChange={onCollectionChange} value={Collection || ''} />
          </InlineField>
        </InlineFieldRow>

        <InlineFieldRow>
          <Field label="Aggregation">
            <TextArea onChange={onAggregationChange} rows={10} cols={100}>
              {Aggregation || ''}
            </TextArea>
          </Field>
        </InlineFieldRow>
      </FieldSet>
    </div>
  );
}
