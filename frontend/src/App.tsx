import React, { useState, useRef, useEffect } from 'react';

interface DescriptorResponse {
  descriptors: string[];
}

const App: React.FC = () => {
  const [descriptors, setDescriptors] = useState<string[]>([]);
  const [count, setCount] = useState<number>(3);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');
  const [copied, setCopied] = useState<boolean>(false);
  const copyTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    return () => {
      if (copyTimeoutRef.current) {
        clearTimeout(copyTimeoutRef.current);
      }
    };
  }, []);

  const fetchDescriptors = async (): Promise<void> => {
    setLoading(true);
    setError('');
    setCopied(false);
    if (copyTimeoutRef.current) {
      clearTimeout(copyTimeoutRef.current);
    }

    try {
      const response = await fetch(`/api/descriptors?n=${count}`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data: DescriptorResponse = await response.json();
      setDescriptors(data.descriptors);
    } catch (err) {
      setError('Failed to fetch descriptors.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = async () => {
    const text = descriptors.join(', ');
    
    if (!navigator.clipboard) {
      setError('Clipboard API not available (requires HTTPS or modern browser).');
      return;
    }

    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      if (copyTimeoutRef.current) {
        clearTimeout(copyTimeoutRef.current);
      }
      copyTimeoutRef.current = setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      setError('Failed to copy to clipboard.');
      console.error('Clipboard error:', err);
    }
  };

  return (
    <div className="app-container">
      <h1 className="app-title">NPC Descriptors</h1>
      
      <div style={{ marginBottom: '1rem' }}>
        <label htmlFor="count-input" style={{ marginRight: '0.5rem' }}>Count (1-10):</label>
        <input
          id="count-input"
          type="number"
          min="1"
          max="10"
          value={count}
          onChange={(e) => setCount(Math.min(10, Math.max(1, parseInt(e.target.value) || 1)))}
          style={{ width: '50px', padding: '5px' }}
        />
      </div>

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
          <button onClick={copyToClipboard} className="btn" style={{ marginTop: '1rem', backgroundColor: '#4CAF50' }}>
            {copied ? 'Copied!' : 'Copy All'}
          </button>
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
