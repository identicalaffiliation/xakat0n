// import client from './client';

// export interface LoginResponse {
//   user_id: string;
//   username: string;
//   token: string;
// }

// export const login = async (username: string): Promise<LoginResponse> => {
//   const res = await client.post<LoginResponse>('/auth/login', { username });
//   return res.data;
// };
import client from './client';
import { USE_MOCK } from './mock';
import { mockLoginResponse } from './mockData';

export interface LoginResponse {
  userId: string;    
  username: string;
  token: string;    
}

export const login = async (username: string): Promise<LoginResponse> => {
  if (USE_MOCK) return mockLoginResponse(username);
  const res = await client.post<LoginResponse>('/auth/login', { username });
  return res.data;
};