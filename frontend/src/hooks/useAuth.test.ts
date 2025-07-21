import { act, renderHook } from '@testing-library/react';
import { useAuthStore } from './useAuth';

describe('useAuthStore', () => {
  beforeEach(() => {
    // Reset the store before each test
    useAuthStore.setState({
      accessToken: null,
      refreshToken: null,
      user: null,
      isAuthenticated: false,
    });
    localStorage.clear(); // Clear localStorage as well
  });

  it('should initialize with default values', () => {
    const { result } = renderHook(() => useAuthStore());
    expect(result.current.accessToken).toBeNull();
    expect(result.current.refreshToken).toBeNull();
    expect(result.current.user).toBeNull();
    expect(result.current.isAuthenticated).toBe(false);
  });

  it('should set tokens correctly', () => {
    const { result } = renderHook(() => useAuthStore());
    const newAccessToken = 'new_access_token';
    const newRefreshToken = 'new_refresh_token';

    act(() => {
      result.current.setTokens(newAccessToken, newRefreshToken);
    });

    expect(result.current.accessToken).toBe(newAccessToken);
    expect(result.current.refreshToken).toBe(newRefreshToken);
    expect(result.current.isAuthenticated).toBe(false); // Setting tokens alone doesn't authenticate
  });

  it('should log in a user', () => {
    const { result } = renderHook(() => useAuthStore());
    const mockUser = { id: 1, email: 'test@example.com', username: 'testuser', first_name: 'Test', last_name: 'User', role: { id: 1, name: 'User' }, is_active: true, created_at: '', last_login_at: '' };

    act(() => {
      result.current.login('access123', 'refresh123', mockUser);
    });

    expect(result.current.accessToken).toBe('access123');
    expect(result.current.refreshToken).toBe('refresh123');
    expect(result.current.user).toEqual(mockUser);
    expect(result.current.isAuthenticated).toBe(true);
  });

  it('should log out a user', () => {
    const { result } = renderHook(() => useAuthStore());
    const mockUser = { id: 1, email: 'test@example.com', username: 'testuser', first_name: 'Test', last_name: 'User', role: { id: 1, name: 'User' }, is_active: true, created_at: '', last_login_at: '' };

    act(() => {
      result.current.login('access123', 'refresh123', mockUser);
    });

    expect(result.current.isAuthenticated).toBe(true);

    act(() => {
      result.current.logout();
    });

    expect(result.current.accessToken).toBeNull();
    expect(result.current.refreshToken).toBeNull();
    expect(result.current.user).toBeNull();
    expect(result.current.isAuthenticated).toBe(false);
  });

  it('should persist state to localStorage', () => {
    const { result } = renderHook(() => useAuthStore());
    const mockUser = { id: 1, email: 'test@example.com', username: 'testuser', first_name: 'Test', last_name: 'User', role: { id: 1, name: 'User' }, is_active: true, created_at: '', last_login_at: '' };

    act(() => {
      result.current.login('access123', 'refresh123', mockUser);
    });

    // Simulate re-rendering or new session by re-rendering the hook
    const { result: newResult } = renderHook(() => useAuthStore());
    expect(newResult.current.accessToken).toBe('access123');
    expect(newResult.current.refreshToken).toBe('refresh123');
    expect(newResult.current.user).toEqual(mockUser);
    expect(newResult.current.isAuthenticated).toBe(true);
  });
});
