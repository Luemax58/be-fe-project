import React from "react";
import { mockTenants } from "../mockData";

export default function Tenants() {
  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold">จัดการผู้เช่า</h2>
        <button className="px-4 py-2 rounded bg-viridian-600 text-white">
          + เพิ่มผู้เช่า
        </button>
      </div>

      <div className="bg-white rounded-lg shadow p-4">
        <table className="w-full text-left">
          <thead>
            <tr className="text-sm text-gray-500 border-b">
              <th className="py-2">ชื่อ</th>
              <th>เบอร์</th>
              <th>ห้อง</th>
              <th>วันเข้าอยู่</th>
              <th>จัดการ</th>
            </tr>
          </thead>
          <tbody>
            {mockTenants.map((t) => (
              <tr key={t.id} className="odd:bg-white even:bg-gray-50">
                <td className="py-2">{t.name}</td>
                <td>{t.phone}</td>
                <td>{t.room}</td>
                <td>{t.moveIn}</td>
                <td>
                  <button className="text-blue-600 mr-2">ดู</button>
                  <button className="text-red-600">ลบ</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
