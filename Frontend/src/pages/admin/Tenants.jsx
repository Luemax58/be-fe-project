import { useEffect, useState } from "react";
import axiosClient from "../../api/axiosClient";

export default function Tenants() {
  const [tenants, setTenants] = useState([]);

  useEffect(() => {
    const loadTenants = async () => {
      try {
        const token = localStorage.getItem("token");

        const res = await axiosClient.get("/tenants", {
          headers: { Authorization: `Bearer ${token}` }
        });

        console.log("TENANT DATA:", res.data); 
        setTenants(res.data);
      } catch (err) {
        console.error("LOAD TENANTS ERROR:", err);
      }
    };

    loadTenants();
  }, []);

  return (
    <div>
      <h1 className="text-xl font-semibold mb-4">รายชื่อผู้เช่า</h1>

      <table className="w-full border">
        <thead className="bg-gray-100 border">
          <tr>
            <th className="p-2 border">ชื่อจริง</th>
            <th className="p-2 border">เบอร์โทร</th>
            <th className="p-2 border">ชื่อผู้ใช้</th>
          </tr>
        </thead>

        <tbody>
          {tenants.map((t) => (
            <tr key={t.UserID} className="border">
              <td className="p-2 border">{t.FullName}</td>
              <td className="p-2 border">{t.Phone}</td>
              <td className="p-2 border">{t.Username}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
