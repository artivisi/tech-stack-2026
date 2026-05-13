import pg from 'pg';
import { required } from './env.js';

export const pool = new pg.Pool({
  connectionString: required('DATABASE_URL'),
});
