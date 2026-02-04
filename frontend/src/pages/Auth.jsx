import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { api, setTokens } from "../api";

export default function Auth({ onAuth }) {
  const [mode, setMode] = useState("login");
  const [form, setForm] = useState({ name: "", email: "", password: "" });
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const submit = async (e) => {
    e.preventDefault();
    setError("");
    try {
      if (mode === "register") {
        await api.register({ name: form.name, email: form.email, password: form.password });
      }
      const res = await api.login({ email: form.email, password: form.password });
      setTokens(res.access_token, res.refresh_token);
      onAuth(res.user);
      navigate("/dashboard");
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div className="page auth">
      <div className="panel">
        <h2>{mode === "login" ? "Вход" : "Регистрация"}</h2>
        <form onSubmit={submit} className="form">
          {mode === "register" && (
            <input
              placeholder="Имя"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              required
            />
          )}
          <input
            type="email"
            placeholder="Email"
            value={form.email}
            onChange={(e) => setForm({ ...form, email: e.target.value })}
            required
          />
          <input
            type="password"
            placeholder="Пароль"
            value={form.password}
            onChange={(e) => setForm({ ...form, password: e.target.value })}
            required
          />
          {error && <div className="error">{error}</div>}
          <button className="btn" type="submit">{mode === "login" ? "Войти" : "Создать аккаунт"}</button>
        </form>
        <button className="ghost" onClick={() => setMode(mode === "login" ? "register" : "login")}>
          {mode === "login" ? "Нет аккаунта? Регистрация" : "Уже есть аккаунт? Войти"}
        </button>
      </div>
    </div>
  );
}
