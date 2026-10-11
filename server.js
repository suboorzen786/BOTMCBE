const express = require('express');
const path = require('path');

const app = express();
const PORT = 8080;

app.use(express.json());
app.use(express.static(path.join(__dirname, 'public')));

app.get('/api/status', (req, res) => {
  res.json({
    app: 'BOTMC PRO',
    type: 'MCPE / Bedrock',
    backend: 'not_connected',
    message: 'Bedrock protocol backend is not configured yet.'
  });
});

app.listen(PORT, '0.0.0.0', () => {
  console.log(`BOTMC PRO website: http://127.0.0.1:${PORT}`);
  console.log('MCPE backend: not configured yet');
});
