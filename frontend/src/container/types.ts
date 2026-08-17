import type { ICard } from '@/cards/types';

export interface IPrints {
  scryfallId: string;
  imageUrls: ICard['imageUrls'];
  amount: number;
}

export interface ICardContainer {
  containerId: string;
  name: string;
  amount: number;
  prints: IPrints[];
}

export interface ICardContainerMatch {
  card: ICard;
  containers: ICardContainer[];
}
