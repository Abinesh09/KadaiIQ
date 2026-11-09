export interface ApiResponse<T = unknown> {
  success: boolean
  message: string
  data?: T
  error?: unknown
}

export interface HealthCheckData {
  service: string
  status: string
  version: string
  timestamp: string
}

const BASE_URL = import.meta.env.VITE_API_URL || '/api/v1'

export const apiService = {
  async getHealth(): Promise<ApiResponse<HealthCheckData>> {
    const res = await fetch(`${BASE_URL}/health`)
    if (!res.ok) {
      throw new Error(`HTTP error ${res.status}: ${res.statusText}`)
    }
    return res.json()
  },
}

