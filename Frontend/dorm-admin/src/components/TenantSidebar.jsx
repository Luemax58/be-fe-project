import React from "react";
import { NavLink } from "react-router-dom";

export default function TenantSidebar() {
  return (
    <aside className="w-64 bg-white border-r border-gray-200 p-4 min-h-screen">
      <h3 className="font-semibold mb-2">ผู้เช่า</h3>
      <nav className="space-y-1">
        <NavLink to="/tenant/room" className="block px-3 py-2 rounded hover:bg-gray-100">My Room</NavLink>
        <NavLink to="/tenant/invoices" className="block px-3 py-2 rounded hover:bg-gray-100">Invoices</NavLink>
        <NavLink to="/tenant/profile" className="block px-3 py-2 rounded hover:bg-gray-100">Profile</NavLink>
      </nav>
    </aside>
  );
}
