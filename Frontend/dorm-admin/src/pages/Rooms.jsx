import React from "react";
import { mockRooms } from "../mockData";

export default function Rooms() {
  return (
    <div>
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold">จัดการห้องพัก</h2>
        <button className="px-4 py-2 rounded bg-viridian-600 text-white">
          + เพิ่มห้องพัก
        </button>
      </div>

      <div className="bg-white rounded-lg shadow p-4">
        <table className="w-full text-left">
          <thead>
            <tr className="text-sm text-gray-500 border-b">
              <th className="py-2">เลขห้อง</th>
              <th>ชั้น</th>
              <th>ประเภท</th>
              <th>ราคา</th>
              <th>สถานะ</th>
              <th>จัดการ</th>
            </tr>
          </thead>
          <tbody>
            {mockRooms.map((r) => (
              <tr key={r.id} className="odd:bg-white even:bg-gray-50">
                <td className="py-2">{r.number}</td>
                <td>{r.floor}</td>
                <td>{r.type}</td>
                <td>{r.price} บาท</td>
                <td>
                  {r.status === "occupied" ? (
                    <span className="text-yellow-600">มีคนเช่า</span>
                  ) : (
                    <span className="text-gray-600">ว่าง</span>
                  )}
                </td>
                <td>
                  <button className="text-blue-600 mr-2">แก้ไข</button>
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
