import React, { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from './AuthContext'
import { Button } from '@/components/ui/button'
import {
  ShieldCheck,
  ArrowRight,
  AlertCircle,
  TrendingUp,
  Boxes,
  Percent,
} from 'lucide-react'

export const LoginPage: React.FC = () => {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      await login({ email, password })
      // Navigate all users to Executive Dashboard
      navigate('/')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Invalid credentials')
    } finally {
      setLoading(false)
    }
  }

  const fillSuperAdmin = () => {
    setEmail('superadmin@kadaiiq.com')
    setPassword('Admin@KadaiIQ2026')
  }

  return (
    <div className="min-h-screen bg-slate-50 flex flex-col justify-center py-12 px-4 sm:px-6 lg:px-8 selection:bg-blue-100 selection:text-blue-900">
      <div className="sm:mx-auto sm:w-full sm:max-w-md">
        {/* Brand Header */}
        <div className="flex flex-col items-center text-center">
          <div className="w-12 h-12 rounded-2xl bg-gradient-to-tr from-slate-900 via-blue-900 to-indigo-700 flex items-center justify-center text-white shadow-md ring-1 ring-blue-500/20 mb-3">
            <span className="font-mono text-emerald-400 font-extrabold text-xl mr-0.5">K</span>
            <span className="text-white text-xs font-bold">IQ</span>
          </div>
          <h1 className="text-2xl sm:text-3xl font-extrabold text-slate-900 tracking-tight">
            KadaiIQ Retail SaaS
          </h1>
          <p className="mt-1.5 text-xs sm:text-sm text-slate-500 max-w-sm">
            Multi-Tenant Footwear Retail Management & Advanced Sales Analytics
          </p>
        </div>

        {/* Card Form */}
        <div className="mt-8 bg-white py-8 px-6 sm:px-10 shadow-lg shadow-slate-200/50 rounded-2xl border border-slate-200/90">
          <form className="space-y-5" onSubmit={handleSubmit}>
            {error && (
              <div className="p-3 bg-rose-50 border border-rose-200 rounded-xl text-xs text-rose-700 flex items-center gap-2">
                <AlertCircle className="w-4 h-4 shrink-0 text-rose-600" />
                <span className="font-medium">{error}</span>
              </div>
            )}

            <div>
              <label className="block text-xs font-bold text-slate-700 uppercase tracking-wider mb-1.5" htmlFor="email">
                Work Email
              </label>
              <input
                id="email"
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="store.manager@retail.com"
                className="w-full px-3.5 py-2.5 bg-slate-50 border border-slate-300 rounded-xl text-xs sm:text-sm text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:bg-white transition-all"
              />
            </div>

            <div>
              <div className="flex items-center justify-between mb-1.5">
                <label className="block text-xs font-bold text-slate-700 uppercase tracking-wider" htmlFor="password">
                  Password
                </label>
                <Link to="/forgot-password" className="text-xs font-semibold text-blue-600 hover:text-blue-700">
                  Forgot password?
                </Link>
              </div>
              <input
                id="password"
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••••••"
                className="w-full px-3.5 py-2.5 bg-slate-50 border border-slate-300 rounded-xl text-xs sm:text-sm text-slate-900 placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:bg-white transition-all font-mono"
              />
            </div>

            <Button
              type="submit"
              disabled={loading}
              className="w-full py-2.5 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-xl shadow-md transition-all flex items-center justify-center gap-2 cursor-pointer"
            >
              {loading ? (
                <span>Authenticating...</span>
              ) : (
                <>
                  <span>Sign In to Store Workspace</span>
                  <ArrowRight className="w-4 h-4" />
                </>
              )}
            </Button>

            {/* Quick Demo Credentials Button */}
            <div className="pt-2 border-t border-slate-100">
              <button
                type="button"
                onClick={fillSuperAdmin}
                className="w-full py-2 px-3 bg-slate-100/80 hover:bg-slate-200 text-slate-700 text-xs font-semibold rounded-xl flex items-center justify-center gap-2 transition-colors cursor-pointer border border-slate-200/60"
              >
                <ShieldCheck className="w-4 h-4 text-purple-600" />
                <span>Auto-fill Super Admin Credentials</span>
              </button>
            </div>
          </form>

          {/* New Store Registration Link */}
          <div className="mt-6 text-center text-xs text-slate-500">
            <span>Looking to onboard your retail business? </span>
            <Link to="/register" className="font-bold text-blue-600 hover:underline">
              Create New Tenant
            </Link>
          </div>
        </div>

        {/* Feature Badges Footer */}
        <div className="mt-8 grid grid-cols-3 gap-2 text-center text-[11px] text-slate-500 font-medium">
          <div className="flex flex-col items-center gap-1 p-2 bg-white rounded-xl border border-slate-200/80 shadow-2xs">
            <TrendingUp className="w-4 h-4 text-emerald-600" />
            <span>Sales Analytics</span>
          </div>
          <div className="flex flex-col items-center gap-1 p-2 bg-white rounded-xl border border-slate-200/80 shadow-2xs">
            <Boxes className="w-4 h-4 text-blue-600" />
            <span>Footwear Matrix</span>
          </div>
          <div className="flex flex-col items-center gap-1 p-2 bg-white rounded-xl border border-slate-200/80 shadow-2xs">
            <Percent className="w-4 h-4 text-indigo-600" />
            <span>GST Compliance</span>
          </div>
        </div>
      </div>
    </div>
  )
}
