export const mockRooms = [
  { id: 1, number: "101", floor: 1, type: "แอร์", price: 3500, status: "occupied" },
  { id: 2, number: "102", floor: 1, type: "พัดลม", price: 3000, status: "vacant" },
];

export const mockTenants = [
  { id: 1, name: "สมชาย ใจดี", phone: "0812345678", room: "101", moveIn: "2024-01-15" },
];

export const mockInvoices = [
  { id: 1, tenant: "สมชาย", room: "101", total: 3500, status: "unpaid" },
];
