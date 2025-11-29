import React from "react";
import { Outlet } from "react-router-dom";
import AdminSidebar from "./AdminSidebar";

function Topbar({ title = "Dorm Admin" }) {
  return (
    <header className="flex items-center justify-between p-4 border-b border-gray-200 bg-white">
      <div className="flex items-center gap-4">
        <h1 className="text-xl font-semibold">{title}</h1>
        <div className="text-sm text-gray-500">ระบบจัดการหอพัก</div>
      </div>
      <div className="flex items-center gap-3">
        <button className="btn-secondary">Settings</button>
        <div className="text-sm text-gray-700">Admin</div>
      </div>
    </header>
  );
}

export default function AdminLayout() {
  return (
    <div className="min-h-screen flex bg-gray-50">
      <AdminSidebar />
      <div className="flex-1 flex flex-col">
        <Topbar />
        <main className="flex-1 overflow-auto p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
