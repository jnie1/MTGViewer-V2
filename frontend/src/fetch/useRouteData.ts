import { useRoute } from 'vue-router';
import type { NavigationGuardNext, RouteMeta } from 'vue-router';
import fetchApi from './api';
import ResponseError from './ResponeError';

export async function loadRouteData(path: string, meta: RouteMeta, next: NavigationGuardNext) {
  try {
    const card = await fetchApi(path);
    meta._data = card;
    next();
  } catch (e) {
    if (e instanceof ResponseError) {
      next(e);
    } else if (e instanceof Error) {
      next(e);
    } else {
      next(false);
    }
  }
}

export function useRouteData<T>() {
  const route = useRoute();
  return route.meta._data as T;
}
