export interface ApiResponse<T = unknown> {
  success: boolean
  message: string
  data?: T
  error?: unknown
}

export interface AppHealthData {
  service: string
  status: string
  version: string
  timestamp: string
}

export interface DbHealthData {
  status: string
  latency_ms: number
  open_connections: number
  in_use: number
  idle: number
  timestamp: string
}

export interface UserProfile {
  id: string
  name: string
  email: string
}

export interface TenantMembershipInfo {
  tenant_id: string
  tenant_name: string
  tenant_slug: string
  role_id: string
  role_name: string
  status: string
  is_super_admin: boolean
}

export interface AuthResponseData {
  access_token: string
  refresh_token: string
  expires_in: number
  user: UserProfile
  active_tenant?: TenantMembershipInfo
  available_tenants: TenantMembershipInfo[]
}

export interface TenantData {
  id: string
  name: string
  slug: string
  business_name: string
  business_type: string
  gstin?: string
  email: string
  phone?: string
  address?: string
  status: string
  plan_name?: string
  member_count: number
  created_at: string
}

export interface MemberData {
  id: string
  user_id: string
  name: string
  email: string
  role_id: string
  role_name: string
  status: string
  created_at: string
}

// Product & Inventory Types
export interface ProductVariantData {
  id: string
  product_id: string
  size: string
  colour: string
  sku: string
  barcode: string
  purchase_price: string | number
  selling_price: string | number
  gst_rate: string | number
  current_stock: number
  min_stock_threshold: number
  stock_status: 'IN_STOCK' | 'LOW_STOCK' | 'OUT_OF_STOCK'
}

export interface ProductData {
  id: string
  name: string
  product_code: string
  sku?: string
  barcode?: string
  brand_name?: string
  category_name?: string
  brand_id?: string
  category_id?: string
  description?: string
  gender: string
  status: string
  total_stock: number
  variants: ProductVariantData[]
  created_at: string
}

export interface StockMovementData {
  id: string
  variant_id: string
  product_name: string
  variant_details: string
  barcode: string
  movement_type: string
  quantity: number
  previous_stock: number
  new_stock: number
  unit_cost: string | number
  reference_number?: string
  notes?: string
  created_by_name: string
  created_at: string
}

export interface InventorySummaryData {
  total_products: number
  total_variants: number
  total_stock_units: number
  low_stock_count: number
  out_of_stock_count: number
  stock_valuation: string | number
}

// Sales & POS Data Types
export interface SaleItemData {
  id: string
  product_name: string
  sku?: string
  barcode?: string
  variant_id?: string
  quantity: number
  unit_price: string | number
  discount: string | number
  tax_rate: string | number
  tax_amount: string | number
  line_total: string | number
}

export interface SaleData {
  id: string
  source_system: string
  source_bill_id: string
  bill_number: string
  sale_date: string
  subtotal: string | number
  discount: string | number
  tax_amount: string | number
  grand_total: string | number
  payment_method: string
  status: string
  inventory_deducted: boolean
  item_count: number
  items?: SaleItemData[]
  created_at: string
}

export interface SaleListData {
  items: SaleData[]
  total_count: number
  page: number
  page_size: number
}

export interface ImportPreviewItem {
  bill_number: string
  date: string
  barcode: string
  product_name: string
  quantity: number
  unit_price: string
  total: string
  matched: boolean
}

export interface ImportPreviewData {
  total_rows: number
  detected_columns: string[]
  sample_rows: ImportPreviewItem[]
}

export interface ImportJobData {
  id: string
  source_system: string
  file_name: string
  total_records: number
  imported_records: number
  failed_records: number
  duplicate_records: number
  status: string
  error_details?: string
  created_at: string
}

// Analytics Types
export interface AnalyticsOverviewData {
  date_range: string
  start_date: string
  end_date: string
  gross_sales: string | number
  total_discounts: string | number
  total_tax: string | number
  net_sales: string | number
  estimated_cost: string | number
  gross_profit: string | number
  profit_margin_pct: string | number
  total_bills: number
  total_items_sold: number
  average_order_value: string | number
  prev_gross_sales: string | number
  prev_net_sales: string | number
  prev_total_bills: number
  sales_growth_pct: string | number
  bills_growth_pct: string | number
}

export interface TrendPointData {
  period: string
  timestamp: string
  sales: string | number
  net_sales: string | number
  profit: string | number
  bills: number
  units_sold: number
}

