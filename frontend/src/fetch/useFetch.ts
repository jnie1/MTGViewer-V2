import { ref, type Ref } from 'vue';
import ResponseError from './ResponeError';
import fetchApi from './api';

type FetchState<T> = {
  data: Ref<T | undefined>;
  error: Ref<ResponseError | undefined>;
};

function useFetch<T>(path: string): FetchState<T> {
  const data = ref<T>();
  const error = ref<ResponseError>();

  fetchApi<T>(path)
    .then((c) => {
      data.value = c;
    })
    .catch((e) => {
      if (e instanceof ResponseError) {
        error.value = e;
      } else {
        throw e;
      }
    });

  return { data, error };
}

export { type FetchState };
export default useFetch;
