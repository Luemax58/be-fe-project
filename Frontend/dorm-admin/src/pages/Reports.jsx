import React from "react";

export default function Reports() {
  return (
    <div>
      <h2 className="text-xl font-semibold mb-4">รายงาน</h2>

      <div className="bg-white rounded-lg shadow p-4 mb-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <div className="text-sm text-gray-500">รายรับรวมเดือนนี้</div>
            <div className="text-2xl font-bold">85,000 บาท</div>
          </div>

          <div>
            <div className="text-sm text-gray-500">จำนวนห้องที่ยังไม่จ่าย</div>
            <div className="text-2xl font-bold text-red-600">3</div>
          </div>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow p-4">
        <div className="h-48 flex items-center justify-center text-gray-400">
          [กราฟ Mock รายรับย้อนหลัง]
        </div>
      </div>
    </div>
  );
}
