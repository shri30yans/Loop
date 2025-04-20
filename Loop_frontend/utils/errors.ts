export class NetworkError extends Error {
  response?: Response;
  responseBody?: any;

  constructor(message?: string, response?: Response, responseBody?: any) {
    super(message || 'Unable to reach the server. Please check your connection.');
    this.name = 'NetworkError';
    this.response = response;
    this.responseBody = responseBody;
  }
}

export class TimeoutError extends Error {
  constructor() {
    super('Request timed out. Please try again.');
    this.name = 'TimeoutError';
  }
}
