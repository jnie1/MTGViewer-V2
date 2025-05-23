export interface ITransaction {
  transaction_id?:number;
	group_id?:number;
	from_container?:number;
	to_container?:number;
	scryfall_id?:number;
	quantity?:number;
	time?:Date;
}