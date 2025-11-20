import React from "react";
import { mockInvoices } from "../mockData";

export default function Invoices() {
  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold">ใบแจ้งหนี้</h2>

        <div className="flex gap-2">
          <button className="px-4 py-2 rounded bg-gray-200 text-gray-700">
            Export Excel
          </button>
          <button className="px-4 py-2 rounded bg-viridian-600 text-white">
            + สร้างใบแจ้งหนี้
          </button>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow p-4">
        <table className="w-full text-left">
          <thead>
            <tr className="text-sm text-gray-500 border-b">
              <th className="py-2">ผู้เช่า</th>
              <th>ห้อง</th>
              <th>เดือน</th>
              <th>ยอดรวม</th>
              <th>สถานะ</th>
              <th>จัดการ</th>
            </tr>
          </thead>
          <tbody>
            {mockInvoices.map((inv) => (
              <tr key={inv.id} className="odd:bg-white even:bg-gray-50">
                <td className="py-2">{inv.tenant}</td>
                <td>{inv.room}</td>
                <td>{inv.month}</td>
                <td>{inv.total} บาท</td>
                <td>
                  {inv.status === "paid" ? (
                    <span className="text-green-600">จ่ายแล้ว</span>
                  ) : (
                    <span className="text-red-600">ยังไม่จ่าย</span>
                  )}
                </td>
                <td>
                  <button className="text-blue-600 mr-2">ดู</button>
                  <button className="text-green-600">ทำเครื่องหมายจ่ายแล้ว</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
