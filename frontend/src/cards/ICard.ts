export interface ICard {
  name: string;
  manaCost: string;
  type: string;
  power: string;
  toughness: string;
  imageUrls: {
    preview: string;
    normal: string;
    full: string;
  };
}
