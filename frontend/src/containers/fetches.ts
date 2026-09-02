import fetchApi from '@/fetch/api';
import { type IContainerPrunePreview } from './types';

export function previewPrune(quantity: number, price: number, abort: AbortSignal) {
  const params = new URLSearchParams({
    size: quantity.toString(),
    price: price.toString(),
  });
  return fetchApi<IContainerPrunePreview>(`containers/prune?${params}`, { signal: abort });
}
