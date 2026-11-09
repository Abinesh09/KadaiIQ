import React from 'react'
import { Bell, User } from 'lucide-react'
import { Badge } from '@/components/ui/badge'

interface HeaderProps {
  tenantName?: string
}

export const Header: React.FC<HeaderProps> = ({ tenantName = 'Demo Retail Organiser' }) => {
  return (
    <header className="h-16 border-b border-slate-200 bg-white px-6 flex items-center justify-between">
      <div className="flex items-center space-x-3">
        <h2 className="text-sm font-semibold text-slate-800">{tenantName}</h2>
        <Badge variant="outline" className="text-xs">Multi-Tenant Store</Badge>
      </div>

      <div className="flex items-center space-x-4">
        <button
          type="button"
          aria-label="Notifications"
          className="p-2 text-slate-500 hover:text-slate-700 hover:bg-slate-100 rounded-full transition-colors"
        >
          <Bell className="w-5 h-5" />
        </button>

        <div className="flex items-center space-x-2 pl-2 border-l border-slate-200">
          <div className="w-8 h-8 rounded-full bg-blue-100 text-blue-700 flex items-center justify-center font-medium text-sm">
            <User className="w-4 h-4" />
          </div>
          <span className="text-sm font-medium text-slate-700">Store Admin</span>
        </div>
      </div>
    </header>
  )
}

