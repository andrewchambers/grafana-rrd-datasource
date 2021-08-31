import React, { useState } from 'react';
import { RRDVariableQuery } from './types';

interface VariableQueryProps {
  query: RRDVariableQuery;
  onChange: (query: RRDVariableQuery, definition: string) => void;
}

export const VariableQueryEditor: React.FC<VariableQueryProps> = ({ onChange, query }) => {
  const [state, setState] = useState(query);

  const saveQuery = () => {
    onChange(state, state.glob);
  };

  const handleChange = (event: React.FormEvent<HTMLInputElement>) =>
    setState({
      ...state,
      [event.currentTarget.name]: event.currentTarget.value,
    });

  return (
    <>
      <div className="gf-form">
        <span className="gf-form-label width-10">RRD glob</span>
        <input
          name="glob"
          placeholder="**"
          className="gf-form-input"
          onBlur={saveQuery}
          onChange={handleChange}
          value={state.glob}
        />
      </div>
    </>
  );
};
