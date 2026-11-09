import React, { createContext, useContext, useState, useEffect } from 'react'
import { apiService } from '@/services/api'
import type {
  UserProfile,
  TenantMembershipInfo,
  AuthResponseData,
} from '@/services/api'

interface AuthContextType {
  user: UserProfile | null
  activeTenant: TenantMembershipInfo | null
  availableTenants: TenantMembershipInfo[]
  isAuthenticated: boolean
  isSuperAdmin: boolean
  isLoading: boolean
  login: (data: { email: string; password: string }) => Promise<AuthResponseData>
  logout: () => Promise<void>
  switchTenant: (tenantId: string) => Promise<void>
  setAuthData: (data: AuthResponseData) => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<UserProfile | null>(null)
  const [activeTenant, setActiveTenant] = useState<TenantMembershipInfo | null>(null)
  const [availableTenants, setAvailableTenants] = useState<TenantMembershipInfo[]>([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    try {
      const storedToken = localStorage.getItem('kadai_token')
      const storedUser = localStorage.getItem('kadai_user')
      const storedTenant = localStorage.getItem('kadai_active_tenant')

      if (storedToken && storedUser) {
        setUser(JSON.parse(storedUser))
        if (storedTenant) {
          setActiveTenant(JSON.parse(storedTenant))
        }
      }
    } catch {
      // Ignored malformed storage
    } finally {
      setIsLoading(false)
    }
  }, [])

  const setAuthData = (data: AuthResponseData) => {
    localStorage.setItem('kadai_token', data.access_token)
    localStorage.setItem('kadai_refresh_token', data.refresh_token)
    localStorage.setItem('kadai_user', JSON.stringify(data.user))
    if (data.active_tenant) {
      localStorage.setItem('kadai_active_tenant', JSON.stringify(data.active_tenant))
    }
    setUser(data.user)
    setActiveTenant(data.active_tenant || null)
    setAvailableTenants(data.available_tenants)
  }

  const login = async (payload: { email: string; password: string }) => {
    const res = await apiService.login(payload)
    if (res.data) {
      setAuthData(res.data)
      return res.data
    }
    throw new Error(res.message || 'Login failed')
  }

  const logout = async () => {
    await apiService.logout()
    setUser(null)
    setActiveTenant(null)
    setAvailableTenants([])
  }

  const switchTenant = async (tenantId: string) => {
    const res = await apiService.switchTenant(tenantId)
    if (res.data) {
      setAuthData(res.data)
    }
  }

  const isSuperAdmin = activeTenant?.is_super_admin || false
  const isAuthenticated = !!user

  return (
    <AuthContext.Provider
      value={{
        user,
        activeTenant,
        availableTenants,
        isAuthenticated,
        isSuperAdmin,
        isLoading,
        login,
        logout,
        switchTenant,
        setAuthData,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
