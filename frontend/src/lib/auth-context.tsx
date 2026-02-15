"use client";

import { createContext, useContext, useState, useEffect, useCallback } from "react";
import type { ReactNode } from "react";

interface AuthTokens {
  access_token: string;
  refresh_token: string;
}

interface AuthContextValue {
  isAuthenticated: boolean;
  accessToken: string | null;
  refreshToken: string | null;
  login: (tokens: AuthTokens) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | null>(null);

const ACCESS_TOKEN_KEY = "access_token";
const REFRESH_TOKEN_KEY = "refresh_token";

export function AuthProvider({ children }: { children: ReactNode }) {
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [refreshToken, setRefreshToken] = useState<string | null>(null);
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    const storedAccess = localStorage.getItem(ACCESS_TOKEN_KEY);
    const storedRefresh = localStorage.getItem(REFRESH_TOKEN_KEY);
    if (storedAccess) {
      setAccessToken(storedAccess);
    }
    if (storedRefresh) {
      setRefreshToken(storedRefresh);
    }
    setLoaded(true);
  }, []);

  const login = useCallback((tokens: AuthTokens) => {
    localStorage.setItem(ACCESS_TOKEN_KEY, tokens.access_token);
    localStorage.setItem(REFRESH_TOKEN_KEY, tokens.refresh_token);
    setAccessToken(tokens.access_token);
    setRefreshToken(tokens.refresh_token);
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    setAccessToken(null);
    setRefreshToken(null);
  }, []);

  if (!loaded) {
    return null;
  }

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated: accessToken !== null,
        accessToken,
        refreshToken,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
