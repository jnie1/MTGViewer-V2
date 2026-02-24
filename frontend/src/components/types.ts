import {ICard} from '@/cards/types';

export interface ILogs {
  groupId?: number;
  time?: string;
}

export interface IContainer {
  containerId: number;
  name: string;
}

export interface ITransactionProps {
  groupId?: number;
  fromContainer?: IContainer;
  toContainer?: IContainer;
  card?: ICard;
  quantity?: number;
}

export interface ITransactionDetailProps {
  boxId?: number;
  transactions?: ITransactionProps;
}