import React, { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { apiService } from '@/services/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { KeyRound, ArrowLeft, CheckCircle2, AlertCircle } from 'lucide-react'

export const ForgotPasswordPage: React.FC = () => {
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [submitted, setSubmitted] = useState(false)
  const [tokenHint, setTokenHint] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      const res = await apiService.forgotPassword(email)
      setSubmitted(true)
      if (res.data?.reset_token) {
        setTokenHint(res.data.reset_token)
      }
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Request failed')
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
            Reset Password
          </h1>
          <p className="mt-1 text-xs text-slate-500">
            Generate a secure, single-use reset token for your account
          </p>
        </div>

        <Card className="rounded-2xl border-slate-200/90 shadow-md">
          <CardHeader>
            <CardTitle className="text-lg flex items-center gap-2">
              <KeyRound className="w-5 h-5 text-blue-600" />
              Reset Your Password
            </CardTitle>
            <CardDescription>
              Enter your email address and we'll generate an expiring single-use reset token.
            </CardDescription>
          </CardHeader>

          {submitted ? (
            <CardContent className="space-y-4">
              <div className="p-4 bg-emerald-50 border border-emerald-200 rounded-lg text-xs text-emerald-800 space-y-2">
                <div className="flex items-center gap-2 font-semibold">
                  <CheckCircle2 className="w-4 h-4 text-emerald-600" />
                  Reset Token Generated
                </div>
                <p>If that email exists in KadaiIQ, a secure reset token has been issued.</p>
                {tokenHint && (
                  <div className="pt-2">
                    <span className="block text-[10px] text-slate-500 font-mono">Demo token:</span>
                    <code className="bg-emerald-100 px-2 py-1 rounded block text-[11px] break-all select-all">
                      {tokenHint}
                    </code>
                  </div>
                )}
              </div>
              <Button
                onClick={() => navigate(`/reset-password${tokenHint ? `?token=${tokenHint}` : ''}`)}
                className="w-full"
              >
                Proceed to Reset Password
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
                  <label className="text-xs font-semibold text-slate-700">Account Email</label>
                  <input
                    type="email"
                    required
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="name@business.com"
                    className="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-600"
                  />
                </div>
              </CardContent>
              <CardFooter className="flex flex-col space-y-3">
                <Button type="submit" disabled={loading} className="w-full">
                  {loading ? 'Generating...' : 'Send Reset Link'}
                </Button>
                <Link to="/login" className="text-xs text-slate-500 flex items-center gap-1 hover:underline">
                  <ArrowLeft className="w-3.5 h-3.5" /> Back to Login
                </Link>
              </CardFooter>
            </form>
          )}
        </Card>
      </div>
    </div>
  )
}

