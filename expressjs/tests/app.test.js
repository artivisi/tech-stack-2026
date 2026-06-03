import request from 'supertest';
import { pool } from '../src/db.js';
import app from '../src/app.js'; 

describe('GET /health', () => {
  it('should respond with a 200 status code', async () => {
    const response = await request(app)
      .get('/health')
      .set('Accept', 'application/json');

    expect(response.statusCode).toBe(200);
    expect(response.body).toEqual({
      status: 'ok'
    });
  });

  afterAll(async () => {
    await pool.end();
  });
});