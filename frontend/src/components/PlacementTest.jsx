import { useEffect, useState } from "react";
import { api } from "../api";

export default function PlacementTest() {
  const [questions, setQuestions] = useState([]);
  const [answers, setAnswers] = useState({});
  const [result, setResult] = useState(null);
  const [error, setError] = useState("");

  useEffect(() => {
    api.placementQuestions().then(setQuestions).catch((e) => setError(e.message));
  }, []);

  const submit = async () => {
    setError("");
    try {
      const res = await api.placementSubmit(answers);
      setResult(res);
    } catch (e) {
      setError(e.message);
    }
  };

  return (
    <div>
      <div className="space-y-4">
        {questions.map((q) => (
          <div key={q.id} className="rounded-xl border border-slate-100 p-4">
            <div className="text-sm font-semibold text-slate-800">{q.question}</div>
            <div className="mt-2 space-y-2">
              {q.options.map((opt, idx) => (
                <label key={idx} className="flex items-center gap-2 text-sm text-slate-600">
                  <input
                    type="radio"
                    name={`pq-${q.id}`}
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
      </div>
      {error && <div className="mt-4 text-sm font-semibold text-red-600">{error}</div>}
      {result && (
        <div className="mt-4 rounded-xl bg-brand-50 p-4 text-sm text-brand-700">
          Вы ответили правильно на {result.score} из {result.total}. Рекомендация: {result.message}
        </div>
      )}
      <button className="mt-4 w-full rounded-xl bg-brand-500 px-4 py-2 text-sm font-semibold text-white" onClick={submit}>
        Получить рекомендацию
      </button>
    </div>
  );
}
