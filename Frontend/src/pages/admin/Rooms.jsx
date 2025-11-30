import { useEffect, useState } from "react";
import axiosClient from "../../api/axiosClient";
import { formatThaiDateTime } from "../../utils/dateFormat";
import AdminLayout from "../../layouts/AdminLayout";

export default function Rooms() {
  const [rooms, setRooms] = useState([]);

  useEffect(() => {
    loadRooms();
  }, []);

  const loadRooms = async () => {
    try {
      const res = await axiosClient.get("/rooms");
      setRooms(res.data || []);
    } catch (error) {
      console.error(error);
      alert("ดึงข้อมูลห้องไม่สำเร็จ");
    }
  };

  return (
    <AdminLayout>
    <div className="p-5">
      <h1 className="text-2xl font-bold mb-5">จัดการห้องพัก (Admin)</h1>

      <div className="bg-white shadow rounded-lg p-4 mb-5">
        <p>จำนวนห้องทั้งหมด: <b>{rooms.length}</b></p>
        <p>ห้องว่าง: <b>{rooms.filter((r) => r.status === "available").length}</b></p>
        <p>ห้องมีผู้เช่า: <b>{rooms.filter((r) => r.status === "occupied").length}</b></p>
      </div>

      {/* ตารางแสดงห้อง */}
      <table className="w-full border text-left border-collapse">
        <thead>
          <tr className="bg-gray-200">
            <th className="border p-2">Room ID</th>
            <th className="border p-2">Room Number</th>
            <th className="border p-2">Floor</th>
            <th className="border p-2">Status</th>
            <th className="border p-2">Tenant</th>
            <th className="border p-2">Base Rent</th>
            <th className="border p-2">Furniture Fee</th>
          </tr>
        </thead>

        <tbody>
          {rooms.map((room) => (
            <tr key={room.room_id} className="border-t">
              <td className="border p-2">{room.room_id}</td>
              <td className="border p-2">{room.room_number}</td>
              <td className="border p-2">{room.floor ?? "-"}</td>

              {/* สีตามสถานะ */}
              <td className="border p-2">
                <span
                  className={
                    room.status === "available"
                      ? "text-green-600 font-semibold"
                      : "text-red-600 font-semibold"
                  }
                >
                  {room.status}
                </span>
              </td>

              {/* ⭐ แก้ตรงนี้ */}
              <td className="border p-2">
                {room.tenant?.full_name || "-"}
              </td>

              <td className="border p-2">{room.base_rent}</td>
              <td className="border p-2">{room.furniture_fee}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
    </AdminLayout>
  );
}
