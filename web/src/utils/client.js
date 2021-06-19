import axios from "axios";


export const client = axios.create({baseURL: "http://localhost:5000/api/", timeout: 5000})

client.interceptors.response.use((res) => res.data)
