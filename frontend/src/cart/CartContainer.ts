import { reactive, watch } from 'vue';

export interface ICartItem {
  scryfallId: string;
  name: string;
  amount: number;
  max: number
}

const STORAGE_KEY: string = 'shoppingCardCart';

function isCartItem(value: unknown): value is ICartItem {
  if (typeof value !== 'object' || value === null) return false;
  const item = value as Record<string, unknown>;
  return (
    typeof item.scryfallId === 'string' &&
    typeof item.name === 'string' &&
    typeof item.amount === 'number' &&
    typeof item.max === 'number'
  );
}

function loadCart(): ICartItem[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return [];

    const parsed: unknown = JSON.parse(raw);
    if (!Array.isArray(parsed)) return [];

    return parsed.filter(isCartItem);
  } catch {
    return [];
  }
}

export const cart = reactive<ICartItem[]>(loadCart());

watch(
  cart,
  (value) => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(value));
  },
  { deep: true }
);

export function removeFromCart(scryfallId: string) {
  const index = cart.findIndex((item) => item.scryfallId === scryfallId);
  if (index !== -1) cart.splice(index, 1);
}

export function addToCart(scryfallId: string, name: string, amount: number, max: number) {
  const existing = cart.find((item) => item.scryfallId === scryfallId);
  if (existing) {
    existing.amount = Math.min(existing.amount + amount, max);
    existing.max = max; // keep max in sync in case data refreshed
  } else {
    cart.push({ scryfallId, name, amount: Math.min(amount, max), max });
  }
}

export function updateAmount(scryfallId: string, newAmount: number) {
  const existing = cart.find((item) => item.scryfallId === scryfallId);
  if (!existing) return;

  if (newAmount <= 0) {
    removeFromCart(scryfallId);
  } else {
    existing.amount = Math.min(newAmount, existing.max);
  }
}