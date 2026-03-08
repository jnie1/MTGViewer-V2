import { type ICard } from '@/cards/types';

export interface ILogs {
  groupId: string;
  time: string;
  amount: number;
}

export interface IContainer {
  containerId: number;
  name: string;
}

export interface ITransactionChange {
  fromContainer?: IContainer;
  toContainer?: IContainer;
  card: ICard;
  quantity: number;
}
