import React from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Settings } from 'lucide-react'

export const SettingsPage: React.FC = () => (
  <div className="space-y-6">
    <div>
      <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Tenant Settings</h1>
      <p className="text-sm text-slate-500 mt-1">Configure store profile, currency, timezone, and business rules</p>
    </div>
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Settings className="w-5 h-5 text-slate-600" />
          Store Preferences
        </CardTitle>
        <CardDescription>Tenant-isolated configuration</CardDescription>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-slate-600">
          Tenant-specific settings are kept strictly isolated from other tenants.
        </p>
      </CardContent>
    </Card>
  </div>
)

