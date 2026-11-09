import React from 'react'
import { Link, useLocation } from 'react-router-dom'
import {
  LayoutDashboard,
  Package,
  Boxes,
  ShoppingCart,
  Receipt,
  TrendingUp,
  Award,
  PieChart,
  FileSpreadsheet,
  Percent,
  DownloadCloud,
  Users,
  Settings,
  ShieldCheck,
  Store,
  X,
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAuth } from '@/features/auth/AuthContext'
import { Badge } from '@/components/ui/badge'

interface NavItem {
  label: string
  path: string
  icon: React.ElementType
}

interface NavSection {
  title: string
  items: NavItem[]
}

const navSections: NavSection[] = [
  {
    title: 'EXECUTIVE',
    items: [
      { label: 'Dashboard', path: '/', icon: LayoutDashboard },
    ],
  },
  {
    title: 'RETAIL OPERATIONS',
    items: [
      { label: 'Products', path: '/products', icon: Package },
      { label: 'Inventory', path: '/inventory', icon: Boxes },
      { label: 'Sales', path: '/sales', icon: ShoppingCart },
      { label: 'Transactions', path: '/sales', icon: Receipt },
    ],
  },
  {
    title: 'BUSINESS INTELLIGENCE',
    items: [
      { label: 'Sales Analytics', path: '/analytics', icon: TrendingUp },
      { label: 'Product Performance', path: '/analytics', icon: Award },
      { label: 'Category & Brand Analytics', path: '/analytics', icon: PieChart },
    ],
  },
  {
    title: 'FINANCIAL',
    items: [
      { label: 'Reports', path: '/reports', icon: FileSpreadsheet },
      { label: 'GST Reports', path: '/gst', icon: Percent },
      { label: 'Exports', path: '/reports', icon: DownloadCloud },
    ],
  },
  {
    title: 'ORGANISATION',
    items: [
      { label: 'Users & Roles', path: '/users', icon: Users },
      { label: 'Settings', path: '/settings', icon: Settings },
    ],
  },
]

interface SidebarProps {
  isOpen?: boolean
  onClose?: () => void
}

