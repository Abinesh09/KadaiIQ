import React, { useState } from 'react'
import { Link } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { apiService, type DetailedProductPerformanceItemData } from '@/services/api'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  TrendingUp,
  Receipt,
  Package,
  FileSpreadsheet,
  ArrowUpRight,
  ArrowDownRight,
  DollarSign,
  ShoppingBag,
  CreditCard,
  RefreshCw,
  ArrowRight,
  Boxes,
  AlertTriangle,
  PieChart,
  Award,
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

  // 1. Sales Performance Summary (Financials, KPIs)
  const { data: perfResp, isLoading: isPerfLoading, refetch: refetchPerf } = useQuery({
    queryKey: ['sales-performance-dashboard', period],
    queryFn: () => apiService.getSalesPerformance({ range: period }),
  })

  // 2. Sales Trends (Charts)
  const { data: trendsResp, isLoading: isTrendsLoading, refetch: refetchTrends } = useQuery({
    queryKey: ['sales-trends-dashboard', period],
    queryFn: () => apiService.getSalesTrends({ range: period === 'today' || period === 'yesterday' ? '7d' : period }),
  })

  // 3. Product Performance (Top Selling & Zero Sales)
  const { data: productPerfResp, refetch: refetchProductPerf } = useQuery({
    queryKey: ['product-performance-dashboard', period],
    queryFn: () => apiService.getDetailedProductPerformance({ range: period }),
  })

  // 4. Category Sales
  const { data: categoryResp, refetch: refetchCategories } = useQuery({
    queryKey: ['category-sales-dashboard', period],
    queryFn: () => apiService.getCategorySales({ range: period }),
  })

  // 5. Payment Methods
  const { data: paymentResp, refetch: refetchPayments } = useQuery({
    queryKey: ['payment-methods-dashboard', period],
    queryFn: () => apiService.getPaymentMethods({ range: period }),
  })

  // 6. Inventory Summary
  const { data: inventoryResp, refetch: refetchInventory } = useQuery({
    queryKey: ['inventory-summary-dashboard'],
    queryFn: apiService.getInventorySummary,
  })

  // 7. Recent Transactions (5 Latest Bills)
  const { data: recentSalesResp, refetch: refetchSales } = useQuery({
    queryKey: ['recent-sales-dashboard'],
    queryFn: () => apiService.getSales({ page: 1, page_size: 5 }),
  })

  const perf = perfResp?.data
  const trends = trendsResp?.data || []
  const invSummary = inventoryResp?.data
  const topProducts = (productPerfResp?.data?.top_selling || []).slice(0, 5)
  const zeroSalesProducts = (productPerfResp?.data?.zero_sales || []).slice(0, 5)
  const categories = categoryResp?.data || []
  const paymentMethods = paymentResp?.data || []
  const recentSales = recentSalesResp?.data?.items || []

  const chartData = trends.map((t) => ({
    date: t.period,
    grossSales: Number(t.sales || 0),
    netSales: Number(t.net_sales || 0),
  }))

  const salesGrowth = Number(perf?.sales_growth_pct || 0)

  const handleRefreshAll = () => {
    refetchPerf()
    refetchTrends()
    refetchProductPerf()
    refetchCategories()
    refetchPayments()
    refetchInventory()
    refetchSales()
  }

  const handleExportExcel = async () => {
    try {
      setIsExporting(true)
      const blob = await apiService.downloadExcelReport({ range: period })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `KadaiIQ_Business_Overview_${period}_${new Date().toISOString().slice(0, 10)}.xlsx`
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
    <div className="space-y-6 pb-12 max-w-[1600px] mx-auto">
      {/* A. Page Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-[#E5EAF0] pb-5">
        <div>
          <h1 className="text-xl sm:text-2xl font-bold tracking-tight text-[#172033]">
            Business Overview
          </h1>
          <p className="text-xs sm:text-sm text-[#64748B] mt-0.5">
            Monitor your shop's sales, revenue and inventory performance.
          </p>
        </div>

        {/* Date Filter & Actions Toolbar */}
        <div className="flex flex-wrap items-center gap-2.5">
          {/* Horizon Selector */}
          <div className="inline-flex items-center bg-[#F5F7FA] border border-[#E5EAF0] p-1 rounded-lg">
            {[
              { id: 'today', label: 'Today' },
              { id: '7d', label: 'Last 7 days' },
              { id: '30d', label: 'Last 30 days' },
              { id: 'this_month', label: 'This month' },
              { id: 'this_year', label: 'This year' },
            ].map((btn) => (
              <button
                key={btn.id}
                onClick={() => setPeriod(btn.id)}
                className={`px-3 py-1.5 rounded-md text-xs font-semibold transition-colors cursor-pointer ${
                  period === btn.id
                    ? 'bg-[#EDE9FE] text-[#5B45B5] shadow-xs border border-[#8B5CF6]/30'
                    : 'text-[#64748B] hover:text-[#5B45B5] hover:bg-[#F5F3FF]'
                }`}
              >
                {btn.label}
              </button>
            ))}
          </div>

          {/* Refresh Action */}
          <Button
            variant="outline"
            size="sm"
            onClick={handleRefreshAll}
            className="flex items-center gap-1.5 text-xs text-[#172033]"
            title="Refresh dashboard data"
          >
            <RefreshCw className="w-3.5 h-3.5 text-[#64748B]" />
            <span className="hidden sm:inline">Refresh</span>
          </Button>

          {/* Export Action */}
          <Button
            variant="accent"
            size="sm"
            onClick={handleExportExcel}
            disabled={isExporting}
            className="flex items-center gap-1.5 text-xs"
          >
            <FileSpreadsheet className="w-3.5 h-3.5" />
            <span>{isExporting ? 'Exporting...' : 'Export Excel'}</span>
          </Button>
        </div>
      </div>

      {/* B. Main KPI Section (4 Primary Cards) */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* KPI 1: Net Sales */}
        <Card className="shadow-card-hover">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <span className="text-xs font-semibold uppercase tracking-wider text-[#64748B]">
              Net Sales
            </span>
            <div className="w-8 h-8 rounded-lg bg-[#16A36A]/10 text-[#16A36A] flex items-center justify-center">
              <DollarSign className="w-4 h-4" />
            </div>
          </CardHeader>
          <CardContent className="space-y-1.5">
            <div className="text-2xl sm:text-3xl font-bold tracking-tight text-[#172033] font-numeric">
              {isPerfLoading ? (
                <div className="h-8 bg-slate-100 animate-pulse rounded w-3/4" />
              ) : (
                `₹${Number(perf?.net_sales || 0).toLocaleString('en-IN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
              )}
            </div>
            <div className="flex items-center justify-between text-xs pt-1 border-t border-[#E5EAF0]">
              <div className="flex items-center gap-1 font-semibold">
                {salesGrowth >= 0 ? (
                  <span className="text-[#16A36A] flex items-center">
                    <ArrowUpRight className="w-3.5 h-3.5" />
                    +{salesGrowth}%
                  </span>
                ) : (
                  <span className="text-[#DC2626] flex items-center">
                    <ArrowDownRight className="w-3.5 h-3.5" />
                    {salesGrowth}%
                  </span>
                )}
                <span className="text-[#64748B] font-normal text-[11px]">vs prior period</span>
              </div>
              <span className="text-[11px] text-[#64748B]">Realized</span>
            </div>
          </CardContent>
        </Card>

        {/* KPI 2: Products Sold */}
        <Card className="shadow-card-hover">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <span className="text-xs font-semibold uppercase tracking-wider text-[#64748B]">
              Products Sold
            </span>
            <div className="w-8 h-8 rounded-lg bg-[#172B4D]/10 text-[#172B4D] flex items-center justify-center">
              <ShoppingBag className="w-4 h-4" />
            </div>
          </CardHeader>
          <CardContent className="space-y-1.5">
            <div className="text-2xl sm:text-3xl font-bold tracking-tight text-[#172033] font-numeric">
              {isPerfLoading ? (
                <div className="h-8 bg-slate-100 animate-pulse rounded w-1/2" />
              ) : (
                `${Number(perf?.total_units_sold || 0).toLocaleString('en-IN')} units`
              )}
            </div>
            <div className="flex items-center justify-between text-xs pt-1 border-t border-[#E5EAF0]">
              <span className="text-[#64748B] text-[11px]">Variant units moved</span>
              <span className="text-[#172B4D] font-semibold text-[11px]">Retail POS</span>
            </div>
          </CardContent>
        </Card>

        {/* KPI 3: Completed Bills */}
        <Card className="shadow-card-hover">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <span className="text-xs font-semibold uppercase tracking-wider text-[#64748B]">
              Completed Bills
            </span>
            <div className="w-8 h-8 rounded-lg bg-[#243B5A]/10 text-[#243B5A] flex items-center justify-center">
              <Receipt className="w-4 h-4" />
            </div>
          </CardHeader>
          <CardContent className="space-y-1.5">
            <div className="text-2xl sm:text-3xl font-bold tracking-tight text-[#172033] font-numeric">
              {isPerfLoading ? (
                <div className="h-8 bg-slate-100 animate-pulse rounded w-1/2" />
              ) : (
                `${Number(perf?.total_bills || 0).toLocaleString('en-IN')} bills`
              )}
            </div>
            <div className="flex items-center justify-between text-xs pt-1 border-t border-[#E5EAF0]">
              <span className="text-[#64748B] text-[11px]">Invoices settled</span>
              <span className="text-[#16A36A] font-semibold text-[11px]">Tax Finalized</span>
            </div>
          </CardContent>
        </Card>

        {/* KPI 4: Average Bill Value */}
        <Card className="shadow-card-hover">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <span className="text-xs font-semibold uppercase tracking-wider text-[#64748B]">
              Average Bill Value
            </span>
            <div className="w-8 h-8 rounded-lg bg-[#D97706]/10 text-[#D97706] flex items-center justify-center">
              <CreditCard className="w-4 h-4" />
            </div>
          </CardHeader>
          <CardContent className="space-y-1.5">
            <div className="text-2xl sm:text-3xl font-bold tracking-tight text-[#172033] font-numeric">
              {isPerfLoading ? (
                <div className="h-8 bg-slate-100 animate-pulse rounded w-3/4" />
              ) : (
                `₹${Number(perf?.average_bill_value || 0).toLocaleString('en-IN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
              )}
            </div>
            <div className="flex items-center justify-between text-xs pt-1 border-t border-[#E5EAF0]">
              <span className="text-[#64748B] text-[11px]">Ticket revenue</span>
              <span className="text-[#D97706] font-semibold text-[11px]">Per customer</span>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* C. Financial Summary (Secondary Breakdown) */}
      <div className="bg-white p-4 sm:p-5 rounded-xl border border-[#E5EAF0] shadow-card">
        <div className="flex items-center justify-between mb-3 border-b border-[#E5EAF0] pb-2.5">
          <div className="flex items-center gap-2">
            <span className="text-xs font-bold text-[#172033] uppercase tracking-wider">
              Financial Summary
            </span>
            <span className="text-[11px] text-[#64748B] hidden sm:inline">
              (Accounting distinctions & verified ledger reconciliation)
            </span>
          </div>
          <span className="text-[11px] font-medium text-[#64748B]">Horizon: {period}</span>
        </div>

        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3">
          {/* Gross Sales */}
          <div className="bg-[#F5F7FA] p-3 rounded-lg border border-[#E5EAF0]">
            <span className="text-[10px] font-semibold uppercase tracking-wider text-[#64748B] block">
              Gross Sales
            </span>
            <span className="text-base sm:text-lg font-bold text-[#172033] font-numeric block mt-0.5">
              ₹{Number(perf?.gross_sales || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
            </span>
            <span className="text-[10px] text-[#64748B]">Before deductions</span>
          </div>

          {/* Discounts */}
          <div className="bg-[#F5F7FA] p-3 rounded-lg border border-[#E5EAF0]">
            <span className="text-[10px] font-semibold uppercase tracking-wider text-[#64748B] block">
              Discounts
            </span>
            <span className="text-base sm:text-lg font-bold text-[#D97706] font-numeric block mt-0.5">
              ₹{Number(perf?.total_discounts || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
            </span>
            <span className="text-[10px] text-[#64748B]">Applied to bills</span>
          </div>

          {/* Returns */}
          <div className="bg-[#F5F7FA] p-3 rounded-lg border border-[#E5EAF0]">
            <span className="text-[10px] font-semibold uppercase tracking-wider text-[#64748B] block">
              Returns
            </span>
            <span className="text-base sm:text-lg font-bold text-[#DC2626] font-numeric block mt-0.5">
              ₹{Number(perf?.total_returns || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
            </span>
            <span className="text-[10px] text-[#64748B]">Returned value</span>
          </div>

          {/* Net Sales */}
          <div className="bg-[#16A36A]/5 p-3 rounded-lg border border-[#16A36A]/20">
            <span className="text-[10px] font-semibold uppercase tracking-wider text-[#16A36A] block">
              Net Sales
            </span>
            <span className="text-base sm:text-lg font-bold text-[#16A36A] font-numeric block mt-0.5">
              ₹{Number(perf?.net_sales || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
            </span>
            <span className="text-[10px] text-[#16A36A]">Gross - Disc - Ret</span>
          </div>

          {/* Amount Collected */}
          <div className="bg-[#F5F7FA] p-3 rounded-lg border border-[#E5EAF0]">
            <span className="text-[10px] font-semibold uppercase tracking-wider text-[#64748B] block">
              Amount Collected
            </span>
            <span className="text-base sm:text-lg font-bold text-[#172033] font-numeric block mt-0.5">
              ₹{Number(perf?.net_sales || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
            </span>
            <span className="text-[10px] text-[#64748B]">Cash, UPI & Card</span>
          </div>

          {/* Gross Profit */}
          <div className="bg-[#172B4D]/5 p-3 rounded-lg border border-[#172B4D]/20">
            <div className="flex items-center justify-between">
              <span className="text-[10px] font-semibold uppercase tracking-wider text-[#172B4D] block">
                Gross Profit
              </span>
              <span className="text-[9px] font-bold px-1.5 py-0.2 bg-[#172B4D] text-white rounded">
                {Number(perf?.profit_margin_pct || 0).toFixed(1)}%
              </span>
            </div>
            <span className="text-base sm:text-lg font-bold text-[#172B4D] font-numeric block mt-0.5">
              ₹{Number(perf?.gross_profit || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
            </span>
            <span className="text-[10px] text-[#64748B]">Based on purchase cost</span>
          </div>
        </div>
      </div>

      {/* D. Main Sales Performance Chart */}
      <Card className="shadow-card">
        <CardHeader className="flex flex-col sm:flex-row sm:items-center justify-between pb-3 gap-2">
          <div>
            <CardTitle className="text-sm sm:text-base font-bold text-[#172033] flex items-center gap-2">
              <TrendingUp className="w-4 h-4 text-[#16A36A]" />
              Revenue & Sales Trend
            </CardTitle>
            <CardDescription className="text-xs text-[#64748B]">
              Daily sales velocity and net revenue progression for the selected timeframe.
            </CardDescription>
          </div>

          <div className="flex items-center gap-4 text-xs">
            <div className="flex items-center gap-1.5">
              <span className="w-3 h-3 rounded-sm bg-[#172B4D]" />
              <span className="text-[#64748B] font-medium">Gross Sales</span>
            </div>
            <div className="flex items-center gap-1.5">
              <span className="w-3 h-3 rounded-sm bg-[#16A36A]" />
              <span className="text-[#64748B] font-medium">Net Sales</span>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {isTrendsLoading ? (
            <div className="h-72 flex items-center justify-center text-xs text-[#64748B]">
              Loading sales trend telemetry...
            </div>
          ) : chartData.length > 0 ? (
            <div className="h-72 sm:h-80 w-full pt-2">
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={chartData} margin={{ top: 10, right: 10, left: -10, bottom: 0 }}>
                  <defs>
                    <linearGradient id="chartGrossSalesGrad" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#172B4D" stopOpacity={0.15} />
                      <stop offset="95%" stopColor="#172B4D" stopOpacity={0.0} />
                    </linearGradient>
                    <linearGradient id="chartNetSalesGrad" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#16A36A" stopOpacity={0.2} />
                      <stop offset="95%" stopColor="#16A36A" stopOpacity={0.0} />
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#E5EAF0" />
                  <XAxis
                    dataKey="date"
                    tick={{ fontSize: 11, fill: '#64748B' }}
                    tickLine={false}
                    stroke="#E5EAF0"
                  />
                  <YAxis
                    tick={{ fontSize: 11, fill: '#64748B' }}
                    tickLine={false}
                    stroke="#E5EAF0"
                    tickFormatter={(v) => `₹${v >= 1000 ? `${(v / 1000).toFixed(0)}k` : v}`}
                  />
                  <Tooltip
                    formatter={(val: unknown, name: unknown) => [
                      `₹${Number(val).toLocaleString('en-IN', { minimumFractionDigits: 2 })}`,
                      name === 'grossSales' ? 'Gross Sales' : 'Net Sales',
                    ]}
                    contentStyle={{
                      backgroundColor: '#FFFFFF',
                      borderRadius: '8px',
                      border: '1px solid #E5EAF0',
                      boxShadow: '0 4px 6px -1px rgba(23, 43, 77, 0.08)',
                      fontSize: '12px',
                    }}
                  />
                  <Area
                    type="monotone"
                    dataKey="grossSales"
                    stroke="#172B4D"
                    strokeWidth={2}
                    fillOpacity={1}
                    fill="url(#chartGrossSalesGrad)"
                  />
                  <Area
                    type="monotone"
                    dataKey="netSales"
                    stroke="#16A36A"
                    strokeWidth={2}
                    fillOpacity={1}
                    fill="url(#chartNetSalesGrad)"
                  />
                </AreaChart>
              </ResponsiveContainer>
            </div>
          ) : (
            <div className="h-56 flex flex-col items-center justify-center text-center p-6 bg-[#F5F7FA] rounded-lg border border-dashed border-[#E5EAF0]">
              <Package className="w-8 h-8 text-[#64748B] mb-2 opacity-50" />
              <p className="text-sm font-semibold text-[#172033]">No Sales Data for This Horizon</p>
              <p className="text-xs text-[#64748B] mt-1 max-w-sm">
                No verified POS transactions were recorded for {period}. Import ZenBiller sales to populate revenue graphs.
              </p>
              <Link to="/sales">
                <Button size="sm" variant="outline" className="mt-3">
                  Import Sales Data
                </Button>
              </Link>
            </div>
          )}
        </CardContent>
      </Card>

      {/* E. Secondary Analytics Sections (Balanced 2-Column Responsive Grid) */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-5">
        {/* 1. Top-Selling Products */}
        <Card className="shadow-card">
          <CardHeader className="flex flex-row items-center justify-between pb-3">
            <div>
              <CardTitle className="text-sm font-bold text-[#172033] flex items-center gap-2">
                <Award className="w-4 h-4 text-[#D97706]" />
                Top-Selling Products
              </CardTitle>
              <CardDescription className="text-xs text-[#64748B]">
                Highest volume product variants sold during this period
              </CardDescription>
            </div>
            <Link to="/products">
              <Button variant="ghost" size="sm" className="text-xs text-[#16A36A] flex items-center gap-1">
                View all
                <ArrowRight className="w-3.5 h-3.5" />
              </Button>
            </Link>
          </CardHeader>
          <CardContent>
            {topProducts.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="w-full text-left text-xs border-collapse">
                  <thead>
                    <tr className="border-b border-[#E5EAF0] text-[#64748B] font-semibold text-[10px] uppercase">
                      <th className="py-2 px-2">Product</th>
                      <th className="py-2 px-2 text-right">Units</th>
                      <th className="py-2 px-2 text-right">Net Revenue</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-[#E5EAF0]">
                    {topProducts.map((p: DetailedProductPerformanceItemData, idx: number) => (
                      <tr key={p.product_id || idx} className="hover:bg-[#F5F7FA]">
                        <td className="py-2.5 px-2">
                          <span className="font-semibold text-[#172033] block">{p.product_name}</span>
                          <span className="text-[10px] text-[#64748B] block">
                            {p.brand_name || 'Retail'} • {p.category_name || 'Footwear'} {p.size ? `• Size ${p.size}` : ''}
                          </span>
                        </td>
                        <td className="py-2.5 px-2 text-right font-numeric font-semibold text-[#172033]">
                          {p.units_sold}
                        </td>
                        <td className="py-2.5 px-2 text-right font-numeric font-bold text-[#16A36A]">
                          ₹{Number(p.net_sales || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ) : (
              <div className="py-8 text-center text-xs text-[#64748B]">
                No product sales recorded in this period.
              </div>
            )}
          </CardContent>
        </Card>

        {/* 2. Category Performance */}
        <Card className="shadow-card">
          <CardHeader className="flex flex-row items-center justify-between pb-3">
            <div>
              <CardTitle className="text-sm font-bold text-[#172033] flex items-center gap-2">
                <PieChart className="w-4 h-4 text-[#172B4D]" />
                Category Performance
              </CardTitle>
              <CardDescription className="text-xs text-[#64748B]">
                Sales contribution across footwear retail categories
              </CardDescription>
            </div>
            <Link to="/analytics">
              <Button variant="ghost" size="sm" className="text-xs text-[#16A36A] flex items-center gap-1">
                Category BI
                <ArrowRight className="w-3.5 h-3.5" />
              </Button>
            </Link>
          </CardHeader>
          <CardContent className="space-y-3.5">
            {categories.length > 0 ? (
              categories.map((c, idx) => {
                const pct = Number(c.percentage || 0)
                return (
                  <div key={idx} className="space-y-1">
                    <div className="flex items-center justify-between text-xs">
                      <span className="font-medium text-[#172033]">{c.category_name}</span>
                      <div className="flex items-center gap-2">
                        <span className="text-[#64748B] font-numeric text-[11px]">{c.units_sold} units</span>
                        <span className="font-bold text-[#172033] font-numeric">
                          ₹{Number(c.revenue || 0).toLocaleString('en-IN', { maximumFractionDigits: 0 })}
                        </span>
                        <span className="text-[10px] font-semibold text-[#16A36A]">({pct.toFixed(1)}%)</span>
                      </div>
                    </div>
                    <div className="w-full bg-[#E5EAF0] h-2 rounded-full overflow-hidden">
                      <div
                        className="bg-[#172B4D] h-full rounded-full transition-all duration-300"
                        style={{ width: `${Math.min(pct, 100)}%` }}
                      />
                    </div>
                  </div>
                )
              })
            ) : (
              <div className="py-8 text-center text-xs text-[#64748B]">
                No category sales recorded.
              </div>
            )}
          </CardContent>
        </Card>

        {/* 3. Payment Method Breakdown */}
        <Card className="shadow-card">
          <CardHeader className="flex flex-row items-center justify-between pb-3">
            <div>
              <CardTitle className="text-sm font-bold text-[#172033] flex items-center gap-2">
                <CreditCard className="w-4 h-4 text-[#243B5A]" />
                Payment Method Breakdown
              </CardTitle>
              <CardDescription className="text-xs text-[#64748B]">
                Collection breakdown across cash, UPI and card modes
              </CardDescription>
            </div>
            <Link to="/sales">
              <Button variant="ghost" size="sm" className="text-xs text-[#16A36A] flex items-center gap-1">
                View bills
                <ArrowRight className="w-3.5 h-3.5" />
              </Button>
            </Link>
          </CardHeader>
          <CardContent>
            {paymentMethods.length > 0 ? (
              <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
                {paymentMethods.map((pm, idx) => (
                  <div key={idx} className="bg-[#F5F7FA] p-3 rounded-lg border border-[#E5EAF0]">
                    <span className="text-[10px] font-semibold uppercase text-[#64748B] block truncate">
                      {pm.method || 'CASH'}
                    </span>
                    <span className="text-base font-bold text-[#172033] font-numeric block mt-0.5">
                      ₹{Number(pm.total_sales || 0).toLocaleString('en-IN', { maximumFractionDigits: 0 })}
                    </span>
                    <div className="flex items-center justify-between text-[10px] text-[#64748B] mt-1">
                      <span>{pm.bill_count} bills</span>
                      <span className="font-semibold text-[#16A36A]">{Number(pm.percentage || 0).toFixed(0)}%</span>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="py-8 text-center text-xs text-[#64748B]">
                No payment transactions recorded.
              </div>
            )}
          </CardContent>
        </Card>

        {/* 4. Inventory Health & Stock Summary */}
        <Card className="shadow-card">
          <CardHeader className="flex flex-row items-center justify-between pb-3">
            <div>
              <CardTitle className="text-sm font-bold text-[#172033] flex items-center gap-2">
                <Boxes className="w-4 h-4 text-[#16A36A]" />
                Inventory Health
              </CardTitle>
              <CardDescription className="text-xs text-[#64748B]">
                Footwear catalogue metrics and physical stock valuation
              </CardDescription>
            </div>
            <Link to="/inventory">
              <Button variant="ghost" size="sm" className="text-xs text-[#16A36A] flex items-center gap-1">
                Stock ledger
                <ArrowRight className="w-3.5 h-3.5" />
              </Button>
            </Link>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 gap-3">
              <div className="bg-[#F5F7FA] p-3 rounded-lg border border-[#E5EAF0]">
                <span className="text-[10px] font-semibold uppercase text-[#64748B] block">
                  Active Variants
                </span>
                <span className="text-lg font-bold text-[#172033] font-numeric block mt-0.5">
                  {invSummary?.total_variants || 0} SKUs
                </span>
                <span className="text-[10px] text-[#64748B]">{invSummary?.total_products || 0} product models</span>
              </div>

              <div className="bg-[#F5F7FA] p-3 rounded-lg border border-[#E5EAF0]">
                <span className="text-[10px] font-semibold uppercase text-[#64748B] block">
                  Units on Hand
                </span>
                <span className="text-lg font-bold text-[#172033] font-numeric block mt-0.5">
                  {Number(invSummary?.total_stock_units || 0).toLocaleString('en-IN')} units
                </span>
                <span className="text-[10px] text-[#16A36A] font-semibold">
                  Valuation: ₹{Number(invSummary?.stock_valuation || 0).toLocaleString('en-IN', { maximumFractionDigits: 0 })}
                </span>
              </div>
            </div>

            {/* Low Stock Watch */}
            <div className="mt-3 p-3 rounded-lg bg-amber-50/60 border border-amber-200/80 flex items-center justify-between">
              <div className="flex items-center gap-2">
                <AlertTriangle className="w-4 h-4 text-[#D97706] shrink-0" />
                <div>
                  <span className="text-xs font-bold text-[#D97706] block">
                    {invSummary?.low_stock_count || 0} Low Stock • {invSummary?.out_of_stock_count || 0} Out of Stock
                  </span>
                  <span className="text-[10px] text-[#64748B]">Variants below minimum reorder threshold</span>
                </div>
              </div>
              <Link to="/inventory">
                <Button size="sm" variant="outline" className="text-xs h-7 px-2 border-amber-300 text-[#D97706] bg-white">
                  Restock
                </Button>
              </Link>
            </div>
          </CardContent>
        </Card>

        {/* 5. Products with Zero Sales (Idle Locked Capital Alert) */}
        <Card className="shadow-card">
          <CardHeader className="flex flex-row items-center justify-between pb-3">
            <div>
              <CardTitle className="text-sm font-bold text-[#172033] flex items-center gap-2">
                <AlertTriangle className="w-4 h-4 text-[#D97706]" />
                Zero-Sales Watch
              </CardTitle>
              <CardDescription className="text-xs text-[#64748B]">
                Footwear products with 0 sales during the selected timeframe
              </CardDescription>
            </div>
            <Link to="/analytics">
              <Button variant="ghost" size="sm" className="text-xs text-[#16A36A] flex items-center gap-1">
                Dead stock BI
                <ArrowRight className="w-3.5 h-3.5" />
              </Button>
            </Link>
          </CardHeader>
          <CardContent>
            {zeroSalesProducts.length > 0 ? (
              <div className="space-y-2">
                <div className="p-2.5 rounded-lg bg-amber-50/50 border border-amber-200/60 flex items-center justify-between text-xs mb-2">
                  <span className="text-[#D97706] font-semibold">
                    {productPerfResp?.data?.zero_sales?.length || zeroSalesProducts.length} Idle SKUs
                  </span>
                  <span className="font-bold text-[#172033] font-numeric">
                    Zero Movement
                  </span>
                </div>
                <div className="divide-y divide-[#E5EAF0]">
                  {zeroSalesProducts.map((p: DetailedProductPerformanceItemData, idx: number) => (
                    <div key={p.product_id || idx} className="py-2 flex items-center justify-between text-xs">
                      <div>
                        <span className="font-semibold text-[#172033] block">{p.product_name}</span>
                        <span className="text-[10px] text-[#64748B]">{p.category_name} {p.size ? `• Size ${p.size}` : ''}</span>
                      </div>
                      <Badge variant="warning" className="text-[10px]">
                        0 units
                      </Badge>
                    </div>
                  ))}
                </div>
              </div>
            ) : (
              <div className="py-8 text-center text-xs text-[#16A36A] font-semibold">
                ✓ All active product models generated sales during this period.
              </div>
            )}
          </CardContent>
        </Card>

        {/* 6. Recent POS Transactions Table */}
        <Card className="shadow-card">
          <CardHeader className="flex flex-row items-center justify-between pb-3">
            <div>
              <CardTitle className="text-sm font-bold text-[#172033] flex items-center gap-2">
                <Receipt className="w-4 h-4 text-[#172B4D]" />
                Recent Transactions
              </CardTitle>
              <CardDescription className="text-xs text-[#64748B]">
                Latest customer invoices settled in the tenant ledger
              </CardDescription>
            </div>
            <Link to="/sales">
              <Button variant="ghost" size="sm" className="text-xs text-[#16A36A] flex items-center gap-1">
                View all POS
                <ArrowRight className="w-3.5 h-3.5" />
              </Button>
            </Link>
          </CardHeader>
          <CardContent>
            {recentSales.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="w-full text-left text-xs border-collapse">
                  <thead>
                    <tr className="border-b border-[#E5EAF0] text-[#64748B] font-semibold text-[10px] uppercase">
                      <th className="py-2 px-2">Bill No.</th>
                      <th className="py-2 px-2">Payment</th>
                      <th className="py-2 px-2 text-right">Amount</th>
                      <th className="py-2 px-2 text-right">Status</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-[#E5EAF0]">
                    {recentSales.map((s) => (
                      <tr key={s.id} className="hover:bg-[#F5F7FA]">
                        <td className="py-2.5 px-2 font-mono font-medium text-[#172033]">
                          {s.bill_number}
                        </td>
                        <td className="py-2.5 px-2 text-[#64748B]">
                          <span className="px-1.5 py-0.5 rounded bg-[#F5F7FA] border border-[#E5EAF0] text-[10px]">
                            {s.payment_method || 'CASH'}
                          </span>
                        </td>
                        <td className="py-2.5 px-2 text-right font-numeric font-bold text-[#172033]">
                          ₹{Number(s.grand_total || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                        </td>
                        <td className="py-2.5 px-2 text-right">
                          <Badge
                            variant={
                              s.status === 'COMPLETED'
                                ? 'success'
                                : s.status === 'CANCELLED'
                                ? 'destructive'
                                : 'outline'
                            }
                            className="text-[10px]"
                          >
                            {s.status}
                          </Badge>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ) : (
              <div className="py-8 text-center text-xs text-[#64748B]">
                No POS transactions recorded yet.
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
