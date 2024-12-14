class ResponseError extends Error {
  status: number;
  type: string;
  url: string;

  constructor(response: Response) {
    super('received an error response');
    this.status = response.status;
    this.type = response.type;
    this.url = response.url;
  }
}

export default ResponseError;
