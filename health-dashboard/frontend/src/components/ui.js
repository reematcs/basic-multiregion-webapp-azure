import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { AlertCircle, CheckCircle, RefreshCw, ArrowRightLeft } from 'lucide-react';
import { Alert, AlertDescription, AlertTitle } from './ui/alert';
import { Button } from './ui/button';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

const Dashboard = () => {
  const [systemInfo, setSystemInfo] = useState(null);
  const [healthStatus, setHealthStatus] = useState(null);
  const [metrics, setMetrics] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const baseUrl = process.env.REACT_APP_API_URL || '';
        
        // Fetch critical data first
        const [systemRes, healthRes] = await Promise.all([
          fetch(`${baseUrl}/api/system`),
          fetch(`${baseUrl}/api/health/status`),
        ]);
    
        // Debug responses
        console.log('System Response:', {
          status: systemRes.status,
          statusText: systemRes.statusText
        });
        console.log('Health Response:', {
          status: healthRes.status,
          statusText: healthRes.statusText
        });
    
        // Check if responses are ok before trying to parse JSON
        if (!systemRes.ok || !healthRes.ok) {
          throw new Error('Failed to fetch critical data');
        }
    
        const [systemInfo, healthStatus] = await Promise.all([
          systemRes.json(),
          healthRes.json(),
        ]);
    
        setSystemInfo(systemInfo);
        setHealthStatus(healthStatus);
    
        // Fetch metrics separately - don't let it block the main data
        try {
          const metricsRes = await fetch(`${baseUrl}/api/metrics`);
          if (metricsRes.ok) {
            const metrics = await metricsRes.json();
            setMetrics(metrics);
          } else {
            console.log('Metrics not available');
            setMetrics([]);  // Set empty metrics
          }
        } catch (metricsError) {
          console.log('Could not fetch metrics:', metricsError);
          setMetrics([]);  // Set empty metrics
        }
    
      } catch (error) {
        console.error('Failed to fetch data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, 30000);
    return () => clearInterval(interval);
  }, []);

  const handleFailover = async () => {
    try {
      await fetch('/api/failover/trigger', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ targetRegion: 'Central US' })
      });
    } catch (error) {
      console.error('Failed to trigger failover:', error);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <RefreshCw className="w-8 h-8 animate-spin text-blue-500" />
      </div>
    );
  }

  return (
    <div className="p-6 space-y-6 max-w-7xl mx-auto">
      {/* System Info */}
      <Card>
        <CardHeader>
          <CardTitle>System Information</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <p className="text-sm font-medium text-gray-500">Region</p>
              <p className="text-2xl">{systemInfo?.region}</p>
            </div>
            <div>
              <p className="text-sm font-medium text-gray-500">Container Version</p>
              <p className="text-2xl">{systemInfo?.containerVersion}</p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Health Status */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Service Health</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {['TrafficManager', 'KeyVault', 'ContainerRegistry'].map((service) => (
                <div key={service} className="flex items-center justify-between">
                  <span className="text-gray-600">{service}</span>
                  {healthStatus?.[service.toLowerCase()] === 'connected' ? (
                    <CheckCircle className="w-5 h-5 text-green-500" />
                  ) : (
                    <AlertCircle className="w-5 h-5 text-red-500" />
                  )}
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Region Status</CardTitle>
          </CardHeader>
          <CardContent>
            <Alert className={healthStatus?.regionStatus === 'healthy' ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'}>
              <AlertTitle>
                {healthStatus?.role === 'primary' ? 'Primary Region' : 'Secondary Region'}
              </AlertTitle>
              <AlertDescription>
                Status: {healthStatus?.regionStatus}
              </AlertDescription>
            </Alert>
            <Button
              onClick={handleFailover}
              className="mt-4 w-full bg-blue-500 hover:bg-blue-600 text-white transition-colors"
            >
              <ArrowRightLeft className="w-4 h-4 mr-2" />
              Trigger Failover
            </Button>
          </CardContent>
        </Card>
      </div>

      {/* Metrics Chart */}
      <Card>
        <CardHeader>
          <CardTitle>Performance Metrics</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="h-[300px] w-full">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={metrics}>
                <CartesianGrid strokeDasharray="3 3" className="opacity-50" />
                <XAxis
                  dataKey="timestamp"
                  className="text-gray-600"
                />
                <YAxis className="text-gray-600" />
                <Tooltip
                  contentStyle={{
                    backgroundColor: 'white',
                    border: '1px solid #e2e8f0',
                    borderRadius: '6px',
                    padding: '8px'
                  }}
                />
                <Line
                  type="monotone"
                  dataKey="latency"
                  stroke="#2563eb"
                  strokeWidth={2}
                  dot={false}
                />
                <Line
                  type="monotone"
                  dataKey="requestCount"
                  stroke="#16a34a"
                  strokeWidth={2}
                  dot={false}
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default Dashboard;