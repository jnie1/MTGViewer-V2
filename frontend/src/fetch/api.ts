import ResponseError from './ResponeError';

const basePath = new URL(import.meta.env.VITE_API_URL);

async function fetchApi<T = undefined>(path: string, init?: RequestInit): Promise<T> {
  const fullPath = new URL(path, basePath);

  const headers: HeadersInit = { ...init?.headers, ['Accept']: 'application/json' };
  const fullInit: RequestInit = { ...init, headers, credentials: 'include' };

  const response = await fetch(fullPath, fullInit);
  if (!response.ok) {
    throw new ResponseError(response);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  const body: T = await response.json();
  return body;
}

export default fetchApi;
