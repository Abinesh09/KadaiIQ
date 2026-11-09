import React from 'react'
import { Link, useLocation } from 'react-router-dom'
import {
  LayoutDashboard,
  Package,
  Boxes,
  ShoppingCart,
  TrendingUp,
  FileSpreadsheet,
  Settings,
  ShieldCheck,
  Percent,
  Users,
  Store,
  X,
  Sparkles,
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAuth } from '@/features/auth/AuthContext'
import { Badge } from '@/components/ui/badge'

interface NavItem {
  label: string
  path: string
  icon: React.ElementType
  badge?: string
}

interface NavSection {
  title: string
  items: NavItem[]
}

const navSections: NavSection[] = [
  {
    title: 'Executive',
    items: [
      { label: 'Dashboard', path: '/', icon: LayoutDashboard },
    ],
  },
  {
    title: 'Retail Operations',
    items: [
      { label: 'Product Catalogue', path: '/products', icon: Package },
      { label: 'Inventory & Stock', path: '/inventory', icon: Boxes },
      { label: 'Sales & POS Records', path: '/sales', icon: ShoppingCart },
    ],
  },
  {
    title: 'Intelligence & Tax',
    items: [
      { label: 'Sales Analytics', path: '/analytics', icon: TrendingUp, badge: 'Insights' },
      { label: 'Reports & Exports', path: '/reports', icon: FileSpreadsheet },
      { label: 'GST Configuration', path: '/gst', icon: Percent },
    ],
  },
  {
    title: 'Organisation',
    items: [
      { label: 'Staff & Roles', path: '/users', icon: Users },
      { label: 'Store Settings', path: '/settings', icon: Settings },
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
          className="fixed inset-0 z-40 bg-slate-900/50 backdrop-blur-xs lg:hidden transition-opacity"
          onClick={onClose}
          aria-hidden="true"
        />
      )}

      {/* Sidebar Container */}
      <aside
        className={cn(
          'fixed inset-y-0 left-0 z-50 w-64 bg-white border-r border-slate-200/90 flex flex-col transition-transform duration-300 ease-in-out lg:static lg:translate-x-0',
          isOpen ? 'translate-x-0 shadow-2xl' : '-translate-x-full lg:shadow-none'
        )}
      >
        {/* Brand Header */}
        <div className="h-16 flex items-center justify-between px-5 border-b border-slate-200/80 bg-slate-50/60">
          <Link to="/" className="flex items-center gap-3 group" onClick={onClose}>
            <div className="w-9 h-9 rounded-xl bg-gradient-to-tr from-slate-900 via-blue-900 to-indigo-700 flex items-center justify-center text-white font-bold text-base shadow-sm ring-1 ring-blue-500/20 group-hover:scale-105 transition-transform">
              <span className="font-mono text-emerald-400 font-extrabold mr-0.5">K</span>
              <span className="text-white text-xs">IQ</span>
            </div>
            <div className="flex flex-col">
              <div className="flex items-center gap-1.5">
                <span className="text-base font-extrabold text-slate-900 tracking-tight leading-none">KadaiIQ</span>
                <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 inline-block animate-pulse" />
              </div>
              <span className="text-[10px] font-semibold text-blue-700 tracking-wider uppercase mt-1">Retail SaaS</span>
            </div>
          </Link>

          {/* Close button on mobile */}
          {onClose && (
            <button
              onClick={onClose}
              className="lg:hidden p-1.5 text-slate-400 hover:text-slate-700 hover:bg-slate-100 rounded-lg transition-colors"
              aria-label="Close menu"
            >
              <X className="w-5 h-5" />
            </button>
          )}
        </div>

        {/* Current Active Tenant Mini-Card */}
        <div className="px-4 pt-3.5 pb-2">
          <div className="bg-slate-50/80 border border-slate-200/70 rounded-xl p-2.5 flex items-center gap-2.5">
            <div className="w-7 h-7 rounded-lg bg-blue-50 border border-blue-100 flex items-center justify-center text-blue-600 shrink-0">
              <Store className="w-3.5 h-3.5" />
            </div>
            <div className="min-w-0 flex-1">
              <span className="text-xs font-semibold text-slate-800 block truncate leading-tight">
                {activeTenant?.tenant_name || 'KadaiIQ Store'}
              </span>
              <span className="text-[10px] font-medium text-slate-400 block truncate">
                {activeTenant?.role_name || (activeTenant?.is_super_admin ? 'Super Admin' : 'Active Store')}
              </span>
            </div>
          </div>
        </div>

        {/* Navigation Sections */}
        <nav className="flex-1 px-3 py-2 space-y-4 overflow-y-auto">
          {navSections.map((section) => (
            <div key={section.title} className="space-y-1">
              <div className="px-3 py-1 text-[10px] font-bold uppercase tracking-wider text-slate-400">
                {section.title}
              </div>
              {section.items.map((item) => {
                const Icon = item.icon
                const isActive = location.pathname === item.path
                return (
                  <Link
                    key={item.path}
                    to={item.path}
                    onClick={onClose}
                    className={cn(
                      'group flex items-center justify-between px-3 py-2 rounded-lg text-xs font-medium transition-all duration-150',
                      isActive
                        ? 'bg-blue-50/90 text-blue-800 font-semibold shadow-xs ring-1 ring-blue-500/20'
                        : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'
                    )}
                  >
                    <div className="flex items-center gap-2.5 min-w-0">
                      <Icon
                        className={cn(
                          'w-4 h-4 transition-colors shrink-0',
                          isActive ? 'text-blue-600' : 'text-slate-400 group-hover:text-slate-600'
                        )}
                      />
                      <span className="truncate">{item.label}</span>
                    </div>
                    {item.badge && (
                      <span className="text-[10px] font-semibold px-1.5 py-0.5 rounded-full bg-emerald-100 text-emerald-800 flex items-center gap-0.5">
                        <Sparkles className="w-2.5 h-2.5 text-emerald-600" />
                        {item.badge}
                      </span>
                    )}
                  </Link>
                )
              })}
            </div>
          ))}

          {/* Platform Administration Section (Super Admin Only) */}
          {isSuperAdmin && (
            <div className="space-y-1 pt-2 border-t border-slate-100">
              <div className="px-3 py-1 text-[10px] font-bold uppercase tracking-wider text-purple-600">
                Platform Admin
              </div>
              <Link
                to="/super-admin"
                onClick={onClose}
                className={cn(
                  'group flex items-center justify-between px-3 py-2 rounded-lg text-xs font-medium transition-all duration-150',
                  location.pathname === '/super-admin'
                    ? 'bg-purple-50 text-purple-900 font-semibold ring-1 ring-purple-500/20'
                    : 'text-slate-600 hover:bg-purple-50/50 hover:text-purple-900'
                )}
              >
                <div className="flex items-center gap-2.5 min-w-0">
                  <ShieldCheck
                    className={cn(
                      'w-4 h-4 transition-colors shrink-0',
                      location.pathname === '/super-admin' ? 'text-purple-600' : 'text-purple-400'
                    )}
                  />
                  <span className="truncate">Super Admin Portal</span>
                </div>
                <Badge variant="outline" className="text-[9px] px-1 py-0 border-purple-200 text-purple-700 bg-purple-50">
                  Root
                </Badge>
              </Link>
            </div>
          )}
        </nav>

        {/* Footer Info Area */}
        <div className="p-3.5 border-t border-slate-200/80 bg-slate-50/70 text-xs">
          <div className="flex items-center justify-between text-slate-500 text-[11px]">
            <span className="font-semibold text-slate-700">KadaiIQ SaaS</span>
            <span className="font-mono text-[10px] px-1.5 py-0.5 bg-slate-200/60 rounded text-slate-600">v1.0.0</span>
          </div>
          <p className="text-[10px] text-slate-400 mt-0.5">Isolated Multi-Tenant Security</p>
        </div>
      </aside>
    </>
  )
}
