import React, { useState, useEffect } from 'react';
import {
  ThemeProvider,
  createTheme,
  CssBaseline,
  Box,
  AppBar,
  Toolbar,
  Typography,
  Container,
  Grid,
  Card,
  CardContent,
  Chip,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Button,
  Drawer,
  ListItemButton,
  Divider,
} from '@mui/material';
import {
  Security,
  Warning,
  Menu,
  Dashboard,
  Rule,
  Settings,
  Assessment,
} from '@mui/icons-material';

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#00bcd4',
    },
    secondary: {
      main: '#ff4081',
    },
    background: {
      default: '#121212',
      paper: '#1e1e1e',
    },
  },
  typography: {
    h4: {
      fontWeight: 600,
    },
    h6: {
      fontWeight: 500,
    },
  },
});

interface Alert {
  id: number;
  created: string;
  severity: string;
  title: string;
  status: string;
}

interface DashboardStats {
  totalAlerts: number;
  activeAlerts: number;
  endpoints: number;
  eventsToday: number;
  riskScore: number;
}

function App() {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [stats, setStats] = useState<DashboardStats>({
    totalAlerts: 0,
    activeAlerts: 0,
    endpoints: 0,
    eventsToday: 0,
    riskScore: 0,
  });
  const [alerts, setAlerts] = useState<Alert[]>([]);

  useEffect(() => {
    fetchDashboardData();
  }, []);

  const fetchDashboardData = async () => {
    try {
      // Mock data for demonstration
      setStats({
        totalAlerts: 47,
        activeAlerts: 12,
        endpoints: 156,
        eventsToday: 2847,
        riskScore: 68,
      });

      setAlerts([
        {
          id: 1,
          created: new Date().toISOString(),
          severity: 'high',
          title: 'Suspicious data exfiltration detected',
          status: 'open',
        },
        {
          id: 2,
          created: new Date(Date.now() - 3600000).toISOString(),
          severity: 'medium',
          title: 'Privilege escalation attempt',
          status: 'investigating',
        },
        {
          id: 3,
          created: new Date(Date.now() - 7200000).toISOString(),
          severity: 'low',
          title: 'Unusual login pattern',
          status: 'resolved',
        },
      ]);
    } catch (error) {
      console.error('Failed to fetch dashboard data:', error);
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'high': return 'error';
      case 'medium': return 'warning';
      case 'low': return 'info';
      default: return 'default';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'open': return '#ff4081';
      case 'investigating': return '#ff9800';
      case 'resolved': return '#4caf50';
      default: return '#9e9e9e';
    }
  };

  const menuItems = [
    { text: 'Dashboard', icon: <Dashboard />, active: true },
    { text: 'Alerts', icon: <Warning /> },
    { text: 'Rules', icon: <Rule /> },
    { text: 'Reports', icon: <Assessment /> },
    { text: 'Settings', icon: <Settings /> },
  ];

  return (
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <Box sx={{ display: 'flex' }}>
        {/* App Bar */}
        <AppBar position="fixed" sx={{ zIndex: (theme: any) => theme.zIndex.drawer + 1 }}>
          <Toolbar>
            <Button
              color="inherit"
              onClick={() => setDrawerOpen(!drawerOpen)}
              sx={{ mr: 2 }}
            >
              <Menu />
            </Button>
            <Security sx={{ mr: 1 }} />
            <Typography variant="h6" noWrap component="div" sx={{ flexGrow: 1 }}>
              INSEC Console
            </Typography>
            <Typography variant="body1">
              Enterprise Insider-Threat Protection
            </Typography>
          </Toolbar>
        </AppBar>

        {/* Navigation Drawer */}
        <Drawer
          variant="temporary"
          open={drawerOpen}
          onClose={() => setDrawerOpen(false)}
          sx={{
            width: 240,
            flexShrink: 0,
            '& .MuiDrawer-paper': {
              width: 240,
              boxSizing: 'border-box',
            },
          }}
        >
          <Toolbar />
          <Box sx={{ overflow: 'auto' }}>
            <List>
              {menuItems.map((item) => (
                <ListItem key={item.text} disablePadding>
                  <ListItemButton selected={item.active}>
                    <ListItemIcon>
                      {item.icon}
                    </ListItemIcon>
                    <ListItemText primary={item.text} />
                  </ListItemButton>
                </ListItem>
              ))}
            </List>
          </Box>
        </Drawer>

        {/* Main Content */}
        <Box component="main" sx={{ flexGrow: 1, p: 3 }}>
          <Toolbar />

          <Container maxWidth="xl">
            <Typography variant="h4" gutterBottom sx={{ mb: 4 }}>
              Security Dashboard
            </Typography>

            {/* Stats Cards */}
            <Grid container spacing={3} sx={{ mb: 4 }}>
              <Grid item xs={12} sm={6} md={3}>
                <Card>
                  <CardContent>
                    <Typography color="textSecondary" gutterBottom>
                      Risk Score
                    </Typography>
                    <Typography variant="h4" component="div" color="primary">
                      {stats.riskScore}%
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      Organization risk level
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>

              <Grid item xs={12} sm={6} md={3}>
                <Card>
                  <CardContent>
                    <Typography color="textSecondary" gutterBottom>
                      Active Alerts
                    </Typography>
                    <Typography variant="h4" component="div" color="error">
                      {stats.activeAlerts}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      Require attention
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>

              <Grid item xs={12} sm={6} md={3}>
                <Card>
                  <CardContent>
                    <Typography color="textSecondary" gutterBottom>
                      Endpoints
                    </Typography>
                    <Typography variant="h4" component="div">
                      {stats.endpoints}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      Protected devices
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>

              <Grid item xs={12} sm={6} md={3}>
                <Card>
                  <CardContent>
                    <Typography color="textSecondary" gutterBottom>
                      Events Today
                    </Typography>
                    <Typography variant="h4" component="div">
                      {stats.eventsToday.toLocaleString()}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      Telemetry events
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>

            {/* Activity Summary */}
            <Grid container spacing={3} sx={{ mb: 4 }}>
              <Grid item xs={12} md={6}>
                <Card>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Activity Summary
                    </Typography>
                    <Typography variant="body1" paragraph>
                      System is monitoring {stats.endpoints} endpoints with {stats.eventsToday.toLocaleString()} events processed today.
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      Last updated: {new Date().toLocaleTimeString()}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>

              <Grid item xs={12} md={6}>
                <Card>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      System Status
                    </Typography>
                    <Typography variant="body1" paragraph>
                      All systems operational. Risk score indicates moderate threat level.
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      Server: Online | Agent: Connected
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>

            {/* Recent Alerts */}
            <Grid container spacing={3}>
              <Grid item xs={12}>
                <Card>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      Recent Alerts
                    </Typography>
                    <List>
                      {alerts.map((alert) => (
                        <React.Fragment key={alert.id}>
                          <ListItem
                            secondaryAction={
                              <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
                                <Chip
                                  label={alert.severity.toUpperCase()}
                                  color={getSeverityColor(alert.severity) as any}
                                  size="small"
                                />
                                <Chip
                                  label={alert.status}
                                  sx={{
                                    backgroundColor: getStatusColor(alert.status),
                                    color: 'white',
                                  }}
                                  size="small"
                                />
                              </Box>
                            }
                          >
                            <ListItemIcon>
                              <Warning color="error" />
                            </ListItemIcon>
                            <ListItemText
                              primary={alert.title}
                              secondary={`Created: ${new Date(alert.created).toLocaleString()}`}
                            />
                          </ListItem>
                          <Divider />
                        </React.Fragment>
                      ))}
                    </List>
                    <Box sx={{ mt: 2, textAlign: 'center' }}>
                      <Button variant="outlined" color="primary">
                        View All Alerts
                      </Button>
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          </Container>
        </Box>
      </Box>
    </ThemeProvider>
  );
}

export default App;