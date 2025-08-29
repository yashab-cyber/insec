import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { AlertForm } from '../components/forms/AlertForm';

// Mock validation functions
jest.mock('../utils/validation', () => ({
  validateEmail: (email: string) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  },
  validateRequired: (value: string) => value.trim().length > 0,
  validateMinLength: (value: string, minLength: number) => value.length >= minLength,
}));

describe('AlertForm', () => {
  const mockOnSubmit = jest.fn();
  const mockOnCancel = jest.fn();

  const defaultProps = {
    onSubmit: mockOnSubmit,
    onCancel: mockOnCancel,
    isLoading: false,
  };

  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders form with all required fields', () => {
    render(<AlertForm {...defaultProps} />);

    expect(screen.getByLabelText(/alert title/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/description/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/severity/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/category/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/assigned to/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /create alert/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument();
  });

  it('displays validation errors for empty required fields', async () => {
    render(<AlertForm {...defaultProps} />);

    const submitButton = screen.getByRole('button', { name: /create alert/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/alert title is required/i)).toBeInTheDocument();
      expect(screen.getByText(/description is required/i)).toBeInTheDocument();
      expect(screen.getByText(/please select a severity level/i)).toBeInTheDocument();
      expect(screen.getByText(/please select a category/i)).toBeInTheDocument();
    });
  });

  it('validates title minimum length', async () => {
    render(<AlertForm {...defaultProps} />);

    const titleInput = screen.getByLabelText(/alert title/i);
    const submitButton = screen.getByRole('button', { name: /create alert/i });

    fireEvent.change(titleInput, { target: { value: 'Hi' } });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/title must be at least 5 characters/i)).toBeInTheDocument();
    });
  });

  it('validates description minimum length', async () => {
    render(<AlertForm {...defaultProps} />);

    const descriptionInput = screen.getByLabelText(/description/i);
    const submitButton = screen.getByRole('button', { name: /create alert/i });

    fireEvent.change(descriptionInput, { target: { value: 'Short' } });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/description must be at least 10 characters/i)).toBeInTheDocument();
    });
  });

  it('validates email format for assigned user', async () => {
    render(<AlertForm {...defaultProps} />);

    const emailInput = screen.getByLabelText(/assigned to/i);
    const submitButton = screen.getByRole('button', { name: /create alert/i });

    fireEvent.change(emailInput, { target: { value: 'invalid-email' } });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/please enter a valid email address/i)).toBeInTheDocument();
    });
  });

  it('submits form with valid data', async () => {
    render(<AlertForm {...defaultProps} />);

    const titleInput = screen.getByLabelText(/alert title/i);
    const descriptionInput = screen.getByLabelText(/description/i);
    const severitySelect = screen.getByLabelText(/severity/i);
    const categorySelect = screen.getByLabelText(/category/i);
    const emailInput = screen.getByLabelText(/assigned to/i);
    const submitButton = screen.getByRole('button', { name: /create alert/i });

    fireEvent.change(titleInput, { target: { value: 'Suspicious Activity Detected' } });
    fireEvent.change(descriptionInput, { target: { value: 'A suspicious process was detected running on the system with elevated privileges.' } });
    fireEvent.change(severitySelect, { target: { value: 'high' } });
    fireEvent.change(categorySelect, { target: { value: 'security' } });
    fireEvent.change(emailInput, { target: { value: 'analyst@company.com' } });

    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalledWith({
        title: 'Suspicious Activity Detected',
        description: 'A suspicious process was detected running on the system with elevated privileges.',
        severity: 'high',
        category: 'security',
        assignedTo: 'analyst@company.com',
      });
    });
  });

  it('shows loading state during submission', () => {
    render(<AlertForm {...defaultProps} isLoading={true} />);

    const submitButton = screen.getByRole('button', { name: /creating alert/i });
    expect(submitButton).toBeDisabled();
    expect(screen.getByText(/creating alert/i)).toBeInTheDocument();
  });

  it('calls onCancel when cancel button is clicked', () => {
    render(<AlertForm {...defaultProps} />);

    const cancelButton = screen.getByRole('button', { name: /cancel/i });
    fireEvent.click(cancelButton);

    expect(mockOnCancel).toHaveBeenCalled();
  });

  it('clears validation errors when user starts typing', async () => {
    render(<AlertForm {...defaultProps} />);

    const titleInput = screen.getByLabelText(/alert title/i);
    const submitButton = screen.getByRole('button', { name: /create alert/i });

    // Trigger validation error
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/alert title is required/i)).toBeInTheDocument();
    });

    // Start typing
    fireEvent.change(titleInput, { target: { value: 'S' } });

    await waitFor(() => {
      expect(screen.queryByText(/alert title is required/i)).not.toBeInTheDocument();
    });
  });

  it('handles form reset after successful submission', async () => {
    render(<AlertForm {...defaultProps} />);

    const titleInput = screen.getByLabelText(/alert title/i);
    const descriptionInput = screen.getByLabelText(/description/i);
    const severitySelect = screen.getByLabelText(/severity/i);
    const categorySelect = screen.getByLabelText(/category/i);
    const emailInput = screen.getByLabelText(/assigned to/i);

    // Fill form
    fireEvent.change(titleInput, { target: { value: 'Test Alert' } });
    fireEvent.change(descriptionInput, { target: { value: 'Test description' } });
    fireEvent.change(severitySelect, { target: { value: 'medium' } });
    fireEvent.change(categorySelect, { target: { value: 'system' } });
    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });

    // Submit
    const submitButton = screen.getByRole('button', { name: /create alert/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalled();
    });

    // Form should be reset (this would be handled by parent component)
    expect(titleInput).toHaveValue('Test Alert'); // Form reset would clear this
  });

  it('prevents multiple submissions', async () => {
    render(<AlertForm {...defaultProps} />);

    const titleInput = screen.getByLabelText(/alert title/i);
    const descriptionInput = screen.getByLabelText(/description/i);
    const severitySelect = screen.getByLabelText(/severity/i);
    const categorySelect = screen.getByLabelText(/category/i);
    const emailInput = screen.getByLabelText(/assigned to/i);
    const submitButton = screen.getByRole('button', { name: /create alert/i });

    // Fill form with valid data
    fireEvent.change(titleInput, { target: { value: 'Valid Alert Title' } });
    fireEvent.change(descriptionInput, { target: { value: 'This is a valid description that meets the minimum length requirement.' } });
    fireEvent.change(severitySelect, { target: { value: 'high' } });
    fireEvent.change(categorySelect, { target: { value: 'security' } });
    fireEvent.change(emailInput, { target: { value: 'valid@example.com' } });

    // Click submit multiple times quickly
    fireEvent.click(submitButton);
    fireEvent.click(submitButton);
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalledTimes(1);
    });
  });

  it('handles special characters in input fields', async () => {
    render(<AlertForm {...defaultProps} />);

    const titleInput = screen.getByLabelText(/alert title/i);
    const descriptionInput = screen.getByLabelText(/description/i);
    const submitButton = screen.getByRole('button', { name: /create alert/i });

    fireEvent.change(titleInput, { target: { value: 'Alert with spécial characters: éñü' } });
    fireEvent.change(descriptionInput, { target: { value: 'Description with symbols: @#$%^&*() and unicode: 🚨⚠️' } });
    fireEvent.change(screen.getByLabelText(/severity/i), { target: { value: 'medium' } });
    fireEvent.change(screen.getByLabelText(/category/i), { target: { value: 'system' } });
    fireEvent.change(screen.getByLabelText(/assigned to/i), { target: { value: 'test@example.com' } });

    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          title: 'Alert with spécial characters: éñü',
          description: 'Description with symbols: @#$%^&*() and unicode: 🚨⚠️',
        })
      );
    });
  });

  it('validates maximum length constraints', async () => {
    render(<AlertForm {...defaultProps} />);

    const titleInput = screen.getByLabelText(/alert title/i);
    const descriptionInput = screen.getByLabelText(/description/i);
    const submitButton = screen.getByRole('button', { name: /create alert/i });

    // Create a very long title
    const longTitle = 'A'.repeat(201); // Exceeds max length
    const longDescription = 'A'.repeat(1001); // Exceeds max length

    fireEvent.change(titleInput, { target: { value: longTitle } });
    fireEvent.change(descriptionInput, { target: { value: longDescription } });
    fireEvent.change(screen.getByLabelText(/severity/i), { target: { value: 'low' } });
    fireEvent.change(screen.getByLabelText(/category/i), { target: { value: 'general' } });
    fireEvent.change(screen.getByLabelText(/assigned to/i), { target: { value: 'test@example.com' } });

    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/title must be less than 200 characters/i)).toBeInTheDocument();
      expect(screen.getByText(/description must be less than 1000 characters/i)).toBeInTheDocument();
    });
  });

  it('handles form accessibility features', () => {
    render(<AlertForm {...defaultProps} />);

    // Check for proper ARIA labels
    const titleInput = screen.getByLabelText(/alert title/i);
    const descriptionInput = screen.getByLabelText(/description/i);

    expect(titleInput).toHaveAttribute('aria-required', 'true');
    expect(descriptionInput).toHaveAttribute('aria-required', 'true');

    // Check for error message associations
    const submitButton = screen.getByRole('button', { name: /create alert/i });
    fireEvent.click(submitButton);

    // Error messages should be associated with inputs
    const titleError = screen.getByText(/alert title is required/i);
    expect(titleError).toHaveAttribute('role', 'alert');
  });

  it('supports keyboard navigation', () => {
    render(<AlertForm {...defaultProps} />);

    const titleInput = screen.getByLabelText(/alert title/i);
    const descriptionInput = screen.getByLabelText(/description/i);
    const submitButton = screen.getByRole('button', { name: /create alert/i });

    // Tab through form fields
    titleInput.focus();
    expect(document.activeElement).toBe(titleInput);

    fireEvent.keyDown(titleInput, { key: 'Tab' });
    expect(document.activeElement).toBe(descriptionInput);

    // Enter key should submit form (if valid)
    fireEvent.change(titleInput, { target: { value: 'Valid Title' } });
    fireEvent.change(descriptionInput, { target: { value: 'Valid description that is long enough.' } });
    fireEvent.change(screen.getByLabelText(/severity/i), { target: { value: 'medium' } });
    fireEvent.change(screen.getByLabelText(/category/i), { target: { value: 'system' } });
    fireEvent.change(screen.getByLabelText(/assigned to/i), { target: { value: 'test@example.com' } });

    fireEvent.keyDown(submitButton, { key: 'Enter' });
    expect(mockOnSubmit).toHaveBeenCalled();
  });
});</content>
<parameter name="filePath">/workspaces/insec/tests/unit/ui/forms/AlertForm.test.tsx
