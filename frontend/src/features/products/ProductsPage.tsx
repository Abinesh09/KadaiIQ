import React from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Package } from 'lucide-react'

export const ProductsPage: React.FC = () => (
  <div className="space-y-6">
    <div>
      <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Products Management</h1>
      <p className="text-sm text-slate-500 mt-1">Manage product catalog, sizes, colours, variants, and barcodes</p>
    </div>
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Package className="w-5 h-5 text-blue-600" />
          Product Catalog & Variants
        </CardTitle>
        <CardDescription>Scheduled for Phase 4 implementation</CardDescription>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-slate-600">
          Will support size/colour variant tracking, individual SKUs, barcodes, and Indian footwear GST rates.
        </p>
      </CardContent>
    </Card>
  </div>
)

