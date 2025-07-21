export interface User {
  id: number;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  role: {
    id: number;
    name: string;
  };
  is_active: boolean;
  created_at: string;
  last_login_at: string;
  permissions: string[];
}
