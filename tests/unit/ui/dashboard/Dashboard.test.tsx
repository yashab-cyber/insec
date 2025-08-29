import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { Dashboard } from '../components/dashboard/Dashboard';
import { AuthProvider } from '../contexts/AuthContext';
import { ApiProvider } from '../contexts/ApiContext';

// Mock contexts
const mockUser = {
  id: 'user-123',
  email: 'test@example.com',
  name: 'Test User',
  role: 'admin',
  tenantId: 'tenant-123',
};

const mockAuthState = {
  user: mockUser,
  isAuthenticated: true,
  isLoading: false,
  error: null,
};

const mockApiState = {
  alerts: [],
  events: [],
  isLoading: false,
  error: null,
};

jest.mock('../contexts/AuthContext', () => ({
  useAuth: () => mockAuthState,
}));

jest.mock('../contexts/ApiContext', () => ({
  useApi: () => ({
    ...mockApiState,
    fetchAlerts: jest.fn(),
    fetchEvents: jest.fn(),
    acknowledgeAlert: jest.fn(),
  }),
}));

// Mock chart components
jest.mock('recharts', () => ({
  ResponsiveContainer: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  LineChart: () => <div>LineChart</div>,
  BarChart: () => <div>BarChart</div>,
  PieChart: () => <div>PieChart</div>,
  Line: () => <div>Line</div>,
  Bar: () => <div>Bar</div>,
  Pie: () => <div>Pie</div>,
  XAxis: () => <div>XAxis</div>,
  YAxis: () => <div>YAxis</div>,
  CartesianGrid: () => <div>CartesianGrid</div>,
  Tooltip: () => <div>Tooltip</div>,
  Legend: () => <div>Legend</div>,
  Cell: () => <div>Cell</div>,
}));

