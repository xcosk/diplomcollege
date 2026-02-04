import { useEffect, useMemo, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { api } from "../api";
import ContentBlocks from "../components/ContentBlocks";

function lessonDuration(orderIndex) {
  const map = { 1: 15, 2: 20, 3: 25, 4: 30, 5: 35, 6: 40, 7: 45 };
  return map[orderIndex] || 20;
}

export default function CourseDetail() {
  const { id } = useParams();
  const [course, setCourse] = useState(null);
  const [lessons, setLessons] = useState([]);
  const [selectedId, setSelectedId] = useState(null);
  const [activeTab, setActiveTab] = useState("materials");
  const [lessonData, setLessonData] = useState(null);
  const [lessonError, setLessonError] = useState("");
  const [lessonLoading, setLessonLoading] = useState(false);
  const [answers, setAnswers] = useState({});
  const [submitResult, setSubmitResult] = useState(null);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  const fetchCourse = () =>
    api.courseDetail(id)
      .then((res) => {
        setCourse(res.course);
        setLessons(res.lessons);
        if (res.lessons?.length && !selectedId) {
          setSelectedId(res.lessons[0].id);
        }
        return res.lessons || [];
      });

  useEffect(() => {
    fetchCourse().catch((e) => setError(e.message));
  }, [id]);

  const selectedLesson = useMemo(
    () => lessons.find((l) => l.id === selectedId) || lessons[0],
    [lessons, selectedId]
  );

  const passedCount = lessons.filter((l) => l.passed).length;

  useEffect(() => {
    if (!selectedLesson) return;
    setAnswers({});
    setSubmitResult(null);
    if (selectedLesson.locked) {
      setLessonData(null);
      setLessonError("");
      return;
    }
    setLessonLoading(true);
    setLessonError("");
    api.lesson(selectedLesson.id)
      .then((res) => setLessonData(res))
      .catch((e) => setLessonError(e.message))
      .finally(() => setLessonLoading(false));
  }, [selectedLesson]);

  const submitTest = async () => {
    if (!selectedLesson) return;
    setSubmitting(true);
    setSubmitResult(null);
    try {
      const res = await api.submitLessonQuiz(selectedLesson.id, answers);
      setSubmitResult(res);
      if (res.passed) {
        const updatedLessons = await fetchCourse();
        const currentIndex = updatedLessons.findIndex((l) => l.id === selectedLesson.id);
        const next = updatedLessons[currentIndex + 1];
        if (next) {
          setSelectedId(next.id);
          setActiveTab("materials");
        }
      }
    } catch (e) {
      setSubmitResult({ error: e.message });
    } finally {
      setSubmitting(false);
    }
  };

  if (!course) return <div className="page">Загрузка...</div>;

  return (
    <div className="bg-gradient-to-b from-[#f6fbff] via-white to-white">
      <div className="mx-auto max-w-6xl px-6 py-10">
        <Link to="/courses" className="text-sm text-slate-500 hover:text-brand-600">← Назад к курсам</Link>
        <h1 className="mt-4 text-3xl font-bold text-slate-900 sm:text-4xl">{course.title}</h1>
        <p className="mt-2 text-slate-500">{course.description}</p>

        <div className="mt-4 flex flex-wrap items-center gap-4 text-sm text-slate-500">
          <div className="flex items-center gap-2">
            <span className="inline-flex h-8 w-8 items-center justify-center rounded-full bg-brand-100 text-brand-600">📘</span>
            {lessons.length} уроков
          </div>
          <div className="flex items-center gap-2">
            <span className="inline-flex h-8 w-8 items-center justify-center rounded-full bg-brand-100 text-brand-600">✅</span>
            {passedCount} пройдено
          </div>
        </div>

        {error && <div className="mt-4 text-sm font-semibold text-red-600">{error}</div>}

        <div className="mt-10 grid gap-6 lg:grid-cols-[320px_1fr]">
          <aside className="rounded-2xl border border-slate-100 bg-white p-6 shadow-card">
            <div className="text-base font-semibold text-slate-900">Уроки курса</div>
            <div className="mt-4 space-y-3">
              {lessons.map((l) => (
                <button
                  key={l.id}
                  className={`flex w-full items-start justify-between rounded-xl border px-4 py-3 text-left transition ${
                    l.id === selectedLesson?.id
                      ? "border-brand-400 bg-brand-50"
                      : "border-slate-100 bg-slate-50/40"
                  } ${l.locked ? "opacity-60" : "hover:border-brand-200"}`}
                  onClick={() => setSelectedId(l.id)}
                >
                  <div>
                    <div className="text-sm font-semibold text-slate-900">Урок {l.order_index}: {l.title}</div>
                    <div className="mt-1 text-xs text-slate-500">{lessonDuration(l.order_index)} мин</div>
                  </div>
                  <div className="text-xs">
                    {l.passed ? (
                      <span className="rounded-full bg-emerald-100 px-2 py-1 text-emerald-700">Пройдено</span>
                    ) : l.locked ? (
                      <span className="text-slate-400">🔒</span>
                    ) : (
                      <span className="text-brand-600">▶</span>
                    )}
                  </div>
                </button>
              ))}
            </div>
          </aside>

          <section className="space-y-6">
            <div className="rounded-2xl border border-slate-100 bg-white p-6 shadow-card">
              <div className="flex items-start justify-between">
                <div>
                  <h2 className="text-lg font-semibold text-slate-900">{selectedLesson?.title || "Выберите урок"}</h2>
                  <div className="mt-1 text-xs text-slate-500">{lessonDuration(selectedLesson?.order_index || 1)} мин</div>
                </div>
                {selectedLesson?.passed && (
                  <span className="rounded-full bg-brand-100 px-3 py-1 text-xs font-semibold text-brand-700">Пройдено</span>
                )}
              </div>
            </div>

            <div className="rounded-2xl border border-slate-100 bg-white p-6 shadow-card">
              <div className="flex items-center gap-6 border-b border-slate-100 pb-4 text-sm font-semibold text-slate-600">
                <button
                  className={activeTab === "materials" ? "text-brand-600" : "hover:text-brand-600"}
                  onClick={() => setActiveTab("materials")}
                >
                  Материалы
                </button>
                <button
                  className={activeTab === "test" ? "text-brand-600" : "hover:text-brand-600"}
                  onClick={() => setActiveTab("test")}
                >
                  Тест
                </button>
              </div>
              <div className="mt-6 rounded-2xl border-2 border-brand-400/60 bg-slate-50 p-6">
                {selectedLesson?.locked ? (
                  <div className="text-sm text-slate-500">
                    Чтобы открыть материалы, пройдите предыдущий урок.
                  </div>
                ) : lessonLoading ? (
                  <div className="text-sm text-slate-500">Загрузка...</div>
                ) : lessonError ? (
                  <div className="text-sm font-semibold text-red-600">{lessonError}</div>
                ) : activeTab === "materials" ? (
                  <ContentBlocks content={lessonData?.lesson?.content || ""} />
                ) : (
                  <div className="space-y-4">
                    {(lessonData?.questions || []).map((q) => (
                      <div key={q.id} className="rounded-xl border border-slate-100 bg-white p-4">
                        <div className="text-sm font-semibold text-slate-800">{q.question}</div>
                        <div className="mt-3 space-y-2 text-sm text-slate-600">
                          {q.options.map((opt, idx) => (
                            <label key={idx} className="flex items-center gap-2">
                              <input
                                type="radio"
                                name={`q-${q.id}`}
                                className="accent-brand-500"
                                checked={answers[q.id] === idx}
                                onChange={() => setAnswers({ ...answers, [q.id]: idx })}
                              />
                              {opt}
                            </label>
                          ))}
                        </div>
                      </div>
                    ))}
                    {(lessonData?.questions || []).length === 0 && (
                      <div className="text-sm text-slate-500">Тест пока не добавлен.</div>
                    )}
                    {submitResult?.error && (
                      <div className="text-sm font-semibold text-red-600">{submitResult.error}</div>
                    )}
                    {submitResult && !submitResult.error && (
                      <div className={`rounded-xl px-4 py-3 text-sm ${submitResult.passed ? "bg-emerald-50 text-emerald-700" : "bg-red-50 text-red-700"}`}>
                        Результат: {submitResult.score}% — {submitResult.passed ? "тест пройден" : "нужно минимум 70%"}.
                      </div>
                    )}
                    <button
                      className="mt-2 inline-flex w-full items-center justify-center rounded-xl bg-brand-500 px-4 py-2 text-sm font-semibold text-white"
                      onClick={submitTest}
                      disabled={submitting || (lessonData?.questions || []).length === 0}
                    >
                      {submitting ? "Проверяем..." : "Отправить ответы"}
                    </button>
                  </div>
                )}
              </div>
            </div>
          </section>
        </div>
      </div>
    </div>
  );
}
