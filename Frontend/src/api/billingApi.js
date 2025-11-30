import axiosClient from "./axiosClient";

const billingApi = {
  // tenant
  myInvoices: () => axiosClient.get("/billing/my-invoices"),
  myPayments: () => axiosClient.get("/billing/my-payments"),
  pay: (data) => axiosClient.post("/billing/pay", data),

  // admin
  generateInvoices: (data) =>
    axiosClient.post("/billing/admin/invoices/generate", data),
  recordUtilities: (data) =>
    axiosClient.post("/billing/admin/utilities/record", data),
  recordPayment: (data) =>
    axiosClient.post("/billing/admin/payments/record", data),
};

export default billingApi;
