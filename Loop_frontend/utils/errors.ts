export const AUTH_ERROR_CODES = {
  BACKEND_DOWN: "BACKEND_DOWN",
  INVALID_CREDENTIALS: "INVALID_CREDENTIALS",
  USER_NOT_FOUND: "USER_NOT_FOUND",
  EMAIL_EXISTS: "EMAIL_EXISTS",
  NETWORK_ERROR: "NETWORK_ERROR"
};

export class AuthError extends Error {
  code: string;
  
  constructor(message: string, code: string) {
    super(message);
    this.name = "AuthError";
    this.code = code;
  }
}

export class NetworkError extends Error {
  response?: Response;
  responseBody?: any;

  constructor(message?: string, response?: Response, responseBody?: any) {
    super(
      message || "Unable to reach the server. Please check your connection.",
    );
    this.name = "NetworkError";
    this.response = response;
    this.responseBody = responseBody;
  }
}

export class TimeoutError extends Error {
  constructor() {
    super("Request timed out. Please try again.");
    this.name = "TimeoutError";
  }
}

export const ERROR_MESSAGES = {
  [AUTH_ERROR_CODES.BACKEND_DOWN]: "Our servers are currently unavailable. Please try again in a few minutes.",
  [AUTH_ERROR_CODES.INVALID_CREDENTIALS]: "Incorrect password. Please check your password and try again.",
  [AUTH_ERROR_CODES.USER_NOT_FOUND]: "No account found with this email. Please check your email or register.",
  [AUTH_ERROR_CODES.EMAIL_EXISTS]: "This email is already registered. Try logging in instead.",
  [AUTH_ERROR_CODES.NETWORK_ERROR]: "Unable to connect. Please check your internet connection."
};
