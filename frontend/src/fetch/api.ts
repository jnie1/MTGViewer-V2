import { accessToken } from './auth';
import ResponseError from './ResponseError';

async function fetchApi<T = unknown>(path: string, init?: RequestInit): Promise<T> {
  const fullPath = new URL(path, import.meta.env.VITE_API_URL);
  fullPath.pathname = '/api' + fullPath.pathname;

  const headers = new Headers(init?.headers);
  headers.set('Accept', 'application/json');

  if (init?.credentials !== 'omit') {
    const token = accessToken.value;
    if (token) {
      headers.set('Authorization', `Bearer ${token}`);
    }
  }

  const response = await fetch(fullPath, {
    ...init,
    headers,
    credentials: 'omit',
  });

  if (!response.ok) {
    throw new ResponseError(response);
  }

  const contentType = response.headers.get('content-type') ?? '';
  if (!contentType.includes('application/json')) {
    return undefined as T;
  }

  const body: T = await response.json();
  return body;
}

export function jsonRequest(value: unknown): RequestInit {
  return {
    body: JSON.stringify(value),
    headers: [['Content-Type', 'application/json']],
  };
}

export default fetchApi;
