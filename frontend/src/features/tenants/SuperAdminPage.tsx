import React from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { ShieldCheck } from 'lucide-react'

export const SuperAdminPage: React.FC = () => (
  <div className="space-y-6">
    <div>
      <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Super Admin Portal</h1>
      <p className="text-sm text-slate-500 mt-1">Platform-level tenant orchestration, subscription plans, and feature flags</p>
    </div>
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <ShieldCheck className="w-5 h-5 text-purple-600" />
          Tenant Management & Provisioning
        </CardTitle>
        <CardDescription>Scheduled for Phase 3 implementation</CardDescription>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-slate-600">
          Super Admin allows provisioning new retail organisations, assigning subscription tiers, enabling features, and monitoring cross-tenant platform health.
        </p>
      </CardContent>
    </Card>
  </div>
)

