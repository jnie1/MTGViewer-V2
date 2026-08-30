import type { ICard, ICardPrint } from '@/cards/types';

export interface ICardContainer {
  containerId: number;
  name: string;
  amount: number;
  prints: ICardPrint[];
}

export interface ICardContainerMatch {
  card: ICard;
  containers: ICardContainer[];
}

export interface IContainer {
  containerId: number;
  name: string;
  used: number;
  capacity: number;
}

export interface IContainerPreview {
  containerId: number;
  name: string;
}
