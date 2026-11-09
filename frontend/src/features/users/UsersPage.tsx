import React, { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiService } from '@/services/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Users, UserPlus, Shield, AlertCircle } from 'lucide-react'

export const UsersPage: React.FC = () => {
  const queryClient = useQueryClient()
  const [showInviteModal, setShowInviteModal] = useState(false)
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [roleName, setRoleName] = useState('Staff')
  const [error, setError] = useState<string | null>(null)

  const { data: membersResp, isLoading } = useQuery({
    queryKey: ['tenant-members'],
    queryFn: apiService.getMembers,
  })

  const members = membersResp?.data || []

  const inviteMutation = useMutation({
    mutationFn: (data: { name: string; email: string; role_id: string }) =>
      apiService.inviteMember(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tenant-members'] })
      setShowInviteModal(false)
      setName('')
      setEmail('')
    },
    onError: (err: Error) => {
      setError(err.message)
    },
  })

  const handleInvite = (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    // Find role ID from members if available, or generate standard lookup
    const existingMemberWithRole = members.find((m) => m.role_name === roleName)
    const roleId = existingMemberWithRole?.role_id || '00000000-0000-0000-0000-000000000000'

    inviteMutation.mutate({
      name,
      email,
      role_id: roleId,
    })
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 tracking-tight flex items-center gap-2">
            <Users className="w-6 h-6 text-blue-600" />
            Organisation Staff & Members
          </h1>
          <p className="text-sm text-slate-500 mt-1">
            Manage store roles, permissions, cashier accounts, and store managers
          </p>
        </div>

        <Button onClick={() => setShowInviteModal(true)}>
          <UserPlus className="w-4 h-4 mr-2" />
          Invite Staff Member
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Active Members ({members.length})</CardTitle>
          <CardDescription>Members who have access to this store workspace</CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="py-8 text-center text-sm text-slate-500">Loading members...</div>
          ) : members.length === 0 ? (
            <div className="py-8 text-center text-sm text-slate-500">
              No staff members found.
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm text-slate-600">
                <thead className="bg-slate-50 text-xs uppercase text-slate-400 font-semibold border-b">
                  <tr>
                    <th className="px-4 py-3">Member</th>
                    <th className="px-4 py-3">Email</th>
                    <th className="px-4 py-3">Assigned Role</th>
                    <th className="px-4 py-3">Status</th>
                    <th className="px-4 py-3">Joined Date</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {members.map((m) => (
                    <tr key={m.id} className="hover:bg-slate-50">
                      <td className="px-4 py-3 font-semibold text-slate-900">{m.name}</td>
                      <td className="px-4 py-3">{m.email}</td>
                      <td className="px-4 py-3">
                        <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium bg-blue-50 text-blue-700 border border-blue-200">
                          <Shield className="w-3 h-3 text-blue-600" />
                          {m.role_name}
                        </span>
                      </td>
                      <td className="px-4 py-3">
                        <Badge variant="success">{m.status}</Badge>
                      </td>
                      <td className="px-4 py-3 text-xs text-slate-400">
                        {new Date(m.created_at).toLocaleDateString()}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Invite Modal */}
      {showInviteModal && (
        <div className="fixed inset-0 bg-slate-900/50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-xl shadow-xl max-w-md w-full p-6 space-y-4">
            <div className="flex items-center justify-between border-b pb-3">
              <h2 className="text-lg font-bold text-slate-900">Invite Staff Member</h2>
              <button
                onClick={() => setShowInviteModal(false)}
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

            <form onSubmit={handleInvite} className="space-y-4">
              <div className="space-y-1">
                <label className="text-xs font-semibold text-slate-700">Full Name</label>
                <input
                  type="text"
                  required
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="e.g. Priya Sharma"
                  className="w-full px-3 py-2 border rounded-lg text-sm"
                />
              </div>

              <div className="space-y-1">
                <label className="text-xs font-semibold text-slate-700">Work Email</label>
                <input
                  type="email"
                  required
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="priya@store.com"
                  className="w-full px-3 py-2 border rounded-lg text-sm"
                />
              </div>

              <div className="space-y-1">
                <label className="text-xs font-semibold text-slate-700">Role</label>
                <select
                  value={roleName}
                  onChange={(e) => setRoleName(e.target.value)}
                  className="w-full px-3 py-2 border rounded-lg text-sm bg-white"
                >
                  <option value="Staff">Staff (Cashier / Floor)</option>
                  <option value="Manager">Manager (Inventory / Sales)</option>
                  <option value="Tenant Admin">Tenant Admin (Store Admin)</option>
                </select>
              </div>

              <div className="flex justify-end gap-3 pt-4 border-t">
                <Button type="button" variant="outline" onClick={() => setShowInviteModal(false)}>
                  Cancel
                </Button>
                <Button type="submit" disabled={inviteMutation.isPending}>
                  {inviteMutation.isPending ? 'Inviting...' : 'Send Invitation'}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  )
}

