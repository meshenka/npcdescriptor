import React, { useState } from 'react';

interface DescriptorResponse {
  descriptors: string[];
}

const App: React.FC = () => {
  const [descriptors, setDescriptors] = useState<string[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');

  const fetchDescriptors = async (): Promise<void> => {
    setLoading(true);
    setError('');

    try {
      const response = await fetch('/api/descriptors');
      const data: DescriptorResponse = await response.json();
      setDescriptors(data.descriptors);
    } catch (err) {
      setError('Failed to fetch descriptors.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="app-container">
      <h1 className="app-title">NPC Descriptors</h1>
      <button
        onClick={() => fetchDescriptors()}
        disabled={loading}
        className="btn"
      >
        {loading ? 'Loading...' : 'Generate Descriptors'}
      </button>

      {descriptors.length > 0 && (
        <div className="sentence-container">
          <ul style={{ listStyle: 'none', padding: 0 }}>
            {descriptors.map((d, i) => (
              <li key={i} className="sentence-text" style={{ marginBottom: '0.5rem' }}>
                {d}
              </li>
            ))}
          </ul>
        </div>
      )}

      {error && (
        <div className="error-container">
          <p>{error}</p>
        </div>
      )}
    </div>
  );
};

export default App;
