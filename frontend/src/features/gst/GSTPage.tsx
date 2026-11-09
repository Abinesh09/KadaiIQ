import React from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Percent } from 'lucide-react'

export const GSTPage: React.FC = () => (
  <div className="space-y-6">
    <div>
      <h1 className="text-2xl font-bold text-slate-900 tracking-tight">GST Configuration</h1>
      <p className="text-sm text-slate-500 mt-1">Tenant-specific tax settings, CGST/SGST/IGST breakdown and product slab rules</p>
    </div>
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Percent className="w-5 h-5 text-purple-600" />
          Indian GST Tax Engine
        </CardTitle>
        <CardDescription>Scheduled for Phase 6 implementation</CardDescription>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-slate-600">
          Enforces per-product GST rates (e.g. footwear price thresholds), preserves transaction-time tax rates, and generates GST filing summaries.
        </p>
      </CardContent>
    </Card>
  </div>
)

