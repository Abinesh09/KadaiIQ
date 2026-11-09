import React, { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiService } from '@/services/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ShieldCheck, Plus, Building2, Users, AlertCircle } from 'lucide-react'

export const SuperAdminPage: React.FC = () => {
  const queryClient = useQueryClient()
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [form, setForm] = useState({
    name: '',
    slug: '',
    business_name: '',
    business_type: 'Footwear & Retail',
    gstin: '',
    email: '',
    phone: '',
    owner_name: '',
    owner_email: '',
    owner_password: '',
  })
  const [error, setError] = useState<string | null>(null)

  const { data: tenantsResp, isLoading } = useQuery({
    queryKey: ['super-admin-tenants'],
    queryFn: apiService.getTenants,
  })

  const tenants = tenantsResp?.data || []

  const createMutation = useMutation({
    mutationFn: apiService.createTenant,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['super-admin-tenants'] })
      setShowCreateModal(false)
      setForm({
        name: '',
        slug: '',
        business_name: '',
        business_type: 'Footwear & Retail',
        gstin: '',
        email: '',
        phone: '',
        owner_name: '',
        owner_email: '',
        owner_password: '',
      })
    },
    onError: (err: Error) => {
      setError(err.message)
    },
  })

  const statusMutation = useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) =>
      apiService.updateTenantStatus(id, status),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['super-admin-tenants'] })
    },
  })

  const handleNameChange = (name: string) => {
    const slug = name
      .toLowerCase()
      .trim()
      .replace(/[^\w\s-]/g, '')
      .replace(/[\s_-]+/g, '-')
      .replace(/^-+|-+$/g, '')
    setForm((prev) => ({ ...prev, name, business_name: name, slug }))
  }

  const handleCreateSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    createMutation.mutate(form)
  }

  return (
    <div className="space-y-6">
      {/* Top Banner */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200/90 shadow-xs">
        <div>
          <h1 className="text-xl sm:text-2xl font-extrabold text-slate-900 tracking-tight flex items-center gap-2.5">
            <div className="w-9 h-9 rounded-xl bg-[#F5F3FF] text-[#7C3AED] flex items-center justify-center">
              <ShieldCheck className="w-5 h-5" />
            </div>
            Super Admin Console
          </h1>
          <p className="text-xs sm:text-sm text-slate-500 mt-1">
            Global organisation provisioning, tenant management, and platform telemetry
          </p>
        </div>

        <Button
          onClick={() => setShowCreateModal(true)}
          className="bg-[#8B5CF6] hover:bg-[#7C3AED] text-white font-semibold rounded-xl shadow-sm cursor-pointer"
        >
          <Plus className="w-4 h-4 mr-2" />
          Provision New Organisation
        </Button>
      </div>

      {/* Super Admin Metric Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <Card className="rounded-2xl border-slate-200/90 shadow-2xs">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <span className="text-xs font-bold text-slate-400 uppercase tracking-widest">
              Total Organisations
            </span>
            <Building2 className="w-4 h-4 text-[#7C3AED]" />
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-black text-slate-900 font-numeric">{tenants.length}</div>
            <p className="text-[11px] text-slate-500 mt-1">Provisioned retail workspaces</p>
          </CardContent>
        </Card>

        <Card className="rounded-2xl border-slate-200/90 shadow-2xs">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <span className="text-xs font-bold text-slate-400 uppercase tracking-widest">
              Active Tenants
            </span>
            <div className="w-2 h-2 rounded-full bg-emerald-500 animate-ping" />
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-black text-emerald-600 font-numeric">
              {tenants.filter((t) => t.status === 'active').length}
            </div>
            <p className="text-[11px] text-slate-500 mt-1">Operating normally</p>
          </CardContent>
        </Card>

        <Card className="rounded-2xl border-slate-200/90 shadow-2xs">
          <CardHeader className="flex flex-row items-center justify-between pb-2 space-y-0">
            <span className="text-xs font-bold text-slate-400 uppercase tracking-widest">
              Platform Architecture
            </span>
            <ShieldCheck className="w-4 h-4 text-[#7C3AED]" />
          </CardHeader>
          <CardContent>
            <div className="text-sm font-bold text-slate-900 mt-1">Shared Schema Isolation</div>
            <p className="text-[11px] text-slate-500 mt-1">PostgreSQL row-level tenant keys</p>
          </CardContent>
        </Card>
      </div>

      {/* Tenants Table */}
      <Card className="rounded-2xl border-slate-200/90 shadow-xs">
        <CardHeader>
          <CardTitle className="text-base flex items-center justify-between">
            <span>Managed Retail Tenants ({tenants.length})</span>
            <span className="text-xs font-normal text-slate-500">PostgreSQL Shared Schema with Isolation</span>
          </CardTitle>
          <CardDescription>All retail tenants configured on the KadaiIQ platform</CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="py-8 text-center text-sm text-slate-500">Loading tenants...</div>
          ) : tenants.length === 0 ? (
            <div className="py-8 text-center text-sm text-slate-500">
              No retail organisations created yet. Click "Provision New Organisation" to create one.
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm text-slate-600">
                <thead className="bg-slate-50 text-xs uppercase text-slate-400 font-semibold border-b border-slate-200">
                  <tr>
                    <th className="px-4 py-3">Organisation</th>
                    <th className="px-4 py-3">Slug</th>
                    <th className="px-4 py-3">Contact</th>
                    <th className="px-4 py-3">GSTIN</th>
                    <th className="px-4 py-3">Status</th>
                    <th className="px-4 py-3">Members</th>
                    <th className="px-4 py-3 text-right">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {tenants.map((t) => (
                    <tr key={t.id} className="hover:bg-slate-50">
                      <td className="px-4 py-3">
                        <div className="font-semibold text-slate-900 flex items-center gap-2">
                          <Building2 className="w-4 h-4 text-[#7C3AED]" />
                          {t.name}
                        </div>
                        <div className="text-xs text-slate-400">{t.business_type}</div>
                      </td>
                      <td className="px-4 py-3">
                        <code className="text-xs bg-slate-100 px-2 py-0.5 rounded font-mono text-slate-700">
                          {t.slug}
                        </code>
                      </td>
                      <td className="px-4 py-3">
                        <div>{t.email}</div>
                        <div className="text-xs text-slate-400">{t.phone || 'No phone'}</div>
                      </td>
                      <td className="px-4 py-3 font-mono text-xs">
                        {t.gstin || <span className="text-slate-400">Not configured</span>}
                      </td>
                      <td className="px-4 py-3">
                        <Badge variant={t.status === 'active' ? 'success' : 'destructive'}>
                          {t.status}
                        </Badge>
                      </td>
                      <td className="px-4 py-3">
                        <span className="inline-flex items-center gap-1 text-xs text-slate-700">
                          <Users className="w-3.5 h-3.5" />
                          {t.member_count}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-right space-x-2">
                        {t.status === 'active' ? (
                          <Button
                            variant="outline"
                            size="sm"
                            className="text-red-600 hover:text-red-700 text-xs"
                            onClick={() => statusMutation.mutate({ id: t.id, status: 'suspended' })}
                          >
                            Suspend
                          </Button>
                        ) : (
                          <Button
                            variant="outline"
                            size="sm"
                            className="text-emerald-600 hover:text-emerald-700 text-xs"
                            onClick={() => statusMutation.mutate({ id: t.id, status: 'active' })}
                          >
                            Activate
                          </Button>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Provisioning Modal */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-slate-900/50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-xl shadow-xl max-w-lg w-full p-6 space-y-4">
            <div className="flex items-center justify-between border-b pb-3">
              <h2 className="text-lg font-bold text-slate-900">Provision Retail Organisation</h2>
              <button
                onClick={() => setShowCreateModal(false)}
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

            <form onSubmit={handleCreateSubmit} className="space-y-4">
              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Store Name</label>
                  <input
                    type="text"
                    required
                    value={form.name}
                    onChange={(e) => handleNameChange(e.target.value)}
                    placeholder="e.g. Royal Footwear"
                    className="w-full px-3 py-1.5 border rounded-lg text-sm"
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Tenant Slug</label>
                  <input
                    type="text"
                    required
                    value={form.slug}
                    onChange={(e) => setForm({ ...form, slug: e.target.value })}
                    placeholder="royal-footwear"
                    className="w-full px-3 py-1.5 border rounded-lg text-sm"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3">
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Business Type</label>
                  <input
                    type="text"
                    value={form.business_type}
                    onChange={(e) => setForm({ ...form, business_type: e.target.value })}
                    className="w-full px-3 py-1.5 border rounded-lg text-sm"
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">GSTIN (Optional)</label>
                  <input
                    type="text"
                    value={form.gstin}
                    onChange={(e) => setForm({ ...form, gstin: e.target.value })}
                    placeholder="33AAAAA0000A1Z5"
                    className="w-full px-3 py-1.5 border rounded-lg text-sm"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-3 border-t pt-3">
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Owner Name</label>
                  <input
                    type="text"
                    required
                    value={form.owner_name}
                    onChange={(e) => setForm({ ...form, owner_name: e.target.value })}
                    placeholder="Owner Full Name"
                    className="w-full px-3 py-1.5 border rounded-lg text-sm"
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Owner Email</label>
                  <input
                    type="email"
                    required
                    value={form.owner_email}
                    onChange={(e) => setForm({ ...form, owner_email: e.target.value, email: e.target.value })}
                    placeholder="owner@store.com"
                    className="w-full px-3 py-1.5 border rounded-lg text-sm"
                  />
                </div>
              </div>

              <div className="space-y-1">
                <label className="text-xs font-semibold text-slate-700">Owner Password</label>
                <input
                  type="password"
                  required
                  value={form.owner_password}
                  onChange={(e) => setForm({ ...form, owner_password: e.target.value })}
                  placeholder="Temporary initial password"
                  className="w-full px-3 py-1.5 border rounded-lg text-sm"
                />
              </div>

              <div className="flex justify-end gap-3 pt-4 border-t">
                <Button type="button" variant="outline" onClick={() => setShowCreateModal(false)}>
                  Cancel
                </Button>
                <Button type="submit" disabled={createMutation.isPending} className="bg-[#8B5CF6] hover:bg-[#7C3AED]">
                  {createMutation.isPending ? 'Provisioning...' : 'Provision Tenant'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}
