import { useEffect, useState } from "react";
import axiosClient from "../../api/axiosClient";
import { formatThaiDateTime } from "../../utils/dateFormat";
import TenantLayout from "../../layouts/TenantLayout";

export default function TenantDashboard() {
  const [profile, setProfile] = useState(null);
  const [latestInvoice, setLatestInvoice] = useState(null);
  const [maintenance, setMaintenance] = useState([]);
  const [room, setRoom] = useState(null);

  // โหลดข้อมูลทั้งหมดตอนเปิดหน้า
  useEffect(() => {
    loadAll();
  }, []);

  const loadAll = async () => {
    try {
      // โหลดโปรไฟล์
      const resProfile = await axiosClient.get("/users/me");
      setProfile(resProfile.data);

      // โหลดบิลล่าสุด
      const resBills = await axiosClient.get("/billing/my-invoices");
      if (resBills.data && resBills.data.length > 0) {
        setLatestInvoice(resBills.data[0]);
      }

      // โหลดคำร้องซ่อมทั้งหมด
      const resMaint = await axiosClient.get("/maintenance/requests");
      setMaintenance(resMaint.data.data || []);
    } catch (error) {
      console.error("TenantDashboard load error:", error);
    }
  };

  // โหลดห้องพักเมื่อมีข้อมูล profile
  useEffect(() => {
    if (profile) {
      loadRoom();
    }
  }, [profile]);

  const loadRoom = async () => {
    try {
      const resRooms = await axiosClient.get("/rooms");

      // หาเฉพาะห้องที่ tenant.user_id === profile.user_id
      const myRoom = resRooms.data.find(
        (r) => r.tenant?.user_id === profile.user_id
      );

      setRoom(myRoom);
    } catch (error) {
      console.error("loadRoom error:", error);
    }
  };

  if (!profile) return <TenantLayout><div className="p-6">Loading...</div></TenantLayout>;

  return (
    <TenantLayout>
      <div className="p-6">
        <h1 className="text-2xl font-bold mb-5">หน้าหลักผู้เช่า</h1>

        {/* ---------- ข้อมูลผู้เช่า ---------- */}
        <div className="bg-white shadow p-4 rounded-lg mb-6">
          <h2 className="text-xl font-semibold mb-2">ข้อมูลผู้เช่า</h2>
          <p><b>ชื่อ:</b> {profile.full_name}</p>
          <p><b>เบอร์โทร:</b> {profile.phone || "-"}</p>
          <p>
            <b>ห้องพัก:</b> {room ? room.room_number : "ไม่พบข้อมูล"}
          </p>
        </div>

        {/* ---------- บิลล่าสุด ---------- */}
        <div className="bg-indigo-50 shadow p-4 rounded-lg mb-6">
          <h2 className="text-xl font-semibold mb-2">บิลล่าสุด</h2>

          {latestInvoice ? (
            <>
              <p><b>เดือน:</b> {formatThaiDateTime(latestInvoice.billing_month)}</p>
              <p><b>ยอดรวม:</b> {latestInvoice.total_utility_bill} บาท</p>
              <p>
                <b>สถานะ:</b>{" "}
                <span
                  className={
                    latestInvoice.status === "paid"
                      ? "text-green-600 font-semibold"
                      : latestInvoice.status === "overdue"
                      ? "text-red-600 font-semibold"
                      : "text-yellow-600 font-semibold"
                  }
                >
                  {latestInvoice.status}
                </span>
              </p>
            </>
          ) : (
            <p>ยังไม่มีบิล</p>
          )}
        </div>

        {/* ---------- คำร้องซ่อมล่าสุด ---------- */}
        <div className="bg-white shadow p-4 rounded-lg">
          <h2 className="text-xl font-semibold mb-2">คำร้องซ่อมล่าสุด</h2>

          <table className="w-full border text-left">
            <thead>
              <tr className="bg-gray-100">
                <th className="p-2">ปัญหา</th>
                <th className="p-2">สถานะ</th>
                <th className="p-2">วันที่แจ้ง</th>
              </tr>
            </thead>
            <tbody>
              {maintenance.slice(0, 5).map((m, idx) => (
                <tr key={idx} className="border-t">
                  <td className="p-2">{m.issue_description}</td>

                  {/* ⭐ สถานะ พร้อมสี */}
                  <td className="p-2">
                    <span
                      className={
                        m.status === "pending"
                          ? "text-yellow-600 font-semibold"
                          : m.status === "in_progress"
                          ? "text-blue-600 font-semibold"
                          : m.status === "completed"
                          ? "text-green-600 font-semibold"
                          : "text-gray-600"
                      }
                    >
                      {m.status}
                    </span>
                  </td>

                  <td className="p-2">{formatThaiDateTime(m.request_date)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </TenantLayout>
  );
}
