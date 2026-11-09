import React, { useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { cn } from '@/lib/utils'
import {
  apiService,
  type SaleData,
  type ImportPreviewData,
  type ImportJobData,
} from '@/services/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Receipt,
  UploadCloud,
  FileSpreadsheet,
  CheckCircle2,
  AlertTriangle,
  Eye,
  RefreshCw,
  Search,
  ArrowRight,
  ShieldCheck,
  PackageCheck,
  TrendingUp,
} from 'lucide-react'

export const SalesPage: React.FC = () => {
  const queryClient = useQueryClient()
  const [searchParams, setSearchParams] = useSearchParams()

  // Tab State: 'bills' | 'imports' (synchronized with URL query)
  const activeTab = searchParams.get('tab') === 'imports' ? 'imports' : 'bills'

  const handleTabChange = (newTab: 'bills' | 'imports') => {
    setSearchParams((prev) => {
      const next = new URLSearchParams(prev)
      if (newTab === 'bills') {
        next.delete('tab')
      } else {
        next.set('tab', newTab)
      }
      return next
    }, { replace: true })
  }

  // Search & Filters
  const [search, setSearch] = useState('')
  const [status, setStatus] = useState('')
  const [dateFrom, setDateFrom] = useState('')
  const [dateTo, setDateTo] = useState('')

  // Modals
  const [showImportModal, setShowImportModal] = useState(false)
  const [selectedSaleId, setSelectedSaleId] = useState<string | null>(null)

  // Import State
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [sourceSystem, setSourceSystem] = useState('ZenBiller')
  const [deductInventory, setDeductInventory] = useState(true)
  const [previewData, setPreviewData] = useState<ImportPreviewData | null>(null)
  const [previewError, setPreviewError] = useState<string | null>(null)
  const [importResult, setImportResult] = useState<ImportJobData | null>(null)

  // Sales Query
  const {
    data: salesData,
    isLoading: isLoadingSales,
    refetch: refetchSales,
  } = useQuery({
    queryKey: ['sales-list', search, status, dateFrom, dateTo],
    queryFn: () =>
      apiService.getSales({
        search,
        status,
        from: dateFrom || undefined,
        to: dateTo || undefined,
      }),
  })

  // Import History Query
  const {
    data: historyData,
    isLoading: isLoadingHistory,
    refetch: refetchHistory,
  } = useQuery({
    queryKey: ['import-history'],
    queryFn: () => apiService.getImportHistory(),
  })

  // Single Sale Query (when modal opened)
  const { data: singleSaleData, isLoading: isLoadingSaleDetail } = useQuery({
    queryKey: ['sale-detail', selectedSaleId],
    queryFn: () => (selectedSaleId ? apiService.getSale(selectedSaleId) : null),
    enabled: !!selectedSaleId,
  })

  // Preview Mutation
  const previewMutation = useMutation({
    mutationFn: (file: File) => apiService.previewSalesImport(file),
    onSuccess: (res) => {
      if (res.data) {
        setPreviewData(res.data)
        setPreviewError(null)
      }
    },
    onError: (err: Error) => {
      setPreviewError(err.message)
      setPreviewData(null)
    },
  })

  // Execute Import Mutation
  const importMutation = useMutation({
    mutationFn: (file: File) =>
      apiService.executeSalesImport({
        file,
        source_system: sourceSystem,
        deduct_inventory: deductInventory,
      }),
    onSuccess: (res) => {
      if (res.data) {
        setImportResult(res.data)
        queryClient.invalidateQueries({ queryKey: ['sales-list'] })
        queryClient.invalidateQueries({ queryKey: ['import-history'] })
        queryClient.invalidateQueries({ queryKey: ['inventory-summary'] })
        queryClient.invalidateQueries({ queryKey: ['products-list'] })
      }
    },
    onError: (err: Error) => {
      setPreviewError(err.message)
    },
  })

  const sales = salesData?.data?.items || []
  const history = historyData?.data || []
  const activeSale = singleSaleData?.data

  // Calculate high-level KPIs from loaded sales & history
  const totalRevenue = sales.reduce((acc, s) => acc + Number(s.grand_total || 0), 0)
  const totalDuplicatesBlocked = history.reduce((acc, h) => acc + (h.duplicate_records || 0), 0)
  const totalBillsCount = salesData?.data?.total_count || sales.length

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      setSelectedFile(file)
      setPreviewData(null)
      setPreviewError(null)
      setImportResult(null)
      previewMutation.mutate(file)
    }
  }

  const handleConfirmImport = () => {
    if (selectedFile) {
      importMutation.mutate(selectedFile)
    }
  }

  const resetImportModal = () => {
    setShowImportModal(false)
    setSelectedFile(null)
    setPreviewData(null)
    setPreviewError(null)
    setImportResult(null)
  }

  return (
    <div className="space-y-6">
      {/* Top Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Sales & POS Records</h1>
          <p className="text-sm text-slate-500 mt-1">
            Import POS transactions (ZenBiller), track sales bills, and verify automated stock deductions.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => {
              refetchSales()
              refetchHistory()
            }}
          >
            <RefreshCw className="w-4 h-4 mr-1.5" />
            Refresh
          </Button>
          <Button
            size="sm"
            className="bg-emerald-600 hover:bg-emerald-700 text-white"
            onClick={() => setShowImportModal(true)}
          >
            <UploadCloud className="w-4 h-4 mr-1.5" />
            Import POS Sales
          </Button>
        </div>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card className="border-slate-200">
          <CardContent className="p-5 flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-slate-500 uppercase tracking-wider">Gross Sales Value</p>
              <h3 className="text-2xl font-bold text-slate-900 mt-1">₹{totalRevenue.toLocaleString('en-IN', { minimumFractionDigits: 2 })}</h3>
              <p className="text-xs text-slate-500 mt-1">Across current view</p>
            </div>
            <div className="p-3 bg-emerald-50 text-emerald-600 rounded-lg">
              <TrendingUp className="w-6 h-6" />
            </div>
          </CardContent>
        </Card>

        <Card className="border-slate-200">
          <CardContent className="p-5 flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-slate-500 uppercase tracking-wider">Total Bills Logged</p>
              <h3 className="text-2xl font-bold text-slate-900 mt-1">{totalBillsCount}</h3>
              <p className="text-xs text-slate-500 mt-1">Multi-tenant protected</p>
            </div>
            <div className="p-3 bg-blue-50 text-blue-600 rounded-lg">
              <Receipt className="w-6 h-6" />
            </div>
          </CardContent>
        </Card>

        <Card className="border-slate-200">
          <CardContent className="p-5 flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-slate-500 uppercase tracking-wider">Stock Deductions</p>
              <h3 className="text-2xl font-bold text-slate-900 mt-1">
                {sales.filter((s) => s.inventory_deducted).length} / {sales.length}
              </h3>
              <p className="text-xs text-slate-500 mt-1">Matched variant ledgers</p>
            </div>
            <div className="p-3 bg-violet-50 text-violet-600 rounded-lg">
              <PackageCheck className="w-6 h-6" />
            </div>
          </CardContent>
        </Card>

        <Card className="border-slate-200">
          <CardContent className="p-5 flex items-center justify-between">
            <div>
              <p className="text-xs font-medium text-slate-500 uppercase tracking-wider">Duplicates Prevented</p>
              <h3 className="text-2xl font-bold text-slate-900 mt-1">{totalDuplicatesBlocked}</h3>
              <p className="text-xs text-emerald-600 font-medium mt-1">Zero double counting</p>
            </div>
            <div className="p-3 bg-amber-50 text-amber-600 rounded-lg">
              <ShieldCheck className="w-6 h-6" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Navigation Tabs (Single-Selection Segmented Control) */}
      <div className="bg-[#F5F7FA] p-1.5 rounded-xl border border-[#E5EAF0] inline-flex items-center gap-1">
        <button
          type="button"
          onClick={() => handleTabChange('bills')}
          className={cn(
            'flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition-all cursor-pointer whitespace-nowrap',
            activeTab === 'bills'
              ? 'bg-white text-[#172B4D] shadow-xs border border-[#E5EAF0]'
              : 'text-[#64748B] hover:text-[#172033] hover:bg-white/60 border border-transparent'
          )}
        >
          <Receipt className={cn('w-3.5 h-3.5', activeTab === 'bills' ? 'text-[#16A36A]' : 'text-[#64748B]')} />
          <span>Sales & Bills Ledger</span>
          <span className="ml-1.5 px-2 py-0.5 rounded-full text-[10px] font-bold bg-[#F5F7FA] text-[#64748B]">
            {totalBillsCount}
          </span>
        </button>
        <button
          type="button"
          onClick={() => handleTabChange('imports')}
          className={cn(
            'flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition-all cursor-pointer whitespace-nowrap',
            activeTab === 'imports'
              ? 'bg-white text-[#172B4D] shadow-xs border border-[#E5EAF0]'
              : 'text-[#64748B] hover:text-[#172033] hover:bg-white/60 border border-transparent'
          )}
        >
          <UploadCloud className={cn('w-3.5 h-3.5', activeTab === 'imports' ? 'text-[#16A36A]' : 'text-[#64748B]')} />
          <span>Import Ledger & Logs</span>
          <span className="ml-1.5 px-2 py-0.5 rounded-full text-[10px] font-bold bg-[#F5F7FA] text-[#64748B]">
            {history.length}
          </span>
        </button>
      </div>

      {/* TAB 1: SALES BILLS LEDGER */}
      {activeTab === 'bills' && (
        <div className="space-y-4">
          {/* Filters Bar */}
          <Card className="border-slate-200">
            <CardContent className="p-4 flex flex-wrap items-center gap-3">
              <div className="relative flex-1 min-w-[200px]">
                <Search className="absolute left-3 top-2.5 h-4 w-4 text-slate-400" />
                <input
                  type="text"
                  placeholder="Search bill number or POS ID..."
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="w-full pl-9 pr-4 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
                />
              </div>

              <select
                value={status}
                onChange={(e) => setStatus(e.target.value)}
                className="px-3 py-2 border border-slate-200 rounded-lg text-sm bg-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
              >
                <option value="">All Statuses</option>
                <option value="COMPLETED">Completed</option>
                <option value="CANCELLED">Cancelled</option>
                <option value="RETURNED">Returned</option>
              </select>

              <div className="flex items-center gap-2">
                <span className="text-xs text-slate-500 font-medium">From:</span>
                <input
                  type="date"
                  value={dateFrom}
                  onChange={(e) => setDateFrom(e.target.value)}
                  className="px-2.5 py-1.5 border border-slate-200 rounded-lg text-xs bg-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
                />
              </div>

              <div className="flex items-center gap-2">
                <span className="text-xs text-slate-500 font-medium">To:</span>
                <input
                  type="date"
                  value={dateTo}
                  onChange={(e) => setDateTo(e.target.value)}
                  className="px-2.5 py-1.5 border border-slate-200 rounded-lg text-xs bg-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
                />
              </div>

              {(search || status || dateFrom || dateTo) && (
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => {
                    setSearch('')
                    setStatus('')
                    setDateFrom('')
                    setDateTo('')
                  }}
                  className="text-xs text-slate-500 hover:text-slate-800"
                >
                  Clear Filters
                </Button>
              )}
            </CardContent>
          </Card>

          {/* Sales Table */}
          <Card className="border-slate-200 overflow-hidden">
            <CardHeader className="py-4 px-6 border-b border-slate-100 flex flex-row items-center justify-between">
              <div>
                <CardTitle className="text-base font-semibold text-slate-800">Invoices & Bills</CardTitle>
                <CardDescription className="text-xs">
                  Showing {sales.length} records. All sales preserve original POS identifiers.
                </CardDescription>
              </div>
            </CardHeader>
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm text-slate-600">
                <thead className="bg-slate-50 text-xs uppercase text-slate-500 border-b border-slate-200">
                  <tr>
                    <th className="px-6 py-3.5 font-medium">Bill Number / POS ID</th>
                    <th className="px-6 py-3.5 font-medium">Sale Date</th>
                    <th className="px-6 py-3.5 font-medium">Payment</th>
                    <th className="px-6 py-3.5 font-medium text-center">Items</th>
                    <th className="px-6 py-3.5 font-medium text-right">Tax (GST)</th>
                    <th className="px-6 py-3.5 font-medium text-right">Grand Total</th>
                    <th className="px-6 py-3.5 font-medium text-center">Status</th>
                    <th className="px-6 py-3.5 font-medium text-center">Inventory</th>
                    <th className="px-6 py-3.5 font-medium text-right">Action</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {isLoadingSales ? (
                    <tr>
                      <td colSpan={9} className="px-6 py-12 text-center text-slate-400">
                        <div className="flex items-center justify-center gap-2">
                          <RefreshCw className="w-4 h-4 animate-spin text-emerald-600" />
                          <span>Loading sales records...</span>
                        </div>
                      </td>
                    </tr>
                  ) : sales.length === 0 ? (
                    <tr>
                      <td colSpan={9} className="px-6 py-12 text-center text-slate-500">
                        <div className="max-w-sm mx-auto space-y-3">
                          <div className="p-3 bg-slate-100 text-slate-400 rounded-full w-12 h-12 mx-auto flex items-center justify-center">
                            <Receipt className="w-6 h-6" />
                          </div>
                          <p className="font-medium text-slate-700">No sales transactions found</p>
                          <p className="text-xs text-slate-400">
                            Import your POS sales export file (CSV or Excel from ZenBiller) to populate this ledger.
                          </p>
                          <Button
                            size="sm"
                            className="bg-emerald-600 hover:bg-emerald-700 text-white mt-2"
                            onClick={() => setShowImportModal(true)}
                          >
                            <UploadCloud className="w-4 h-4 mr-1.5" />
                            Import POS File
                          </Button>
                        </div>
                      </td>
                    </tr>
                  ) : (
                    sales.map((sale: SaleData) => (
                      <tr key={sale.id} className="hover:bg-slate-50/70 transition-colors">
                        <td className="px-6 py-4">
                          <div className="font-semibold text-slate-900 flex items-center gap-2">
                            <span>{sale.bill_number}</span>
                            <Badge variant="outline" className="text-[10px] text-slate-500 font-normal">
                              {sale.source_system}
                            </Badge>
                          </div>
                          {sale.source_bill_id && sale.source_bill_id !== sale.bill_number && (
                            <div className="text-xs text-slate-400 font-mono">ID: {sale.source_bill_id}</div>
                          )}
                        </td>
                        <td className="px-6 py-4 text-xs">
                          <div className="text-slate-800 font-medium">
                            {new Date(sale.sale_date).toLocaleDateString()}
                          </div>
                          <div className="text-slate-400">
                            {new Date(sale.sale_date).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <Badge variant="secondary" className="text-xs">
                            {sale.payment_method || 'Cash'}
                          </Badge>
                        </td>
                        <td className="px-6 py-4 text-center font-medium text-slate-700">
                          {sale.item_count}
                        </td>
                        <td className="px-6 py-4 text-right text-xs text-slate-500 font-mono">
                          ₹{Number(sale.tax_amount || 0).toFixed(2)}
                        </td>
                        <td className="px-6 py-4 text-right font-semibold text-slate-900 font-mono">
                          ₹{Number(sale.grand_total || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                        </td>
                        <td className="px-6 py-4 text-center">
                          {sale.status === 'COMPLETED' ? (
                            <Badge className="bg-emerald-100 text-emerald-800 hover:bg-emerald-100 border-none text-[11px]">
                              Completed
                            </Badge>
                          ) : sale.status === 'CANCELLED' ? (
                            <Badge className="bg-rose-100 text-rose-800 hover:bg-rose-100 border-none text-[11px]">
                              Cancelled
                            </Badge>
                          ) : (
                            <Badge className="bg-amber-100 text-amber-800 hover:bg-amber-100 border-none text-[11px]">
                              {sale.status}
                            </Badge>
                          )}
                        </td>
                        <td className="px-6 py-4 text-center">
                          {sale.inventory_deducted ? (
                            <span className="inline-flex items-center gap-1 text-xs text-emerald-600 font-medium">
                              <CheckCircle2 className="w-3.5 h-3.5" />
                              Deducted
                            </span>
                          ) : (
                            <span className="text-xs text-slate-400">Untracked</span>
                          )}
                        </td>
                        <td className="px-6 py-4 text-right">
                          <Button
                            variant="ghost"
                            size="sm"
                            className="text-xs text-slate-600 hover:text-emerald-700 hover:bg-emerald-50"
                            onClick={() => setSelectedSaleId(sale.id)}
                          >
                            <Eye className="w-3.5 h-3.5 mr-1" />
                            View
                          </Button>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </Card>
        </div>
      )}

      {/* TAB 2: IMPORT HISTORY & LOGS */}
      {activeTab === 'imports' && (
        <Card className="border-slate-200 overflow-hidden">
          <CardHeader className="py-4 px-6 border-b border-slate-100 flex flex-row items-center justify-between">
            <div>
              <CardTitle className="text-base font-semibold text-slate-800">POS Import Jobs Audit Log</CardTitle>
              <CardDescription className="text-xs">
                History of uploaded POS batches, deduplication metrics, and error traces.
              </CardDescription>
            </div>
          </CardHeader>
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm text-slate-600">
              <thead className="bg-slate-50 text-xs uppercase text-slate-500 border-b border-slate-200">
                <tr>
                  <th className="px-6 py-3.5 font-medium">File Name</th>
                  <th className="px-6 py-3.5 font-medium">Source POS</th>
                  <th className="px-6 py-3.5 font-medium">Imported At</th>
                  <th className="px-6 py-3.5 font-medium text-center">Total Rows</th>
                  <th className="px-6 py-3.5 font-medium text-center">Imported</th>
                  <th className="px-6 py-3.5 font-medium text-center">Duplicates Skipped</th>
                  <th className="px-6 py-3.5 font-medium text-center">Failed</th>
                  <th className="px-6 py-3.5 font-medium text-center">Status</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {isLoadingHistory ? (
                  <tr>
                    <td colSpan={8} className="px-6 py-12 text-center text-slate-400">
                      Loading import history...
                    </td>
                  </tr>
                ) : history.length === 0 ? (
                  <tr>
                    <td colSpan={8} className="px-6 py-12 text-center text-slate-500">
                      No import jobs recorded yet.
                    </td>
                  </tr>
                ) : (
                  history.map((job: ImportJobData) => (
                    <tr key={job.id} className="hover:bg-slate-50/70">
                      <td className="px-6 py-4 font-medium text-slate-900 flex items-center gap-2">
                        <FileSpreadsheet className="w-4 h-4 text-emerald-600" />
                        <span>{job.file_name}</span>
                      </td>
                      <td className="px-6 py-4">
                        <Badge variant="outline" className="text-xs">
                          {job.source_system}
                        </Badge>
                      </td>
                      <td className="px-6 py-4 text-xs text-slate-500">
                        {new Date(job.created_at).toLocaleString()}
                      </td>
                      <td className="px-6 py-4 text-center font-medium text-slate-700">
                        {job.total_records}
                      </td>
                      <td className="px-6 py-4 text-center font-semibold text-emerald-600">
                        {job.imported_records}
                      </td>
                      <td className="px-6 py-4 text-center font-medium text-blue-600">
                        {job.duplicate_records > 0 ? (
                          <span className="inline-flex items-center gap-1">
                            <ShieldCheck className="w-3.5 h-3.5 text-blue-500" />
                            {job.duplicate_records}
                          </span>
                        ) : (
                          '0'
                        )}
                      </td>
                      <td className="px-6 py-4 text-center font-medium text-rose-600">
                        {job.failed_records}
                      </td>
                      <td className="px-6 py-4 text-center">
                        {job.status === 'COMPLETED' ? (
                          <Badge className="bg-emerald-100 text-emerald-800 hover:bg-emerald-100 border-none text-[11px]">
                            Completed
                          </Badge>
                        ) : job.status === 'PARTIAL' ? (
                          <Badge className="bg-amber-100 text-amber-800 hover:bg-amber-100 border-none text-[11px]">
                            Partial
                          </Badge>
                        ) : (
                          <Badge className="bg-rose-100 text-rose-800 hover:bg-rose-100 border-none text-[11px]">
                            Failed
                          </Badge>
                        )}
                        {job.error_details && (
                          <div className="text-[11px] text-rose-500 mt-1 max-w-xs truncate" title={job.error_details}>
                            {job.error_details}
                          </div>
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

      {/* BILL DETAILS MODAL */}
      {selectedSaleId && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="bg-white rounded-xl shadow-xl max-w-3xl w-full max-h-[90vh] overflow-y-auto">
            <div className="p-6 border-b border-slate-100 flex items-center justify-between">
              <div>
                <div className="flex items-center gap-2">
                  <h3 className="text-lg font-bold text-slate-900">
                    Bill Details: {activeSale?.bill_number || 'Loading...'}
                  </h3>
                  {activeSale && (
                    <Badge variant="outline" className="text-xs font-normal">
                      {activeSale.source_system}
                    </Badge>
                  )}
                </div>
                <p className="text-xs text-slate-500 mt-0.5">
                  Transaction Date: {activeSale ? new Date(activeSale.sale_date).toLocaleString() : ''}
                </p>
              </div>
              <button
                onClick={() => setSelectedSaleId(null)}
                className="text-slate-400 hover:text-slate-600 text-lg font-semibold"
              >
                ✕
              </button>
            </div>

            <div className="p-6 space-y-6">
              {isLoadingSaleDetail ? (
                <div className="py-12 text-center text-slate-400">Loading line items...</div>
              ) : activeSale ? (
                <>
                  {/* Meta summary badges */}
                  <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 bg-slate-50 p-3.5 rounded-lg text-xs">
                    <div>
                      <span className="text-slate-400 block">Payment Mode:</span>
                      <span className="font-semibold text-slate-800">{activeSale.payment_method || 'Cash'}</span>
                    </div>
                    <div>
                      <span className="text-slate-400 block">Status:</span>
                      <span className="font-semibold text-emerald-600">{activeSale.status}</span>
                    </div>
                    <div>
                      <span className="text-slate-400 block">Inventory Adjusted:</span>
                      <span className="font-semibold text-slate-800">
                        {activeSale.inventory_deducted ? 'Yes (Ledger movement recorded)' : 'No'}
                      </span>
                    </div>
                    <div>
                      <span className="text-slate-400 block">Total Items:</span>
                      <span className="font-semibold text-slate-800">{activeSale.item_count}</span>
                    </div>
                  </div>

                  {/* Line Items Table */}
                  <div>
                    <h4 className="text-sm font-semibold text-slate-800 mb-2">Itemized Breakdown</h4>
                    <div className="border border-slate-200 rounded-lg overflow-hidden">
                      <table className="w-full text-left text-xs text-slate-600">
                        <thead className="bg-slate-50 text-[11px] uppercase text-slate-500 border-b border-slate-200">
                          <tr>
                            <th className="px-4 py-2.5 font-medium">Product / Item</th>
                            <th className="px-4 py-2.5 font-medium">Barcode / SKU</th>
                            <th className="px-4 py-2.5 font-medium text-center">Qty</th>
                            <th className="px-4 py-2.5 font-medium text-right">Price</th>
                            <th className="px-4 py-2.5 font-medium text-right">Disc.</th>
                            <th className="px-4 py-2.5 font-medium text-right">GST %</th>
                            <th className="px-4 py-2.5 font-medium text-right">Total</th>
                          </tr>
                        </thead>
                        <tbody className="divide-y divide-slate-100 font-mono">
                          {activeSale.items && activeSale.items.length > 0 ? (
                            activeSale.items.map((item) => (
                              <tr key={item.id} className="hover:bg-slate-50/50">
                                <td className="px-4 py-3 font-sans font-medium text-slate-900">
                                  {item.product_name}
                                  {item.variant_id && (
                                    <span className="ml-1.5 inline-flex items-center text-[10px] text-emerald-600 font-medium">
                                      ✓ Matched
                                    </span>
                                  )}
                                </td>
                                <td className="px-4 py-3 text-slate-500">{item.barcode || item.sku || '-'}</td>
                                <td className="px-4 py-3 text-center font-sans font-semibold text-slate-800">
                                  {item.quantity}
                                </td>
                                <td className="px-4 py-3 text-right">₹{Number(item.unit_price).toFixed(2)}</td>
                                <td className="px-4 py-3 text-right text-rose-600">
                                  {Number(item.discount) > 0 ? `-₹${Number(item.discount).toFixed(2)}` : '₹0.00'}
                                </td>
                                <td className="px-4 py-3 text-right text-slate-500">
                                  {Number(item.tax_rate).toFixed(1)}%
                                </td>
                                <td className="px-4 py-3 text-right font-semibold text-slate-900">
                                  ₹{Number(item.line_total).toFixed(2)}
                                </td>
                              </tr>
                            ))
                          ) : (
                            <tr>
                              <td colSpan={7} className="px-4 py-4 text-center text-slate-400">
                                No line item details available
                              </td>
                            </tr>
                          )}
                        </tbody>
                      </table>
                    </div>
                  </div>

                  {/* Financial Summary */}
                  <div className="flex justify-end">
                    <div className="w-64 space-y-2 text-xs">
                      <div className="flex justify-between text-slate-600">
                        <span>Subtotal:</span>
                        <span className="font-mono">₹{Number(activeSale.subtotal).toFixed(2)}</span>
                      </div>
                      <div className="flex justify-between text-rose-600">
                        <span>Discounts:</span>
                        <span className="font-mono">-₹{Number(activeSale.discount).toFixed(2)}</span>
                      </div>
                      <div className="flex justify-between text-slate-600">
                        <span>Tax / GST:</span>
                        <span className="font-mono">+₹{Number(activeSale.tax_amount).toFixed(2)}</span>
                      </div>
                      <div className="flex justify-between text-sm font-bold text-slate-900 pt-2 border-t border-slate-200">
                        <span>Grand Total:</span>
                        <span className="font-mono">₹{Number(activeSale.grand_total).toFixed(2)}</span>
                      </div>
                    </div>
                  </div>
                </>
              ) : null}
            </div>

            <div className="p-4 bg-slate-50 border-t border-slate-100 flex justify-end">
              <Button variant="outline" size="sm" onClick={() => setSelectedSaleId(null)}>
                Close
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* POS IMPORT MODAL */}
      {showImportModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
          <div className="bg-white rounded-xl shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
            <div className="p-6 border-b border-slate-100 flex items-center justify-between">
              <div>
                <h3 className="text-lg font-bold text-slate-900 flex items-center gap-2">
                  <UploadCloud className="w-5 h-5 text-emerald-600" />
                  Import POS Sales (ZenBiller / CSV)
                </h3>
                <p className="text-xs text-slate-500 mt-0.5">
                  Upload daily sales export files to sync sales and automate inventory deductions.
                </p>
              </div>
              <button onClick={resetImportModal} className="text-slate-400 hover:text-slate-600 text-lg font-semibold">
                ✕
              </button>
            </div>

            <div className="p-6 space-y-5">
              {importResult ? (
                // SUCCESS RESULT VIEW
                <div className="text-center py-6 space-y-4">
                  <div className="w-12 h-12 bg-emerald-100 text-emerald-600 rounded-full flex items-center justify-center mx-auto">
                    <CheckCircle2 className="w-7 h-7" />
                  </div>
                  <div>
                    <h4 className="text-lg font-bold text-slate-900">Import Batch Finished</h4>
                    <p className="text-xs text-slate-500 mt-1">File: {importResult.file_name}</p>
                  </div>

                  <div className="grid grid-cols-3 gap-3 max-w-md mx-auto pt-2">
                    <div className="p-3 bg-emerald-50 rounded-lg">
                      <span className="text-xs text-emerald-600 block">Imported</span>
                      <span className="text-xl font-bold text-emerald-700">{importResult.imported_records}</span>
                    </div>
                    <div className="p-3 bg-blue-50 rounded-lg">
                      <span className="text-xs text-blue-600 block">Duplicates Prevented</span>
                      <span className="text-xl font-bold text-blue-700">{importResult.duplicate_records}</span>
                    </div>
                    <div className="p-3 bg-rose-50 rounded-lg">
                      <span className="text-xs text-rose-600 block">Failed</span>
                      <span className="text-xl font-bold text-rose-700">{importResult.failed_records}</span>
                    </div>
                  </div>

                  <div className="p-3 bg-slate-50 text-xs text-slate-600 rounded-lg max-w-md mx-auto text-left">
                    <p className="flex items-center gap-1.5 font-medium text-emerald-700 mb-1">
                      <ShieldCheck className="w-4 h-4" />
                      Idempotent Deduplication Guarantee
                    </p>
                    <p className="text-slate-500">
                      Duplicate bills were safely skipped without double stock deduction.
                    </p>
                  </div>

                  <Button
                    size="sm"
                    className="bg-emerald-600 hover:bg-emerald-700 text-white mt-4"
                    onClick={resetImportModal}
                  >
                    Done & View Sales
                  </Button>
                </div>
              ) : (
                // UPLOAD & PREVIEW STEP
                <>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div>
                      <label className="text-xs font-semibold text-slate-700 block mb-1">Source POS System</label>
                      <select
                        value={sourceSystem}
                        onChange={(e) => setSourceSystem(e.target.value)}
                        className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm bg-white focus:outline-none focus:ring-2 focus:ring-emerald-500"
                      >
                        <option value="ZenBiller">ZenBiller POS</option>
                        <option value="Tally">Tally Prime</option>
                        <option value="PineLabs">Pine Labs / MSwipe</option>
                        <option value="RetailPOS">Standard Retail POS</option>
                      </select>
                    </div>

                    <div className="flex items-end">
                      <label className="flex items-center gap-2 text-xs font-medium text-slate-700 cursor-pointer pb-2.5">
                        <input
                          type="checkbox"
                          checked={deductInventory}
                          onChange={(e) => setDeductInventory(e.target.checked)}
                          className="rounded border-slate-300 text-emerald-600 focus:ring-emerald-500 w-4 h-4"
                        />
                        <span>Automatically deduct variant stock</span>
                      </label>
                    </div>
                  </div>

                  {/* File Upload Zone */}
                  <div>
                    <label className="text-xs font-semibold text-slate-700 block mb-1">
                      Upload File (.csv or .xlsx)
                    </label>
                    <div className="border-2 border-dashed border-slate-200 hover:border-emerald-400 rounded-lg p-6 text-center cursor-pointer transition-colors bg-slate-50/50">
                      <input
                        type="file"
                        accept=".csv, .xlsx, .xls"
                        onChange={handleFileChange}
                        className="hidden"
                        id="pos-file-upload"
                      />
                      <label htmlFor="pos-file-upload" className="cursor-pointer">
                        <UploadCloud className="w-8 h-8 text-slate-400 mx-auto mb-2" />
                        {selectedFile ? (
                          <div>
                            <span className="text-sm font-semibold text-emerald-700">{selectedFile.name}</span>
                            <span className="text-xs text-slate-400 block mt-0.5">
                              {(selectedFile.size / 1024).toFixed(1)} KB — Click to change
                            </span>
                          </div>
                        ) : (
                          <div>
                            <span className="text-sm font-medium text-slate-700">Click to select export file</span>
                            <span className="text-xs text-slate-400 block mt-1">Supports CSV, XLSX or XLS</span>
                          </div>
                        )}
                      </label>
                    </div>
                  </div>

                  {/* Preview Loading or Error */}
                  {previewMutation.isPending && (
                    <div className="p-4 bg-emerald-50 text-emerald-800 rounded-lg text-xs flex items-center gap-2">
                      <RefreshCw className="w-4 h-4 animate-spin text-emerald-600" />
                      <span>Parsing file headers and matching variant barcodes...</span>
                    </div>
                  )}

                  {previewError && (
                    <div className="p-3 bg-rose-50 border border-rose-200 text-rose-700 rounded-lg text-xs flex items-center gap-2">
                      <AlertTriangle className="w-4 h-4 flex-shrink-0" />
                      <span>{previewError}</span>
                    </div>
                  )}

                  {/* Preview Table */}
                  {previewData && (
                    <div className="space-y-3">
                      <div className="flex items-center justify-between text-xs">
                        <span className="font-semibold text-slate-800">
                          Preview ({previewData.total_rows} rows detected)
                        </span>
                        <span className="text-slate-500">
                          Detected columns: {previewData.detected_columns?.join(', ')}
                        </span>
                      </div>

                      <div className="border border-slate-200 rounded-lg overflow-hidden max-h-48 overflow-y-auto">
                        <table className="w-full text-left text-xs text-slate-600">
                          <thead className="bg-slate-50 text-[10px] uppercase text-slate-500 border-b border-slate-200">
                            <tr>
                              <th className="px-3 py-2 font-medium">Bill #</th>
                              <th className="px-3 py-2 font-medium">Item Name</th>
                              <th className="px-3 py-2 font-medium">Barcode</th>
                              <th className="px-3 py-2 font-medium text-center">Qty</th>
                              <th className="px-3 py-2 font-medium text-right">Price</th>
                              <th className="px-3 py-2 font-medium text-center">Match</th>
                            </tr>
                          </thead>
                          <tbody className="divide-y divide-slate-100">
                            {(previewData.sample_rows || []).map((row, idx) => (
                              <tr key={idx}>
                                <td className="px-3 py-1.5 font-mono">{row.bill_number}</td>
                                <td className="px-3 py-1.5 font-medium text-slate-900 truncate max-w-[120px]">
                                  {row.product_name}
                                </td>
                                <td className="px-3 py-1.5 font-mono text-slate-500">{row.barcode || '-'}</td>
                                <td className="px-3 py-1.5 text-center">{row.quantity}</td>
                                <td className="px-3 py-1.5 text-right font-mono">₹{row.unit_price}</td>
                                <td className="px-3 py-1.5 text-center">
                                  {row.matched ? (
                                    <span className="text-[10px] text-emerald-600 font-semibold">✓ Found</span>
                                  ) : (
                                    <span className="text-[10px] text-slate-400">Unmatched</span>
                                  )}
                                </td>
                              </tr>
                            ))}
                          </tbody>
                        </table>
                      </div>

                      {/* Deduplication Guarantee callout */}
                      <div className="p-3 bg-blue-50 border border-blue-100 rounded-lg text-xs text-blue-900 flex items-start gap-2">
                        <ShieldCheck className="w-4 h-4 text-blue-600 flex-shrink-0 mt-0.5" />
                        <div>
                          <span className="font-semibold block">Strict Deduplication Guaranteed</span>
                          Bills already imported under ({sourceSystem}) will be ignored without creating duplicate records or double deducting stock.
                        </div>
                      </div>
                    </div>
                  )}
                </>
              )}
            </div>

            {!importResult && (
              <div className="p-4 bg-slate-50 border-t border-slate-100 flex items-center justify-between">
                <Button variant="outline" size="sm" onClick={resetImportModal}>
                  Cancel
                </Button>
                <Button
                  size="sm"
                  className="bg-emerald-600 hover:bg-emerald-700 text-white"
                  disabled={!previewData || importMutation.isPending}
                  onClick={handleConfirmImport}
                >
                  {importMutation.isPending ? (
                    <>
                      <RefreshCw className="w-4 h-4 animate-spin mr-1.5" />
                      Importing...
                    </>
                  ) : (
                    <>
                      <ArrowRight className="w-4 h-4 mr-1.5" />
                      Confirm & Import
                    </>
                  )}
                </Button>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
