import React, { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiService, type GSTSettingsData, type GSTR1SummaryData, type GSTR3BTable31Data } from '@/services/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Percent,
  CheckCircle2,
  AlertTriangle,
  Building2,
  FileText,
  Calendar,
  Save,
  RefreshCw,
} from 'lucide-react'

export const GSTPage: React.FC = () => {
  const queryClient = useQueryClient()
  const [activeTab, setActiveTab] = useState<'filing' | 'gstr3b' | 'settings'>('filing')

  // Return period selector
  const now = new Date()
  const [month, setMonth] = useState<string>((now.getMonth() + 1).toString())
  const [year, setYear] = useState<string>(now.getFullYear().toString())

  // Form State for GST Settings
  const [gstin, setGstin] = useState('')
  const [legalName, setLegalName] = useState('')
  const [tradeName, setTradeName] = useState('')
  const [regType, setRegType] = useState('REGULAR')
  const [defaultHsn, setDefaultHsn] = useState('6403')
  const [defaultRate, setDefaultRate] = useState(12)
  const [isComposition, setIsComposition] = useState(false)
  const [formSuccess, setFormSuccess] = useState<string | null>(null)
  const [formError, setFormError] = useState<string | null>(null)

  // 1. Fetch GST Settings
  const { data: settingsResp } = useQuery({
    queryKey: ['gst-settings'],
    queryFn: async () => {
      const res = await apiService.getGSTSettings()
      if (res.data) {
        setGstin(res.data.gstin)
        setLegalName(res.data.legal_name)
        setTradeName(res.data.trade_name || '')
        setRegType(res.data.registration_type || 'REGULAR')
        setDefaultHsn(res.data.default_footwear_hsn || '6403')
        setDefaultRate(Number(res.data.default_gst_rate || 12))
        setIsComposition(res.data.enable_composition || false)
      }
      return res
    },
  })

  // 2. Fetch GSTR-1
  const {
    data: gstr1Resp,
    isLoading: isLoadingGstr1,
    refetch: refetchGstr1,
  } = useQuery({
    queryKey: ['gst-gstr1', month, year],
    queryFn: () => apiService.getGSTR1({ month, year }),
  })

  // 3. Fetch GSTR-3B
  const {
    data: gstr3bResp,
    isLoading: isLoadingGstr3b,
    refetch: refetchGstr3b,
  } = useQuery({
    queryKey: ['gst-gstr3b', month, year],
    queryFn: () => apiService.getGSTR3B({ month, year }),
  })

  // Update Settings Mutation
  const updateMutation = useMutation({
    mutationFn: (payload: Partial<GSTSettingsData>) => apiService.updateGSTSettings(payload),
    onSuccess: (res) => {
      setFormSuccess('GST settings and registration profile updated successfully')
      setFormError(null)
      queryClient.setQueryData(['gst-settings'], res)
      queryClient.invalidateQueries({ queryKey: ['gst-gstr1'] })
      queryClient.invalidateQueries({ queryKey: ['gst-gstr3b'] })
    },
    onError: (err: Error) => {
      setFormError(err.message)
      setFormSuccess(null)
    },
  })

  const handleSaveSettings = (e: React.FormEvent) => {
    e.preventDefault()
    setFormSuccess(null)
    setFormError(null)
    updateMutation.mutate({
      gstin: gstin.toUpperCase(),
      legal_name: legalName,
      trade_name: tradeName,
      registration_type: regType,
      default_footwear_hsn: defaultHsn,
      default_gst_rate: defaultRate,
      enable_composition: isComposition,
    })
  }

  const currentSetting = settingsResp?.data
  const gstr1: GSTR1SummaryData | undefined = gstr1Resp?.data
  const gstr3b: GSTR3BTable31Data | undefined = gstr3bResp?.data

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight">GST & Tax Compliance</h1>
          <p className="text-sm text-slate-500 mt-1">
            Indian GST calculations, GSTR-1 outward summaries, HSN classification, and GSTR-3B liabilities.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => {
              refetchGstr1()
              refetchGstr3b()
            }}
          >
            <RefreshCw className="w-4 h-4 mr-1.5" />
            Refresh
          </Button>
        </div>
      </div>

      {/* Profile Banner */}
      <Card className="border-slate-200 bg-gradient-to-r from-slate-50 to-white">
        <CardContent className="p-4 flex flex-wrap items-center justify-between gap-4 text-xs">
          <div className="flex items-center gap-3">
            <div className="p-2.5 bg-purple-100 text-purple-700 rounded-lg">
              <Building2 className="w-5 h-5" />
            </div>
            <div>
              <span className="font-semibold text-slate-900 block text-sm">
                {currentSetting?.legal_name || 'Retail Enterprise'}
              </span>
              <span className="text-slate-500">
                GSTIN: <span className="font-mono font-medium text-slate-800">{currentSetting?.gstin || '-'}</span> •{' '}
                {currentSetting?.state_name} (Code {currentSetting?.state_code})
              </span>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Badge variant="outline" className="bg-white text-xs">
              {currentSetting?.registration_type || 'REGULAR'}
            </Badge>
            <Badge variant="secondary" className="text-xs">
              Footwear HSN: {currentSetting?.default_footwear_hsn || '6403'}
            </Badge>
          </div>
        </CardContent>
      </Card>

      {/* Navigation Tabs (Single-Selection Soft Lavender & Deep Purple) */}
      <div className="bg-[#F5F7FA] p-1.5 rounded-xl border border-[#E5EAF0] inline-flex items-center gap-1.5">
        <button
          type="button"
          onClick={() => setActiveTab('filing')}
          className={`flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition-all cursor-pointer whitespace-nowrap ${
            activeTab === 'filing'
              ? 'bg-[#EDE9FE] text-[#5B45B5] shadow-xs border border-[#8B5CF6]/30'
              : 'text-[#64748B] hover:text-[#5B45B5] hover:bg-[#F5F3FF] border border-transparent'
          }`}
        >
          <FileText className={`w-3.5 h-3.5 ${activeTab === 'filing' ? 'text-[#8B5CF6]' : 'text-[#64748B]'}`} />
          <span>GSTR-1 Monthly Return</span>
        </button>
        <button
          type="button"
          onClick={() => setActiveTab('gstr3b')}
          className={`flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition-all cursor-pointer whitespace-nowrap ${
            activeTab === 'gstr3b'
              ? 'bg-[#EDE9FE] text-[#5B45B5] shadow-xs border border-[#8B5CF6]/30'
              : 'text-[#64748B] hover:text-[#5B45B5] hover:bg-[#F5F3FF] border border-transparent'
          }`}
        >
          <Percent className={`w-3.5 h-3.5 ${activeTab === 'gstr3b' ? 'text-[#8B5CF6]' : 'text-[#64748B]'}`} />
          <span>GSTR-3B Tax Liability</span>
        </button>
        <button
          type="button"
          onClick={() => setActiveTab('settings')}
          className={`flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition-all cursor-pointer whitespace-nowrap ${
            activeTab === 'settings'
              ? 'bg-[#EDE9FE] text-[#5B45B5] shadow-xs border border-[#8B5CF6]/30'
              : 'text-[#64748B] hover:text-[#5B45B5] hover:bg-[#F5F3FF] border border-transparent'
          }`}
        >
          <Building2 className={`w-3.5 h-3.5 ${activeTab === 'settings' ? 'text-[#8B5CF6]' : 'text-[#64748B]'}`} />
          <span>GST Settings & Rates</span>
        </button>
      </div>

      {/* TAB 1: GSTR-1 MONTHLY RETURN */}
      {activeTab === 'filing' && (
        <div className="space-y-6">
          {/* Month / Year Selector */}
          <Card className="border-slate-200">
            <CardContent className="p-4 flex flex-wrap items-center justify-between gap-3 text-xs">
              <div className="flex items-center gap-2">
                <Calendar className="w-4 h-4 text-slate-400" />
                <span className="font-semibold text-slate-700">Filing Period:</span>
                <select
                  value={month}
                  onChange={(e) => setMonth(e.target.value)}
                  className="px-2.5 py-1.5 border border-slate-200 rounded-lg bg-white text-xs focus:ring-2 focus:ring-purple-500"
                >
                  <option value="1">January</option>
                  <option value="2">February</option>
                  <option value="3">March</option>
                  <option value="4">April</option>
                  <option value="5">May</option>
                  <option value="6">June</option>
                  <option value="7">July</option>
                  <option value="8">August</option>
                  <option value="9">September</option>
                  <option value="10">October</option>
                  <option value="11">November</option>
                  <option value="12">December</option>
                </select>

                <select
                  value={year}
                  onChange={(e) => setYear(e.target.value)}
                  className="px-2.5 py-1.5 border border-slate-200 rounded-lg bg-white text-xs focus:ring-2 focus:ring-purple-500"
                >
                  <option value="2025">2025</option>
                  <option value="2026">2026</option>
                  <option value="2027">2027</option>
                </select>
              </div>

              <div className="text-slate-500">
                Return Period: <span className="font-mono font-semibold text-slate-800">{gstr1?.return_period}</span>
              </div>
            </CardContent>
          </Card>

          {/* GSTR-1 Summary KPI Cards */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            <Card className="border-slate-200">
              <CardContent className="p-5">
                <span className="text-xs font-medium text-slate-500 uppercase tracking-wider">Total Taxable Value</span>
                <h3 className="text-2xl font-bold text-slate-900 font-mono mt-1">
                  ₹{Number(gstr1?.total_taxable_value || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                </h3>
                <p className="text-xs text-slate-400 mt-1">{gstr1?.b2c_small_count || 0} Invoices (B2C Small)</p>
              </CardContent>
            </Card>

            <Card className="border-slate-200">
              <CardContent className="p-5">
                <span className="text-xs font-medium text-slate-500 uppercase tracking-wider">Central Tax (CGST)</span>
                <h3 className="text-2xl font-bold text-[#7C3AED] font-mono mt-1">
                  ₹{Number(gstr1?.total_cgst || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                </h3>
                <p className="text-xs text-slate-400 mt-1">50% Intra-State Rate Share</p>
              </CardContent>
            </Card>

            <Card className="border-slate-200">
              <CardContent className="p-5">
                <span className="text-xs font-medium text-slate-500 uppercase tracking-wider">State Tax (SGST)</span>
                <h3 className="text-2xl font-bold text-[#8B5CF6] font-mono mt-1">
                  ₹{Number(gstr1?.total_sgst || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                </h3>
                <p className="text-xs text-slate-400 mt-1">50% Intra-State Rate Share</p>
              </CardContent>
            </Card>

            <Card className="border-slate-200">
              <CardContent className="p-5">
                <span className="text-xs font-medium text-slate-500 uppercase tracking-wider">Gross Tax Liability</span>
                <h3 className="text-2xl font-bold text-purple-700 font-mono mt-1">
                  ₹{Number(gstr1?.total_tax_liability || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                </h3>
                <p className="text-xs text-slate-400 mt-1">CGST + SGST + IGST</p>
              </CardContent>
            </Card>
          </div>

          {/* Rate Slab Matrix */}
          <Card className="border-slate-200 overflow-hidden">
            <CardHeader className="py-4 px-6 border-b border-slate-100">
              <CardTitle className="text-base font-semibold text-slate-800">
                Tax Rate Slab Breakdown (Table 7: B2C Small)
              </CardTitle>
              <CardDescription className="text-xs">
                Tax liability calculated per rate slab. Footwear typically categorized under 5% or 12%.
              </CardDescription>
            </CardHeader>
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm text-slate-600">
                <thead className="bg-slate-50 text-xs uppercase text-slate-500 border-b border-slate-200 font-mono">
                  <tr>
                    <th className="px-6 py-3.5 font-medium">GST Rate Slab</th>
                    <th className="px-6 py-3.5 font-medium text-right">Taxable Value (₹)</th>
                    <th className="px-6 py-3.5 font-medium text-right">CGST (₹)</th>
                    <th className="px-6 py-3.5 font-medium text-right">SGST (₹)</th>
                    <th className="px-6 py-3.5 font-medium text-right">IGST (₹)</th>
                    <th className="px-6 py-3.5 font-medium text-right">Total Tax (₹)</th>
                    <th className="px-6 py-3.5 font-medium text-center">Invoices</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100 font-mono">
                  {isLoadingGstr1 ? (
                    <tr>
                      <td colSpan={7} className="px-6 py-8 text-center text-slate-400">Loading tax slabs...</td>
                    </tr>
                  ) : (gstr1?.rate_breakdown || []).length === 0 ? (
                    <tr>
                      <td colSpan={7} className="px-6 py-8 text-center text-slate-400">No taxable supplies recorded for this period.</td>
                    </tr>
                  ) : (
                    gstr1?.rate_breakdown.map((slab, idx) => (
                      <tr key={idx} className="hover:bg-slate-50/70">
                        <td className="px-6 py-4 font-sans font-semibold text-slate-900">
                          <Badge variant="outline" className="text-xs">
                            {Number(slab.rate_pct).toFixed(0)}% GST Slab
                          </Badge>
                        </td>
                        <td className="px-6 py-4 text-right font-medium text-slate-900">
                          ₹{Number(slab.taxable_value).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                        </td>
                        <td className="px-6 py-4 text-right text-[#7C3AED]">
                          ₹{Number(slab.cgst).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                        </td>
                        <td className="px-6 py-4 text-right text-[#8B5CF6]">
                          ₹{Number(slab.sgst).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                        </td>
                        <td className="px-6 py-4 text-right text-slate-400">₹0.00</td>
                        <td className="px-6 py-4 text-right font-bold text-purple-700">
                          ₹{Number(slab.total_tax).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                        </td>
                        <td className="px-6 py-4 text-center font-sans text-slate-700">{slab.invoices}</td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </Card>

          {/* Table 12: HSN Summary */}
          <Card className="border-slate-200 overflow-hidden">
            <CardHeader className="py-4 px-6 border-b border-slate-100">
              <CardTitle className="text-base font-semibold text-slate-800">
                Table 12: HSN-wise Summary of Outward Supplies
              </CardTitle>
              <CardDescription className="text-xs">
                Mandatory for GSTR-1 e-filing for retail footwear and goods.
              </CardDescription>
            </CardHeader>
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm text-slate-600">
                <thead className="bg-slate-50 text-xs uppercase text-slate-500 border-b border-slate-200 font-mono">
                  <tr>
                    <th className="px-6 py-3.5 font-medium">HSN Code</th>
                    <th className="px-6 py-3.5 font-medium font-sans">Description</th>
                    <th className="px-6 py-3.5 font-medium text-center">UQC</th>
                    <th className="px-6 py-3.5 font-medium text-center">Total Qty</th>
                    <th className="px-6 py-3.5 font-medium text-right">Taxable Value (₹)</th>
                    <th className="px-6 py-3.5 font-medium text-right">Central Tax (₹)</th>
                    <th className="px-6 py-3.5 font-medium text-right">State Tax (₹)</th>
                    <th className="px-6 py-3.5 font-medium text-right">Total Tax (₹)</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100 font-mono">
                  {(gstr1?.hsn_summary || []).length === 0 ? (
                    <tr>
                      <td colSpan={8} className="px-6 py-8 text-center text-slate-400">No HSN entries recorded.</td>
                    </tr>
                  ) : (
                    gstr1?.hsn_summary.map((hsn, idx) => (
                      <tr key={idx} className="hover:bg-slate-50/70">
                        <td className="px-6 py-4 font-bold text-slate-900">{hsn.hsn_code}</td>
                        <td className="px-6 py-4 font-sans text-xs text-slate-600 max-w-xs">{hsn.description}</td>
                        <td className="px-6 py-4 text-center">
                          <Badge variant="secondary" className="text-[11px]">{hsn.uqc}</Badge>
                        </td>
                        <td className="px-6 py-4 text-center font-bold text-slate-900">{hsn.total_qty}</td>
                        <td className="px-6 py-4 text-right">
                          ₹{Number(hsn.taxable_value).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                        </td>
                        <td className="px-6 py-4 text-right text-[#7C3AED]">
                          ₹{Number(hsn.cgst).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                        </td>
                        <td className="px-6 py-4 text-right text-[#8B5CF6]">
                          ₹{Number(hsn.sgst).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                        </td>
                        <td className="px-6 py-4 text-right font-bold text-purple-700">
                          ₹{Number(hsn.total_tax).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
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

      {/* TAB 2: GSTR-3B TABLE 3.1 */}
      {activeTab === 'gstr3b' && (
        <Card className="border-slate-200 overflow-hidden">
          <CardHeader className="py-4 px-6 border-b border-slate-100">
            <CardTitle className="text-base font-semibold text-slate-800">
              GSTR-3B — Table 3.1: Details of Outward Supplies
            </CardTitle>
            <CardDescription className="text-xs">
              Summary values for filing monthly return GSTR-3B on the GST portal.
            </CardDescription>
          </CardHeader>
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm text-slate-600">
              <thead className="bg-slate-50 text-xs uppercase text-slate-500 border-b border-slate-200">
                <tr>
                  <th className="px-6 py-3.5 font-medium">Nature of Supplies</th>
                  <th className="px-6 py-3.5 font-medium text-right font-mono">Total Taxable Value (₹)</th>
                  <th className="px-6 py-3.5 font-medium text-right font-mono">Integrated Tax (₹)</th>
                  <th className="px-6 py-3.5 font-medium text-right font-mono">Central Tax (₹)</th>
                  <th className="px-6 py-3.5 font-medium text-right font-mono">State/UT Tax (₹)</th>
                  <th className="px-6 py-3.5 font-medium text-right font-mono">Cess (₹)</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100 font-mono">
                {isLoadingGstr3b ? (
                  <tr>
                    <td colSpan={6} className="px-6 py-8 text-center text-slate-400">Loading GSTR-3B...</td>
                  </tr>
                ) : (
                  <tr className="hover:bg-slate-50/70">
                    <td className="px-6 py-4 font-sans font-medium text-slate-900 text-xs">
                      {gstr3b?.nature_of_supplies}
                    </td>
                    <td className="px-6 py-4 text-right font-bold text-slate-900">
                      ₹{Number(gstr3b?.taxable_value || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                    </td>
                    <td className="px-6 py-4 text-right text-slate-400">₹0.00</td>
                    <td className="px-6 py-4 text-right font-semibold text-[#7C3AED]">
                      ₹{Number(gstr3b?.central_tax || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                    </td>
                    <td className="px-6 py-4 text-right font-semibold text-[#8B5CF6]">
                      ₹{Number(gstr3b?.state_ut_tax || 0).toLocaleString('en-IN', { minimumFractionDigits: 2 })}
                    </td>
                    <td className="px-6 py-4 text-right text-slate-400">₹0.00</td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </Card>
      )}

      {/* TAB 3: GST SETTINGS & CONFIGURATION */}
      {activeTab === 'settings' && (
        <Card className="border-slate-200">
          <CardHeader className="py-4 px-6 border-b border-slate-100">
            <CardTitle className="text-base font-semibold text-slate-800">
              Tenant GST Profile & Tax Parameters
            </CardTitle>
            <CardDescription className="text-xs">
              State code is automatically derived from the first two digits of your Indian GSTIN.
            </CardDescription>
          </CardHeader>
          <CardContent className="p-6">
            <form onSubmit={handleSaveSettings} className="space-y-4 max-w-2xl">
              {formSuccess && (
                <div className="p-3 bg-emerald-50 border border-emerald-200 text-emerald-800 rounded-lg text-xs flex items-center gap-2">
                  <CheckCircle2 className="w-4 h-4 text-emerald-600 flex-shrink-0" />
                  <span>{formSuccess}</span>
                </div>
              )}

              {formError && (
                <div className="p-3 bg-rose-50 border border-rose-200 text-rose-800 rounded-lg text-xs flex items-center gap-2">
                  <AlertTriangle className="w-4 h-4 text-rose-600 flex-shrink-0" />
                  <span>{formError}</span>
                </div>
              )}

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="text-xs font-semibold text-slate-700 block mb-1">
                    GSTIN (15 Characters) *
                  </label>
                  <input
                    type="text"
                    maxLength={15}
                    value={gstin}
                    onChange={(e) => setGstin(e.target.value.toUpperCase())}
                    placeholder="e.g. 33AAAAA0000A1Z5"
                    className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm font-mono uppercase focus:ring-2 focus:ring-purple-500 focus:outline-none"
                    required
                  />
                  <span className="text-[11px] text-slate-400 mt-0.5 block">Format: 2 state digits + 10 PAN + 1 entity + Z + 1 check digit</span>
                </div>

                <div>
                  <label className="text-xs font-semibold text-slate-700 block mb-1">
                    Registration Type
                  </label>
                  <select
                    value={regType}
                    onChange={(e) => setRegType(e.target.value)}
                    className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm bg-white focus:ring-2 focus:ring-purple-500 focus:outline-none"
                  >
                    <option value="REGULAR">Regular Taxpayer</option>
                    <option value="COMPOSITION">Composition Scheme</option>
                    <option value="UNREGISTERED">Unregistered Retailer</option>
                  </select>
                </div>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="text-xs font-semibold text-slate-700 block mb-1">Legal Registered Name *</label>
                  <input
                    type="text"
                    value={legalName}
                    onChange={(e) => setLegalName(e.target.value)}
                    placeholder="Legal business name"
                    className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:ring-2 focus:ring-purple-500 focus:outline-none"
                    required
                  />
                </div>

                <div>
                  <label className="text-xs font-semibold text-slate-700 block mb-1">Trade Name (Brand Name)</label>
                  <input
                    type="text"
                    value={tradeName}
                    onChange={(e) => setTradeName(e.target.value)}
                    placeholder="Trade name"
                    className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:ring-2 focus:ring-purple-500 focus:outline-none"
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="text-xs font-semibold text-slate-700 block mb-1">Default Footwear HSN</label>
                  <select
                    value={defaultHsn}
                    onChange={(e) => setDefaultHsn(e.target.value)}
                    className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm bg-white focus:ring-2 focus:ring-purple-500 focus:outline-none"
                  >
                    <option value="6403">6403 — Leather Sole & Upper</option>
                    <option value="6404">6404 — Rubber / Plastic Sole & Textile Upper</option>
                    <option value="6402">6402 — Rubber / Plastic Sole & Upper</option>
                    <option value="6401">6401 — Waterproof Footwear</option>
                    <option value="6405">6405 — Other Footwear</option>
                  </select>
                </div>

                <div>
                  <label className="text-xs font-semibold text-slate-700 block mb-1">Default Footwear GST Rate (%)</label>
                  <select
                    value={defaultRate}
                    onChange={(e) => setDefaultRate(Number(e.target.value))}
                    className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm bg-white focus:ring-2 focus:ring-purple-500 focus:outline-none"
                  >
                    <option value={12}>12.00% (Footwear standard rate)</option>
                    <option value={5}>5.00% (Footwear sale value ≤ ₹1000)</option>
                    <option value={18}>18.00% (Accessories & synthetic items)</option>
                  </select>
                </div>
              </div>

              <div className="pt-2">
                <Button
                  type="submit"
                  size="sm"
                  className="bg-[#8B5CF6] hover:bg-[#7C3AED] text-white shadow-xs"
                  disabled={updateMutation.isPending}
                >
                  <Save className="w-4 h-4 mr-1.5" />
                  {updateMutation.isPending ? 'Saving...' : 'Save GST Configuration'}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
