import { ICard } from "@/cards/types";

export interface ISearchResult {
  HasNextPage: boolean;
  CardResult: ICard[];
}