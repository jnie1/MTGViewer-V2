import type { ICard } from '@/cards/types';

export interface ICardContainer {
  containerId: string;
  containerName: string;
  amount: number;
}
export interface ICardContainerMatch {
  card: ICard;
  containers: ICardContainer[];
}
