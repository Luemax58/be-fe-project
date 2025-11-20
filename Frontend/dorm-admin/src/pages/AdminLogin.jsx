import React, { useState } from "react";
import { useNavigate } from "react-router-dom";

export default function LoginAdmin() {
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = (e) => {
    e.preventDefault();

    // Simple hardcoded login
    if (username === "admin" && password === "1234") {
      navigate("/admin/dashboard");
    } else {
      alert("ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง");
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">

      <div className="bg-white p-10 rounded-xl shadow-lg w-96">
        <h1 className="text-2xl font-bold mb-6 text-center">เข้าสู่ระบบเจ้าของหอ</h1>

        <form onSubmit={handleLogin} className="space-y-4">
          <div>
            <label className="text-gray-700">Username</label>
            <input
              type="text"
              className="w-full px-4 py-2 border rounded-lg mt-1"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </div>

          <div>
            <label className="text-gray-700">Password</label>
            <input
              type="password"
              className="w-full px-4 py-2 border rounded-lg mt-1"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </div>

          <button
            type="submit"
            className="w-full py-2 bg-viridian-600 text-white rounded-lg hover:bg-viridian-700"
          >
            เข้าสู่ระบบ
          </button>
        </form>

        <button
          onClick={() => navigate("/")}
          className="w-full mt-4 py-2 bg-gray-200 text-gray-900 rounded-lg hover:bg-gray-300"
        >
          กลับหน้าแรก
        </button>
      </div>
    </div>
  );
}
