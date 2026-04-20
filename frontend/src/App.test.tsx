import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import App from './App';

// Explicitly define global fetch for jsdom if missing
if (!global.fetch) {
  (global as any).fetch = jest.fn();
}

describe('App Component', () => {
  let fetchSpy: jest.SpyInstance;

  beforeEach(() => {
    fetchSpy = jest.spyOn(global, 'fetch');
    localStorage.clear();
  });

  afterEach(() => {
    fetchSpy.mockRestore();
  });

  test('renders title and button', () => {
    render(<App />);
    expect(screen.getByText(/NPC Descriptors/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Generate Descriptors/i })).toBeInTheDocument();
  });

  test('fetches and displays descriptors on button click', async () => {
    const mockDescriptors = ['Brave', 'Cunning', 'Tall'];
    fetchSpy.mockResolvedValue({
      ok: true,
      json: jest.fn().mockResolvedValue({ descriptors: mockDescriptors }),
    } as any);

    render(<App />);
    const button = screen.getByRole('button', { name: /Generate Descriptors/i });
    fireEvent.click(button);

    expect(button).toBeDisabled();
    expect(screen.getByText(/Loading.../i)).toBeInTheDocument();

    await waitFor(() => {
      mockDescriptors.forEach(d => {
        expect(screen.getByText(d)).toBeInTheDocument();
      });
    });
  });

  test('displays error message on fetch failure', async () => {
    const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
    fetchSpy.mockRejectedValue(new Error('API Error'));

    render(<App />);
    const button = screen.getByRole('button', { name: /Generate Descriptors/i });
    fireEvent.click(button);

    await waitFor(() => {
      expect(screen.getByText(/Failed to fetch descriptors/i)).toBeInTheDocument();
    });
    consoleSpy.mockRestore();
  });

  test('displays error message on HTTP error status', async () => {
    const consoleSpy = jest.spyOn(console, 'error').mockImplementation(() => {});
    fetchSpy.mockResolvedValue({
      ok: false,
      status: 500,
    } as any);

    render(<App />);
    const button = screen.getByRole('button', { name: /Generate Descriptors/i });
    fireEvent.click(button);

    await waitFor(() => {
      expect(screen.getByText(/Failed to fetch descriptors/i)).toBeInTheDocument();
    });
    consoleSpy.mockRestore();
  });

  test('fetches with correct query parameter n and default lang', async () => {
    fetchSpy.mockResolvedValue({
      ok: true,
      json: jest.fn().mockResolvedValue({ descriptors: ['Stoic'] }),
    } as any);

    render(<App />);
    const input = screen.getByLabelText(/Count \(1-10\):/i);
    const button = screen.getByRole('button', { name: /Generate Descriptors/i });

    // Change input to 5
    fireEvent.change(input, { target: { value: '5' } });
    fireEvent.click(button);

    expect(fetchSpy).toHaveBeenCalledWith('/api/descriptors?n=5&lang=en');
    
    // Wait for state updates to finish
    await waitFor(() => expect(button).not.toBeDisabled());
  });

  test('switches language and fetches with correct lang parameter', async () => {
    fetchSpy.mockResolvedValue({
      ok: true,
      json: jest.fn().mockResolvedValue({ descriptors: ['Stoïque'] }),
    } as any);

    render(<App />);
    
    // Verify initial state
    const enButton = screen.getByRole('button', { name: 'English' });
    const frButton = screen.getByRole('button', { name: 'Français' });
    expect(enButton).toHaveAttribute('aria-pressed', 'true');
    expect(frButton).toHaveAttribute('aria-pressed', 'false');

    // Switch to French
    fireEvent.click(frButton);
    expect(enButton).toHaveAttribute('aria-pressed', 'false');
    expect(frButton).toHaveAttribute('aria-pressed', 'true');

    // Check UI strings updated to French
    expect(screen.getByText(/Descripteurs de PNJ/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Nombre \(1-10\) :/i)).toBeInTheDocument();
    const genButton = screen.getByRole('button', { name: /Générer les descripteurs/i });
    expect(genButton).toBeInTheDocument();

    // Fetch in French
    fireEvent.click(genButton);
    expect(fetchSpy).toHaveBeenCalledWith('/api/descriptors?n=3&lang=fr');

    await waitFor(() => {
      expect(screen.getByText('Stoïque')).toBeInTheDocument();
      expect(genButton).not.toBeDisabled();
    });

    // Switch back to English
    fireEvent.click(enButton);
    expect(screen.getByText(/NPC Descriptors/i)).toBeInTheDocument();
  });

  test('copies descriptors to clipboard and shows success message', async () => {
    const mockDescriptors = ['Brave', 'Cunning'];
    fetchSpy.mockResolvedValue({
      ok: true,
      json: jest.fn().mockResolvedValue({ descriptors: mockDescriptors }),
    } as any);

    const writeTextMock = jest.fn().mockResolvedValue(undefined);
    Object.assign(navigator, {
      clipboard: {
        writeText: writeTextMock,
      },
    });

    render(<App />);
    const genButton = screen.getByRole('button', { name: /Generate Descriptors/i });
    fireEvent.click(genButton);

    const copyButton = await screen.findByRole('button', { name: /Copy All/i });
    fireEvent.click(copyButton);

    expect(writeTextMock).toHaveBeenCalledWith('Brave, Cunning');
    expect(await screen.findByText(/Copied!/i)).toBeInTheDocument();
  });

  test('maintains history of previous rolls', async () => {
    const roll1 = ['Brave', 'Cunning'];
    const roll2 = ['Tall', 'Short'];
    
    fetchSpy.mockResolvedValueOnce({
      ok: true,
      json: jest.fn().mockResolvedValue({ descriptors: roll1 }),
    } as any).mockResolvedValueOnce({
      ok: true,
      json: jest.fn().mockResolvedValue({ descriptors: roll2 }),
    } as any);

    render(<App />);
    const button = screen.getByRole('button', { name: /Generate Descriptors/i });

    // First roll
    fireEvent.click(button);
    await waitFor(() => expect(screen.getByText('Brave')).toBeInTheDocument());

    // Second roll
    fireEvent.click(button);
    await waitFor(() => expect(screen.getByText('Tall')).toBeInTheDocument());

    // Check history (should contain roll1)
    expect(screen.getByText('Recent Rolls')).toBeInTheDocument();
    expect(screen.getByText('Brave, Cunning')).toBeInTheDocument();
  });

  test('validates and trims history on load', () => {
    const validHistory = [
      { id: '1', descriptors: ['Current'] },
      { id: '2', descriptors: ['Previous 1'] },
      { id: '3', descriptors: ['Previous 2'] }
    ];
    const invalidHistory = [{ id: 1, descriptors: 'not an array' }];
    
    // Test valid load
    localStorage.setItem('npc_history', JSON.stringify(validHistory));
    const { unmount } = render(<App />);
    // Current is in descriptors, Previous 1 is in history list
    expect(screen.getByText('Current')).toBeInTheDocument();
    expect(screen.getByText('Previous 1')).toBeInTheDocument();
    unmount();

    // Test invalid load
    localStorage.clear();
    localStorage.setItem('npc_history', JSON.stringify(invalidHistory));
    render(<App />);
    expect(screen.queryByText('Recent Rolls')).not.toBeInTheDocument();
  });
});
