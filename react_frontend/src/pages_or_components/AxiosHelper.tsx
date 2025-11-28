import axios from "axios"

const axiosHelper = axios.create({
    baseURL: "https://similarsongs-ee56.onrender.com" + "/api/v1",
    headers: {
        "Content-Type" : "application/json",
    },
    timeout: 5000,
})

axiosHelper.interceptors.request.use((config) => {
    const token = localStorage.getItem("token")
    if (token){
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})

export default axiosHelper