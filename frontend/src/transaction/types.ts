import { type ICard } from '@/cards/types';
import type { IContainer } from '@/collection/types';

export interface ICardTransaction {
  groupId: string;
  time: string;
  total: number;
}

export interface ICardLog {
  fromContainer?: IContainer;
  toContainer?: IContainer;
  card: ICard;
  amount: number;
}
