// src/components/TenantLayout.jsx
import { Outlet } from "react-router-dom";
import TenantSidebar from "./TenantSidebar";
import TenantTopbar from "./TenantTopbar";

export default function TenantLayout() {
  return (
    <div className="flex min-h-screen bg-gray-100">
      <TenantSidebar />
      <div className="flex-1 flex flex-col">
        <TenantTopbar />
        <main className="p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
