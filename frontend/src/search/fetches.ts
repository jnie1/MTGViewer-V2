import fetchApi from '@/fetch/api';
import { type ISearchResult } from './types';

export function searchCards(search: string, page: number, abort: AbortSignal) {
  const params = new URLSearchParams({
    q: search,
    page: page.toString(),
  });
  return fetchApi<ISearchResult>(`/cards/search?${params}`, { signal: abort });
}
