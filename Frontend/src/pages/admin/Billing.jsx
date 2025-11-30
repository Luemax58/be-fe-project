import { useEffect, useState } from "react";
import axiosClient from "../../api/axiosClient";
import { formatThaiDateTime } from "../../utils/dateFormat";

export default function Billing() {
  const [billings, setBillings] = useState([]);
  const [loading, setLoading] = useState(true);

  // ฟอร์มสร้างบิล
  const [month, setMonth] = useState("");
  const [dueDays, setDueDays] = useState(7);

  // ฟอร์มแก้น้ำไฟ
  const [editBill, setEditBill] = useState(null);
  const [waterUnits, setWaterUnits] = useState("");
  const [electricityUnits, setElectricityUnits] = useState("");

  // ฟอร์มรับเงิน
  const [payBill, setPayBill] = useState(null);
  const [amountPaid, setAmountPaid] = useState("");

  useEffect(() => {
    loadBillings();
  }, []);

  const loadBillings = async () => {
    try {
      const res = await axiosClient.get("/billing/admin/all"); // ถ้ามี API นี้
      setBillings(res.data || []);
    } catch (err) {
      console.error(err);
      alert("ดึงข้อมูลบิลไม่สำเร็จ");
    } finally {
      setLoading(false);
    }
  };

  // สร้างบิล
  const generateInvoices = async () => {
    try {
      await axiosClient.post("/billing/admin/invoices/generate", {
        month,
        due_days_after: dueDays,
      });
      alert("สร้างบิลสำเร็จ");
      loadBillings();
    } catch (err) {
      console.error(err);
      alert("สร้างบิลล้มเหลว");
    }
  };

  // อัปเดตค่าน้ำไฟ
  const recordUtility = async () => {
    try {
      await axiosClient.post("/billing/admin/utilities/record", {
        billing_id: editBill.billing_id,
        water_units: Number(waterUnits),
        electricity_units: Number(electricityUnits),
      });
      alert("บันทึกสำเร็จ");
      setEditBill(null);
      loadBillings();
    } catch (err) {
      console.error(err);
      alert("บันทึกน้ำไฟล้มเหลว");
    }
  };

  // รับเงิน
  const recordPayment = async () => {
    try {
      await axiosClient.post("/billing/admin/payments/record", {
        billing_id: payBill.billing_id,
        tenant_id: payBill.tenant_id,
        amount_paid: Number(amountPaid),
        payment_method: "cash",
      });
      alert("รับเงินสำเร็จ");
      setPayBill(null);
      loadBillings();
    } catch (err) {
      console.error(err);
      alert("รับเงินล้มเหลว");
    }
  };

  if (loading) return <div className="p-5">Loading...</div>;

  return (
    <div className="p-5">
      <h1 className="text-2xl font-bold mb-4">จัดการ Billing (Admin)</h1>

      {/* ฟอร์มสร้างบิล */}
      <div className="bg-white p-4 shadow mb-6 rounded-lg">
        <h2 className="text-lg font-semibold mb-2">สร้างบิลรายเดือน</h2>
        <div className="flex items-center gap-4">
          <input
            type="month"
            value={month}
            onChange={(e) => setMonth(e.target.value)}
            className="border p-2 rounded"
          />

          <input
            type="number"
            value={dueDays}
            onChange={(e) => setDueDays(e.target.value)}
            className="border p-2 rounded"
          />

          <button
            className="bg-blue-600 text-white px-4 py-2 rounded"
            onClick={generateInvoices}
          >
            สร้างบิล
          </button>
        </div>
      </div>

      {/* ตารางบิล */}
      <table className="w-full border">
        <thead className="bg-gray-200">
          <tr>
            <th className="p-2">Room</th>
            <th className="p-2">Month</th>
            <th className="p-2">Water</th>
            <th className="p-2">Electric</th>
            <th className="p-2">Total</th>
            <th className="p-2">Status</th>
            <th className="p-2">Action</th>
          </tr>
        </thead>

        <tbody>
          {billings.map((b) => (
            <tr key={b.billing_id} className="border-t">
              <td className="p-2">{b.room_id}</td>
              <td className="p-2">{formatThaiDateTime(b.billing_month)}</td>
              <td className="p-2">{b.water_units}</td>
              <td className="p-2">{b.electricity_units}</td>
              <td className="p-2">{b.total_utility_bill}</td>
              <td className="p-2">{b.status}</td>

              <td className="p-2 flex gap-2">
                <button
                  className="bg-yellow-500 text-white px-3 py-1 rounded"
                  onClick={() => setEditBill(b)}
                >
                  น้ำ/ไฟ
                </button>
                <button
                  className="bg-green-600 text-white px-3 py-1 rounded"
                  onClick={() => setPayBill(b)}
                >
                  รับเงิน
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* Modal บันทึกน้ำไฟ */}
      {editBill && (
        <div className="fixed inset-0 bg-black bg-opacity-40 flex justify-center items-center">
          <div className="bg-white p-5 rounded shadow w-96">
            <h2 className="text-lg font-bold mb-2">บันทึกน้ำ-ไฟ</h2>

            <input
              type="number"
              placeholder="Water"
              value={waterUnits}
              onChange={(e) => setWaterUnits(e.target.value)}
              className="border p-2 w-full mb-2"
            />

            <input
              type="number"
              placeholder="Electric"
              value={electricityUnits}
              onChange={(e) => setElectricityUnits(e.target.value)}
              className="border p-2 w-full mb-4"
            />

            <button
              onClick={recordUtility}
              className="bg-blue-600 text-white px-4 py-2 rounded"
            >
              บันทึก
            </button>

            <button
              onClick={() => setEditBill(null)}
              className="ml-2 text-gray-600"
            >
              ยกเลิก
            </button>
          </div>
        </div>
      )}

      {/* Modal รับเงิน */}
      {payBill && (
        <div className="fixed inset-0 bg-black bg-opacity-40 flex justify-center items-center">
          <div className="bg-white p-5 rounded shadow w-96">
            <h2 className="text-lg font-bold mb-3">รับเงินบิล</h2>

            <input
              type="number"
              placeholder="จำนวนเงิน"
              value={amountPaid}
              onChange={(e) => setAmountPaid(e.target.value)}
              className="border p-2 w-full mb-3"
            />

            <button
              onClick={recordPayment}
              className="bg-green-600 text-white px-4 py-2 rounded"
            >
              บันทึกการจ่ายเงิน
            </button>

            <button
              onClick={() => setPayBill(null)}
              className="ml-2 text-gray-600"
            >
              ยกเลิก
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
