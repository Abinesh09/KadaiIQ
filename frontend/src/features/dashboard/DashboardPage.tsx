import React, { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { apiService } from '@/services/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
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
  FileSpreadsheet,
  ArrowUpRight,
  ArrowDownRight,
  DollarSign,
  Calendar,
} from 'lucide-react'
import {
  ResponsiveContainer,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
} from 'recharts'

export const DashboardPage: React.FC = () => {
  const [period, setPeriod] = useState<string>('today')
  const [isExporting, setIsExporting] = useState<boolean>(false)

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

  const { data: perfResp, isLoading: isPerfLoading } = useQuery({
    queryKey: ['sales-performance-dashboard', period],
    queryFn: () => apiService.getSalesPerformance({ range: period }),
  })

  const { data: trendsResp } = useQuery({
    queryKey: ['sales-trends-dashboard', period],
    queryFn: () => apiService.getSalesTrends({ range: period === 'today' || period === 'yesterday' ? '7d' : period }),
  })

  const invSummary = inventoryResp?.data
  const perf = perfResp?.data
  const trends = trendsResp?.data || []

  const chartData = trends.map((t) => ({
    date: t.period,
    sales: Number(t.sales || 0),
    netSales: Number(t.net_sales || 0),
  }))

  const salesGrowth = Number(perf?.sales_growth_pct || 0)

  const handleExportExcel = async () => {
    try {
      setIsExporting(true)
      const blob = await apiService.downloadExcelReport({ range: period })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `KadaiIQ_Sales_Overview_${period}_${new Date().toISOString().slice(0, 10)}.xlsx`
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
    } catch (err) {
      console.error('Failed to export Excel report:', err)
    } finally {
      setIsExporting(false)
    }
  }

  return (
    <div className="space-y-6">
      {/* Top Banner & Health Status */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white p-6 rounded-xl border border-slate-200 shadow-sm">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">KadaiIQ Executive Overview</h1>
          <p className="text-sm text-slate-500 mt-1">
            Multi-Tenant Retail Management & Advanced Sales Analytics SaaS
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-3">
          {/* Export Excel Button */}
          <Button
            variant="default"
            size="sm"
            onClick={handleExportExcel}
            disabled={isExporting}
            className="bg-emerald-600 hover:bg-emerald-700 text-white shadow-sm flex items-center gap-1.5"
          >
            <FileSpreadsheet className="w-4 h-4" />
            {isExporting ? 'Exporting...' : 'Export Report (.xlsx)'}
          </Button>

          {/* API Health */}
          <div className="flex items-center space-x-2 bg-slate-50 px-3 py-1.5 rounded-lg border border-slate-200">
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
          <div className="flex items-center space-x-2 bg-slate-50 px-3 py-1.5 rounded-lg border border-slate-200">
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

      {/* Sales Dashboard Section */}
      <div className="bg-white p-6 rounded-xl border border-slate-200 shadow-sm space-y-6">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-100 pb-4">
          <div>
            <h2 className="text-lg font-bold text-slate-900 tracking-tight flex items-center gap-2">
              <TrendingUp className="w-5 h-5 text-emerald-600" />
              Sales Overview
            </h2>
            <p className="text-xs text-slate-500 mt-0.5">
              Select time horizon to view earnings, products sold, completed bills, and average bill value.
            </p>
          </div>

          {/* Horizon Switcher */}
          <div className="flex items-center bg-slate-100 p-1 rounded-lg">
            {[
              { id: 'today', label: 'Today' },
              { id: 'yesterday', label: 'Yesterday' },
              { id: 'this_month', label: 'This Month' },
              { id: 'this_year', label: 'This Year' },
            ].map((btn) => (
              <button
                key={btn.id}
                onClick={() => setPeriod(btn.id)}
                className={`px-3 py-1.5 rounded-md text-xs font-medium transition-all ${
                  period === btn.id
                    ? 'bg-white text-slate-900 shadow-sm font-semibold'
                    : 'text-slate-600 hover:text-slate-900'
                }`}
              >
                {btn.label}
              </button>
            ))}
          </div>
        </div>

        {/* 4 Core Retail Sales Cards */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {/* 1. Total Sales Amount */}
          <div className="p-4 rounded-xl border border-slate-100 bg-slate-50/50">
            <span className="text-xs text-slate-500 font-medium block">Total Sales Amount</span>
            <div className="text-2xl font-bold text-slate-900 font-mono mt-1">
              {isPerfLoading ? '...' : `₹${Number(perf?.gross_sales || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}`}
            </div>
            <div className="flex items-center gap-1.5 mt-2 text-xs">
              {salesGrowth >= 0 ? (
                <span className="inline-flex items-center text-emerald-600 font-semibold">
                  <ArrowUpRight className="w-3.5 h-3.5 mr-0.5" />
                  +{salesGrowth}%
                </span>
              ) : (
                <span className="inline-flex items-center text-rose-600 font-semibold">
                  <ArrowDownRight className="w-3.5 h-3.5 mr-0.5" />
                  {salesGrowth}%
                </span>
              )}
              <span className="text-slate-400">vs prior period</span>
            </div>
          </div>

          {/* 2. Products Sold */}
          <div className="p-4 rounded-xl border border-slate-100 bg-slate-50/50">
            <span className="text-xs text-slate-500 font-medium block">Products Sold</span>
            <div className="text-2xl font-bold text-slate-900 font-mono mt-1">
              {isPerfLoading ? '...' : `${perf?.total_units_sold || 0}`}
            </div>
            <p className="text-xs text-slate-500 mt-2">Units moved across bills</p>
          </div>

          {/* 3. Completed Bills */}
          <div className="p-4 rounded-xl border border-slate-100 bg-slate-50/50">
            <span className="text-xs text-slate-500 font-medium block">Completed Bills</span>
            <div className="text-2xl font-bold text-slate-900 font-mono mt-1">
              {isPerfLoading ? '...' : `${perf?.total_bills || 0}`}
            </div>
            <p className="text-xs text-slate-500 mt-2">Invoices settled & finalized</p>
          </div>

          {/* 4. Average Bill Value */}
          <div className="p-4 rounded-xl border border-slate-100 bg-slate-50/50">
            <span className="text-xs text-slate-500 font-medium block">Average Bill Value</span>
            <div className="text-2xl font-bold text-slate-900 font-mono mt-1">
              {isPerfLoading ? '...' : `₹${Number(perf?.average_bill_value || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}`}
            </div>
            <p className="text-xs text-slate-500 mt-2">Net revenue per customer bill</p>
          </div>
        </div>

        {/* Financial Reconciliations (Net Sales, COGS, Profit) */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 pt-2 border-t border-slate-100">
          <div className="flex items-center justify-between p-3 rounded-lg bg-emerald-50/40 border border-emerald-100">
            <div>
              <span className="text-[11px] text-emerald-800 font-medium block">Net Sales (After Returns & Discounts)</span>
              <span className="text-lg font-bold text-emerald-900 font-mono">
                ₹{Number(perf?.net_sales || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
              </span>
            </div>
            <DollarSign className="w-5 h-5 text-emerald-600" />
          </div>

          <div className="flex items-center justify-between p-3 rounded-lg bg-blue-50/40 border border-blue-100">
            <div>
              <span className="text-[11px] text-blue-800 font-medium block">Gross Profit</span>
              <span className="text-lg font-bold text-blue-900 font-mono">
                ₹{Number(perf?.gross_profit || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
              </span>
            </div>
            <Badge className="bg-blue-600 text-white text-[10px]">
              {Number(perf?.profit_margin_pct || 0).toFixed(1)}% Margin
            </Badge>
          </div>

          <div className="flex items-center justify-between p-3 rounded-lg bg-indigo-50/40 border border-indigo-100">
            <div>
              <span className="text-[11px] text-indigo-800 font-medium block">Average Daily Sales</span>
              <span className="text-lg font-bold text-indigo-900 font-mono">
                ₹{Number(perf?.average_daily_sales || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
              </span>
            </div>
            <Calendar className="w-5 h-5 text-indigo-600" />
          </div>
        </div>

        {/* Sales Trend Chart */}
        {chartData.length > 0 && (
          <div className="pt-2">
            <span className="text-xs font-semibold text-slate-700 block mb-3">Sales Trend Trajectory</span>
            <div className="h-60 w-full">
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={chartData} margin={{ top: 5, right: 10, left: 0, bottom: 0 }}>
                  <defs>
                    <linearGradient id="dashSalesGrad" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#10b981" stopOpacity={0.4} />
                      <stop offset="95%" stopColor="#10b981" stopOpacity={0} />
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#f1f5f9" />
                  <XAxis dataKey="date" tick={{ fontSize: 10 }} tickLine={false} stroke="#94a3b8" />
                  <YAxis
                    tick={{ fontSize: 10 }}
                    tickLine={false}
                    stroke="#94a3b8"
                    tickFormatter={(v) => `₹${v >= 1000 ? `${(v / 1000).toFixed(0)}k` : v}`}
                  />
                  <Tooltip
                    formatter={(val: any) => [`₹${Number(val).toLocaleString('en-IN')}`, 'Sales']}
                    contentStyle={{ backgroundColor: '#ffffff', borderRadius: '8px', border: '1px solid #e2e8f0', fontSize: '11px' }}
                  />
                  <Area type="monotone" dataKey="sales" stroke="#10b981" fillOpacity={1} fill="url(#dashSalesGrad)" strokeWidth={2} />
                </AreaChart>
              </ResponsiveContainer>
            </div>
          </div>
        )}
      </div>

      {/* Inventory KPI Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
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
            <p className="text-xs text-slate-500 mt-1">Low Stock / Out of Stock alerts</p>
          </CardContent>
        </Card>
      </div>

      {/* Product & Variant Footwear Architecture */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle className="text-base flex items-center gap-2">
              <Package className="w-5 h-5 text-blue-600" />
              Product & Variant Architecture
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
