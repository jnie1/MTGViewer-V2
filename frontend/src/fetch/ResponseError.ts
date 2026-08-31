class ResponseError extends Error {
  status: number;
  type: string;
  url: string;

  constructor(response: Response, message?: string) {
    super(message ?? 'received an error response');
    this.status = response.status;
    this.type = response.type;
    this.url = response.url;
  }
}

export function isStatusError(err: unknown, statusCode: number): err is ResponseError {
  return err instanceof ResponseError && err.status === statusCode;
}

export default ResponseError;
