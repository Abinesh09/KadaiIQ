import React, { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  apiService,
  type TrendPointData,
  type PaymentMethodData,
  type FootwearIntelligenceData,
  type SalesPerformanceSummaryData,
  type ProductPerformanceResponseData,
  type DetailedProductPerformanceItemData,
  type DailyBreakdownItemData,
  type MonthlyBreakdownItemData,
  type YearlyBreakdownItemData,
} from '@/services/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  TrendingUp,
  Receipt,
  Percent,
  ShoppingBag,
  RefreshCw,
  CreditCard,
  Layers,
  ArrowUpRight,
  ArrowDownRight,
  Calendar,
  FileSpreadsheet,
  AlertTriangle,
  Award,
  CalendarDays,
  DollarSign,
  Search,
} from 'lucide-react'
import {
  ResponsiveContainer,
  AreaChart,
  Area,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
} from 'recharts'

const COLORS = ['#10b981', '#3b82f6', '#8b5cf6', '#f59e0b', '#ef4444', '#06b6d4', '#ec4899']

export const AnalyticsPage: React.FC = () => {
  const [dateRange, setDateRange] = useState<string>('this_month')
  const [dateFrom, setDateFrom] = useState<string>('')
  const [dateTo, setDateTo] = useState<string>('')
  const [isExporting, setIsExporting] = useState<boolean>(false)
  const [activeTab, setActiveTab] = useState<'trends' | 'daily' | 'monthly' | 'yearly' | 'products' | 'footwear' | 'payments'>('trends')
  const [productSubTab, setProductSubTab] = useState<'top' | 'slow' | 'zero'>('top')
  const [productSearch, setProductSearch] = useState<string>('')

  const queryParams = {
    range: dateRange,
    from: dateRange === 'custom' ? dateFrom : undefined,
    to: dateRange === 'custom' ? dateTo : undefined,
  }

  // 1. Sales Performance Summary (Gross, Net, Discounts, Returns, COGS, Profit, Velocity)
  const {
    data: perfResp,
    isLoading: isLoadingPerf,
    refetch: refetchPerf,
  } = useQuery({
    queryKey: ['sales-performance', queryParams],
    queryFn: () => apiService.getSalesPerformance(queryParams),
  })

  // 2. Daily Sales Trends (for charts)
  const { data: trendsResp, isLoading: isLoadingTrends, refetch: refetchTrends } = useQuery({
    queryKey: ['analytics-trends', queryParams],
    queryFn: () => apiService.getSalesTrends(queryParams),
  })

  // 3. Daily Breakdown (for daily view table)
  const { data: dailyResp, isLoading: isLoadingDaily, refetch: refetchDaily } = useQuery({
    queryKey: ['sales-daily-breakdown', queryParams],
    queryFn: () => apiService.getDailyBreakdown(queryParams),
    enabled: activeTab === 'daily',
  })

  // 4. Monthly Breakdown (for monthly view table)
  const { data: monthlyResp, isLoading: isLoadingMonthly, refetch: refetchMonthly } = useQuery({
    queryKey: ['sales-monthly-breakdown', queryParams],
    queryFn: () => apiService.getMonthlyBreakdown(queryParams),
    enabled: activeTab === 'monthly',
  })

  // 4b. Yearly Breakdown (for yearly view table)
  const { data: yearlyResp, isLoading: isLoadingYearly, refetch: refetchYearly } = useQuery({
    queryKey: ['sales-yearly-breakdown', queryParams],
    queryFn: () => apiService.getYearlyBreakdown(queryParams),
    enabled: activeTab === 'yearly',
  })

  // 5. Detailed Product Performance & Zero-Sales Alerts
  const { data: productPerfResp, isLoading: isLoadingProdPerf, refetch: refetchProdPerf } = useQuery({
    queryKey: ['product-detailed-perf', queryParams],
    queryFn: () => apiService.getDetailedProductPerformance(queryParams),
    enabled: activeTab === 'products',
  })

  // 6. Payment Methods Distribution
  const { data: paymentsResp, refetch: refetchPayments } = useQuery({
    queryKey: ['analytics-payments', queryParams],
    queryFn: () => apiService.getPaymentMethods(queryParams),
  })

  // 7. Footwear Intelligence
  const { data: footwearResp, refetch: refetchFootwear } = useQuery({
    queryKey: ['analytics-footwear', queryParams],
    queryFn: () => apiService.getFootwearIntelligence(queryParams),
  })

  const refreshAll = () => {
    refetchPerf()
    refetchTrends()
    refetchDaily()
    refetchMonthly()
    refetchYearly()
    refetchProdPerf()
    refetchPayments()
    refetchFootwear()
  }

  const handleExportExcel = async () => {
    try {
      setIsExporting(true)
      const blob = await apiService.downloadExcelReport(queryParams)
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `KadaiIQ_Sales_Analytics_${dateRange}_${new Date().toISOString().slice(0, 10)}.xlsx`
      document.body.appendChild(a)
      a.click()
      window.URL.revokeObjectURL(url)
      document.body.removeChild(a)
    } catch (err) {
      console.error('Failed to export Excel analytics:', err)
    } finally {
      setIsExporting(false)
    }
  }

  const perf: SalesPerformanceSummaryData | undefined = perfResp?.data
  const trends: TrendPointData[] = trendsResp?.data || []
  const dailyBreakdown: DailyBreakdownItemData[] = dailyResp?.data || []
  const monthlyBreakdown: MonthlyBreakdownItemData[] = monthlyResp?.data || []
  const yearlyBreakdown: YearlyBreakdownItemData[] = yearlyResp?.data || []
  const productPerf: ProductPerformanceResponseData | undefined = productPerfResp?.data
  const payments: PaymentMethodData[] = paymentsResp?.data || []
  const footwear: FootwearIntelligenceData | undefined = footwearResp?.data

  // Recharts trend data
  const chartTrendsData = trends.map((t) => ({
    date: t.period,
    sales: Number(t.sales || 0),
    netSales: Number(t.net_sales || 0),
    bills: t.bills,
  }))

  // Footwear size data
  const sizesData = (footwear?.top_sizes || []).map((s) => ({
    size: s.attribute_value,
    units: s.units_sold,
    revenue: Number(s.revenue || 0),
    pct: Number(s.percentage || 0),
  }))

  // Footwear colour data
  const coloursData = (footwear?.top_colours || []).map((c) => ({
    name: c.attribute_value,
    value: c.units_sold,
    revenue: Number(c.revenue || 0),
  }))

  const salesGrowth = Number(perf?.sales_growth_pct || 0)
  const billsGrowth = Number(perf?.bills_growth_pct || 0)

  // Filter products by search term
  const filterList = (list?: DetailedProductPerformanceItemData[]) => {
    if (!list) return []
    if (!productSearch.trim()) return list
    const term = productSearch.toLowerCase()
    return list.filter(
      (p) =>
        p.product_name.toLowerCase().includes(term) ||
        p.product_code.toLowerCase().includes(term) ||
        p.brand_name.toLowerCase().includes(term) ||
        (p.barcode && p.barcode.toLowerCase().includes(term))
    )
  }

  return (
    <div className="space-y-6">
      {/* Top Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Sales Performance & Revenue Analytics</h1>
          <p className="text-sm text-slate-500 mt-1">
            Track daily, monthly and yearly sales, net revenue, discounts, returns, margins, and download multi-sheet Excel reports.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="default"
            size="sm"
            onClick={handleExportExcel}
            disabled={isExporting}
            className="bg-emerald-600 hover:bg-emerald-700 text-white shadow-sm flex items-center gap-1.5"
          >
            <FileSpreadsheet className="w-4 h-4" />
            {isExporting ? 'Generating Excel...' : 'Export Excel (.xlsx)'}
          </Button>
          <Button variant="outline" size="sm" onClick={refreshAll}>
            <RefreshCw className="w-4 h-4 mr-1.5" />
            Refresh
          </Button>
        </div>
      </div>

      {/* Date Range Selector Bar */}
      <Card className="border-slate-200">
        <CardContent className="p-3.5 flex flex-wrap items-center justify-between gap-3">
          <div className="flex flex-wrap items-center gap-1.5">
            {[
              { id: 'today', label: 'Today' },
              { id: 'yesterday', label: 'Yesterday' },
              { id: 'daily', label: 'Daily (30D)' },
              { id: 'this_month', label: 'This Month' },
              { id: 'monthly', label: 'Monthly' },
              { id: 'this_year', label: 'This Year' },
              { id: 'yearly', label: 'Yearly' },
              { id: 'custom', label: 'Custom Range' },
            ].map((p) => (
              <button
                key={p.id}
                onClick={() => setDateRange(p.id)}
                className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-colors ${
                  dateRange === p.id
                    ? 'bg-slate-900 text-white shadow-sm'
                    : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                }`}
              >
                {p.label}
              </button>
            ))}
          </div>

          {dateRange === 'custom' && (
            <div className="flex items-center gap-2 text-xs">
              <Calendar className="w-4 h-4 text-slate-400" />
              <input
                type="date"
                value={dateFrom}
                onChange={(e) => setDateFrom(e.target.value)}
                className="px-2.5 py-1 border border-slate-200 rounded-lg text-xs bg-white"
              />
              <span className="text-slate-400">to</span>
              <input
                type="date"
                value={dateTo}
                onChange={(e) => setDateTo(e.target.value)}
                className="px-2.5 py-1 border border-slate-200 rounded-lg text-xs bg-white"
              />
            </div>
          )}
        </CardContent>
      </Card>

      {/* Executive Financial KPI Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Total Sales Amount */}
        <Card className="border-slate-200">
          <CardContent className="p-5">
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium text-slate-500 uppercase tracking-wider">Total Sales Amount</span>
              <div className="p-2 bg-emerald-50 text-emerald-600 rounded-lg">
                <TrendingUp className="w-5 h-5" />
              </div>
            </div>
            <div className="mt-2">
              <div className="text-2xl font-bold text-slate-900 font-mono">
                {isLoadingPerf ? '...' : `₹${Number(perf?.gross_sales || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}`}
              </div>
              <div className="flex items-center gap-1.5 mt-1 text-xs">
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
          </CardContent>
        </Card>

        {/* Net Sales */}
        <Card className="border-slate-200">
          <CardContent className="p-5">
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium text-slate-500 uppercase tracking-wider">Net Sales</span>
              <div className="p-2 bg-blue-50 text-blue-600 rounded-lg">
                <DollarSign className="w-5 h-5" />
              </div>
            </div>
            <div className="mt-2">
              <div className="text-2xl font-bold text-slate-900 font-mono">
                {isLoadingPerf ? '...' : `₹${Number(perf?.net_sales || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}`}
              </div>
              <div className="flex items-center gap-1 mt-1 text-xs text-slate-500">
                <span>Discounts: ₹{Number(perf?.total_discounts || 0).toFixed(0)}</span>
                <span>•</span>
                <span>Returns: ₹{Number(perf?.total_returns || 0).toFixed(0)}</span>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Products Sold */}
        <Card className="border-slate-200">
          <CardContent className="p-5">
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium text-slate-500 uppercase tracking-wider">Products Sold</span>
              <div className="p-2 bg-indigo-50 text-indigo-600 rounded-lg">
                <ShoppingBag className="w-5 h-5" />
              </div>
            </div>
            <div className="mt-2">
              <div className="text-2xl font-bold text-slate-900 font-mono">
                {isLoadingPerf ? '...' : `${perf?.total_units_sold || 0} units`}
              </div>
              <div className="flex items-center gap-1.5 mt-1 text-xs text-slate-500">
                <span>Across {perf?.total_bills || 0} completed invoices</span>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Completed Bills & Average Bill Value */}
        <Card className="border-slate-200">
          <CardContent className="p-5">
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium text-slate-500 uppercase tracking-wider">Average Bill Value</span>
              <div className="p-2 bg-purple-50 text-purple-600 rounded-lg">
                <Receipt className="w-5 h-5" />
              </div>
            </div>
            <div className="mt-2">
              <div className="text-2xl font-bold text-slate-900 font-mono">
                {isLoadingPerf ? '...' : `₹${Number(perf?.average_bill_value || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}`}
              </div>
              <div className="flex items-center gap-1.5 mt-1 text-xs text-slate-500">
                <span>{perf?.total_bills || 0} Completed Bills</span>
                {billsGrowth !== 0 && (
                  <span className={billsGrowth > 0 ? 'text-emerald-600 font-medium' : 'text-rose-600 font-medium'}>
                    ({billsGrowth > 0 ? `+${billsGrowth}%` : `${billsGrowth}%`})
                  </span>
                )}
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Secondary Benchmark Cards (Profit, Margins, Daily Average, Best Day/Month) */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        {/* Gross Profit & Margins */}
        <Card className="border-slate-200 bg-slate-50/50">
          <CardContent className="p-4 flex items-center justify-between">
            <div>
              <span className="text-xs text-slate-500 font-medium block">Gross Profit (Net Sales - COGS)</span>
              <span className="text-lg font-bold text-slate-900 font-mono">
                ₹{Number(perf?.gross_profit || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
              </span>
              <p className="text-[11px] text-slate-500 mt-0.5">
                COGS: ₹{Number(perf?.cogs || 0).toLocaleString('en-IN', { minimumFractionDigits: 0 })}
              </p>
            </div>
            <Badge variant="outline" className="bg-emerald-50 text-emerald-700 border-emerald-200 text-xs font-semibold px-2.5 py-1">
              <Percent className="w-3 h-3 mr-1" />
              {Number(perf?.profit_margin_pct || 0).toFixed(1)}% Margin
            </Badge>
          </CardContent>
        </Card>

        {/* Daily Velocity */}
        <Card className="border-slate-200 bg-slate-50/50">
          <CardContent className="p-4 flex items-center justify-between">
            <div>
              <span className="text-xs text-slate-500 font-medium block">Average Daily Sales</span>
              <span className="text-lg font-bold text-slate-900 font-mono">
                ₹{Number(perf?.average_daily_sales || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
              </span>
              <p className="text-[11px] text-slate-500 mt-0.5">Calculated across selected period window</p>
            </div>
            <div className="p-2 bg-blue-50 text-blue-600 rounded-lg">
              <CalendarDays className="w-5 h-5" />
            </div>
          </CardContent>
        </Card>

        {/* Best Performance Benchmarks */}
        <Card className="border-slate-200 bg-slate-50/50">
          <CardContent className="p-4 flex items-center justify-between">
            <div>
              <span className="text-xs text-slate-500 font-medium block">Best-Selling Milestone</span>
              {perf?.best_selling_day ? (
                <div>
                  <span className="text-sm font-bold text-slate-900">{perf.best_selling_day.label}</span>
                  <p className="text-[11px] text-emerald-600 font-semibold mt-0.5">
                    ₹{Number(perf.best_selling_day.sales).toLocaleString('en-IN')} ({perf.best_selling_day.units} units)
                  </p>
                </div>
              ) : perf?.best_selling_month ? (
                <div>
                  <span className="text-sm font-bold text-slate-900">{perf.best_selling_month.label}</span>
                  <p className="text-[11px] text-emerald-600 font-semibold mt-0.5">
                    ₹{Number(perf.best_selling_month.sales).toLocaleString('en-IN')}
                  </p>
                </div>
              ) : (
                <span className="text-xs text-slate-400">No milestones recorded</span>
              )}
            </div>
            <div className="p-2 bg-amber-50 text-amber-600 rounded-lg">
              <Award className="w-5 h-5" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Analytics Tabs Bar */}
      <div className="flex border-b border-slate-200 gap-2 overflow-x-auto pb-1">
        {[
          { id: 'trends', label: 'Sales Trends', icon: TrendingUp },
          { id: 'daily', label: 'Daily Breakdown', icon: CalendarDays },
          { id: 'monthly', label: 'Monthly Breakdown', icon: Calendar },
          { id: 'yearly', label: 'Yearly Summary', icon: Calendar },
          { id: 'products', label: 'Product Intelligence', icon: ShoppingBag },
          { id: 'footwear', label: 'Footwear Matrix', icon: Layers },
          { id: 'payments', label: 'Payment Tenders', icon: CreditCard },
        ].map((tab) => {
          const Icon = tab.icon
          const isActive = activeTab === tab.id
          return (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id as typeof activeTab)}
              className={`pb-3 px-3 text-xs sm:text-sm font-medium transition-colors border-b-2 flex items-center gap-1.5 whitespace-nowrap ${
                isActive
                  ? 'border-emerald-600 text-emerald-600'
                  : 'border-transparent text-slate-500 hover:text-slate-700'
              }`}
            >
              <Icon className="w-4 h-4" />
              {tab.label}
            </button>
          )
        })}
      </div>

      {/* TAB 1: Sales Trends Chart */}
      {activeTab === 'trends' && (
        <Card className="border-slate-200">
          <CardHeader>
            <CardTitle className="text-base">Sales & Net Revenue Trajectory</CardTitle>
            <CardDescription>Daily revenue trends, invoice volumes, and net bill amounts</CardDescription>
          </CardHeader>
          <CardContent>
            {isLoadingTrends ? (
              <div className="h-72 flex items-center justify-center text-slate-400 text-sm">Loading trend metrics...</div>
            ) : chartTrendsData.length === 0 ? (
              <div className="h-72 flex items-center justify-center text-slate-400 text-sm">
                No sales records found for this period.
              </div>
            ) : (
              <div className="h-80 w-full">
                <ResponsiveContainer width="100%" height="100%">
                  <AreaChart data={chartTrendsData} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
                    <defs>
                      <linearGradient id="colorSales" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#10b981" stopOpacity={0.8} />
                        <stop offset="95%" stopColor="#10b981" stopOpacity={0} />
                      </linearGradient>
                      <linearGradient id="colorNet" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.8} />
                        <stop offset="95%" stopColor="#3b82f6" stopOpacity={0} />
                      </linearGradient>
                    </defs>
                    <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#f1f5f9" />
                    <XAxis dataKey="date" tick={{ fontSize: 11 }} tickLine={false} stroke="#94a3b8" />
                    <YAxis
                      tick={{ fontSize: 11 }}
                      tickLine={false}
                      stroke="#94a3b8"
                      tickFormatter={(val) => `₹${val >= 1000 ? `${(val / 1000).toFixed(0)}k` : val}`}
                    />
                    <Tooltip
                      formatter={(value: any, name: any) => [
                        `₹${Number(value).toLocaleString('en-IN')}`,
                        name === 'sales' ? 'Gross Sales' : 'Net Sales',
                      ]}
                      contentStyle={{ backgroundColor: '#ffffff', borderRadius: '8px', border: '1px solid #e2e8f0', fontSize: '12px' }}
                    />
                    <Legend />
                    <Area type="monotone" dataKey="sales" name="Gross Sales" stroke="#10b981" fillOpacity={1} fill="url(#colorSales)" strokeWidth={2} />
                    <Area type="monotone" dataKey="netSales" name="Net Sales" stroke="#3b82f6" fillOpacity={1} fill="url(#colorNet)" strokeWidth={2} />
                  </AreaChart>
                </ResponsiveContainer>
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* TAB 2: Daily Breakdown Table */}
      {activeTab === 'daily' && (
        <Card className="border-slate-200">
          <CardHeader>
            <CardTitle className="text-base">Daily Sales Ledger & Breakdown</CardTitle>
            <CardDescription>Day-by-day record of invoices, items sold, gross revenue, discounts, returns and net proceeds</CardDescription>
          </CardHeader>
          <CardContent>
            {isLoadingDaily ? (
              <div className="py-12 text-center text-sm text-slate-400">Loading daily transactions...</div>
            ) : dailyBreakdown.length === 0 ? (
              <div className="py-12 text-center text-sm text-slate-400">No daily records for the chosen window.</div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-xs text-left">
                  <thead className="bg-slate-50 text-slate-600 uppercase font-semibold border-b">
                    <tr>
                      <th className="py-3 px-3">Date</th>
                      <th className="py-3 px-3 text-right">Bills</th>
                      <th className="py-3 px-3 text-right">Units Sold</th>
                      <th className="py-3 px-3 text-right">Gross Sales</th>
                      <th className="py-3 px-3 text-right">Discounts</th>
                      <th className="py-3 px-3 text-right">Returns</th>
                      <th className="py-3 px-3 text-right">Taxes</th>
                      <th className="py-3 px-3 text-right">Net Sales</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100 font-mono">
                    {dailyBreakdown.map((row) => (
                      <tr key={row.date} className="hover:bg-slate-50">
                        <td className="py-2.5 px-3 font-medium text-slate-900 font-sans">{row.date}</td>
                        <td className="py-2.5 px-3 text-right">{row.bills}</td>
                        <td className="py-2.5 px-3 text-right font-semibold text-indigo-600">{row.units_sold}</td>
                        <td className="py-2.5 px-3 text-right">₹{Number(row.gross_sales).toLocaleString('en-IN', { minimumFractionDigits: 2 })}</td>
                        <td className="py-2.5 px-3 text-right text-rose-500">₹{Number(row.discounts).toFixed(2)}</td>
                        <td className="py-2.5 px-3 text-right text-amber-600">₹{Number(row.returns).toFixed(2)}</td>
                        <td className="py-2.5 px-3 text-right text-slate-500">₹{Number(row.tax_amount).toFixed(2)}</td>
                        <td className="py-2.5 px-3 text-right font-bold text-emerald-600">
                          ₹{Number(row.net_sales).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* TAB 3: Monthly Breakdown Table */}
      {activeTab === 'monthly' && (
        <Card className="border-slate-200">
          <CardHeader>
            <CardTitle className="text-base">Monthly Sales Performance</CardTitle>
            <CardDescription>Month-on-month revenue progression and unit sales</CardDescription>
          </CardHeader>
          <CardContent>
            {isLoadingMonthly ? (
              <div className="py-12 text-center text-sm text-slate-400">Loading monthly metrics...</div>
            ) : monthlyBreakdown.length === 0 ? (
              <div className="py-12 text-center text-sm text-slate-400">No monthly data available.</div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-xs text-left">
                  <thead className="bg-slate-50 text-slate-600 uppercase font-semibold border-b">
                    <tr>
                      <th className="py-3 px-3">Year / Month</th>
                      <th className="py-3 px-3 text-right">Invoices</th>
                      <th className="py-3 px-3 text-right">Units Sold</th>
                      <th className="py-3 px-3 text-right">Gross Sales</th>
                      <th className="py-3 px-3 text-right">Discounts</th>
                      <th className="py-3 px-3 text-right">Returns</th>
                      <th className="py-3 px-3 text-right">Net Sales</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100 font-mono">
                    {monthlyBreakdown.map((m) => (
                      <tr key={`${m.year}-${m.month}`} className="hover:bg-slate-50">
                        <td className="py-2.5 px-3 font-medium text-slate-900 font-sans">
                          {m.month_name} {m.year}
                        </td>
                        <td className="py-2.5 px-3 text-right">{m.bills}</td>
                        <td className="py-2.5 px-3 text-right font-semibold text-indigo-600">{m.units_sold}</td>
                        <td className="py-2.5 px-3 text-right">₹{Number(m.gross_sales).toLocaleString('en-IN', { minimumFractionDigits: 2 })}</td>
                        <td className="py-2.5 px-3 text-right text-rose-500">₹{Number(m.discounts).toFixed(2)}</td>
                        <td className="py-2.5 px-3 text-right text-amber-600">₹{Number(m.returns).toFixed(2)}</td>
                        <td className="py-2.5 px-3 text-right font-bold text-emerald-600">
                          ₹{Number(m.net_sales).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* TAB 3b: Yearly Summary Table */}
      {activeTab === 'yearly' && (
        <Card className="border-slate-200">
          <CardHeader>
            <CardTitle className="text-base">Yearly Sales Progression</CardTitle>
            <CardDescription>Multi-year revenue, bill volumes, and net performance</CardDescription>
          </CardHeader>
          <CardContent>
            {isLoadingYearly ? (
              <div className="py-12 text-center text-sm text-slate-400">Loading yearly summary...</div>
            ) : yearlyBreakdown.length === 0 ? (
              <div className="py-12 text-center text-sm text-slate-400">No yearly records found.</div>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-xs text-left">
                  <thead className="bg-slate-50 text-slate-600 uppercase font-semibold border-b">
                    <tr>
                      <th className="py-3 px-3">Calendar Year</th>
                      <th className="py-3 px-3 text-right">Invoices</th>
                      <th className="py-3 px-3 text-right">Units Sold</th>
                      <th className="py-3 px-3 text-right">Gross Sales</th>
                      <th className="py-3 px-3 text-right">Discounts</th>
                      <th className="py-3 px-3 text-right">Returns</th>
                      <th className="py-3 px-3 text-right">Net Sales</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-100 font-mono">
                    {yearlyBreakdown.map((y) => (
                      <tr key={y.year} className="hover:bg-slate-50">
                        <td className="py-2.5 px-3 font-semibold text-slate-900 font-sans">{y.year}</td>
                        <td className="py-2.5 px-3 text-right">{y.bills}</td>
                        <td className="py-2.5 px-3 text-right font-semibold text-indigo-600">{y.units_sold}</td>
                        <td className="py-2.5 px-3 text-right">₹{Number(y.gross_sales).toLocaleString('en-IN', { minimumFractionDigits: 2 })}</td>
                        <td className="py-2.5 px-3 text-right text-rose-500">₹{Number(y.discounts).toFixed(2)}</td>
                        <td className="py-2.5 px-3 text-right text-amber-600">₹{Number(y.returns).toFixed(2)}</td>
                        <td className="py-2.5 px-3 text-right font-bold text-emerald-600">
                          ₹{Number(y.net_sales).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* TAB 4: Product Intelligence (Top Selling, Slow Moving, Zero Sales Alert) */}
      {activeTab === 'products' && (
        <div className="space-y-4">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <div className="flex items-center gap-2">
              <button
                onClick={() => setProductSubTab('top')}
                className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-colors ${
                  productSubTab === 'top'
                    ? 'bg-emerald-600 text-white shadow-sm'
                    : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                }`}
              >
                Top-Selling Products ({productPerf?.top_selling.length || 0})
              </button>
              <button
                onClick={() => setProductSubTab('slow')}
                className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-colors ${
                  productSubTab === 'slow'
                    ? 'bg-amber-600 text-white shadow-sm'
                    : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                }`}
              >
                Slow-Moving ({productPerf?.slow_moving.length || 0})
              </button>
              <button
                onClick={() => setProductSubTab('zero')}
                className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-colors flex items-center gap-1 ${
                  productSubTab === 'zero'
                    ? 'bg-rose-600 text-white shadow-sm'
                    : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
                }`}
              >
                <AlertTriangle className="w-3.5 h-3.5" />
                Zero-Sales Alert ({productPerf?.zero_sales.length || 0})
              </button>
            </div>

            <div className="relative">
              <Search className="w-3.5 h-3.5 absolute left-2.5 top-2.5 text-slate-400" />
              <input
                type="text"
                placeholder="Search product, SKU, barcode..."
                value={productSearch}
                onChange={(e) => setProductSearch(e.target.value)}
                className="pl-8 pr-3 py-1.5 border border-slate-200 rounded-lg text-xs bg-white w-64 focus:outline-none focus:ring-1 focus:ring-emerald-500"
              />
            </div>
          </div>

          <Card className="border-slate-200">
            <CardContent className="p-0">
              {isLoadingProdPerf ? (
                <div className="py-12 text-center text-sm text-slate-400">Loading product catalog performance...</div>
              ) : (
                <div className="overflow-x-auto">
                  <table className="w-full text-xs text-left">
                    <thead className="bg-slate-50 text-slate-600 uppercase font-semibold border-b">
                      <tr>
                        <th className="py-3 px-4">Product Details</th>
                        <th className="py-3 px-3">Variant (Size / Colour)</th>
                        <th className="py-3 px-3">Barcode</th>
                        <th className="py-3 px-3 text-right">Units Sold</th>
                        <th className="py-3 px-3 text-right">Gross Sales</th>
                        <th className="py-3 px-3 text-right">Net Sales</th>
                        <th className="py-3 px-3 text-right">Gross Profit</th>
                        <th className="py-3 px-3 text-right">Margin %</th>
                        <th className="py-3 px-3 text-right">Stock on Hand</th>
                        <th className="py-3 px-4 text-center">Status</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100">
                      {filterList(
                        productSubTab === 'top'
                          ? productPerf?.top_selling
                          : productSubTab === 'slow'
                          ? productPerf?.slow_moving
                          : productPerf?.zero_sales
                      ).map((p, idx) => (
                        <tr key={`${p.product_code}-${p.barcode || idx}`} className="hover:bg-slate-50 font-mono">
                          <td className="py-3 px-4 font-sans">
                            <span className="font-semibold text-slate-900 block">{p.product_name}</span>
                            <span className="text-[11px] text-slate-500">
                              {p.product_code} • {p.brand_name} ({p.category_name})
                            </span>
                          </td>
                          <td className="py-3 px-3 font-sans">
                            {p.size || p.colour ? (
                              <Badge variant="outline" className="text-[10px]">
                                {p.size ? `Size ${p.size}` : ''} {p.colour ? `• ${p.colour}` : ''}
                              </Badge>
                            ) : (
                              <span className="text-slate-400 text-xs">Standard</span>
                            )}
                          </td>
                          <td className="py-3 px-3 text-slate-500">{p.barcode || '-'}</td>
                          <td className="py-3 px-3 text-right font-bold text-indigo-600 font-sans">{p.units_sold}</td>
                          <td className="py-3 px-3 text-right">₹{Number(p.gross_sales).toLocaleString('en-IN', { minimumFractionDigits: 2 })}</td>
                          <td className="py-3 px-3 text-right font-semibold text-emerald-600">
                            ₹{Number(p.net_sales).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                          </td>
                          <td className="py-3 px-3 text-right font-semibold text-slate-900">
                            ₹{Number(p.gross_profit).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                          </td>
                          <td className="py-3 px-3 text-right text-slate-600">{Number(p.profit_margin).toFixed(1)}%</td>
                          <td className="py-3 px-3 text-right font-sans">
                            <span className={p.stock_on_hand <= 5 ? 'text-rose-600 font-bold' : 'text-slate-700'}>
                              {p.stock_on_hand} units
                            </span>
                          </td>
                          <td className="py-3 px-4 text-center font-sans">
                            {p.status === 'BEST_SELLER' && (
                              <Badge className="bg-emerald-100 text-emerald-800 border-none text-[10px]">Best Seller</Badge>
                            )}
                            {p.status === 'SLOW_MOVING' && (
                              <Badge className="bg-amber-100 text-amber-800 border-none text-[10px]">Slow Moving</Badge>
                            )}
                            {p.status === 'ZERO_SALES' && (
                              <Badge className="bg-rose-100 text-rose-800 border-none text-[10px]">Zero Sales</Badge>
                            )}
                            {p.status === 'ACTIVE' && (
                              <Badge variant="outline" className="text-[10px]">
                                Active
                              </Badge>
                            )}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      )}

      {/* TAB 5: Footwear Matrix */}
      {activeTab === 'footwear' && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <Card className="border-slate-200">
            <CardHeader>
              <CardTitle className="text-base">Footwear Size Distribution</CardTitle>
              <CardDescription>Units sold and revenue grouped by shoe size curve</CardDescription>
            </CardHeader>
            <CardContent>
              {sizesData.length === 0 ? (
                <div className="h-64 flex items-center justify-center text-slate-400 text-sm">No size data recorded</div>
              ) : (
                <div className="h-72 w-full">
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={sizesData} margin={{ top: 10, right: 10, left: 0, bottom: 0 }}>
                      <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#f1f5f9" />
                      <XAxis dataKey="size" tick={{ fontSize: 11 }} tickLine={false} stroke="#94a3b8" />
                      <YAxis tick={{ fontSize: 11 }} tickLine={false} stroke="#94a3b8" />
                      <Tooltip
                        formatter={(val: any) => [`${val} Units`, 'Sold']}
                        contentStyle={{ backgroundColor: '#ffffff', borderRadius: '8px', border: '1px solid #e2e8f0', fontSize: '12px' }}
                      />
                      <Bar dataKey="units" fill="#3b82f6" radius={[4, 4, 0, 0]} />
                    </BarChart>
                  </ResponsiveContainer>
                </div>
              )}
            </CardContent>
          </Card>

          <Card className="border-slate-200">
            <CardHeader>
              <CardTitle className="text-base">Footwear Colour Preference</CardTitle>
              <CardDescription>Sales distribution by product colour variant</CardDescription>
            </CardHeader>
            <CardContent>
              {coloursData.length === 0 ? (
                <div className="h-64 flex items-center justify-center text-slate-400 text-sm">No colour data recorded</div>
              ) : (
                <div className="h-72 w-full">
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie data={coloursData} cx="50%" cy="50%" outerRadius={85} dataKey="value" nameKey="name" label={(entry) => entry.name}>
                        {coloursData.map((_, index) => (
                          <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                        ))}
                      </Pie>
                      <Tooltip formatter={(val: any) => [`${val} Units`, 'Sold']} />
                      <Legend />
                    </PieChart>
                  </ResponsiveContainer>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      )}

      {/* TAB 6: Payment Tenders */}
      {activeTab === 'payments' && (
        <Card className="border-slate-200">
          <CardHeader>
            <CardTitle className="text-base">Payment Method Breakdown</CardTitle>
            <CardDescription>Collections by Cash, UPI, Card, and other tender modes</CardDescription>
          </CardHeader>
          <CardContent>
            {payments.length === 0 ? (
              <div className="py-12 text-center text-sm text-slate-400">No payment records available.</div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6 items-center">
                <div className="h-64 w-full">
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie
                        data={payments.map((p) => ({ name: p.method, value: Number(p.total_sales) }))}
                        cx="50%"
                        cy="50%"
                        outerRadius={80}
                        dataKey="value"
                        label={(e) => `${e.name}`}
                      >
                        {payments.map((_, idx) => (
                          <Cell key={`pay-${idx}`} fill={COLORS[idx % COLORS.length]} />
                        ))}
                      </Pie>
                      <Tooltip formatter={(val: any) => [`₹${Number(val).toLocaleString('en-IN')}`, 'Collected']} />
                      <Legend />
                    </PieChart>
                  </ResponsiveContainer>
                </div>

                <div className="space-y-3">
                  {payments.map((pm, idx) => (
                    <div key={pm.method} className="flex items-center justify-between p-3 rounded-lg border border-slate-100 bg-slate-50/50">
                      <div className="flex items-center gap-2.5">
                        <div className="w-3 h-3 rounded-full" style={{ backgroundColor: COLORS[idx % COLORS.length] }} />
                        <span className="text-xs font-semibold text-slate-800">{pm.method}</span>
                        <span className="text-[11px] text-slate-500">({pm.bill_count} bills)</span>
                      </div>
                      <div className="text-right">
                        <span className="text-xs font-bold text-slate-900 font-mono block">
                          ₹{Number(pm.total_sales).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                        </span>
                        <span className="text-[10px] text-slate-500 font-semibold">{Number(pm.percentage).toFixed(1)}% Share</span>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  )
}
