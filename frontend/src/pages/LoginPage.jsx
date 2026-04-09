import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { DUMMY_USER } from "../config/constants";
import { useAuth } from "../context/AuthContext";
import whiteLogo from "../assets/white-logo.svg";

export default function LoginPage() {
  const [username, setUsername] = useState(DUMMY_USER.username);
  const [password, setPassword] = useState(DUMMY_USER.password);
  const [error, setError] = useState("");
  const { login } = useAuth();
  const navigate = useNavigate();

  const onSubmit = (e) => {
    e.preventDefault();
    const result = login(username.trim(), password);
    if (!result.ok) {
      setError(result.message || "Login failed");
      return;
    }
    navigate("/dashboard", { replace: true });
  };

  return (
    <div className="login-page">
      <img src={whiteLogo} alt="Best Infra" className="brand-logo" />
      <p className="brand-subtitle">Smart Meter Management System</p>

      <form className="login-card" onSubmit={onSubmit}>
        <h2>Welcome Back</h2>
        <p className="sub">Sign in to your account</p>

        <label>
          Username
          <input value={username} onChange={(e) => setUsername(e.target.value)} required />
        </label>

        <label>
          Password
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </label>

        {error ? <p className="error">{error}</p> : null}

        <button type="submit" className="signin-btn">
          Sign In
        </button>
      </form>
    </div>
  );
}

