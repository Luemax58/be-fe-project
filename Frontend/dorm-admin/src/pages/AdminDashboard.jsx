import React from "react";
import { mockRooms, mockTenants, mockInvoices } from "../mockData";

function StatCard({ title, value, color = "text-gray-900" }) {
  return (
    <div className="bg-white rounded-lg shadow p-4">
      <div className="text-sm text-gray-500">{title}</div>
      <div className={`text-2xl font-bold ${color}`}>{value}</div>
    </div>
  );
}

export default function AdminDashboard() {
  const totalTenants = mockTenants.length;
  const totalRooms = mockRooms.length;
  const unpaid = mockInvoices.filter((i) => i.status !== "paid").length;

  return (
    <div className="space-y-6">

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <StatCard title="จำนวนผู้เช่า" value={totalTenants} />
        <StatCard title="ห้องทั้งหมด" value={totalRooms} />
        <StatCard title="ค้างชำระ" value={unpaid} color="text-red-600" />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div className="bg-white rounded-lg shadow p-4">
          <div className="text-sm font-medium text-gray-600 mb-3">
            รายรับ 6 เดือนย้อนหลัง
          </div>
          <div className="h-40 flex items-center justify-center text-gray-400">
            [Mock Chart Placeholder]
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-4">
          <div className="text-sm font-medium text-gray-600 mb-3">
            ใบแจ้งหนี้ล่าสุด
          </div>

          <table className="w-full text-left">
            <thead>
              <tr className="text-sm text-gray-500 border-b">
                <th className="py-2">ผู้เช่า</th>
                <th>ห้อง</th>
                <th>ยอดรวม</th>
                <th>สถานะ</th>
              </tr>
            </thead>
            <tbody>
              {mockInvoices.map((inv) => (
                <tr key={inv.id} className="odd:bg-white even:bg-gray-50">
                  <td className="py-2">{inv.tenant}</td>
                  <td>{inv.room}</td>
                  <td>{inv.total} บาท</td>
                  <td>
                    {inv.status === "paid" ? (
                      <span className="text-green-600">จ่ายแล้ว</span>
                    ) : (
                      <span className="text-red-600">ยังไม่จ่าย</span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

      </div>
    </div>
  );
}
