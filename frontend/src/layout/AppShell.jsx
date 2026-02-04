import { Link } from "react-router-dom";

export default function AppShell({ user, loading, onLogout, children }) {
  return (
    <div className="min-h-screen bg-[#f6fbff] text-slate-900">
      <header className="border-b border-slate-100 bg-white/80 backdrop-blur">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
          <div className="text-lg font-semibold text-brand-600">GoLearn</div>
          <nav className="hidden items-center gap-6 text-sm font-medium text-slate-600 md:flex">
            <Link className="hover:text-brand-600" to="/">Главная</Link>
            {user && <Link className="hover:text-brand-600" to="/courses">Курсы</Link>}
            {user && <Link className="hover:text-brand-600" to="/dashboard">Личный кабинет</Link>}
            {user?.is_admin && <Link className="hover:text-brand-600" to="/admin">Админ</Link>}
          </nav>
          <div className="flex items-center gap-4 text-sm">
            {user ? (
              <>
                <span className="hidden text-slate-500 sm:inline">Привет, {user.name}</span>
                <button className="rounded-xl border border-slate-200 px-4 py-2 text-sm font-semibold text-slate-600" onClick={onLogout}>
                  Выйти
                </button>
              </>
            ) : (
              <Link className="rounded-xl bg-brand-500 px-4 py-2 text-sm font-semibold text-white" to="/auth">
                Войти
              </Link>
            )}
          </div>
        </div>
      </header>
      <main>{loading ? <div className="page">Загрузка...</div> : children}</main>
    </div>
  );
}
