export function isAbortError(value: unknown): value is DOMException {
  return value instanceof DOMException && value.name === 'AbortError';
}

export function timeout(ms: number, signal: AbortSignal): Promise<void> {
  return new Promise((resolve, reject) => {
    const onTimeout = () => {
      signal.removeEventListener('abort', onAbort);
      resolve();
    };

    const onAbort = () => {
      clearTimeout(timeoutHandle);
      reject(signal.reason);
    };

    signal.addEventListener('abort', onAbort, { once: true });
    const timeoutHandle = setTimeout(onTimeout, ms);
  });
}
