// Response umum untuk API (tanpa pagination)
export type APIResponse<T = any> = {
  rc?: number;
  success: boolean;
  message: string;
  pagination?: APIPaginationResponse;
  data: T;
};

// Response pagination
export type APIPaginationResponse = {
  page: number;
  page_size: number;
  total_items: number;
  total_pages: number;
  has_more: boolean;
  from: number;
  to: number;
};

// Response khusus untuk list (ada pagination + array data)
export type APIListResponse<T = any> = {
  rc?: number;
  success: boolean;
  message: string;
  pagination: APIPaginationResponse;
  data: T[];
};

export type APIErrorResponse = {
  rc?: number;
  success: false;
  message: string;
  errors?: any;
  error_code?: string;
  dev_message?: string;
  trace?: string;
};

