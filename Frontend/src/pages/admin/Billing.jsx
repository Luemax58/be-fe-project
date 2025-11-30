import { useEffect, useState } from "react";
import AdminLayout from "../../layouts/AdminLayout";
import { formatThaiDateTime } from "../../utils/dateFormat";
import axiosClient from "../../api/axiosClient";

export default function BillingPage() {
  const [bills, setBills] = useState([]);
  const [month, setMonth] = useState("");
  const [year, setYear] = useState(new Date().getFullYear() + 543);
  const [loading, setLoading] = useState(false);

  // โหลดบิลทั้งหมด
  useEffect(() => {
    loadBills();
  }, []);

  const loadBills = async () => {
    try {
      const res = await axiosClient.get("/billing/admin/all");

      if (Array.isArray(res.data)) {
        setBills(res.data);
      } else {
        console.error("Billing format unexpected:", res.data);
        setBills([]);
      }

    } catch (err) {
      console.error("loadBills error:", err);
    }
  };

  // สร้างบิลรายเดือน
  const createBilling = async () => {
    if (!month || !year) {
      alert("กรุณาเลือกปีและเดือนก่อนสร้างบิล");
      return;
    }

    if (!window.confirm(`ยืนยันการสร้างบิล เดือน ${month}/${year}?`)) return;

    setLoading(true);

    try {
      await axiosClient.post("/billing/admin/invoices/generate", {
        month: Number(month),
        year: Number(year),
      });

      alert("สร้างบิลสำเร็จ");
      loadBills();
    } catch (err) {
      alert("สร้างบิลล้มเหลว");
      console.error("createBilling error →", err);
    }

    setLoading(false);
  };

  return (
    <AdminLayout>
      <div className="p-6">

        <h1 className="text-3xl font-bold mb-6">จัดการ Billing (Admin)</h1>

        {/* กล่องสร้างบิล */}
        <div className="bg-white p-5 rounded-xl shadow mb-8">
          <h2 className="text-xl font-semibold mb-3">สร้างบิลรายเดือน</h2>

          <div className="flex items-center space-x-4">

            {/* ปี */}
            <select
              value={year}
              onChange={(e) => setYear(e.target.value)}
              className="border p-2 rounded min-w-[120px]"
            >
              <option value="">เลือกปี</option>
              <option value="2568">2568</option>
              <option value="2567">2567</option>
              <option value="2566">2566</option>
              <option value="2565">2565</option>
            </select>

            {/* เดือน */}
            <select
              value={month}
              onChange={(e) => setMonth(e.target.value)}
              className="border p-2 rounded min-w-[120px]"
            >
              <option value="">เลือกเดือน</option>
              <option value="1">มกราคม</option>
              <option value="2">กุมภาพันธ์</option>
              <option value="3">มีนาคม</option>
              <option value="4">เมษายน</option>
              <option value="5">พฤษภาคม</option>
              <option value="6">มิถุนายน</option>
              <option value="7">กรกฎาคม</option>
              <option value="8">สิงหาคม</option>
              <option value="9">กันยายน</option>
              <option value="10">ตุลาคม</option>
              <option value="11">พฤศจิกายน</option>
              <option value="12">ธันวาคม</option>
            </select>

            <button
              onClick={createBilling}
              className={`px-4 py-2 rounded text-white ${
                loading ? "bg-gray-400" : "bg-blue-600 hover:bg-blue-700"
              }`}
            >
              {loading ? "กำลังสร้าง..." : "สร้างบิล"}
            </button>
          </div>
        </div>

        {/* ตาราง */}
        <div className="bg-white p-5 rounded-xl shadow">
          <h2 className="text-xl font-semibold mb-4">รายการบิลทั้งหมด</h2>

          <table className="min-w-full table-auto border-collapse bg-white">
            <thead>
              <tr className="bg-gray-100 text-left">
                <th className="p-3">Room</th>
                <th className="p-3">Month</th>
                <th className="p-3">Water</th>
                <th className="p-3">Electric</th>
                <th className="p-3">Total</th>
                <th className="p-3">Status</th>
                <th className="p-3">Action</th>
              </tr>
            </thead>

            <tbody>
              {bills.map((b) => (
                <tr key={b.billing_id} className="border-t">
                  <td className="p-2">{b.room_number}</td>
                  <td className="p-2">{formatThaiDateTime(b.billing_month)}</td>
                  <td className="p-2">{b.water_bill}</td>
                  <td className="p-2">{b.electricity_bill}</td>
                  <td className="p-2">{b.total_utility_bill}</td>

                  <td className="p-2">
                    <span
                      className={
                        b.status === "paid"
                          ? "text-green-600"
                          : b.status === "overdue"
                          ? "text-red-600"
                          : "text-yellow-600"
                      }
                    >
                      {b.status}
                    </span>
                  </td>

                  <td className="p-2">
                    <button className="px-3 py-1 bg-yellow-500 text-white rounded mr-2">
                      น้ำ/ไฟ
                    </button>
                    <button className="px-3 py-1 bg-green-600 text-white rounded">
                      รับเงิน
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

        </div>
      </div>
    </AdminLayout>
  );
}
