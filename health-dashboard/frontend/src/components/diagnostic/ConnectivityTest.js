import React, { useState, useEffect, useCallback } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '../common/card';
import { Alert, AlertDescription } from '../common/alert';

const ConnectivityTest = ({ onComplete }) => {
    const [endpoints, setEndpoints] = useState([
        { name: 'System Info', url: '/api/system', status: 'checking' },
        { name: 'Health Status', url: '/api/health/status', status: 'checking' },
        { name: 'Metrics', url: '/api/metrics', status: 'checking' },
        { name: 'Failover History', url: '/api/failover/history', status: 'checking' }
    ]);

    const [overallStatus, setOverallStatus] = useState('checking');

    const checkEndpoint = useCallback(async (endpoint) => {
        try {
            console.log(`⏳ Starting check for ${endpoint.name}`);

            // Add a delay to make sure we see each request clearly
            await new Promise(resolve => setTimeout(resolve, 500));

            const response = await fetch(endpoint.url, {
                method: 'GET',
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json'
                }
            });

            console.log(`📥 ${endpoint.name} Response:`, {
                status: response.status,
                ok: response.ok,
                statusText: response.statusText,
                type: response.type,
                headers: Object.fromEntries(response.headers)
            });

            let responseData;
            if (response.ok) {
                responseData = await response.json();
                console.log(`✅ ${endpoint.name} Data:`, responseData);
                return 'connected';
            } else {
                console.warn(`⚠️ ${endpoint.name} Failed:`, response.statusText);
                return 'error';
            }
        } catch (err) {
            console.error(`❌ ${endpoint.name} Error:`, {
                name: err.name,
                message: err.message,
                cause: err.cause,
                stack: err.stack
            });
            return 'error';
        }
    }, []);

    const runTests = useCallback(async () => {
        console.log('🚀 Starting connectivity tests...');

        // Check endpoints sequentially instead of all at once
        const results = [];
        for (const endpoint of endpoints) {
            const status = await checkEndpoint(endpoint);
            results.push({ ...endpoint, status });
            // Wait a bit between checks
            await new Promise(resolve => setTimeout(resolve, 1000));
        }

        setEndpoints(results);

        const hasErrors = results.some(endpoint => endpoint.status === 'error');
        const newOverallStatus = hasErrors ? 'error' : 'connected';
        console.log('🏁 Test results:', {
            results: results.map(r => ({ name: r.name, status: r.status })),
            overallStatus: newOverallStatus
        });

        setOverallStatus(newOverallStatus);
        onComplete?.(newOverallStatus === 'connected');
    }, [endpoints, checkEndpoint, onComplete]);

    useEffect(() => {
        runTests();
    }, [runTests]);

    const getStatusColor = (status) => {
        switch (status) {
            case 'connected':
                return 'bg-green-100 text-green-800 border-green-200';
            case 'checking':
                return 'bg-yellow-100 text-yellow-800 border-yellow-200';
            case 'error':
                return 'bg-red-100 text-red-800 border-red-200';
            default:
                return 'bg-gray-100 text-gray-800 border-gray-200';
        }
    };

    return (
        <Card>
            <CardHeader>
                <CardTitle>API Connectivity Test</CardTitle>
            </CardHeader>
            <CardContent>
                <Alert className={getStatusColor(overallStatus)}>
                    <AlertDescription>
                        Overall Status: {overallStatus.toUpperCase()}
                    </AlertDescription>
                </Alert>

                <div className="mt-4 space-y-2">
                    {endpoints.map((endpoint) => (
                        <div
                            key={endpoint.url}
                            className="flex items-center justify-between p-2 border rounded"
                        >
                            <span className="font-medium">{endpoint.name}</span>
                            <div className="flex items-center">
                                <span className="mr-2 text-sm text-gray-600">{endpoint.url}</span>
                                <span
                                    className={`px-2 py-1 rounded text-sm ${getStatusColor(endpoint.status)}`}
                                >
                                    {endpoint.status}
                                </span>
                            </div>
                        </div>
                    ))}
                </div>
            </CardContent>
        </Card>
    );
};

export default ConnectivityTest;