export interface ApiResponse<T> {
  data: T | null;
  error: string | null;
}

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

export async function apiClient<T>(
  endpoint: string,
  options?: RequestInit,
): Promise<T> {
  const { headers, ...rest } = options ?? {};

  const response = await fetch(`${API_URL}${endpoint}`, {
    ...rest,
    headers: {
      "Content-Type": "application/json",
      ...headers,
    },
  });

  let body: ApiResponse<T>;
  try {
    body = await response.json();
  } catch {
    throw new ApiError("An unexpected error occurred", response.status);
  }

  if (!response.ok) {
    throw new ApiError(
      body.error ?? "An unexpected error occurred",
      response.status,
    );
  }

  if (body.data === null) {
    throw new ApiError("No data returned", response.status);
  }

  return body.data;
}
