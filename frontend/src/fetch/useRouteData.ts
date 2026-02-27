import { useRoute } from 'vue-router';
import type { NavigationGuardNext, RouteMeta } from 'vue-router';
import fetchApi from './api';
import ResponseError from './ResponseError';

export async function loadRouteData(path: string, meta: RouteMeta, next: NavigationGuardNext) {
  try {
    meta._data = await fetchApi(path);
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
