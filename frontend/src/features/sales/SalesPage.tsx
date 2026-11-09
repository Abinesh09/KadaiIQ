import React from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { TrendingUp } from 'lucide-react'

export const SalesPage: React.FC = () => (
  <div className="space-y-6">
    <div>
      <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Sales & POS Records</h1>
      <p className="text-sm text-slate-500 mt-1">Import external POS data (e.g. ZenBiller), bills, and transactions</p>
    </div>
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <TrendingUp className="w-5 h-5 text-emerald-600" />
          POS Data Import & Sales Ledger
        </CardTitle>
        <CardDescription>Scheduled for Phase 5 implementation</CardDescription>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-slate-600">
          Maintains source POS transaction IDs, handles discounts/taxes/net amounts, and prevents duplicate imports.
        </p>
      </CardContent>
    </Card>
  </div>
)

