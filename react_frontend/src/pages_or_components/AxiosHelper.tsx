import axios from "axios"

const axiosHelper = axios.create({
    baseURL: "https://song-recommendations-web-app-7jyz.onrender.com" + "/api/v1",
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