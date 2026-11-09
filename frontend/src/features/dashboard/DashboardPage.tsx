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
  AlertCircle,
} from 'lucide-react'

export const DashboardPage: React.FC = () => {
  const { data: health, isLoading: isAppLoading, isError: isAppError } = useQuery({
    queryKey: ['system-health'],
    queryFn: apiService.getHealth,
    refetchInterval: 30000,
  })

  const { data: dbHealth, isLoading: isDbLoading, isError: isDbError } = useQuery({
    queryKey: ['database-health'],
    queryFn: apiService.getDbHealth,
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

        <div className="flex flex-wrap items-center gap-3">
          {/* API Health */}
          <div className="flex items-center space-x-2 bg-slate-50 px-3 py-2 rounded-lg border border-slate-200">
            <Activity className="w-4 h-4 text-blue-600 animate-pulse" />
            <div className="text-xs">
              <span className="text-slate-500 block text-[10px]">API Server:</span>
              {isAppLoading && <span className="font-medium text-amber-600">Connecting...</span>}
              {isAppError && <Badge variant="destructive" className="text-[10px] px-1.5 py-0">Offline</Badge>}
              {health?.success && (
                <span className="font-semibold text-emerald-600 flex items-center gap-1 text-xs">
                  <CheckCircle2 className="w-3 h-3 inline" /> v{health.data?.version}
                </span>
              )}
            </div>
          </div>

          {/* DB Health */}
          <div className="flex items-center space-x-2 bg-slate-50 px-3 py-2 rounded-lg border border-slate-200">
            <Database className="w-4 h-4 text-indigo-600" />
            <div className="text-xs">
              <span className="text-slate-500 block text-[10px]">PostgreSQL:</span>
              {isDbLoading && <span className="font-medium text-amber-600">Checking...</span>}
              {isDbError && (
                <span className="font-semibold text-amber-700 flex items-center gap-1 text-xs">
                  <AlertCircle className="w-3 h-3 inline" /> Standby
                </span>
              )}
              {dbHealth?.success && (
                <span className="font-semibold text-emerald-600 flex items-center gap-1 text-xs">
                  <CheckCircle2 className="w-3 h-3 inline" /> {dbHealth.data?.latency_ms}ms
                </span>
              )}
            </div>
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

      {/* Architecture & Multi-Tenant Isolation Status */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg flex items-center gap-2">
              <Database className="w-5 h-5 text-blue-600" />
              Database & Multi-Tenant Models (Phase 2)
            </CardTitle>
            <CardDescription>
              PostgreSQL schema with versioned migrations and repository tenant boundary checks
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3 text-sm text-slate-600">
            <div className="flex items-center justify-between border-b border-slate-100 pb-2">
              <span className="font-medium text-slate-700">Schema Migration Engine</span>
              <Badge variant="outline">000001_initial_schema</Badge>
            </div>
            <div className="flex items-center justify-between border-b border-slate-100 pb-2">
              <span className="font-medium text-slate-700">Tenants & Memberships</span>
              <Badge variant="success">Composite Unique Indexes</Badge>
            </div>
            <div className="flex items-center justify-between border-b border-slate-100 pb-2">
              <span className="font-medium text-slate-700">RBAC Foundation</span>
              <Badge variant="outline">Super Admin, Owner, Admin, Manager, Staff</Badge>
            </div>
            <div className="flex items-center justify-between">
              <span className="font-medium text-slate-700">Tenant Isolation Scope</span>
              <Badge variant="success">BaseTenantRepository Verified</Badge>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-lg flex items-center gap-2">
              <ShieldAlert className="w-5 h-5 text-amber-600" />
              Connection Health & Diagnostics
            </CardTitle>
            <CardDescription>
              Dual health endpoints providing live readiness & database availability probes
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3 text-sm text-slate-600">
            <div className="flex items-center justify-between border-b border-slate-100 pb-2">
              <span className="font-medium text-slate-700">Service Probe</span>
              <code className="text-xs bg-slate-100 px-2 py-0.5 rounded">GET /api/v1/health</code>
            </div>
            <div className="flex items-center justify-between border-b border-slate-100 pb-2">
              <span className="font-medium text-slate-700">Database Probe</span>
              <code className="text-xs bg-slate-100 px-2 py-0.5 rounded">GET /api/v1/health/db</code>
            </div>
            <div className="flex items-center justify-between border-b border-slate-100 pb-2">
              <span className="font-medium text-slate-700">Connection Pooling</span>
              <span className="text-xs text-slate-500">MaxOpen: 25 | MaxIdle: 10</span>
            </div>
            <div className="flex items-center justify-between">
              <span className="font-medium text-slate-700">Status Code Handling</span>
              <span className="text-xs text-slate-500">200 OK / 503 Unavailable</span>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
