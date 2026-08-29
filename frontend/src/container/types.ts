import type { ICard } from '@/cards/types';

export interface IPrints {
  scryfallId: string;
  imageUrls: ICard['imageUrls'];
  amount: number;
}

export interface ICardContainer {
  containerId: number;
  name: string;
  amount: number;
  prints: IPrints[];
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