export interface PaymentMethodData {
  method: string
  total_sales: string | number
  bill_count: number
  percentage: string | number
}

export interface ProductSalesData {
  product_id?: string
  product_name: string
  brand_name?: string
  category_name?: string
  units_sold: number
  revenue: string | number
  gross_profit: string | number
}

export interface CategorySalesData {
  category_name: string
  revenue: string | number
  units_sold: number
  percentage: string | number
}

export interface VariantAttributeData {
  attribute_type: string
  attribute_value: string
  units_sold: number
  revenue: string | number
  percentage: string | number
}

export interface FootwearIntelligenceData {
  top_sizes: VariantAttributeData[]
  top_colours: VariantAttributeData[]
  top_brands: VariantAttributeData[]
}

// GST Types
export interface GSTSettingsData {
  id: string
  tenant_id: string
  gstin: string
  legal_name: string
  trade_name: string
  state_code: string
  state_name: string
  registration_type: string
  default_footwear_hsn: string
  default_gst_rate: string | number
  enable_composition: boolean
  updated_at: string
}

export interface GSTR1RateSlabData {
  rate_pct: string | number
  taxable_value: string | number
  cgst: string | number
  sgst: string | number
  igst: string | number
  total_tax: string | number
  invoices: number
}

export interface GSTR1HSNRowData {
  hsn_code: string
  description: string
  uqc: string
  total_qty: number
  taxable_value: string | number
  cgst: string | number
  sgst: string | number
  igst: string | number
  total_tax: string | number
}

export interface GSTR1SummaryData {
  return_period: string
  start_date: string
  end_date: string
  gstin: string
  legal_name: string
  state_code: string
  total_taxable_value: string | number
  total_cgst: string | number
  total_sgst: string | number
  total_igst: string | number
  total_tax_liability: string | number
  b2c_small_count: number
  rate_breakdown: GSTR1RateSlabData[]
  hsn_summary: GSTR1HSNRowData[]
}

export interface GSTR3BTable31Data {
  nature_of_supplies: string
  taxable_value: string | number
  integrated_tax: string | number
  central_tax: string | number
  state_ut_tax: string | number
  cess: string | number
}

// Reports Types
export interface SalesSummaryRowData {
  bill_number: string
  sale_date: string
  payment_method: string
  item_count: number
  subtotal: string | number
  discount: string | number
  tax_amount: string | number
  grand_total: string | number
  status: string
}

export interface StockValuationRowData {
  product_name: string
  product_code: string
  size: string
  colour: string
  sku: string
  barcode: string
  current_stock: number
  purchase_price: string | number
  selling_price: string | number
  cost_valuation: string | number
  retail_valuation: string | number
  stock_status: string
}

export interface StockMovementReportRowData {
  date: string
  product_name: string
  barcode: string
  variant_details: string
  movement_type: string
  quantity: number
  previous_stock: number
  new_stock: number
  unit_cost: string | number
  reference_number?: string
  notes?: string
  created_by: string
}

export interface ReportContainerData<T> {
  report_name: string
  generated_at: string
  date_from?: string
  date_to?: string
  total_records: number
  summary_value?: string | number
  data: T[]
}

const BASE_URL = import.meta.env.VITE_API_URL || '/api/v1'

const getHeaders = (isFormData = false): Record<string, string> => {
  const headers: Record<string, string> = {}
  if (!isFormData) {
    headers['Content-Type'] = 'application/json'
  }
  const token = localStorage.getItem('kadai_token')
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  return headers
}

