import axiosClient from "./axiosClient";

const roomApi = {
  getRooms: () => axiosClient.get("/rooms"),
};

export default roomApi;
