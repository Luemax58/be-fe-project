import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";

import Login from "./pages/Login";

// Admin pages
import AdminDashboard from "./pages/admin/Dashboard";
import Rooms from "./pages/admin/Rooms";
import Billing from "./pages/admin/Billing";
import AdminMaintenance from "./pages/admin/Maintenance";

// Tenant pages
import TenantDashboard from "./pages/tenant/Dashboard";
import TenantBilling from "./pages/tenant/Billing";
import TenantMaintenance from "./pages/tenant/Maintenance";

// Guards
import AdminRoute from "./routes/AdminRoute";
import TenantRoute from "./routes/TenantRoute";

export default function App() {
  return (
    <BrowserRouter>
      <Routes>

        {/* Redirect "/" -> "/login" */}
        <Route path="/" element={<Navigate to="/login" />} />

        {/* Public */}
        <Route path="/login" element={<Login />} />

        {/* Admin Routes */}
        <Route element={<AdminRoute />}>
          <Route path="/admin/dashboard" element={<AdminDashboard />} />
          <Route path="/admin/rooms" element={<Rooms />} />
          <Route path="/admin/billing" element={<Billing />} />
          <Route path="/admin/maintenance" element={<AdminMaintenance />} />
        </Route>

        {/* Tenant Routes */}
        <Route element={<TenantRoute />}>
          <Route path="/tenant/dashboard" element={<TenantDashboard />} />
          <Route path="/tenant/billing" element={<TenantBilling />} />
          <Route path="/tenant/maintenance" element={<TenantMaintenance />} />
        </Route>

      </Routes>
    </BrowserRouter>
  );
}
