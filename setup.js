const http = require('http');
const fs = require('fs');
const path = require('path');
const { exec } = require('child_process');

const CONFIG_FILE = path.join(__dirname, '.openclaw', 'openclaw.json');
const PORT = 8080;

// Read config
function readConfig() {
    try {
        const content = fs.readFileSync(CONFIG_FILE, 'utf8');
        return JSON.parse(content);
    } catch (e) {
        return null;
    }
}

// Write config
function writeConfig(config) {
    fs.writeFileSync(CONFIG_FILE, JSON.stringify(config, null, 2), 'utf8');
}

// Check if API key is configured (not placeholder ***)
function hasValidApiKey(config) {
    if (!config || !config.models || !config.models.providers) return false;
    const providers = config.models.providers;
    for (const key in providers) {
        const apiKey = providers[key].apiKey;
        if (apiKey && apiKey !== '***' && apiKey.trim() !== '') {
            return true;
        }
    }
    return false;
}

// HTML page
const HTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>OpenClaw 激活</title>
<style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body { 
        font-family: 'Microsoft YaHei', Arial, sans-serif; 
        background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
        min-height: 100vh;
        display: flex;
        justify-content: center;
        align-items: center;
        padding: 20px;
    }
    .container {
        background: #0f3460;
        border-radius: 20px;
        padding: 40px;
        max-width: 500px;
        width: 100%;
        box-shadow: 0 20px 60px rgba(0,0,0,0.5);
    }
    h1 { 
        color: #e94560; 
        text-align: center; 
        margin-bottom: 30px;
        font-size: 28px;
    }
    .notice {
        background: #1a1a2e;
        border-radius: 12px;
        padding: 20px;
        margin-bottom: 25px;
        border-left: 4px solid #e94560;
    }
    .notice p {
        color: #ccc;
        font-size: 14px;
        line-height: 1.8;
        margin-bottom: 10px;
    }
    .notice a {
        color: #00d9ff;
        text-decoration: none;
        font-weight: bold;
    }
    .notice a:hover {
        text-decoration: underline;
    }
    .notice ol {
        color: #ccc;
        font-size: 14px;
        line-height: 2;
        padding-left: 20px;
    }
    .input-group {
        margin-bottom: 20px;
    }
    label {
        color: #fff;
        display: block;
        margin-bottom: 8px;
        font-size: 14px;
    }
    input[type="text"] {
        width: 100%;
        padding: 15px;
        border: 2px solid #1a1a2e;
        border-radius: 10px;
        font-size: 16px;
        background: #1a1a2e;
        color: #fff;
        outline: none;
        transition: border-color 0.3s;
    }
    input[type="text"]:focus {
        border-color: #e94560;
    }
    .btn {
        width: 100%;
        padding: 15px;
        background: #e94560;
        color: #fff;
        border: none;
        border-radius: 10px;
        font-size: 18px;
        font-weight: bold;
        cursor: pointer;
        transition: background 0.3s, transform 0.2s;
    }
    .btn:hover {
        background: #d1364f;
        transform: scale(1.02);
    }
    .btn:disabled {
        background: #666;
        cursor: not-allowed;
    }
    .error {
        background: #ff4757;
        color: #fff;
        padding: 15px;
        border-radius: 10px;
        margin-bottom: 20px;
        text-align: center;
        display: none;
    }
    .success {
        background: #2ed573;
        color: #fff;
        padding: 15px;
        border-radius: 10px;
        margin-bottom: 20px;
        text-align: center;
        display: none;
    }
    .spinner {
        display: none;
        text-align: center;
        color: #fff;
        margin-bottom: 20px;
    }
    .spinner::after {
        content: '';
        display: inline-block;
        width: 20px;
        height: 20px;
        border: 3px solid #fff;
        border-top-color: transparent;
        border-radius: 50%;
        animation: spin 1s linear infinite;
        margin-left: 10px;
        vertical-align: middle;
    }
    @keyframes spin {
        to { transform: rotate(360deg); }
    }
    .tip {
        color: #888;
        font-size: 12px;
        margin-top: 10px;
        text-align: center;
    }
</style>
</head>
<body>
<div class="container">
    <h1>OpenClaw U盘版</h1>
    
    <div class="notice">
        <p><strong>还没账号？先注册：</strong></p>
        <ol>
            <li>打开 <a href="http://3295b30e.r8.cpolar.cn/" target="_blank">http://3295b30e.r8.cpolar.cn/</a></li>
            <li>注册账号并登录</li>
            <li>充值获得Token额度</li>
            <li>在后台获取你的API Key</li>
        </ol>
        <p style="margin-top:15px;">💰 充值后即可获得额度，开始使用OpenClaw</p>
    </div>
    
    <div id="error" class="error"></div>
    <div id="success" class="success"></div>
    
    <div class="input-group">
        <label>请输入你的API Key：</label>
        <input type="text" id="apiKey" placeholder="输入从平台获取的API Key">
    </div>
    
    <button class="btn" id="validateBtn" onclick="validateKey()">验证并启动</button>
    
    <div class="spinner" id="spinner">验证中...</div>
    
    <p class="tip">API Key验证通过后会自动启动OpenClaw</p>
