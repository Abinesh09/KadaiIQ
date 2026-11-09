import React from 'react'
import { useQuery } from '@tanstack/react-query'
import { apiService } from '@/services/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
  Activity,
  Layers,
  Database,
  ShieldAlert,
  Boxes,
  TrendingUp,
  Receipt,
  CheckCircle2,
} from 'lucide-react'

export const DashboardPage: React.FC = () => {
  const { data: health, isLoading, isError } = useQuery({
    queryKey: ['system-health'],
    queryFn: apiService.getHealth,
    refetchInterval: 30000,
  })

  return (
    <div className="space-y-6">
      {/* Top Banner & Health Status */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white p-6 rounded-xl border border-slate-200">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">KadaiIQ Executive Overview</h1>
          <p className="text-sm text-slate-500 mt-1">
            Multi-Tenant Retail Management & Advanced Sales Analytics SaaS
          </p>
        </div>

        <div className="flex items-center space-x-3 bg-slate-50 px-4 py-2 rounded-lg border border-slate-200">
          <Activity className="w-5 h-5 text-blue-600 animate-pulse" />
          <div className="text-xs">
            <span className="text-slate-500 block">Backend Fiber API:</span>
            {isLoading && <span className="font-semibold text-amber-600">Connecting...</span>}
            {isError && <Badge variant="destructive">Disconnected</Badge>}
            {health?.success && (
              <span className="font-semibold text-emerald-600 flex items-center gap-1">
                <CheckCircle2 className="w-3.5 h-3.5 inline" /> Operational (v{health.data?.version})
              </span>
            )}
          </div>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-slate-600">Total Gross Sales</CardTitle>
            <TrendingUp className="w-4 h-4 text-emerald-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900">₹0.00</div>
            <p className="text-xs text-slate-500 mt-1">Ready for POS Import</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-slate-600">Active SKUs</CardTitle>
            <Boxes className="w-4 h-4 text-blue-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900">0</div>
            <p className="text-xs text-slate-500 mt-1">Product variant tracking</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-slate-600">Stock Units</CardTitle>
            <Layers className="w-4 h-4 text-indigo-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900">0 Units</div>
            <p className="text-xs text-slate-500 mt-1">Stock movement ledger</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-slate-600">Tax Collected (GST)</CardTitle>
            <Receipt className="w-4 h-4 text-purple-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900">₹0.00</div>
            <p className="text-xs text-slate-500 mt-1">CGST / SGST / IGST</p>
          </CardContent>
        </Card>
      </div>

      {/* Architecture Verification Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg flex items-center gap-2">
              <Database className="w-5 h-5 text-blue-600" />
              Multi-Tenant Architecture
            </CardTitle>
            <CardDescription>
              Backend Go Fiber + GORM PostgreSQL with strict tenant boundary isolation
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3 text-sm text-slate-600">
            <div className="flex items-center justify-between border-b border-slate-100 pb-2">
              <span className="font-medium text-slate-700">Database Layer</span>
              <Badge variant="outline">PostgreSQL 16</Badge>
            </div>
            <div className="flex items-center justify-between border-b border-slate-100 pb-2">
              <span className="font-medium text-slate-700">API Gateway</span>
              <Badge variant="outline">Go Fiber v2 (/api/v1)</Badge>
            </div>
            <div className="flex items-center justify-between border-b border-slate-100 pb-2">
              <span className="font-medium text-slate-700">Tenant Isolation</span>
              <Badge variant="success">Tenant Context Enforced</Badge>
            </div>
            <div className="flex items-center justify-between">
              <span className="font-medium text-slate-700">Security Model</span>
              <Badge variant="outline">JWT Access & Refresh Tokens</Badge>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-lg flex items-center gap-2">
              <ShieldAlert className="w-5 h-5 text-amber-600" />
              Core Capabilities in Phase 1
            </CardTitle>
            <CardDescription>
              Monorepo foundation, Docker containerization, and service communication
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3 text-sm text-slate-600">
            <div className="flex items-center justify-between border-b border-slate-100 pb-2">
              <span className="font-medium text-slate-700">Frontend Port</span>
              <Badge variant="outline">7070 (Configured)</Badge>
            </div>
            <div className="flex items-center justify-between border-b border-slate-100 pb-2">
              <span className="font-medium text-slate-700">Backend Port</span>
              <Badge variant="outline">8080 (Go Fiber)</Badge>
            </div>
            <div className="flex items-center justify-between border-b border-slate-100 pb-2">
              <span className="font-medium text-slate-700">Component System</span>
              <Badge variant="outline">shadcn/ui + Tailwind CSS</Badge>
            </div>
            <div className="flex items-center justify-between">
              <span className="font-medium text-slate-700">State Management</span>
              <Badge variant="outline">TanStack Query v5</Badge>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

