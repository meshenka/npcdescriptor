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
});
