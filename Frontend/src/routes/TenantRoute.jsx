import { Navigate, Outlet } from "react-router-dom";

export default function TenantRoute() {
  const token = localStorage.getItem("token");
  const role = localStorage.getItem("role");

  console.log("TENANT ROUTE:", token, role);

  // ถ้าไม่มี token หรือ role ไม่ใช่ tenant → ไป login
  if (!token || role !== "tenant") {
    return <Navigate to="/login" replace />;
  }

  // ถ้าต้องการ loading ขณะรอข้อมูล role (optional)
  if (role === null || role === undefined) {
    return <div>Loading...</div>;
  }

  // ถ้า role เป็น tenant → ให้ Router แสดง nested routes ผ่าน <Outlet/>
  if (role === "tenant") {
    return <Outlet />;
  }

  if (role === "owner") {
    return <Navigate to="/admin/dashboard" />;
  }

  return <Navigate to="/login" />;
}
