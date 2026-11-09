import React, { useState } from 'react'
import { Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { apiService } from '@/services/api'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/features/auth/AuthContext'
import {
  Activity,
  Boxes,
  TrendingUp,
  Receipt,
  AlertCircle,
  Package,
  FileSpreadsheet,
  ArrowUpRight,
  ArrowDownRight,
  DollarSign,
  Calendar,
  Layers,
  ShoppingBag,
  CreditCard,
  RefreshCw,
  ArrowRight,
  Percent,
  Clock,
  Sparkles,
  Store,
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
  const { activeTenant, user } = useAuth()
  const [period, setPeriod] = useState<string>('today')
  const [isExporting, setIsExporting] = useState<boolean>(false)

  // 1. System Health
  const { data: health, refetch: refetchHealth } = useQuery({
    queryKey: ['system-health'],
    queryFn: apiService.getHealth,
    refetchInterval: 30000,
  })

  // 2. Database Health
  const { data: dbHealth, refetch: refetchDb } = useQuery({
    queryKey: ['database-health'],
    queryFn: apiService.getDbHealth,
    refetchInterval: 30000,
  })

  // 3. Inventory Summary
  const { data: inventoryResp, refetch: refetchInventory } = useQuery({
    queryKey: ['inventory-summary-dashboard'],
    queryFn: apiService.getInventorySummary,
  })

  // 4. Sales Performance
  const { data: perfResp, isLoading: isPerfLoading, refetch: refetchPerf } = useQuery({
    queryKey: ['sales-performance-dashboard', period],
    queryFn: () => apiService.getSalesPerformance({ range: period }),
  })

  // 5. Sales Trends (Charts)
  const { data: trendsResp, refetch: refetchTrends } = useQuery({
    queryKey: ['sales-trends-dashboard', period],
    queryFn: () => apiService.getSalesTrends({ range: period === 'today' || period === 'yesterday' ? '7d' : period }),
  })

  // 6. Recent Sales Transactions (5 latest orders)
  const { data: recentSalesResp, refetch: refetchSales } = useQuery({
    queryKey: ['recent-sales-dashboard'],
    queryFn: () => apiService.getSales({ page: 1, page_size: 5 }),
  })

  const invSummary = inventoryResp?.data
  const perf = perfResp?.data
  const trends = trendsResp?.data || []
  const recentSales = recentSalesResp?.data?.items || []

  const chartData = trends.map((t) => ({
    date: t.period,
    sales: Number(t.sales || 0),
    netSales: Number(t.net_sales || 0),
  }))

  const salesGrowth = Number(perf?.sales_growth_pct || 0)

  const handleRefreshAll = () => {
    refetchHealth()
    refetchDb()
    refetchInventory()
    refetchPerf()
    refetchTrends()
    refetchSales()
  }

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

  // Current formatted date
  const todayFormatted = new Intl.DateTimeFormat('en-IN', {
    weekday: 'long',
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  }).format(new Date())

  return (
    <div className="space-y-6 pb-12">
      {/* 1. Executive Top Banner with Store Greeting & Quick Actions */}
      <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-slate-900 via-slate-800 to-blue-950 p-6 sm:p-8 text-white shadow-md border border-slate-800">
        <div className="absolute right-0 top-0 -mt-6 -mr-6 w-96 h-96 rounded-full bg-blue-500/10 blur-3xl pointer-events-none" />
        <div className="absolute right-32 bottom-0 w-64 h-64 rounded-full bg-emerald-500/10 blur-2xl pointer-events-none" />

        <div className="relative z-10 flex flex-col lg:flex-row lg:items-center justify-between gap-6">
          <div className="space-y-2">
            <div className="flex flex-wrap items-center gap-2.5">
              <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full bg-blue-500/20 text-blue-300 text-xs font-semibold border border-blue-400/30">
                <Store className="w-3.5 h-3.5" />
                {activeTenant?.tenant_name || 'KadaiIQ Retail'}
              </span>
              <span className="text-slate-400 text-xs flex items-center gap-1">
                <Clock className="w-3.5 h-3.5" />
                {todayFormatted}
              </span>
            </div>

            <h1 className="text-2xl sm:text-3xl font-extrabold tracking-tight text-white leading-tight">
              Good day, {user?.name?.split(' ')[0] || 'Store Owner'}
            </h1>
            <p className="text-slate-300 text-xs sm:text-sm max-w-2xl leading-relaxed">
              Real-time multi-tenant retail operations, footwear inventory ledger, and revenue performance analytics.
            </p>
          </div>

          <div className="flex flex-wrap items-center gap-3">
            <Button
              variant="outline"
              size="sm"
              onClick={handleRefreshAll}
              className="bg-white/10 hover:bg-white/20 text-white border-white/20 shadow-xs flex items-center gap-1.5 backdrop-blur-xs cursor-pointer"
            >
              <RefreshCw className="w-3.5 h-3.5" />
              <span className="hidden sm:inline">Refresh</span>
            </Button>

            <Button
              variant="default"
              size="sm"
              onClick={handleExportExcel}
              disabled={isExporting}
              className="bg-emerald-600 hover:bg-emerald-700 text-white shadow-sm flex items-center gap-1.5 font-semibold cursor-pointer"
            >
              <FileSpreadsheet className="w-4 h-4 text-emerald-100" />
              {isExporting ? 'Exporting...' : 'Export Excel (.xlsx)'}
            </Button>

            {/* Health Pill Mini */}
            <div className="hidden sm:flex items-center gap-2 bg-slate-800/80 border border-slate-700/80 px-3 py-1.5 rounded-xl text-xs backdrop-blur-xs">
              <Activity className="w-3.5 h-3.5 text-blue-400 animate-pulse" />
              <span className="text-slate-300 font-mono text-[11px]">
                {health?.success ? `v${health.data?.version}` : 'Connecting'}
              </span>
              <span className="text-slate-600">|</span>
              <span className="text-emerald-400 font-mono text-[11px]">
                {dbHealth?.success ? `${dbHealth.data?.latency_ms}ms` : 'DB'}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* 2. Sales Horizon Filter Bar */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 bg-white p-2.5 sm:p-3 rounded-xl border border-slate-200/90 shadow-xs">
        <div className="flex items-center gap-2 px-2">
          <TrendingUp className="w-4 h-4 text-emerald-600" />
          <span className="text-xs font-bold text-slate-800 tracking-tight uppercase">
            Performance Horizon:
          </span>
        </div>

        <div className="flex flex-wrap items-center gap-1 bg-slate-100/90 p-1 rounded-lg">
          {[
            { id: 'today', label: 'Today' },
            { id: 'yesterday', label: 'Yesterday' },
            { id: '7d', label: 'Last 7 Days' },
            { id: '30d', label: 'Last 30 Days' },
            { id: 'this_month', label: 'This Month' },
            { id: 'this_year', label: 'This Year' },
          ].map((btn) => (
            <button
              key={btn.id}
              onClick={() => setPeriod(btn.id)}
              className={`px-3 py-1.5 rounded-md text-xs font-medium transition-all cursor-pointer ${
                period === btn.id
                  ? 'bg-white text-blue-700 shadow-xs font-bold ring-1 ring-slate-200/60'
                  : 'text-slate-600 hover:text-slate-900 hover:bg-slate-200/60'
              }`}
            >
              {btn.label}
            </button>
          ))}
        </div>
      </div>

      {/* 3. Core Retail Performance Grid (4 Hero KPI Cards) */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* KPI 1: Total Sales Amount */}
        <Card className="card-hover border-slate-200/90 relative overflow-hidden">
          <div className="absolute top-0 right-0 w-24 h-24 bg-blue-500/5 rounded-full -mr-6 -mt-6 pointer-events-none" />
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">
              Total Sales Amount
            </span>
            <div className="w-9 h-9 rounded-xl bg-blue-50 text-blue-600 flex items-center justify-center shadow-xs">
              <TrendingUp className="w-4 h-4" />
            </div>
          </CardHeader>
          <CardContent className="space-y-2">
            <div className="text-2xl sm:text-3xl font-extrabold text-slate-900 font-numeric tracking-tight">
              {isPerfLoading ? (
                <div className="h-8 bg-slate-100 animate-pulse rounded-md w-3/4" />
              ) : (
                `₹${Number(perf?.gross_sales || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}`
              )}
            </div>

            <div className="flex items-center justify-between text-xs pt-1 border-t border-slate-100">
              <div className="flex items-center gap-1 font-semibold">
                {salesGrowth >= 0 ? (
                  <span className="text-emerald-600 flex items-center">
                    <ArrowUpRight className="w-3.5 h-3.5" />
                    +{salesGrowth}%
                  </span>
                ) : (
                  <span className="text-rose-600 flex items-center">
                    <ArrowDownRight className="w-3.5 h-3.5" />
                    {salesGrowth}%
                  </span>
                )}
                <span className="text-slate-400 font-normal text-[11px]">vs prior</span>
              </div>
              <span className="text-[11px] text-slate-500 font-medium">
                Net: ₹{Number(perf?.net_sales || 0).toLocaleString('en-IN', { maximumFractionDigits: 0 })}
              </span>
            </div>
          </CardContent>
        </Card>

        {/* KPI 2: Products Sold */}
        <Card className="card-hover border-slate-200/90 relative overflow-hidden">
          <div className="absolute top-0 right-0 w-24 h-24 bg-indigo-500/5 rounded-full -mr-6 -mt-6 pointer-events-none" />
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">
              Products Sold
            </span>
            <div className="w-9 h-9 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center shadow-xs">
              <ShoppingBag className="w-4 h-4" />
            </div>
          </CardHeader>
          <CardContent className="space-y-2">
            <div className="text-2xl sm:text-3xl font-extrabold text-slate-900 font-numeric tracking-tight">
              {isPerfLoading ? (
                <div className="h-8 bg-slate-100 animate-pulse rounded-md w-1/2" />
              ) : (
                `${Number(perf?.total_units_sold || 0).toLocaleString('en-IN')} units`
              )}
            </div>

            <div className="flex items-center justify-between text-xs pt-1 border-t border-slate-100">
              <span className="text-slate-500 text-[11px]">Variant units moved</span>
              <span className="text-indigo-600 font-semibold text-[11px]">Retail POS</span>
            </div>
          </CardContent>
        </Card>

        {/* KPI 3: Completed Bills */}
        <Card className="card-hover border-slate-200/90 relative overflow-hidden">
          <div className="absolute top-0 right-0 w-24 h-24 bg-emerald-500/5 rounded-full -mr-6 -mt-6 pointer-events-none" />
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">
              Completed Bills
            </span>
            <div className="w-9 h-9 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center shadow-xs">
              <Receipt className="w-4 h-4" />
            </div>
          </CardHeader>
          <CardContent className="space-y-2">
            <div className="text-2xl sm:text-3xl font-extrabold text-slate-900 font-numeric tracking-tight">
              {isPerfLoading ? (
                <div className="h-8 bg-slate-100 animate-pulse rounded-md w-1/2" />
              ) : (
                `${Number(perf?.total_bills || 0).toLocaleString('en-IN')} bills`
              )}
            </div>

            <div className="flex items-center justify-between text-xs pt-1 border-t border-slate-100">
              <span className="text-slate-500 text-[11px]">Invoices settled</span>
              <span className="text-emerald-600 font-semibold text-[11px]">100% Tax Finalized</span>
            </div>
          </CardContent>
        </Card>

        {/* KPI 4: Average Bill Value */}
        <Card className="card-hover border-slate-200/90 relative overflow-hidden">
          <div className="absolute top-0 right-0 w-24 h-24 bg-amber-500/5 rounded-full -mr-6 -mt-6 pointer-events-none" />
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <span className="text-xs font-semibold text-slate-500 uppercase tracking-wider">
              Average Bill Value
            </span>
            <div className="w-9 h-9 rounded-xl bg-amber-50 text-amber-600 flex items-center justify-center shadow-xs">
              <CreditCard className="w-4 h-4" />
            </div>
          </CardHeader>
          <CardContent className="space-y-2">
            <div className="text-2xl sm:text-3xl font-extrabold text-slate-900 font-numeric tracking-tight">
              {isPerfLoading ? (
                <div className="h-8 bg-slate-100 animate-pulse rounded-md w-3/4" />
              ) : (
                `₹${Number(perf?.average_bill_value || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}`
              )}
            </div>

            <div className="flex items-center justify-between text-xs pt-1 border-t border-slate-100">
              <span className="text-slate-500 text-[11px]">Ticket revenue</span>
              <span className="text-amber-700 font-semibold text-[11px]">Per customer</span>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* 4. Financial Health & Margins Ribbon */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        {/* Net Realized Sales */}
        <div className="flex items-center justify-between p-4 rounded-xl bg-emerald-50/50 border border-emerald-200/70 shadow-xs">
          <div>
            <span className="text-[11px] font-bold text-emerald-800 uppercase tracking-wider block">
              Net Realized Sales
            </span>
            <div className="text-xl font-extrabold text-emerald-950 font-numeric mt-0.5">
              ₹{Number(perf?.net_sales || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
            </div>
            <p className="text-[11px] text-emerald-700 mt-0.5">After deductions & verified returns</p>
          </div>
          <div className="w-10 h-10 rounded-xl bg-emerald-100 text-emerald-700 flex items-center justify-center">
            <DollarSign className="w-5 h-5" />
          </div>
        </div>

        {/* Gross Profit & Margin */}
        <div className="flex items-center justify-between p-4 rounded-xl bg-blue-50/50 border border-blue-200/70 shadow-xs">
          <div>
            <div className="flex items-center gap-1.5">
              <span className="text-[11px] font-bold text-blue-800 uppercase tracking-wider">
                Gross Profit
              </span>
              <span className="text-[10px] font-bold px-1.5 py-0.2 rounded-full bg-blue-600 text-white">
                {Number(perf?.profit_margin_pct || 0).toFixed(1)}% Margin
              </span>
            </div>
            <div className="text-xl font-extrabold text-blue-950 font-numeric mt-0.5">
              ₹{Number(perf?.gross_profit || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
            </div>
            <p className="text-[11px] text-blue-700 mt-0.5">
              COGS: ₹{Number(perf?.cogs || 0).toLocaleString('en-IN', { maximumFractionDigits: 0 })}
            </p>
          </div>
          <div className="w-10 h-10 rounded-xl bg-blue-100 text-blue-700 flex items-center justify-center">
            <Percent className="w-5 h-5" />
          </div>
        </div>

        {/* Average Daily Velocity */}
        <div className="flex items-center justify-between p-4 rounded-xl bg-indigo-50/50 border border-indigo-200/70 shadow-xs">
          <div>
            <span className="text-[11px] font-bold text-indigo-800 uppercase tracking-wider block">
              Daily Sales Velocity
            </span>
            <div className="text-xl font-extrabold text-indigo-950 font-numeric mt-0.5">
              ₹{Number(perf?.average_daily_sales || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
            </div>
            <p className="text-[11px] text-indigo-700 mt-0.5">Prorated daily revenue trajectory</p>
          </div>
          <div className="w-10 h-10 rounded-xl bg-indigo-100 text-indigo-700 flex items-center justify-center">
            <Calendar className="w-5 h-5" />
          </div>
        </div>
      </div>

      {/* 5. Sales Performance Trend Chart */}
      <Card className="border-slate-200/90 shadow-xs">
        <CardHeader className="flex flex-col sm:flex-row sm:items-center justify-between pb-4 gap-2">
          <div>
            <CardTitle className="text-base font-bold text-slate-900 flex items-center gap-2">
              <TrendingUp className="w-4 h-4 text-emerald-600" />
              Sales Revenue Trajectory
            </CardTitle>
            <CardDescription className="text-xs text-slate-500">
              Verified daily sales and net revenue trend for the selected horizon
            </CardDescription>
          </div>
          <div className="flex items-center gap-3 text-xs">
            <div className="flex items-center gap-1.5">
              <span className="w-3 h-3 rounded-full bg-blue-600" />
              <span className="text-slate-600 font-medium">Gross Sales</span>
            </div>
            <div className="flex items-center gap-1.5">
              <span className="w-3 h-3 rounded-full bg-emerald-500" />
              <span className="text-slate-600 font-medium">Net Sales</span>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {chartData.length > 0 ? (
            <div className="h-72 w-full pt-2">
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={chartData} margin={{ top: 10, right: 10, left: -10, bottom: 0 }}>
                  <defs>
                    <linearGradient id="salesGrad" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#2563eb" stopOpacity={0.3} />
                      <stop offset="95%" stopColor="#2563eb" stopOpacity={0.0} />
                    </linearGradient>
                    <linearGradient id="netSalesGrad" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#10b981" stopOpacity={0.3} />
                      <stop offset="95%" stopColor="#10b981" stopOpacity={0.0} />
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#f1f5f9" />
                  <XAxis
                    dataKey="date"
                    tick={{ fontSize: 11, fill: '#64748b' }}
                    tickLine={false}
                    stroke="#cbd5e1"
                  />
                  <YAxis
                    tick={{ fontSize: 11, fill: '#64748b' }}
                    tickLine={false}
                    stroke="#cbd5e1"
                    tickFormatter={(v) => `₹${v >= 1000 ? `${(v / 1000).toFixed(0)}k` : v}`}
                  />
                  <Tooltip
                    formatter={(val: unknown, name: unknown) => [
                      `₹${Number(val).toLocaleString('en-IN', { minimumFractionDigits: 2 })}`,
                      name === 'sales' ? 'Gross Sales' : 'Net Sales',
                    ]}
                    contentStyle={{
                      backgroundColor: '#ffffff',
                      borderRadius: '10px',
                      border: '1px solid #e2e8f0',
                      boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.08)',
                      fontSize: '12px',
                    }}
                  />
                  <Area
                    type="monotone"
                    dataKey="sales"
                    stroke="#2563eb"
                    strokeWidth={2.5}
                    fillOpacity={1}
                    fill="url(#salesGrad)"
                  />
                  <Area
                    type="monotone"
                    dataKey="netSales"
                    stroke="#10b981"
                    strokeWidth={2}
                    strokeDasharray="4 4"
                    fillOpacity={1}
                    fill="url(#netSalesGrad)"
                  />
                </AreaChart>
              </ResponsiveContainer>
            </div>
          ) : (
            <div className="h-48 flex flex-col items-center justify-center text-center p-6 bg-slate-50/50 rounded-xl border border-dashed border-slate-200">
              <Sparkles className="w-8 h-8 text-slate-300 mb-2" />
              <p className="text-sm font-semibold text-slate-700">No Sales Transactions for This Period</p>
              <p className="text-xs text-slate-400 mt-1 max-w-sm">
                Import sales from ZenBiller POS or create new retail bills to generate revenue trends.
              </p>
              <Link to="/sales">
                <Button size="sm" variant="outline" className="mt-3 cursor-pointer">
                  Import Sales Data
                </Button>
              </Link>
            </div>
          )}
        </CardContent>
      </Card>

      {/* 6. Inventory & Stock Health Overview */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card className="border-slate-200/90 shadow-xs">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-xs font-semibold text-slate-500 uppercase tracking-wider">
              Active Footwear SKUs
            </CardTitle>
            <div className="w-8 h-8 rounded-lg bg-blue-50 text-blue-600 flex items-center justify-center">
              <Boxes className="w-4 h-4" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900 font-numeric">
              {Number(invSummary?.total_variants || 0).toLocaleString('en-IN')}
            </div>
            <p className="text-xs text-slate-500 mt-1">
              Across {invSummary?.total_products || 0} base product models
            </p>
          </CardContent>
        </Card>

        <Card className="border-slate-200/90 shadow-xs">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-xs font-semibold text-slate-500 uppercase tracking-wider">
              Stock Units on Hand
            </CardTitle>
            <div className="w-8 h-8 rounded-lg bg-indigo-50 text-indigo-600 flex items-center justify-center">
              <Layers className="w-4 h-4" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900 font-numeric">
              {Number(invSummary?.total_stock_units || 0).toLocaleString('en-IN')} Units
            </div>
            <p className="text-xs text-slate-500 mt-1">
              Valuation: ₹{Number(invSummary?.stock_valuation || 0).toLocaleString('en-IN')}
            </p>
          </CardContent>
        </Card>

        <Card className="border-slate-200/90 shadow-xs">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <CardTitle className="text-xs font-semibold text-slate-500 uppercase tracking-wider">
              Stock Alert Watch
            </CardTitle>
            <div className="w-8 h-8 rounded-lg bg-rose-50 text-rose-600 flex items-center justify-center">
              <AlertCircle className="w-4 h-4" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900 font-numeric flex items-center gap-2">
              <span className="text-amber-600">{invSummary?.low_stock_count || 0} Low</span>
              <span className="text-slate-300">/</span>
              <span className="text-rose-600">{invSummary?.out_of_stock_count || 0} Out</span>
            </div>
            <div className="flex items-center justify-between mt-1 text-xs">
              <span className="text-slate-500 text-[11px]">Reorder threshold alert</span>
              <Link to="/inventory" className="text-blue-600 hover:underline font-semibold text-[11px]">
                Review stock →
              </Link>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* 7. Recent POS Transactions Table & Enterprise Architecture */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left (2 Cols): Recent Retail Transactions */}
        <Card className="lg:col-span-2 border-slate-200/90 shadow-xs">
          <CardHeader className="flex flex-row items-center justify-between pb-3">
            <div>
              <CardTitle className="text-sm sm:text-base font-bold text-slate-900 flex items-center gap-2">
                <Receipt className="w-4 h-4 text-blue-600" />
                Recent POS Transactions
              </CardTitle>
              <CardDescription className="text-xs text-slate-500">
                Latest customer bills settled and recorded in the tenant ledger
              </CardDescription>
            </div>
            <Link to="/sales">
              <Button variant="ghost" size="sm" className="text-xs text-blue-600 hover:text-blue-700 flex items-center gap-1 cursor-pointer">
                View all records
                <ArrowRight className="w-3.5 h-3.5" />
              </Button>
            </Link>
          </CardHeader>
          <CardContent>
            {recentSales.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="w-full text-left border-collapse text-xs">
                  <thead>
                    <tr className="border-b border-slate-100 text-slate-400 font-semibold uppercase tracking-wider text-[10px]">
                      <th className="py-2.5 px-3">Bill Number</th>
                      <th className="py-2.5 px-3">Items</th>
                      <th className="py-2.5 px-3">Payment</th>
                      <th className="py-2.5 px-3 text-right">Amount</th>
                      <th className="py-2.5 px-3 text-right">Status</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100">
                    {recentSales.map((sale) => (
                      <tr key={sale.id} className="hover:bg-slate-50/70 transition-colors">
                        <td className="py-3 px-3 font-semibold text-slate-900 font-mono">
                          {sale.bill_number}
                        </td>
                        <td className="py-3 px-3 text-slate-600">
                          {sale.item_count || 1} units
                        </td>
                        <td className="py-3 px-3">
                          <span className="px-2 py-0.5 rounded bg-slate-100 text-slate-700 font-medium text-[11px]">
                            {sale.payment_method || 'CASH'}
                          </span>
                        </td>
                        <td className="py-3 px-3 text-right font-bold text-slate-900 font-numeric">
                          ₹{Number(sale.grand_total || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                        </td>
                        <td className="py-3 px-3 text-right">
                          <Badge
                            variant={
                              sale.status === 'COMPLETED'
                                ? 'success'
                                : sale.status === 'CANCELLED'
                                ? 'destructive'
                                : 'outline'
                            }
                            className="text-[10px] uppercase font-bold"
                          >
                            {sale.status}
                          </Badge>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ) : (
              <div className="py-8 text-center text-slate-400 text-xs">
                No POS transactions recorded yet.
              </div>
            )}
          </CardContent>
        </Card>

        {/* Right (1 Col): Footwear Retail Architecture & Audit Guarantees */}
        <Card className="border-slate-200/90 shadow-xs flex flex-col justify-between">
          <CardHeader className="pb-3">
            <CardTitle className="text-sm sm:text-base font-bold text-slate-900 flex items-center gap-2">
              <Package className="w-4 h-4 text-indigo-600" />
              KadaiIQ Integrity Guarantees
            </CardTitle>
            <CardDescription className="text-xs text-slate-500">
              Core multi-tenant architectural standards
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3.5 text-xs">
            <div className="flex items-center justify-between p-2.5 rounded-lg bg-slate-50 border border-slate-100">
              <div>
                <span className="font-semibold text-slate-800 block">Footwear Size Matrix</span>
                <span className="text-[11px] text-slate-500">Sizes, Colours, Barcodes & SKUs</span>
              </div>
              <Badge variant="outline" className="text-[10px] font-semibold text-indigo-700 border-indigo-200 bg-indigo-50">
                Enforced
              </Badge>
            </div>

            <div className="flex items-center justify-between p-2.5 rounded-lg bg-slate-50 border border-slate-100">
              <div>
                <span className="font-semibold text-slate-800 block">Immutable Stock Ledger</span>
                <span className="text-[11px] text-slate-500">Zero silent inventory adjustments</span>
              </div>
              <Badge variant="success" className="text-[10px] font-semibold">
                Guaranteed
              </Badge>
            </div>

            <div className="flex items-center justify-between p-2.5 rounded-lg bg-slate-50 border border-slate-100">
              <div>
                <span className="font-semibold text-slate-800 block">Multi-Tenant Isolation</span>
                <span className="text-[11px] text-slate-500">Server-enforced repository scopes</span>
              </div>
              <Badge variant="outline" className="text-[10px] font-semibold text-emerald-700 border-emerald-200 bg-emerald-50">
                Isolated
              </Badge>
            </div>

            <div className="flex items-center justify-between p-2.5 rounded-lg bg-slate-50 border border-slate-100">
              <div>
                <span className="font-semibold text-slate-800 block">Indian GST Compliance</span>
                <span className="text-[11px] text-slate-500">CGST, SGST, IGST multi-rate audit</span>
              </div>
              <Badge variant="outline" className="text-[10px] font-semibold text-blue-700 border-blue-200 bg-blue-50">
                Ready
              </Badge>
            </div>
          </CardContent>
          <div className="p-4 pt-0">
            <Link to="/analytics" className="w-full block">
              <Button variant="outline" size="sm" className="w-full text-xs flex items-center justify-center gap-1.5 cursor-pointer">
                Deep Sales Analytics
                <ArrowRight className="w-3.5 h-3.5" />
              </Button>
            </Link>
          </div>
        </Card>
      </div>
    </div>
  )
}
