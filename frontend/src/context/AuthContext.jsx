import { createContext, useContext, useMemo, useState } from "react";
import { DUMMY_USER } from "../config/constants";

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);

  const value = useMemo(
    () => ({
      user,
      isAuthenticated: Boolean(user),
      login: (username, password) => {
        if (username === DUMMY_USER.username && password === DUMMY_USER.password) {
          setUser({ username, fullName: DUMMY_USER.fullName });
          return { ok: true };
        }
        return { ok: false, message: "Invalid username or password" };
      },
      logout: () => setUser(null),
    }),
    [user]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}

