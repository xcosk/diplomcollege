import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { api } from "../api";

export default function Dashboard() {
  const [data, setData] = useState([]);
  const [error, setError] = useState("");

  useEffect(() => {
    api.progress().then(setData).catch((e) => setError(e.message));
  }, []);

  return (
    <div className="bg-gradient-to-b from-[#f6fbff] via-white to-white">
      <div className="mx-auto max-w-6xl px-6 py-10">
        <h1 className="text-3xl font-bold text-slate-900">Личный кабинет</h1>
        <p className="mt-2 text-slate-500">Отслеживайте свой прогресс и продолжайте обучение</p>
        {error && <div className="mt-4 text-sm font-semibold text-red-600">{error}</div>}

        <div className="mt-10 grid gap-6 lg:grid-cols-[2fr_1fr]">
          <div className="space-y-6">
            <h2 className="text-lg font-semibold text-slate-900">Мои курсы</h2>
            {data.map((item, idx) => (
              <div key={item.course_id} className="rounded-2xl border border-slate-100 bg-white p-6 shadow-card">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="text-base font-semibold text-slate-900">{item.title}</h3>
                    <p className="mt-1 text-xs text-slate-500">Урок {Math.max(1, item.passed + 1)}: Продолжение</p>
                  </div>
                  <span className="rounded-full bg-brand-100 px-3 py-1 text-xs font-semibold text-brand-700">
                    {item.progress}%
                  </span>
                </div>
                <div className="mt-4">
                  <div className="text-xs text-slate-500">Прогресс</div>
                  <div className="mt-2 h-2 rounded-full bg-slate-100">
                    <div className="h-2 rounded-full bg-brand-500" style={{ width: `${item.progress}%` }} />
                  </div>
                </div>
                <div className="mt-4 flex items-center justify-between text-xs text-slate-500">
                  <span>⏱ {12 + idx * 2} часов</span>
                  <Link className="inline-flex items-center gap-2 rounded-xl bg-brand-500 px-4 py-2 text-xs font-semibold text-white" to={`/courses/${item.course_id}`}>
                    Продолжить →
                  </Link>
                </div>
              </div>
            ))}

            <div className="rounded-2xl border border-dashed border-slate-200 bg-white p-6">
              <h3 className="text-base font-semibold text-slate-900">Начать новый курс</h3>
              <p className="mt-1 text-sm text-slate-500">Расширьте свои знания Go</p>
              <Link className="mt-4 inline-flex w-full items-center justify-center rounded-xl border border-slate-200 px-4 py-2 text-sm font-semibold text-slate-600" to="/courses">
                Посмотреть все курсы
              </Link>
            </div>
          </div>

          <aside className="space-y-6">
            <div className="rounded-2xl border border-slate-100 bg-white p-6 shadow-card">
              <h3 className="text-base font-semibold text-slate-900">Ваши награды</h3>
              <div className="mt-4 space-y-3">
                {[
                  { title: "Первые шаги", text: "Начал первый курс" },
                  { title: "Настойчивость", text: "7 дней подряд" },
                  { title: "Половина пути", text: "50% курса" }
                ].map((item) => (
                  <div key={item.title} className="rounded-xl bg-brand-50 px-4 py-3">
                    <div className="text-sm font-semibold text-slate-900">{item.title}</div>
                    <div className="text-xs text-slate-500">{item.text}</div>
                  </div>
                ))}
              </div>
            </div>

            <div className="rounded-2xl bg-gradient-to-r from-brand-500 to-brand-300 p-6 text-white shadow-soft">
              <h3 className="text-base font-semibold">Совет дня</h3>
              <p className="mt-2 text-sm text-white/90">
                Практикуйтесь каждый день хотя бы 30 минут. Регулярность важнее продолжительности!
              </p>
            </div>
          </aside>
        </div>
      </div>
    </div>
  );
}
