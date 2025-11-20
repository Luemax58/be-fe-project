import React from "react";
import { useNavigate } from "react-router-dom";

export default function Login() {
  const nav = useNavigate();
  const login = (e) => {
    e.preventDefault();
    // mock login -> redirect to admin
    nav("/admin/dashboard");
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <form className="w-full max-w-md bg-white p-8 rounded-lg shadow" onSubmit={login}>
        <h2 className="text-2xl font-semibold mb-4">เข้าสู่ระบบ</h2>
        <div className="mb-3">
          <label className="text-sm text-gray-600">Email</label>
          <input className="input-field mt-1" defaultValue="admin@example.com" />
        </div>
        <div className="mb-4">
          <label className="text-sm text-gray-600">Password</label>
          <input className="input-field mt-1" type="password" defaultValue="password" />
        </div>
        <button className="btn-primary w-full">Login</button>
      </form>
    </div>
  );
}
