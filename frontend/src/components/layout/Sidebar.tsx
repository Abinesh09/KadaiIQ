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
      { label: 'Executive Overview', path: '/', icon: LayoutDashboard },
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
          className="fixed inset-0 z-40 bg-slate-900/60 backdrop-blur-xs lg:hidden transition-opacity"
          onClick={onClose}
          aria-hidden="true"
        />
      )}

      {/* Sidebar Container */}
      <aside
        className={cn(
          'fixed inset-y-0 left-0 z-50 w-72 bg-white border-r border-slate-200/90 flex flex-col transition-transform duration-300 ease-in-out lg:static lg:translate-x-0 select-none shadow-[1px_0_4px_0_rgba(15,23,42,0.03)]',
          isOpen ? 'translate-x-0 shadow-2xl' : '-translate-x-full lg:shadow-none'
        )}
      >
        {/* Brand Header */}
        <div className="h-18 flex items-center justify-between px-6 border-b border-slate-200/80 bg-slate-50/50">
          <Link to="/" className="flex items-center gap-3.5 group" onClick={onClose}>
            <div className="w-10 h-10 rounded-xl bg-slate-900 text-white flex items-center justify-center font-extrabold text-lg shadow-md ring-2 ring-slate-800 group-hover:scale-105 transition-transform">
              <span className="font-mono text-emerald-400 font-black tracking-tighter">K</span>
              <span className="text-white text-xs font-bold tracking-tight">IQ</span>
            </div>
            <div className="flex flex-col">
              <div className="flex items-center gap-1.5">
                <span className="text-lg font-black text-slate-950 tracking-tight leading-none">KadaiIQ</span>
                <span className="w-2 h-2 rounded-full bg-emerald-500 inline-block animate-pulse ring-2 ring-emerald-200" />
              </div>
              <span className="text-[10px] font-bold text-blue-600 tracking-wider uppercase mt-1">Retail SaaS Platform</span>
            </div>
          </Link>

          {/* Close button on mobile */}
          {onClose && (
            <button
              onClick={onClose}
              className="lg:hidden p-2 text-slate-400 hover:text-slate-800 hover:bg-slate-100 rounded-xl transition-colors cursor-pointer"
              aria-label="Close menu"
            >
              <X className="w-5 h-5" />
            </button>
          )}
        </div>

        {/* Current Active Tenant Card */}
        <div className="px-5 pt-4 pb-2">
          <div className="bg-slate-50 border border-slate-200/80 rounded-xl p-3 flex items-center gap-3 shadow-2xs">
            <div className="w-8 h-8 rounded-lg bg-blue-600 text-white flex items-center justify-center shrink-0 shadow-xs">
              <Store className="w-4 h-4" />
            </div>
            <div className="min-w-0 flex-1">
              <span className="text-xs font-bold text-slate-900 block truncate leading-tight">
                {activeTenant?.tenant_name || 'KadaiIQ Store'}
              </span>
              <span className="text-[10px] font-semibold text-slate-400 block truncate mt-0.5">
                {activeTenant?.role_name || (activeTenant?.is_super_admin ? 'Super Admin' : 'Store Owner')}
              </span>
            </div>
          </div>
        </div>

        {/* Navigation Sections */}
        <nav className="flex-1 px-4 py-3 space-y-5 overflow-y-auto">
          {navSections.map((section) => (
            <div key={section.title} className="space-y-1">
              <div className="px-3 py-1 text-[10px] font-extrabold uppercase tracking-widest text-slate-400">
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
                      'group flex items-center justify-between px-3.5 py-2.5 rounded-xl text-xs font-semibold transition-all duration-150',
                      isActive
                        ? 'bg-slate-900 text-white shadow-md shadow-slate-900/10'
                        : 'text-slate-600 hover:bg-slate-100 hover:text-slate-950'
                    )}
                  >
                    <div className="flex items-center gap-3 min-w-0">
                      <Icon
                        className={cn(
                          'w-4 h-4 transition-colors shrink-0',
                          isActive ? 'text-emerald-400' : 'text-slate-400 group-hover:text-slate-700'
                        )}
                      />
                      <span className="truncate">{item.label}</span>
                    </div>
                    {item.badge && (
                      <span
                        className={cn(
                          'text-[10px] font-bold px-2 py-0.5 rounded-full flex items-center gap-1',
                          isActive
                            ? 'bg-emerald-500/20 text-emerald-300 border border-emerald-400/30'
                            : 'bg-emerald-50 text-emerald-700 border border-emerald-200'
                        )}
                      >
                        <Sparkles className="w-2.5 h-2.5" />
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
            <div className="space-y-1 pt-3 border-t border-slate-200/80">
              <div className="px-3 py-1 text-[10px] font-extrabold uppercase tracking-widest text-purple-600">
                Platform Admin
              </div>
              <Link
                to="/super-admin"
                onClick={onClose}
                className={cn(
                  'group flex items-center justify-between px-3.5 py-2.5 rounded-xl text-xs font-semibold transition-all duration-150',
                  location.pathname === '/super-admin'
                    ? 'bg-purple-900 text-white shadow-md'
                    : 'text-slate-600 hover:bg-purple-50 hover:text-purple-900'
                )}
              >
                <div className="flex items-center gap-3 min-w-0">
                  <ShieldCheck
                    className={cn(
                      'w-4 h-4 transition-colors shrink-0',
                      location.pathname === '/super-admin' ? 'text-purple-300' : 'text-purple-500'
                    )}
                  />
                  <span className="truncate">Super Admin Portal</span>
                </div>
                <Badge variant="outline" className="text-[9px] px-1.5 py-0 border-purple-300 text-purple-700 bg-purple-50">
                  Root
                </Badge>
              </Link>
            </div>
          )}
        </nav>

        {/* Footer Info Area */}
        <div className="p-4 border-t border-slate-200/80 bg-slate-50 text-xs">
          <div className="flex items-center justify-between text-slate-500 text-[11px]">
            <span className="font-bold text-slate-700">KadaiIQ SaaS</span>
            <span className="font-mono text-[10px] font-semibold px-2 py-0.5 bg-slate-200 text-slate-700 rounded-md">v1.0.0</span>
          </div>
          <p className="text-[10px] text-slate-400 mt-1">Shared PostgreSQL • Multi-Tenant Isolation</p>
        </div>
      </aside>
    </>
  )
}
