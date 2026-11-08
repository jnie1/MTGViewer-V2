import fetchApi from '@/fetch/api';
import { type IPrunePreview } from './types';

export function previewPrune(quantity: number, price: number, abort: AbortSignal) {
  const params = new URLSearchParams({
    size: quantity.toString(),
    price: price.toString(),
  });
  return fetchApi<IPrunePreview>(`containers/prune?${params}`, { signal: abort });
}

export function commitPrune(size: number, price: number, signal?: AbortSignal) {
  const query = new URLSearchParams({ size: String(size), price: String(price) });
  return fetchApi(`/containers/prune?${query}`, { method: 'POST', signal });
}
