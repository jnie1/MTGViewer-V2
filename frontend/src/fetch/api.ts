import ResponseError from './ResponseError';

async function fetchApi<T = unknown>(path: string, init?: RequestInit): Promise<T> {
  const basePath = import.meta.env.VITE_API_URL;
  const fullPath = new URL(path, basePath);

  const headers: HeadersInit = { ...init?.headers, ['Accept']: 'application/json' };
  const fullInit: RequestInit = { ...init, headers, credentials: 'include' };

  const response = await fetch(fullPath, fullInit);
  if (!response.ok) {
    throw new ResponseError(response);
  }

  const contentType = response.headers.get('content-type') ?? '';
  if (!contentType.includes('application/json')) {
    throw new ResponseError(response, 'content does not contain json');
  }

  if (response.status === 204) {
    return undefined as T;
  }

  const body: T = await response.json();
  return body;
}

export default fetchApi;
