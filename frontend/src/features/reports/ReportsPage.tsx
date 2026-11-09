import React from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { FileSpreadsheet } from 'lucide-react'

export const ReportsPage: React.FC = () => (
  <div className="space-y-6">
    <div>
      <h1 className="text-2xl font-bold text-slate-900 tracking-tight">Reports & Exports</h1>
      <p className="text-sm text-slate-500 mt-1">Export sales, stock, and tax records to Excel, CSV, and PDF</p>
    </div>
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <FileSpreadsheet className="w-5 h-5 text-emerald-600" />
          Scheduled & On-Demand Reports
        </CardTitle>
        <CardDescription>Scheduled for Phase 6 implementation</CardDescription>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-slate-600">
          Provides daily, weekly, monthly sales summaries, low-stock reports, and GST summaries.
        </p>
      </CardContent>
    </Card>
  </div>
)

