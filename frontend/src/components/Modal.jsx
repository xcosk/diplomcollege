export default function Modal({ open, title, onClose, children }) {
  if (!open) return null;
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/40 px-4">
      <div className="w-full max-w-3xl rounded-2xl bg-white shadow-card">
        <div className="flex items-center justify-between border-b border-slate-100 px-6 py-4">
          <h3 className="text-lg font-semibold text-slate-900">{title}</h3>
          <button className="rounded-full px-3 py-1 text-sm text-slate-500 hover:bg-slate-100" onClick={onClose}>
            Закрыть
          </button>
        </div>
        <div className="max-h-[70vh] overflow-y-auto px-6 py-6">
          {children}
        </div>
      </div>
    </div>
  );
}
