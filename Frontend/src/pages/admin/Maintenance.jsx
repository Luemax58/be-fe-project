import { useEffect, useState } from "react";
import axiosClient from "../../api/axiosClient";
import { formatThaiDateTime } from "../../utils/dateFormat";
import AdminLayout from "../../layouts/AdminLayout";

export default function AdminMaintenance() {
  const [requests, setRequests] = useState([]);
  const [filter, setFilter] = useState("all");
  const [selectedReq, setSelectedReq] = useState(null);

  useEffect(() => {
    loadRequests();
  }, []);

  const loadRequests = async () => {
    try {
      const res = await axiosClient.get("/maintenance/requests");
      setRequests(res.data.data || []);
    } catch (err) {
      console.error(err);
      alert("โหลดข้อมูลคำร้องซ่อมไม่สำเร็จ");
    }
  };

  const updateStatus = async (status) => {
    try {
      await axiosClient.put(`/maintenance/update/${selectedReq.request_id}`, {
        status,
      });

      alert("อัปเดตสถานะสำเร็จ");
      setSelectedReq(null);
      loadRequests();
    } catch (err) {
      console.error(err);
      alert("อัปเดตสถานะไม่สำเร็จ");
    }
  };

  const filteredRequests =
    filter === "all"
      ? requests
      : requests.filter((r) => r.status === filter);

  return (
    <AdminLayout>
    <div className="p-5">
      <h1 className="text-2xl font-bold mb-5">คำร้องซ่อมบำรุง (Admin)</h1>

      {/* ตัวกรอง */}
      <div className="flex gap-4 mb-4">
        {["all", "pending", "in_progress", "completed"].map((s) => (
          <button
            key={s}
            onClick={() => setFilter(s)}
            className={`px-4 py-2 rounded ${
              filter === s
                ? "bg-blue-600 text-white"
                : "bg-gray-200 text-gray-700"
            }`}
          >
            {s === "all"
              ? "ทั้งหมด"
              : s === "pending"
              ? "รอดำเนินการ"
              : s === "in_progress"
              ? "กำลังดำเนินการ"
              : "เสร็จสิ้น"}
          </button>
        ))}
      </div>

      {/* ตาราง */}
      <table className="w-full border text-left border-collapse">
        <thead>
          <tr className="bg-gray-200">
            <th className="border p-2">Room</th>
            <th className="border p-2">Issue</th>
            <th className="border p-2">Status</th>
            <th className="border p-2">Date</th>
            <th className="border p-2">Action</th>
          </tr>
        </thead>

        <tbody>
          {filteredRequests.map((r) => (
            <tr key={r.request_id} className="border-t">
              <td className="border p-2">{r.room_id}</td>
              <td className="border p-2">{r.issue_description}</td>

              {/* สีสถานะ */}
              <td className="border p-2">
                <span
                  className={
                    r.status === "pending"
                      ? "text-yellow-600 font-semibold"
                      : r.status === "in_progress"
                      ? "text-blue-600 font-semibold"
                      : "text-green-600 font-semibold"
                  }
                >
                  {r.status}
                </span>
              </td>

              <td className="border p-2">{formatThaiDateTime(r.request_date)}</td>

              <td className="border p-2">
                <button
                  onClick={() => setSelectedReq(r)}
                  className="bg-blue-600 text-white px-3 py-1 rounded"
                >
                  จัดการ
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* Modal อัปเดตสถานะ */}
      {selectedReq && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center">
          <div className="bg-white p-5 rounded-lg shadow w-96">

            <h2 className="text-lg font-bold mb-3">
              อัปเดตสถานะคำร้อง #{selectedReq.request_id}
            </h2>

            <div className="flex flex-col gap-2">
              <button
                className="bg-yellow-500 text-white px-4 py-2 rounded"
                onClick={() => updateStatus("pending")}
              >
                รอดำเนินการ
              </button>

              <button
                className="bg-blue-600 text-white px-4 py-2 rounded"
                onClick={() => updateStatus("in_progress")}
              >
                กำลังดำเนินการ
              </button>

              <button
                className="bg-green-600 text-white px-4 py-2 rounded"
                onClick={() => updateStatus("completed")}
              >
                เสร็จสิ้น
              </button>

              <button
                className="mt-2 text-gray-600"
                onClick={() => setSelectedReq(null)}
              >
                ยกเลิก
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
    </AdminLayout>
  );
}
