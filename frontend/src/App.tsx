import React, { useState, useRef, useEffect, useCallback, useMemo } from 'react';
import { Locale, translations } from './i18n';
import usFlag from './assets/us.png';
import frFlag from './assets/fr.png';

interface DescriptorResponse {
  descriptors: string[];
}

interface HistoryEntry {
  id: string;
  descriptors: string[];
}

const App: React.FC = () => {
  const [locale, setLocale] = useState<Locale>(() => {
    try {
      const saved = localStorage.getItem('npc_locale');
      return (saved === 'en' || saved === 'fr') ? saved : 'en';
    } catch (e) {
      return 'en';
    }
  });
  const [history, setHistory] = useState<HistoryEntry[]>(() => {
    try {
      const saved = localStorage.getItem('npc_history');
      if (!saved) return [];
      const parsed = JSON.parse(saved);
      if (Array.isArray(parsed) && parsed.every(entry => 
        entry && typeof entry === 'object' && 
        typeof entry.id === 'string' && 
        Array.isArray(entry.descriptors) && 
        entry.descriptors.every((d: any) => typeof d === 'string')
      )) {
        return parsed.slice(0, 10);
      }
      return [];
    } catch (e) {
      return [];
    }
  });
  const [descriptors, setDescriptors] = useState<string[]>(() => history[0]?.descriptors || []);
  const [count, setCount] = useState<number>(() => {
    try {
      const saved = localStorage.getItem('npc_count');
      const parsed = parseInt(saved || '');
      return (!isNaN(parsed) && parsed >= 1 && parsed <= 10) ? parsed : 3;
    } catch (e) {
      return 3;
    }
  });
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');
  const [copied, setCopied] = useState<boolean>(false);
  const copyTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const t = useMemo(() => translations[locale], [locale]);

  useEffect(() => {
    try {
      localStorage.setItem('npc_locale', locale);
    } catch (e) {
      // Ignore storage errors
    }
  }, [locale]);

  useEffect(() => {
    try {
      localStorage.setItem('npc_count', count.toString());
    } catch (e) {
      // Ignore storage errors
    }
  }, [count]);

  useEffect(() => {
    try {
      localStorage.setItem('npc_history', JSON.stringify(history));
    } catch (e) {
      // Ignore storage errors
    }
  }, [history]);

  useEffect(() => {
    return () => {
      if (copyTimeoutRef.current) {
        clearTimeout(copyTimeoutRef.current);
      }
    };
  }, []);

  const fetchDescriptors = useCallback(async (): Promise<void> => {
    const startTime = Date.now();
    setLoading(true);
    setError('');
    setCopied(false);
    if (copyTimeoutRef.current) {
      clearTimeout(copyTimeoutRef.current);
    }

    try {
      const response = await fetch(`/api/descriptors?n=${count}&lang=${locale}`);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data: DescriptorResponse = await response.json();
      
      // Prevent blink by ensuring at least 300ms loading state
      const elapsed = Date.now() - startTime;
      if (elapsed < 300) {
        await new Promise(resolve => setTimeout(resolve, 300 - elapsed));
      }
      
      setDescriptors(data.descriptors);
      setHistory(prev => [{ id: Date.now().toString(), descriptors: data.descriptors }, ...prev].slice(0, 10));
    } catch (err) {
      setError(t.fetchError);
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [count, locale, t.fetchError]);

  const copyToClipboard = async () => {
    const text = descriptors.join(', ');
    
    if (!navigator.clipboard) {
      setError(t.clipboardError);
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
      setError(t.copyError);
      console.error('Clipboard error:', err);
    }
  };

  return (
    <div className="app-container">
      <div className="locale-selector" style={{ position: 'absolute', top: '1rem', right: '1rem', display: 'flex', gap: '0.5rem' }}>
        <button 
          onClick={() => setLocale('en')} 
          style={{ 
            background: 'none', 
            border: locale === 'en' ? '2px solid #3b82f6' : '2px solid transparent', 
            borderRadius: '4px',
            cursor: 'pointer',
            padding: '2px',
            display: 'flex',
            alignItems: 'center'
          }}
          title="English"
          aria-label="English"
          aria-pressed={locale === 'en'}
        >
          <img src={usFlag} width="30" height="20" style={{ aspectRatio: '3/2' }} alt="USA Flag" />
        </button>
        <button 
          onClick={() => setLocale('fr')} 
          style={{ 
            background: 'none', 
            border: locale === 'fr' ? '2px solid #3b82f6' : '2px solid transparent', 
            borderRadius: '4px',
            cursor: 'pointer',
            padding: '2px',
            display: 'flex',
            alignItems: 'center'
          }}
          title="Français"
          aria-label="Français"
          aria-pressed={locale === 'fr'}
        >
          <img src={frFlag} width="30" height="20" style={{ aspectRatio: '3/2' }} alt="French Flag" />
        </button>
      </div>

      <h1 className="app-title">{t.title}</h1>
      
      <div style={{ marginBottom: '1rem' }}>
        <label htmlFor="count-input" style={{ marginRight: '0.5rem' }}>{t.countLabel}</label>
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
        {loading ? t.loading : t.generateBtn}
      </button>

      <div className="sentence-container">
        {descriptors.length > 0 ? (
          <>
            <ul style={{ listStyle: 'none', padding: 0 }}>
              {descriptors.map((d, i) => (
                <li key={i} className="sentence-text" style={{ marginBottom: '0.5rem' }}>
                  {d}
                </li>
              ))}
            </ul>
            <button onClick={copyToClipboard} className="btn" style={{ marginTop: '1rem', backgroundColor: '#4CAF50', display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '0.5rem' }}>
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <rect width="14" height="14" x="8" y="8" rx="2" ry="2"/><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"/>
              </svg>
              {copied ? t.copied : t.copyBtn}
            </button>
          </>
        ) : (
          <p style={{ textAlign: 'center', color: '#999', fontStyle: 'italic' }}>
            {/* Empty state placeholder */}
          </p>
        )}
      </div>

      {error && (
        <div className="error-container">
          <p>{error}</p>
        </div>
      )}

      {history.length > 1 && (
        <div className="history-container" style={{ marginTop: '2rem', borderTop: '1px solid #ccc', paddingTop: '1rem', width: '100%', maxWidth: '500px' }}>
          <h3 style={{ marginBottom: '1rem' }}>{t.historyTitle}</h3>
          <ul style={{ listStyle: 'none', padding: 0 }}>
            {history.slice(1).map((entry) => (
              <li key={entry.id} style={{ marginBottom: '0.5rem', padding: '0.5rem', borderBottom: '1px solid #eee', fontSize: '0.9rem', color: '#666' }}>
                {entry.descriptors.join(', ')}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default App;
