import { apiClient } from "./client";

export interface RegisterRequest {
  full_name: string;
  email: string;
  phone: string;
  password: string;
}

export interface RegisterResponse {
  access_token: string;
  refresh_token: string;
}

export function registerCustomer(data: RegisterRequest): Promise<RegisterResponse> {
  return apiClient<RegisterResponse>("/api/v1/auth/register", {
    method: "POST",
    body: JSON.stringify(data),
  });
}
