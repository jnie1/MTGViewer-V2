import { accessToken } from './auth';
import ResponseError from './ResponseError';

async function fetchBlob(path: string, init?: RequestInit): Promise<Blob> {
  const fullPath = new URL(path, import.meta.env.VITE_API_URL);
  fullPath.pathname = '/api' + fullPath.pathname;

  const headers = new Headers(init?.headers);

  if (init?.credentials !== 'omit') {
    const token = accessToken.value;
    if (token) {
      headers.set('Authorization', `Bearer ${token}`);
    }
  }

  const fullInit: RequestInit = { ...init, headers, credentials: 'omit' };
  const response = await fetch(fullPath, fullInit);
  if (!response.ok) {
    throw new ResponseError(response);
  }
  return response.blob();
}

export default fetchBlob;
