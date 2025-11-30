import { Navigate, Outlet } from "react-router-dom";

export default function AdminRoute() {
  const token = localStorage.getItem("token");
  const role = localStorage.getItem("role");

  if (!token) return <Navigate to="/login" />;

  // ถ้า role ไม่ใช่ owner → เด้งออก
  if (role !== "owner") return <Navigate to="/login" />;

  // ถ้า role เป็น owner → ผ่าน
  return <Outlet />;
}
