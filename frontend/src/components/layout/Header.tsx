import React from 'react'
import { useNavigate, useLocation } from 'react-router-dom'
import {
  Bell,
  LogOut,
  ChevronDown,
  Building2,
  Menu,
} from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { useAuth } from '@/features/auth/AuthContext'

interface HeaderProps {
  onToggleSidebar?: () => void
}

const pageTitles: Record<string, { title: string; subtitle: string }> = {
  '/': { title: 'Executive Overview', subtitle: 'Store performance, sales metrics, and stock overview' },
  '/products': { title: 'Product Catalogue', subtitle: 'Manage footwear models, variants, barcodes, and SKUs' },
  '/inventory': { title: 'Inventory Management', subtitle: 'Real-time variant stock tracking and immutable movement ledger' },
  '/sales': { title: 'Sales & POS Records', subtitle: 'Verified retail transactions, bill imports, and order history' },
  '/analytics': { title: 'Sales Analytics', subtitle: 'Revenue intelligence, product matrix, and date-range comparisons' },
  '/reports': { title: 'Reports & Exports', subtitle: 'Downloadable GST summaries, Excel exports, and financial logs' },
  '/gst': { title: 'GST Configuration', subtitle: 'Indian GST rates, tax slabs, and store compliance' },
  '/users': { title: 'Staff & Roles', subtitle: 'Manage store team members, RBAC roles, and access permissions' },
  '/settings': { title: 'Store Settings', subtitle: 'Organisation profile, billing preferences, and tenant configuration' },
  '/super-admin': { title: 'Super Admin Console', subtitle: 'Cross-tenant administration, subscription tiers, and system health' },
}

export const Header: React.FC<HeaderProps> = ({ onToggleSidebar }) => {
  const { user, activeTenant, availableTenants, switchTenant, logout } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

  const currentMeta = pageTitles[location.pathname] || {
    title: 'KadaiIQ Platform',
    subtitle: 'Retail Management and Sales Analytics SaaS',
  }

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
    <header className="h-16 border-b border-slate-200/90 bg-white/95 backdrop-blur-xs px-4 sm:px-6 flex items-center justify-between sticky top-0 z-30 shadow-[0_1px_2px_0_rgba(15,23,42,0.03)]">
      {/* Left: Mobile Toggle & Page Title */}
      <div className="flex items-center gap-3 min-w-0">
        <button
          type="button"
          onClick={onToggleSidebar}
          className="lg:hidden p-2 text-slate-500 hover:text-slate-800 hover:bg-slate-100 rounded-lg transition-colors focus:outline-none"
          aria-label="Open sidebar menu"
        >
          <Menu className="w-5 h-5" />
        </button>

        <div className="min-w-0">
          <div className="flex items-center gap-2">
            <h1 className="text-sm sm:text-base font-bold text-slate-900 tracking-tight truncate leading-tight">
              {currentMeta.title}
            </h1>
            <span className="hidden md:inline-flex items-center px-1.5 py-0.2 text-[10px] font-medium bg-slate-100 text-slate-600 rounded">
              v1.0
            </span>
          </div>
          <p className="hidden sm:block text-[11px] text-slate-400 truncate max-w-xs md:max-w-md">
            {currentMeta.subtitle}
          </p>
        </div>
      </div>

      {/* Right: Tenant Switcher, System Status & Profile */}
      <div className="flex items-center gap-2.5 sm:gap-4 shrink-0">
        {/* Active Store Switcher */}
        <div className="flex items-center gap-2 bg-slate-50/80 border border-slate-200/80 rounded-lg px-2.5 py-1">
          <Building2 className="w-4 h-4 text-blue-600 shrink-0" />
          {availableTenants.length > 1 ? (
            <div className="relative">
              <select
                value={activeTenant?.tenant_id || ''}
                onChange={handleTenantChange}
                className="appearance-none bg-transparent text-xs font-semibold text-slate-800 pr-5 focus:outline-none cursor-pointer"
              >
                {availableTenants.map((t) => (
                  <option key={t.tenant_id} value={t.tenant_id}>
                    {t.tenant_name}
                  </option>
                ))}
              </select>
              <ChevronDown className="w-3 h-3 text-slate-400 absolute right-0 top-1 pointer-events-none" />
            </div>
          ) : (
            <span className="text-xs font-semibold text-slate-800 truncate max-w-[120px] sm:max-w-[160px]">
              {activeTenant?.tenant_name || 'KadaiIQ Store'}
            </span>
          )}
          <Badge variant="outline" className="text-[10px] py-0 px-1.5 border-slate-200 text-slate-600 hidden sm:inline-flex">
            {activeTenant?.role_name || (activeTenant?.is_super_admin ? 'Super Admin' : 'Admin')}
          </Badge>
        </div>

        {/* Live System Indicator */}
        <div className="hidden md:flex items-center gap-1.5 text-[11px] font-medium text-slate-500 bg-emerald-50/60 border border-emerald-100 rounded-lg px-2 py-1">
          <span className="relative flex h-2 w-2">
            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span className="relative inline-flex rounded-full h-2 w-2 bg-emerald-500"></span>
          </span>
          <span className="text-emerald-700 text-[10px] font-semibold">Online</span>
        </div>

        {/* Notification Bell */}
        <button
          type="button"
          aria-label="System notifications"
          className="relative p-2 text-slate-400 hover:text-slate-700 hover:bg-slate-100 rounded-lg transition-colors focus:outline-none"
        >
          <Bell className="w-4 h-4" />
          <span className="absolute top-1.5 right-1.5 w-1.5 h-1.5 bg-blue-600 rounded-full" />
        </button>

        {/* User Account & Logout */}
        <div className="flex items-center gap-2 pl-2 sm:pl-3 border-l border-slate-200">
          <div className="w-7 h-7 sm:w-8 sm:h-8 rounded-lg bg-gradient-to-tr from-blue-600 to-indigo-600 text-white flex items-center justify-center font-bold text-xs shadow-xs">
            {user?.name ? user.name.charAt(0).toUpperCase() : 'U'}
          </div>
          <div className="hidden lg:block text-left">
            <span className="text-xs font-semibold text-slate-800 block leading-tight">
              {user?.name || 'Store Manager'}
            </span>
            <span className="text-[10px] text-slate-400 block truncate max-w-[120px]">
              {user?.email}
            </span>
          </div>

          <button
            onClick={handleLogout}
            title="Log Out of KadaiIQ"
            className="p-1.5 text-slate-400 hover:text-rose-600 hover:bg-rose-50 rounded-lg transition-colors ml-1 focus:outline-none cursor-pointer"
            aria-label="Logout"
          >
            <LogOut className="w-4 h-4" />
          </button>
        </div>
      </div>
    </header>
  )
}
