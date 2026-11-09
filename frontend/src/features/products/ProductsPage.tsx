import React, { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiService, type ProductVariantData } from '@/services/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Package,
  Plus,
  Search,
  Barcode,
  Layers,
  ChevronDown,
  ChevronRight,
  Archive,
  AlertCircle,
  Scan,
} from 'lucide-react'

export const ProductsPage: React.FC = () => {
  const queryClient = useQueryClient()
  const [search, setSearch] = useState('')
  const [gender, setGender] = useState('')
  const [showAddModal, setShowAddModal] = useState(false)
  const [expandedProduct, setExpandedProduct] = useState<string | null>(null)
  const [barcodeInput, setBarcodeInput] = useState('')
  const [scannedVariant, setScannedVariant] = useState<ProductVariantData | null>(null)
  const [scanError, setScanError] = useState<string | null>(null)

  // Form State
  const [name, setName] = useState('')
  const [productCode, setProductCode] = useState('')
  const [genderChoice, setGenderChoice] = useState('Men')
  const [description, setDescription] = useState('')
  const [selectedSizes, setSelectedSizes] = useState<string[]>(['UK 7', 'UK 8', 'UK 9'])
  const [colour, setColour] = useState('Black')
  const [purchasePrice, setPurchasePrice] = useState(1200)
  const [sellingPrice, setSellingPrice] = useState(2499)
  const [gstRate, setGstRate] = useState(12)
  const [openingStock, setOpeningStock] = useState(10)
  const [formError, setFormError] = useState<string | null>(null)

  const { data: prodData, isLoading } = useQuery({
    queryKey: ['products-list', search, gender],
    queryFn: () => apiService.getProducts({ search, gender }),
  })

  const products = prodData?.data?.items || []

  const createMutation = useMutation({
    mutationFn: apiService.createProduct,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products-list'] })
      queryClient.invalidateQueries({ queryKey: ['inventory-summary'] })
      setShowAddModal(false)
      resetForm()
    },
    onError: (err: Error) => {
      setFormError(err.message)
    },
  })

  const archiveMutation = useMutation({
    mutationFn: apiService.archiveProduct,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['products-list'] })
    },
  })

  const resetForm = () => {
    setName('')
    setProductCode('')
    setDescription('')
    setFormError(null)
  }

  const handleBarcodeSearch = async (e: React.FormEvent) => {
    e.preventDefault()
    setScanError(null)
    setScannedVariant(null)
    if (!barcodeInput.trim()) return

    try {
      const res = await apiService.lookupBarcode(barcodeInput.trim())
      if (res.data) {
        setScannedVariant(res.data)
      }
    } catch (err: unknown) {
      setScanError(err instanceof Error ? err.message : 'Barcode not found')
    }
  }

  const handleCreateSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setFormError(null)

    if (selectedSizes.length === 0) {
      setFormError('Please select at least one footwear size')
      return
    }

    const cleanCode = productCode.trim().toUpperCase()
    const variants = selectedSizes.map((size) => {
      const sizeTag = size.replace(/\s+/g, '')
      return {
        size,
        colour,
        sku: `${cleanCode}-${colour.toUpperCase().slice(0, 3)}-${sizeTag}`,
        barcode: `${Math.floor(100000000000 + Math.random() * 900000000000)}`,
        purchase_price: purchasePrice,
        selling_price: sellingPrice,
        gst_rate: gstRate,
        opening_stock: openingStock,
        min_stock_threshold: 5,
      }
    })

    createMutation.mutate({
      name,
      product_code: cleanCode,
      gender: genderChoice,
      description,
      variants,
    })
  }

  const toggleSize = (size: string) => {
    if (selectedSizes.includes(size)) {
      setSelectedSizes(selectedSizes.filter((s) => s !== size))
    } else {
      setSelectedSizes([...selectedSizes, size])
    }
  }

  return (
    <div className="space-y-6">
      {/* Top Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
            <Package className="w-6 h-6 text-blue-600" />
            Product Catalogue & Variants
          </h1>
          <p className="text-sm text-slate-500 mt-1">
            Footwear & retail products, size/colour matrix, and barcode tracking
          </p>
        </div>

        <Button onClick={() => setShowAddModal(true)} className="bg-blue-600 hover:bg-blue-700">
          <Plus className="w-4 h-4 mr-2" />
          Add Footwear Product
        </Button>
      </div>

      {/* Barcode Scanner Bar */}
      <Card className="bg-gradient-to-r from-blue-50 to-indigo-50 border-blue-200">
        <CardContent className="p-4 flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 rounded-lg bg-blue-600 text-white flex items-center justify-center">
              <Scan className="w-5 h-5" />
            </div>
            <div>
              <span className="text-sm font-semibold text-slate-900 block">Barcode Quick Scanner</span>
              <span className="text-xs text-slate-500">Scan product label or type barcode to inspect variant</span>
            </div>
          </div>

          <form onSubmit={handleBarcodeSearch} className="flex items-center gap-2 w-full md:w-auto">
            <input
              type="text"
              value={barcodeInput}
              onChange={(e) => setBarcodeInput(e.target.value)}
              placeholder="Scan or enter barcode..."
              className="px-3 py-1.5 border border-slate-300 rounded-lg text-sm bg-white focus:outline-none focus:ring-2 focus:ring-blue-600 font-mono w-64"
            />
            <Button type="submit" size="sm" variant="outline">
              Lookup
            </Button>
          </form>
        </CardContent>

        {scannedVariant && (
          <div className="px-6 pb-4 pt-1 border-t border-blue-200 flex items-center justify-between text-xs text-slate-700">
            <div>
              <span className="font-semibold text-slate-900">Found Variant: </span>
              Size {scannedVariant.size} ({scannedVariant.colour}) | SKU: <code className="bg-white px-1.5 py-0.5 rounded border">{scannedVariant.sku}</code>
            </div>
            <div className="flex items-center gap-4">
              <span>MRP: ₹{scannedVariant.selling_price}</span>
              <Badge variant={scannedVariant.stock_status === 'IN_STOCK' ? 'success' : 'destructive'}>
                Stock: {scannedVariant.current_stock}
              </Badge>
            </div>
          </div>
        )}

        {scanError && (
          <div className="px-6 pb-3 pt-1 text-xs text-red-600 font-medium">
            {scanError}
          </div>
        )}
      </Card>

      {/* Filter and Search Bar */}
      <div className="flex flex-col sm:flex-row items-center gap-3">
        <div className="relative flex-1 w-full">
          <Search className="w-4 h-4 text-slate-400 absolute left-3 top-2.5" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search by product name, code, SKU..."
            className="w-full pl-9 pr-4 py-2 border border-slate-200 rounded-lg text-sm bg-white focus:outline-none focus:ring-2 focus:ring-blue-600"
          />
        </div>

        <select
          value={gender}
          onChange={(e) => setGender(e.target.value)}
          className="px-3 py-2 border border-slate-200 rounded-lg text-sm bg-white text-slate-700"
        >
          <option value="">All Genders</option>
          <option value="Men">Men</option>
          <option value="Women">Women</option>
          <option value="Kids">Kids</option>
          <option value="Unisex">Unisex</option>
        </select>
      </div>

      {/* Products Listing Table */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base flex items-center justify-between">
            <span>Catalog Items ({products.length})</span>
            <span className="text-xs font-normal text-slate-500">Tracked at variant-level</span>
          </CardTitle>
          <CardDescription>Click any product row to view its sizes, barcodes, and variant stock</CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="py-8 text-center text-sm text-slate-500">Loading catalog...</div>
          ) : products.length === 0 ? (
            <div className="py-8 text-center text-sm text-slate-500">
              No products found. Click "Add Footwear Product" to add shoes with size variants.
            </div>
          ) : (
            <div className="divide-y divide-slate-100">
              {products.map((p) => {
                const isExpanded = expandedProduct === p.id
                return (
                  <div key={p.id} className="py-3">
                    <div
                      onClick={() => setExpandedProduct(isExpanded ? null : p.id)}
                      className="flex items-center justify-between p-3 rounded-lg hover:bg-slate-50 cursor-pointer transition-colors"
                    >
                      <div className="flex items-center space-x-3">
                        <button className="text-slate-400">
                          {isExpanded ? <ChevronDown className="w-4 h-4" /> : <ChevronRight className="w-4 h-4" />}
                        </button>
                        <div className="w-9 h-9 rounded-lg bg-blue-100 text-blue-700 flex items-center justify-center font-bold text-xs">
                          {p.product_code.slice(0, 3)}
                        </div>
                        <div>
                          <div className="font-semibold text-slate-900 text-sm flex items-center gap-2">
                            {p.name}
                            <Badge variant="outline" className="text-[10px]">{p.gender}</Badge>
                          </div>
                          <div className="text-xs text-slate-400 font-mono">
                            Code: {p.product_code} • {p.variants.length} Variants
                          </div>
                        </div>
                      </div>

                      <div className="flex items-center space-x-4">
                        <div className="text-right">
                          <span className="text-xs font-semibold text-slate-900 block">
                            {p.total_stock} Units in Stock
                          </span>
                          <span className="text-[10px] text-slate-400">
                            {p.variants.some((v) => v.current_stock <= 0) ? (
                              <span className="text-amber-600 font-medium">Some variants out of stock</span>
                            ) : (
                              'Healthy Stock'
                            )}
                          </span>
                        </div>

                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={(e) => {
                            e.stopPropagation()
                            if (confirm(`Archive product "${p.name}"?`)) {
                              archiveMutation.mutate(p.id)
                            }
                          }}
                          className="text-slate-400 hover:text-red-600"
                        >
                          <Archive className="w-4 h-4" />
                        </Button>
                      </div>
                    </div>

                    {/* Expandable Variants Table */}
                    {isExpanded && (
                      <div className="mt-2 ml-10 p-4 bg-slate-50 rounded-xl border border-slate-200 space-y-3">
                        <div className="flex items-center justify-between text-xs font-semibold text-slate-600 border-b pb-2">
                          <span className="flex items-center gap-1.5">
                            <Layers className="w-3.5 h-3.5 text-blue-600" />
                            Footwear Size & Colour Variants
                          </span>
                          <span>GST Rate: 12% (Footwear Standard)</span>
                        </div>

                        <div className="overflow-x-auto">
                          <table className="w-full text-left text-xs text-slate-600">
                            <thead className="text-[11px] uppercase text-slate-400 font-semibold">
                              <tr>
                                <th className="pb-2">Size</th>
                                <th className="pb-2">Colour</th>
                                <th className="pb-2">SKU</th>
                                <th className="pb-2">Barcode</th>
                                <th className="pb-2">Purchase Price</th>
                                <th className="pb-2">Selling Price (MRP)</th>
                                <th className="pb-2">Current Stock</th>
                                <th className="pb-2">Status</th>
                              </tr>
                            </thead>
                            <tbody className="divide-y divide-slate-200">
                              {p.variants.map((v) => (
                                <tr key={v.id} className="hover:bg-white/60">
                                  <td className="py-2 font-bold text-slate-900">{v.size}</td>
                                  <td className="py-2">{v.colour}</td>
                                  <td className="py-2 font-mono text-[11px]">{v.sku}</td>
                                  <td className="py-2 font-mono text-[11px] flex items-center gap-1">
                                    <Barcode className="w-3.5 h-3.5 text-slate-400 inline" />
                                    {v.barcode}
                                  </td>
                                  <td className="py-2">₹{v.purchase_price}</td>
                                  <td className="py-2 font-semibold text-slate-900">₹{v.selling_price}</td>
                                  <td className="py-2 font-bold text-blue-700">{v.current_stock}</td>
                                  <td className="py-2">
                                    <Badge
                                      variant={
                                        v.stock_status === 'IN_STOCK'
                                          ? 'success'
                                          : v.stock_status === 'LOW_STOCK'
                                          ? 'secondary'
                                          : 'destructive'
                                      }
                                      className="text-[10px]"
                                    >
                                      {v.stock_status.replace('_', ' ')}
                                    </Badge>
                                  </td>
                                </tr>
                              ))}
                            </tbody>
                          </table>
                        </div>
                      </div>
                    )}
                  </div>
                )
              })}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Add Footwear Product Modal */}
      {showAddModal && (
        <div className="fixed inset-0 bg-slate-900/50 flex items-center justify-center p-4 z-50 overflow-y-auto">
          <div className="bg-white rounded-xl shadow-xl max-w-2xl w-full p-6 space-y-4 my-8">
            <div className="flex items-center justify-between border-b pb-3">
              <h2 className="text-lg font-bold text-slate-900 flex items-center gap-2">
                <Package className="w-5 h-5 text-blue-600" />
                Add Footwear Product & Variants
              </h2>
              <button
                onClick={() => setShowAddModal(false)}
                className="text-slate-400 hover:text-slate-600 text-lg"
              >
                &times;
              </button>
            </div>

            {formError && (
              <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-xs text-red-700 flex items-center gap-2">
                <AlertCircle className="w-4 h-4 flex-shrink-0" />
                <span>{formError}</span>
              </div>
            )}

            <form onSubmit={handleCreateSubmit} className="space-y-4">
              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Product Name</label>
                  <input
                    type="text"
                    required
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder="e.g. Air Walker Running Shoes"
                    className="w-full px-3 py-1.5 border rounded-lg text-sm"
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Product Code</label>
                  <input
                    type="text"
                    required
                    value={productCode}
                    onChange={(e) => setProductCode(e.target.value)}
                    placeholder="e.g. AW-001"
                    className="w-full px-3 py-1.5 border rounded-lg text-sm font-mono uppercase"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Target Segment</label>
                  <select
                    value={genderChoice}
                    onChange={(e) => setGenderChoice(e.target.value)}
                    className="w-full px-3 py-1.5 border rounded-lg text-sm bg-white"
                  >
                    <option value="Men">Men</option>
                    <option value="Women">Women</option>
                    <option value="Kids">Kids</option>
                    <option value="Unisex">Unisex</option>
                  </select>
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Primary Colour</label>
                  <input
                    type="text"
                    required
                    value={colour}
                    onChange={(e) => setColour(e.target.value)}
                    placeholder="e.g. Black / White"
                    className="w-full px-3 py-1.5 border rounded-lg text-sm"
                  />
                </div>
              </div>

              {/* Footwear Size Multi-Select */}
              <div className="space-y-2 border-t pt-3">
                <label className="text-xs font-semibold text-slate-700 block">
                  Select Footwear Sizes (Generates individual variants & SKUs)
                </label>
                <div className="flex flex-wrap gap-2">
                  {['UK 6', 'UK 7', 'UK 8', 'UK 9', 'UK 10', 'UK 11'].map((sz) => {
                    const isSelected = selectedSizes.includes(sz)
                    return (
                      <button
                        type="button"
                        key={sz}
                        onClick={() => toggleSize(sz)}
                        className={`px-3 py-1.5 rounded-lg text-xs font-semibold border transition-colors ${
                          isSelected
                            ? 'bg-blue-600 text-white border-blue-600'
                            : 'bg-white text-slate-700 border-slate-300 hover:bg-slate-50'
                        }`}
                      >
                        {sz}
                      </button>
                    )
                  })}
                </div>
              </div>

              {/* Pricing & Opening Stock Matrix */}
              <div className="grid grid-cols-4 gap-3 border-t pt-3">
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Purchase Price (₹)</label>
                  <input
                    type="number"
                    required
                    min="0"
                    value={purchasePrice}
                    onChange={(e) => setPurchasePrice(Number(e.target.value))}
                    className="w-full px-3 py-1.5 border rounded-lg text-sm"
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Selling Price MRP (₹)</label>
                  <input
                    type="number"
                    required
                    min="0"
                    value={sellingPrice}
                    onChange={(e) => setSellingPrice(Number(e.target.value))}
                    className="w-full px-3 py-1.5 border rounded-lg text-sm"
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">GST Slab (%)</label>
                  <select
                    value={gstRate}
                    onChange={(e) => setGstRate(Number(e.target.value))}
                    className="w-full px-3 py-1.5 border rounded-lg text-sm bg-white"
                  >
                    <option value={5}>5%</option>
                    <option value={12}>12% (Standard)</option>
                    <option value={18}>18%</option>
                  </select>
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Opening Stock (per size)</label>
                  <input
                    type="number"
                    required
                    min="0"
                    value={openingStock}
                    onChange={(e) => setOpeningStock(Number(e.target.value))}
                    className="w-full px-3 py-1.5 border rounded-lg text-sm"
                  />
                </div>
              </div>

              <div className="p-3 bg-blue-50 rounded-lg text-xs text-blue-800">
                💡 Submitting will generate <strong>{selectedSizes.length} variants</strong> with individual barcodes, SKUs, and record <strong>{selectedSizes.length * openingStock} total units</strong> into the auditable stock ledger.
              </div>

              <div className="flex justify-end gap-3 pt-4 border-t">
                <Button type="button" variant="outline" onClick={() => setShowAddModal(false)}>
                  Cancel
                </Button>
                <Button type="submit" disabled={createMutation.isPending} className="bg-blue-600 hover:bg-blue-700">
                  {createMutation.isPending ? 'Saving Product...' : 'Create Product & Variants'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
