import React from 'react'
import { Link, useLocation } from 'react-router-dom'
import {
  LayoutDashboard,
  Package,
  Boxes,
  TrendingUp,
  FileSpreadsheet,
  Settings,
  ShieldCheck,
  Percent,
} from 'lucide-react'
import { cn } from '@/lib/utils'

const navItems = [
  { label: 'Dashboard', path: '/', icon: LayoutDashboard },
  { label: 'Products', path: '/products', icon: Package },
  { label: 'Inventory', path: '/inventory', icon: Boxes },
  { label: 'Sales & POS', path: '/sales', icon: TrendingUp },
  { label: 'Analytics', path: '/analytics', icon: TrendingUp },
  { label: 'Reports', path: '/reports', icon: FileSpreadsheet },
  { label: 'GST Config', path: '/gst', icon: Percent },
  { label: 'Settings', path: '/settings', icon: Settings },
]

export const Sidebar: React.FC = () => {
  const location = useLocation()

  return (
    <aside className="w-64 border-r border-slate-200 bg-white flex flex-col h-screen">
      <div className="h-16 flex items-center px-6 border-b border-slate-200">
        <Link to="/" className="flex items-center space-x-2">
          <div className="w-8 h-8 rounded-lg bg-blue-600 flex items-center justify-center text-white font-bold text-lg">
            K
          </div>
          <div>
            <span className="text-lg font-bold text-slate-900 tracking-tight">KadaiIQ</span>
            <span className="text-[10px] text-blue-600 font-semibold block -mt-1 tracking-wider uppercase">Retail SaaS</span>
          </div>
        </Link>
      </div>

      <nav className="flex-1 px-4 py-4 space-y-1 overflow-y-auto">
        <div className="px-2 py-1 text-xs font-semibold uppercase tracking-wider text-slate-400">
          Tenant Workspace
        </div>
        {navItems.map((item) => {
          const Icon = item.icon
          const isActive = location.pathname === item.path
          return (
            <Link
              key={item.path}
              to={item.path}
              className={cn(
                'flex items-center space-x-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors',
                isActive
                  ? 'bg-blue-50 text-blue-700'
                  : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'
              )}
            >
              <Icon className={cn('w-4 h-4', isActive ? 'text-blue-600' : 'text-slate-400')} />
              <span>{item.label}</span>
            </Link>
          )
        })}

        <div className="pt-4 px-2 py-1 text-xs font-semibold uppercase tracking-wider text-slate-400">
          Administration
        </div>
        <Link
          to="/super-admin"
          className="flex items-center space-x-3 px-3 py-2 rounded-lg text-sm font-medium text-slate-600 hover:bg-slate-50 hover:text-slate-900 transition-colors"
        >
          <ShieldCheck className="w-4 h-4 text-purple-600" />
          <span>Super Admin</span>
        </Link>
      </nav>

      <div className="p-4 border-t border-slate-200 bg-slate-50 text-xs text-slate-500">
        <p className="font-semibold text-slate-700">KadaiIQ v1.0.0</p>
        <p>Retail Multi-Tenant Platform</p>
      </div>
    </aside>
  )
}

