import { z } from 'zod';

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
const NAME_RE = /^[\p{L}\p{M}\s.'\-]+$/u;
const PHONE_RE = /^[+0-9 ()\-]+$/;

export const RegistrationSchema = z.object({
  email: z
    .string()
    .trim()
    .toLowerCase()
    .min(3, 'valid email is required')
    .max(254, 'email is too long')
    .regex(EMAIL_RE, 'valid email is required'),
  fullName: z
    .string()
    .trim()
    .min(2, 'full name must be 2-100 characters')
    .max(100, 'full name must be 2-100 characters')
    .regex(NAME_RE, 'full name contains invalid characters'),
  phone: z
    .string()
    .trim()
    .min(7, 'phone must be 7-20 characters')
    .max(20, 'phone must be 7-20 characters')
    .regex(PHONE_RE, 'phone contains invalid characters'),
});

export function collectFieldErrors(zodError) {
  const errors = {};
  for (const issue of zodError.issues) {
    const key = issue.path[0];
    if (!errors[key]) errors[key] = issue.message;
  }
  return errors;
}
