import React, { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiService, type ProductData } from '@/services/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Boxes,
  PlusCircle,
  TrendingDown,
  AlertTriangle,
  History,
  Layers,
  Receipt,
  AlertCircle,
  ArrowDownRight,
  ArrowUpRight,
} from 'lucide-react'

export const InventoryPage: React.FC = () => {
  const queryClient = useQueryClient()
  const [showAdjustModal, setShowAdjustModal] = useState(false)
  const [selectedVariantId, setSelectedVariantId] = useState('')
  const [movementType, setMovementType] = useState('PURCHASE')
  const [quantity, setQuantity] = useState(10)
  const [unitCost, setUnitCost] = useState<number | undefined>(undefined)
  const [reference, setReference] = useState('')
  const [notes, setNotes] = useState('')
  const [error, setError] = useState<string | null>(null)

  // Fetch Inventory KPI Summary
  const { data: summaryResp, isLoading: isSummaryLoading } = useQuery({
    queryKey: ['inventory-summary'],
    queryFn: apiService.getInventorySummary,
  })

  // Fetch Stock Movement History
  const { data: movementsResp, isLoading: isMovementsLoading } = useQuery({
    queryKey: ['inventory-movements'],
    queryFn: () => apiService.getStockMovements(),
  })

  // Fetch Products & Variants to populate adjustment selector
  const { data: productsResp } = useQuery({
    queryKey: ['products-list-for-adjust'],
    queryFn: () => apiService.getProducts(),
  })

  const summary = summaryResp?.data
  const movements = movementsResp?.data?.items || []
  const products: ProductData[] = productsResp?.data?.items || []

  const adjustMutation = useMutation({
    mutationFn: apiService.recordStockAdjustment,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['inventory-summary'] })
      queryClient.invalidateQueries({ queryKey: ['inventory-movements'] })
      queryClient.invalidateQueries({ queryKey: ['products-list'] })
      setShowAdjustModal(false)
      setReference('')
      setNotes('')
    },
    onError: (err: Error) => {
      setError(err.message)
    },
  })

  const handleAdjustSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    if (!selectedVariantId) {
      setError('Please select a product variant to adjust')
      return
    }

    adjustMutation.mutate({
      variant_id: selectedVariantId,
      movement_type: movementType,
      quantity,
      unit_cost: unitCost,
      reference_number: reference,
      notes,
    })
  }

  return (
    <div className="space-y-6">
      {/* Page Title & Quick Action */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
            <Boxes className="w-6 h-6 text-indigo-600" />
            Inventory & Stock Movement Ledger
          </h1>
          <p className="text-sm text-slate-500 mt-1">
            Real-time variant stock tracking, auditable movement ledger, and valuation
          </p>
        </div>

        <Button onClick={() => setShowAdjustModal(true)} className="bg-indigo-600 hover:bg-indigo-700">
          <PlusCircle className="w-4 h-4 mr-2" />
          Record Stock Adjustment
        </Button>
      </div>

      {/* KPI Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-xs font-semibold uppercase text-slate-500">Total Stock Units</CardTitle>
            <Boxes className="w-4 h-4 text-indigo-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900">
              {isSummaryLoading ? '...' : `${summary?.total_stock_units || 0} Units`}
            </div>
            <p className="text-xs text-slate-500 mt-1">Across {summary?.total_variants || 0} active variants</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-xs font-semibold uppercase text-slate-500">Stock Valuation (Cost)</CardTitle>
            <Receipt className="w-4 h-4 text-emerald-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-slate-900">
              {isSummaryLoading ? '...' : `₹${Number(summary?.stock_valuation || 0).toLocaleString('en-IN')}`}
            </div>
            <p className="text-xs text-slate-500 mt-1">Based on variant purchase price</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-xs font-semibold uppercase text-slate-500">Low Stock Variants</CardTitle>
            <AlertTriangle className="w-4 h-4 text-amber-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-amber-600">
              {isSummaryLoading ? '...' : summary?.low_stock_count || 0}
            </div>
            <p className="text-xs text-slate-500 mt-1">Below minimum threshold</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-xs font-semibold uppercase text-slate-500">Out of Stock</CardTitle>
            <TrendingDown className="w-4 h-4 text-red-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">
              {isSummaryLoading ? '...' : summary?.out_of_stock_count || 0}
            </div>
            <p className="text-xs text-slate-500 mt-1">Requires immediate restocking</p>
          </CardContent>
        </Card>
      </div>

      {/* Stock Movement Ledger Table */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base flex items-center justify-between">
            <span className="flex items-center gap-2">
              <History className="w-4 h-4 text-indigo-600" />
              Stock Movement Ledger ({movements.length})
            </span>
            <span className="text-xs font-normal text-slate-500">
              Immutable chronological record of all stock changes
            </span>
          </CardTitle>
          <CardDescription>
            Enforces that stock quantities are never modified without a traceable movement entry
          </CardDescription>
        </CardHeader>
        <CardContent>
          {isMovementsLoading ? (
            <div className="py-8 text-center text-sm text-slate-500">Loading movement ledger...</div>
          ) : movements.length === 0 ? (
            <div className="py-8 text-center text-sm text-slate-500">
              No stock movements recorded yet. Create products with opening stock or record adjustments.
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm text-slate-600">
                <thead className="bg-slate-50 text-xs uppercase text-slate-400 font-semibold border-b">
                  <tr>
                    <th className="px-4 py-3">Timestamp</th>
                    <th className="px-4 py-3">Product / Variant</th>
                    <th className="px-4 py-3">Movement Type</th>
                    <th className="px-4 py-3 text-right">Change</th>
                    <th className="px-4 py-3 text-right">Balance</th>
                    <th className="px-4 py-3">Reference / Notes</th>
                    <th className="px-4 py-3">Logged By</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {movements.map((m) => {
                    const isAddition = m.quantity > 0 && m.new_stock > m.previous_stock
                    return (
                      <tr key={m.id} className="hover:bg-slate-50 text-xs">
                        <td className="px-4 py-3 text-slate-500 whitespace-nowrap">
                          {new Date(m.created_at).toLocaleString()}
                        </td>
                        <td className="px-4 py-3 font-medium text-slate-900">
                          <div>{m.product_name}</div>
                          <div className="text-[11px] text-slate-400 font-mono">
                            {m.variant_details} • {m.barcode}
                          </div>
                        </td>
                        <td className="px-4 py-3">
                          <Badge variant="outline" className="text-[10px] font-semibold">
                            {m.movement_type}
                          </Badge>
                        </td>
                        <td className="px-4 py-3 text-right font-semibold">
                          <span className={`inline-flex items-center gap-0.5 ${isAddition ? 'text-emerald-600' : 'text-red-600'}`}>
                            {isAddition ? <ArrowUpRight className="w-3.5 h-3.5" /> : <ArrowDownRight className="w-3.5 h-3.5" />}
                            {isAddition ? `+${m.quantity}` : `-${m.quantity}`}
                          </span>
                        </td>
                        <td className="px-4 py-3 text-right font-bold text-slate-800">
                          {m.previous_stock} → {m.new_stock}
                        </td>
                        <td className="px-4 py-3">
                          <div className="font-mono text-slate-700">{m.reference_number || '-'}</div>
                          <div className="text-slate-400 text-[10px] truncate max-w-xs">{m.notes || ''}</div>
                        </td>
                        <td className="px-4 py-3 text-slate-600">{m.created_by_name}</td>
                      </tr>
                    )
                  })}
                </tbody>
              </table>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Record Stock Adjustment Modal */}
      {showAdjustModal && (
        <div className="fixed inset-0 bg-slate-900/50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-xl shadow-xl max-w-md w-full p-6 space-y-4">
            <div className="flex items-center justify-between border-b pb-3">
              <h2 className="text-lg font-bold text-slate-900 flex items-center gap-2">
                <Layers className="w-5 h-5 text-indigo-600" />
                Record Stock Movement
              </h2>
              <button
                onClick={() => setShowAdjustModal(false)}
                className="text-slate-400 hover:text-slate-600 text-lg"
              >
                &times;
              </button>
            </div>

            {error && (
              <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-xs text-red-700 flex items-center gap-2">
                <AlertCircle className="w-4 h-4 flex-shrink-0" />
                <span>{error}</span>
              </div>
            )}

            <form onSubmit={handleAdjustSubmit} className="space-y-4">
              <div className="space-y-1">
                <label className="text-xs font-semibold text-slate-700">Select Footwear Variant</label>
                <select
                  value={selectedVariantId}
                  onChange={(e) => setSelectedVariantId(e.target.value)}
                  required
                  className="w-full px-3 py-2 border rounded-lg text-xs bg-white"
                >
                  <option value="">-- Choose Product Variant --</option>
                  {products.flatMap((p) =>
                    p.variants.map((v) => (
                      <option key={v.id} value={v.id}>
                        {p.name} - {v.size} ({v.colour}) [Stock: {v.current_stock}]
                      </option>
                    ))
                  )}
                </select>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Movement Type</label>
                  <select
                    value={movementType}
                    onChange={(e) => setMovementType(e.target.value)}
                    className="w-full px-3 py-1.5 border rounded-lg text-xs bg-white"
                  >
                    <option value="PURCHASE">Purchase / Restock (+)</option>
                    <option value="CUSTOMER_RETURN">Sales Return (+)</option>
                    <option value="SALE">Sale Deduction (-)</option>
                    <option value="MANUAL_ADJUSTMENT">Manual Adjustment</option>
                  </select>
                </div>

                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Units Quantity</label>
                  <input
                    type="number"
                    required
                    min="1"
                    value={quantity}
                    onChange={(e) => setQuantity(Number(e.target.value))}
                    className="w-full px-3 py-1.5 border rounded-lg text-sm"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Unit Cost (₹ Optional)</label>
                  <input
                    type="number"
                    min="0"
                    value={unitCost || ''}
                    onChange={(e) => setUnitCost(e.target.value ? Number(e.target.value) : undefined)}
                    placeholder="Defaults to purchase cost"
                    className="w-full px-3 py-1.5 border rounded-lg text-sm"
                  />
                </div>

                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Invoice / Ref No.</label>
                  <input
                    type="text"
                    value={reference}
                    onChange={(e) => setReference(e.target.value)}
                    placeholder="e.g. PO-9812 / INV-04"
                    className="w-full px-3 py-1.5 border rounded-lg text-sm"
                  />
                </div>
              </div>

              <div className="space-y-1">
                <label className="text-xs font-semibold text-slate-700">Adjustment Notes</label>
                <textarea
                  rows={2}
                  value={notes}
                  onChange={(e) => setNotes(e.target.value)}
                  placeholder="Reason for adjustment, supplier name, or receipt note..."
                  className="w-full px-3 py-1.5 border rounded-lg text-xs"
                />
              </div>

              <div className="flex justify-end gap-3 pt-3 border-t">
                <Button type="button" variant="outline" onClick={() => setShowAdjustModal(false)}>
                  Cancel
                </Button>
                <Button type="submit" disabled={adjustMutation.isPending} className="bg-indigo-600 hover:bg-indigo-700">
                  {adjustMutation.isPending ? 'Recording Movement...' : 'Save Movement'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
