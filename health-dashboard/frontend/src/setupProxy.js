const { createProxyMiddleware } = require('http-proxy-middleware');

console.log('Loading proxy configuration...');

module.exports = function (app) {
    console.log('Configuring proxy for /api...');

    app.use('/api', createProxyMiddleware({
        target: 'http://west-us:8080',
        changeOrigin: true,
        secure: false,
        logLevel: 'debug',
        onProxyReq: (proxyReq, req, res) => {
            console.log('🔄 Proxying request:', {
                method: req.method,
                path: req.path,
                headers: proxyReq.getHeaders()
            });
        },
        onProxyRes: (proxyRes, req, res) => {
            console.log('✅ Proxy response:', {
                status: proxyRes.statusCode,
                path: req.path
            });
        },
        onError: (err, req, res) => {
            console.error('❌ Proxy error:', err.message);
            res.status(500).json({
                error: 'Proxy Error',
                message: err.message,
                code: err.code
            });
        },
        router: {
            // Forward to west-us service
            '/api': 'http://west-us:8080',
        }
    }));

    console.log('Proxy configuration complete');
};
// const { createProxyMiddleware } = require('http-proxy-middleware');

// module.exports = function (app) {
//     app.use(
//         '/api',
//         createProxyMiddleware({
//             target: 'http://west-us:8080',
//             changeOrigin: true,
//             pathRewrite: {
//                 '^/api': '/api'
//             },
//             onError: (err, req, res) => {
//                 console.error('Proxy error:', err);
//                 res.status(500).send('Proxy Error');
//             }
//         })
//     );
// };