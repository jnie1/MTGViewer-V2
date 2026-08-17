import type { ICard } from '@/cards/types';

export interface IPrints {
  scryfallId: string;
  imagesUrls: ICard['imageUrls'];
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
