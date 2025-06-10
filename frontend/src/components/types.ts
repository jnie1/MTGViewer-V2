export interface ITransaction {
  transactionId?: number;
  groupId?: number;
  fromContainer?: number;
  toContainer?: number;
  scryfallId?: number;
  quantity?: number;
  time?: string;
}
