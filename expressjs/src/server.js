import { createApp } from './app.js';
import { required } from './env.js';
import { pool } from './db.js';
import app from './app.js';

const port = Number(required('HTTP_PORT'));

const server = app.listen(port, () => {
  console.log(`registration-expressjs listening on http://localhost:${port}`);
});

function shutdown(signal) {
  console.log(`${signal} received — shutting down`);
  server.close(async (err) => {
    if (err) {
      console.error('http server close error:', err);
      process.exit(1);
    }
    try {
      await pool.end();
      process.exit(0);
    } catch (poolErr) {
      console.error('pool end error:', poolErr);
      process.exit(1);
    }
  });
}

process.on('SIGINT', () => shutdown('SIGINT'));
process.on('SIGTERM', () => shutdown('SIGTERM'));
