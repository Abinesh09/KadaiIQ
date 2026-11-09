import React from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Boxes } from 'lucide-react'

export const InventoryPage: React.FC = () => (
  <div className="space-y-6">
    <div>
      <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Inventory Management</h1>
      <p className="text-sm text-slate-500 mt-1">Stock movement ledger, variant-level tracking, and low-stock alerts</p>
    </div>
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Boxes className="w-5 h-5 text-indigo-600" />
          Stock Movement Ledger
        </CardTitle>
        <CardDescription>Scheduled for Phase 4 implementation</CardDescription>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-slate-600">
          Enforces mandatory movement history for opening stock, restocks, sales, returns, and adjustments.
        </p>
      </CardContent>
    </Card>
  </div>
)

