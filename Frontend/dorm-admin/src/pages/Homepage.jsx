import React from "react";
import { useNavigate } from "react-router-dom";

export default function Homepage() {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50 text-center px-6">
      <h1 className="text-4xl font-bold mb-4 text-viridian-700">
        ระบบจัดการหอพัก
      </h1>

      <p className="text-gray-600 mb-8 text-lg">
        โปรดเลือกเมนูเพื่อเข้าสู่ระบบ
      </p>

      <button
        onClick={() => navigate("/admin/login")}
        className="px-6 py-3 bg-viridian-600 text-white rounded-lg text-lg font-medium
                   hover:bg-viridian-700 transition shadow"
      >
        เข้าสู่ระบบเจ้าของหอ
      </button>
    </div>
  );
}
