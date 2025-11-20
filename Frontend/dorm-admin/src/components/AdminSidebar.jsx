import React from "react";
import { NavLink, useNavigate } from "react-router-dom";

const NavItem = ({ to, children }) => (
  <NavLink
    to={to}
    className={({ isActive }) =>
      `block px-4 py-2 rounded-lg font-medium transition ${
        isActive
          ? "bg-viridian-600 text-white shadow"
          : "text-gray-700 hover:bg-gray-100"
      }`
    }
  >
    {children}
  </NavLink>
);

export default function AdminSidebar() {
  const navigate = useNavigate();

  const handleLogout = () => {
    // ลบ token ถ้ามีในอนาคต
    localStorage.removeItem("adminToken");

    // กลับไปหน้า homepage
    navigate("/");
  };

  return (
    <aside className="w-64 bg-white border-r border-gray-200 min-h-screen p-6 select-none">
      <h1 className="text-2xl font-bold text-viridian-700 mb-1">Dorm Admin</h1>
      <p className="text-sm text-gray-500 mb-8">ระบบจัดการเจ้าของหอพัก</p>

      <nav className="space-y-1">
        <NavItem to="/admin/dashboard">Dashboard</NavItem>
        <NavItem to="/admin/rooms">Rooms</NavItem>
        <NavItem to="/admin/tenants">Tenants</NavItem>
        <NavItem to="/admin/invoices">Invoices</NavItem>
        <NavItem to="/admin/reports">Reports</NavItem>
      </nav>

      {/* เส้นแบ่งเมนู */}
      <div className="border-t pt-6 mt-6"></div>

      {/* Add Item Button */}
      <button className="w-full bg-viridian-600 text-white py-2 rounded-lg shadow hover:bg-viridian-700 transition mb-3">
        + เพิ่มรายการ
      </button>

      {/* Logout Button */}
      <button
        onClick={handleLogout}
        className="w-full bg-red-500 text-white py-2 rounded-lg shadow hover:bg-red-600 transition"
      >
        ออกจากระบบ
      </button>
    </aside>
  );
}
