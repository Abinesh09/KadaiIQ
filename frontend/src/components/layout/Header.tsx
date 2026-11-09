import React from 'react'
import { useNavigate } from 'react-router-dom'
import { Bell, User, LogOut, ChevronDown, Building2 } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { useAuth } from '@/features/auth/AuthContext'

export const Header: React.FC = () => {
  const { user, activeTenant, availableTenants, switchTenant, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  const handleTenantChange = async (e: React.ChangeEvent<HTMLSelectElement>) => {
    const selectedTenantId = e.target.value
    if (selectedTenantId && selectedTenantId !== activeTenant?.tenant_id) {
      await switchTenant(selectedTenantId)
    }
  }

  return (
    <header className="h-16 border-b border-slate-200 bg-white px-6 flex items-center justify-between">
      {/* Active Organisation / Switcher */}
      <div className="flex items-center space-x-3">
        <Building2 className="w-5 h-5 text-blue-600" />
        {availableTenants.length > 1 ? (
          <div className="relative">
            <select
              value={activeTenant?.tenant_id || ''}
              onChange={handleTenantChange}
              className="appearance-none bg-slate-50 border border-slate-200 text-xs font-semibold text-slate-800 rounded-lg pl-3 pr-8 py-1.5 focus:outline-none focus:ring-1 focus:ring-blue-600"
            >
              {availableTenants.map((t) => (
                <option key={t.tenant_id} value={t.tenant_id}>
                  {t.tenant_name} ({t.role_name})
                </option>
              ))}
            </select>
            <ChevronDown className="w-3.5 h-3.5 text-slate-400 absolute right-2.5 top-2.5 pointer-events-none" />
          </div>
        ) : (
          <div>
            <span className="text-sm font-semibold text-slate-800">
              {activeTenant?.tenant_name || 'KadaiIQ Store'}
            </span>
          </div>
        )}
        <Badge variant="outline" className="text-xs">
          {activeTenant?.role_name || (activeTenant?.is_super_admin ? 'Super Admin' : 'Member')}
        </Badge>
      </div>

      {/* User Actions */}
      <div className="flex items-center space-x-4">
        <button
          type="button"
          aria-label="Notifications"
          className="p-2 text-slate-500 hover:text-slate-700 hover:bg-slate-100 rounded-full transition-colors"
        >
          <Bell className="w-4 h-4" />
        </button>

        <div className="flex items-center space-x-3 pl-3 border-l border-slate-200">
          <div className="flex items-center space-x-2">
            <div className="w-8 h-8 rounded-full bg-blue-100 text-blue-700 flex items-center justify-center font-medium text-xs">
              <User className="w-4 h-4" />
            </div>
            <div className="text-left hidden sm:block">
              <span className="text-xs font-semibold text-slate-800 block">
                {user?.name || 'Store User'}
              </span>
              <span className="text-[10px] text-slate-400 block truncate max-w-[140px]">
                {user?.email}
              </span>
            </div>
          </div>

          <button
            onClick={handleLogout}
            title="Log Out"
            className="p-1.5 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
          >
            <LogOut className="w-4 h-4" />
          </button>
        </div>
      </div>
    </header>
  )
}
