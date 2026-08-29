export type Rarity = 'common' | 'uncommon' | 'rare' | 'mythic' | 'special' | 'bonus';

export interface ICard {
  scryfallId: string;
  name: string;
  manaCost?: string;
  type: string;
  rarity: Rarity;
  power?: string;
  toughness?: string;
  imageUrls: {
    preview?: string;
    normal?: string;
    full?: string;
  };
  amount?: number;
}

export interface ISearchResult {
  totalCards: boolean;
  cards: ICard[];
  page: number;
  hasMore: boolean;
}
