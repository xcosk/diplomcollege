const API_BASE = import.meta.env.VITE_API_BASE || "http://localhost:8080/api";
const ACCESS_KEY = "access_token";
const REFRESH_KEY = "refresh_token";

function getAccessToken() {
  return localStorage.getItem(ACCESS_KEY) || "";
}

function getRefreshToken() {
  return localStorage.getItem(REFRESH_KEY) || "";
}

export function setTokens(accessToken, refreshToken) {
  if (accessToken) localStorage.setItem(ACCESS_KEY, accessToken);
  if (refreshToken) localStorage.setItem(REFRESH_KEY, refreshToken);
}

export function clearTokens() {
  localStorage.removeItem(ACCESS_KEY);
  localStorage.removeItem(REFRESH_KEY);
}

export function hasAccessToken() {
  return Boolean(getAccessToken());
}

async function refreshAccessToken() {
  const refreshToken = getRefreshToken();
  if (!refreshToken) return "";
  const res = await fetch(`${API_BASE}/refresh`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refresh_token: refreshToken })
  });
  if (!res.ok) {
    clearTokens();
    return "";
  }
  const data = await res.json();
  if (data.access_token) {
    setTokens(data.access_token, "");
    return data.access_token;
  }
  return "";
}

async function request(path, options = {}, retry = true) {
  const headers = { "Content-Type": "application/json", ...(options.headers || {}) };
  const token = getAccessToken();
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers
  });

  if (res.status === 401 && retry) {
    const newToken = await refreshAccessToken();
    if (newToken) {
      return request(path, options, false);
    }
  }

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: "Unknown error" }));
    throw new Error(err.error || "Request failed");
  }

  return res.json();
}

export const api = {
  register: (data) => request("/register", { method: "POST", body: JSON.stringify(data) }),
  login: (data) => request("/login", { method: "POST", body: JSON.stringify(data) }),
  logout: (refreshToken) => request("/logout", { method: "POST", body: JSON.stringify({ refresh_token: refreshToken }) }),
  me: () => request("/me"),
  placementQuestions: () => request("/placement-test"),
  placementSubmit: (answers) => request("/placement-test/submit", { method: "POST", body: JSON.stringify({ answers }) }),
  courses: () => request("/courses"),
  courseDetail: (id) => request(`/courses/${id}`),
  lesson: (id) => request(`/lessons/${id}`),
  submitLessonQuiz: (id, answers) => request(`/lessons-quiz/${id}`, { method: "POST", body: JSON.stringify({ answers }) }),
  progress: () => request("/progress"),

  adminCourses: () => request("/admin/courses"),
  adminCreateCourse: (data) => request("/admin/courses", { method: "POST", body: JSON.stringify(data) }),
  adminUpdateCourse: (id, data) => request(`/admin/courses/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  adminDeleteCourse: (id) => request(`/admin/courses/${id}`, { method: "DELETE" }),

  adminLessons: (courseId) => request(`/admin/courses/${courseId}/lessons`),
  adminCreateLesson: (courseId, data) => request(`/admin/courses/${courseId}/lessons`, { method: "POST", body: JSON.stringify(data) }),
  adminUpdateLesson: (id, data) => request(`/admin/lessons/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  adminDeleteLesson: (id) => request(`/admin/lessons/${id}`, { method: "DELETE" }),

  adminLessonQuiz: (lessonId) => request(`/admin/lessons/${lessonId}/quiz`),
  adminCreateLessonQuiz: (lessonId, data) => request(`/admin/lessons/${lessonId}/quiz`, { method: "POST", body: JSON.stringify(data) }),
  adminUpdateQuiz: (id, data) => request(`/admin/quiz/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  adminDeleteQuiz: (id) => request(`/admin/quiz/${id}`, { method: "DELETE" }),

  adminPlacement: () => request("/admin/placement"),
  adminCreatePlacement: (data) => request("/admin/placement", { method: "POST", body: JSON.stringify(data) }),
  adminUpdatePlacement: (id, data) => request(`/admin/placement/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  adminDeletePlacement: (id) => request(`/admin/placement/${id}`, { method: "DELETE" })
};

export function getStoredRefreshToken() {
  return getRefreshToken();
}
