import React, { useEffect, useState } from "react";
import axiosClient from "../../api/axiosClient";

export default function Rooms() {
  const [rooms, setRooms] = useState([]);

  useEffect(() => {
    const loadRooms = async () => {
      try {
        const token = localStorage.getItem("token");

        const res = await axiosClient.get("/rooms", {
          headers: { Authorization: `Bearer ${token}` },
        });

        setRooms(res.data);
      } catch (err) {
        console.error("โหลดห้องผิดพลาด", err);
      }
    };

    loadRooms();
  }, []);

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
              <th>สถานะ</th>
            </tr>
          </thead>

          <tbody>
            {rooms.map((room) => (
              <tr key={room.room_id} className="odd:bg-white even:bg-gray-50">
                <td className="py-2">{room.room_number}</td>
                <td>{room.status === "available" ? "ว่าง" : "มีผู้เช่า"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
