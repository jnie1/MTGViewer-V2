import ResponseError from './ResponeError';

const basePath = new URL(import.meta.env.VITE_API_URL);

async function fetchApi<T = unknown>(path: string, init?: RequestInit): Promise<T> {
    const fullPath = new URL(path, basePath);
    const fullInit: RequestInit = {
        ...init,
        credentials: 'include',
    };

    const response = await fetch(fullPath, fullInit);
    if (!response.ok) {
        throw new ResponseError(response);
    }

    const body: T = await response.json();
    return body;
}

export default fetchApi;
