import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { api } from "../api";
import ContentBlocks from "../components/ContentBlocks";

export default function LessonView() {
  const { id } = useParams();
  const [lesson, setLesson] = useState(null);
  const [questions, setQuestions] = useState([]);
  const [answers, setAnswers] = useState({});
  const [result, setResult] = useState(null);
  const [error, setError] = useState("");

  useEffect(() => {
    api.lesson(id)
      .then((res) => {
        setLesson(res.lesson);
        setQuestions(res.questions);
      })
      .catch((e) => setError(e.message));
  }, [id]);

  const submit = async () => {
    setError("");
    try {
      const res = await api.submitLessonQuiz(id, answers);
      setResult(res);
    } catch (e) {
      setError(e.message);
    }
  };

  if (!lesson) return <div className="page">Загрузка...</div>;

  return (
    <div className="page">
      <h2>{lesson.title}</h2>
      <article className="content">
        <ContentBlocks content={lesson.content} />
      </article>

      <section className="quiz">
        <h3>Контрольный тест</h3>
        {questions.map((q) => (
          <div key={q.id} className="question">
            <div className="question-title">{q.question}</div>
            <div className="options">
              {q.options.map((opt, idx) => (
                <label key={idx} className="option">
                  <input
                    type="radio"
                    name={`q-${q.id}`}
                    checked={answers[q.id] === idx}
                    onChange={() => setAnswers({ ...answers, [q.id]: idx })}
                  />
                  {opt}
                </label>
              ))}
            </div>
          </div>
        ))}
        {error && <div className="error">{error}</div>}
        {result && (
          <div className={`result ${result.passed ? "pass" : "fail"}`}>
            Результат: {result.score}%. {result.passed ? "Тест пройден, доступ открыт!" : "Нужно минимум 70%"}
          </div>
        )}
        <button className="btn" onClick={submit}>Отправить ответы</button>
      </section>
    </div>
  );
}
