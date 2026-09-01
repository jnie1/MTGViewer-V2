import { type ICard } from '@/cards/types';
import type { IContainerPreview } from '@/containers/types';

export interface ICardTransaction {
  groupId: string;
  time: string;
  total: number;
  description?: string;
}

export interface ICardLog {
  fromContainer?: IContainerPreview;
  toContainer?: IContainerPreview;
  card: ICard;
  amount: number;
}

export interface IContainerTransfers {
  containerId: number;
  containerName: string;
  total: number;
  cards: ICardTransfer[];
}

export interface ICardTransfer {
  scryfallId: string;
  name: string;
  imageUrls: ICard['imageUrls'];
  delta: number;
  withContainerId?: number;
}
