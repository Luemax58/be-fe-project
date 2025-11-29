import React, { useEffect, useState } from "react";
import axiosClient from "../../api/axiosClient";

export default function PaymentForm() {
  const [invoices, setInvoices] = useState([]);

  useEffect(() => {
    const loadInvoices = async () => {
      try {
        const token = localStorage.getItem("token");
        const res = await axiosClient.get("/tenant/invoices", {
          headers: { Authorization: `Bearer ${token}` },
        });
        setInvoices(res.data);
      } catch (err) {
        console.error(err);
      }
    };

    loadInvoices();
  }, []);

  return (
    <div>
      <h2 className="text-xl font-semibold mb-4">ข้อมูลใบแจ้งหนี้</h2>

      {invoices.length === 0 ? (
        <p>ยังไม่มีข้อมูล</p>
      ) : (
        invoices.map((inv) => (
          <div key={inv.billing_id} className="p-3 border-b">
            ห้อง {inv.Room?.room_number ?? "-"}  
            | ยอดค่าน้ำไฟ: {inv.total_utility_bill} บาท
            | สถานะ: {inv.status}
          </div>
        ))
      )}
    </div>
  );
}
