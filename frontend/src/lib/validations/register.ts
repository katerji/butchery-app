export const UAE_MOBILE_REGEX = /^(\+971|0)(5[0-9])\s?\d{3}\s?\d{4}$/;

export interface RegisterFormValues {
  full_name: string;
  email: string;
  phone: string;
  password: string;
}
