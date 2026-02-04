import { useEffect, useMemo, useState } from "react";
import { api } from "../api";

function optionsToText(options = []) {
  return options.join("\n");
}

function textToOptions(text) {
  return text
    .split("\n")
    .map((line) => line.trim())
    .filter(Boolean);
}

export default function Admin({ user }) {
  const [courses, setCourses] = useState([]);
  const [selectedCourse, setSelectedCourse] = useState("");
  const [lessons, setLessons] = useState([]);
  const [selectedLesson, setSelectedLesson] = useState("");
  const [quiz, setQuiz] = useState([]);
  const [placement, setPlacement] = useState([]);
  const [error, setError] = useState("");

  const [courseForm, setCourseForm] = useState({ level: "base", title: "", description: "" });
  const [lessonForm, setLessonForm] = useState({ title: "", content: "", order_index: 0 });
  const [quizForm, setQuizForm] = useState({ question: "", options: "", correct_index: 0, explanation: "" });
  const [placementForm, setPlacementForm] = useState({ question: "", options: "", correct_index: 0 });

  const levels = useMemo(() => ["base", "mid", "pro"], []);

  useEffect(() => {
    if (!user?.is_admin) return;
    loadAll();
  }, [user]);

  const loadAll = async () => {
    setError("");
    try {
      const data = await api.adminCourses();
      setCourses(data);
      const placementData = await api.adminPlacement();
      setPlacement(placementData);
    } catch (e) {
      setError(e.message);
    }
  };

  const loadLessons = async (courseId) => {
    setError("");
    try {
      const data = await api.adminLessons(courseId);
      setLessons(data);
    } catch (e) {
      setError(e.message);
    }
  };

  const loadQuiz = async (lessonId) => {
    setError("");
    try {
      const data = await api.adminLessonQuiz(lessonId);
      setQuiz(data);
    } catch (e) {
      setError(e.message);
    }
  };

  const createCourse = async () => {
    setError("");
    try {
      await api.adminCreateCourse(courseForm);
      setCourseForm({ level: "base", title: "", description: "" });
      loadAll();
    } catch (e) {
      setError(e.message);
    }
  };

  const updateCourse = async (course) => {
    setError("");
    try {
      await api.adminUpdateCourse(course.id, course);
      loadAll();
    } catch (e) {
      setError(e.message);
    }
  };

  const deleteCourse = async (id) => {
    setError("");
    try {
      await api.adminDeleteCourse(id);
      if (String(selectedCourse) === String(id)) {
        setSelectedCourse("");
        setLessons([]);
        setSelectedLesson("");
        setQuiz([]);
      }
      loadAll();
    } catch (e) {
      setError(e.message);
    }
  };

  const createLesson = async () => {
    if (!selectedCourse) return;
    setError("");
    try {
      await api.adminCreateLesson(selectedCourse, lessonForm);
      setLessonForm({ title: "", content: "", order_index: 0 });
      loadLessons(selectedCourse);
    } catch (e) {
      setError(e.message);
    }
  };

  const updateLesson = async (lesson) => {
    setError("");
    try {
      await api.adminUpdateLesson(lesson.id, lesson);
      loadLessons(selectedCourse);
    } catch (e) {
      setError(e.message);
    }
  };

  const deleteLesson = async (id) => {
    setError("");
    try {
      await api.adminDeleteLesson(id);
      if (String(selectedLesson) === String(id)) {
        setSelectedLesson("");
        setQuiz([]);
      }
      loadLessons(selectedCourse);
    } catch (e) {
      setError(e.message);
    }
  };

  const createQuizQuestion = async () => {
    if (!selectedLesson) return;
    setError("");
    try {
      await api.adminCreateLessonQuiz(selectedLesson, {
        question: quizForm.question,
        options: textToOptions(quizForm.options),
        correct_index: Number(quizForm.correct_index),
        explanation: quizForm.explanation
      });
      setQuizForm({ question: "", options: "", correct_index: 0, explanation: "" });
      loadQuiz(selectedLesson);
    } catch (e) {
      setError(e.message);
    }
  };

  const updateQuizQuestion = async (q) => {
    setError("");
    try {
      await api.adminUpdateQuiz(q.id, {
        question: q.question,
        options: q.options,
        correct_index: Number(q.correct_index),
        explanation: q.explanation
      });
      loadQuiz(selectedLesson);
    } catch (e) {
      setError(e.message);
    }
  };

  const deleteQuizQuestion = async (id) => {
    setError("");
    try {
      await api.adminDeleteQuiz(id);
      loadQuiz(selectedLesson);
    } catch (e) {
      setError(e.message);
    }
  };

  const createPlacementQuestion = async () => {
    setError("");
    try {
      await api.adminCreatePlacement({
        question: placementForm.question,
        options: textToOptions(placementForm.options),
        correct_index: Number(placementForm.correct_index)
      });
      setPlacementForm({ question: "", options: "", correct_index: 0 });
      loadAll();
    } catch (e) {
      setError(e.message);
    }
  };

  const updatePlacementQuestion = async (q) => {
    setError("");
    try {
      await api.adminUpdatePlacement(q.id, {
        question: q.question,
        options: q.options,
        correct_index: Number(q.correct_index)
      });
      loadAll();
    } catch (e) {
      setError(e.message);
    }
  };

  const deletePlacementQuestion = async (id) => {
    setError("");
    try {
      await api.adminDeletePlacement(id);
      loadAll();
    } catch (e) {
      setError(e.message);
    }
  };

  if (!user?.is_admin) {
    return (
      <div className="page">
        <h2>Админ‑панель</h2>
        <p>Доступ только для администраторов.</p>
      </div>
    );
  }

  return (
    <div className="page">
      <h2>Админ‑панель</h2>
      <p>Управление курсами, уроками и тестами.</p>
      {error && <div className="error">{error}</div>}

      <section className="card" style={{ marginTop: 16 }}>
        <h3>Курсы</h3>
        <div className="form">
          <select value={courseForm.level} onChange={(e) => setCourseForm({ ...courseForm, level: e.target.value })}>
            {levels.map((l) => (
              <option key={l} value={l}>{l}</option>
            ))}
          </select>
          <input placeholder="Название" value={courseForm.title} onChange={(e) => setCourseForm({ ...courseForm, title: e.target.value })} />
          <input placeholder="Описание" value={courseForm.description} onChange={(e) => setCourseForm({ ...courseForm, description: e.target.value })} />
          <button className="btn" onClick={createCourse}>Добавить курс</button>
        </div>
        <div className="lesson-list">
          {courses.map((c) => (
            <div key={c.id} className="lesson-item">
              <div style={{ flex: 1 }}>
                <input value={c.level} onChange={(e) => setCourses(courses.map((x) => x.id === c.id ? { ...x, level: e.target.value } : x))} />
                <input value={c.title} onChange={(e) => setCourses(courses.map((x) => x.id === c.id ? { ...x, title: e.target.value } : x))} />
                <input value={c.description} onChange={(e) => setCourses(courses.map((x) => x.id === c.id ? { ...x, description: e.target.value } : x))} />
              </div>
              <div className="lesson-actions">
                <button className="btn small" onClick={() => updateCourse(c)}>Сохранить</button>
                <button className="ghost" onClick={() => deleteCourse(c.id)}>Удалить</button>
                <button className="ghost" onClick={() => { setSelectedCourse(String(c.id)); loadLessons(c.id); }}>Уроки</button>
              </div>
            </div>
          ))}
        </div>
      </section>

      <section className="card" style={{ marginTop: 16 }}>
        <h3>Уроки</h3>
        <div className="form">
          <select value={selectedCourse} onChange={(e) => { setSelectedCourse(e.target.value); if (e.target.value) loadLessons(e.target.value); }}>
            <option value="">Выберите курс</option>
            {courses.map((c) => (
              <option key={c.id} value={c.id}>{c.title}</option>
            ))}
          </select>
          <input placeholder="Название" value={lessonForm.title} onChange={(e) => setLessonForm({ ...lessonForm, title: e.target.value })} />
          <input placeholder="Контент" value={lessonForm.content} onChange={(e) => setLessonForm({ ...lessonForm, content: e.target.value })} />
          <input type="number" placeholder="Порядок (0 = авто)" value={lessonForm.order_index} onChange={(e) => setLessonForm({ ...lessonForm, order_index: Number(e.target.value) })} />
          <button className="btn" onClick={createLesson} disabled={!selectedCourse}>Добавить урок</button>
        </div>
        <div className="lesson-list">
          {lessons.map((l) => (
            <div key={l.id} className="lesson-item">
              <div style={{ flex: 1 }}>
                <input value={l.title} onChange={(e) => setLessons(lessons.map((x) => x.id === l.id ? { ...x, title: e.target.value } : x))} />
                <input value={l.content} onChange={(e) => setLessons(lessons.map((x) => x.id === l.id ? { ...x, content: e.target.value } : x))} />
                <input type="number" value={l.order_index} onChange={(e) => setLessons(lessons.map((x) => x.id === l.id ? { ...x, order_index: Number(e.target.value) } : x))} />
              </div>
              <div className="lesson-actions">
                <button className="btn small" onClick={() => updateLesson(l)}>Сохранить</button>
                <button className="ghost" onClick={() => deleteLesson(l.id)}>Удалить</button>
                <button className="ghost" onClick={() => { setSelectedLesson(String(l.id)); loadQuiz(l.id); }}>Тест</button>
              </div>
            </div>
          ))}
        </div>
      </section>

      <section className="card" style={{ marginTop: 16 }}>
        <h3>Тест урока</h3>
        <div className="form">
          <select value={selectedLesson} onChange={(e) => { setSelectedLesson(e.target.value); if (e.target.value) loadQuiz(e.target.value); }}>
            <option value="">Выберите урок</option>
            {lessons.map((l) => (
              <option key={l.id} value={l.id}>{l.title}</option>
            ))}
          </select>
          <input placeholder="Вопрос" value={quizForm.question} onChange={(e) => setQuizForm({ ...quizForm, question: e.target.value })} />
          <textarea placeholder="Варианты (каждый с новой строки)" value={quizForm.options} onChange={(e) => setQuizForm({ ...quizForm, options: e.target.value })} />
          <input type="number" placeholder="Индекс правильного (с 0)" value={quizForm.correct_index} onChange={(e) => setQuizForm({ ...quizForm, correct_index: Number(e.target.value) })} />
          <input placeholder="Пояснение" value={quizForm.explanation} onChange={(e) => setQuizForm({ ...quizForm, explanation: e.target.value })} />
          <button className="btn" onClick={createQuizQuestion} disabled={!selectedLesson}>Добавить вопрос</button>
        </div>
        <div className="lesson-list">
          {quiz.map((q) => (
            <div key={q.id} className="lesson-item">
              <div style={{ flex: 1 }}>
                <input value={q.question} onChange={(e) => setQuiz(quiz.map((x) => x.id === q.id ? { ...x, question: e.target.value } : x))} />
                <textarea value={optionsToText(q.options)} onChange={(e) => setQuiz(quiz.map((x) => x.id === q.id ? { ...x, options: textToOptions(e.target.value) } : x))} />
                <input type="number" value={q.correct_index} onChange={(e) => setQuiz(quiz.map((x) => x.id === q.id ? { ...x, correct_index: Number(e.target.value) } : x))} />
                <input value={q.explanation} onChange={(e) => setQuiz(quiz.map((x) => x.id === q.id ? { ...x, explanation: e.target.value } : x))} />
              </div>
              <div className="lesson-actions">
                <button className="btn small" onClick={() => updateQuizQuestion(q)}>Сохранить</button>
                <button className="ghost" onClick={() => deleteQuizQuestion(q.id)}>Удалить</button>
              </div>
            </div>
          ))}
        </div>
      </section>

      <section className="card" style={{ marginTop: 16 }}>
        <h3>Тест подбора уровня</h3>
        <div className="form">
          <input placeholder="Вопрос" value={placementForm.question} onChange={(e) => setPlacementForm({ ...placementForm, question: e.target.value })} />
          <textarea placeholder="Варианты (каждый с новой строки)" value={placementForm.options} onChange={(e) => setPlacementForm({ ...placementForm, options: e.target.value })} />
          <input type="number" placeholder="Индекс правильного (с 0)" value={placementForm.correct_index} onChange={(e) => setPlacementForm({ ...placementForm, correct_index: Number(e.target.value) })} />
          <button className="btn" onClick={createPlacementQuestion}>Добавить вопрос</button>
        </div>
        <div className="lesson-list">
          {placement.map((q) => (
            <div key={q.id} className="lesson-item">
              <div style={{ flex: 1 }}>
                <input value={q.question} onChange={(e) => setPlacement(placement.map((x) => x.id === q.id ? { ...x, question: e.target.value } : x))} />
                <textarea value={optionsToText(q.options)} onChange={(e) => setPlacement(placement.map((x) => x.id === q.id ? { ...x, options: textToOptions(e.target.value) } : x))} />
                <input type="number" value={q.correct_index} onChange={(e) => setPlacement(placement.map((x) => x.id === q.id ? { ...x, correct_index: Number(e.target.value) } : x))} />
              </div>
              <div className="lesson-actions">
                <button className="btn small" onClick={() => updatePlacementQuestion(q)}>Сохранить</button>
                <button className="ghost" onClick={() => deletePlacementQuestion(q.id)}>Удалить</button>
              </div>
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
