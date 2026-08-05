import { type ICard } from '@/cards/types';

export interface ISearchResult {
  totalCards: boolean;
  cards: ICard[];
  page: number;
  hasMore: boolean;
}
