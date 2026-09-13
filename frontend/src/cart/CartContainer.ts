import { computed, reactive, readonly, watch } from 'vue';

export interface ICartItem {
  scryfallId: string;
  name: string;
  amount: number;
  max: number;
  containerId: number;
}

export interface IScryfallAmount {
  card: {
    scryfallId: string;
  };
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

const items = reactive(loadCart());
export const cart: readonly ICartItem[] = readonly(items);

export const withdrawals = computed(() => {
  const grouped = new Map<number, IScryfallAmount[]>();

  for (const { scryfallId, containerId, amount } of items) {
    let amounts = grouped.get(containerId);
    if (!amounts) {
      amounts = [];
      grouped.set(containerId, amounts);
    }
    amounts.push({ card: { scryfallId }, amount });
  }

  return Object.fromEntries(grouped);
});

watch(
  items,
  (value) => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(value));
  },
  { deep: true },
);

function findItem(scryfallId: string, containerId: number) {
  return items.find((i) => i.scryfallId === scryfallId && i.containerId === containerId);
}

export function addToCart(scryfallId: string, containerId: number, name: string, max: number) {
  const existing = findItem(scryfallId, containerId);
  if (existing) {
    existing.amount = Math.min(existing.amount + 1, max);
    existing.max = max; // keep max in sync in case data refreshed
  } else {
    items.push({
      scryfallId,
      name,
      amount: 1,
      max,
      containerId,
    });
  }
}

export function removeFromCart(scryfallId: string, containerId: number) {
  const existing = findItem(scryfallId, containerId);
  if (!existing) return;

  if (existing.amount > 1) {
    existing.amount -= 1;
  } else {
    removeItem(scryfallId, containerId);
  }
}

export function removeItem(scryfallId: string, containerId: number) {
  const index = items.findIndex(
    (item) => item.scryfallId === scryfallId && item.containerId === containerId,
  );
  if (index !== -1) items.splice(index, 1);
}

export function removeAllCards() {
  items.splice(0);
}
