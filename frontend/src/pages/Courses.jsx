import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { api } from "../api";

export default function Courses() {
  const [courses, setCourses] = useState([]);
  const [progress, setProgress] = useState([]);
  const [error, setError] = useState("");

  const params = new URLSearchParams(window.location.search);
  const filterLevel = params.get("level");

  useEffect(() => {
    Promise.all([api.courses(), api.progress()])
      .then(([coursesData, progressData]) => {
        setCourses(coursesData);
        setProgress(progressData || []);
      })
      .catch((e) => setError(e.message));
  }, []);

  const progressMap = Object.fromEntries(progress.map((p) => [p.course_id, p]));
  const filteredCourses = useMemo(() => {
    if (!filterLevel) return courses;
    return courses.filter((c) => c.level === filterLevel);
  }, [courses, filterLevel]);

  return (
    <div className="bg-gradient-to-b from-[#f6fbff] via-white to-white">
      <div className="mx-auto max-w-6xl px-6 py-10">
        <h1 className="text-3xl font-bold text-slate-900">Курсы</h1>
        <p className="mt-2 text-slate-500">Выберите курс и начните обучение.</p>
        {error && <div className="mt-4 text-sm font-semibold text-red-600">{error}</div>}

        <div className="mt-8 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {filteredCourses.map((c, idx) => (
            <div key={c.id} className="overflow-hidden rounded-2xl border border-slate-100 bg-white shadow-card">
              <div className="h-40 bg-gradient-to-r from-brand-200 via-brand-300 to-brand-200" />
              <div className="p-6">
                <div className="inline-flex rounded-full bg-brand-100 px-3 py-1 text-xs font-semibold text-brand-700">
                  {idx === 0 ? "Новичок" : idx === 1 ? "Средний" : "Профи"}
                </div>
                <h3 className="mt-3 text-lg font-semibold text-slate-900">{c.title}</h3>
                <p className="mt-2 text-sm text-slate-500">{c.description}</p>
                {progressMap[c.id] && (
                  <div className="mt-4">
                    <div className="flex items-center justify-between text-xs text-slate-500">
                      <span>Прогресс</span>
                      <span>{progressMap[c.id].progress}%</span>
                    </div>
                    <div className="mt-2 h-2 rounded-full bg-slate-100">
                      <div
                        className="h-2 rounded-full bg-brand-500"
                        style={{ width: `${progressMap[c.id].progress}%` }}
                      />
                    </div>
                    <div className="mt-2 text-xs text-slate-400">
                      Пройдено {progressMap[c.id].passed} из {progressMap[c.id].total} уроков
                    </div>
                    <div className="mt-3 flex items-center gap-4 text-xs text-slate-500">
                      <span className="inline-flex items-center gap-1">
                        <span className="text-brand-500">⏱</span>
                        {12 + idx * 2} часов
                      </span>
                      <span className="inline-flex items-center gap-1">
                        <span className="text-brand-500">✅</span>
                        {progressMap[c.id].passed} пройдено
                      </span>
                    </div>
                  </div>
                )}
                <Link className="mt-4 inline-flex w-full items-center justify-center rounded-xl bg-brand-500 px-4 py-2 text-sm font-semibold text-white" to={`/courses/${c.id}`}>
                  Начать обучение →
                </Link>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
