import { useEffect } from "react";
import { useNavigate } from "react-router-dom";

export default function Protected({ user, children }) {
  const navigate = useNavigate();
  useEffect(() => {
    if (!user) navigate("/auth");
  }, [user, navigate]);
  if (!user) return null;
  return children;
}
