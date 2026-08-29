import { reactive, watch } from 'vue';
import fetchApi from '@/fetch/api';

export interface ICartItem {
  scryfallId: string;
  name: string;
  amount: number;
  max: number;
  containerId: number;
}

interface IScryfallId {
  scryfallId: string;
}

interface IScryfallAmount {
  card: IScryfallId;
  amount: number;
}

const STORAGE_KEY: string = 'shoppingCardCart';

function isCartItem(value: unknown): value is ICartItem {
  if (typeof value !== 'object' || value === null) return false;
  const item = value as Record<string, unknown>;
  return (
    typeof item.scryfallId === 'string' &&
    typeof item.name === 'string' &&
    typeof item.amount === 'number' &&
    typeof item.max === 'number' &&
    typeof item.containerId === 'string'
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
  { deep: true },
);

function findItem(scryfallId: string, containerId: number) {
  return cart.find((item) => item.scryfallId === scryfallId && item.containerId === containerId);
}

export function removeFromCart(scryfallId: string, containerId: number) {
  const index = cart.findIndex(
    (item) => item.scryfallId === scryfallId && item.containerId === containerId,
  );
  if (index !== -1) cart.splice(index, 1);
}

export function addToCart(
  scryfallId: string,
  containerId: number,
  name: string,
  amount: number,
  max: number,
) {
  const existing = findItem(scryfallId, containerId);
  if (existing) {
    existing.amount = Math.min(existing.amount + amount, max);
    existing.max = max; // keep max in sync in case data refreshed
  } else {
    cart.push({
      scryfallId,
      name,
      amount: Math.min(amount, max),
      max,
      containerId,
    });
  }
}

export function updateAmount(scryfallId: string, containerId: number, newAmount: number) {
  const existing = findItem(scryfallId, containerId);
  if (!existing) return;

  if (newAmount <= 0) {
    removeFromCart(scryfallId, containerId);
  } else {
    existing.amount = Math.min(newAmount, existing.max);
  }
}

export function removeAllCards() {
  cart.splice(0);
}

export async function submitAllCards() {
  const grouped = new Map<number, IScryfallAmount[]>();

  for (const { scryfallId, containerId, amount } of cart) {
    let amounts = grouped.get(containerId);
    if (!amounts) {
      amounts = [];
      grouped.set(containerId, amounts);
    }
    const card: IScryfallId = { scryfallId };
    amounts.push({ card, amount });
  }

  const payload = Object.fromEntries(grouped);

  await fetchApi('/cards/withdraw', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });

  removeAllCards();
}
