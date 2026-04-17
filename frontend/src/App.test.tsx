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

    expect(button).not.toBeDisabled();
  });

  test('displays error message on fetch failure', async () => {
    fetchSpy.mockRejectedValue(new Error('API Error'));

    render(<App />);
    const button = screen.getByRole('button', { name: /Generate Descriptors/i });
    fireEvent.click(button);

    await waitFor(() => {
      expect(screen.getByText(/Failed to fetch descriptors/i)).toBeInTheDocument();
    });
  });

  test('displays error message on HTTP error status', async () => {
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
});
