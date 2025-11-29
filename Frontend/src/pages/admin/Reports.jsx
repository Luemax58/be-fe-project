import React, { useEffect, useState } from "react";
import axiosClient from "../../api/axiosClient";

export default function Reports() {
  const [invoices, setInvoices] = useState([]);

  useEffect(() => {
    const load = async () => {
      try {
        const token = localStorage.getItem("token");
        const res = await axiosClient.get("/billing/invoices", {
          headers: { Authorization: `Bearer ${token}` },
        });
        setInvoices(res.data);
      } catch (err) {
        console.error(err);
      }
    };

    load();
  }, []);

  const totalBilling = invoices.reduce(
    (sum, i) => sum + (i.total_utility_bill ?? 0),
    0
  );

  const totalPaid = invoices.reduce((sum, i) => {
    const payments = i.Payments || [];
    const paidAmount = payments.reduce(
      (p, t) => p + (t.amount_paid ?? 0),
      0
    );
    return sum + paidAmount;
  }, 0);

  const outstanding = totalBilling - totalPaid;

  return (
    <div>
      <h2 className="text-xl font-semibold mb-4">รายงานภาพรวม</h2>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-white p-4 rounded shadow">
          <p className="text-gray-500">ยอดบิลทั้งหมด</p>
          <p className="text-2xl font-bold">{totalBilling} บาท</p>
        </div>

        <div className="bg-white p-4 rounded shadow">
          <p className="text-gray-500">ยอดที่ชำระแล้ว</p>
          <p className="text-2xl font-bold">{totalPaid} บาท</p>
        </div>

        <div className="bg-white p-4 rounded shadow">
          <p className="text-gray-500">ยอดค้างชำระ</p>
          <p className="text-2xl font-bold">{outstanding} บาท</p>
        </div>
      </div>
    </div>
  );
}