</div>

<script>
function showError(msg) {
    const el = document.getElementById('error');
    el.textContent = msg;
    el.style.display = 'block';
    document.getElementById('success').style.display = 'none';
}

function showSuccess(msg) {
    const el = document.getElementById('success');
    el.textContent = msg;
    el.style.display = 'block';
    document.getElementById('error').style.display = 'none';
}

function hideMsg() {
    document.getElementById('error').style.display = 'none';
    document.getElementById('success').style.display = 'none';
}

async function validateKey() {
    const apiKey = document.getElementById('apiKey').value.trim();
    if (!apiKey) {
        showError('请输入API Key');
        return;
    }
    
    hideMsg();
    document.getElementById('validateBtn').disabled = true;
    document.getElementById('spinner').style.display = 'block';
    
    try {
        // Try to call the API to validate key - test models endpoint
        const response = await fetch('http://3295b30e.r8.cpolar.cn/v1/models?key=' + encodeURIComponent(apiKey));
        
        if (response.ok) {
            // Key is valid, save to config
            showSuccess('验证成功！正在启动...');
            
            // Save to config via server endpoint
            const configRes = await fetch('/save-config', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ apiKey })
            });
            
            // Wait a moment then redirect to OpenClaw
            setTimeout(() => {
                window.location.href = '/start-openclaw';
            }, 1500);
        } else {
            showError('当前APIKey无法启动OpenClaw');
            document.getElementById('validateBtn').disabled = false;
        }
    } catch (e) {
        showError('当前APIKey无法启动OpenClaw');
        document.getElementById('validateBtn').disabled = false;
    }
    
    document.getElementById('spinner').style.display = 'none';
}

// Allow Enter key to submit
document.getElementById('apiKey').addEventListener('keypress', function(e) {
    if (e.key === 'Enter') validateKey();
});
</script>
</body>
</html>`;

// Save config endpoint
function saveConfig(apiKey, callback) {
    const config = readConfig() || {};
    if (!config.models) config.models = {};
    if (!config.models.providers) config.models.providers = {};
    if (!config.models.providers['token-platform']) {
        config.models.providers['token-platform'] = {
            baseUrl: 'http://3295b30e.r8.cpolar.cn/v1',
            apiKey: apiKey,
            api: 'openai-chat',
            models: [{ id: 'deepseek-ai/DeepSeek-V4-Flash', name: 'DeepSeek V4 Flash' }]
        };
    } else {
        config.models.providers['token-platform'].apiKey = apiKey;
    }
    
    try {
        writeConfig(config);
        callback(true);
    } catch (e) {
        callback(false);
    }
}

// Start OpenClaw
function startOpenClaw() {
    const exePath = path.join(__dirname, 'node', 'node.exe');
    const openclawPath = path.join(__dirname, 'openclaw', 'openclaw.mjs');
    
    exec(`"${exePath}" "${openclawPath}" gateway start`, {
        cwd: __dirname
    }, (err, stdout, stderr) => {
        if (err) {
            console.error('Failed to start:', err);
        }
    });
}

// HTTP server
const server = http.createServer((req, res) => {
    const url = req.url;
    
    if (url === '/' || url === '/setup') {
        res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8' });
        res.end(HTML);
    } else if (url === '/save-config') {
        let body = '';
        req.on('data', d => body += d);
        req.on('end', () => {
            try {
                const { apiKey } = JSON.parse(body);
                saveConfig(apiKey, (ok) => {
                    res.writeHead(200, { 'Content-Type': 'application/json' });
                    res.end(JSON.stringify({ ok }));
                });
            } catch (e) {
                res.writeHead(400, { 'Content-Type': 'application/json' });
                res.end(JSON.stringify({ ok: false }));
            }
        });
    } else if (url === '/start-openclaw') {
        res.writeHead(200, { 'Content-Type': 'text/html; charset=utf-8' });
        res.end(`<!DOCTYPE html><html><head><meta charset="UTF-8"><title>启动中</title></head><body style="font-family:Microsoft YaHei;background:#1a1a2e;color:#fff;display:flex;justify-content:center;align-items:center;min-height:100vh;"><div style="text-align:center;"><h1>OpenClaw 启动中...</h1><p>请稍候</p></div></body></html><script>setTimeout(() => window.close(), 3000);<\/script>`);
        startOpenClaw();
    } else {
        res.writeHead(404);
        res.end('Not found');
    }
});

server.listen(PORT, () => {
    console.log('Setup server running at http://localhost:' + PORT);
    
    // Check if already has valid key
    const config = readConfig();
    if (hasValidApiKey(config)) {
        console.log('Valid API key found, starting OpenClaw...');
        startOpenClaw();
        setTimeout(() => process.exit(0), 1000);
    } else {
        console.log('No valid API key, showing setup page...');
        // Open browser
        exec(`start http://localhost:${PORT}`);
    }
});