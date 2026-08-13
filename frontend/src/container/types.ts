import type { ICard } from '@/cards/types';

export interface IPrints {
  scryfallId: string;
  amount: number;
  images: {
    preview?: string;
    normal?: string;
    full?: string;
  };
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
