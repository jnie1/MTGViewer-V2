import { reactive, watch } from 'vue';

export interface ICartItem {
  scryfallId: string;
  name: string;
  amount: number;
  max: number
}

const STORAGE_KEY = 'cart';

function loadCart(): ICartItem[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

const cart = reactive<ICartItem[]>(loadCart());

watch(
  cart,
  (value) => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(value));
  },
  { deep: true }
);

export function useCart() {
  

  function removeFromCart(scryfallId: string) {
    const index = cart.findIndex((item) => item.scryfallId === scryfallId);
    if (index !== -1) cart.splice(index, 1);
  }
  
  function addToCart(scryfallId: string, name: string, amount: number, max: number) {
    const existing = cart.find((item) => item.scryfallId === scryfallId);
    if (existing) {
      existing.amount = Math.min(existing.amount + amount, max);
      existing.max = max; // keep max in sync in case data refreshed
    } else {
      cart.push({ scryfallId, name, amount: Math.min(amount, max), max });
    }
  }

  function updateAmount(scryfallId: string, newAmount: number) {
    const existing = cart.find((item) => item.scryfallId === scryfallId);
    if (!existing) return;

    if (newAmount <= 0) {
      removeFromCart(scryfallId);
    } else {
      existing.amount = Math.min(newAmount, existing.max);
    }
  }
  
  return { cart, addToCart, removeFromCart, updateAmount };
}