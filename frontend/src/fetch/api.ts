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

async function request(path: string, init?: RequestInit): Promise<Response> {
  const fullPath = new URL(path, import.meta.env.VITE_API_URL);
  fullPath.pathname = '/api' + fullPath.pathname;

  const headers: HeadersInit = { ...init?.headers, ['Accept']: 'application/json' };
  const fullInit: RequestInit = { ...init, headers, credentials: 'omit' };

  const response = await fetch(fullPath, fullInit);
  if (!response.ok) {
    throw new ResponseError(response);
  }
  return response;
}
export async function fetchApiBlob(path: string, init?: RequestInit): Promise<Blob> {
  const response = await request(path, init);
  return response.blob();
}
export default fetchApi;
