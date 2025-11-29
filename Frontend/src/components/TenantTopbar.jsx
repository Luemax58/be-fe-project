// src/components/TenantTopbar.jsx
import { useNavigate } from "react-router-dom";

export default function TenantTopbar() {
  const nav = useNavigate();
  const logout = () => {
    localStorage.clear();
    nav("/");
  };

  return (
    <header className="bg-white border-b h-16 flex items-center justify-between px-6 shadow-sm">
      <h1 className="text-xl font-semibold">ระบบผู้เช่าหอพัก</h1>

      <button
        onClick={logout}
        className="bg-red-500 text-white px-3 py-1 rounded"
      >
        ออกจากระบบ
      </button>
    </header>
  );
}
