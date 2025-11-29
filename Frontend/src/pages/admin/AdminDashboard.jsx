import React, { useEffect, useState } from "react";
import axiosClient from "../../api/axiosClient";

function StatCard({ title, value, color = "text-gray-900" }) {
  return (
    <div className="bg-white rounded-lg shadow p-4">
      <div className="text-sm text-gray-500">{title}</div>
      <div className={`text-2xl font-bold ${color}`}>{value}</div>
    </div>
  );
}

export default function AdminDashboard() {
  const [rooms, setRooms] = useState([]);
  const [invoices, setInvoices] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadData = async () => {
      try {
        const token = localStorage.getItem("token");

        // ดึงห้องทั้งหมด
        const roomRes = await axiosClient.get("/rooms", {
          headers: { Authorization: `Bearer ${token}` },
        });

        setRooms(roomRes.data);

        // TODO: BE ยังไม่มี API สำหรับ invoices → ต้องเพิ่มเองภายหลัง
        setInvoices([]);
      } catch (err) {
        console.error("โหลดข้อมูลผิดพลาด", err);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, []);

  if (loading) return <div>กำลังโหลด...</div>;

  return (
    <div className="space-y-6">

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <StatCard title="ห้องทั้งหมด" value={rooms.length} />
        <StatCard title="ห้องที่ว่าง" value={rooms.filter((r) => r.status === "available").length} />
        <StatCard title="ห้องที่มีผู้เช่า" value={rooms.filter((r) => r.status === "occupied").length} color="text-green-600" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div className="bg-white rounded-lg shadow p-4">
          <div className="text-sm font-medium text-gray-600 mb-3">
            รายรับ 6 เดือนย้อนหลัง (กำลังพัฒนา)
          </div>
          <div className="h-40 flex items-center justify-center text-gray-400">
            [ยังไม่มี API]
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-4">
          <div className="text-sm font-medium text-gray-600 mb-3">
            ใบแจ้งหนี้ล่าสุด (กำลังพัฒนา)
          </div>

          <div className="h-20 flex items-center justify-center text-gray-400">
            [BE ยังไม่มี API ใบแจ้งหนี้]
          </div>
        </div>
      </div>

    </div>
  );
}
