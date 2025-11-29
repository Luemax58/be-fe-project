// src/App.js
import { BrowserRouter, Routes, Route } from "react-router-dom";

// Login
import Login from "./pages/Login";

// Admin layouts + pages
import AdminLayout from "./components/AdminLayout";
import AdminDashboard from "./pages/admin/AdminDashboard";
import Rooms from "./pages/admin/Rooms";
import Tenants from "./pages/admin/Tenants";
import Invoices from "./pages/admin/Invoices";
import Reports from "./pages/admin/Reports";

// Tenant layouts + pages
import TenantLayout from "./components/TenantLayout";
import Dashboard from "./pages/tenant/Dashboard";
import RoomDetail from "./pages/tenant/RoomDetail";
import MaintenanceList from "./pages/tenant/MaintenanceList";
import MaintenanceCreate from "./pages/tenant/MaintenanceCreate";
import PaymentForm from "./pages/tenant/PaymentForm";

// Protected route
import ProtectedRoute from "./components/ProtectedRoute";

export default function App() {
  return (
    <BrowserRouter>
      <Routes>

        {/* หน้า Login รวมทุกฝ่าย */}
        <Route path="/" element={<Login />} />

        {/* ADMIN ROUTES */}
        <Route
          path="/admin"
          element={
            <ProtectedRoute role="owner">
              <AdminLayout />
            </ProtectedRoute>
          }
        >
          <Route path="dashboard" element={<AdminDashboard />} />
          <Route path="rooms" element={<Rooms />} />
          <Route path="tenants" element={<Tenants />} />
          <Route path="invoices" element={<Invoices />} />
          <Route path="reports" element={<Reports />} />
        </Route>

        {/* TENANT ROUTES */}
        <Route
          path="/tenant"
          element={
            <ProtectedRoute role="tenant">
              <TenantLayout />
            </ProtectedRoute>
          }
        >
          <Route path="dashboard" element={<Dashboard />} />
          <Route path="room" element={<RoomDetail />} />
          <Route path="maintenance" element={<MaintenanceList />} />
          <Route path="maintenance/create" element={<MaintenanceCreate />} />
          <Route path="payment" element={<PaymentForm />} />
        </Route>

      </Routes>
    </BrowserRouter>
  );
}
