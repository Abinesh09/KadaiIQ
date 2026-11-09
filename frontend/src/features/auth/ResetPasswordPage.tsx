import React, { useState, useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { apiService } from '@/services/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Lock, CheckCircle2, AlertCircle } from 'lucide-react'

export const ResetPasswordPage: React.FC = () => {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const [token, setToken] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [success, setSuccess] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    const urlToken = searchParams.get('token')
    if (urlToken) {
      setToken(urlToken)
    }
  }, [searchParams])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    if (newPassword !== confirmPassword) {
      setError('Passwords do not match')
      return
    }
    if (newPassword.length < 8) {
      setError('Password must be at least 8 characters')
      return
    }

    setLoading(true)
    try {
      await apiService.resetPassword({ token, new_password: newPassword })
      setSuccess(true)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Reset failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-slate-50 flex flex-col justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-md space-y-6">
        <div className="flex flex-col items-center text-center">
          <div className="w-12 h-12 rounded-2xl bg-gradient-to-tr from-slate-900 via-blue-900 to-indigo-700 flex items-center justify-center text-white shadow-md ring-1 ring-blue-500/20 mb-3">
            <span className="font-mono text-emerald-400 font-extrabold text-xl mr-0.5">K</span>
            <span className="text-white text-xs font-bold">IQ</span>
          </div>
          <h1 className="text-2xl font-extrabold text-slate-900 tracking-tight">
            New Password
          </h1>
          <p className="mt-1 text-xs text-slate-500">
            Set your new secure password to access your store workspace
          </p>
        </div>

        <Card className="rounded-2xl border-slate-200/90 shadow-md">
          <CardHeader>
            <CardTitle className="text-lg flex items-center gap-2">
              <Lock className="w-5 h-5 text-[#7C3AED]" />
              Set New Password
            </CardTitle>
            <CardDescription>Enter the token received and choose your new password</CardDescription>
          </CardHeader>

          {success ? (
            <CardContent className="space-y-4">
              <div className="p-4 bg-emerald-50 border border-emerald-200 rounded-lg text-xs text-emerald-800 space-y-2">
                <div className="flex items-center gap-2 font-semibold">
                  <CheckCircle2 className="w-4 h-4 text-emerald-600" />
                  Password Updated Successfully
                </div>
                <p>All previous sessions have been invalidated. You can now log in with your new password.</p>
              </div>
              <Button onClick={() => navigate('/login')} className="w-full">
                Proceed to Sign In
              </Button>
            </CardContent>
          ) : (
            <form onSubmit={handleSubmit}>
              <CardContent className="space-y-4">
                {error && (
                  <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-xs text-red-700 flex items-center gap-2">
                    <AlertCircle className="w-4 h-4" />
                    <span>{error}</span>
                  </div>
                )}

                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Reset Token</label>
                  <input
                    type="text"
                    required
                    value={token}
                    onChange={(e) => setToken(e.target.value)}
                    placeholder="Enter reset token"
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm font-mono focus:outline-none focus:ring-2 focus:ring-[#8B5CF6]"
                  />
                </div>

                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">New Password</label>
                  <input
                    type="password"
                    required
                    value={newPassword}
                    onChange={(e) => setNewPassword(e.target.value)}
                    placeholder="••••••••"
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[#8B5CF6]"
                  />
                </div>

                <div className="space-y-1">
                  <label className="text-xs font-semibold text-slate-700">Confirm New Password</label>
                  <input
                    type="password"
                    required
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    placeholder="••••••••"
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[#8B5CF6]"
                  />
                </div>
              </CardContent>
              <CardFooter>
                <Button type="submit" disabled={loading} className="w-full">
                  {loading ? 'Updating Password...' : 'Save New Password'}
                </Button>
              </CardFooter>
            </form>
          )}
        </Card>
      </div>
    </div>
  )
}
