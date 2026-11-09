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

export interface UserProfile {
  id: string
  name: string
  email: string
}

export interface TenantMembershipInfo {
  tenant_id: string
  tenant_name: string
  tenant_slug: string
  role_id: string
  role_name: string
  status: string
  is_super_admin: boolean
}

export interface AuthResponseData {
  access_token: string
  refresh_token: string
  expires_in: number
  user: UserProfile
  active_tenant?: TenantMembershipInfo
  available_tenants: TenantMembershipInfo[]
}

export interface TenantData {
  id: string
  name: string
  slug: string
  business_name: string
  business_type: string
  gstin?: string
  email: string
  phone?: string
  address?: string
  status: string
  plan_name?: string
  member_count: number
  created_at: string
}

export interface MemberData {
  id: string
  user_id: string
  name: string
  email: string
  role_id: string
  role_name: string
  status: string
  created_at: string
}

const BASE_URL = import.meta.env.VITE_API_URL || '/api/v1'

const getHeaders = (): Record<string, string> => {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }
  const token = localStorage.getItem('kadai_token')
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  return headers
}

export const apiService = {
  // System Health
  async getHealth(): Promise<ApiResponse<AppHealthData>> {
    const res = await fetch(`${BASE_URL}/health`)
    if (!res.ok) throw new Error(`Health check failed: HTTP ${res.status}`)
    return res.json()
  },

  async getDbHealth(): Promise<ApiResponse<DbHealthData>> {
    const res = await fetch(`${BASE_URL}/health/db`)
    const data: ApiResponse<DbHealthData> = await res.json()
    if (!res.ok) throw new Error(data.message || `Database check failed: HTTP ${res.status}`)
    return data
  },

  // Auth Endpoints
  async login(payload: { email: string; password: string }): Promise<ApiResponse<AuthResponseData>> {
    const res = await fetch(`${BASE_URL}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Login failed')
    return data
  },

  async register(payload: {
    organization_name: string
    slug: string
    owner_name: string
    email: string
    password: string
    phone?: string
  }): Promise<ApiResponse<AuthResponseData>> {
    const res = await fetch(`${BASE_URL}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Registration failed')
    return data
  },

  async forgotPassword(email: string): Promise<ApiResponse<{ reset_token?: string }>> {
    const res = await fetch(`${BASE_URL}/auth/forgot-password`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email }),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Forgot password failed')
    return data
  },

  async resetPassword(payload: { token: string; new_password: string }): Promise<ApiResponse<null>> {
    const res = await fetch(`${BASE_URL}/auth/reset-password`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Password reset failed')
    return data
  },

  async switchTenant(tenant_id: string): Promise<ApiResponse<AuthResponseData>> {
    const res = await fetch(`${BASE_URL}/auth/switch-tenant`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify({ tenant_id }),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to switch organisation')
    return data
  },

  async logout(): Promise<void> {
    const refreshToken = localStorage.getItem('kadai_refresh_token')
    if (refreshToken) {
      await fetch(`${BASE_URL}/auth/logout`, {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify({ refresh_token: refreshToken }),
      }).catch(() => {})
    }
    localStorage.removeItem('kadai_token')
    localStorage.removeItem('kadai_refresh_token')
    localStorage.removeItem('kadai_user')
    localStorage.removeItem('kadai_active_tenant')
  },

  // Super Admin Endpoints
  async getTenants(): Promise<ApiResponse<TenantData[]>> {
    const res = await fetch(`${BASE_URL}/super-admin/tenants`, {
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch organisations')
    return data
  },

  async createTenant(payload: {
    name: string
    slug: string
    business_name: string
    business_type: string
    gstin?: string
    email: string
    phone?: string
    address?: string
    owner_name: string
    owner_email: string
    owner_password: string
  }): Promise<ApiResponse<TenantData>> {
    const res = await fetch(`${BASE_URL}/super-admin/tenants`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify(payload),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to create organisation')
    return data
  },

  async updateTenantStatus(id: string, status: string): Promise<ApiResponse<null>> {
    const res = await fetch(`${BASE_URL}/super-admin/tenants/${id}/status`, {
      method: 'PATCH',
      headers: getHeaders(),
      body: JSON.stringify({ status }),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to update status')
    return data
  },

  // Tenant Admin Endpoints
  async getMembers(): Promise<ApiResponse<MemberData[]>> {
    const res = await fetch(`${BASE_URL}/tenant/members`, {
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch members')
    return data
  },

  async inviteMember(payload: { name: string; email: string; role_id: string }): Promise<ApiResponse<MemberData>> {
    const res = await fetch(`${BASE_URL}/tenant/members`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify(payload),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to invite member')
    return data
  },
}
