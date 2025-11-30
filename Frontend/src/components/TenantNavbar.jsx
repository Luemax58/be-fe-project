import { NavLink, useNavigate } from "react-router-dom";

export default function TenantNavbar() {
  const navigate = useNavigate();
  const role = localStorage.getItem("role");

  const logout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("role");
    navigate("/login");
  };

  const linkClass = ({ isActive }) =>
    `px-4 py-2 rounded-lg transition ${
      isActive
        ? "bg-blue-600 text-white font-semibold shadow"
        : "text-gray-700 hover:bg-gray-100"
    }`;

  return (
    <nav className="w-full bg-white shadow p-4 flex items-center justify-between">
      <h1 className="text-xl font-bold text-blue-700">Ovaltin System</h1>

      <div className="flex space-x-2">
        <NavLink to="/tenant/dashboard" className={linkClass}>แดชบอร์ด</NavLink>
        <NavLink to="/tenant/billing" className={linkClass}>บิลของฉัน</NavLink>
        <NavLink to="/tenant/maintenance" className={linkClass}>แจ้งซ่อม</NavLink>
      </div>

      <div className="flex items-center space-x-4">
        <span className="text-gray-600">Tenant</span>
        <button
          onClick={logout}
          className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600"
        >
          ออกจากระบบ
        </button>
      </div>
    </nav>
  );
}
