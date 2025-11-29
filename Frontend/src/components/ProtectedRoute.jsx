// src/components/ProtectedRoute.jsx
import { Navigate } from "react-router-dom";

export default function ProtectedRoute({ children, role }) {
  const token = localStorage.getItem("token");
  const userRole = localStorage.getItem("role");

  if (!token || !userRole) {
    return <Navigate to="/" replace />;
  }

  if (role && userRole !== role) {
    if (userRole === "owner") return <Navigate to="/admin/dashboard" replace />;
    if (userRole === "tenant") return <Navigate to="/tenant/dashboard" replace />;
    return <Navigate to="/" replace />;
  }

  return children;
}