export const Sidebar: React.FC<SidebarProps> = ({ isOpen = false, onClose }) => {
  const location = useLocation()
  const { isSuperAdmin, activeTenant } = useAuth()

  return (
    <>
      {/* Mobile Backdrop Overlay */}
      {isOpen && (
        <div
          className="fixed inset-0 z-40 bg-[#172B4D]/40 backdrop-blur-xs lg:hidden transition-opacity"
          onClick={onClose}
          aria-hidden="true"
        />
      )}

      {/* Sidebar Container */}
      <aside
        className={cn(
          'fixed inset-y-0 left-0 z-50 w-64 bg-white border-r border-[#E5EAF0] flex flex-col transition-transform duration-200 ease-in-out lg:static lg:translate-x-0 select-none shadow-card',
          isOpen ? 'translate-x-0 shadow-xl' : '-translate-x-full lg:shadow-none'
        )}
      >
        {/* Brand Header */}
        <div className="h-16 flex items-center justify-between px-5 border-b border-[#E5EAF0] bg-white">
          <Link to="/" className="flex items-center gap-2.5 group" onClick={onClose}>
            <div className="w-8 h-8 rounded-lg bg-[#172B4D] text-white flex items-center justify-center font-bold text-sm shadow-2xs group-hover:bg-[#243B5A] transition-colors">
              <span className="font-mono text-[#16A36A] font-extrabold mr-0.5">K</span>
              <span className="text-white text-xs">IQ</span>
            </div>
            <div className="flex flex-col">
              <div className="flex items-center gap-1.5">
                <span className="text-base font-bold text-[#172033] tracking-tight leading-none">KadaiIQ</span>
                <span className="w-1.5 h-1.5 rounded-full bg-[#16A36A] inline-block" />
              </div>
              <span className="text-[10px] font-medium text-[#64748B] tracking-wider uppercase mt-0.5">Retail Analytics</span>
            </div>
          </Link>

          {/* Close button on mobile */}
          {onClose && (
            <button
              onClick={onClose}
              className="lg:hidden p-1.5 text-[#64748B] hover:text-[#172033] hover:bg-[#F5F7FA] rounded-md transition-colors cursor-pointer"
              aria-label="Close menu"
            >
              <X className="w-4 h-4" />
            </button>
          )}
        </div>

        {/* Current Active Tenant Card */}
        <div className="px-3.5 pt-3 pb-1">
          <div className="bg-[#F5F7FA] border border-[#E5EAF0] rounded-lg p-2.5 flex items-center gap-2.5">
            <div className="w-7 h-7 rounded-md bg-[#172B4D] text-white flex items-center justify-center shrink-0">
              <Store className="w-3.5 h-3.5" />
            </div>
            <div className="min-w-0 flex-1">
              <span className="text-xs font-semibold text-[#172033] block truncate leading-tight">
                {activeTenant?.tenant_name || 'KadaiIQ Store'}
              </span>
              <span className="text-[10px] font-medium text-[#64748B] block truncate mt-0.5">
                {activeTenant?.role_name || (activeTenant?.is_super_admin ? 'Super Admin' : 'Store Owner')}
              </span>
            </div>
          </div>
        </div>

        {/* Navigation Sections */}
        <nav className="flex-1 px-3 py-2 space-y-3.5 overflow-y-auto">
          {navSections.map((section) => (
            <div key={section.title} className="space-y-0.5">
              <div className="px-2.5 py-1 text-[10px] font-semibold uppercase tracking-wider text-[#64748B]">
                {section.title}
              </div>
              {section.items.map((item) => {
                const Icon = item.icon
                const isActive = location.pathname === item.path
                return (
                  <Link
                    key={`${section.title}-${item.label}`}
                    to={item.path}
                    onClick={onClose}
                    className={cn(
                      'group flex items-center justify-between px-2.5 py-1.5 rounded-md text-xs font-medium transition-colors',
                      isActive
                        ? 'bg-[#172B4D] text-white font-semibold shadow-2xs'
                        : 'text-[#64748B] hover:bg-[#F5F7FA] hover:text-[#172033]'
                    )}
                  >
                    <div className="flex items-center gap-2.5 min-w-0">
                      <Icon
                        className={cn(
                          'w-4 h-4 transition-colors shrink-0',
                          isActive ? 'text-[#16A36A]' : 'text-[#64748B] group-hover:text-[#172033]'
                        )}
                      />
                      <span className="truncate">{item.label}</span>
                    </div>
                  </Link>
                )
              })}
            </div>
          ))}

          {/* Platform Administration Section (Super Admin Only) */}
          {isSuperAdmin && (
            <div className="space-y-0.5 pt-2 border-t border-[#E5EAF0]">
              <div className="px-2.5 py-1 text-[10px] font-semibold uppercase tracking-wider text-[#64748B]">
                PLATFORM ADMIN
              </div>
              <Link
                to="/super-admin"
                onClick={onClose}
                className={cn(
                  'group flex items-center justify-between px-2.5 py-1.5 rounded-md text-xs font-medium transition-colors',
                  location.pathname === '/super-admin'
                    ? 'bg-[#172B4D] text-white font-semibold'
                    : 'text-[#64748B] hover:bg-[#F5F7FA] hover:text-[#172033]'
                )}
              >
                <div className="flex items-center gap-2.5 min-w-0">
                  <ShieldCheck
                    className={cn(
                      'w-4 h-4 transition-colors shrink-0',
                      location.pathname === '/super-admin' ? 'text-[#16A36A]' : 'text-[#64748B]'
                    )}
                  />
                  <span className="truncate">Super Admin</span>
                </div>
                <Badge variant="outline" className="text-[9px] px-1.5 py-0 border-[#E5EAF0] text-[#64748B] bg-white">
                  Root
                </Badge>
              </Link>
            </div>
          )}
        </nav>

        {/* Footer Info Area */}
        <div className="p-3 border-t border-[#E5EAF0] bg-white text-xs">
          <div className="flex items-center justify-between text-[#64748B] text-[11px]">
            <span className="font-semibold text-[#172033]">KadaiIQ</span>
            <span className="font-mono text-[10px] px-1.5 py-0.2 bg-[#F5F7FA] text-[#64748B] rounded border border-[#E5EAF0]">v1.0.0</span>
          </div>
          <p className="text-[10px] text-[#64748B] mt-0.5">Retail Business Intelligence</p>
        </div>
      </aside>
    </>
  )
}
