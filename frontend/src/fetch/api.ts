import ResponseError from './ResponseError';

async function fetchApi<T = unknown>(path: string, init?: RequestInit): Promise<T> {
  const fullPath = new URL(path, import.meta.env.VITE_API_URL);
  fullPath.pathname = '/api' + fullPath.pathname;

  const headers: HeadersInit = { ...init?.headers, ['Accept']: 'application/json' };
  const fullInit: RequestInit = { ...init, headers, credentials: 'include' };

  const response = await fetch(fullPath, fullInit);
  if (!response.ok) {
    throw new ResponseError(response);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  const contentType = response.headers.get('content-type') ?? '';
  if (!contentType.includes('application/json')) {
    throw new ResponseError(response, 'content does not contain json');
  }

  const body: T = await response.json();
  return body;
}

export default fetchApi;
