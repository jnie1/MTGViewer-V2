type TimeoutId = Parameters<typeof clearTimeout>[0];

export function timeout(ms: number, signal: AbortSignal): Promise<void> {
  return new Promise((resolve, reject) => {
    let timeoutHandle: TimeoutId = undefined;

    const onAbort = () => {
      clearTimeout(timeoutHandle);
      reject(signal.reason);
    };

    const onTimeout = () => {
      signal.removeEventListener('abort', onAbort);
      resolve();
    };

    signal.addEventListener('abort', onAbort, { once: true });
    timeoutHandle = setTimeout(onTimeout, ms);
  });
}

export function isAbortError(value: unknown): value is Error {
  return value instanceof Error && value.name === 'AbortError';
}
