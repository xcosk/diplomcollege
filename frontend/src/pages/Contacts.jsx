export default function Contacts() {
  return (
    <div className="bg-gradient-to-b from-[#f6fbff] via-white to-white">
      <div className="mx-auto max-w-4xl px-6 py-12">
        <h1 className="text-3xl font-bold text-slate-900">Контакты</h1>
        <p className="mt-4 text-base text-slate-600">
          Напишите нам, если у вас есть вопросы по обучению или работе платформы.
        </p>
        <div className="mt-6 grid gap-6 md:grid-cols-2">
          <div className="rounded-2xl border border-slate-100 bg-white p-6 shadow-card">
            <div className="text-sm text-slate-600">Email поддержки</div>
            <div className="mt-1 text-base font-semibold text-slate-900">support@golearn.io</div>
            <div className="mt-4 text-sm text-slate-600">Время ответа</div>
            <div className="mt-1 text-base font-semibold text-slate-900">Пн–Пт, 10:00–19:00</div>
          </div>
          <div className="rounded-2xl border border-slate-100 bg-white p-6 shadow-card">
            <div className="text-sm text-slate-600">Сообщество Telegram</div>
            <div className="mt-1 text-base font-semibold text-slate-900">
              <a className="text-brand-600 hover:text-brand-700" href="https://t.me/+WZb1Id_hd8RmZTUy" target="_blank" rel="noreferrer">
                Перейти в чат
              </a>
            </div>
            <div className="mt-4 text-sm text-slate-600">Поддержка и обсуждения</div>
            <div className="mt-1 text-base font-semibold text-slate-900">24/7</div>
          </div>
        </div>
        <div className="mt-8 rounded-2xl bg-gradient-to-r from-brand-500 to-brand-300 p-6 text-white shadow-soft">
          <h2 className="text-lg font-semibold">Нужна помощь быстрее?</h2>
          <p className="mt-2 text-sm text-white/90">
            В Telegram мы отвечаем быстрее и помогаем с домашними заданиями.
          </p>
        </div>
      </div>
    </div>
  );
}
