import type { LocationQuery, RouteMeta } from 'vue-router';
import fetchApi from './api';
import ResponseError from './ResponseError';

export async function loadRouteData(meta: RouteMeta, ...paths: string[]) {
  try {
    const loads: Promise<void>[] = [];
    for (const path of new Set(paths)) {
      loads.push(loadPath(meta, path));
    }
    await Promise.all(loads);
  } catch (e) {
    if (e instanceof ResponseError && e.status === 401) {
      return { name: 'login', replace: true };
    }
  }
}

async function loadPath(meta: RouteMeta, path: string) {
  meta[path] = await fetchApi(path);
}

export function routeData<T>(meta: RouteMeta, path: string) {
  return meta[path] as T;
}

export function toQueryString(value: LocationQuery[string]) {
  return value?.toString() ?? '';
}
