export default function About() {
  return (
    <div className="bg-gradient-to-b from-[#f6fbff] via-white to-white">
      <div className="mx-auto max-w-4xl px-6 py-12">
        <h1 className="text-3xl font-bold text-slate-900">О платформе GoLearn</h1>
        <p className="mt-4 text-base text-slate-600">
          GoLearn — это практическая платформа обучения Golang с упором на реальные навыки.
          Мы даём структурированные курсы, задания и тесты, чтобы вы уверенно перешли от
          базовых понятий к production‑подходам.
        </p>
        <div className="mt-6 rounded-2xl border border-slate-100 bg-white p-6 shadow-card">
          <h2 className="text-lg font-semibold text-slate-900">Что внутри</h2>
          <ul className="mt-3 list-disc space-y-2 pl-6 text-sm text-slate-600">
            <li>Базовый, средний и профессиональный уровни</li>
            <li>Материалы с примерами кода и практикой</li>
            <li>Контрольные тесты после каждого урока</li>
            <li>Прогресс и рекомендации курсов</li>
          </ul>
        </div>
      </div>
    </div>
  );
}
