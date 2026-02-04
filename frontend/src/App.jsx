import { useEffect, useState } from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import AppShell from "./layout/AppShell";
import Protected from "./components/Protected";
import Home from "./pages/Home";
import Auth from "./pages/Auth";
import Dashboard from "./pages/Dashboard";
import Courses from "./pages/Courses";
import CourseDetail from "./pages/CourseDetail";
import LessonView from "./pages/LessonView";
import Admin from "./pages/Admin";
import About from "./pages/About";
import Contacts from "./pages/Contacts";
import Modal from "./components/Modal";
import PlacementTest from "./components/PlacementTest";
import { api, clearTokens, getStoredRefreshToken, hasAccessToken } from "./api";
import "./styles.css";

export default function App() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [showTest, setShowTest] = useState(false);
  const [hasProgress, setHasProgress] = useState(false);

  useEffect(() => {
    if (!hasAccessToken()) {
      setLoading(false);
      return;
    }
    api.me()
      .then((data) => {
        setUser(data);
        return api.progress();
      })
      .then((progress) => {
        const anyProgress = (progress || []).some((p) => (p.attempts || 0) > 0);
        setHasProgress(anyProgress);
      })
      .catch(() => clearTokens())
      .finally(() => setLoading(false));
  }, []);

  const onLogout = async () => {
    try {
      await api.logout(getStoredRefreshToken());
    } catch {
      // ignore
    }
    clearTokens();
    setUser(null);
    setHasProgress(false);
  };

  const handleStart = () => {
    if (user && !hasProgress) {
      setShowTest(true);
    }
  };

  return (
    <BrowserRouter>
      <AppShell user={user} loading={loading} onLogout={onLogout}>
        <Routes>
          <Route path="/" element={<Home onStartTest={handleStart} showTestButton={!user || !hasProgress} user={user} />} />
          <Route path="/auth" element={<Auth onAuth={setUser} />} />
          <Route path="/dashboard" element={<Protected user={user}><Dashboard /></Protected>} />
          <Route path="/courses" element={<Protected user={user}><Courses /></Protected>} />
          <Route path="/courses/:id" element={<Protected user={user}><CourseDetail /></Protected>} />
          <Route path="/lessons/:id" element={<Protected user={user}><LessonView /></Protected>} />
          <Route path="/admin" element={<Protected user={user}><Admin user={user} /></Protected>} />
          <Route path="/about" element={<About />} />
          <Route path="/contacts" element={<Contacts />} />
        </Routes>
      </AppShell>

      <Modal open={showTest} title="Тест подбора курса" onClose={() => setShowTest(false)}>
        <PlacementTest />
      </Modal>
    </BrowserRouter>
  );
}
