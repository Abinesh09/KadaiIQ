import React from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { AppLayout } from '@/components/layout/AppLayout'
import { DashboardPage } from '@/features/dashboard/DashboardPage'
import { ProductsPage } from '@/features/products/ProductsPage'
import { InventoryPage } from '@/features/inventory/InventoryPage'
import { SalesPage } from '@/features/sales/SalesPage'
import { AnalyticsPage } from '@/features/analytics/AnalyticsPage'
import { ReportsPage } from '@/features/reports/ReportsPage'
import { GSTPage } from '@/features/gst/GSTPage'
import { SettingsPage } from '@/features/settings/SettingsPage'
import { SuperAdminPage } from '@/features/tenants/SuperAdminPage'

export const AppRoutes: React.FC = () => {
  return (
    <Routes>
      <Route path="/" element={<AppLayout />}>
        <Route index element={<DashboardPage />} />
        <Route path="products" element={<ProductsPage />} />
        <Route path="inventory" element={<InventoryPage />} />
        <Route path="sales" element={<SalesPage />} />
        <Route path="analytics" element={<AnalyticsPage />} />
        <Route path="reports" element={<ReportsPage />} />
        <Route path="gst" element={<GSTPage />} />
        <Route path="settings" element={<SettingsPage />} />
        <Route path="super-admin" element={<SuperAdminPage />} />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