export const apiService = {
  // System Health
  async getHealth(): Promise<ApiResponse<AppHealthData>> {
    const res = await fetch(`${BASE_URL}/health`)
    if (!res.ok) throw new Error(`Health check failed: HTTP ${res.status}`)
    return res.json()
  },

  async getDbHealth(): Promise<ApiResponse<DbHealthData>> {
    const res = await fetch(`${BASE_URL}/health/db`)
    const data: ApiResponse<DbHealthData> = await res.json()
    if (!res.ok) throw new Error(data.message || `Database check failed: HTTP ${res.status}`)
    return data
  },

  // Auth Endpoints
  async login(payload: { email: string; password: string }): Promise<ApiResponse<AuthResponseData>> {
    const res = await fetch(`${BASE_URL}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Login failed')
    return data
  },

  async register(payload: {
    organization_name: string
    slug: string
    owner_name: string
    email: string
    password: string
    phone?: string
  }): Promise<ApiResponse<AuthResponseData>> {
    const res = await fetch(`${BASE_URL}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Registration failed')
    return data
  },

  async forgotPassword(email: string): Promise<ApiResponse<{ reset_token?: string }>> {
    const res = await fetch(`${BASE_URL}/auth/forgot-password`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email }),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Forgot password failed')
    return data
  },

  async resetPassword(payload: { token: string; new_password: string }): Promise<ApiResponse<null>> {
    const res = await fetch(`${BASE_URL}/auth/reset-password`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Password reset failed')
    return data
  },

  async switchTenant(tenant_id: string): Promise<ApiResponse<AuthResponseData>> {
    const res = await fetch(`${BASE_URL}/auth/switch-tenant`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify({ tenant_id }),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to switch organisation')
    return data
  },

  async logout(): Promise<void> {
    const refreshToken = localStorage.getItem('kadai_refresh_token')
    if (refreshToken) {
      await fetch(`${BASE_URL}/auth/logout`, {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify({ refresh_token: refreshToken }),
      }).catch(() => {})
    }
    localStorage.removeItem('kadai_token')
    localStorage.removeItem('kadai_refresh_token')
    localStorage.removeItem('kadai_user')
    localStorage.removeItem('kadai_active_tenant')
  },

  // Super Admin Endpoints
  async getTenants(): Promise<ApiResponse<TenantData[]>> {
    const res = await fetch(`${BASE_URL}/super-admin/tenants`, {
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch organisations')
    return data
  },

  async createTenant(payload: {
    name: string
    slug: string
    business_name: string
    business_type: string
    gstin?: string
    email: string
    phone?: string
    address?: string
    owner_name: string
    owner_email: string
    owner_password: string
  }): Promise<ApiResponse<TenantData>> {
    const res = await fetch(`${BASE_URL}/super-admin/tenants`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify(payload),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to create organisation')
    return data
  },

  async updateTenantStatus(id: string, status: string): Promise<ApiResponse<null>> {
    const res = await fetch(`${BASE_URL}/super-admin/tenants/${id}/status`, {
      method: 'PATCH',
      headers: getHeaders(),
      body: JSON.stringify({ status }),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to update status')
    return data
  },

  // Tenant Staff Endpoints
  async getMembers(): Promise<ApiResponse<MemberData[]>> {
    const res = await fetch(`${BASE_URL}/tenant/members`, {
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch members')
    return data
  },

  async inviteMember(payload: { name: string; email: string; role_id: string }): Promise<ApiResponse<MemberData>> {
    const res = await fetch(`${BASE_URL}/tenant/members`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify(payload),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to invite member')
    return data
  },

  // Products Endpoints
  async getProducts(params?: { search?: string; gender?: string; page?: number }): Promise<ApiResponse<{ items: ProductData[]; total_count: number }>> {
    const query = new URLSearchParams()
    if (params?.search) query.set('search', params.search)
    if (params?.gender) query.set('gender', params.gender)
    if (params?.page) query.set('page', params.page.toString())

    const res = await fetch(`${BASE_URL}/products?${query.toString()}`, {
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch products')
    return data
  },

  async createProduct(payload: {
    name: string
    product_code: string
    gender: string
    description?: string
    variants: Array<{
      size: string
      colour: string
      sku: string
      barcode: string
      purchase_price: number
      selling_price: number
      gst_rate: number
      opening_stock: number
      min_stock_threshold: number
    }>
  }): Promise<ApiResponse<ProductData>> {
    const res = await fetch(`${BASE_URL}/products`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify(payload),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to create product')
    return data
  },

  async archiveProduct(id: string): Promise<ApiResponse<null>> {
    const res = await fetch(`${BASE_URL}/products/${id}`, {
      method: 'DELETE',
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to archive product')
    return data
  },

  async lookupBarcode(barcode: string): Promise<ApiResponse<ProductVariantData>> {
    const res = await fetch(`${BASE_URL}/products/barcode/${encodeURIComponent(barcode)}`, {
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Barcode not found')
    return data
  },

  // Inventory Endpoints
  async getInventorySummary(): Promise<ApiResponse<InventorySummaryData>> {
    const res = await fetch(`${BASE_URL}/inventory/summary`, {
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch inventory summary')
    return data
  },

  async getStockMovements(params?: { variant_id?: string; page?: number }): Promise<ApiResponse<{ items: StockMovementData[]; total_count: number }>> {
    const query = new URLSearchParams()
    if (params?.variant_id) query.set('variant_id', params.variant_id)
    if (params?.page) query.set('page', params.page.toString())

    const res = await fetch(`${BASE_URL}/inventory/movements?${query.toString()}`, {
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch stock movements')
    return data
  },

  async recordStockAdjustment(payload: {
    variant_id: string
    movement_type: string
    quantity: number
    unit_cost?: number
    reference_number?: string
    notes?: string
  }): Promise<ApiResponse<StockMovementData>> {
    const res = await fetch(`${BASE_URL}/inventory/adjustments`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify(payload),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to adjust stock')
    return data
  },

  // Sales & Import Endpoints
  async getSales(params?: {
    search?: string
    status?: string
    from?: string
    to?: string
    page?: number
    page_size?: number
  }): Promise<ApiResponse<SaleListData>> {
    const query = new URLSearchParams()
    if (params?.search) query.set('search', params.search)
    if (params?.status) query.set('status', params.status)
    if (params?.from) query.set('from', params.from)
    if (params?.to) query.set('to', params.to)
    if (params?.page) query.set('page', params.page.toString())
    if (params?.page_size) query.set('page_size', params.page_size.toString())

    const res = await fetch(`${BASE_URL}/sales?${query.toString()}`, {
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch sales')
    return data
  },

  async getSale(id: string): Promise<ApiResponse<SaleData>> {
    const res = await fetch(`${BASE_URL}/sales/${id}`, {
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch sale details')
    return data
  },

  async previewSalesImport(file: File): Promise<ApiResponse<ImportPreviewData>> {
    const formData = new FormData()
    formData.append('file', file)

    const res = await fetch(`${BASE_URL}/sales/preview`, {
      method: 'POST',
      headers: getHeaders(true),
      body: formData,
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to preview sales file')
    return data
  },

  async executeSalesImport(params: {
    file: File
    source_system?: string
    deduct_inventory?: boolean
  }): Promise<ApiResponse<ImportJobData>> {
    const formData = new FormData()
    formData.append('file', params.file)
    formData.append('source_system', params.source_system || 'ZenBiller')
    formData.append('deduct_inventory', params.deduct_inventory ? 'true' : 'false')

    const res = await fetch(`${BASE_URL}/sales/import`, {
      method: 'POST',
      headers: getHeaders(true),
      body: formData,
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to import sales data')
    return data
  },

  async getImportHistory(): Promise<ApiResponse<ImportJobData[]>> {
    const res = await fetch(`${BASE_URL}/sales/import/history`, {
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch import history')
    return data
  },

  // Analytics Endpoints
  async getAnalyticsOverview(params?: { range?: string; from?: string; to?: string }): Promise<ApiResponse<AnalyticsOverviewData>> {
    const query = new URLSearchParams()
    if (params?.range) query.set('range', params.range)
    if (params?.from) query.set('from', params.from)
    if (params?.to) query.set('to', params.to)

    const res = await fetch(`${BASE_URL}/analytics/overview?${query.toString()}`, {
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch analytics overview')
    return data
  },

  async getSalesTrends(params?: { range?: string; from?: string; to?: string }): Promise<ApiResponse<TrendPointData[]>> {
    const query = new URLSearchParams()
    if (params?.range) query.set('range', params.range)
    if (params?.from) query.set('from', params.from)
    if (params?.to) query.set('to', params.to)

    const res = await fetch(`${BASE_URL}/analytics/trends?${query.toString()}`, {
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch sales trends')
    return data
  },

  async getPaymentMethods(params?: { range?: string; from?: string; to?: string }): Promise<ApiResponse<PaymentMethodData[]>> {
    const query = new URLSearchParams()
    if (params?.range) query.set('range', params.range)
    if (params?.from) query.set('from', params.from)
    if (params?.to) query.set('to', params.to)

    const res = await fetch(`${BASE_URL}/analytics/payment-methods?${query.toString()}`, {
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch payment methods')
    return data
  },

  async getTopProducts(params?: { range?: string; from?: string; to?: string; limit?: number }): Promise<ApiResponse<ProductSalesData[]>> {
    const query = new URLSearchParams()
    if (params?.range) query.set('range', params.range)
    if (params?.from) query.set('from', params.from)
    if (params?.to) query.set('to', params.to)
    if (params?.limit) query.set('limit', params.limit.toString())

    const res = await fetch(`${BASE_URL}/analytics/products?${query.toString()}`, {
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch top products')
    return data
  },

  async getCategorySales(params?: { range?: string; from?: string; to?: string }): Promise<ApiResponse<CategorySalesData[]>> {
    const query = new URLSearchParams()
    if (params?.range) query.set('range', params.range)
    if (params?.from) query.set('from', params.from)
    if (params?.to) query.set('to', params.to)

    const res = await fetch(`${BASE_URL}/analytics/categories?${query.toString()}`, {
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch category sales')
    return data
  },

  async getFootwearIntelligence(params?: { range?: string; from?: string; to?: string }): Promise<ApiResponse<FootwearIntelligenceData>> {
    const query = new URLSearchParams()
    if (params?.range) query.set('range', params.range)
    if (params?.from) query.set('from', params.from)
    if (params?.to) query.set('to', params.to)

    const res = await fetch(`${BASE_URL}/analytics/footwear?${query.toString()}`, {
      headers: getHeaders(),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch footwear intelligence')
    return data
  },

  // GST Endpoints
  async getGSTSettings(): Promise<ApiResponse<GSTSettingsData>> {
    const res = await fetch(`${BASE_URL}/gst/settings`, { headers: getHeaders() })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch GST settings')
    return data
  },

  async updateGSTSettings(payload: Partial<GSTSettingsData>): Promise<ApiResponse<GSTSettingsData>> {
    const res = await fetch(`${BASE_URL}/gst/settings`, {
      method: 'PUT',
      headers: getHeaders(),
      body: JSON.stringify(payload),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to update GST settings')
    return data
  },

  async getGSTR1(params?: { month?: string; year?: string }): Promise<ApiResponse<GSTR1SummaryData>> {
    const query = new URLSearchParams()
    if (params?.month) query.set('month', params.month)
    if (params?.year) query.set('year', params.year)

    const res = await fetch(`${BASE_URL}/gst/gstr1?${query.toString()}`, { headers: getHeaders() })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch GSTR-1 summary')
    return data
  },

  async getGSTR3B(params?: { month?: string; year?: string }): Promise<ApiResponse<GSTR3BTable31Data>> {
    const query = new URLSearchParams()
    if (params?.month) query.set('month', params.month)
    if (params?.year) query.set('year', params.year)

    const res = await fetch(`${BASE_URL}/gst/gstr3b?${query.toString()}`, { headers: getHeaders() })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch GSTR-3B table')
    return data
  },

  // Reports Endpoints
  async getSalesReport(params?: { from?: string; to?: string }): Promise<ApiResponse<ReportContainerData<SalesSummaryRowData>>> {
    const query = new URLSearchParams()
    if (params?.from) query.set('from', params.from)
    if (params?.to) query.set('to', params.to)

    const res = await fetch(`${BASE_URL}/reports/sales?${query.toString()}`, { headers: getHeaders() })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch sales report')
    return data
  },

  async getInventoryReport(): Promise<ApiResponse<ReportContainerData<StockValuationRowData>>> {
    const res = await fetch(`${BASE_URL}/reports/inventory`, { headers: getHeaders() })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch inventory report')
    return data
  },

  async getMovementsReport(params?: { from?: string; to?: string }): Promise<ApiResponse<ReportContainerData<StockMovementReportRowData>>> {
    const query = new URLSearchParams()
    if (params?.from) query.set('from', params.from)
    if (params?.to) query.set('to', params.to)

    const res = await fetch(`${BASE_URL}/reports/movements?${query.toString()}`, { headers: getHeaders() })
    const data = await res.json()
    if (!res.ok) throw new Error(data.message || 'Failed to fetch movements report')
    return data
  },

  getReportDownloadUrl(type: string, from?: string, to?: string): string {
    const query = new URLSearchParams()
    query.set('type', type)
    if (from) query.set('from', from)
    if (to) query.set('to', to)
    return `${BASE_URL}/reports/download?${query.toString()}`
  },
}
