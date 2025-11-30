import axiosClient from "./axiosClient";

const authApi = {
  login: (data) => axiosClient.post("/login", data),

  getProfile: () => axiosClient.get("/users/me"),
};

export default authApi;
