// src/components/TenantSidebar.jsx
import { NavLink, useNavigate } from "react-router-dom";

export default function TenantSidebar() {
  const nav = useNavigate();

  const logout = () => {
    localStorage.clear();
    nav("/");
  };

  return (
    <aside className="w-64 bg-white min-h-screen border-r p-4">
      <h2 className="text-xl font-bold mb-6">Dorm User</h2>

      <nav className="space-y-2">
        <NavLink to="/tenant/dashboard" className="block px-3 py-2 rounded hover:bg-gray-200">Dashboard</NavLink>
        <NavLink to="/tenant/room" className="block px-3 py-2 rounded hover:bg-gray-200">ห้องพักของฉัน</NavLink>
        <NavLink to="/tenant/maintenance/create" className="block px-3 py-2 rounded hover:bg-gray-200">แจ้งซ่อม</NavLink>
        <NavLink to="/tenant/maintenance" className="block px-3 py-2 rounded hover:bg-gray-200">ประวัติแจ้งซ่อม</NavLink>
        <NavLink to="/tenant/payment" className="block px-3 py-2 rounded hover:bg-gray-200">ชำระเงิน</NavLink>
      </nav>

      <button
        onClick={logout}
        className="mt-6 bg-red-500 text-white w-full py-2 rounded"
      >
        ออกจากระบบ
      </button>
    </aside>
  );
}
