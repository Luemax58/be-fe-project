import { useState, useEffect } from "react";
import axiosClient from "../../api/axiosClient";
import { formatThaiDateTime } from "../../utils/dateFormat";

export default function TenantMaintenance() {
  const [requests, setRequests] = useState([]);
  const [roomID, setRoomID] = useState("");
  const [issue, setIssue] = useState("");

  useEffect(() => {
    loadRequests();
    loadProfileRoom();
  }, []);

  const loadProfileRoom = async () => {
    try {
      // โหลดโปรไฟล์ผู้ใช้
      const me = await axiosClient.get("/users/me");

      // โหลดห้องทั้งหมด
      const resRooms = await axiosClient.get("/rooms");

      // หาเฉพาะห้องที่ tenant.user_id ตรงกับ profile.user_id
      const myRoom = resRooms.data.find(
        (r) => r.tenant?.user_id === me.data.user_id
      );

      setRoomID(myRoom?.room_id || "");
    } catch (error) {
      console.error("loadProfileRoom error:", error);
    }
  };

  const loadRequests = async () => {
    const res = await axiosClient.get("/maintenance/requests");
    setRequests(res.data.data || []);
  };

  const submit = async () => {
    if (!issue) return alert("กรุณากรอกปัญหา");

    try {
      await axiosClient.post("/maintenance/creates", {
        room_id: roomID,
        issue_description: issue,
      });

      alert("ส่งคำร้องแล้ว");

      setIssue("");
      loadRequests();
    } catch (error) {
      alert(error.response?.data?.error || "ส่งคำร้องไม่สำเร็จ");
      console.error(error);
    }
  };

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-5">คำร้องซ่อมของฉัน</h1>

      {/* ฟอร์ม */}
      <div className="bg-white shadow p-4 rounded mb-6">
        <h2 className="text-xl font-semibold mb-2">ส่งคำร้องใหม่</h2>

        <textarea
          className="w-full border p-2 rounded mb-3"
          rows="3"
          placeholder="อธิบายปัญหาที่พบ"
          value={issue}
          onChange={(e) => setIssue(e.target.value)}
        />

        <button
          onClick={submit}
          className="px-4 py-2 bg-blue-600 text-white rounded"
        >
          ส่งคำร้อง
        </button>
      </div>

      {/* ตาราง */}
      <table className="w-full border text-left">
        <thead>
          <tr className="bg-gray-100">
            <th className="p-2">ปัญหา</th>
            <th className="p-2">สถานะ</th>
            <th className="p-2">วันที่</th>
          </tr>
        </thead>
        <tbody>
          {requests.map((r, idx) => (
            <tr key={idx} className="border-t">
              <td className="p-2">{r.issue_description}</td>
              <td className="p-2">
                <span
                  className={
                    r.status === "pending"
                      ? "text-yellow-600"
                      : r.status === "in_progress"
                      ? "text-blue-600"
                      : "text-green-600"
                  }
                >
                  {r.status}
                </span>
              </td>
              <td className="p-2">{formatThaiDateTime(r.request_date)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
