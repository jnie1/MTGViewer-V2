import { type ICard } from '@/cards/types';
import type { IContainerPreview } from '@/collection/types';

export interface ICardTransaction {
  groupId: string;
  time: string;
  total: number;
}

export interface ICardLog {
  fromContainer?: IContainerPreview;
  toContainer?: IContainerPreview;
  card: ICard;
  amount: number;
}
