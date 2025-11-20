// src/App.js
import { BrowserRouter, Routes, Route } from "react-router-dom";

import Homepage from "./pages/Homepage";
import AdminLayout from "./components/AdminLayout";

import AdminDashboard from "./pages/AdminDashboard";
import Rooms from "./pages/Rooms";
import Tenants from "./pages/Tenants";
import Invoices from "./pages/Invoices";
import Reports from "./pages/Reports";
import AdminLogin from "./pages/AdminLogin";

function App() {
  return (
    <BrowserRouter>
      <Routes>

        {/* หน้าแรกของเว็บ */}
        <Route path="/" element={<Homepage />} />

        <Route path="/admin/login" element={<AdminLogin />} />
        
        {/* เส้นทางของ Admin ที่มี Sidebar */}
        <Route path="/admin" element={<AdminLayout />}>
          <Route path="dashboard" element={<AdminDashboard />} />
          <Route path="rooms" element={<Rooms />} />
          <Route path="tenants" element={<Tenants />} />
          <Route path="invoices" element={<Invoices />} />
          <Route path="reports" element={<Reports />} />
        </Route>

      </Routes>
    </BrowserRouter>
  );
}

export default App;
