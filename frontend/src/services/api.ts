export interface ApiResponse<T = unknown> {
  success: boolean
  message: string
  data?: T
  error?: unknown
}

export interface AppHealthData {
  service: string
  status: string
  version: string
  timestamp: string
}

export interface DbHealthData {
  status: string
  latency_ms: number
  open_connections: number
  in_use: number
  idle: number
  timestamp: string
}

const BASE_URL = import.meta.env.VITE_API_URL || '/api/v1'

export const apiService = {
  async getHealth(): Promise<ApiResponse<AppHealthData>> {
    const res = await fetch(`${BASE_URL}/health`)
    if (!res.ok) {
      throw new Error(`Health check failed: HTTP ${res.status}`)
    }
    return res.json()
  },

  async getDbHealth(): Promise<ApiResponse<DbHealthData>> {
    const res = await fetch(`${BASE_URL}/health/db`)
    const data: ApiResponse<DbHealthData> = await res.json()
    if (!res.ok) {
      throw new Error(data.message || `Database check failed: HTTP ${res.status}`)
    }
    return data
  },
}
