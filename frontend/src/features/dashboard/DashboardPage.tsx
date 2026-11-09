import React from 'react'
import { useQuery } from '@tanstack/react-query'
import { apiService } from '@/services/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
  Activity,
  Layers,
  Database,
  Boxes,
  TrendingUp,
  Receipt,
  CheckCircle2,
  AlertCircle,
  Package,
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

  const { data: inventoryResp } = useQuery({
    queryKey: ['inventory-summary-dashboard'],
    queryFn: apiService.getInventorySummary,
  })

  const { data: analyticsResp } = useQuery({
    queryKey: ['analytics-overview-dashboard'],
    queryFn: () => apiService.getAnalyticsOverview({ range: '30d' }),
  })

  const invSummary = inventoryResp?.data
  const analytics = analyticsResp?.data

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
            <CardTitle className="text-sm font-medium text-slate-600">30-Day Gross Revenue</CardTitle>
            <TrendingUp className="w-4 h-4 text-emerald-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900 font-mono">
              ₹{Number(analytics?.gross_sales || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
            </div>
            <p className="text-xs text-slate-500 mt-1">
              {analytics?.total_bills || 0} Bills • {Number(analytics?.profit_margin_pct || 0).toFixed(1)}% Est. Margin
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-slate-600">Active SKUs</CardTitle>
            <Boxes className="w-4 h-4 text-blue-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900">{invSummary?.total_variants || 0}</div>
            <p className="text-xs text-slate-500 mt-1">{invSummary?.total_products || 0} product models</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-slate-600">Stock Units on Hand</CardTitle>
            <Layers className="w-4 h-4 text-indigo-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900">{invSummary?.total_stock_units || 0} Units</div>
            <p className="text-xs text-slate-500 mt-1">
              Valuation: ₹{Number(invSummary?.stock_valuation || 0).toLocaleString('en-IN')}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-slate-600">Low & Out of Stock</CardTitle>
            <Receipt className="w-4 h-4 text-purple-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900">
              {invSummary?.low_stock_count || 0} / {invSummary?.out_of_stock_count || 0}
            </div>
            <p className="text-xs text-slate-500 mt-1">Low Stock / Out of Stock</p>
          </CardContent>
        </Card>
      </div>

      {/* Product & Variant Footwear Architecture */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle className="text-base flex items-center gap-2">
              <Package className="w-5 h-5 text-blue-600" />
              Product & Variant Management (Phase 4)
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3 text-sm text-slate-600">
            <div className="flex items-center justify-between border-b pb-2">
              <span>Footwear Variants</span>
              <Badge variant="outline">Size / Colour / SKU Matrix</Badge>
            </div>
            <div className="flex items-center justify-between border-b pb-2">
              <span>Barcode Support</span>
              <Badge variant="success">EAN-13 / Scanner Ready</Badge>
            </div>
            <div className="flex items-center justify-between border-b pb-2">
              <span>GST Compliance</span>
              <Badge variant="outline">Footwear Standard (12%)</Badge>
            </div>
            <div className="flex items-center justify-between">
              <span>Multi-Tenant Boundaries</span>
              <Badge variant="success">Tenant Scope Enforced</Badge>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-base flex items-center gap-2">
              <Boxes className="w-5 h-5 text-indigo-600" />
              Stock Movement Ledger
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3 text-sm text-slate-600">
            <div className="flex items-center justify-between border-b pb-2">
              <span>Movement Types</span>
              <span className="text-xs font-mono text-slate-600">OPENING, PURCHASE, SALE, RETURN</span>
            </div>
            <div className="flex items-center justify-between border-b pb-2">
              <span>Audit Guarantee</span>
              <span className="text-xs text-emerald-600 font-semibold">Zero silent updates</span>
            </div>
            <div className="flex items-center justify-between border-b pb-2">
              <span>Concurrency Control</span>
              <span className="text-xs text-slate-600">Row-level locking on adjustment</span>
            </div>
            <div className="flex items-center justify-between">
              <span>Overselling Defense</span>
              <span className="text-xs text-emerald-600 font-semibold">Validated deduction bounds</span>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
