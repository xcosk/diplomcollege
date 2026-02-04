import { Link } from "react-router-dom";

export default function Home({ onStartTest, showTestButton, user }) {
  return (
    <div className="bg-gradient-to-b from-[#f1f9ff] via-white to-white">
      <section className="mx-auto flex max-w-6xl flex-col gap-10 px-6 pb-16 pt-12 lg:flex-row lg:items-center">
        <div className="flex-1">
          <div className="inline-flex items-center gap-2 rounded-full bg-brand-100 px-4 py-2 text-sm font-semibold text-brand-700">
            🚀 Более 500 студентов уже с нами
          </div>
          <h1 className="mt-6 text-4xl font-extrabold leading-tight text-slate-900 sm:text-5xl">
            Освой Golang<br />
            <span className="text-brand-500">от A до Z</span>
          </h1>
          <p className="mt-4 max-w-xl text-base text-slate-600 sm:text-lg">
            Структурированные курсы, практические задания и тесты. Все что нужно для старта
            успешной карьеры Go-разработчика.
          </p>
          <div className="mt-6 flex flex-wrap gap-4">
            {showTestButton ? (
              user ? (
                <button
                  className="inline-flex items-center gap-2 rounded-xl bg-brand-500 px-6 py-3 text-sm font-semibold text-white shadow-glow transition hover:bg-brand-600"
                  onClick={onStartTest}
                >
                  Начать обучение →
                </button>
              ) : (
                <Link
                  className="inline-flex items-center gap-2 rounded-xl bg-brand-500 px-6 py-3 text-sm font-semibold text-white shadow-glow transition hover:bg-brand-600"
                  to="/auth"
                >
                  Начать обучение →
                </Link>
              )
            ) : (
              <Link
                className="inline-flex items-center gap-2 rounded-xl bg-brand-500 px-6 py-3 text-sm font-semibold text-white shadow-glow transition hover:bg-brand-600"
                to="/dashboard"
              >
                Продолжить обучение →
              </Link>
            )}
            <Link
              className="inline-flex items-center gap-2 rounded-xl border border-brand-200 bg-white px-6 py-3 text-sm font-semibold text-brand-700"
              to="/courses"
            >
              Посмотреть курсы
            </Link>
          </div>
          <div className="mt-8 flex items-center gap-4 text-sm text-slate-500">
            <div className="flex -space-x-2">
              {[1, 2, 3, 4, 5].map((i) => (
                <div key={i} className="h-8 w-8 rounded-full border-2 border-white bg-brand-100" />
              ))}
            </div>
            <span>500+ студентов</span>
            <div className="ml-4 flex items-center gap-1 text-brand-500">
              {Array.from({ length: 5 }).map((_, i) => (
                <span key={i}>★</span>
              ))}
              <span className="ml-2 text-slate-500">4.9/5</span>
            </div>
          </div>
        </div>
        <div className="flex-1">
          <div className="relative mx-auto h-64 w-full max-w-xl rounded-[32px] bg-gradient-to-r from-brand-200 via-brand-300 to-brand-200 shadow-soft">
            <div className="absolute -right-6 top-10 h-40 w-40 rounded-full bg-brand-300/60 blur-3xl" />
            <div className="absolute bottom-6 left-10 h-28 w-28 rounded-full bg-brand-100 blur-2xl" />
          </div>
        </div>
      </section>

      <section className="mx-auto max-w-6xl px-6 pb-14">
        <div className="text-center">
          <span className="rounded-full bg-brand-100 px-4 py-1 text-xs font-semibold text-brand-700">Преимущества</span>
          <h2 className="mt-4 text-3xl font-bold text-slate-900">Почему выбирают нас</h2>
          <p className="mt-2 text-slate-500">Мы создали идеальную среду для изучения Go</p>
        </div>
        <div className="mt-10 grid gap-6 md:grid-cols-2 lg:grid-cols-4">
          {[
            { title: "Быстрый старт", text: "Начните изучать Go с первого дня. Понятные объяснения и практика." },
            { title: "Практический подход", text: "Реальные проекты и задачи. Применяйте знания сразу." },
            { title: "Поддержка сообщества", text: "Общайтесь с другими студентами и получайте помощь." },
            { title: "Сертификат", text: "Получите сертификат после прохождения курса." }
          ].map((item) => (
            <div key={item.title} className="rounded-2xl border border-slate-100 bg-white p-6 text-center shadow-card">
              <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-2xl bg-brand-100 text-brand-600">★</div>
              <h3 className="mt-4 text-base font-semibold text-slate-900">{item.title}</h3>
              <p className="mt-2 text-sm text-slate-500">{item.text}</p>
            </div>
          ))}
        </div>
      </section>

      <section className="mx-auto max-w-6xl px-6 pb-16">
        <div className="text-center">
          <span className="rounded-full bg-brand-100 px-4 py-1 text-xs font-semibold text-brand-700">Процесс обучения</span>
          <h2 className="mt-4 text-3xl font-bold text-slate-900">Как это работает</h2>
        </div>
        <div className="mt-10 grid gap-6 md:grid-cols-3">
          {[
            { step: "01", title: "Выберите уровень", text: "Определите свой текущий уровень и начните с подходящего курса" },
            { step: "02", title: "Практикуйтесь", text: "Выполняйте практические задания и укрепляйте навыки" },
            { step: "03", title: "Проходите тесты", text: "Проверяйте знания и получайте обратную связь" }
          ].map((item, index) => (
            <div key={item.step} className="relative rounded-2xl border border-slate-100 bg-white p-6 text-center shadow-card">
              <div className="text-3xl font-extrabold text-brand-300">{item.step}</div>
              <h3 className="mt-2 text-base font-semibold text-slate-900">{item.title}</h3>
              <p className="mt-2 text-sm text-slate-500">{item.text}</p>
              {index < 2 && <div className="absolute right-0 top-1/2 hidden h-0.5 w-12 bg-brand-200 md:block" />}
            </div>
          ))}
        </div>
      </section>

      <section className="mx-auto max-w-6xl px-6 pb-16">
        <div className="text-center">
          <span className="rounded-full bg-brand-100 px-4 py-1 text-xs font-semibold text-brand-700">Наши курсы</span>
          <h2 className="mt-4 text-3xl font-bold text-slate-900">Выберите свой путь</h2>
          <p className="mt-2 text-slate-500">Каждый курс содержит видеоуроки, практические материалы и тесты</p>
        </div>
        <div className="mt-10 grid gap-6 md:grid-cols-3">
          {[
            {
              label: "Новичок",
              title: "Golang для новичков",
              text: "Начните свой путь в программировании с Go. Основы языка и первые проекты.",
              level: "base"
            },
            {
              label: "Средний",
              title: "Средний уровень Go",
              text: "Углубленное изучение Go: горутины, каналы, базы данных.",
              level: "mid"
            },
            {
              label: "Профи",
              title: "Профессиональный Go",
              text: "Микросервисы, production-ready код, оптимизация и архитектурные паттерны.",
              level: "pro"
            }
          ].map((course) => (
            <div key={course.label} className="overflow-hidden rounded-2xl border border-slate-100 bg-white shadow-card">
              <div className="h-36 bg-gradient-to-r from-brand-200 via-brand-300 to-brand-200" />
              <div className="p-6">
                <div className="inline-flex rounded-full bg-brand-100 px-3 py-1 text-xs font-semibold text-brand-700">{course.label}</div>
                <h3 className="mt-3 text-lg font-semibold text-slate-900">{course.title}</h3>
                <p className="mt-2 text-sm text-slate-500">{course.text}</p>
                <Link
                  className="mt-4 inline-flex w-full items-center justify-center rounded-xl bg-brand-500 px-4 py-2 text-sm font-semibold text-white"
                  to={`/courses?level=${course.level}`}
                >
                  Начать обучение →
                </Link>
              </div>
            </div>
          ))}
        </div>
      </section>

      <section className="mx-auto max-w-6xl px-6 pb-16">
        <div className="text-center">
          <span className="rounded-full bg-brand-100 px-4 py-1 text-xs font-semibold text-brand-700">Отзывы</span>
          <h2 className="mt-4 text-3xl font-bold text-slate-900">Что говорят студенты</h2>
        </div>
        <div className="mt-8 grid gap-6 md:grid-cols-3">
          {[
            { name: "Алексей М.", role: "Junior Developer", text: "Отличный курс для новичков! За 2 месяца получил работу Go-разработчиком." },
            { name: "Мария К.", role: "Backend Developer", text: "Курс среднего уровня помог углубить знания. Особенно понравились разделы про конкурентность." },
            { name: "Дмитрий П.", role: "Tech Lead", text: "Профессиональный уровень — именно то, что нужно для production. Рекомендую!" }
          ].map((item) => (
            <div key={item.name} className="rounded-2xl border border-slate-100 bg-white p-6 shadow-card">
              <div className="text-brand-500">★★★★★</div>
              <p className="mt-3 text-sm text-slate-600">"{item.text}"</p>
              <div className="mt-4 flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-full bg-brand-100 text-brand-600">
                  {item.name[0]}
                </div>
                <div>
                  <div className="text-sm font-semibold text-slate-900">{item.name}</div>
                  <div className="text-xs text-slate-500">{item.role}</div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </section>

      <section className="bg-gradient-to-r from-brand-500 to-brand-300 py-10">
        <div className="mx-auto max-w-6xl px-6">
          <div className="grid gap-6 md:grid-cols-4">
            {[
              { value: "500+", label: "Активных студентов" },
              { value: "96", label: "Видео уроков" },
              { value: "4.9", label: "Средний рейтинг" },
              { value: "89%", label: "Нашли работу" }
            ].map((item) => (
              <div key={item.label} className="text-center text-white">
                <div className="text-3xl font-bold">{item.value}</div>
                <div className="mt-1 text-sm text-white/80">{item.label}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="mx-auto max-w-6xl px-6 py-14">
        <div className="rounded-3xl bg-gradient-to-r from-brand-500 to-brand-300 px-8 py-10 text-center text-white shadow-soft">
          <h2 className="text-3xl font-bold">Готовы начать обучение?</h2>
          <p className="mt-2 text-white/90">Присоединяйтесь к нашему сообществу и станьте профессиональным Go-разработчиком</p>
          <div className="mt-6 flex flex-wrap justify-center gap-4">
            <Link className="rounded-xl bg-white px-6 py-3 text-sm font-semibold text-brand-700" to="/auth">
              Зарегистрироваться бесплатно
            </Link>
            <Link className="rounded-xl border border-white/40 px-6 py-3 text-sm font-semibold text-white" to="/courses">
              Посмотреть курсы
            </Link>
          </div>
        </div>
      </section>

      <footer className="border-t border-slate-100 bg-white">
        <div className="mx-auto grid max-w-6xl grid-cols-1 gap-8 px-6 py-10 md:grid-cols-4">
          <div>
            <div className="text-lg font-semibold text-brand-600">GoLearn</div>
            <p className="mt-2 text-sm text-slate-500">Профессиональная платформа для изучения Golang.</p>
          </div>
          <div>
            <div className="text-sm font-semibold text-slate-900">Курсы</div>
            <ul className="mt-3 space-y-2 text-sm text-slate-500">
              <li>Для новичков</li>
              <li>Средний уровень</li>
              <li>Профессиональный</li>
            </ul>
          </div>
          <div>
            <div className="text-sm font-semibold text-slate-900">О платформе</div>
            <ul className="mt-3 space-y-2 text-sm text-slate-500">
              <li><Link to="/about" className="hover:text-brand-600">О нас</Link></li>
              <li>Преподаватели</li>
              <li>Отзывы</li>
              <li><Link to="/contacts" className="hover:text-brand-600">Контакты</Link></li>
            </ul>
          </div>
          <div>
            <div className="text-sm font-semibold text-slate-900">Поддержка</div>
            <ul className="mt-3 space-y-2 text-sm text-slate-500">
              <li>Помощь</li>
              <li>FAQ</li>
              <li>Политика</li>
              <li>Условия</li>
            </ul>
          </div>
        </div>
      </footer>
    </div>
  );
}
