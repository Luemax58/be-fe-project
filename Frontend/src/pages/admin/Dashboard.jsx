import { useEffect, useState } from "react";
import axiosClient from "../../api/axiosClient";
import { formatThaiDateTime } from "../../utils/dateFormat";
import AdminLayout from "../../layouts/AdminLayout";

export default function AdminDashboard() {
  const [rooms, setRooms] = useState(null);
  const [maintenance, setMaintenance] = useState(null);

  useEffect(() => {
    fetchDashData();
  }, []);

  const fetchDashData = async () => {
    try {
      // ดึงข้อมูลห้อง
      const roomRes = await axiosClient.get("/rooms");
      setRooms(roomRes.data || []);

      // ดึงคำร้องซ่อม
      const maintRes = await axiosClient.get("/maintenance/requests");
      setMaintenance(maintRes.data?.data || []);

    } catch (error) {
      console.error(error);
      setRooms([]);
      setMaintenance([]);
    }
  };

  if (!rooms || !maintenance) {
    return <div className="p-5">Loading Dashboard...</div>;
  }

  const availableRooms = rooms.filter((r) => r.status === "available").length;
  const occupiedRooms = rooms.filter((r) => r.status === "occupied").length;
  const pendingRequests = maintenance.filter((m) => m.status === "pending").length;

  return (
    <AdminLayout>
    <div className="p-5">
      <h1 className="text-2xl font-bold mb-5">Admin Dashboard</h1>

      <div className="grid grid-cols-3 gap-4 mb-8">
        <div className="bg-blue-200 p-4 rounded-lg">
          <h2 className="text-lg font-semibold">ห้องทั้งหมด</h2>
          <p className="text-3xl">{rooms.length}</p>
        </div>

        <div className="bg-green-200 p-4 rounded-lg">
          <h2 className="text-lg font-semibold">ห้องมีผู้เช่า</h2>
          <p className="text-3xl">{occupiedRooms}</p>
        </div>

        <div className="bg-yellow-200 p-4 rounded-lg">
          <h2 className="text-lg font-semibold">คำร้องซ่อมที่รอดำเนินการ</h2>
          <p className="text-3xl">{pendingRequests}</p>
        </div>
      </div>

      <h2 className="text-xl font-semibold mb-3">คำร้องซ่อมล่าสุด</h2>
      <table className="w-full border text-left">
        <thead>
          <tr className="bg-gray-200">
            <th className="p-2">Room</th>
            <th className="p-2">Description</th>
            <th className="p-2">Status</th>
            <th className="p-2">Date</th>
          </tr>
        </thead>
        <tbody>
          {maintenance.slice(0, 5).map((m, i) => (
            <tr key={i} className="border-t">
              <td className="p-2">{m.room_id}</td>
              <td className="p-2">{m.issue_description}</td>
              <td className="p-2">{m.status}</td>
              <td className="p-2">{formatThaiDateTime(m.request_date)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
    </AdminLayout>
  );
}
