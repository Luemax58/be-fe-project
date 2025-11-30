import { useState, useEffect } from "react";
import axiosClient from "../../api/axiosClient";
import { formatThaiDateTime } from "../../utils/dateFormat";
import TenantLayout from "../../layouts/TenantLayout";

export default function TenantBilling() {
  const [invoices, setInvoices] = useState([]);

  useEffect(() => {
    loadInvoices();
  }, []);

  const loadInvoices = async () => {
    try {
      const res = await axiosClient.get("/billing/my-invoices");
      setInvoices(res.data || []);
    } catch (error) {
      console.error(error);
    }
  };

  const payInvoice = async (billingId) => {
    try {
      await axiosClient.post("/billing/pay", {
        billing_id: billingId,
        amount_paid: 999,   // คุณกำหนดเอง/หรือให้ผู้ใช้ใส่ทีหลัง
        method: "transfer",
        notes: "",
      });

      alert("ชำระเงินสำเร็จ");
      loadInvoices();
    } catch (err) {
      alert("ชำระเงินล้มเหลว");
      console.error(err);
    }
  };

  return (
    <TenantLayout>
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-5">บิลของฉัน</h1>

      <table className="w-full border text-left">
        <thead>
          <tr className="bg-gray-100">
            <th className="p-2">เดือน</th>
            <th className="p-2">ยอดรวม</th>
            <th className="p-2">สถานะ</th>
            <th className="p-2">จัดการ</th>
          </tr>
        </thead>
        <tbody>
          {invoices.map((b) => (
            <tr key={b.billing_id ?? b.billing_month} className="border-t">
            
              {/* เดือน */}
              <td className="p-2">{formatThaiDateTime(b.billing_month)}</td>

              {/* ยอดรวม */}
              <td className="p-2">{b.total_utility_bill} บาท</td>

              {/* สถานะ */}
              <td className="p-2">
                <span
                  className={
                    b.status === "paid"
                      ? "text-green-600 font-semibold"
                      : b.status === "overdue"
                      ? "text-red-600 font-semibold"
                      : "text-yellow-600 font-semibold"
                  }
                >
                  {b.status}
                </span>
              </td>

              {/* ปุ่มชำระเงิน */}
              <td className="p-2">
                {b.status !== "paid" ? (
                  <button
                    onClick={() => payInvoice(b.billing_id)}
                    className="px-3 py-1 bg-blue-600 text-white rounded"
                  >
                    ชำระเงิน
                  </button>
                ) : (
                  "-"
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
    </TenantLayout>
  );
}