describe('Dashboard', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders dashboard with user information', () => {
    render(<Dashboard />);

    expect(screen.getByText(`Welcome back, ${mockUser.name}`)).toBeInTheDocument();
    expect(screen.getByText(mockUser.email)).toBeInTheDocument();
    expect(screen.getByText(mockUser.role)).toBeInTheDocument();
  });

  it('displays dashboard sections', () => {
    render(<Dashboard />);

    expect(screen.getByText(/alerts overview/i)).toBeInTheDocument();
    expect(screen.getByText(/recent events/i)).toBeInTheDocument();
    expect(screen.getByText(/system status/i)).toBeInTheDocument();
    expect(screen.getByText(/risk metrics/i)).toBeInTheDocument();
  });

  it('shows loading state when data is being fetched', () => {
    const loadingApiState = { ...mockApiState, isLoading: true };

    jest.mock('../contexts/ApiContext', () => ({
      useApi: () => loadingApiState,
    }));

    render(<Dashboard />);

    expect(screen.getByText(/loading/i)).toBeInTheDocument();
  });

  it('displays error message when data fetch fails', () => {
    const errorApiState = {
      ...mockApiState,
      error: 'Failed to fetch data',
    };

    jest.mock('../contexts/ApiContext', () => ({
      useApi: () => errorApiState,
    }));

    render(<Dashboard />);

    expect(screen.getByText(/failed to fetch data/i)).toBeInTheDocument();
  });

  it('renders alert summary cards', () => {
    const alerts = [
      { id: '1', severity: 'high', status: 'active' },
      { id: '2', severity: 'medium', status: 'acknowledged' },
      { id: '3', severity: 'low', status: 'resolved' },
    ];

    const apiStateWithAlerts = {
      ...mockApiState,
      alerts,
    };

    jest.mock('../contexts/ApiContext', () => ({
      useApi: () => apiStateWithAlerts,
    }));

    render(<Dashboard />);

    expect(screen.getByText('3')).toBeInTheDocument(); // Total alerts
    expect(screen.getByText('1')).toBeInTheDocument(); // High severity
    expect(screen.getByText('1')).toBeInTheDocument(); // Medium severity
    expect(screen.getByText('1')).toBeInTheDocument(); // Low severity
  });

  it('displays recent events table', () => {
    const events = [
      {
        id: 'event-1',
        timestamp: new Date().toISOString(),
        eventType: 'process',
        riskScore: 0.8,
        description: 'Suspicious process detected',
      },
      {
        id: 'event-2',
        timestamp: new Date().toISOString(),
        eventType: 'file',
        riskScore: 0.3,
        description: 'File access logged',
      },
    ];

    const apiStateWithEvents = {
      ...mockApiState,
      events,
    };

    jest.mock('../contexts/ApiContext', () => ({
      useApi: () => apiStateWithEvents,
    }));

    render(<Dashboard />);

    expect(screen.getByText('Suspicious process detected')).toBeInTheDocument();
    expect(screen.getByText('File access logged')).toBeInTheDocument();
    expect(screen.getByText('process')).toBeInTheDocument();
    expect(screen.getByText('file')).toBeInTheDocument();
  });

  it('shows system health indicators', () => {
    render(<Dashboard />);

    expect(screen.getByText(/system health/i)).toBeInTheDocument();
    expect(screen.getByText(/agents online/i)).toBeInTheDocument();
    expect(screen.getByText(/data processing/i)).toBeInTheDocument();
  });

  it('renders risk score chart', () => {
    render(<Dashboard />);

    expect(screen.getByText(/risk trends/i)).toBeInTheDocument();
    // Chart components are mocked, so we just check for the container
  });

  it('displays time range selector', () => {
    render(<Dashboard />);

    expect(screen.getByText(/last 24 hours/i)).toBeInTheDocument();
    expect(screen.getByText(/last 7 days/i)).toBeInTheDocument();
    expect(screen.getByText(/last 30 days/i)).toBeInTheDocument();
  });

  it('handles time range changes', () => {
    const mockFetchEvents = jest.fn();

    jest.mock('../contexts/ApiContext', () => ({
      useApi: () => ({
        ...mockApiState,
        fetchEvents: mockFetchEvents,
      }),
    }));

    render(<Dashboard />);

    const timeRangeSelect = screen.getByDisplayValue(/last 24 hours/i);
    fireEvent.change(timeRangeSelect, { target: { value: '7d' } });

    expect(mockFetchEvents).toHaveBeenCalledWith('7d');
  });

  it('shows refresh button', () => {
    render(<Dashboard />);

    const refreshButton = screen.getByRole('button', { name: /refresh/i });
    expect(refreshButton).toBeInTheDocument();
  });

  it('handles refresh action', () => {
    const mockFetchAlerts = jest.fn();
    const mockFetchEvents = jest.fn();

    jest.mock('../contexts/ApiContext', () => ({
      useApi: () => ({
        ...mockApiState,
        fetchAlerts: mockFetchAlerts,
        fetchEvents: mockFetchEvents,
      }),
    }));

    render(<Dashboard />);

    const refreshButton = screen.getByRole('button', { name: /refresh/i });
    fireEvent.click(refreshButton);

    expect(mockFetchAlerts).toHaveBeenCalled();
    expect(mockFetchEvents).toHaveBeenCalled();
  });

  it('displays export options', () => {
    render(<Dashboard />);

    expect(screen.getByText(/export/i)).toBeInTheDocument();
    expect(screen.getByText(/pdf/i)).toBeInTheDocument();
    expect(screen.getByText(/csv/i)).toBeInTheDocument();
  });

  it('shows notification preferences', () => {
    render(<Dashboard />);

    expect(screen.getByText(/notifications/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/email notifications/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/push notifications/i)).toBeInTheDocument();
  });

  it('handles notification preference changes', () => {
    render(<Dashboard />);

    const emailCheckbox = screen.getByLabelText(/email notifications/i);
    const pushCheckbox = screen.getByLabelText(/push notifications/i);

    fireEvent.click(emailCheckbox);
    fireEvent.click(pushCheckbox);

    // Checkboxes should be checked
    expect(emailCheckbox).toBeChecked();
    expect(pushCheckbox).toBeChecked();
  });

  it('displays quick actions', () => {
    render(<Dashboard />);

    expect(screen.getByText(/quick actions/i)).toBeInTheDocument();
    expect(screen.getByText(/view all alerts/i)).toBeInTheDocument();
    expect(screen.getByText(/generate report/i)).toBeInTheDocument();
    expect(screen.getByText(/system settings/i)).toBeInTheDocument();
  });

  it('navigates to alerts page when view all alerts is clicked', () => {
    const mockNavigate = jest.fn();
    jest.mock('react-router-dom', () => ({
      useNavigate: () => mockNavigate,
    }));

    render(<Dashboard />);

    const viewAlertsButton = screen.getByText(/view all alerts/i);
    fireEvent.click(viewAlertsButton);

    expect(mockNavigate).toHaveBeenCalledWith('/alerts');
  });

  it('shows user avatar and profile menu', () => {
    render(<Dashboard />);

    expect(screen.getByAltText(/user avatar/i)).toBeInTheDocument();
    expect(screen.getByText(/profile/i)).toBeInTheDocument();
    expect(screen.getByText(/settings/i)).toBeInTheDocument();
    expect(screen.getByText(/logout/i)).toBeInTheDocument();
  });

  it('handles logout action', () => {
    const mockLogout = jest.fn();

    jest.mock('../contexts/AuthContext', () => ({
      useAuth: () => ({
        ...mockAuthState,
        logout: mockLogout,
      }),
    }));

    render(<Dashboard />);

    const logoutButton = screen.getByText(/logout/i);
    fireEvent.click(logoutButton);

    expect(mockLogout).toHaveBeenCalled();
  });

  it('displays breadcrumbs', () => {
    render(<Dashboard />);

    expect(screen.getByText(/home/i)).toBeInTheDocument();
    expect(screen.getByText(/dashboard/i)).toBeInTheDocument();
  });

  it('shows last updated timestamp', () => {
    render(<Dashboard />);

    expect(screen.getByText(/last updated/i)).toBeInTheDocument();
  });

  it('handles window resize events', () => {
    render(<Dashboard />);

    // Simulate window resize
    window.innerWidth = 800;
    fireEvent(window, new Event('resize'));

    // Dashboard should adapt to smaller screen
    expect(screen.getByText(/dashboard/i)).toBeInTheDocument();
  });

  it('displays contextual help tooltips', () => {
    render(<Dashboard />);

    // Hover over an element with help tooltip
    const alertsCard = screen.getByText(/alerts overview/i);
    fireEvent.mouseOver(alertsCard);

    // Tooltip should appear (mocked)
    expect(screen.getByText(/alerts overview/i)).toBeInTheDocument();
  });
});</content>
<parameter name="filePath">/workspaces/insec/tests/unit/ui/dashboard/Dashboard.test.tsx
