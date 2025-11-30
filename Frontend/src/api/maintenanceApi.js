import axiosClient from "./axiosClient";

const maintenanceApi = {
  create: (data) => axiosClient.post("/maintenance/creates", data),
  getAll: () => axiosClient.get("/maintenance/requests"),
};
export default maintenanceApi;
