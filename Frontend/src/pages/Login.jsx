import { useState } from "react";
import { useNavigate } from "react-router-dom";
import authApi from "../api/authApi";

export default function Login() {
  const navigate = useNavigate();
  const [form, setForm] = useState({ username: "", password: "" });
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState("");

  const handleChange = (e) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleLogin = async () => {
    setErrorMsg("");

    if (!form.username || !form.password) {
      return setErrorMsg("กรุณากรอกชื่อผู้ใช้และรหัสผ่าน");
    }

    try {
      setLoading(true);

      // 🔥 1) ส่ง username + password ไปยัง Backend
      const res = await authApi.login(form);

      const token = res.data.token;
      if (!token) {
        setErrorMsg("ไม่พบ token จาก server");
        return;
      }

      localStorage.setItem("token", token);

      // 🔥 2) ดึงข้อมูลผู้ใช้
      const profile = await authApi.getProfile();
      const role = profile.data.role;

      if (!role) {
        setErrorMsg("ไม่พบ role ของผู้ใช้");
        return;
      }

      localStorage.setItem("role", role);

      // 🔥 3) Redirect ตาม role
      if (role === "owner") navigate("/admin/dashboard");
      if (role === "tenant") navigate("/tenant/dashboard");

    } catch (err) {
      console.error(err);
      setErrorMsg("เข้าสู่ระบบล้มเหลว โปรดตรวจสอบข้อมูลอีกครั้ง");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 p-4">
      <div className="bg-white shadow-xl rounded-2xl p-8 w-full max-w-md">

        <h1 className="text-3xl font-bold text-center mb-6">เข้าสู่ระบบ</h1>

        {/* Error Box */}
        {errorMsg && (
          <div className="bg-red-100 border border-red-300 text-red-700 text-sm p-3 mb-4 rounded">
            {errorMsg}
          </div>
        )}

        {/* Username */}
        <div className="mb-4">
          <label className="block text-gray-700 mb-1">ชื่อผู้ใช้</label>
          <input
            name="username"
            className="w-full border p-2 rounded-lg focus:ring-2 focus:ring-blue-400 outline-none"
            placeholder="กรอกชื่อผู้ใช้"
            value={form.username}
            onChange={handleChange}
          />
        </div>

        {/* Password */}
        <div className="mb-6">
          <label className="block text-gray-700 mb-1">รหัสผ่าน</label>
          <input
            name="password"
            type="password"
            className="w-full border p-2 rounded-lg focus:ring-2 focus:ring-blue-400 outline-none"
            placeholder="กรอกรหัสผ่าน"
            value={form.password}
            onChange={handleChange}
          />
        </div>

        {/* Login Button */}
        <button
          onClick={handleLogin}
          disabled={loading}
          className={`w-full py-2 rounded-lg text-white font-semibold transition ${
            loading ? "bg-blue-300" : "bg-blue-600 hover:bg-blue-700"
          }`}
        >
          {loading ? "กำลังเข้าสู่ระบบ..." : "เข้าสู่ระบบ"}
        </button>
      </div>
    </div>
  );
}
