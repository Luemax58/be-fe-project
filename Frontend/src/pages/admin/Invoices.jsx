import { useEffect, useState } from "react";
import axiosClient from "../../api/axiosClient";

export default function Invoices() {
  const [invoices, setInvoices] = useState([]);

  useEffect(() => {
    const loadInvoices = async () => {
      try {
        const res = await axiosClient.get("/billing/invoices", {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
        });

        console.log("INVOICE DATA:", res.data);
        setInvoices(res.data);
      } catch (err) {
        console.error("Failed to load invoices:", err);
      }
    };

    loadInvoices();
  }, []);

  return (
    <div>
      <h1 className="text-xl font-semibold mb-4">รายการบิลค่าน้ำ/ไฟ</h1>

      <table className="w-full border">
        <thead className="bg-gray-100 border">
          <tr>
            <th className="p-2 border">เดือนบิล</th>
            <th className="p-2 border">ห้อง</th>
            <th className="p-2 border">ผู้เช่า</th>
            <th className="p-2 border text-right">ยอดรวม (บาท)</th>
            <th className="p-2 border text-center">สถานะ</th>
          </tr>
        </thead>

        <tbody>
          {invoices.map((inv) => (
            <tr key={inv.billing_id} className="border">
              {/* เดือนบิล */}
              <td className="p-2 border">
                {inv.billing_month
                  ? new Date(inv.billing_month).toLocaleDateString("th-TH", {
                      year: "numeric",
                      month: "long",
                    })
                  : "-"}
              </td>

              {/* ห้อง */}
              <td className="p-2 border">
                {inv.room_number ? `ห้อง ${inv.room_number}` : "-"}
              </td>

              {/* ผู้เช่า */}
              <td className="p-2 border">{inv.full_name || "-"}</td>

              {/* ยอดรวม */}
              <td className="p-2 border text-right">
                {inv.total_utility_bill
                  ? inv.total_utility_bill.toLocaleString()
                  : "-"}
              </td>

              {/* สถานะ */}
              <td className="p-2 border text-center">
                <span
                  className={
                    inv.status === "paid"
                      ? "text-green-600"
                      : inv.status === "overdue"
                      ? "text-red-600 font-bold"
                      : "text-orange-500"
                  }
                >
                  {inv.status || "-"}
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
