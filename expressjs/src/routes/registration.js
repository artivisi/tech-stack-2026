import { Router } from 'express';
import { randomUUID } from 'node:crypto';
import { pool } from '../db.js';
import { RegistrationSchema, collectFieldErrors } from '../validation.js';

const router = Router();

router.get('/health', async (req, res) => {
  try {
    await pool.query('SELECT 1');
    res.json({ status: 'ok' });
  } catch (err) {
    console.error('health check failed:', err);
    res.status(503).json({
      status: 'error',
      code: err.code,
      error: String(err),
    });
  }
});

router.get('/', (req, res) => {
  res.render('form', { values: {}, errors: {} });
});

router.post('/register', async (req, res, next) => {
  const submitted = {
    email: req.body.email ?? '',
    full_name: req.body.full_name ?? '',
    phone: req.body.phone ?? '',
  };

  const result = RegistrationSchema.safeParse(req.body);

  if (!result.success) {
    return res.status(400).render('form', {
      values: submitted,
      errors: collectFieldErrors(result.error),
    });
  }

  const { email, full_name, phone } = result.data;

  try {
    await pool.query(
      'INSERT INTO registration (id, email, full_name, phone, created_at) VALUES ($1, $2, $3, $4, $5)',
      [randomUUID(), email, full_name, phone, new Date()],
    );
    res.redirect('/registrations');
  } catch (err) {
    if (err.code === '23505') {
      return res.status(409).render('form', {
        values: submitted,
        errors: { email: 'email is already registered' },
      });
    }
    next(err);
  }
});

router.get('/registrations', async (req, res, next) => {
  try {
    const { rows } = await pool.query(
      'SELECT id, email, full_name, phone, created_at FROM registration ORDER BY created_at DESC',
    );
    res.render('list', { registrations: rows, count: rows.length });
  } catch (err) {
    next(err);
  }
});

export default router;
