import React, { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  apiService,
  type AnalyticsOverviewData,
  type TrendPointData,
  type PaymentMethodData,
  type ProductSalesData,
  type FootwearIntelligenceData,
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
  const [dateRange, setDateRange] = useState<string>('30d')
  const [dateFrom, setDateFrom] = useState<string>('')
  const [dateTo, setDateTo] = useState<string>('')
  const [activeTab, setActiveTab] = useState<'trends' | 'footwear' | 'products' | 'payments'>('trends')

  const queryParams = {
    range: dateRange,
    from: dateRange === 'custom' ? dateFrom : undefined,
    to: dateRange === 'custom' ? dateTo : undefined,
  }

  // 1. Overview KPIs
  const {
    data: overviewResp,
    isLoading: isLoadingOverview,
    refetch: refetchOverview,
  } = useQuery({
    queryKey: ['analytics-overview', queryParams],
    queryFn: () => apiService.getAnalyticsOverview(queryParams),
  })

  // 2. Daily Sales Trends
  const { data: trendsResp, isLoading: isLoadingTrends, refetch: refetchTrends } = useQuery({
    queryKey: ['analytics-trends', queryParams],
    queryFn: () => apiService.getSalesTrends(queryParams),
  })

  // 3. Payment Methods
  const { data: paymentsResp, isLoading: isLoadingPayments, refetch: refetchPayments } = useQuery({
    queryKey: ['analytics-payments', queryParams],
    queryFn: () => apiService.getPaymentMethods(queryParams),
  })

  // 4. Top Products
  const { data: productsResp, isLoading: isLoadingProducts, refetch: refetchProducts } = useQuery({
    queryKey: ['analytics-products', queryParams],
    queryFn: () => apiService.getTopProducts({ ...queryParams, limit: 10 }),
  })

  // 5. Footwear Intelligence
  const { data: footwearResp, isLoading: isLoadingFootwear, refetch: refetchFootwear } = useQuery({
    queryKey: ['analytics-footwear', queryParams],
    queryFn: () => apiService.getFootwearIntelligence(queryParams),
  })

  const refreshAll = () => {
    refetchOverview()
    refetchTrends()
    refetchPayments()
    refetchProducts()
    refetchFootwear()
  }

  const overview: AnalyticsOverviewData | undefined = overviewResp?.data
  const trends: TrendPointData[] = trendsResp?.data || []
  const payments: PaymentMethodData[] = paymentsResp?.data || []
  const topProducts: ProductSalesData[] = productsResp?.data || []
  const footwear: FootwearIntelligenceData | undefined = footwearResp?.data

  // Format Recharts trend data
  const chartTrendsData = trends.map((t) => ({
    date: t.period,
    sales: Number(t.sales || 0),
    netSales: Number(t.net_sales || 0),
    bills: t.bills,
  }))

  // Format Footwear size data
  const sizesData = (footwear?.top_sizes || []).map((s) => ({
    size: s.attribute_value,
    units: s.units_sold,
    revenue: Number(s.revenue || 0),
    pct: Number(s.percentage || 0),
  }))

  // Format Footwear colour data
  const coloursData = (footwear?.top_colours || []).map((c) => ({
    name: c.attribute_value,
    value: c.units_sold,
    revenue: Number(c.revenue || 0),
  }))

  const salesGrowth = Number(overview?.sales_growth_pct || 0)
  const billsGrowth = Number(overview?.bills_growth_pct || 0)

  return (
    <div className="space-y-6">
      {/* Top Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Advanced Retail Analytics</h1>
          <p className="text-sm text-slate-500 mt-1">
            Real-time performance analytics, margin intelligence, and footwear variant trends.
          </p>
        </div>
        <div className="flex items-center gap-2">
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
              { id: '7d', label: 'Last 7 Days' },
              { id: '30d', label: 'Last 30 Days' },
              { id: 'this_month', label: 'This Month' },
              { id: 'last_month', label: 'Last Month' },
              { id: 'custom', label: 'Custom' },
            ].map((p) => (
              <button
                key={p.id}
                onClick={() => setDateRange(p.id)}
                className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-colors ${
                  dateRange === p.id
                    ? 'bg-emerald-600 text-white shadow-sm'
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

      {/* Primary KPI Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {/* Gross Sales */}
        <Card className="border-slate-200">
          <CardContent className="p-5">
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium text-slate-500 uppercase tracking-wider">Gross Sales</span>
              <div className="p-2 bg-emerald-50 text-emerald-600 rounded-lg">
                <TrendingUp className="w-5 h-5" />
              </div>
            </div>
            <div className="mt-2">
              <div className="text-2xl font-bold text-slate-900 font-mono">
                {isLoadingOverview ? '...' : `₹${Number(overview?.gross_sales || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}`}
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
                <span className="text-slate-400">vs prev period</span>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Gross Profit & Margin */}
        <Card className="border-slate-200">
          <CardContent className="p-5">
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium text-slate-500 uppercase tracking-wider">Estimated Profit</span>
              <div className="p-2 bg-blue-50 text-blue-600 rounded-lg">
                <Percent className="w-5 h-5" />
              </div>
            </div>
            <div className="mt-2">
              <div className="text-2xl font-bold text-slate-900 font-mono">
                {isLoadingOverview ? '...' : `₹${Number(overview?.gross_profit || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}`}
              </div>
              <div className="flex items-center gap-1.5 mt-1 text-xs">
                <Badge variant="outline" className="text-[11px] font-semibold text-blue-600 border-blue-200 bg-blue-50/50">
                  {Number(overview?.profit_margin_pct || 0).toFixed(1)}% Margin
                </Badge>
                <span className="text-slate-400">Net after COGS</span>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Total Invoices / Orders */}
        <Card className="border-slate-200">
          <CardContent className="p-5">
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium text-slate-500 uppercase tracking-wider">Bills / Invoices</span>
              <div className="p-2 bg-violet-50 text-violet-600 rounded-lg">
                <Receipt className="w-5 h-5" />
              </div>
            </div>
            <div className="mt-2">
              <div className="text-2xl font-bold text-slate-900 font-mono">
                {isLoadingOverview ? '...' : overview?.total_bills || 0}
              </div>
              <div className="flex items-center gap-1.5 mt-1 text-xs">
                {billsGrowth >= 0 ? (
                  <span className="inline-flex items-center text-emerald-600 font-semibold">
                    <ArrowUpRight className="w-3.5 h-3.5 mr-0.5" />
                    +{billsGrowth}%
                  </span>
                ) : (
                  <span className="inline-flex items-center text-rose-600 font-semibold">
                    <ArrowDownRight className="w-3.5 h-3.5 mr-0.5" />
                    {billsGrowth}%
                  </span>
                )}
                <span className="text-slate-400">{overview?.total_items_sold || 0} total units sold</span>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Average Order Value (AOV) */}
        <Card className="border-slate-200">
          <CardContent className="p-5">
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium text-slate-500 uppercase tracking-wider">Avg Order Value</span>
              <div className="p-2 bg-amber-50 text-amber-600 rounded-lg">
                <ShoppingBag className="w-5 h-5" />
              </div>
            </div>
            <div className="mt-2">
              <div className="text-2xl font-bold text-slate-900 font-mono">
                {isLoadingOverview ? '...' : `₹${Number(overview?.average_order_value || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}`}
              </div>
              <div className="flex items-center gap-1.5 mt-1 text-xs text-slate-500">
                <span>Discounts: ₹{Number(overview?.total_discounts || 0).toFixed(0)}</span>
                <span>•</span>
                <span>GST: ₹{Number(overview?.total_tax || 0).toFixed(0)}</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Analytics Sub-Views Tabs */}
      <div className="flex border-b border-slate-200 gap-4">
        <button
          onClick={() => setActiveTab('trends')}
          className={`pb-3 text-sm font-medium transition-colors border-b-2 flex items-center gap-1.5 ${
            activeTab === 'trends'
              ? 'border-emerald-600 text-emerald-600'
              : 'border-transparent text-slate-500 hover:text-slate-700'
          }`}
        >
          <TrendingUp className="w-4 h-4" />
          Revenue & Sales Trends
        </button>
        <button
          onClick={() => setActiveTab('footwear')}
          className={`pb-3 text-sm font-medium transition-colors border-b-2 flex items-center gap-1.5 ${
            activeTab === 'footwear'
              ? 'border-emerald-600 text-emerald-600'
              : 'border-transparent text-slate-500 hover:text-slate-700'
          }`}
        >
          <Layers className="w-4 h-4" />
          Footwear Variant Intelligence
        </button>
        <button
          onClick={() => setActiveTab('products')}
          className={`pb-3 text-sm font-medium transition-colors border-b-2 flex items-center gap-1.5 ${
            activeTab === 'products'
              ? 'border-emerald-600 text-emerald-600'
              : 'border-transparent text-slate-500 hover:text-slate-700'
          }`}
        >
          <ShoppingBag className="w-4 h-4" />
          Top Products & Brands
        </button>
        <button
          onClick={() => setActiveTab('payments')}
          className={`pb-3 text-sm font-medium transition-colors border-b-2 flex items-center gap-1.5 ${
            activeTab === 'payments'
              ? 'border-emerald-600 text-emerald-600'
              : 'border-transparent text-slate-500 hover:text-slate-700'
          }`}
        >
          <CreditCard className="w-4 h-4" />
          Payment Mode Distribution
        </button>
      </div>

      {/* TAB 1: REVENUE & SALES TRENDS */}
      {activeTab === 'trends' && (
        <div className="space-y-6">
          <Card className="border-slate-200">
            <CardHeader className="py-4 px-6 border-b border-slate-100 flex flex-row items-center justify-between">
              <div>
                <CardTitle className="text-base font-semibold text-slate-800">
                  Gross & Net Sales Timeline
                </CardTitle>
                <CardDescription className="text-xs">
                  Daily revenue performance across selected date range.
                </CardDescription>
              </div>
            </CardHeader>
            <CardContent className="p-6">
              {isLoadingTrends ? (
                <div className="h-72 flex items-center justify-center text-slate-400">
                  <RefreshCw className="w-5 h-5 animate-spin mr-2" /> Loading trends...
                </div>
              ) : chartTrendsData.length === 0 ? (
                <div className="h-72 flex flex-col items-center justify-center text-slate-400">
                  <Receipt className="w-8 h-8 mb-2 opacity-50" />
                  <p className="text-sm font-medium">No sales recorded in this period</p>
                  <p className="text-xs text-slate-400 mt-0.5">Try selecting a different date range</p>
                </div>
              ) : (
                <div className="h-72 w-full">
                  <ResponsiveContainer width="100%" height="100%">
                    <AreaChart data={chartTrendsData} margin={{ top: 10, right: 30, left: 10, bottom: 0 }}>
                      <defs>
                        <linearGradient id="colorSales" x1="0" y1="0" x2="0" y2="1">
                          <stop offset="5%" stopColor="#10b981" stopOpacity={0.4} />
                          <stop offset="95%" stopColor="#10b981" stopOpacity={0} />
                        </linearGradient>
                        <linearGradient id="colorNet" x1="0" y1="0" x2="0" y2="1">
                          <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.4} />
                          <stop offset="95%" stopColor="#3b82f6" stopOpacity={0} />
                        </linearGradient>
                      </defs>
                      <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#f1f5f9" />
                      <XAxis dataKey="date" stroke="#94a3b8" fontSize={11} tickLine={false} />
                      <YAxis
                        stroke="#94a3b8"
                        fontSize={11}
                        tickLine={false}
                        tickFormatter={(val) => `₹${val}`}
                      />
                      <Tooltip
                        contentStyle={{
                          backgroundColor: '#1e293b',
                          borderColor: '#334155',
                          borderRadius: '8px',
                          color: '#fff',
                          fontSize: '12px',
                        }}
                        formatter={(val: unknown) => [
                          `₹${Number(val || 0).toLocaleString('en-IN')}`,
                          '',
                        ]}
                      />
                      <Legend wrapperStyle={{ fontSize: '12px', paddingTop: '10px' }} />
                      <Area
                        type="monotone"
                        dataKey="sales"
                        name="Gross Sales (₹)"
                        stroke="#10b981"
                        strokeWidth={2}
                        fillOpacity={1}
                        fill="url(#colorSales)"
                      />
                      <Area
                        type="monotone"
                        dataKey="netSales"
                        name="Net Sales (₹)"
                        stroke="#3b82f6"
                        strokeWidth={2}
                        fillOpacity={1}
                        fill="url(#colorNet)"
                      />
                    </AreaChart>
                  </ResponsiveContainer>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      )}

      {/* TAB 2: FOOTWEAR VARIANT INTELLIGENCE */}
      {activeTab === 'footwear' && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Size Distribution */}
          <Card className="border-slate-200">
            <CardHeader className="py-4 px-6 border-b border-slate-100">
              <CardTitle className="text-base font-semibold text-slate-800">
                Footwear Size Demand (Units Sold)
              </CardTitle>
              <CardDescription className="text-xs">
                Helps optimize inventory restocking by shoe size curve.
              </CardDescription>
            </CardHeader>
            <CardContent className="p-6">
              {isLoadingFootwear ? (
                <div className="h-64 flex items-center justify-center text-slate-400">Loading sizes...</div>
              ) : sizesData.length === 0 ? (
                <div className="h-64 flex flex-col items-center justify-center text-slate-400">
                  <p className="text-sm">No footwear variant data available</p>
                </div>
              ) : (
                <div className="h-64 w-full">
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={sizesData} margin={{ top: 10, right: 20, left: 0, bottom: 0 }}>
                      <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#f1f5f9" />
                      <XAxis dataKey="size" stroke="#64748b" fontSize={11} />
                      <YAxis stroke="#64748b" fontSize={11} />
                      <Tooltip
                        contentStyle={{
                          backgroundColor: '#1e293b',
                          borderColor: '#334155',
                          borderRadius: '8px',
                          color: '#fff',
                          fontSize: '12px',
                        }}
                        formatter={(val: unknown) => [`${val} units`, 'Demand']}
                      />
                      <Bar dataKey="units" fill="#3b82f6" radius={[4, 4, 0, 0]} name="Units Sold" />
                    </BarChart>
                  </ResponsiveContainer>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Colour Breakdown */}
          <Card className="border-slate-200">
            <CardHeader className="py-4 px-6 border-b border-slate-100">
              <CardTitle className="text-base font-semibold text-slate-800">
                Footwear Colour Distribution
              </CardTitle>
              <CardDescription className="text-xs">
                Popularity breakdown by shoe colourway.
              </CardDescription>
            </CardHeader>
            <CardContent className="p-6">
              {isLoadingFootwear ? (
                <div className="h-64 flex items-center justify-center text-slate-400">Loading colours...</div>
              ) : coloursData.length === 0 ? (
                <div className="h-64 flex flex-col items-center justify-center text-slate-400">
                  <p className="text-sm">No colour distribution available</p>
                </div>
              ) : (
                <div className="h-64 w-full">
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie
                        data={coloursData}
                        cx="50%"
                        cy="50%"
                        innerRadius={55}
                        outerRadius={80}
                        paddingAngle={5}
                        dataKey="value"
                        nameKey="name"
                      >
                        {coloursData.map((_, index) => (
                          <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                        ))}
                      </Pie>
                      <Tooltip
                        contentStyle={{
                          backgroundColor: '#1e293b',
                          borderColor: '#334155',
                          borderRadius: '8px',
                          color: '#fff',
                          fontSize: '12px',
                        }}
                        formatter={(val: unknown) => [`${val} units`, 'Volume']}
                      />
                      <Legend wrapperStyle={{ fontSize: '12px' }} />
                    </PieChart>
                  </ResponsiveContainer>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      )}

      {/* TAB 3: TOP PRODUCTS & BRANDS */}
      {activeTab === 'products' && (
        <Card className="border-slate-200 overflow-hidden">
          <CardHeader className="py-4 px-6 border-b border-slate-100">
            <CardTitle className="text-base font-semibold text-slate-800">
              Top Performing Footwear & Products
            </CardTitle>
            <CardDescription className="text-xs">
              Ranked by total revenue with units sold and estimated gross margin.
            </CardDescription>
          </CardHeader>
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm text-slate-600">
              <thead className="bg-slate-50 text-xs uppercase text-slate-500 border-b border-slate-200">
                <tr>
                  <th className="px-6 py-3 font-medium">Rank</th>
                  <th className="px-6 py-3 font-medium">Product Name</th>
                  <th className="px-6 py-3 font-medium">Brand</th>
                  <th className="px-6 py-3 font-medium">Category</th>
                  <th className="px-6 py-3 font-medium text-center">Units Sold</th>
                  <th className="px-6 py-3 font-medium text-right">Revenue (₹)</th>
                  <th className="px-6 py-3 font-medium text-right">Estimated Profit (₹)</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {isLoadingProducts ? (
                  <tr>
                    <td colSpan={7} className="px-6 py-8 text-center text-slate-400">
                      Loading products...
                    </td>
                  </tr>
                ) : topProducts.length === 0 ? (
                  <tr>
                    <td colSpan={7} className="px-6 py-8 text-center text-slate-400">
                      No product sales recorded in this period.
                    </td>
                  </tr>
                ) : (
                  topProducts.map((p, idx) => (
                    <tr key={idx} className="hover:bg-slate-50/70">
                      <td className="px-6 py-3.5 font-bold text-slate-400">#{idx + 1}</td>
                      <td className="px-6 py-3.5 font-semibold text-slate-900">{p.product_name}</td>
                      <td className="px-6 py-3.5">
                        <Badge variant="outline" className="text-xs">
                          {p.brand_name || 'Generic'}
                        </Badge>
                      </td>
                      <td className="px-6 py-3.5 text-xs text-slate-500">{p.category_name || 'Footwear'}</td>
                      <td className="px-6 py-3.5 text-center font-semibold text-slate-800">
                        {p.units_sold}
                      </td>
                      <td className="px-6 py-3.5 text-right font-mono font-semibold text-slate-900">
                        ₹{Number(p.revenue).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                      </td>
                      <td className="px-6 py-3.5 text-right font-mono font-semibold text-emerald-600">
                        ₹{Number(p.gross_profit).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </Card>
      )}

      {/* TAB 4: PAYMENT METHODS */}
      {activeTab === 'payments' && (
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <Card className="border-slate-200">
            <CardHeader className="py-4 px-6 border-b border-slate-100">
              <CardTitle className="text-base font-semibold text-slate-800">
                Payment Mode Breakdown
              </CardTitle>
              <CardDescription className="text-xs">
                Distribution of sales revenue by transaction tender.
              </CardDescription>
            </CardHeader>
            <CardContent className="p-6">
              {isLoadingPayments ? (
                <div className="h-64 flex items-center justify-center text-slate-400">Loading payments...</div>
              ) : payments.length === 0 ? (
                <div className="h-64 flex flex-col items-center justify-center text-slate-400">
                  <p className="text-sm">No payment records available</p>
                </div>
              ) : (
                <div className="h-64 w-full">
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie
                        data={payments.map((p) => ({
                          name: p.method,
                          value: Number(p.total_sales),
                        }))}
                        cx="50%"
                        cy="50%"
                        innerRadius={55}
                        outerRadius={80}
                        paddingAngle={5}
                        dataKey="value"
                        nameKey="name"
                      >
                        {payments.map((_, index) => (
                          <Cell key={`pay-cell-${index}`} fill={COLORS[index % COLORS.length]} />
                        ))}
                      </Pie>
                      <Tooltip
                        contentStyle={{
                          backgroundColor: '#1e293b',
                          borderColor: '#334155',
                          borderRadius: '8px',
                          color: '#fff',
                          fontSize: '12px',
                        }}
                        formatter={(val: unknown) => [
                          `₹${Number(val || 0).toLocaleString('en-IN')}`,
                          'Sales',
                        ]}
                      />
                      <Legend wrapperStyle={{ fontSize: '12px' }} />
                    </PieChart>
                  </ResponsiveContainer>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Payment Method Stats Table */}
          <Card className="border-slate-200">
            <CardHeader className="py-4 px-6 border-b border-slate-100">
              <CardTitle className="text-base font-semibold text-slate-800">Tender Summary</CardTitle>
              <CardDescription className="text-xs">Volume and revenue per tender mode</CardDescription>
            </CardHeader>
            <CardContent className="p-4">
              <div className="space-y-3">
                {payments.map((p, idx) => (
                  <div key={idx} className="flex items-center justify-between p-3 bg-slate-50 rounded-lg text-xs">
                    <div className="flex items-center gap-2">
                      <div
                        className="w-3 h-3 rounded-full"
                        style={{ backgroundColor: COLORS[idx % COLORS.length] }}
                      />
                      <span className="font-semibold text-slate-800">{p.method}</span>
                      <span className="text-slate-400">({p.bill_count} bills)</span>
                    </div>
                    <div className="text-right">
                      <div className="font-mono font-bold text-slate-900">
                        ₹{Number(p.total_sales).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                      </div>
                      <div className="text-[11px] text-slate-500">{Number(p.percentage).toFixed(1)}% of total</div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  )
}
