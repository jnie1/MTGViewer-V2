import ResponseError from './ResponeError';

const basePath = new URL(import.meta.env.VITE_API_URL);

async function fetchApi<T = unknown>(path: string, init?: RequestInit): Promise<T> {
  const fullPath = new URL(path, basePath);
  const fullInit: RequestInit = {
    ...init,
    credentials: 'include',
  };

  const repsonse = await fetch(fullPath, fullInit);

  if (!repsonse.ok) {
    throw new ResponseError(repsonse);
  }

  const body: T = await repsonse.json();

  return body;
}

export default fetchApi;
