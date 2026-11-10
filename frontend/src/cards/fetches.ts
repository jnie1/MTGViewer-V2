import fetchApi from '@/fetch/api';
import fetchBlob from '@/fetch/blob';
import { type ISearchResult } from './types';

export function searchCards(search: string, page: number, abort?: AbortSignal) {
  const params = new URLSearchParams({
    q: search,
    page: page.toString(),
  });
  return fetchApi<ISearchResult>(`/cards/search?${params}`, { signal: abort });
}

export function prePrune(file: File, size: number, price: number, signal?: AbortSignal) {
  const query = new URLSearchParams({ size: String(size), price: String(price) });
  const formData = new FormData();
  formData.append('file', file);
  return fetchBlob(`/cards/preprune?${query}`, { method: 'POST', body: formData, signal });
}
