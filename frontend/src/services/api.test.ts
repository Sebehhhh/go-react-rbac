import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useAuthStore } from '../hooks/useAuth';

// Mock localStorage
const localStorageMock = (() => {
  let store: { [key: string]: string } = {};
  return {
    getItem: vi.fn((key: string) => store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
      store[key] = value.toString();
    }),
    removeItem: vi.fn((key: string) => {
      delete store[key];
    }),
    clear: vi.fn(() => {
      store = {};
    }),
  };
})();

Object.defineProperty(window, 'localStorage', { value: localStorageMock });

// Define the mock axios instance structure
const mockAxiosInstance = {
  interceptors: {
    request: {
      use: vi.fn(),
    },
    response: {
      use: vi.fn(),
    },
  },
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn(),
};

// Mock axios directly. This will replace the actual axios module.
// This factory function runs once when the test file is loaded.
vi.mock('axios', () => ({
  __esModule: true,
  default: {
    create: vi.fn(() => mockAxiosInstance),
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));

// Import `api` and `axios` after mocking.
// `api.ts` will now use the mocked axios.
import api from './api';
import axios from 'axios';

// Cast axios to its mocked type for easier access to mock properties
const mockedAxios = axios as vi.Mocked<typeof axios>;

describe('api interceptors', () => {
  beforeEach(() => {
    // Clear all mocks on the mock instance and top-level axios
    vi.clearAllMocks();
    localStorageMock.clear(); // Clear localStorage mock

    // Reset the auth store state
    useAuthStore.setState({
      accessToken: null,
      refreshToken: null,
      user: null,
      isAuthenticated: false,
    });

    // Reset mock implementations for axios instance methods
    mockAxiosInstance.get.mockReset();
    mockAxiosInstance.post.mockReset();
    mockAxiosInstance.put.mockReset();
    mockAxiosInstance.delete.mockReset();

    // Reset mock implementations for interceptor `use` methods
    mockAxiosInstance.interceptors.request.use.mockReset();
    mockAxiosInstance.interceptors.response.use.mockReset();
  });

  it('should add Authorization header if accessToken exists', async () => {
    useAuthStore.setState({ accessToken: 'test_access_token' });
    
    // Mock the actual request that the interceptor will eventually make
    mockAxiosInstance.get.mockResolvedValue({ data: {} });

    // Call the API method. The interceptor in api.ts will be triggered.
    await api.get('/test');

    // Assert that the request to the mocked axios instance was made with the correct headers
    expect(mockAxiosInstance.get).toHaveBeenCalledWith('/test', {
      headers: {
        'Content-Type': 'application/json',
        Authorization: 'Bearer test_access_token',
      },
    });
  });

  it('should not add Authorization header if accessToken does not exist', async () => {
    mockAxiosInstance.get.mockResolvedValue({ data: {} });

    await api.get('/test');

    expect(mockAxiosInstance.get).toHaveBeenCalledWith('/test', {
      headers: {
        'Content-Type': 'application/json',
      },
    });
  });

  it('should refresh token and retry original request on 401', async () => {
    useAuthStore.setState({ refreshToken: 'old_refresh_token' });

    // Mock the first request to fail with 401
    mockAxiosInstance.get.mockRejectedValueOnce({
      response: { status: 401 },
      config: { url: '/protected', _retry: false },
    });

    // Mock the refresh token request to succeed (using top-level axios.post for refresh endpoint)
    mockedAxios.post.mockResolvedValueOnce({
      data: {
        access_token: 'new_access_token',
        refresh_token: 'new_refresh_token',
      },
    });

    // Mock the retried original request to succeed
    mockAxiosInstance.get.mockResolvedValueOnce({ data: { message: 'Success' } });

    // Call the API method. The response interceptor in api.ts will be triggered.
    const promise = api.get('/protected');

    const response = await promise;

    expect(mockedAxios.post).toHaveBeenCalledWith('/auth/refresh', {
      refresh_token: 'old_refresh_token',
    });
    expect(useAuthStore.getState().accessToken).toBe('new_access_token');
    expect(response.data).toEqual({ message: 'Success' });
    expect(mockAxiosInstance.get).toHaveBeenCalledTimes(2); // Original and retried
  });

  it('should logout if refresh token fails', async () => {
    useAuthStore.setState({ refreshToken: 'old_refresh_token' });

    mockAxiosInstance.get.mockRejectedValueOnce({
      response: { status: 401 },
      config: { url: '/protected', _retry: false },
    });

    mockedAxios.post.mockRejectedValueOnce(new Error('Refresh failed'));

    const logoutSpy = vi.spyOn(useAuthStore.getState(), 'logout');

    await expect(api.get('/protected')).rejects.toThrow('Refresh failed');

    expect(logoutSpy).toHaveBeenCalled();
    expect(useAuthStore.getState().accessToken).toBeNull();
  });

  it('should logout if no refresh token on 401', async () => {
    useAuthStore.setState({ accessToken: 'expired_token', refreshToken: null });

    mockAxiosInstance.get.mockRejectedValueOnce({
      response: { status: 401 },
      config: { url: '/protected', _retry: false },
    });

    const logoutSpy = vi.spyOn(useAuthStore.getState(), 'logout');

    await expect(api.get('/protected')).rejects.toHaveProperty('response.status', 401);

    expect(logoutSpy).toHaveBeenCalled();
    expect(useAuthStore.getState().accessToken).toBeNull();
  });
});