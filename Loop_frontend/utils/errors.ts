export class NetworkError extends Error {
  constructor() {
    super('Unable to reach the server. Please check your connection.');
    this.name = 'NetworkError';
  }
}

export class TimeoutError extends Error {
  constructor() {
    super('Request timed out. Please try again.');
    this.name = 'TimeoutError';
  }
}
