import React, { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import {
  apiService,
  type SalesSummaryRowData,
  type StockValuationRowData,
  type StockMovementReportRowData,
} from '@/services/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Download,
  Calendar,
  RefreshCw,
  Receipt,
  Boxes,
  ArrowRightLeft,
  Search,
} from 'lucide-react'

export const ReportsPage: React.FC = () => {
  const [reportType, setReportType] = useState<'sales' | 'inventory' | 'movements'>('sales')

  // Date filters
  const [dateFrom, setDateFrom] = useState('')
  const [dateTo, setDateTo] = useState('')
  const [searchFilter, setSearchFilter] = useState('')

  // 1. Sales Report Query
  const {
    data: salesReportResp,
    isLoading: isLoadingSales,
    refetch: refetchSales,
  } = useQuery({
    queryKey: ['report-sales', dateFrom, dateTo],
    queryFn: () => apiService.getSalesReport({ from: dateFrom || undefined, to: dateTo || undefined }),
    enabled: reportType === 'sales',
  })

  // 2. Inventory Valuation Report Query
  const {
    data: invReportResp,
    isLoading: isLoadingInv,
    refetch: refetchInv,
  } = useQuery({
    queryKey: ['report-inventory'],
    queryFn: () => apiService.getInventoryReport(),
    enabled: reportType === 'inventory',
  })

  // 3. Stock Movements Report Query
  const {
    data: movReportResp,
    isLoading: isLoadingMov,
    refetch: refetchMov,
  } = useQuery({
    queryKey: ['report-movements', dateFrom, dateTo],
    queryFn: () => apiService.getMovementsReport({ from: dateFrom || undefined, to: dateTo || undefined }),
    enabled: reportType === 'movements',
  })

  const handleDownloadCSV = () => {
    const downloadUrl = apiService.getReportDownloadUrl(
      reportType,
      dateFrom || undefined,
      dateTo || undefined
    )
    // Create an authenticated download anchor
    const token = localStorage.getItem('kadai_token')
    fetch(downloadUrl, {
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    })
      .then((res) => res.blob())
      .then((blob) => {
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `kadaiiq_${reportType}_report_${new Date().toISOString().slice(0, 10)}.csv`
        document.body.appendChild(a)
        a.click()
        a.remove()
        window.URL.revokeObjectURL(url)
      })
      .catch((err) => alert('Download failed: ' + err.message))
  }

  const salesData: SalesSummaryRowData[] = salesReportResp?.data?.data || []
  const invData: StockValuationRowData[] = invReportResp?.data?.data || []
  const movData: StockMovementReportRowData[] = movReportResp?.data?.data || []

  // Client-side search filter
  const filteredSales = salesData.filter(
    (s) => s.bill_number.toLowerCase().includes(searchFilter.toLowerCase()) || s.payment_method.toLowerCase().includes(searchFilter.toLowerCase())
  )

  const filteredInv = invData.filter(
    (i) =>
      i.product_name.toLowerCase().includes(searchFilter.toLowerCase()) ||
      i.barcode.includes(searchFilter) ||
      i.size.toLowerCase().includes(searchFilter.toLowerCase())
  )

  const filteredMov = movData.filter(
    (m) =>
      m.product_name.toLowerCase().includes(searchFilter.toLowerCase()) ||
      m.barcode.includes(searchFilter) ||
      m.movement_type.toLowerCase().includes(searchFilter.toLowerCase())
  )

  return (
    <div className="space-y-6">
      {/* Top Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Reports & Data Exports</h1>
          <p className="text-sm text-slate-500 mt-1">
            Generate and export sales ledgers, stock valuations, and audit trail reports in RFC-4180 CSV format.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => {
              if (reportType === 'sales') refetchSales()
              if (reportType === 'inventory') refetchInv()
              if (reportType === 'movements') refetchMov()
            }}
          >
            <RefreshCw className="w-4 h-4 mr-1.5" />
            Refresh
          </Button>
          <Button
            size="sm"
            className="bg-emerald-600 hover:bg-emerald-700 text-white"
            onClick={handleDownloadCSV}
          >
            <Download className="w-4 h-4 mr-1.5" />
            Export to CSV
          </Button>
        </div>
      </div>

      {/* Report Cards Selector */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <button
          onClick={() => setReportType('sales')}
          className={`p-5 rounded-xl border text-left transition-all ${
            reportType === 'sales'
              ? 'border-emerald-600 bg-emerald-50/50 shadow-sm'
              : 'border-slate-200 bg-white hover:border-slate-300'
          }`}
        >
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold uppercase tracking-wider text-slate-500">Invoices & Sales</span>
            <div className={`p-2 rounded-lg ${reportType === 'sales' ? 'bg-emerald-600 text-white' : 'bg-slate-100 text-slate-600'}`}>
              <Receipt className="w-5 h-5" />
            </div>
          </div>
          <h3 className="text-base font-bold text-slate-900 mt-2">Sales Ledger Report</h3>
          <p className="text-xs text-slate-500 mt-1">
            Bill numbers, payment tenders, discounts, GST tax, and grand totals.
          </p>
        </button>

        <button
          onClick={() => setReportType('inventory')}
          className={`p-5 rounded-xl border text-left transition-all ${
            reportType === 'inventory'
              ? 'border-blue-600 bg-blue-50/50 shadow-sm'
              : 'border-slate-200 bg-white hover:border-slate-300'
          }`}
        >
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold uppercase tracking-wider text-slate-500">Stock & Valuation</span>
            <div className={`p-2 rounded-lg ${reportType === 'inventory' ? 'bg-blue-600 text-white' : 'bg-slate-100 text-slate-600'}`}>
              <Boxes className="w-5 h-5" />
            </div>
          </div>
          <h3 className="text-base font-bold text-slate-900 mt-2">Footwear Valuation Report</h3>
          <p className="text-xs text-slate-500 mt-1">
            Current stock quantities, purchase cost, selling price, and stock valuation.
          </p>
        </button>

        <button
          onClick={() => setReportType('movements')}
          className={`p-5 rounded-xl border text-left transition-all ${
            reportType === 'movements'
              ? 'border-violet-600 bg-violet-50/50 shadow-sm'
              : 'border-slate-200 bg-white hover:border-slate-300'
          }`}
        >
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold uppercase tracking-wider text-slate-500">Audit Trail</span>
            <div className={`p-2 rounded-lg ${reportType === 'movements' ? 'bg-violet-600 text-white' : 'bg-slate-100 text-slate-600'}`}>
              <ArrowRightLeft className="w-5 h-5" />
            </div>
          </div>
          <h3 className="text-base font-bold text-slate-900 mt-2">Stock Movement Audit</h3>
          <p className="text-xs text-slate-500 mt-1">
            Full audit log of adjustments, POS sales deductions, and restocks.
          </p>
        </button>
      </div>

      {/* Date Filters & Search Toolbar */}
      <Card className="border-slate-200">
        <CardContent className="p-4 flex flex-wrap items-center justify-between gap-3 text-xs">
          <div className="flex items-center gap-2 flex-1 min-w-[200px]">
            <Search className="w-4 h-4 text-slate-400" />
            <input
              type="text"
              value={searchFilter}
              onChange={(e) => setSearchFilter(e.target.value)}
              placeholder={`Search ${reportType} report records...`}
              className="w-full px-3 py-1.5 border border-slate-200 rounded-lg text-xs focus:ring-2 focus:ring-emerald-500 focus:outline-none"
            />
          </div>

          {reportType !== 'inventory' && (
            <div className="flex items-center gap-2">
              <Calendar className="w-4 h-4 text-slate-400" />
              <span className="font-semibold text-slate-600">From:</span>
              <input
                type="date"
                value={dateFrom}
                onChange={(e) => setDateFrom(e.target.value)}
                className="px-2.5 py-1 border border-slate-200 rounded-lg text-xs bg-white"
              />
              <span className="font-semibold text-slate-600">To:</span>
              <input
                type="date"
                value={dateTo}
                onChange={(e) => setDateTo(e.target.value)}
                className="px-2.5 py-1 border border-slate-200 rounded-lg text-xs bg-white"
              />
              {(dateFrom || dateTo) && (
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => {
                    setDateFrom('')
                    setDateTo('')
                  }}
                  className="text-xs text-slate-500"
                >
                  Clear
                </Button>
              )}
            </div>
          )}
        </CardContent>
      </Card>

      {/* REPORT 1: SALES LEDGER TABLE */}
      {reportType === 'sales' && (
        <Card className="border-slate-200 overflow-hidden">
          <CardHeader className="py-4 px-6 border-b border-slate-100 flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-base font-semibold text-slate-800">
                Sales Ledger Summary ({filteredSales.length} records)
              </CardTitle>
              <CardDescription className="text-xs">
                Total Value:{' '}
                <span className="font-mono font-bold text-slate-900">
                  ₹{Number(salesReportResp?.data?.summary_value || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                </span>
              </CardDescription>
            </div>
          </CardHeader>
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm text-slate-600">
              <thead className="bg-slate-50 text-xs uppercase text-slate-500 border-b border-slate-200 font-mono">
                <tr>
                  <th className="px-6 py-3 font-medium">Bill Number</th>
                  <th className="px-6 py-3 font-medium">Date & Time</th>
                  <th className="px-6 py-3 font-medium">Tender</th>
                  <th className="px-6 py-3 font-medium text-center">Items</th>
                  <th className="px-6 py-3 font-medium text-right">Subtotal</th>
                  <th className="px-6 py-3 font-medium text-right">Discount</th>
                  <th className="px-6 py-3 font-medium text-right">Tax (GST)</th>
                  <th className="px-6 py-3 font-medium text-right">Grand Total</th>
                  <th className="px-6 py-3 font-medium text-center">Status</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 font-mono text-xs">
                {isLoadingSales ? (
                  <tr>
                    <td colSpan={9} className="px-6 py-8 text-center text-slate-400">Loading sales report...</td>
                  </tr>
                ) : filteredSales.length === 0 ? (
                  <tr>
                    <td colSpan={9} className="px-6 py-8 text-center text-slate-400">No sales records in selected period.</td>
                  </tr>
                ) : (
                  filteredSales.map((s, idx) => (
                    <tr key={idx} className="hover:bg-slate-50/70">
                      <td className="px-6 py-3.5 font-bold text-slate-900">{s.bill_number}</td>
                      <td className="px-6 py-3.5 text-slate-500">{new Date(s.sale_date).toLocaleString()}</td>
                      <td className="px-6 py-3.5 font-sans">
                        <Badge variant="secondary" className="text-[11px]">{s.payment_method}</Badge>
                      </td>
                      <td className="px-6 py-3.5 text-center font-bold text-slate-800">{s.item_count}</td>
                      <td className="px-6 py-3.5 text-right font-medium text-slate-900">
                        ₹{Number(s.subtotal).toFixed(2)}
                      </td>
                      <td className="px-6 py-3.5 text-right text-rose-600">
                        -₹{Number(s.discount).toFixed(2)}
                      </td>
                      <td className="px-6 py-3.5 text-right text-slate-500">
                        +₹{Number(s.tax_amount).toFixed(2)}
                      </td>
                      <td className="px-6 py-3.5 text-right font-bold text-slate-900">
                        ₹{Number(s.grand_total).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                      </td>
                      <td className="px-6 py-3.5 text-center font-sans">
                        <Badge className="bg-emerald-100 text-emerald-800 border-none text-[10px]">
                          {s.status}
                        </Badge>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </Card>
      )}

      {/* REPORT 2: STOCK VALUATION TABLE */}
      {reportType === 'inventory' && (
        <Card className="border-slate-200 overflow-hidden">
          <CardHeader className="py-4 px-6 border-b border-slate-100 flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-base font-semibold text-slate-800">
                Footwear Stock Valuation Report ({filteredInv.length} variants)
              </CardTitle>
              <CardDescription className="text-xs">
                Total Valuation at Cost:{' '}
                <span className="font-mono font-bold text-slate-900">
                  ₹{Number(invReportResp?.data?.summary_value || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                </span>
              </CardDescription>
            </div>
          </CardHeader>
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm text-slate-600">
              <thead className="bg-slate-50 text-xs uppercase text-slate-500 border-b border-slate-200 font-mono">
                <tr>
                  <th className="px-6 py-3 font-medium font-sans">Product Name</th>
                  <th className="px-6 py-3 font-medium">Size</th>
                  <th className="px-6 py-3 font-medium font-sans">Colour</th>
                  <th className="px-6 py-3 font-medium">Barcode</th>
                  <th className="px-6 py-3 font-medium text-center">Current Stock</th>
                  <th className="px-6 py-3 font-medium text-right">Cost (₹)</th>
                  <th className="px-6 py-3 font-medium text-right">Retail (₹)</th>
                  <th className="px-6 py-3 font-medium text-right">Valuation at Cost</th>
                  <th className="px-6 py-3 font-medium text-center font-sans">Status</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 font-mono text-xs">
                {isLoadingInv ? (
                  <tr>
                    <td colSpan={9} className="px-6 py-8 text-center text-slate-400">Loading stock valuation...</td>
                  </tr>
                ) : filteredInv.length === 0 ? (
                  <tr>
                    <td colSpan={9} className="px-6 py-8 text-center text-slate-400">No inventory records found.</td>
                  </tr>
                ) : (
                  filteredInv.map((item, idx) => (
                    <tr key={idx} className="hover:bg-slate-50/70">
                      <td className="px-6 py-3.5 font-sans font-semibold text-slate-900">
                        {item.product_name}
                        <span className="block text-[11px] text-slate-400 font-mono">{item.product_code}</span>
                      </td>
                      <td className="px-6 py-3.5 font-bold text-slate-800">{item.size}</td>
                      <td className="px-6 py-3.5 font-sans text-slate-600">{item.colour}</td>
                      <td className="px-6 py-3.5 text-slate-500">{item.barcode}</td>
                      <td className="px-6 py-3.5 text-center font-bold text-slate-900">{item.current_stock}</td>
                      <td className="px-6 py-3.5 text-right">₹{Number(item.purchase_price).toFixed(2)}</td>
                      <td className="px-6 py-3.5 text-right font-medium text-slate-900">
                        ₹{Number(item.selling_price).toFixed(2)}
                      </td>
                      <td className="px-6 py-3.5 text-right font-bold text-emerald-700">
                        ₹{Number(item.cost_valuation).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                      </td>
                      <td className="px-6 py-3.5 text-center font-sans">
                        {item.stock_status === 'IN_STOCK' ? (
                          <Badge className="bg-emerald-100 text-emerald-800 border-none text-[10px]">In Stock</Badge>
                        ) : item.stock_status === 'LOW_STOCK' ? (
                          <Badge className="bg-amber-100 text-amber-800 border-none text-[10px]">Low Stock</Badge>
                        ) : (
                          <Badge className="bg-rose-100 text-rose-800 border-none text-[10px]">Out of Stock</Badge>
                        )}
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </Card>
      )}

      {/* REPORT 3: STOCK MOVEMENTS TABLE */}
      {reportType === 'movements' && (
        <Card className="border-slate-200 overflow-hidden">
          <CardHeader className="py-4 px-6 border-b border-slate-100">
            <CardTitle className="text-base font-semibold text-slate-800">
              Stock Movement Audit Ledger ({filteredMov.length} transactions)
            </CardTitle>
            <CardDescription className="text-xs">
              Immutable ledger trail tracking every addition, restock, sale deduction, and return.
            </CardDescription>
          </CardHeader>
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm text-slate-600">
              <thead className="bg-slate-50 text-xs uppercase text-slate-500 border-b border-slate-200 font-mono">
                <tr>
                  <th className="px-6 py-3 font-medium">Timestamp</th>
                  <th className="px-6 py-3 font-medium font-sans">Product Name</th>
                  <th className="px-6 py-3 font-medium font-sans">Details</th>
                  <th className="px-6 py-3 font-medium font-sans">Movement Type</th>
                  <th className="px-6 py-3 font-medium text-center">Qty</th>
                  <th className="px-6 py-3 font-medium text-center">Prev → New</th>
                  <th className="px-6 py-3 font-medium">Reference #</th>
                  <th className="px-6 py-3 font-medium font-sans">Operator</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 font-mono text-xs">
                {isLoadingMov ? (
                  <tr>
                    <td colSpan={8} className="px-6 py-8 text-center text-slate-400">Loading audit trail...</td>
                  </tr>
                ) : filteredMov.length === 0 ? (
                  <tr>
                    <td colSpan={8} className="px-6 py-8 text-center text-slate-400">No stock movements recorded in this period.</td>
                  </tr>
                ) : (
                  filteredMov.map((m, idx) => (
                    <tr key={idx} className="hover:bg-slate-50/70">
                      <td className="px-6 py-3.5 text-slate-500">{new Date(m.date).toLocaleString()}</td>
                      <td className="px-6 py-3.5 font-sans font-semibold text-slate-900">{m.product_name}</td>
                      <td className="px-6 py-3.5 font-sans text-slate-600">{m.variant_details}</td>
                      <td className="px-6 py-3.5 font-sans">
                        <Badge
                          variant="outline"
                          className={
                            m.movement_type === 'SALE'
                              ? 'text-rose-700 bg-rose-50 border-rose-200 text-[10px]'
                              : m.movement_type === 'PURCHASE'
                              ? 'text-emerald-700 bg-emerald-50 border-emerald-200 text-[10px]'
                              : 'text-blue-700 bg-blue-50 border-blue-200 text-[10px]'
                          }
                        >
                          {m.movement_type}
                        </Badge>
                      </td>
                      <td className="px-6 py-3.5 text-center font-bold text-slate-900">
                        {m.quantity > 0 ? `+${m.quantity}` : m.quantity}
                      </td>
                      <td className="px-6 py-3.5 text-center text-slate-500">
                        {m.previous_stock} → {m.new_stock}
                      </td>
                      <td className="px-6 py-3.5 text-slate-600">{m.reference_number || '-'}</td>
                      <td className="px-6 py-3.5 font-sans text-slate-500">{m.created_by}</td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </Card>
      )}
    </div>
  )
}
