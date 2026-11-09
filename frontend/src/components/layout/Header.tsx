import React from 'react'
import { useNavigate, useLocation, Link } from 'react-router-dom'
import {
  LogOut,
  ChevronDown,
  Building2,
  Menu,
  ChevronRight,
} from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { useAuth } from '@/features/auth/AuthContext'

interface HeaderProps {
  onToggleSidebar?: () => void
}

const pageTitles: Record<string, { section: string; title: string }> = {
  '/': { section: 'Executive', title: 'Business Overview' },
  '/products': { section: 'Retail Operations', title: 'Product Catalogue' },
  '/inventory': { section: 'Retail Operations', title: 'Inventory Management' },
  '/sales': { section: 'Retail Operations', title: 'Sales & POS Records' },
  '/analytics': { section: 'Business Intelligence', title: 'Sales Analytics' },
  '/reports': { section: 'Financial', title: 'Reports & Exports' },
  '/gst': { section: 'Financial', title: 'GST Configuration' },
  '/users': { section: 'Organisation', title: 'Users & Roles' },
  '/settings': { section: 'Organisation', title: 'Settings' },
  '/super-admin': { section: 'Platform Admin', title: 'Super Admin' },
}

export const Header: React.FC<HeaderProps> = ({ onToggleSidebar }) => {
  const { user, activeTenant, availableTenants, switchTenant, logout } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

  let currentMeta = pageTitles[location.pathname] || {
    section: 'Workspace',
    title: 'KadaiIQ Platform',
  }

  const searchParams = new URLSearchParams(location.search)
  if (location.pathname === '/analytics') {
    const tab = searchParams.get('tab')
    if (tab === 'products') {
      currentMeta = { section: 'Business Intelligence', title: 'Product Performance' }
    } else if (tab === 'footwear') {
      currentMeta = { section: 'Business Intelligence', title: 'Category & Brand Analytics' }
    } else if (tab === 'daily') {
      currentMeta = { section: 'Business Intelligence', title: 'Daily Breakdown' }
    } else if (tab === 'monthly') {
      currentMeta = { section: 'Business Intelligence', title: 'Monthly Breakdown' }
    } else if (tab === 'yearly') {
      currentMeta = { section: 'Business Intelligence', title: 'Yearly Summary' }
    } else if (tab === 'payments') {
      currentMeta = { section: 'Business Intelligence', title: 'Payment Tenders' }
    }
  } else if (location.pathname === '/sales') {
    const tab = searchParams.get('tab')
    if (tab === 'imports') {
      currentMeta = { section: 'Retail Operations', title: 'Import Ledger' }
    }
  } else if (location.pathname === '/reports') {
    const type = searchParams.get('type')
    if (type === 'inventory') {
      currentMeta = { section: 'Financial', title: 'Inventory Valuation' }
    } else if (type === 'movements') {
      currentMeta = { section: 'Financial', title: 'Stock Movement Audit' }
    }
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
    <header className="h-14 border-b border-[#E5EAF0] bg-white px-4 sm:px-6 flex items-center justify-between sticky top-0 z-30">
      {/* Left: Mobile Toggle & Breadcrumbs */}
      <div className="flex items-center gap-3 min-w-0">
        <button
          type="button"
          onClick={onToggleSidebar}
          className="lg:hidden p-1.5 text-[#64748B] hover:text-[#172033] hover:bg-[#F5F7FA] rounded-md transition-colors cursor-pointer"
          aria-label="Open sidebar menu"
        >
          <Menu className="w-5 h-5" />
        </button>

        {/* Clean Breadcrumb */}
        <nav aria-label="Breadcrumb" className="flex items-center gap-1.5 text-xs">
          <Link to="/" className="text-[#64748B] hover:text-[#172033] font-medium hidden sm:inline">
            KadaiIQ
          </Link>
          <ChevronRight className="w-3.5 h-3.5 text-[#64748B] hidden sm:inline" />
          <span className="text-[#64748B] hidden md:inline">{currentMeta.section}</span>
          <ChevronRight className="w-3.5 h-3.5 text-[#64748B] hidden md:inline" />
          <span className="font-semibold text-[#172033] truncate">{currentMeta.title}</span>
        </nav>
      </div>

      {/* Right: Tenant Switcher, Profile & Actions */}
      <div className="flex items-center gap-3 shrink-0">
        {/* Tenant Switcher */}
        <div className="flex items-center gap-2 bg-[#F5F7FA] border border-[#E5EAF0] rounded-md px-2.5 py-1">
          <Building2 className="w-3.5 h-3.5 text-[#7C3AED] shrink-0" />
          {availableTenants.length > 1 ? (
            <div className="relative">
              <select
                value={activeTenant?.tenant_id || ''}
                onChange={handleTenantChange}
                className="appearance-none bg-transparent text-xs font-medium text-[#172033] pr-4 focus:outline-none cursor-pointer"
              >
                {availableTenants.map((t) => (
                  <option key={t.tenant_id} value={t.tenant_id}>
                    {t.tenant_name}
                  </option>
                ))}
              </select>
              <ChevronDown className="w-3 h-3 text-[#64748B] absolute right-0 top-1 pointer-events-none" />
            </div>
          ) : (
            <span className="text-xs font-semibold text-[#172033] truncate max-w-[120px] sm:max-w-[160px]">
              {activeTenant?.tenant_name || 'KadaiIQ Store'}
            </span>
          )}
          <Badge variant="outline" className="text-[10px] py-0 px-1.5 border-[#EDE9FE] text-[#5B45B5] bg-[#F5F3FF] hidden sm:inline-flex font-semibold">
            {activeTenant?.role_name || (activeTenant?.is_super_admin ? 'Super Admin' : 'Admin')}
          </Badge>
        </div>

        {/* User Account & Logout */}
        <div className="flex items-center gap-2.5 pl-2 border-l border-[#E5EAF0]">
          <div className="w-7 h-7 rounded-md bg-[#EDE9FE] text-[#5B45B5] border border-[#8B5CF6]/20 flex items-center justify-center font-bold text-xs">
            {user?.name ? user.name.charAt(0).toUpperCase() : 'U'}
          </div>
          <div className="hidden lg:block text-left">
            <span className="text-xs font-semibold text-[#172033] block leading-tight">
              {user?.name || 'Store User'}
            </span>
            <span className="text-[10px] text-[#64748B] block truncate max-w-[110px]">
              {user?.email}
            </span>
          </div>

          <button
            onClick={handleLogout}
            title="Log out"
            className="p-1.5 text-[#64748B] hover:text-[#DC2626] hover:bg-red-50 rounded-md transition-colors cursor-pointer"
            aria-label="Logout"
          >
            <LogOut className="w-4 h-4" />
          </button>
        </div>
      </div>
    </header>
  )
}
