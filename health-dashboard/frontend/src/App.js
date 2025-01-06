import React, { useState } from 'react';
import ConnectivityTest from './components/diagnostic/ConnectivityTest';
import Dashboard from './components/dashboard/Dashboard';  // We'll move ui.js here

const App = () => {
    const [isConnected, setIsConnected] = useState(null);

    if (isConnected === null) {
        return (
            <div className="min-h-screen bg-background">
                <ConnectivityTest onComplete={setIsConnected} />
            </div>
        );
    }

    if (!isConnected) {
        return (
            <div className="min-h-screen bg-background p-6">
                <h2 className="text-xl font-bold text-destructive">Connection Error</h2>
                <p className="mt-2">Unable to connect to required services. Please check your configuration.</p>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-background">
            <Dashboard />
        </div>
    );
};

export default App;