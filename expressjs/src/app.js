import express from 'express';
import { engine } from 'express-handlebars';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import registrationRoutes from './routes/registration.js';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

function formatDate(value) {
  if (!value) return '';
  const d = value instanceof Date ? value : new Date(value);
  return d.toISOString().replace('T', ' ').slice(0, 19) + ' UTC';
}

function requestLogger(req, res, next) {
  const start = Date.now();
  res.on('finish', () => {
    console.log(
      `[${new Date().toISOString()}] ${req.method} ${req.originalUrl} ${res.statusCode} ${Date.now() - start}ms`,
    );
  });
  next();
}

export function createApp() {
  const app = express();

  app.engine(
    'hbs',
    engine({
      extname: '.hbs',
      defaultLayout: 'main',
      helpers: { formatDate },
    }),
  );
  app.set('view engine', 'hbs');
  app.set('views', path.join(__dirname, 'views'));

  app.locals.nodeVersion = process.versions.node;

  app.use(requestLogger);
  app.use(express.urlencoded({ extended: false }));
  app.use(express.static(path.join(__dirname, '..', 'public')));

  app.use('/', registrationRoutes);

  app.use((err, req, res, next) => {
    console.error(err);
    res.status(500).render('error', { message: err.message });
  });

  return app;
}
